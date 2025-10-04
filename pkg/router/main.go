package router

import (
	"context"
	"fmt"

	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ErrorResponse represents a standard error response for the API
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"error message"`
}

type router struct {
	engine *gin.Engine
	repo   *repositories.ApplicationRepository
}

func NewRouter(repo *repositories.ApplicationRepository) *router {
	return &router{
		engine: gin.Default(),
		repo:   repo,
	}
}

func Run(ctx context.Context, port int, repo *repositories.ApplicationRepository) {
	router := NewRouter(repo)
	initializeRepo(ctx, repo)
	setupDashboard(router, ctx)
	getRoutes(router)

	router.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.engine.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func setupDashboard(router *router, ctx context.Context) {
	dashboardPath := viper.GetString("front-end-path")
	router.engine.GET("/", func(c *gin.Context) {
		if c.Query("swagger") == "true" {
			c.Redirect(302, "/swagger/index.html")
			return
		}
		c.File(dashboardPath + "/index.html")
	})
	router.engine.Static("/assets", dashboardPath+"/assets")
	router.engine.GET("/favicon.ico", func(c *gin.Context) {
		c.File(dashboardPath + "/favicon.ico")
	})
}

func getRoutes(router *router) {
	v1 := router.engine.Group("/v1")
	addV1Routes(router, v1)
}

func initializeRepo(ctx context.Context, repo *repositories.ApplicationRepository) {
	repo.ServiceTokens.CreateIndices(ctx)
	repo.Buckets.CreateIndices(ctx)
	repo.Files.CreateIndices(ctx)
}

func addV1Routes(router *router, v1 *gin.RouterGroup) {
	AddServiceTokenRoutes(router, v1)
	AddBucketsRoutes(router, v1)
	AddFileRoutes(v1)
	AddFileBlobRoutes(v1)
}
