package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/saibaend/template-svc/docs"
	"github.com/saibaend/template-svc/internal/app/someModule/delivery"
	"github.com/saibaend/template-svc/internal/env"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
)

// @title						Example Swagger API
// @version					1.0
// @BasePath					/api
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Format: Bearer {userId} or Bearer {access token}
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("loading .env file: %v", err)
	}

	config, err := env.New()
	if err != nil {
		log.Fatalf("prepare configs from env: %v", err)
	}

	router := initRouter()

	httpSrv := &http.Server{
		Addr:         config.HttpPort,
		Handler:      router,
		ReadTimeout:  config.HttpClientTimeout,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Printf("http server listening on port %s\n", config.HttpPort)

	if err := httpSrv.ListenAndServe(); err != nil {
		panic(err)
	}

}

func initRouter() *gin.Engine {

	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// healthChecker
	router.GET("/actuator/health")

	//some module route attach
	delivery.AttachRoutes(router, nil)

	return router
}
