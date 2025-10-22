package router

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/argon-chat/KineticaFS/pkg/repositories"
	"github.com/gin-contrib/cors"
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
	ginRouter := gin.Default()
	var allowedOrigins []string
	var allowedHeaders []string
	if envOrigins := os.Getenv("KINETICAFS_CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		for _, origin := range strings.Split(envOrigins, ",") {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(origin))
		}
	} else {
		allowedOrigins = strings.Split(viper.GetString("cors-allowed-origins"), ",")
	}
	if envHeaders := os.Getenv("KINETICAFS_CORS_ALLOWED_HEADERS"); envHeaders != "" {
		for _, header := range strings.Split(envHeaders, ",") {
			allowedHeaders = append(allowedHeaders, strings.TrimSpace(header))
		}
	} else {
		allowedHeaders = strings.Split(viper.GetString("cors-allowed-headers"), ",")
	}

	for i, header := range allowedHeaders {
		allowedHeaders[i] = strings.TrimSpace(header)
	}
	fmt.Printf("%+v", allowedOrigins)
	fmt.Printf("%+v", allowedHeaders)
	corsConfig := cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: allowedHeaders,
		MaxAge:       24 * time.Hour,
	}
	ginRouter.Use(cors.New(corsConfig))
	return &router{
		engine: ginRouter,
		repo:   repo,
		port:   port,
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
	v1 := router.engine.Group("/api/v1")
	addV1Routes(router, v1)
}

func addV1Routes(router *router, v1 *gin.RouterGroup) {
	AddServiceTokenRoutes(router, v1)
	AddBucketsRoutes(router, v1)
	AddFileRoutes(router, v1)
	AddFileBlobRoutes(router, v1)
}
