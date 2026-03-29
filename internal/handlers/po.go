package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type POHandler struct {
	svc *services.POService
}

func NewPOHandler(svc *services.POService) *POHandler {
	return &POHandler{svc: svc}
}

func (h *POHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	vendorID, _ := strconv.ParseUint(c.Query("vendor_id"), 10, 64)

	pos, total, err := h.svc.List(page, pageSize, c.Query("status"), uint(vendorID), c.GetUint("entity_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "PO_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"items": pos, "total": total, "page": page, "page_size": pageSize})
}

func (h *POHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	po, err := h.svc.GetByIDScoped(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "PO_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, po)
}

func (h *POHandler) Create(c *gin.Context) {
	var req services.CreatePORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	po, err := h.svc.Create(req, c.GetUint("entity_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "PO_CREATE_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, po)
}

func (h *POHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "status", Message: err.Error()},
		})
		return
	}

	po, err := h.svc.UpdateStatus(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"), body.Status)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "PO_STATUS_UPDATE_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, po)
}
