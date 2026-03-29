package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itvico/e-proc-api/internal/config"
	"github.com/itvico/e-proc-api/internal/middleware"
	"github.com/itvico/e-proc-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserDTO   `json:"user"`
}

type UserDTO struct {
	ID             uint   `json:"id"`
	VendorID       uint   `json:"vendor_id,omitempty"`
	EntityID       uint   `json:"entity_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	RoleCode       string `json:"role_code"`
	RoleName       string `json:"role_name"`
	ScopeType      string `json:"scope_type"`
	SubjectType    string `json:"subject_type"`
	PortalType     string `json:"portal_type"`
	ForceChange    bool   `json:"force_change_password"`
	DepartmentName string `json:"department_name,omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	var user models.User
	err := s.db.
		Preload("Department").
		Preload("PrimaryRole").
		Preload("UserRoles", "is_primary = ? AND status = ?", true, "active").
		Preload("UserRoles.Role").
		Where("username = ? AND status = ?", req.Username, "active").
		First(&user).Error
	if err == nil {
		return s.loginInternalUser(&user, req.Password)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var vendorUser models.VendorUser
	if err := s.db.Where("email = ? AND status = ?", req.Username, "active").First(&vendorUser).Error; err != nil {
		recordAuditLog(s.db, AuditEntry{
			ModuleCode: "AUTH",
			ObjectType: "AUTH",
			Action:     "login_failed",
			EventCode:  "AUTH_LOGIN_FAILED",
			ActorType:  "anonymous",
			DataAfter:  fmt.Sprintf(`{"username":"%s"}`, req.Username),
		})
		return nil, errors.New("invalid credentials")
	}

	return s.loginVendorUser(&vendorUser, req.Password)
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *AuthService) ChangePassword(userID uint, req ChangePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is invalid")
	}
	if err := validatePasswordPolicy(req.NewPassword); err != nil {
		return err
	}

	newHash, err := s.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := s.db.Model(&user).Updates(map[string]interface{}{
		"password_hash":         newHash,
		"force_change_password": false,
		"failed_login_count":    0,
		"locked_until":          nil,
	}).Error; err != nil {
		return err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(user.EntityID),
		ModuleCode: "AUTH",
		ObjectType: "USER",
		ObjectID:   user.ID,
		EventCode:  "AUTH_CHANGE_PASSWORD",
		ActorType:  "internal_user",
		ActorID:    uintPtr(user.ID),
		Action:     "change_password",
	})
	return nil
}

func (s *AuthService) loginInternalUser(user *models.User, password string) (*LoginResponse, error) {
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, errors.New("account temporarily locked")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		failedCount := user.FailedLoginCount + 1
		updates := map[string]interface{}{"failed_login_count": failedCount}
		if failedCount >= 5 {
			updates["locked_until"] = time.Now().Add(30 * time.Minute)
		}
		_ = s.db.Model(user).Updates(updates).Error
		recordAuditLog(s.db, AuditEntry{
			EntityID:   uintPtr(user.EntityID),
			ModuleCode: "AUTH",
			ObjectType: "USER",
			ObjectID:   user.ID,
			EventCode:  "AUTH_LOGIN_FAILED",
			ActorType:  "internal_user",
			ActorID:    uintPtr(user.ID),
			Action:     "login_failed",
		})
		return nil, errors.New("invalid credentials")
	}

	roleCode := ""
	roleName := ""
	scopeType := "own_entity"

	if len(user.UserRoles) > 0 {
		roleCode = user.UserRoles[0].Role.RoleCode
		roleName = user.UserRoles[0].Role.RoleName
		scopeType = user.UserRoles[0].ScopeType
	} else if user.PrimaryRole != nil {
		roleCode = user.PrimaryRole.RoleCode
		roleName = user.PrimaryRole.RoleName
	}

	departmentName := ""
	if user.Department != nil {
		departmentName = user.Department.Name
	}

	tokenStr, refreshToken, expiresAt, err := s.issueTokens(middleware.Claims{
		UserID:      user.ID,
		EntityID:    user.EntityID,
		Username:    user.Username,
		RoleCode:    roleCode,
		RoleName:    roleName,
		ScopeType:   scopeType,
		SubjectType: "internal_user",
		PortalType:  "internal",
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.db.Model(user).Updates(map[string]interface{}{
		"last_login_at":      &now,
		"failed_login_count": 0,
		"locked_until":       nil,
	}).Error

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(user.EntityID),
		ModuleCode: "AUTH",
		ObjectType: "USER",
		ObjectID:   user.ID,
		EventCode:  "AUTH_LOGIN_SUCCESS",
		ActorType:  "internal_user",
		ActorID:    uintPtr(user.ID),
		Action:     "login",
	})

	return &LoginResponse{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserDTO{
			ID:             user.ID,
			EntityID:       user.EntityID,
			Username:       user.Username,
			Email:          user.Email,
			FullName:       user.FullName,
			RoleCode:       roleCode,
			RoleName:       roleName,
			ScopeType:      scopeType,
			SubjectType:    "internal_user",
			PortalType:     "internal",
			ForceChange:    user.ForceChangePassword,
			DepartmentName: departmentName,
		},
	}, nil
}

