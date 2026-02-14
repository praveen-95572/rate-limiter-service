package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/praveen-95572/rate-limiter-service/internal/adaptor/persistance"
	"github.com/praveen-95572/rate-limiter-service/internal/limiter"
)

func main() {
	ctx := context.Background()

	redisStore := persistance.NewRedisStore("localhost:6379")

	if err := redisStore.Ping(ctx); err != nil {
		log.Fatalf("Redis not available: %v", err)
	}

	// tokenBucket := limiter.NewTokenBucket(redisStore, 2, 1)
	tokenBucket := limiter.NewSlidingWindow(redisStore, 2, 1)
	r := gin.Default()
	r.Use(limiter.RateLimitMiddleware(tokenBucket))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	log.Println("Server started on :8080")
	r.Run(":8080")
}
