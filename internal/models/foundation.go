package models

import "time"

type VendorUser struct {
	BaseModel
	VendorID            uint       `gorm:"not null;index" json:"vendor_id"`
	FullName            string     `gorm:"not null;size:150" json:"full_name"`
	Email               string     `gorm:"uniqueIndex;not null;size:150" json:"email"`
	PasswordHash        string     `gorm:"not null" json:"-"`
	Status              string     `gorm:"not null;size:20;default:'active'" json:"status"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
	ForceChangePassword bool       `gorm:"not null;default:false" json:"force_change_password"`
}

type VendorBlacklist struct {
	BaseModel
	VendorID      uint       `gorm:"not null;index" json:"vendor_id"`
	EntityID      *uint      `gorm:"index" json:"entity_id,omitempty"`
	BlacklistType string     `gorm:"not null;size:20;default:'group'" json:"blacklist_type"`
	Reason        string     `gorm:"type:text" json:"reason"`
	StartAt       time.Time  `gorm:"not null" json:"start_at"`
	EndAt         *time.Time `json:"end_at,omitempty"`
	Status        string     `gorm:"not null;size:20;default:'active'" json:"status"`
}

type ReferencePrice struct {
	BaseModel
	EntityID        uint      `gorm:"not null;index" json:"entity_id"`
	ItemName        string    `gorm:"not null;size:200" json:"item_name"`
	ItemCategory    string    `gorm:"size:100" json:"item_category,omitempty"`
	UOM             string    `gorm:"not null;size:50" json:"uom"`
	ReferencePrice  float64   `gorm:"type:decimal(18,2);not null" json:"reference_price"`
	SourceType      string    `gorm:"not null;size:30" json:"source_type"`
	SourceReference string    `gorm:"size:255" json:"source_reference,omitempty"`
	EffectiveDate   time.Time `gorm:"not null" json:"effective_date"`
	Status          string    `gorm:"not null;size:20;default:'active'" json:"status"`
}

type Budget struct {
	BaseModel
	EntityID            uint    `gorm:"not null;index" json:"entity_id"`
	FiscalYear          int     `gorm:"not null;index" json:"fiscal_year"`
	DepartmentCode      string  `gorm:"size:50;index" json:"department_code,omitempty"`
	ProcurementCategory string  `gorm:"size:50;index" json:"procurement_category,omitempty"`
	BudgetMode          string  `gorm:"not null;size:20;default:'limited'" json:"budget_mode"`
	AllocatedAmount     float64 `gorm:"type:decimal(18,2);default:0" json:"allocated_amount"`
	ReservedAmount      float64 `gorm:"type:decimal(18,2);default:0" json:"reserved_amount"`
	UsedAmount          float64 `gorm:"type:decimal(18,2);default:0" json:"used_amount"`
	Status              string  `gorm:"not null;size:20;default:'active'" json:"status"`
}

type ProcurementPolicyRule struct {
	BaseModel
	EntityID          *uint   `gorm:"index" json:"entity_id,omitempty"`
	PolicyCode        string  `gorm:"not null;size:50;index" json:"policy_code"`
	GoodsOrServices   string  `gorm:"not null;size:20" json:"goods_or_services"`
	RoutineType       string  `gorm:"not null;size:20" json:"routine_type"`
	BudgetStatus      string  `gorm:"not null;size:20" json:"budget_status"`
	MinAmount         float64 `gorm:"type:decimal(18,2);default:0" json:"min_amount"`
	MaxAmount         float64 `gorm:"type:decimal(18,2);default:0" json:"max_amount"`
	RecommendedMethod string  `gorm:"not null;size:30" json:"recommended_method"`
	ApprovalModelID   *uint   `gorm:"index" json:"approval_model_id,omitempty"`
	Status            string  `gorm:"not null;size:20;default:'active'" json:"status"`
}

type ApprovalModel struct {
	BaseModel
	EntityID    *uint  `gorm:"index" json:"entity_id,omitempty"`
	ModelCode   string `gorm:"not null;size:50;index" json:"model_code"`
	ModelName   string `gorm:"not null;size:100" json:"model_name"`
	ObjectType  string `gorm:"not null;size:20" json:"object_type"`
	Status      string `gorm:"not null;size:20;default:'active'" json:"status"`
	Description string `gorm:"size:255" json:"description,omitempty"`
}

type ApprovalMatrix struct {
	BaseModel
	ApprovalModelID  uint   `gorm:"not null;index" json:"approval_model_id"`
	LevelNo          int    `gorm:"not null" json:"level_no"`
	ApproverRoleCode string `gorm:"not null;size:50" json:"approver_role_code"`
	ApproverScope    string `gorm:"not null;size:20" json:"approver_scope"`
	LogicType        string `gorm:"not null;size:20;default:'single'" json:"logic_type"`
	SLADays          int    `gorm:"not null;default:2" json:"sla_days"`
	IsMandatory      bool   `gorm:"not null;default:true" json:"is_mandatory"`
	Status           string `gorm:"not null;size:20;default:'active'" json:"status"`
}

type Notification struct {
	BaseModel
	ChannelType   string     `gorm:"not null;size:20" json:"channel_type"`
	RecipientType string     `gorm:"not null;size:20" json:"recipient_type"`
	RecipientID   uint       `gorm:"not null;index" json:"recipient_id"`
	EntityID      *uint      `gorm:"index" json:"entity_id,omitempty"`
	EventCode     string     `gorm:"not null;size:50;index" json:"event_code"`
	Subject       string     `gorm:"size:255" json:"subject,omitempty"`
	Body          string     `gorm:"type:text" json:"body,omitempty"`
	Status        string     `gorm:"not null;size:20;default:'queued'" json:"status"`
	SentAt        *time.Time `json:"sent_at,omitempty"`
	ErrorMessage  string     `gorm:"type:text" json:"error_message,omitempty"`
}

type AppLog struct {
	BaseModel
	ServiceName      string `gorm:"not null;size:50;index" json:"service_name"`
	Level            string `gorm:"not null;size:20" json:"level"`
	TraceID          string `gorm:"size:100;index" json:"trace_id,omitempty"`
	Message          string `gorm:"not null;type:text" json:"message"`
	Payload          string `gorm:"type:longtext" json:"payload,omitempty"`
	ShippedToElastic bool   `gorm:"not null;default:false" json:"shipped_to_elastic"`
}
