package services

import (
	"errors"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type EntityService struct {
	db *gorm.DB
}

func NewEntityService(db *gorm.DB) *EntityService {
	return &EntityService{db: db}
}

type EntityListResult struct {
	Items []models.Entity `json:"items"`
	Total int64           `json:"total"`
}

type CreateEntityRequest struct {
	EntityCode        string `json:"entity_code" binding:"required"`
	EntityName        string `json:"entity_name" binding:"required"`
	ParentEntityID    *uint  `json:"parent_entity_id"`
	EntityType        string `json:"entity_type"`
	ApprovalModelCode string `json:"approval_model_code"`
	GovernanceMode    string `json:"governance_mode"`
	Status            string `json:"status"`
}

func (s *EntityService) List(actorEntityID uint, scopeType string) (*EntityListResult, error) {
	query := s.db.Model(&models.Entity{})
	query = applyEntityScope(query, "id", actorEntityID, scopeType)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []models.Entity
	if err := query.Order("entity_name ASC").Find(&items).Error; err != nil {
		return nil, err
	}

	return &EntityListResult{Items: items, Total: total}, nil
}

func (s *EntityService) GetByID(id, actorEntityID uint, scopeType string) (*models.Entity, error) {
	var entity models.Entity
	query := applyEntityScope(s.db, "id", actorEntityID, scopeType)
	if err := query.First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("entity not found")
		}
		return nil, err
	}
	return &entity, nil
}

func (s *EntityService) Create(req CreateEntityRequest) (*models.Entity, error) {
	entityType := req.EntityType
	if entityType == "" {
		entityType = "subsidiary"
	}

	governanceMode := req.GovernanceMode
	if governanceMode == "" {
		governanceMode = "entity_only"
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	entity := &models.Entity{
		EntityCode:        req.EntityCode,
		EntityName:        req.EntityName,
		ParentEntityID:    req.ParentEntityID,
		EntityType:        entityType,
		Status:            status,
		ApprovalModelCode: req.ApprovalModelCode,
		GovernanceMode:    governanceMode,
	}

	if err := s.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
