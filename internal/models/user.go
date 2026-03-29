package models

import "time"

type Entity struct {
	BaseModel
	EntityCode        string  `gorm:"uniqueIndex;not null;size:50" json:"entity_code"`
	EntityName        string  `gorm:"not null;size:200" json:"entity_name"`
	ParentEntityID    *uint   `gorm:"index" json:"parent_entity_id,omitempty"`
	ParentEntity      *Entity `gorm:"foreignKey:ParentEntityID" json:"parent_entity,omitempty"`
	EntityType        string  `gorm:"not null;size:30;default:'subsidiary'" json:"entity_type"`
	Status            string  `gorm:"not null;size:20;default:'active'" json:"status"`
	ApprovalModelCode string  `gorm:"size:50" json:"approval_model_code,omitempty"`
	GovernanceMode    string  `gorm:"not null;size:30;default:'entity_only'" json:"governance_mode"`
}

type Department struct {
	BaseModel
	EntityID uint   `gorm:"not null;index" json:"entity_id"`
	Name     string `gorm:"not null;size:100" json:"name"`
	Code     string `gorm:"not null;size:50" json:"code"`
}

type Role struct {
	BaseModel
	RoleCode    string `gorm:"uniqueIndex;not null;size:50" json:"role_code"`
	RoleName    string `gorm:"not null;size:100" json:"role_name"`
	PortalType  string `gorm:"not null;size:20;default:'internal'" json:"portal_type"`
	IsActive    bool   `gorm:"not null;default:true" json:"is_active"`
	Description string `gorm:"size:255" json:"description,omitempty"`
}

type User struct {
	BaseModel
	EntityID            uint        `gorm:"not null;index" json:"entity_id"`
	Entity              Entity      `gorm:"foreignKey:EntityID" json:"entity"`
	DepartmentID        *uint       `gorm:"index" json:"department_id,omitempty"`
	Department          *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	FullName            string      `gorm:"not null;size:150" json:"full_name"`
	Email               string      `gorm:"uniqueIndex;not null;size:150" json:"email"`
	Username            string      `gorm:"uniqueIndex;not null;size:100" json:"username"`
	PasswordHash        string      `gorm:"not null" json:"-"`
	Status              string      `gorm:"not null;size:20;default:'active'" json:"status"`
	ForceChangePassword bool        `gorm:"not null;default:false" json:"force_change_password"`
	LastLoginAt         *time.Time  `json:"last_login_at,omitempty"`
	FailedLoginCount    int         `gorm:"not null;default:0" json:"failed_login_count"`
	LockedUntil         *time.Time  `json:"locked_until,omitempty"`
	UserRoles           []UserRole  `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
	PrimaryRoleID       *uint       `gorm:"index" json:"primary_role_id,omitempty"`
	PrimaryRole         *Role       `gorm:"foreignKey:PrimaryRoleID" json:"primary_role,omitempty"`
}

type UserRole struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"-"`
	RoleID    uint   `gorm:"not null;index" json:"role_id"`
	Role      Role   `gorm:"foreignKey:RoleID" json:"role"`
	EntityID  uint   `gorm:"not null;index" json:"entity_id"`
	Entity    Entity `gorm:"foreignKey:EntityID" json:"entity"`
	ScopeType string `gorm:"not null;size:30;default:'own_entity'" json:"scope_type"`
	IsPrimary bool   `gorm:"not null;default:false" json:"is_primary"`
	Status    string `gorm:"not null;size:20;default:'active'" json:"status"`
}

type DelegateApprover struct {
	BaseModel
	EntityID       uint      `gorm:"not null;index" json:"entity_id"`
	OriginalUserID uint      `gorm:"not null;index" json:"original_user_id"`
	DelegateUserID uint      `gorm:"not null;index" json:"delegate_user_id"`
	StartAt        time.Time `gorm:"not null" json:"start_at"`
	EndAt          time.Time `gorm:"not null" json:"end_at"`
	Reason         string    `gorm:"type:text" json:"reason"`
	Status         string    `gorm:"not null;size:20;default:'active'" json:"status"`
}
