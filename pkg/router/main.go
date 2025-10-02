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

var router = gin.Default()

var applicationRepository *repositories.ApplicationRepository

func Run(ctx context.Context, port int, repo *repositories.ApplicationRepository) {
	initializeRepo(ctx, repo)
	setupDashboard(ctx)
	getRoutes()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func setupDashboard(ctx context.Context) {
	dashboardPath := viper.GetString("front-end-path")
	router.GET("/", func(c *gin.Context) {
		if c.Query("swagger") == "true" {
			c.Redirect(302, "/swagger/index.html")
			return
		}
		c.File(dashboardPath + "/index.html")
	})
	router.Static("/assets", dashboardPath+"/assets")
}

func getRoutes() {
	v1 := router.Group("/v1")
	addV1Routes(v1)
}

func initializeRepo(ctx context.Context, repo *repositories.ApplicationRepository) {
	applicationRepository = repo
	applicationRepository.ServiceTokens.CreateIndices(ctx)
	applicationRepository.Buckets.CreateIndices(ctx)
	applicationRepository.Files.CreateIndices(ctx)
}

func addV1Routes(v1 *gin.RouterGroup) {
	AddServiceTokenRoutes(v1)
	AddBucketsRoutes(v1)
	AddFileRoutes(v1)
	AddFileBlobRoutes(v1)
}
