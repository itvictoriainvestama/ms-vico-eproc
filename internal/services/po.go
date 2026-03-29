package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type POService struct {
	db *gorm.DB
}

func NewPOService(db *gorm.DB) *POService {
	return &POService{db: db}
}

type CreatePORequest struct {
	PRID            *uint          `json:"pr_id"`
	RFQID           *uint          `json:"rfq_id"`
	DAID            *uint          `json:"da_id"`
	VendorID        uint           `json:"vendor_id" binding:"required"`
	PODate          time.Time      `json:"po_date" binding:"required"`
	ExpectedDate    *time.Time     `json:"expected_date"`
	DeliveryAddress string         `json:"delivery_address" binding:"required"`
	PaymentTerms    string         `json:"payment_terms"`
	Notes           string         `json:"notes"`
	Items           []CreatePOItem `json:"items" binding:"required,min=1"`
}

type CreatePOItem struct {
	PRItemID      *uint   `json:"pr_item_id"`
	ItemName      string  `json:"item_name" binding:"required"`
	Specification string  `json:"specification"`
	Qty           float64 `json:"qty" binding:"required,gt=0"`
	UOM           string  `json:"uom" binding:"required"`
	UnitPrice     float64 `json:"unit_price" binding:"required,gt=0"`
}

func (s *POService) List(page, pageSize int, status string, vendorID, entityID uint) ([]models.PurchaseOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&models.PurchaseOrder{}).Preload("Vendor").Preload("Items")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if vendorID > 0 {
		query = query.Where("vendor_id = ?", vendorID)
	}
	if entityID > 0 {
		query = query.Where("entity_id = ?", entityID)
	}

	var total int64
	query.Count(&total)

	var pos []models.PurchaseOrder
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&pos).Error
	return pos, total, err
}

func (s *POService) GetByID(id uint) (*models.PurchaseOrder, error) {
	return s.GetByIDScoped(id, 0, ScopeCrossEntity)
}

func (s *POService) GetByIDScoped(id, actorEntityID uint, scopeType string) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	query := s.db.Preload("Vendor").Preload("Items").Preload("Approvals")
	query = applyEntityScope(query, "entity_id", actorEntityID, scopeType)

	err := query.First(&po, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("purchase order not found")
	}
	return &po, err
}

func (s *POService) Create(req CreatePORequest, entityID uint) (*models.PurchaseOrder, error) {
	if req.PRID != nil {
		var pr models.PurchaseRequisition
		if err := s.db.Select("id", "entity_id").First(&pr, *req.PRID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("purchase request not found")
			}
			return nil, err
		}
		if pr.EntityID != entityID {
			return nil, errors.New("purchase request is outside your entity scope")
		}
	}

	if req.RFQID != nil {
		var rfq models.RFQ
		if err := s.db.Select("id", "entity_id").First(&rfq, *req.RFQID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("rfq not found")
			}
			return nil, err
		}
		if rfq.EntityID != entityID {
			return nil, errors.New("rfq is outside your entity scope")
		}
	}

	poNumber, err := s.generatePONumber()
	if err != nil {
		return nil, err
	}

	var total float64
	items := make([]models.POItem, 0, len(req.Items))
	for _, item := range req.Items {
		lineTotal := item.Qty * item.UnitPrice
		total += lineTotal
		items = append(items, models.POItem{
			PRItemID:      item.PRItemID,
			ItemName:      item.ItemName,
			Specification: item.Specification,
			Qty:           item.Qty,
			UOM:           item.UOM,
			UnitPrice:     item.UnitPrice,
			TotalPrice:    lineTotal,
		})
	}

	po := &models.PurchaseOrder{
		EntityID:        entityID,
		PONumber:        poNumber,
		PRID:            req.PRID,
		RFQID:           req.RFQID,
		DAID:            req.DAID,
		VendorID:        req.VendorID,
		CurrencyCode:    "IDR",
		TotalAmount:     total,
		Status:          models.POStatusDraft,
		PODate:          req.PODate,
		ExpectedDate:    req.ExpectedDate,
		DeliveryAddress: req.DeliveryAddress,
		PaymentTerms:    req.PaymentTerms,
		Notes:           req.Notes,
		Items:           items,
	}

	if err := s.db.Create(po).Error; err != nil {
		return nil, err
	}
	return s.GetByIDScoped(po.ID, entityID, "")
}

func (s *POService) UpdateStatus(id uint, actorEntityID uint, scopeType, newStatus string) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	if err := s.db.First(&po, id).Error; err != nil {
		return nil, errors.New("purchase order not found")
	}
	if err := ensureEntityAccess(po.EntityID, actorEntityID, scopeType); err != nil {
		return nil, err
	}
	if err := s.db.Model(&po).Update("status", newStatus).Error; err != nil {
		return nil, err
	}
	return s.GetByIDScoped(id, actorEntityID, scopeType)
}

func (s *POService) generatePONumber() (string, error) {
	var count int64
	now := time.Now()
	s.db.Model(&models.PurchaseOrder{}).
		Where("YEAR(created_at) = ?", now.Year()).
		Count(&count)
	return fmt.Sprintf("PO-%d-%04d", now.Year(), count+1), nil
}