func (s *AuthService) loginVendorUser(user *models.VendorUser, password string) (*LoginResponse, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		recordAuditLog(s.db, AuditEntry{
			ModuleCode: "AUTH",
			ObjectType: "VENDOR_USER",
			ObjectID:   user.ID,
			EventCode:  "AUTH_LOGIN_FAILED",
			ActorType:  "vendor_user",
			ActorID:    uintPtr(user.ID),
			Action:     "login_failed",
		})
		return nil, errors.New("invalid credentials")
	}

	tokenStr, refreshToken, expiresAt, err := s.issueTokens(middleware.Claims{
		UserID:      user.ID,
		VendorID:    user.VendorID,
		Username:    user.Email,
		RoleCode:    "VENDOR_ADMIN",
		RoleName:    "Vendor Admin",
		ScopeType:   "vendor_self",
		SubjectType: "vendor_user",
		PortalType:  "vendor",
	})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_ = s.db.Model(user).Updates(map[string]interface{}{
		"last_login_at": &now,
	}).Error

	recordAuditLog(s.db, AuditEntry{
		ModuleCode: "AUTH",
		ObjectType: "VENDOR_USER",
		ObjectID:   user.ID,
		EventCode:  "AUTH_LOGIN_SUCCESS",
		ActorType:  "vendor_user",
		ActorID:    uintPtr(user.ID),
		Action:     "login",
	})

	return &LoginResponse{
		Token:        tokenStr,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User: UserDTO{
			ID:          user.ID,
			VendorID:    user.VendorID,
			Username:    user.Email,
			Email:       user.Email,
			FullName:    user.FullName,
			RoleCode:    "VENDOR_ADMIN",
			RoleName:    "Vendor Admin",
			ScopeType:   "vendor_self",
			SubjectType: "vendor_user",
			PortalType:  "vendor",
			ForceChange: user.ForceChangePassword,
		},
	}, nil
}

func (s *AuthService) issueTokens(baseClaims middleware.Claims) (string, string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(s.cfg.JWT.ExpiryHours) * time.Hour)
	baseClaims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &baseClaims)
	tokenStr, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshClaims := baseClaims
	refreshClaims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(s.cfg.JWT.RefreshExpHours) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", time.Time{}, err
	}

	return tokenStr, refreshTokenStr, expiresAt, nil
}

func validatePasswordPolicy(password string) error {
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	switch {
	case len(password) < 8:
		return errors.New("password must be at least 8 characters")
	case !hasUpper:
		return errors.New("password must contain at least one uppercase letter")
	case !hasLower:
		return errors.New("password must contain at least one lowercase letter")
	case !hasDigit:
		return errors.New("password must contain at least one number")
	case !hasSpecial:
		return errors.New("password must contain at least one special character")
	default:
		return nil
	}
}
