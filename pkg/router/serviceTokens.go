package router

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateServiceTokenRequestDto struct {
	Name string `json:"name" binding:"required" example:"my-token"`
}

// AddServiceTokenRoutes sets up the service token endpoints.
func AddServiceTokenRoutes(router *router, v1 *gin.RouterGroup) {
	st := v1.Group("/st")
	st.GET("/first-run", router.FirstRunCheckHandler)
	st.POST("/bootstrap", router.CreateAdminServiceTokenHandler)
	st.GET("/", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.ListAllServiceTokens)
	st.POST("/", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.CreateServiceTokenHandler)
	st.GET("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.GetServiceTokenHandler)
	st.DELETE("/:id", AuthMiddleware(router.repo), AdminOnlyMiddleware, router.DeleteServiceTokenHandler)
}

// @Summary Check if admin token has already been created
// @Description Returns whether the admin token exists. Used to determine if setup is required.
// @Tags service-tokens
// @Produce json
// @Success 200 {object} map[string]bool "first_run: true if no admin token exists, false otherwise"
// @Failure 500 {object} router.ErrorResponse "Internal server error"
// @Router /v1/st/first-run [get]
// @Id FirstRunCheck
func (r *router) FirstRunCheckHandler(c *gin.Context) {
	ctx := c.Request.Context()
	existingToken, err := r.repo.ServiceTokens.GetServiceTokenByName(ctx, "admin")
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to check existing admin token: %v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"first_run": existingToken == nil})
}

// @Summary Bootstrap admin token
// @Description Create the initial admin service token. Only allowed if no admin token exists. Used for first-time setup.
// @Tags service-tokens
// @Produce json
// @Success 201 {object} models.ServiceToken
// @Failure 400 {object} router.ErrorResponse
// @Failure 409 {object} router.ErrorResponse "Admin token already exists"
// @Router /v1/st/bootstrap [post]
// @Id BootstrapAdminToken
func (r *router) CreateAdminServiceTokenHandler(c *gin.Context) {
	ctx := c.Request.Context()
	existingToken, err := r.repo.ServiceTokens.GetServiceTokenByName(ctx, "admin")
	if err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to check existing admin token: %v", err))
		return
	}
	if existingToken != nil {
		writeError(c, http.StatusConflict, "Admin token already exists")
		return
	}
	token := models.ServiceToken{
		Name:      "admin",
		AccessKey: fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
		TokenType: models.AdminToken | models.UserToken,
	}
	err = r.repo.ServiceTokens.CreateServiceToken(ctx, &token)
	if err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to create admin token: %v", err))
		return
	}
	c.JSON(http.StatusCreated, token)
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
// @Id ListAllServiceTokens
func (r *router) ListAllServiceTokens(c *gin.Context) {
	tokens, err := r.repo.ServiceTokens.GetAllServiceTokens(c.Request.Context())
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to list service tokens: %v", err))
		return
	}
	c.JSON(http.StatusOK, tokens)
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
// @Id CreateServiceToken
func (r *router) CreateServiceTokenHandler(c *gin.Context) {
	var req CreateServiceTokenRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}
	ctx := c.Request.Context()
	existingToken, err := r.repo.ServiceTokens.GetServiceTokenByName(ctx, req.Name)
	if err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to check existing token: %v", err))
		return
	}
	if existingToken != nil {
		writeError(c, http.StatusConflict, "Service token with the same name already exists")
		return
	}
	token := models.ServiceToken{
		Name:      req.Name,
		AccessKey: fmt.Sprintf("%x", sha256.Sum256([]byte(uuid.NewString()))),
		TokenType: models.UserToken,
	}
	err = r.repo.ServiceTokens.CreateServiceToken(ctx, &token)
	if err != nil {
		writeError(c, http.StatusBadRequest, fmt.Sprintf("failed to create service token: %v", err))
		return
	}
	c.JSON(http.StatusCreated, token)
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
// @Id GetServiceToken
func (r *router) GetServiceTokenHandler(c *gin.Context) {
	id := c.Param("id")
	token, err := r.repo.ServiceTokens.GetServiceTokenById(c.Request.Context(), id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get service token: %v", err))
		return
	}
	if token == nil {
		writeError(c, http.StatusNotFound, "Service token not found")
		return
	}
	c.JSON(200, token)
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
// @Id DeleteServiceToken
func (r *router) DeleteServiceTokenHandler(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	token, err := r.repo.ServiceTokens.GetServiceTokenById(ctx, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to get service token: %v", err))
		return
	}
	if token == nil {
		writeError(c, http.StatusNotFound, "Service token not found")
		return
	}
	if token.TokenType&models.AdminToken == models.AdminToken {
		writeError(c, http.StatusForbidden, "Cannot delete admin token")
		return
	}
	err = r.repo.ServiceTokens.RevokeServiceToken(ctx, id)
	if err != nil {
		writeError(c, http.StatusInternalServerError, fmt.Sprintf("failed to delete service token: %v", err))
		return
	}
	c.Status(http.StatusNoContent)
}
