package main

import (
	"fmt"
	"log"
	"log/slog"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/saibaend/template-svc/docs"
	"github.com/saibaend/template-svc/internal/app"
	"github.com/saibaend/template-svc/internal/env"
	"github.com/saibaend/template-svc/pkg/conductor"
	srvHttp "github.com/saibaend/template-svc/pkg/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title						Template Service API
// @version					1.0
// @BasePath					/api
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Format: Bearer {access token}
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("loading .env file: %v", err)
	}

	config, err := env.New()
	if err != nil {
		log.Fatalf("prepare configs from env: %v", err)
	}

	logger := app.NewLogger(config.LogLevel)

	application, err := app.New(config, logger)
	if err != nil {
		log.Fatalf("initialize application: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			logger.Error("close application", slog.Any("error", err))
		}
	}()

	application.Router.GET("/actuator/health", application.HealthHandler())
	application.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := normalizeHTTPAddr(config.HttpPort)
	httpServer := srvHttp.New(
		addr,
		logger,
		srvHttp.WithRouter(application.Router),
		srvHttp.WithReadTimeout(config.HttpClientTimeout),
	)

	cdc := conductor.New(logger, httpServer)
	logger.Info("http server starting", slog.String("addr", addr))
	cdc.Shutdown(cdc.Run())
}

func normalizeHTTPAddr(port string) string {
	if strings.Contains(port, ":") {
		return port
	}

	return fmt.Sprintf(":%s", port)
}
