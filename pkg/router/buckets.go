package router

import "github.com/gin-gonic/gin"

func addBucketsRoutes(v1 *gin.RouterGroup) {
	bucket := v1.Group("/bucket", authenticateUser, authorizeUser)

	bucket.POST("/", func(c *gin.Context) {
	})
	bucket.GET("/", func(c *gin.Context) {
	})
	bucket.GET("/:id", func(c *gin.Context) {
	})
	bucket.PATCH("/:id", func(c *gin.Context) {
	})
	bucket.DELETE("/:id", func(c *gin.Context) {
	})
}
