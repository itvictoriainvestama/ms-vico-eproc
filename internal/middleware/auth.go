package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/itvico/e-proc-api/internal/config"
	"github.com/itvico/e-proc-api/internal/httpapi"
)

type Claims struct {
	UserID      uint   `json:"user_id"`
	VendorID    uint   `json:"vendor_id"`
	EntityID    uint   `json:"entity_id"`
	Username    string `json:"username"`
	RoleCode    string `json:"role_code"`
	RoleName    string `json:"role_name"`
	ScopeType   string `json:"scope_type"`
	SubjectType string `json:"subject_type"`
	PortalType  string `json:"portal_type"`
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httpapi.RespondError(c, http.StatusUnauthorized, "Authorization header required", "AUTH_REQUIRED", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			httpapi.RespondError(c, http.StatusUnauthorized, "Invalid authorization format", "AUTH_INVALID_FORMAT", nil)
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			httpapi.RespondError(c, http.StatusUnauthorized, "Invalid or expired token", "AUTH_INVALID_TOKEN", nil)
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("vendor_id", claims.VendorID)
		c.Set("entity_id", claims.EntityID)
		c.Set("username", claims.Username)
		c.Set("role_code", claims.RoleCode)
		c.Set("role_name", claims.RoleName)
		c.Set("scope_type", claims.ScopeType)
		c.Set("subject_type", claims.SubjectType)
		c.Set("portal_type", claims.PortalType)
		c.Next()
	}
}

func RequireRole(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got := c.GetString("role_code")
		for _, roleCode := range roleCodes {
			if got == roleCode {
				c.Next()
				return
			}
		}

		httpapi.RespondError(c, http.StatusForbidden, "Insufficient permissions", "AUTH_FORBIDDEN", nil)
		c.Abort()
	}
}

func RequireSubjectType(subjectTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		got := c.GetString("subject_type")
		for _, subjectType := range subjectTypes {
			if got == subjectType {
				c.Next()
				return
			}
		}

		httpapi.RespondError(c, http.StatusForbidden, "Insufficient permissions", "AUTH_FORBIDDEN", nil)
		c.Abort()
	}
}
