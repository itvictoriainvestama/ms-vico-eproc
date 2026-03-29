package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type ApprovalHandler struct {
	svc *services.ApprovalService
}

func NewApprovalHandler(svc *services.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{svc: svc}
}

func (h *ApprovalHandler) MyTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tasks, total, err := h.svc.GetTasksByUser(c.GetUint("user_id"), c.GetUint("entity_id"), c.GetString("scope_type"), page, pageSize)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "APPROVAL_TASK_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"items": tasks, "total": total, "page": page, "page_size": pageSize})
}

func (h *ApprovalHandler) Approve(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var req services.ApproveRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.svc.Approve(uint(id), c.GetUint("user_id"), c.GetUint("entity_id"), c.GetString("scope_type"), req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "APPROVAL_APPROVE_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"task_id": id, "status": "approved"})
}

func (h *ApprovalHandler) Reject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	var req services.RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "notes", Message: err.Error()},
		})
		return
	}

	if err := h.svc.Reject(uint(id), c.GetUint("user_id"), c.GetUint("entity_id"), c.GetString("scope_type"), req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "APPROVAL_REJECT_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, gin.H{"task_id": id, "status": "rejected"})
}
