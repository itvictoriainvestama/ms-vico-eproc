package services

import (
	"errors"
	"fmt"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type VendorService struct {
	db *gorm.DB
}

func NewVendorService(db *gorm.DB) *VendorService {
	return &VendorService{db: db}
}

type CreateVendorRequest struct {
	VendorName string `json:"vendor_name" binding:"required"`
	TaxID      string `json:"tax_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
}

func (s *VendorService) List(page, pageSize int, activeOnly bool) ([]models.Vendor, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&models.Vendor{})
	if activeOnly {
		query = query.Where("approved_status = ?", "approved")
	}

	var total int64
	query.Count(&total)

	var vendors []models.Vendor
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("vendor_name ASC").Find(&vendors).Error
	return vendors, total, err
}

func (s *VendorService) GetByID(id uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := s.db.First(&vendor, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("vendor not found")
	}
	return &vendor, err
}

func (s *VendorService) Create(req CreateVendorRequest) (*models.Vendor, error) {
	code, err := s.generateVendorCode()
	if err != nil {
		return nil, err
	}

	vendor := &models.Vendor{
		VendorCode:        code,
		VendorName:        req.VendorName,
		TaxID:             req.TaxID,
		Email:             req.Email,
		Phone:             req.Phone,
		Address:           req.Address,
		ApprovedStatus:    "approved",
		BlacklistStatus:   false,
		EligibilityStatus: "eligible",
	}

	if err := s.db.Create(vendor).Error; err != nil {
		return nil, err
	}
	return vendor, nil
}

func (s *VendorService) Update(id uint, req CreateVendorRequest) (*models.Vendor, error) {
	var vendor models.Vendor
	if err := s.db.First(&vendor, id).Error; err != nil {
		return nil, errors.New("vendor not found")
	}

	vendor.VendorName = req.VendorName
	vendor.TaxID = req.TaxID
	vendor.Email = req.Email
	vendor.Phone = req.Phone
	vendor.Address = req.Address

	if err := s.db.Save(&vendor).Error; err != nil {
		return nil, err
	}
	return &vendor, nil
}

func (s *VendorService) generateVendorCode() (string, error) {
	var count int64
	s.db.Model(&models.Vendor{}).Count(&count)
	return fmt.Sprintf("V-%04d", count+1), nil
}
