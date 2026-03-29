package services

import (
	"errors"

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
