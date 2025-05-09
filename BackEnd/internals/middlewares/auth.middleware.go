package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/pkg"
)

// AuthMiddleware verifies JWT tokens from cookies and authorizes requests
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from cookie
		token, err := c.Cookie("authToken")
		if err != nil {
			c.JSON(401, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Verify token
		payload := pkg.NewPayload("", "")
		jwtErr := payload.VerifyToken(token)
		if jwtErr.Err != nil {
			c.JSON(401, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user ID for controllers to use
		c.Set("userID", payload.Id)
		c.Set("userRole", payload.Role)
		c.Next()
	}
}
