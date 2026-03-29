package models

type AuditLog struct {
	BaseModel
	EntityID   *uint  `gorm:"index" json:"entity_id,omitempty"`
	ModuleCode string `gorm:"not null;size:50;index" json:"module_code"`
	ObjectType string `gorm:"not null;size:30;index" json:"object_type"`
	ObjectID   uint   `gorm:"not null;index" json:"object_id"`
	EventCode  string `gorm:"size:50;index" json:"event_code,omitempty"`
	ActorType  string `gorm:"not null;size:20;default:'internal_user'" json:"actor_type"`
	ActorID    *uint  `gorm:"index" json:"actor_id,omitempty"`
	Action     string `gorm:"not null;size:50" json:"action"`
	DataBefore string `gorm:"type:longtext" json:"data_before,omitempty"`
	DataAfter  string `gorm:"type:longtext" json:"data_after,omitempty"`
	IPAddress  string `gorm:"size:100" json:"ip_address,omitempty"`
	UserAgent  string `gorm:"type:text" json:"user_agent,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
