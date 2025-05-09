package middlewares

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/pkg"
)

// LoggerMiddleware creates a middleware that logs HTTP requests
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Process request
		c.Next()

		// Calculate request duration
		duration := time.Since(startTime)
		status := c.Writer.Status()

		// Log request details
		pkg.LogHTTPRequest(method, path, clientIP, status, duration)
	}
}

// RecoveryMiddleware creates a middleware that recovers from panics
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				stackTrace := string(buf[:n])

				fields := map[string]interface{}{
					"stack": stackTrace,
					"path":  c.Request.URL.Path,
				}
				pkg.ErrorWithFields("PANIC RECOVERED", fmt.Errorf("%v", err), fields)

				// Return 500 error
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
