package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestTraceIDKey = "trace_id"

func RequestContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}

		c.Set(requestTraceIDKey, traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}
