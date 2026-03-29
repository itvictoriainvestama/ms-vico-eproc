package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type PRService struct {
	db *gorm.DB
}

func NewPRService(db *gorm.DB) *PRService {
	return &PRService{db: db}
}

type PRListParams struct {
	Page           int
	PageSize       int
	Status         string
	EntityID       uint
	DepartmentCode string
	RequestorID    uint
}

type PRListResult struct {
	Items    []models.PurchaseRequisition `json:"items"`
	Total    int64                        `json:"total"`
	Page     int                          `json:"page"`
	PageSize int                          `json:"page_size"`
}

type CreatePRRequest struct {
	Title           string         `json:"title" binding:"required"`
	Description     string         `json:"description" binding:"required"`
	DepartmentCode  string         `json:"department_code" binding:"required"`
	ProcurementType string         `json:"procurement_type" binding:"required"`
	RoutineType     string         `json:"routine_type" binding:"required"`
	BudgetStatus    string         `json:"budget_status"`
	NeedDate        time.Time      `json:"need_date" binding:"required"`
	Items           []CreatePRItem `json:"items" binding:"required,min=1"`
}

type CreatePRItem struct {
	ItemName           string  `json:"item_name" binding:"required"`
	Specification      string  `json:"specification"`
	Qty                float64 `json:"qty" binding:"required,gt=0"`
	UOM                string  `json:"uom" binding:"required"`
	EstimatedUnitPrice float64 `json:"estimated_unit_price" binding:"required,gt=0"`
}

func (s *PRService) List(params PRListParams) (*PRListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	query := s.db.Model(&models.PurchaseRequisition{}).
		Preload("Requestor").
		Preload("Entity").
		Preload("Items")

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.EntityID > 0 {
		query = query.Where("entity_id = ?", params.EntityID)
	}
	if params.DepartmentCode != "" {
		query = query.Where("department_code = ?", params.DepartmentCode)
	}
	if params.RequestorID > 0 {
		query = query.Where("requestor_id = ?", params.RequestorID)
	}

	var total int64
	query.Count(&total)

	var items []models.PurchaseRequisition
	err := query.Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		Order("created_at DESC").
		Find(&items).Error

	return &PRListResult{Items: items, Total: total, Page: params.Page, PageSize: params.PageSize}, err
}

func (s *PRService) GetByID(id uint) (*models.PurchaseRequisition, error) {
	return s.GetByIDScoped(id, 0, ScopeCrossEntity)
}

func (s *PRService) GetByIDScoped(id, actorEntityID uint, scopeType string) (*models.PurchaseRequisition, error) {
	var pr models.PurchaseRequisition
	query := s.db.
		Preload("Entity").
		Preload("Requestor").
		Preload("Items").
		Preload("Approvals").
		Preload("Attachments")
	query = applyEntityScope(query, "entity_id", actorEntityID, scopeType)

	err := query.First(&pr, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("purchase request not found")
	}
	return &pr, err
}

func (s *PRService) Create(req CreatePRRequest, requestorID, entityID uint) (*models.PurchaseRequisition, error) {
	prNumber, err := s.generatePRNumber()
	if err != nil {
		return nil, err
	}

	budgetStatus := req.BudgetStatus
	if budgetStatus == "" {
		budgetStatus = models.BudgetWithin
	}

	var total float64
	items := make([]models.PRItem, 0, len(req.Items))
	for _, item := range req.Items {
		lineTotal := item.Qty * item.EstimatedUnitPrice
		total += lineTotal
		items = append(items, models.PRItem{
			EntityID:            entityID,
			ItemName:            item.ItemName,
			Specification:       item.Specification,
			Qty:                 item.Qty,
			UOM:                 item.UOM,
			EstimatedUnitPrice:  item.EstimatedUnitPrice,
			EstimatedTotalPrice: lineTotal,
		})
	}

	pr := &models.PurchaseRequisition{
		EntityID:        entityID,
		PRNumber:        prNumber,
		RequestorID:     requestorID,
		DepartmentCode:  req.DepartmentCode,
		Title:           req.Title,
		Description:     req.Description,
		ProcurementType: req.ProcurementType,
		RoutineType:     req.RoutineType,
		EstimatedAmount: total,
		BudgetStatus:    budgetStatus,
		NeedDate:        req.NeedDate,
		Status:          models.PRStatusDraft,
		Items:           items,
	}

	if err := s.db.Create(pr).Error; err != nil {
		return nil, err
	}
	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(pr.EntityID),
		ModuleCode: "PROCUREMENT",
		ObjectType: "PR",
		ObjectID:   pr.ID,
		EventCode:  "PR_CREATED",
		ActorType:  "internal_user",
		ActorID:    uintPtr(requestorID),
		Action:     "create",
	})
	return s.GetByID(pr.ID)
}

func (s *PRService) Submit(id uint, actorID, actorEntityID uint, scopeType string) (*models.PurchaseRequisition, error) {
	var pr models.PurchaseRequisition
	if err := s.db.First(&pr, id).Error; err != nil {
		return nil, errors.New("purchase request not found")
	}
	if err := ensureEntityAccess(pr.EntityID, actorEntityID, scopeType); err != nil {
		return nil, err
	}
	if pr.Status != models.PRStatusDraft && pr.Status != models.PRStatusRevised {
		return nil, errors.New("only draft or revised purchase requests can be submitted")
	}

	if err := s.db.Transaction(func(tx *gorm.DB) error {
		assignee, originalUserID, err := resolveEntityApprover(tx, pr.EntityID, "APPROVER", "ENTITY_ADMIN", "SUPER_ADMIN")
		if err != nil {
			return errors.New("no approver is configured for this entity")
		}

		if err := tx.Model(&pr).Update("status", models.PRStatusPendingApproval).Error; err != nil {
			return err
		}

		task := models.ApprovalTask{
			EntityID:         pr.EntityID,
			AssigneeID:       assignee.ID,
			OriginalUserID:   originalUserID,
			DocumentType:     "PR",
			DocumentID:       pr.ID,
			RefNumber:        pr.PRNumber,
			Summary:          pr.Title,
			Amount:           pr.EstimatedAmount,
			Priority:         "Normal",
			ApprovalLevel:    1,
			ApproverRoleCode: primaryRoleCode(assignee),
			Status:           models.ApprovalStatusPending,
		}
		if err := tx.Create(&task).Error; err != nil {
			return err
		}

		approval := models.PRApproval{
			PRID:               pr.ID,
			EntityID:           pr.EntityID,
			ApprovalLevel:      1,
			AssignedApproverID: &assignee.ID,
		}
		return tx.Create(&approval).Error
	}); err != nil {
		return nil, err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(pr.EntityID),
		ModuleCode: "PROCUREMENT",
		ObjectType: "PR",
		ObjectID:   pr.ID,
		EventCode:  "PR_SUBMITTED",
		ActorType:  "internal_user",
		ActorID:    uintPtr(actorID),
		Action:     "submit",
	})

	return s.GetByIDScoped(id, actorEntityID, scopeType)
}

func (s *PRService) generatePRNumber() (string, error) {
	var count int64
	now := time.Now()
	s.db.Model(&models.PurchaseRequisition{}).
		Where("YEAR(created_at) = ?", now.Year()).
		Count(&count)
	return fmt.Sprintf("PR-%d-%04d", now.Year(), count+1), nil
}
