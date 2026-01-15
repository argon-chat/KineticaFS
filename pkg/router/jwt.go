package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// JWTClaims represents the claims extracted from JWT token
type JWTClaims struct {
	Sub     string  `json:"sub"`     // GUID - user identifier
	SpaceId *string `json:"spaceId"` // GUID nullable - optional space identifier
	Type    string  `json:"type"`    // string - scope/type of upload (avatar, profileHeader, stickers, etc.)
	jwt.RegisteredClaims
}

// JWTAuthMiddleware validates JWT tokens from Authorization header
// It extracts sub (GUID), spaceId (nullable GUID), and type (string) fields
func JWTAuthMiddleware() GinMiddleware {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Missing Authorization header",
			})
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid Authorization header format. Expected: Bearer <token>",
			})
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Get JWT secret from configuration
			secret := viper.GetString("jwt_secret")
			if secret == "" {
				return nil, fmt.Errorf("JWT secret not configured")
			}

			return []byte(secret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid JWT token: " + err.Error(),
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid JWT token",
			})
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid JWT claims",
			})
			return
		}

		// Validate required fields
		if claims.Sub == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Missing 'sub' claim in JWT token",
			})
			return
		}

		// Validate that sub is a valid GUID
		if _, err := uuid.Parse(claims.Sub); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid 'sub' claim: must be a valid GUID",
			})
			return
		}

		// Validate spaceId if present
		if claims.SpaceId != nil && *claims.SpaceId != "" {
			if _, err := uuid.Parse(*claims.SpaceId); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
					Code:    http.StatusUnauthorized,
					Message: "Invalid 'spaceId' claim: must be a valid GUID",
				})
				return
			}
		}

		if claims.Type == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Missing 'type' claim in JWT token",
			})
			return
		}

		// Store JWT claims in context for later use
		c.Set("jwtClaims", claims)
		c.Next()
	}
}

// CombinedAuthMiddleware accepts either service token (x-api-token) or JWT token (Authorization Bearer)
func CombinedAuthMiddleware(repo *repositories.ApplicationRepository) GinMiddleware {
	return func(c *gin.Context) {
		// Check for x-api-token header first (service token)
		serviceTokenHeader := c.GetHeader("x-api-token")
		if serviceTokenHeader != "" {
			// Use existing service token validation
			serviceToken, err := repo.ServiceTokens.GetServiceTokenByAccessKey(c.Request.Context(), serviceTokenHeader)
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
			return
		}

		// Check for Authorization header (JWT token)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Use JWT validation
			JWTAuthMiddleware()(c)
			return
		}

		// No authentication provided
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Missing authentication: provide either x-api-token or Authorization header",
		})
	}
}
