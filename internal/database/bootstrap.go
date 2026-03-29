package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itvico/e-proc-api/internal/config"
	"github.com/itvico/e-proc-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ResetDatabase(cfg *config.Config) error {
	serverDB, err := sql.Open("mysql", cfg.Database.ServerDSN())
	if err != nil {
		return err
	}
	defer serverDB.Close()

	if _, err := serverDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", cfg.Database.Name)); err != nil {
		return err
	}

	if _, err := serverDB.Exec(fmt.Sprintf("CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.Name)); err != nil {
		return err
	}

	log.Printf("Database %s reset successfully", cfg.Database.Name)
	return nil
}

func EnsureDatabaseExists(cfg *config.Config) error {
	serverDB, err := sql.Open("mysql", cfg.Database.ServerDSN())
	if err != nil {
		return err
	}
	defer serverDB.Close()

	if _, err := serverDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.Database.Name)); err != nil {
		return err
	}

	return nil
}

func SeedMasterData(db *gorm.DB, cfg *config.Config) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.Bootstrap.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		entity := models.Entity{
			EntityCode:     cfg.Bootstrap.DefaultEntityCode,
			EntityName:     cfg.Bootstrap.DefaultEntityName,
			EntityType:     "head_office",
			Status:         "active",
			GovernanceMode: "entity_only",
		}
		if err := upsertEntity(tx, &entity); err != nil {
			return err
		}

		department := models.Department{
			EntityID: entity.ID,
			Code:     cfg.Bootstrap.DefaultDepartmentCode,
			Name:     cfg.Bootstrap.DefaultDepartmentName,
		}
		if err := upsertDepartment(tx, &department); err != nil {
			return err
		}

		roleSeeds := []models.Role{
			{RoleCode: "SUPER_ADMIN", RoleName: "Super Admin", PortalType: "internal", IsActive: true, Description: "Global internal administrator"},
			{RoleCode: "ENTITY_ADMIN", RoleName: "Entity Admin", PortalType: "internal", IsActive: true, Description: "Entity-level administrator"},
			{RoleCode: "PROCUREMENT_ADMIN", RoleName: "Procurement Admin", PortalType: "internal", IsActive: true, Description: "Procurement operations administrator"},
			{RoleCode: "REQUESTER", RoleName: "Requester", PortalType: "internal", IsActive: true, Description: "Purchase request creator"},
			{RoleCode: "APPROVER", RoleName: "Approver", PortalType: "internal", IsActive: true, Description: "Approval workflow actor"},
			{RoleCode: "VENDOR_ADMIN", RoleName: "Vendor Admin", PortalType: "vendor", IsActive: true, Description: "Vendor portal administrator"},
		}

		roleByCode := make(map[string]models.Role, len(roleSeeds))
		for i := range roleSeeds {
			role := roleSeeds[i]
			if err := upsertRole(tx, &role); err != nil {
				return err
			}
			roleByCode[role.RoleCode] = role
		}

		admin := models.User{
			EntityID:            entity.ID,
			DepartmentID:        &department.ID,
			FullName:            "System Administrator",
			Email:               "admin@eproc.local",
			Username:            "admin",
			PasswordHash:        string(passwordHash),
			Status:              "active",
			ForceChangePassword: true,
		}
		if err := upsertAdminUser(tx, &admin); err != nil {
			return err
		}

		adminRole := roleByCode["SUPER_ADMIN"]
		userRole := models.UserRole{
			UserID:    admin.ID,
			RoleID:    adminRole.ID,
			EntityID:  entity.ID,
			ScopeType: "cross_entity",
			IsPrimary: true,
			Status:    "active",
		}
		if err := upsertUserRole(tx, &userRole); err != nil {
			return err
		}

		if err := tx.Model(&models.User{}).Where("id = ?", admin.ID).Updates(map[string]interface{}{
			"primary_role_id": adminRole.ID,
			"department_id":   department.ID,
			"entity_id":       entity.ID,
			"status":          "active",
		}).Error; err != nil {
			return err
		}

		log.Printf("Seeded baseline data for entity=%s, admin_username=%s", entity.EntityCode, admin.Username)
		return nil
	})
}

func upsertEntity(tx *gorm.DB, entity *models.Entity) error {
	var existing models.Entity
	err := tx.Where("entity_code = ?", entity.EntityCode).First(&existing).Error
	if err == nil {
		entity.ID = existing.ID
		return tx.Model(&existing).Updates(map[string]interface{}{
			"entity_name":         entity.EntityName,
			"entity_type":         entity.EntityType,
			"status":              entity.Status,
			"approval_model_code": entity.ApprovalModelCode,
			"governance_mode":     entity.GovernanceMode,
		}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(entity).Error
}

func upsertDepartment(tx *gorm.DB, department *models.Department) error {
	var existing models.Department
	err := tx.Where("entity_id = ? AND code = ?", department.EntityID, department.Code).First(&existing).Error
	if err == nil {
		department.ID = existing.ID
		return tx.Model(&existing).Updates(map[string]interface{}{
			"name": department.Name,
		}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(department).Error
}

func upsertRole(tx *gorm.DB, role *models.Role) error {
	var existing models.Role
	err := tx.Where("role_code = ?", role.RoleCode).First(&existing).Error
	if err == nil {
		role.ID = existing.ID
		return tx.Model(&existing).Updates(map[string]interface{}{
			"role_name":   role.RoleName,
			"portal_type": role.PortalType,
			"is_active":   role.IsActive,
			"description": role.Description,
		}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(role).Error
}

func upsertAdminUser(tx *gorm.DB, user *models.User) error {
	var existing models.User
	err := tx.Where("username = ?", user.Username).First(&existing).Error
	if err == nil {
		user.ID = existing.ID
		return tx.Model(&existing).Updates(map[string]interface{}{
			"entity_id":             user.EntityID,
			"department_id":         user.DepartmentID,
			"full_name":             user.FullName,
			"email":                 user.Email,
			"password_hash":         user.PasswordHash,
			"status":                user.Status,
			"force_change_password": user.ForceChangePassword,
			"failed_login_count":    0,
			"locked_until":          nil,
		}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(user).Error
}

func upsertUserRole(tx *gorm.DB, userRole *models.UserRole) error {
	var existing models.UserRole
	err := tx.Where("user_id = ? AND role_id = ? AND entity_id = ?", userRole.UserID, userRole.RoleID, userRole.EntityID).First(&existing).Error
	if err == nil {
		userRole.ID = existing.ID
		return tx.Model(&existing).Updates(map[string]interface{}{
			"scope_type": userRole.ScopeType,
			"is_primary": userRole.IsPrimary,
			"status":     userRole.Status,
		}).Error
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return tx.Create(userRole).Error
}
