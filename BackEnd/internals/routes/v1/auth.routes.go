package v1

import (
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/internals/handlers"
	"github.com/redha28/blogku/internals/middlewares"
)

func SetupAuthRoutes(router *gin.RouterGroup, db *sql.DB) {
	authController := handlers.NewAuthController(db)

	auth := router.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/logout", authController.Logout)

		// Private admin creation route - using a special API key
		// Should only be used internally or by a super admin
		admin := auth.Group("/admin")
		admin.Use(func(c *gin.Context) {
			apiKey := c.GetHeader("X-API-Key")
			if apiKey != os.Getenv("ADMIN_API_KEY") {
				c.JSON(403, gin.H{"error": "Forbidden"})
				c.Abort()
				return
			}
			c.Next()
		})
		{
			admin.POST("/create", authController.CreateAdmin)
		}

		// Route that requires authentication
		profile := auth.Group("/profile")
		profile.Use(middlewares.AuthMiddleware())
		{
			// Add profile routes here if needed
		}
	}
}
