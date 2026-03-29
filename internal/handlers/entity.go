package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type EntityHandler struct {
	svc *services.EntityService
}

func NewEntityHandler(svc *services.EntityService) *EntityHandler {
	return &EntityHandler{svc: svc}
}

func (h *EntityHandler) List(c *gin.Context) {
	result, err := h.svc.List(c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "ENTITY_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, result)
}

func (h *EntityHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	entity, err := h.svc.GetByID(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "ENTITY_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, entity)
}

func (h *EntityHandler) Create(c *gin.Context) {
	var req services.CreateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	entity, err := h.svc.Create(req)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "ENTITY_CREATE_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, entity)
}
