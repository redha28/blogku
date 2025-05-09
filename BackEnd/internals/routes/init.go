package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redha28/blogku/internals/middlewares"
	v1 "github.com/redha28/blogku/internals/routes/v1"
	"github.com/redis/go-redis/v9"
)

// InitRouter initializes all routes for the application
// @title Blog CMS API
// @version 1.0
// @description This is a Blog CMS API server.
// @BasePath /
func InitRouter(mySql *sql.DB, rdb *redis.Client) *gin.Engine {
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middlewares.CORSMiddleware())

	router.Static("/public", "./public")
	v1.InitRouter(router, mySql, rdb)
	return router
}
