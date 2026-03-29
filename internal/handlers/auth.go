package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itvico/e-proc-api/internal/httpapi"
	"github.com/itvico/e-proc-api/internal/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpapi.RespondError(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", []httpapi.ValidationError{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	resp, err := h.svc.Login(req)
	if err != nil {
		httpapi.RespondError(c, http.StatusUnauthorized, err.Error(), "AUTH_LOGIN_FAILED", nil)
		return
	}

	httpapi.RespondOK(c, resp)
}

func (h *AuthHandler) Me(c *gin.Context) {
	httpapi.RespondOK(c, gin.H{
		"user_id":      c.GetUint("user_id"),
		"entity_id":    c.GetUint("entity_id"),
		"username":     c.GetString("username"),
		"role_code":    c.GetString("role_code"),
		"role_name":    c.GetString("role_name"),
		"scope_type":   c.GetString("scope_type"),
		"subject_type": c.GetString("subject_type"),
	})
}
