package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func AccessLog(c *gin.Context) {
	t := time.Now()

	c.Next()

	if strings.Contains(c.Request.URL.Path, "docs") {
		return
	}

	req := c.Request
	latency := time.Since(t)
	logger := log.With().
		Str("method", req.Method).
		Str("uri", req.RequestURI).
		Str("latency", latency.String()).
		Str("ip", req.Host).
		Int("status", c.Writer.Status()).Logger()
	if len(c.Errors) != 0 {
		logger.Error().Msgf("http access log: %v", req.RequestURI)
		return
	}
	logger.Info().Msgf("http access log: %v", req.RequestURI)
}
