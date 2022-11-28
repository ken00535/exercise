package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(limiterRate int, bucket int, timeout int) gin.HandlerFunc {
	r := rate.Every(time.Duration(limiterRate) * time.Millisecond)
	l := rate.NewLimiter(r, bucket)
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, time.Duration(timeout)*time.Millisecond)
		defer cancel()

		if err := l.Wait(ctx); err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}
		c.Next()
	}
}
