package models

import "time"

const (
	PRStatusDraft           = "Draft"
	PRStatusSubmitted       = "Submitted"
	PRStatusPendingApproval = "Pending Approval"
	PRStatusApproved        = "Approved"
	PRStatusRejected        = "Rejected"
	PRStatusRevised         = "Revised"
	PRStatusCancelled       = "Cancelled"
)

const (
	ProcurementGoods    = "barang"
	ProcurementServices = "jasa"
)

const (
	RoutineRutin    = "rutin"
	RoutineNonRutin = "non-rutin"
)

const (
	BudgetWithin = "within_budget"
	BudgetOver   = "over_budget"
	BudgetNon    = "non_budget"
)

const (
	MethodBidding           = "bidding"
	MethodDirectAppointment = "direct_appointment"
)

type PurchaseRequisition struct {
	BaseModel
	EntityID          uint           `gorm:"not null;index" json:"entity_id"`
	Entity            Entity         `gorm:"foreignKey:EntityID" json:"entity"`
	PRNumber          string         `gorm:"uniqueIndex;not null;size:50" json:"pr_number"`
	RequestorID       uint           `gorm:"not null;index" json:"requestor_id"`
	Requestor         User           `gorm:"foreignKey:RequestorID" json:"requestor"`
	DepartmentCode    string         `gorm:"size:50;index" json:"department_code,omitempty"`
	Title             string         `gorm:"not null;size:200" json:"title"`
	Description       string         `gorm:"type:text" json:"description"`
	ProcurementType   string         `gorm:"not null;size:20" json:"procurement_type"`
	RoutineType       string         `gorm:"not null;size:20" json:"routine_type"`
	EstimatedAmount   float64        `gorm:"type:decimal(18,2);default:0" json:"estimated_amount"`
	BudgetStatus      string         `gorm:"not null;size:20;default:'within_budget'" json:"budget_status"`
	NeedDate          time.Time      `gorm:"not null" json:"need_date"`
	ProcurementMethod *string        `gorm:"size:30" json:"procurement_method,omitempty"`
	Status            string         `gorm:"not null;size:30;default:'Draft'" json:"status"`
	Items             []PRItem       `gorm:"foreignKey:PRID" json:"items,omitempty"`
	Attachments       []PRAttachment `gorm:"foreignKey:PRID" json:"attachments,omitempty"`
	Approvals         []PRApproval   `gorm:"foreignKey:PRID" json:"approvals,omitempty"`
}

func (PurchaseRequisition) TableName() string {
	return "purchase_requests"
}

type PRItem struct {
	BaseModel
	PRID                uint    `gorm:"not null;index" json:"pr_id"`
	EntityID            uint    `gorm:"not null;index" json:"entity_id"`
	ItemName            string  `gorm:"not null;size:200" json:"item_name"`
	Specification       string  `gorm:"type:text" json:"specification,omitempty"`
	Qty                 float64 `gorm:"type:decimal(18,4);not null" json:"qty"`
	UOM                 string  `gorm:"not null;size:50" json:"uom"`
	EstimatedUnitPrice  float64 `gorm:"type:decimal(18,2);default:0" json:"estimated_unit_price"`
	EstimatedTotalPrice float64 `gorm:"type:decimal(18,2);default:0" json:"estimated_total_price"`
	ReferencePriceID    *uint   `gorm:"index" json:"reference_price_id,omitempty"`
}

func (PRItem) TableName() string {
	return "purchase_request_items"
}

type PRAttachment struct {
	BaseModel
	PRID             uint   `gorm:"not null;index" json:"pr_id"`
	EntityID         uint   `gorm:"not null;index" json:"entity_id"`
	FileNameOriginal string `gorm:"not null;size:255" json:"file_name_original"`
	FileNameStored   string `gorm:"size:255" json:"file_name_stored,omitempty"`
	BucketName       string `gorm:"size:100" json:"bucket_name,omitempty"`
	ObjectKey        string `gorm:"not null;size:255" json:"object_key"`
	MimeType         string `gorm:"size:100" json:"mime_type,omitempty"`
	FileSize         int64  `gorm:"default:0" json:"file_size"`
	UploadedBy       *uint  `gorm:"index" json:"uploaded_by,omitempty"`
}

func (PRAttachment) TableName() string {
	return "pr_attachments"
}

type PRApproval struct {
	BaseModel
	PRID               uint       `gorm:"not null;index" json:"pr_id"`
	EntityID           uint       `gorm:"not null;index" json:"entity_id"`
	ApprovalLevel      int        `gorm:"not null" json:"approval_level"`
	AssignedApproverID *uint      `gorm:"index" json:"assigned_approver_id,omitempty"`
	ActedByUserID      *uint      `gorm:"index" json:"acted_by_user_id,omitempty"`
	OnBehalfOfUserID   *uint      `gorm:"index" json:"on_behalf_of_user_id,omitempty"`
	Action             string     `gorm:"size:20" json:"action,omitempty"`
	Remarks            string     `gorm:"type:text" json:"remarks,omitempty"`
	ActedAt            *time.Time `json:"acted_at,omitempty"`
	StatusAfterAction  string     `gorm:"size:30" json:"status_after_action,omitempty"`
}

func (PRApproval) TableName() string {
	return "pr_approvals"
}
