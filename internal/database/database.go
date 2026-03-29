package database

import (
	"log"

	"github.com/itvico/e-proc-api/internal/config"
	"github.com/itvico/e-proc-api/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) (*gorm.DB, error) {
	logLevel := logger.Info
	if cfg.App.Env == "production" {
		logLevel = logger.Error
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logLevel),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	DB = db
	log.Println("Database connected successfully")
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Entity{},
		&models.Department{},
		&models.Role{},
		&models.User{},
		&models.UserRole{},
		&models.DelegateApprover{},
		&models.Vendor{},
		&models.VendorUser{},
		&models.VendorBlacklist{},
		&models.ReferencePrice{},
		&models.Budget{},
		&models.ProcurementPolicyRule{},
		&models.ApprovalModel{},
		&models.ApprovalMatrix{},
		&models.PurchaseRequisition{},
		&models.PRItem{},
		&models.PRAttachment{},
		&models.PRApproval{},
		&models.RFQ{},
		&models.RFQVendor{},
		&models.VendorBid{},
		&models.BidItem{},
		&models.VendorEvaluation{},
		&models.BAFORound{},
		&models.VendorSelection{},
		&models.DirectAppointment{},
		&models.PurchaseOrder{},
		&models.POItem{},
		&models.POApproval{},
		&models.VendorConfirmation{},
		&models.ApprovalTask{},
		&models.Notification{},
		&models.AuditLog{},
		&models.AppLog{},
	)
}
