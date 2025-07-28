package router

import "github.com/gin-gonic/gin"

func addSessionRoutes(v1 *gin.RouterGroup) {
	session := v1.Group("/identity")

	session.POST("/", func(c *gin.Context) {
	})
	session.DELETE("/", authenticateUser, func(c *gin.Context) {
	})
}

func authenticateUser(c *gin.Context) {
	c.Next()
}

func authorizeUser(c *gin.Context) {
	c.Next()
}
