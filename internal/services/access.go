package services

import (
	"errors"

	"gorm.io/gorm"
)

const ScopeCrossEntity = "cross_entity"

func applyEntityScope(query *gorm.DB, entityColumn string, actorEntityID uint, scopeType string) *gorm.DB {
	if scopeType == ScopeCrossEntity || actorEntityID == 0 {
		return query
	}
	return query.Where(entityColumn+" = ?", actorEntityID)
}

func ensureEntityAccess(resourceEntityID, actorEntityID uint, scopeType string) error {
	if scopeType == ScopeCrossEntity || actorEntityID == 0 || resourceEntityID == actorEntityID {
		return nil
	}
	return errors.New("resource is outside your entity scope")
}
