package models

import "time"

const (
	ApprovalStatusPending  = "pending"
	ApprovalStatusApproved = "approved"
	ApprovalStatusRejected = "rejected"
)

type ApprovalTask struct {
	BaseModel
	EntityID         uint       `gorm:"not null;index" json:"entity_id"`
	AssigneeID       uint       `gorm:"not null;index" json:"assignee_id"`
	Assignee         User       `gorm:"foreignKey:AssigneeID" json:"assignee"`
	OriginalUserID   *uint      `gorm:"index" json:"original_user_id,omitempty"`
	DocumentType     string     `gorm:"not null;size:20;index" json:"document_type"`
	DocumentID       uint       `gorm:"not null;index" json:"document_id"`
	RefNumber        string     `gorm:"not null;size:50" json:"ref_number"`
	Summary          string     `gorm:"size:255" json:"summary,omitempty"`
	Amount           float64    `gorm:"type:decimal(18,2);default:0" json:"amount"`
	Priority         string     `gorm:"not null;size:20;default:'Normal'" json:"priority"`
	ApprovalLevel    int        `gorm:"not null;default:1" json:"approval_level"`
	ApproverRoleCode string     `gorm:"size:50" json:"approver_role_code,omitempty"`
	Status           string     `gorm:"not null;size:20;default:'pending'" json:"status"`
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	Deadline         *time.Time `json:"deadline,omitempty"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
}
