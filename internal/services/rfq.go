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
	return s.GetByIDScoped(id, actorEntityID, scopeType)
}

func (s *RFQService) generateRFQNumber() (string, error) {
	var count int64
	now := time.Now()
	s.db.Model(&models.RFQ{}).
		Where("YEAR(created_at) = ?", now.Year()).
		Count(&count)
	return fmt.Sprintf("RFQ-%d-%04d", now.Year(), count+1), nil
}
