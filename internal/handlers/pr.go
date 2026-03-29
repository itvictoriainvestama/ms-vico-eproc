package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type PRHandler struct {
	svc *services.PRService
}

func NewPRHandler(svc *services.PRService) *PRHandler {
	return &PRHandler{svc: svc}
}

func (h *PRHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	params := services.PRListParams{
		Page:           page,
		PageSize:       pageSize,
		Status:         c.Query("status"),
		EntityID:       c.GetUint("entity_id"),
		DepartmentCode: c.Query("department_code"),
	}

	result, err := h.svc.List(params)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "PR_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, result)
}

func (h *PRHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	pr, err := h.svc.GetByIDScoped(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "PR_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, pr)
}

func (h *PRHandler) Create(c *gin.Context) {
	var req services.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	pr, err := h.svc.Create(req, c.GetUint("user_id"), c.GetUint("entity_id"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "PR_CREATE_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, pr)
}

func (h *PRHandler) Submit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	pr, err := h.svc.Submit(uint(id), c.GetUint("user_id"), c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "PR_SUBMIT_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, pr)
}
