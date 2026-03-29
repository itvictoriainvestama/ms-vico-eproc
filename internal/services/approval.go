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
			"notes":        req.Notes,
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}

		switch task.DocumentType {
		case "PR":
			if err := tx.Model(&models.PRApproval{}).
				Where("pr_id = ? AND approval_level = ?", task.DocumentID, task.ApprovalLevel).
				Updates(map[string]interface{}{
					"acted_by_user_id":     approverID,
					"on_behalf_of_user_id": task.OriginalUserID,
					"action":               "approve",
					"remarks":              req.Notes,
					"acted_at":             &now,
					"status_after_action":  models.PRStatusApproved,
				}).Error; err != nil {
				return err
			}
			return tx.Model(&models.PurchaseRequisition{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.PRStatusApproved).Error
		case "PO":
			if err := tx.Model(&models.POApproval{}).
				Where("po_id = ? AND approval_level = ?", task.DocumentID, task.ApprovalLevel).
				Updates(map[string]interface{}{
					"acted_by_user_id":     approverID,
					"on_behalf_of_user_id": task.OriginalUserID,
					"action":               "approve",
					"remarks":              req.Notes,
					"acted_at":             &now,
					"status_after_action":  models.POStatusApproved,
				}).Error; err != nil {
				return err
			}
			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.POStatusApproved).Error
		default:
			return nil
		}
	}); err != nil {
		return err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(task.EntityID),
		ModuleCode: "APPROVAL",
		ObjectType: task.DocumentType,
		ObjectID:   task.DocumentID,
		EventCode:  "APPROVAL_APPROVED",
		ActorType:  "internal_user",
		ActorID:    uintPtr(approverID),
		Action:     "approve",
		DataAfter:  req.Notes,
	})
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
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&task).Updates(map[string]interface{}{
			"status":       models.ApprovalStatusRejected,
			"notes":        req.Notes,
			"completed_at": &now,
		}).Error; err != nil {
			return err
		}

		switch task.DocumentType {
		case "PR":
			if err := tx.Model(&models.PRApproval{}).
				Where("pr_id = ? AND approval_level = ?", task.DocumentID, task.ApprovalLevel).
				Updates(map[string]interface{}{
					"acted_by_user_id":     approverID,
					"on_behalf_of_user_id": task.OriginalUserID,
					"action":               "reject",
					"remarks":              req.Notes,
					"acted_at":             &now,
					"status_after_action":  models.PRStatusRejected,
				}).Error; err != nil {
				return err
			}
			return tx.Model(&models.PurchaseRequisition{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.PRStatusRejected).Error
		case "PO":
			if err := tx.Model(&models.POApproval{}).
				Where("po_id = ? AND approval_level = ?", task.DocumentID, task.ApprovalLevel).
				Updates(map[string]interface{}{
					"acted_by_user_id":     approverID,
					"on_behalf_of_user_id": task.OriginalUserID,
					"action":               "reject",
					"remarks":              req.Notes,
					"acted_at":             &now,
					"status_after_action":  models.POStatusRejected,
				}).Error; err != nil {
				return err
			}
			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", task.DocumentID).
				Update("status", models.POStatusRejected).Error
		default:
			return nil
		}
	}); err != nil {
		return err
	}

	recordAuditLog(s.db, AuditEntry{
		EntityID:   uintPtr(task.EntityID),
		ModuleCode: "APPROVAL",
		ObjectType: task.DocumentType,
		ObjectID:   task.DocumentID,
		EventCode:  "APPROVAL_REJECTED",
		ActorType:  "internal_user",
		ActorID:    uintPtr(approverID),
		Action:     "reject",
		DataAfter:  req.Notes,
	})
	return nil
}
