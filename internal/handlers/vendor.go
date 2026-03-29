package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type VendorHandler struct {
	svc *services.VendorService
}

func NewVendorHandler(svc *services.VendorService) *VendorHandler {
	return &VendorHandler{svc: svc}
}

func (h *VendorHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	activeOnly := c.Query("active") != "false"

	vendors, total, err := h.svc.List(page, pageSize, activeOnly)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "VENDOR_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"items": vendors, "total": total, "page": page, "page_size": pageSize})
}

func (h *VendorHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	vendor, err := h.svc.GetByID(uint(id))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "VENDOR_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, vendor)
}

func (h *VendorHandler) Create(c *gin.Context) {
	var req services.CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	vendor, err := h.svc.Create(req)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "VENDOR_CREATE_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, vendor)
}

func (h *VendorHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var req services.CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	vendor, err := h.svc.Update(uint(id), req)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "VENDOR_UPDATE_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, vendor)
}

func (h *VendorHandler) Blacklist(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var req services.VendorBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	vendor, err := h.svc.Blacklist(uint(id), c.GetUint("entity_id"), req)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "VENDOR_BLACKLIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, vendor)
}

func (h *VendorHandler) Unblacklist(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var req services.VendorBlacklistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	vendor, err := h.svc.Unblacklist(uint(id), c.GetUint("entity_id"), req)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "VENDOR_UNBLACKLIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, vendor)
}
