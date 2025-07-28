package router

import "github.com/gin-gonic/gin"

func addFileRoutes(v1 *gin.RouterGroup) {
	files := v1.Group("/file", authenticateUser, authorizeUser)

	files.POST("/", func(c *gin.Context) {
	})
	files.GET("/:id", func(c *gin.Context) {
	})
	files.PATCH("/:id", func(c *gin.Context) {
	})
	files.DELETE("/:id", func(c *gin.Context) {
	})
}
