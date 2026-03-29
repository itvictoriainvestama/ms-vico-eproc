package models

import "time"

const (
	POStatusDraft           = "Draft"
	POStatusPendingApproval = "Pending Approval"
	POStatusApproved        = "Approved"
	POStatusRejected        = "Rejected"
	POStatusSentToVendor    = "Sent to Vendor"
	POStatusVendorConfirmed = "Vendor Confirmed"
	POStatusCompleted       = "Completed"
	POStatusVoided          = "Voided"
)

type PurchaseOrder struct {
	BaseModel
	EntityID          uint         `gorm:"not null;index" json:"entity_id"`
	PONumber          string       `gorm:"uniqueIndex;not null;size:50" json:"po_number"`
	PRID              *uint        `gorm:"index" json:"pr_id,omitempty"`
	RFQID             *uint        `gorm:"index" json:"rfq_id,omitempty"`
	DAID              *uint        `gorm:"column:da_id;index" json:"da_id,omitempty"`
	VendorID          uint         `gorm:"not null;index" json:"vendor_id"`
	Vendor            Vendor       `gorm:"foreignKey:VendorID" json:"vendor"`
	CurrencyCode      string       `gorm:"not null;size:10;default:'IDR'" json:"currency_code"`
	TotalAmount       float64      `gorm:"type:decimal(18,2);default:0" json:"total_amount"`
	Status            string       `gorm:"not null;size:30;default:'Draft'" json:"status"`
	SentAt            *time.Time   `json:"sent_at,omitempty"`
	VendorConfirmedAt *time.Time   `json:"vendor_confirmed_at,omitempty"`
	PODate            time.Time    `gorm:"not null" json:"po_date"`
	ExpectedDate      *time.Time   `json:"expected_date,omitempty"`
	DeliveryAddress   string       `gorm:"type:text" json:"delivery_address,omitempty"`
	PaymentTerms      string       `gorm:"size:255" json:"payment_terms,omitempty"`
	Notes             string       `gorm:"type:text" json:"notes,omitempty"`
	Items             []POItem     `gorm:"foreignKey:POID" json:"items,omitempty"`
	Approvals         []POApproval `gorm:"foreignKey:POID" json:"approvals,omitempty"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

type POItem struct {
	BaseModel
	POID          uint    `gorm:"not null;index" json:"po_id"`
	PRItemID      *uint   `gorm:"index" json:"pr_item_id,omitempty"`
	ItemName      string  `gorm:"not null;size:200" json:"item_name"`
	Specification string  `gorm:"type:text" json:"specification,omitempty"`
	Qty           float64 `gorm:"type:decimal(18,4);not null" json:"qty"`
	UOM           string  `gorm:"not null;size:50" json:"uom"`
	UnitPrice     float64 `gorm:"type:decimal(18,2);not null" json:"unit_price"`
	TotalPrice    float64 `gorm:"type:decimal(18,2);not null" json:"total_price"`
}

func (POItem) TableName() string {
	return "purchase_order_items"
}

type POApproval struct {
	BaseModel
	POID               uint       `gorm:"not null;index" json:"po_id"`
	ApprovalLevel      int        `gorm:"not null" json:"approval_level"`
	AssignedApproverID *uint      `gorm:"index" json:"assigned_approver_id,omitempty"`
	ActedByUserID      *uint      `gorm:"index" json:"acted_by_user_id,omitempty"`
	OnBehalfOfUserID   *uint      `gorm:"index" json:"on_behalf_of_user_id,omitempty"`
	Action             string     `gorm:"size:20" json:"action,omitempty"`
	Remarks            string     `gorm:"type:text" json:"remarks,omitempty"`
	ActedAt            *time.Time `json:"acted_at,omitempty"`
	StatusAfterAction  string     `gorm:"size:30" json:"status_after_action,omitempty"`
}

func (POApproval) TableName() string {
	return "po_approvals"
}

type VendorConfirmation struct {
	BaseModel
	POID                    uint       `gorm:"not null;index" json:"po_id"`
	VendorID                uint       `gorm:"not null;index" json:"vendor_id"`
	ConfirmedByVendorUserID *uint      `gorm:"index" json:"confirmed_by_vendor_user_id,omitempty"`
	ConfirmationStatus      string     `gorm:"not null;size:20" json:"confirmation_status"`
	Remarks                 string     `gorm:"type:text" json:"remarks,omitempty"`
	ConfirmedAt             *time.Time `json:"confirmed_at,omitempty"`
}

func (VendorConfirmation) TableName() string {
	return "vendor_confirmations"
}
