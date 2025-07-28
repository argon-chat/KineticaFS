package router

import (
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func Run() {
	getRoutes()
	router.Run(":3000")
}

func getRoutes() {
	v1 := router.Group("/v1")
	addFileRoutes(v1)
}
