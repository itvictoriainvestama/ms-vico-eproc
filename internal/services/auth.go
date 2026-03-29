package services

import (
	"errors"
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
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserDTO   `json:"user"`
}

type UserDTO struct {
	ID             uint   `json:"id"`
	EntityID       uint   `json:"entity_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FullName       string `json:"full_name"`
	RoleCode       string `json:"role_code"`
	RoleName       string `json:"role_name"`
	ScopeType      string `json:"scope_type"`
	DepartmentName string `json:"department_name,omitempty"`
}

func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	var user models.User
	err := s.db.
		Preload("Department").
		Preload("PrimaryRole").
		Preload("UserRoles", "is_primary = ? AND status = ?", true, "active").
		Preload("UserRoles.Role").
		Where("username = ? AND status = ?", req.Username, "active").
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		return nil, errors.New("account temporarily locked")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
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

	expiresAt := time.Now().Add(time.Duration(s.cfg.JWT.ExpiryHours) * time.Hour)
	claims := &middleware.Claims{
		UserID:      user.ID,
		EntityID:    user.EntityID,
		Username:    user.Username,
		RoleCode:    roleCode,
		RoleName:    roleName,
		ScopeType:   scopeType,
		SubjectType: "internal_user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	now := time.Now()
	s.db.Model(&user).Updates(map[string]interface{}{
		"last_login_at":      &now,
		"failed_login_count": 0,
	})

	return &LoginResponse{
		Token:     tokenStr,
		ExpiresAt: expiresAt,
		User: UserDTO{
			ID:             user.ID,
			EntityID:       user.EntityID,
			Username:       user.Username,
			Email:          user.Email,
			FullName:       user.FullName,
			RoleCode:       roleCode,
			RoleName:       roleName,
			ScopeType:      scopeType,
			DepartmentName: departmentName,
		},
	}, nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}
