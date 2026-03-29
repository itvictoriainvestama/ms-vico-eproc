package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type RFQService struct {
	db *gorm.DB
}

func NewRFQService(db *gorm.DB) *RFQService {
	return &RFQService{db: db}
}

type CreateRFQRequest struct {
	PRID                  uint      `json:"pr_id" binding:"required"`
	Title                 string    `json:"title" binding:"required"`
	TechnicalRequirement  string    `json:"technical_requirement"`
	CommercialRequirement string    `json:"commercial_requirement"`
	MinimumVendorCount    int       `json:"minimum_vendor_count"`
	DeadlineAt            time.Time `json:"deadline_at" binding:"required"`
	VendorIDs             []uint    `json:"vendor_ids" binding:"required,min=1"`
}

type VendorQuotationRequest struct {
	PaymentTerms  string                `json:"payment_terms"`
	DeliveryTerms string                `json:"delivery_terms"`
	ValidUntil    *time.Time            `json:"valid_until"`
	Items         []VendorQuotationItem `json:"items" binding:"required,min=1"`
}

type VendorQuotationItem struct {
	PRItemID  *uint   `json:"pr_item_id"`
	Qty       float64 `json:"qty" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gt=0"`
	Note      string  `json:"note"`
}

func (s *RFQService) List(page, pageSize int, status string, entityID uint) ([]models.RFQ, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&models.RFQ{}).Preload("PR").Preload("Vendors").Preload("Vendors.Vendor")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if entityID > 0 {
		query = query.Where("entity_id = ?", entityID)
	}

	var total int64
	query.Count(&total)

	var rfqs []models.RFQ
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&rfqs).Error
	return rfqs, total, err
}

func (s *RFQService) GetByID(id uint) (*models.RFQ, error) {
	return s.GetByIDScoped(id, 0, ScopeCrossEntity)
}

func (s *RFQService) GetByIDScoped(id, actorEntityID uint, scopeType string) (*models.RFQ, error) {
	var rfq models.RFQ
	query := s.db.
		Preload("PR").
		Preload("Vendors").
		Preload("Vendors.Vendor")
	query = applyEntityScope(query, "entity_id", actorEntityID, scopeType)

	err := query.First(&rfq, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("rfq not found")
	}
	return &rfq, err
}

func (s *RFQService) Create(req CreateRFQRequest, entityID uint) (*models.RFQ, error) {
	var pr models.PurchaseRequisition
	if err := s.db.Select("id", "entity_id").First(&pr, req.PRID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("purchase request not found")
		}
		return nil, err
	}
	if pr.EntityID != entityID {
		return nil, errors.New("purchase request is outside your entity scope")
	}

	rfqNumber, err := s.generateRFQNumber()
	if err != nil {
		return nil, err
	}

	minimumVendorCount := req.MinimumVendorCount
	if minimumVendorCount <= 0 {
		minimumVendorCount = 1
	}

	vendors := make([]models.RFQVendor, 0, len(req.VendorIDs))
	for _, vendorID := range req.VendorIDs {
		var vendor models.Vendor
		if err := s.db.Select("id", "blacklist_status", "eligibility_status").First(&vendor, vendorID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("vendor not found")
			}
			return nil, err
		}
		if vendor.BlacklistStatus || vendor.EligibilityStatus != "eligible" {
			return nil, errors.New("all invited vendors must be eligible and not blacklisted")
		}

		vendors = append(vendors, models.RFQVendor{
			VendorID:            vendorID,
			InvitationType:      "restricted",
			EligibilityStatus:   "eligible",
			ParticipationStatus: "invited",
		})
	}

	rfq := &models.RFQ{
		EntityID:              entityID,
		PRID:                  req.PRID,
		RFQNumber:             rfqNumber,
		Title:                 req.Title,
		TechnicalRequirement:  req.TechnicalRequirement,
		CommercialRequirement: req.CommercialRequirement,
		MinimumVendorCount:    minimumVendorCount,
		DeadlineAt:            req.DeadlineAt,
		Status:                models.RFQStatusCreated,
		Vendors:               vendors,
	}

	if err := s.db.Create(rfq).Error; err != nil {
		return nil, err
	}
	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(entityID),
		ModuleCode: "RFQ",
		ObjectType: "RFQ",
		ObjectID:   rfq.ID,
		EventCode:  "RFQ_CREATED",
		ActorType:  "internal_user",
		Action:     "create",
	})
	return s.GetByIDScoped(rfq.ID, entityID, "")
}

func (s *RFQService) UpdateStatus(id uint, actorEntityID uint, scopeType, newStatus string) (*models.RFQ, error) {
	var rfq models.RFQ
	if err := s.db.First(&rfq, id).Error; err != nil {
		return nil, errors.New("rfq not found")
	}
	if err := ensureEntityAccess(rfq.EntityID, actorEntityID, scopeType); err != nil {
		return nil, err
	}
	if err := s.db.Model(&rfq).Update("status", newStatus).Error; err != nil {
		return nil, err
	}
	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(rfq.EntityID),
		ModuleCode: "RFQ",
		ObjectType: "RFQ",
		ObjectID:   rfq.ID,
		EventCode:  "RFQ_STATUS_UPDATED",
		ActorType:  "internal_user",
		Action:     "update_status",
		DataAfter:  newStatus,
	})
	return s.GetByIDScoped(id, actorEntityID, scopeType)
}

