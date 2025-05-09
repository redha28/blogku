package v1

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/internals/handlers"
	"github.com/redha28/blogku/internals/middlewares"
	"github.com/redis/go-redis/v9"
)

func SetupBlogRoutes(router *gin.RouterGroup, db *sql.DB, rdb *redis.Client) {
	blogController := handlers.NewBlogController(db, rdb)

	// Public routes
	router.GET("/blogs", blogController.GetAllBlogs)
	router.GET("/blogs/:slug", blogController.GetBlogBySlug)

	// Protected routes
	adminBlogs := router.Group("/admin/blogs")
	adminBlogs.Use(middlewares.AuthMiddleware())
	{
		adminBlogs.POST("", blogController.CreateBlog)
		adminBlogs.PATCH("/:id", blogController.UpdateBlog)
		adminBlogs.DELETE("/:id", blogController.DeleteBlog)
	}
}
