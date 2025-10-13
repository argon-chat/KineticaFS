package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	port   int
}

func (r *router) Run(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	setupDashboard(r)
	getRoutes(r)

	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", r.port),
		Handler: r.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	log.Printf("Listening and serving HTTP on %s\n", srv.Addr)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Printf("Server on port %s stopped gracefully", srv.Addr)
	}

	return nil
}

func NewRouter(repo *repositories.ApplicationRepository, port int) *router {
	return &router{
		engine: gin.Default(),
		repo:   repo,
		port:   port,
	}
}

func RunRouter(ctx context.Context, wg *sync.WaitGroup, port int, repo *repositories.ApplicationRepository) {
	defer wg.Done()

	router := NewRouter(repo, port)
	setupDashboard(router)
	getRoutes(router)

	router.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	log.Printf("Listening and serving HTTP on %s\n", srv.Addr)
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Printf("Server on port %s stopped gracefully", srv.Addr)
	}
}

func setupDashboard(router *router) {
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

func addV1Routes(router *router, v1 *gin.RouterGroup) {
	AddServiceTokenRoutes(router, v1)
	AddBucketsRoutes(router, v1)
	AddFileRoutes(v1)
	AddFileBlobRoutes(v1)
}
