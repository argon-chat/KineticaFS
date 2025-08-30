package router

import (
	"github.com/gin-gonic/gin"
)

func addServiceTokenRoutes(v1 *gin.RouterGroup) {
	st := v1.Group("/st")

	st.POST("/", func(c *gin.Context) {
	})
	st.GET("/:id", func(c *gin.Context) {
	})
	st.DELETE("/:id", func(c *gin.Context) {
	})
}
