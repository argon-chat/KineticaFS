package router

import "github.com/gin-gonic/gin"

func addUserRoutes(v1 *gin.RouterGroup) {
	users := v1.Group("/users")

	users.POST("/", func(c *gin.Context) {
	})
	users.GET("/:id", func(c *gin.Context) {
	})
	users.PATCH("/:id", authenticateUser, authorizeUser, func(c *gin.Context) {
	})
	users.DELETE("/:id", authenticateUser, authorizeUser, func(c *gin.Context) {
	})
}
