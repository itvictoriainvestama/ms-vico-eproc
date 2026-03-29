package main

import (
	"log"

	"github.com/itvico/e-proc-api/internal/config"
	"github.com/itvico/e-proc-api/internal/database"
	"github.com/itvico/e-proc-api/internal/handlers"
	"github.com/itvico/e-proc-api/internal/router"
	"github.com/itvico/e-proc-api/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	if cfg.Bootstrap.ResetDatabase {
		log.Println("Resetting database...")
		if err := database.ResetDatabase(cfg); err != nil {
			log.Fatalf("Database reset failed: %v", err)
		}
	}

	if err := database.EnsureDatabaseExists(cfg); err != nil {
		log.Fatalf("Database bootstrap failed: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if cfg.Bootstrap.Migrate {
		log.Println("Running database migrations...")
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed")
	}

	if cfg.Bootstrap.SeedMasterData {
		log.Println("Seeding baseline master data...")
		if err := database.SeedMasterData(db, cfg); err != nil {
			log.Fatalf("Seed failed: %v", err)
		}
		log.Println("Seed completed")
	}

	// Wire services
	authSvc := services.NewAuthService(db, cfg)
	prSvc := services.NewPRService(db)
	rfqSvc := services.NewRFQService(db)
	poSvc := services.NewPOService(db)
	vendorSvc := services.NewVendorService(db)
	approvalSvc := services.NewApprovalService(db)
	entitySvc := services.NewEntityService(db)
	userSvc := services.NewUserService(db)

	// Wire handlers
	h := router.Handlers{
		Auth:     handlers.NewAuthHandler(authSvc),
		PR:       handlers.NewPRHandler(prSvc),
		RFQ:      handlers.NewRFQHandler(rfqSvc),
		PO:       handlers.NewPOHandler(poSvc),
		Vendor:   handlers.NewVendorHandler(vendorSvc),
		Approval: handlers.NewApprovalHandler(approvalSvc),
		Entity:   handlers.NewEntityHandler(entitySvc),
		User:     handlers.NewUserHandler(userSvc),
	}

	r := router.New(cfg, h)

	addr := ":" + cfg.App.Port
	log.Printf("Server starting on %s (env: %s)", addr, cfg.App.Env)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
