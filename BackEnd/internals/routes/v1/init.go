package v1

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitRouter(router *gin.Engine, mySql *sql.DB, rdb *redis.Client) {
	v1 := router.Group("/api/v1")

	// Setup routes
	SetupAuthRoutes(v1, mySql)
	SetupBlogRoutes(v1, mySql, rdb)
}
