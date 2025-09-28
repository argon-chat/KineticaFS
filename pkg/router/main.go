package router

import (
	"fmt"

	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var router = gin.Default()

var applicationRepository *repositories.ApplicationRepository

func Run(port int, repo *repositories.ApplicationRepository) {
	initializeRepo(repo)
	getRoutes()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func getRoutes() {
	v1 := router.Group("/v1")
	addV1Routes(v1)
}

func initializeRepo(repo *repositories.ApplicationRepository) {
	applicationRepository = repo
	applicationRepository.ServiceTokens.CreateIndices()
	applicationRepository.Buckets.CreateIndices()
	applicationRepository.Files.CreateIndices()
}

func addV1Routes(v1 *gin.RouterGroup) {
	AddServiceTokenRoutes(v1)
	AddBucketsRoutes(v1)
	AddFileRoutes(v1)
	AddFileBlobRoutes(v1)
}
