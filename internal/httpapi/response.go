package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	TraceID string `json:"trace_id"`
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

type ErrorResponse struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	ErrorCode string            `json:"error_code"`
	Errors    []ValidationError `json:"errors,omitempty"`
	Meta      Meta              `json:"meta"`
}

func Respond(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    Meta{TraceID: TraceID(c)},
	})
}

func RespondOK(c *gin.Context, data interface{}) {
	Respond(c, http.StatusOK, "OK", data)
}

func RespondCreated(c *gin.Context, data interface{}) {
	Respond(c, http.StatusCreated, "Created", data)
}

func RespondError(c *gin.Context, status int, message, errorCode string, errs []ValidationError) {
	c.JSON(status, ErrorResponse{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Errors:    errs,
		Meta:      Meta{TraceID: TraceID(c)},
	})
}

func TraceID(c *gin.Context) string {
	traceID := c.GetString("trace_id")
	if traceID == "" {
		return "unknown"
	}
	return traceID
}
