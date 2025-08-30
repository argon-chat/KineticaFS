package router

import (
	"fmt"

	"github.com/argon-chat/KineticaFS/pkg/migrator"
	"github.com/argon-chat/KineticaFS/pkg/models"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run(port int) {
	getRoutes()
	migrator.MigrationTypes = []models.ApplicationRecord{
		models.ServiceToken{},
	}
	migrator.Migrate()
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
}

func getRoutes() {
	v1 := router.Group("/v1")
	addV1Routes(v1)
}

func addV1Routes(v1 *gin.RouterGroup) {
	addServiceTokenRoutes(v1)
	addBucketsRoutes(v1)
	addFileRoutes(v1)
}
