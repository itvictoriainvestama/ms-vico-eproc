package services

import (
	"errors"
	"time"

	"github.com/itvico/e-proc-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type UserListResult struct {
	Items []models.User `json:"items"`
	Total int64         `json:"total"`
}

type UserListParams struct {
	EntityID uint
	Status   string
}

type CreateUserRequest struct {
	EntityID            uint   `json:"entity_id" binding:"required"`
	DepartmentID        *uint  `json:"department_id"`
	FullName            string `json:"full_name" binding:"required"`
	Email               string `json:"email" binding:"required,email"`
	Username            string `json:"username" binding:"required"`
	Password            string `json:"password" binding:"required,min=8"`
	RoleCode            string `json:"role_code" binding:"required"`
	ScopeType           string `json:"scope_type"`
	Status              string `json:"status"`
	ForceChangePassword bool   `json:"force_change_password"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

type CreateDelegateApproverRequest struct {
	OriginalUserID uint   `json:"original_user_id" binding:"required"`
	DelegateUserID uint   `json:"delegate_user_id" binding:"required"`
	StartAt        string `json:"start_at" binding:"required"`
	EndAt          string `json:"end_at" binding:"required"`
	Reason         string `json:"reason"`
}

func (s *UserService) List(actorEntityID uint, scopeType string, params UserListParams) (*UserListResult, error) {
	query := s.db.Model(&models.User{}).
		Preload("Department").
		Preload("PrimaryRole").
		Preload("UserRoles", "is_primary = ? AND status = ?", true, "active").
		Preload("UserRoles.Role")

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	effectiveEntityID := params.EntityID
	if scopeType != ScopeCrossEntity {
		effectiveEntityID = actorEntityID
	}
	if effectiveEntityID > 0 {
		query = query.Where("entity_id = ?", effectiveEntityID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []models.User
	if err := query.Order("full_name ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return &UserListResult{Items: items, Total: total}, nil
}

func (s *UserService) GetByID(id, actorEntityID uint, scopeType string) (*models.User, error) {
	var user models.User
	query := s.db.
		Preload("Department").
		Preload("PrimaryRole").
		Preload("UserRoles").
		Preload("UserRoles.Role")
	query = applyEntityScope(query, "entity_id", actorEntityID, scopeType)

	if err := query.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) Create(req CreateUserRequest, actorEntityID uint, actorScopeType string) (*models.User, error) {
	if actorScopeType != ScopeCrossEntity && req.EntityID != actorEntityID {
		return nil, errors.New("cannot create user outside your entity scope")
	}

	var entity models.Entity
	if err := s.db.First(&entity, req.EntityID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entity not found")
		}
		return nil, err
	}

	if req.DepartmentID != nil {
		var department models.Department
		if err := s.db.First(&department, *req.DepartmentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("department not found")
			}
			return nil, err
		}
		if department.EntityID != req.EntityID {
			return nil, errors.New("department is outside the selected entity")
		}
	}

	var role models.Role
	if err := s.db.Where("role_code = ? AND is_active = ?", req.RoleCode, true).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	roleScopeType := req.ScopeType
	if roleScopeType == "" {
		roleScopeType = "own_entity"
	}
	if actorScopeType != ScopeCrossEntity && roleScopeType == ScopeCrossEntity {
		return nil, errors.New("cannot assign cross-entity scope from your current role")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		EntityID:            req.EntityID,
		DepartmentID:        req.DepartmentID,
		FullName:            req.FullName,
		Email:               req.Email,
		Username:            req.Username,
		PasswordHash:        string(passwordHash),
		Status:              status,
		ForceChangePassword: req.ForceChangePassword,
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		userRole := models.UserRole{
			UserID:    user.ID,
			RoleID:    role.ID,
			EntityID:  req.EntityID,
			ScopeType: roleScopeType,
			IsPrimary: true,
			Status:    "active",
		}

		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}

		return tx.Model(user).Update("primary_role_id", role.ID).Error
	}); err != nil {
		return nil, err
	}

	return s.GetByID(user.ID, actorEntityID, actorScopeType)
}

func (s *UserService) ResetPassword(id, actorEntityID uint, actorScopeType string, req ResetPasswordRequest) error {
	user, err := s.GetByID(id, actorEntityID, actorScopeType)
	if err != nil {
		return err
	}

	newPassword := req.NewPassword
	if newPassword == "" {
		newPassword = "Temp123!"
	}
	if err := validatePasswordPolicy(newPassword); err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.db.Model(user).Updates(map[string]interface{}{
		"password_hash":         string(passwordHash),
		"force_change_password": true,
		"failed_login_count":    0,
		"locked_until":          nil,
	}).Error; err != nil {
		return err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(user.EntityID),
		ModuleCode: "USER_MANAGEMENT",
		ObjectType: "USER",
		ObjectID:   user.ID,
		EventCode:  "USER_RESET_PASSWORD",
		ActorType:  "internal_user",
		ActorID:    uintPtr(actorEntityID),
		Action:     "reset_password",
	})
	return nil
}

func (s *UserService) ListDelegates(actorEntityID uint, scopeType string) ([]models.DelegateApprover, error) {
	query := s.db.Model(&models.DelegateApprover{})
	query = applyEntityScope(query, "entity_id", actorEntityID, scopeType)

	var delegates []models.DelegateApprover
	if err := query.Order("start_at DESC").Find(&delegates).Error; err != nil {
		return nil, err
	}
	return delegates, nil
}

func (s *UserService) CreateDelegate(req CreateDelegateApproverRequest, actorEntityID uint, scopeType string) (*models.DelegateApprover, error) {
	startAt, err := time.Parse(time.RFC3339, req.StartAt)
	if err != nil {
		return nil, errors.New("start_at must use RFC3339 format")
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		return nil, errors.New("end_at must use RFC3339 format")
	}
	if !endAt.After(startAt) {
		return nil, errors.New("end_at must be after start_at")
	}

	originalUser, err := s.GetByID(req.OriginalUserID, actorEntityID, scopeType)
	if err != nil {
		return nil, err
	}
	delegateUser, err := s.GetByID(req.DelegateUserID, actorEntityID, scopeType)
	if err != nil {
		return nil, err
	}
	if originalUser.EntityID != delegateUser.EntityID {
		return nil, errors.New("delegate approver must be in the same entity")
	}

	delegate := &models.DelegateApprover{
		EntityID:       originalUser.EntityID,
		OriginalUserID: originalUser.ID,
		DelegateUserID: delegateUser.ID,
		StartAt:        startAt,
		EndAt:          endAt,
		Reason:         req.Reason,
		Status:         "active",
	}
	if err := s.db.Create(delegate).Error; err != nil {
		return nil, err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(delegate.EntityID),
		ModuleCode: "USER_MANAGEMENT",
		ObjectType: "DELEGATE_APPROVER",
		ObjectID:   delegate.ID,
		EventCode:  "USER_DELEGATE_CREATE",
		ActorType:  "internal_user",
		ActorID:    uintPtr(actorEntityID),
		Action:     "create_delegate",
	})
	return delegate, nil
}
