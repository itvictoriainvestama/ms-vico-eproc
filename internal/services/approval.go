package services

import (
	"errors"
	"time"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type ApprovalService struct {
	db *gorm.DB
}

func NewApprovalService(db *gorm.DB) *ApprovalService {
	return &ApprovalService{db: db}
}

type ApproveRequest struct {
	Notes string `json:"notes"`
}

type RejectRequest struct {
	Notes string `json:"notes" binding:"required"`
}

func (s *ApprovalService) GetTasksByUser(userID, actorEntityID uint, scopeType string, page, pageSize int) ([]models.ApprovalTask, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&models.ApprovalTask{}).Where("assignee_id = ?", userID)
	query = applyEntityScope(query, "entity_id", actorEntityID, scopeType)

	var total int64
	query.Count(&total)

	var tasks []models.ApprovalTask
	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error
	return tasks, total, err
}

func (s *ApprovalService) Approve(taskID, approverID, actorEntityID uint, scopeType string, req ApproveRequest) error {
	var task models.ApprovalTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return errors.New("approval task not found")
	}
	if err := ensureEntityAccess(task.EntityID, actorEntityID, scopeType); err != nil {
		return err
	}
	if task.AssigneeID != approverID {
		return errors.New("you are not assigned to this approval task")
	}
	if task.Status != models.ApprovalStatusPending {
		return errors.New("approval task is not pending")
	}

	now := time.Now()
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"status":       models.ApprovalStatusApproved,
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}

		switch task.DocumentType {
		case "PR":
			return tx.Model(&models.PurchaseRequisition{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.PRStatusApproved).Error
		case "PO":
			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.POStatusApproved).Error
		default:
			return nil
		}
	}); err != nil {
		return err
	}

	_ = req
	return nil
}

func (s *ApprovalService) Reject(taskID, approverID, actorEntityID uint, scopeType string, req RejectRequest) error {
	var task models.ApprovalTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		return errors.New("approval task not found")
	}
	if err := ensureEntityAccess(task.EntityID, actorEntityID, scopeType); err != nil {
		return err
	}
	if task.AssigneeID != approverID {
		return errors.New("you are not assigned to this approval task")
	}
	if task.Status != models.ApprovalStatusPending {
		return errors.New("approval task is not pending")
	}

	now := time.Now()
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"status":       models.ApprovalStatusRejected,
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}

		switch task.DocumentType {
		case "PR":
			return tx.Model(&models.PurchaseRequisition{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.PRStatusRejected).Error
		case "PO":
			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.POStatusRejected).Error
		default:
			return nil
		}
	})
}
