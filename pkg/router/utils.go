package router

import "github.com/gin-gonic/gin"

func writeError(c *gin.Context, code int, msg string) {
	c.JSON(code, ErrorResponse{
		Code:    code,
		Message: msg,
	})
}
