package router

import (
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("x-api-token")
	if token == "" {
		c.AbortWithStatusJSON(401, ErrorResponse{
			Code:    401,
			Message: "Missing x-api-token header",
		})
		return
	}

	serviceToken, err := applicationRepository.ServiceTokens.GetServiceTokenByAccessKey(c.Request.Context(), token)
	if err != nil {
		c.AbortWithStatusJSON(500, ErrorResponse{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}
	if serviceToken == nil {
		c.AbortWithStatusJSON(401, ErrorResponse{
			Code:    401,
			Message: "Invalid API token",
		})
		return
	}
	c.Set("serviceToken", serviceToken)
	c.Next()
}

func AdminOnlyMiddleware(c *gin.Context) {
	serviceTokenIface, exists := c.Get("serviceToken")
	if !exists {
		c.AbortWithStatusJSON(500, ErrorResponse{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}
	serviceToken, ok := serviceTokenIface.(*models.ServiceToken)
	if !ok {
		c.AbortWithStatusJSON(500, ErrorResponse{
			Code:    500,
			Message: "Internal server error",
		})
		return
	}
	if serviceToken.TokenType&models.AdminToken != models.AdminToken {
		c.AbortWithStatusJSON(403, ErrorResponse{
			Code:    403,
			Message: "Forbidden",
		})
		return
	}
	c.Next()
}
