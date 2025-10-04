package router

import "github.com/gin-gonic/gin"

func writeError(c *gin.Context, code int, msg string) {
	c.JSON(400, ErrorResponse{
		Code:    400,
		Message: msg,
	})
}
