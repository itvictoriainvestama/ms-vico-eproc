package models

import "time"

const (
	RFQStatusCreated          = "Created"
	RFQStatusPublished        = "Published"
	RFQStatusVendorSubmission = "Vendor Submission"
	RFQStatusClosed           = "Closed"
	RFQStatusReopened         = "Reopened"
	RFQStatusEvaluation       = "Evaluation"
	RFQStatusBAFO             = "BAFO"
	RFQStatusVendorSelected   = "Vendor Selected"
	RFQStatusCancelled        = "Cancelled"
)

type RFQ struct {
	BaseModel
	EntityID              uint                `gorm:"not null;index" json:"entity_id"`
	PRID                  uint                `gorm:"not null;index" json:"pr_id"`
	PR                    PurchaseRequisition `gorm:"foreignKey:PRID" json:"pr"`
	RFQNumber             string              `gorm:"uniqueIndex;not null;size:50" json:"rfq_number"`
	Title                 string              `gorm:"not null;size:255" json:"title"`
	TechnicalRequirement  string              `gorm:"type:text" json:"technical_requirement,omitempty"`
	CommercialRequirement string              `gorm:"type:text" json:"commercial_requirement,omitempty"`
	MinimumVendorCount    int                 `gorm:"not null;default:1" json:"minimum_vendor_count"`
	PublishAt             *time.Time          `json:"publish_at,omitempty"`
	DeadlineAt            time.Time           `gorm:"not null" json:"deadline_at"`
	Status                string              `gorm:"not null;size:30;default:'Created'" json:"status"`
	Vendors               []RFQVendor         `gorm:"foreignKey:RFQID" json:"vendors,omitempty"`
}

func (RFQ) TableName() string {
	return "rfqs"
}

type RFQVendor struct {
	BaseModel
	RFQID               uint       `gorm:"not null;index" json:"rfq_id"`
	VendorID            uint       `gorm:"not null;index" json:"vendor_id"`
	Vendor              Vendor     `gorm:"foreignKey:VendorID" json:"vendor"`
	InvitationType      string     `gorm:"not null;size:20;default:'restricted'" json:"invitation_type"`
	EligibilityStatus   string     `gorm:"not null;size:20;default:'eligible'" json:"eligibility_status"`
	InvitedAt           *time.Time `json:"invited_at,omitempty"`
	ViewedAt            *time.Time `json:"viewed_at,omitempty"`
	ParticipationStatus string     `gorm:"not null;size:20;default:'invited'" json:"participation_status"`
}

func (RFQVendor) TableName() string {
	return "rfq_vendors"
}

type VendorBid struct {
	BaseModel
	RFQID           uint       `gorm:"not null;index" json:"rfq_id"`
	VendorID        uint       `gorm:"not null;index" json:"vendor_id"`
	EntityID        uint       `gorm:"not null;index" json:"entity_id"`
	QuotationNumber string     `gorm:"size:50" json:"quotation_number,omitempty"`
	QuotationTotal  float64    `gorm:"type:decimal(18,2);default:0" json:"quotation_total"`
	PaymentTerms    string     `gorm:"size:255" json:"payment_terms,omitempty"`
	DeliveryTerms   string     `gorm:"size:255" json:"delivery_terms,omitempty"`
	ValidUntil      *time.Time `json:"valid_until,omitempty"`
	SubmittedAt     *time.Time `json:"submitted_at,omitempty"`
	Status          string     `gorm:"not null;size:30;default:'Draft'" json:"status"`
	Items           []BidItem  `gorm:"foreignKey:BidID" json:"items,omitempty"`
}

func (VendorBid) TableName() string {
	return "quotations"
}

type BidItem struct {
	BaseModel
	BidID      uint    `gorm:"not null;index" json:"quotation_id"`
	PRItemID   *uint   `gorm:"index" json:"pr_item_id,omitempty"`
	UnitPrice  float64 `gorm:"type:decimal(18,2);not null" json:"unit_price"`
	Qty        float64 `gorm:"type:decimal(18,4);not null" json:"qty"`
	TotalPrice float64 `gorm:"type:decimal(18,2);not null" json:"total_price"`
	Note       string  `gorm:"type:text" json:"note,omitempty"`
}

func (BidItem) TableName() string {
	return "quotation_items"
}

type VendorEvaluation struct {
	BaseModel
	RFQID             uint       `gorm:"not null;index" json:"rfq_id"`
	VendorID          uint       `gorm:"not null;index" json:"vendor_id"`
	TechnicalScore    float64    `gorm:"type:decimal(8,2);default:0" json:"technical_score"`
	CommercialScore   float64    `gorm:"type:decimal(8,2);default:0" json:"commercial_score"`
	WeightedScore     float64    `gorm:"type:decimal(8,2);default:0" json:"weighted_score"`
	RankingNo         *int       `json:"ranking_no,omitempty"`
	EvaluationSummary string     `gorm:"type:text" json:"evaluation_summary,omitempty"`
	EvaluatedBy       *uint      `gorm:"index" json:"evaluated_by,omitempty"`
	EvaluatedAt       *time.Time `json:"evaluated_at,omitempty"`
}

func (VendorEvaluation) TableName() string {
	return "vendor_evaluations"
}

type BAFORound struct {
	BaseModel
	RFQID      uint      `gorm:"not null;index" json:"rfq_id"`
	RoundNo    int       `gorm:"not null" json:"round_no"`
	DeadlineAt time.Time `gorm:"not null" json:"deadline_at"`
	Status     string    `gorm:"not null;size:20;default:'Open'" json:"status"`
}

func (BAFORound) TableName() string {
	return "bafo_rounds"
}

type VendorSelection struct {
	BaseModel
	RFQID            uint    `gorm:"not null;index" json:"rfq_id"`
	SelectedVendorID uint    `gorm:"not null;index" json:"selected_vendor_id"`
	SelectionMethod  string  `gorm:"not null;size:30" json:"selection_method"`
	SelectionReason  string  `gorm:"type:text" json:"selection_reason"`
	FinalPrice       float64 `gorm:"type:decimal(18,2);default:0" json:"final_price"`
	SelectedBy       *uint   `gorm:"index" json:"selected_by,omitempty"`
}

func (VendorSelection) TableName() string {
	return "vendor_selections"
}

type DirectAppointment struct {
	BaseModel
	EntityID              uint    `gorm:"not null;index" json:"entity_id"`
	PRID                  uint    `gorm:"not null;index" json:"pr_id"`
	VendorID              uint    `gorm:"not null;index" json:"vendor_id"`
	Justification         string  `gorm:"type:text;not null" json:"justification"`
	ReferenceDocumentNote string  `gorm:"type:text" json:"reference_document_note,omitempty"`
	EstimatedValue        float64 `gorm:"type:decimal(18,2);default:0" json:"estimated_value"`
	Status                string  `gorm:"not null;size:30;default:'Created'" json:"status"`
}

func (DirectAppointment) TableName() string {
	return "direct_appointments"
}
