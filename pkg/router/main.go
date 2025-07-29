package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run(port int) {
	getRoutes()
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
