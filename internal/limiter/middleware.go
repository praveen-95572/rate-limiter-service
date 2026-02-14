package limiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/praveen-95572/rate-limiter-service/internal/core"
)

func RateLimitMiddleware(tb core.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			userID = "global"
		}

		key := "rate_limit:" + userID

		allowed, err := tb.Allow(c, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter failure",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
