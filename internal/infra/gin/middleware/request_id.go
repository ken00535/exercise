package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func requestIDGenerator() string {
	return uuid.New().String()
}

// RequestID set X-Request-ID to request and response header
func RequestID(c *gin.Context) {
	req := c.Request
	rid := req.Header.Get("X-Request-ID")
	if rid == "" {
		rid = requestIDGenerator()
		req.Header.Set("X-Request-ID", rid)
	}
	c.Header("X-Request-ID", rid)

	c.Next()
}
