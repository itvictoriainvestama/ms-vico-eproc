package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type UserHandler struct {
	svc *services.UserService
}

func NewUserHandler(svc *services.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) List(c *gin.Context) {
	entityID, _ := strconv.ParseUint(c.Query("entity_id"), 10, 64)
	params := services.UserListParams{
		EntityID: uint(entityID),
		Status:   c.Query("status"),
	}

	result, err := h.svc.List(c.GetUint("entity_id"), c.GetString("scope_type"), params)
	if err != nil {
		httpapi.RespondError(c, http.StatusInternalServerError, err.Error(), "USER_LIST_FAILED", nil)
		return
	}
	httpapi.RespondOK(c, result)
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Invalid id", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "id", Message: "must be a positive integer"},
		})
		return
	}

	user, err := h.svc.GetByID(uint(id), c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusNotFound, err.Error(), "USER_NOT_FOUND", nil)
		return
	}
	httpapi.RespondOK(c, user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req services.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	user, err := h.svc.Create(req, c.GetUint("entity_id"), c.GetString("scope_type"))
	if err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, err.Error(), "USER_CREATE_FAILED", nil)
		return
	}
	httpapi.RespondCreated(c, user)
}
