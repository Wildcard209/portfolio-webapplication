package routes

import (
	"net/http"
	"os"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/config"
	"github.com/Wildcard209/portfolio-webapplication/handlers"
	"github.com/Wildcard209/portfolio-webapplication/middleware"
	"github.com/Wildcard209/portfolio-webapplication/repository"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, authService *auth.AuthService) {
	// Add basic CORS for development
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.Use(middleware.LoggingMiddleware())

	api := r.Group("/api")
	{
		api.GET("/test", handlers.Hello)

		// Configure Swagger without explicit URL to use relative paths
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		if cfg.DB != nil {
			setupAdminRoutes(api, cfg, authService)
		}
	}
}

func setupAdminRoutes(api *gin.RouterGroup, cfg *config.Config, authService *auth.AuthService) {
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		adminToken = "1234"
	}

	adminRepo := repository.NewAdminRepository(cfg.DB)
	loginAttemptRepo := repository.NewLoginAttemptRepository(cfg.DB)

	adminHandler := handlers.NewAdminHandler(authService, adminRepo, loginAttemptRepo)

	adminGroup := api.Group("/:adminToken/admin")
	adminGroup.Use(middleware.AdminTokenValidationMiddleware(adminToken))
	adminGroup.Use(middleware.ValidateContentTypeMiddleware())
	{
		adminGroup.POST("/login",
			middleware.RateLimitMiddleware(handlers.GetRateLimiterForLogin()),
			adminHandler.Login,
		)

		protected := adminGroup.Group("")
		protected.Use(middleware.AuthMiddleware(authService, adminRepo))
		{
			protected.POST("/logout", adminHandler.Logout)
		}
	}
}
