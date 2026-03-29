package services

import (
	"errors"
	"time"

	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/gorm"
)

type AuditEntry struct {
	EntityID   *uint
	ModuleCode string
	ObjectType string
	ObjectID   uint
	EventCode  string
	ActorType  string
	ActorID    *uint
	Action     string
	DataBefore string
	DataAfter  string
}

func recordAuditLog(db *gorm.DB, entry AuditEntry) {
	if db == nil {
		return
	}

	actorType := entry.ActorType
	if actorType == "" {
		actorType = "internal_user"
	}

	log := models.AuditLog{
		EntityID:   entry.EntityID,
		ModuleCode: entry.ModuleCode,
		ObjectType: entry.ObjectType,
		ObjectID:   entry.ObjectID,
		EventCode:  entry.EventCode,
		ActorType:  actorType,
		ActorID:    entry.ActorID,
		Action:     entry.Action,
		DataBefore: entry.DataBefore,
		DataAfter:  entry.DataAfter,
	}

	_ = db.Create(&log).Error
}

func uintPtr(v uint) *uint {
	return &v
}

func resolveActiveDelegate(tx *gorm.DB, entityID, userID uint, at time.Time) (uint, *uint, error) {
	var delegate models.DelegateApprover
	err := tx.
		Where("entity_id = ? AND original_user_id = ? AND status = ? AND start_at <= ? AND end_at >= ?", entityID, userID, "active", at, at).
		Order("id DESC").
		First(&delegate).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userID, nil, nil
		}
		return 0, nil, err
	}

	return delegate.DelegateUserID, &delegate.OriginalUserID, nil
}

func resolveEntityApprover(tx *gorm.DB, entityID uint, roleCodes ...string) (*models.User, *uint, error) {
	var userRole models.UserRole
	err := tx.
		Preload("User").
		Preload("User.PrimaryRole").
		Preload("Role").
		Joins("JOIN users ON users.id = user_roles.user_id").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.entity_id = ? AND user_roles.status = ? AND user_roles.is_primary = ? AND users.status = ? AND roles.role_code IN ?", entityID, "active", true, "active", roleCodes).
		Order("user_roles.id ASC").
		First(&userRole).Error
	if err != nil {
		return nil, nil, err
	}

	assigneeID, originalUserID, err := resolveActiveDelegate(tx, entityID, userRole.UserID, time.Now())
	if err != nil {
		return nil, nil, err
	}

	if assigneeID == userRole.UserID {
		return &userRole.User, nil, nil
	}

	var delegateUser models.User
	if err := tx.Preload("PrimaryRole").First(&delegateUser, assigneeID).Error; err != nil {
		return nil, nil, err
	}

	return &delegateUser, originalUserID, nil
}

func primaryRoleCode(user *models.User) string {
	if user == nil || user.PrimaryRole == nil {
		return ""
	}
	return user.PrimaryRole.RoleCode
}
