package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LogginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		log.Printf(
			"Method: %s | Path: %s | Duration: %v | Status: %d",
			c.Request.Method,
			c.Request.URL.Path,
			duration,
			c.Writer.Status(),
		)
	}
}
