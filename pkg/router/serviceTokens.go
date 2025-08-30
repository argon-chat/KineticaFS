package router

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response for the API
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"error message"`
}

// AddServiceTokenRoutes sets up the service token endpoints.
func AddServiceTokenRoutes(v1 *gin.RouterGroup) {
	st := v1.Group("/st")
	st.POST("/", CreateServiceTokenHandler)
	st.GET("/:id", GetServiceTokenHandler)
	st.DELETE("/:id", DeleteServiceTokenHandler)
}

// CreateServiceTokenHandler creates a new service token
// @Summary Create service token
// @Description Create a new service token. Only one admin token can exist. Only admin can create/delete other tokens. Admin token cannot be deleted.
// @Tags service-tokens
// @Accept json
// @Produce json
// @Param token body models.ServiceToken true "Service Token"
// @Success 201 {object} models.ServiceToken
// @Failure 400 {object} router.ErrorResponse
// @Router /v1/st/ [post]
func CreateServiceTokenHandler(c *gin.Context) {
	// Implementation goes here
}

// GetServiceTokenHandler gets a service token by ID
// @Summary Get service token
// @Description Get a service token by ID
// @Tags service-tokens
// @Produce json
// @Param id path string true "Token ID"
// @Success 200 {object} models.ServiceToken
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/st/{id} [get]
func GetServiceTokenHandler(c *gin.Context) {
	// Implementation goes here
}

// DeleteServiceTokenHandler deletes a service token by ID
// @Summary Delete service token
// @Description Delete a service token by ID. Only admin can delete other tokens. Admin token cannot be deleted.
// @Tags service-tokens
// @Param id path string true "Token ID"
// @Success 204 {object} nil
// @Failure 403 {object} router.ErrorResponse
// @Failure 404 {object} router.ErrorResponse
// @Router /v1/st/{id} [delete]
func DeleteServiceTokenHandler(c *gin.Context) {
	// Implementation goes here
}
