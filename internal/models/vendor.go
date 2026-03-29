package models

type Vendor struct {
	BaseModel
	VendorCode        string `gorm:"uniqueIndex;not null;size:50" json:"vendor_code"`
	VendorName        string `gorm:"not null;size:200" json:"vendor_name"`
	TaxID             string `gorm:"size:100" json:"tax_id,omitempty"`
	Email             string `gorm:"size:150" json:"email,omitempty"`
	Phone             string `gorm:"size:50" json:"phone,omitempty"`
	Address           string `gorm:"type:text" json:"address,omitempty"`
	ApprovedStatus    string `gorm:"not null;size:20;default:'approved'" json:"approved_status"`
	BlacklistStatus   bool   `gorm:"not null;default:false" json:"blacklist_status"`
	EligibilityStatus string `gorm:"not null;size:20;default:'eligible'" json:"eligibility_status"`
}
