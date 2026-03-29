package router

import (
	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/config"
	"github.com/itvico/e-proc-api/internal/handlers"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/middleware"
)

type Handlers struct {
	Auth     *handlers.AuthHandler
	PR       *handlers.PRHandler
	RFQ      *handlers.RFQHandler
	PO       *handlers.POHandler
	Vendor   *handlers.VendorHandler
	Approval *handlers.ApprovalHandler
	Entity   *handlers.EntityHandler
	User     *handlers.UserHandler
}

func New(cfg *config.Config, h Handlers) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.RequestContext())
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		httpapi.RespondOK(c, gin.H{
			"status":  "ok",
			"service": "e-proc-api",
			"version": "phase-1-foundation",
		})
	})

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Auth.Login)
		}

		internal := v1.Group("/internal")
		internal.Use(middleware.AuthMiddleware(cfg))
		{
			internal.GET("/auth/me", h.Auth.Me)

			purchaseRequests := internal.Group("/purchase-requests")
			purchaseRequests.Use(middleware.RequireRole("SUPER_ADMIN", "ENTITY_ADMIN", "PROCUREMENT_ADMIN", "REQUESTER"))
			{
				purchaseRequests.GET("", h.PR.List)
				purchaseRequests.POST("", h.PR.Create)
				purchaseRequests.GET("/:id", h.PR.Get)
				purchaseRequests.POST("/:id/submit", h.PR.Submit)
			}

			rfqs := internal.Group("/rfqs")
			rfqs.Use(middleware.RequireRole("SUPER_ADMIN", "ENTITY_ADMIN", "PROCUREMENT_ADMIN"))
			{
				rfqs.GET("", h.RFQ.List)
				rfqs.POST("", h.RFQ.Create)
				rfqs.GET("/:id", h.RFQ.Get)
				rfqs.PATCH("/:id/status", h.RFQ.UpdateStatus)
			}

			purchaseOrders := internal.Group("/purchase-orders")
			purchaseOrders.Use(middleware.RequireRole("SUPER_ADMIN", "ENTITY_ADMIN", "PROCUREMENT_ADMIN"))
			{
				purchaseOrders.GET("", h.PO.List)
				purchaseOrders.POST("", h.PO.Create)
				purchaseOrders.GET("/:id", h.PO.Get)
				purchaseOrders.PATCH("/:id/status", h.PO.UpdateStatus)
			}

			vendors := internal.Group("/vendors")
			vendors.Use(middleware.RequireRole("SUPER_ADMIN", "ENTITY_ADMIN", "PROCUREMENT_ADMIN"))
			{
				vendors.GET("", h.Vendor.List)
				vendors.POST("", h.Vendor.Create)
				vendors.GET("/:id", h.Vendor.Get)
				vendors.PUT("/:id", h.Vendor.Update)
			}

			approvals := internal.Group("/approvals")
			approvals.Use(middleware.RequireRole("SUPER_ADMIN", "ENTITY_ADMIN", "PROCUREMENT_ADMIN", "APPROVER"))
			{
				approvals.GET("/tasks", h.Approval.MyTasks)
				approvals.POST("/tasks/:id/approve", h.Approval.Approve)
				approvals.POST("/tasks/:id/reject", h.Approval.Reject)
			}
		}

		vendor := v1.Group("/vendor")
		vendor.Use(middleware.AuthMiddleware(cfg))
		{
			vendor.GET("/health", func(c *gin.Context) {
				httpapi.RespondOK(c, gin.H{"portal": "vendor"})
			})
		}

		files := v1.Group("/files")
		files.Use(middleware.AuthMiddleware(cfg))
		{
			files.GET("/health", func(c *gin.Context) {
				httpapi.RespondOK(c, gin.H{"module": "files"})
			})
		}

		reports := v1.Group("/reports")
		reports.Use(middleware.AuthMiddleware(cfg))
		{
			reports.GET("/health", func(c *gin.Context) {
				httpapi.RespondOK(c, gin.H{"module": "reports"})
			})
		}

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))
		admin.Use(middleware.RequireRole("SUPER_ADMIN", "ENTITY_ADMIN"))
		{
			admin.GET("/health", func(c *gin.Context) {
				httpapi.RespondOK(c, gin.H{"module": "admin"})
			})

			admin.GET("/entities", h.Entity.List)
			admin.GET("/entities/:id", h.Entity.Get)

			admin.GET("/users", h.User.List)
			admin.POST("/users", h.User.Create)
			admin.GET("/users/:id", h.User.Get)
		}

		superAdmin := admin.Group("")
		superAdmin.Use(middleware.RequireRole("SUPER_ADMIN"))
		{
			superAdmin.POST("/entities", h.Entity.Create)
		}
	}

	return r
}