func (s *RFQService) ListVendorTenders(vendorID uint) ([]models.RFQ, error) {
	var rfqs []models.RFQ
	err := s.db.
		Joins("JOIN rfq_vendors ON rfq_vendors.rfq_id = rfqs.id").
		Where("rfq_vendors.vendor_id = ? AND rfq_vendors.eligibility_status = ? AND rfq_vendors.participation_status IN ? AND rfqs.status IN ?", vendorID, "eligible", []string{"invited", "viewed", "submitted"}, []string{models.RFQStatusPublished, models.RFQStatusVendorSubmission}).
		Preload("PR").
		Find(&rfqs).Error
	return rfqs, err
}

func (s *RFQService) GetVendorTender(id, vendorID uint) (*models.RFQ, error) {
	var rfq models.RFQ
	if err := s.db.
		Joins("JOIN rfq_vendors ON rfq_vendors.rfq_id = rfqs.id").
		Where("rfqs.id = ? AND rfq_vendors.vendor_id = ?", id, vendorID).
		Preload("PR").
		Preload("Vendors", "vendor_id = ?", vendorID).
		First(&rfq).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found")
		}
		return nil, err
	}

	now := time.Now()
	_ = s.db.Model(&models.RFQVendor{}).
		Where("rfq_id = ? AND vendor_id = ?", id, vendorID).
		Updates(map[string]interface{}{
			"viewed_at":            &now,
			"participation_status": "viewed",
		}).Error
	return &rfq, nil
}

func (s *RFQService) SubmitQuotation(rfqID, vendorID, vendorUserID uint, req VendorQuotationRequest) (*models.VendorBid, error) {
	var vendor models.Vendor
	if err := s.db.First(&vendor, vendorID).Error; err != nil {
		return nil, errors.New("vendor not found")
	}
	if vendor.BlacklistStatus || vendor.EligibilityStatus != "eligible" {
		return nil, errors.New("vendor is not eligible for tender participation")
	}

	var rfq models.RFQ
	if err := s.db.First(&rfq, rfqID).Error; err != nil {
		return nil, errors.New("rfq not found")
	}
	if time.Now().After(rfq.DeadlineAt) {
		return nil, errors.New("batas waktu bidding telah berakhir")
	}
	if rfq.Status != models.RFQStatusPublished && rfq.Status != models.RFQStatusVendorSubmission {
		return nil, errors.New("rfq is not open for vendor submission")
	}

	var invitation models.RFQVendor
	if err := s.db.Where("rfq_id = ? AND vendor_id = ?", rfqID, vendorID).First(&invitation).Error; err != nil {
		return nil, errors.New("vendor is not invited to this rfq")
	}

	var total float64
	items := make([]models.BidItem, 0, len(req.Items))
	for _, item := range req.Items {
		lineTotal := item.Qty * item.UnitPrice
		total += lineTotal
		items = append(items, models.BidItem{
			PRItemID:   item.PRItemID,
			Qty:        item.Qty,
			UnitPrice:  item.UnitPrice,
			TotalPrice: lineTotal,
			Note:       item.Note,
		})
	}

	submittedAt := time.Now()
	bid := &models.VendorBid{
		RFQID:          rfqID,
		VendorID:       vendorID,
		EntityID:       rfq.EntityID,
		QuotationTotal: total,
		PaymentTerms:   req.PaymentTerms,
		DeliveryTerms:  req.DeliveryTerms,
		ValidUntil:     req.ValidUntil,
		SubmittedAt:    &submittedAt,
		Status:         "Submitted",
		Items:          items,
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(bid).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.RFQVendor{}).
			Where("rfq_id = ? AND vendor_id = ?", rfqID, vendorID).
			Update("participation_status", "submitted").Error; err != nil {
			return err
		}
		return tx.Model(&rfq).Where("id = ?", rfqID).Update("status", models.RFQStatusVendorSubmission).Error
	}); err != nil {
		return nil, err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(rfq.EntityID),
		ModuleCode: "RFQ",
		ObjectType: "QUOTATION",
		ObjectID:   bid.ID,
		EventCode:  "QUOTATION_SUBMITTED",
		ActorType:  "vendor_user",
		ActorID:    uintPtr(vendorUserID),
		Action:     "submit_quotation",
	})
	return bid, nil
}

func (s *RFQService) generateRFQNumber() (string, error) {
	var count int64
	now := time.Now()
	s.db.Model(&models.RFQ{}).
		Where("YEAR(created_at) = ?", now.Year()).
		Count(&count)
	return fmt.Sprintf("RFQ-%d-%04d", now.Year(), count+1), nil
}
