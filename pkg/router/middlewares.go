package router

import (
	"net/http"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/gin-gonic/gin"
)

type GinMiddleware = func(c *gin.Context)

func AuthMiddleware(repo *repositories.ApplicationRepository) GinMiddleware {
	return func(c *gin.Context) {
		token := c.GetHeader("x-api-token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Missing x-api-token header",
			})
			return
		}

		serviceToken, err := repo.ServiceTokens.GetServiceTokenByAccessKey(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
			return
		}
		if serviceToken == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid API token",
			})
			return
		}
		c.Set("serviceToken", serviceToken)
		c.Next()
	}
}

func AdminOnlyMiddleware(c *gin.Context) {
	serviceTokenIface, exists := c.Get("serviceToken")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}
	serviceToken, ok := serviceTokenIface.(*models.ServiceToken)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}
	if serviceToken.TokenType&models.AdminToken != models.AdminToken {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
			Code:    http.StatusForbidden,
			Message: "Forbidden",
		})
		return
	}
	c.Next()
}
