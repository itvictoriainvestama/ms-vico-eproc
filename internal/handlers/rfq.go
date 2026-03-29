package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type RFQHandler struct {
	svc *services.RFQService
}

func NewRFQHandler(svc *services.RFQService) *RFQHandler {
	return &RFQHandler{svc: svc}
}

func (h *RFQHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	rfqs, total, err := h.svc.List(page, pageSize, c.Query("status"), c.GetUint("entity_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "RFQ_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"items": rfqs, "total": total, "page": page, "page_size": pageSize})
}

func (h *RFQHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	rfq, err := h.svc.GetByIDScoped(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "RFQ_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, rfq)
}

func (h *RFQHandler) Create(c *gin.Context) {
	var req services.CreateRFQRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	rfq, err := h.svc.Create(req, c.GetUint("entity_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "RFQ_CREATE_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, rfq)
}

func (h *RFQHandler) UpdateStatus(c *gin.Context) {
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

	rfq, err := h.svc.UpdateStatus(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"), body.Status)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "RFQ_STATUS_UPDATE_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, rfq)
}

func (h *RFQHandler) VendorList(c *gin.Context) {
	rfqs, err := h.svc.ListVendorTenders(c.GetUint("vendor_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "VENDOR_TENDER_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"items": rfqs})
}

func (h *RFQHandler) VendorGet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	rfq, err := h.svc.GetVendorTender(uint(id), c.GetUint("vendor_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "VENDOR_TENDER_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, rfq)
}

func (h *RFQHandler) VendorSubmitQuotation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var req services.VendorQuotationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	bid, err := h.svc.SubmitQuotation(uint(id), c.GetUint("vendor_id"), c.GetUint("user_id"), req)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "VENDOR_QUOTATION_SUBMIT_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, bid)
}
