// Package main Portfolio Web Application API
//
// This is a RESTful API for the portfolio web application.
// It provides endpoints for managing portfolio content and user authentication.
//
// Terms Of Service: N/A
//
// Schemes: http, https
// Host: localhost
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/config"
	_ "github.com/Wildcard209/portfolio-webapplication/docs"
	"github.com/Wildcard209/portfolio-webapplication/routes"
	"github.com/Wildcard209/portfolio-webapplication/services"
	"github.com/gin-gonic/gin"
)

// @title Portfolio Web Application API
// @version 1.0
// @description This is a RESTful API for the portfolio web application
// @termsOfService N/A

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost
// @BasePath /api
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}
	defer cfg.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		if os.Getenv("GIN_MODE") == "release" {
			log.Fatal("JWT_SECRET environment variable must be set in production")
		} else {
			jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
			log.Println("Warning: Using default JWT secret. Please set JWT_SECRET environment variable.")
		}
	}

	authService := auth.NewAuthService(jwtSecret, 1*time.Hour)

	if cfg.DB != nil {
		adminService := services.NewAdminService(cfg.DB, authService)
		if err := adminService.InitializeAdminSystem(); err != nil {
			log.Fatalf("Failed to initialize admin system: %v", err)
		}

		adminService.StartMaintenanceTasks()
	} else {
		log.Println("Warning: Database not available, skipping admin system initialization")
	}

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	routes.SetupRoutes(r, cfg, authService)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server is running on port %s", cfg.Port)
		log.Printf("Swagger documentation available at http://localhost/api/swagger/index.html")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
