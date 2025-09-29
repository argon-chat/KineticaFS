package router

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrorResponse represents a standard error response for the API
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"error message"`
}

type CreateServiceTokenRequestDto struct {
	Name string `json:"name" binding:"required" example:"my-token"`
}

// AddServiceTokenRoutes sets up the service token endpoints.
func AddServiceTokenRoutes(v1 *gin.RouterGroup) {
	st := v1.Group("/st")
	st.POST("/bootstrap", CreateAdminServiceTokenHandler)
	st.GET("/", AuthMiddleware, AdminOnlyMiddleware, ListAllServiceTokens)
	st.POST("/", AuthMiddleware, AdminOnlyMiddleware, CreateServiceTokenHandler)
	st.GET("/:id", AuthMiddleware, AdminOnlyMiddleware, GetServiceTokenHandler)
	st.DELETE("/:id", AuthMiddleware, AdminOnlyMiddleware, DeleteServiceTokenHandler)
}

// @Summary Bootstrap admin token
// @Description Create the initial admin service token. Only allowed if no admin token exists. Used for first-time setup.
// @Tags service-tokens
// @Produce json
// @Success 201 {object} models.ServiceToken
// @Failure 400 {object} router.ErrorResponse
// @Failure 409 {object} router.ErrorResponse "Admin token already exists"
// @Router /v1/st/bootstrap [post]
func CreateAdminServiceTokenHandler(c *gin.Context) {
	existingToken, err := applicationRepository.ServiceTokens.GetServiceTokenByName("admin")
	if err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("failed to check existing admin token: %v", err),
		})
		return
	}
	if existingToken != nil {
		c.JSON(409, ErrorResponse{
			Code:    409,
			Message: "Admin token already exists",
		})
		return
	}
	token := models.ServiceToken{
		Name:      "admin",
		AccessKey: fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
		TokenType: models.AdminToken | models.UserToken,
	}
	token.ID = uuid.NewString()
	token.CreatedAt = time.Now()
	token.UpdatedAt = token.CreatedAt
	err = applicationRepository.ServiceTokens.CreateServiceToken(&token)
	if err != nil {
		c.JSON(400, ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("failed to create admin token: %v", err),
		})
		return
	}
	c.JSON(201, token)
}

// @Summary List all service tokens
// @Description List all service tokens (admin only).
// @Tags service-tokens
// @Param x-api-token header string true "API Token"
// @Produce json
// @Success 200 {array} models.ServiceToken
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Router /v1/st/ [get]
func ListAllServiceTokens(c *gin.Context) {
	tokens, err := applicationRepository.ServiceTokens.GetAllServiceTokens()
	if err != nil {
		c.JSON(500, ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to list service tokens: %v", err),
		})
		return
	}
	c.JSON(200, tokens)
}

// CreateServiceTokenHandler creates a new service token
// @Summary Create service token
// @Description Create a new service token. Only one admin token can exist. Only admin can create/delete other tokens. Admin token cannot be deleted.
// @Tags service-tokens
// @Param x-api-token header string true "API Token"
// @Accept json
// @Produce json
// @Param request body CreateServiceTokenRequestDto true "Service Token Request"
// @Success 201 {object} models.ServiceToken
// @Failure 400 {object} router.ErrorResponse
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Router /v1/st/ [post]
func CreateServiceTokenHandler(c *gin.Context) {
	// Implementation goes here
}

// GetServiceTokenHandler gets a service token by ID
// @Summary Get service token
// @Description Get a service token by ID
// @Tags service-tokens
// @Param x-api-token header string true "API Token"
// @Produce json
// @Param id path string true "Token ID"
// @Success 200 {object} models.ServiceToken
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/st/{id} [get]
func GetServiceTokenHandler(c *gin.Context) {
	// Implementation goes here
}

// DeleteServiceTokenHandler deletes a service token by ID
// @Summary Delete service token
// @Description Delete a service token by ID. Only admin can delete other tokens. Admin token cannot be deleted.
// @Tags service-tokens
// @Param x-api-token header string true "API Token"
// @Param id path string true "Token ID"
// @Success 204 {object} nil
// @Failure 401 {object} router.ErrorResponse "Unauthorized"
// @Failure 403 {object} router.ErrorResponse "Forbidden - Admin only"
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/st/{id} [delete]
func DeleteServiceTokenHandler(c *gin.Context) {
	// Implementation goes here
}
