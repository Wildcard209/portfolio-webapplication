package routes

import (
	"os"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/config"
	"github.com/Wildcard209/portfolio-webapplication/handlers"
	"github.com/Wildcard209/portfolio-webapplication/middleware"
	"github.com/Wildcard209/portfolio-webapplication/repository"
	"github.com/Wildcard209/portfolio-webapplication/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, authService *auth.AuthService) {
	// Apply secure CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Apply security headers middleware
	r.Use(middleware.SecurityHeadersMiddleware())

	r.Use(middleware.LoggingMiddleware())

	api := r.Group("/api")
	{
		api.GET("/test", handlers.Hello)

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Asset routes (public access)
		if cfg.MinioClient != nil {
			setupAssetRoutes(api, cfg)
		}

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
	{
		adminGroup.POST("/login",
			middleware.RateLimitMiddleware(handlers.GetRateLimiterForLogin()),
			middleware.ValidateContentTypeMiddleware(),
			adminHandler.Login,
		)

		protected := adminGroup.Group("")
		protected.Use(middleware.AuthMiddleware(authService, adminRepo))
		{
			protected.POST("/logout", adminHandler.Logout)
		}
	}
}

func setupAssetRoutes(api *gin.RouterGroup, cfg *config.Config) {
	assetService := services.NewAssetService(cfg.MinioClient)
	assetHandler := handlers.NewAssetHandler(assetService)

	// Public asset routes
	assetsGroup := api.Group("/assets")
	{
		assetsGroup.GET("/hero-banner", assetHandler.GetHeroBanner)
		assetsGroup.GET("/info", assetHandler.GetAssetInfo)
	}

	// Protected asset routes (require admin authentication)
	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		adminToken = "1234"
	}

	adminAssetGroup := api.Group("/:adminToken/admin/assets")
	adminAssetGroup.Use(middleware.AdminTokenValidationMiddleware(adminToken))

	if cfg.DB != nil {
		adminRepo := repository.NewAdminRepository(cfg.DB)
		authService := auth.NewAuthService(os.Getenv("JWT_SECRET"), 1*time.Hour)

		protected := adminAssetGroup.Group("")
		protected.Use(middleware.AuthMiddleware(authService, adminRepo))
		{
			protected.POST("/hero-banner", assetHandler.UploadHeroBanner)
		}
	}
}
