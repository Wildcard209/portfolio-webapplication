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
	r.Use(middleware.CORSMiddleware())

	r.Use(middleware.SecurityHeadersMiddleware(cfg))

	r.Use(middleware.RequestBodySizeLimitMiddleware(middleware.GetRequestBodySizeLimit()))

	r.Use(middleware.LoggingMiddleware())

	r.Use(middleware.RateLimitViolationMiddleware())

	api := r.Group("/api")
	{
		api.GET("/test",
			middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitPublic, cfg.RateLimit),
			handlers.Hello,
		)

		api.POST("/csp-report", handlers.CSPReportHandler)

		api.GET("/health", handlers.HealthCheckHandler)

		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		if cfg.MinioClient != nil {
			setupAssetRoutes(api, cfg)
		}

		if cfg.DB != nil {
			setupAdminRoutes(api, cfg, authService)
		}
	}
}

func setupAdminRoutes(api *gin.RouterGroup, cfg *config.Config, authService *auth.AuthService) {
	adminRepo := repository.NewAdminRepository(cfg.DB)
	loginAttemptRepo := repository.NewLoginAttemptRepository(cfg.DB)

	adminHandler := handlers.NewAdminHandler(authService, adminRepo, loginAttemptRepo)

	adminGroup := api.Group("/admin")
	{
		adminGroup.POST("/login",
			middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitLogin, cfg.RateLimit),
			middleware.ValidateContentTypeMiddleware(),
			adminHandler.Login,
		)

		adminGroup.POST("/refresh",
			middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitRefresh, cfg.RateLimit),
			adminHandler.RefreshToken,
		)

		protected := adminGroup.Group("")
		protected.Use(middleware.AuthMiddleware(authService, adminRepo))
		{
			protected.POST("/logout",
				middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitAdmin, cfg.RateLimit),
				adminHandler.Logout,
			)
		}
	}
}

func setupAssetRoutes(api *gin.RouterGroup, cfg *config.Config) {
	assetService := services.NewAssetService(cfg.MinioClient)
	assetHandler := handlers.NewAssetHandler(assetService)

	assetsGroup := api.Group("/assets")
	{
		assetsGroup.GET("/hero-banner",
			middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitPublic, cfg.RateLimit),
			assetHandler.GetHeroBanner,
		)
		assetsGroup.GET("/info",
			middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitPublic, cfg.RateLimit),
			assetHandler.GetAssetInfo,
		)
	}

	adminAssetGroup := api.Group("/admin/assets")

	if cfg.DB != nil {
		adminRepo := repository.NewAdminRepository(cfg.DB)
		authService := auth.NewAuthService(os.Getenv("JWT_SECRET"), 1*time.Hour)

		protected := adminAssetGroup.Group("")
		protected.Use(middleware.AuthMiddleware(authService, adminRepo))
		protected.Use(middleware.FileUploadSizeLimitMiddleware(middleware.GetFileUploadSizeLimit()))
		{
			protected.POST("/hero-banner",
				middleware.RateLimitMiddlewareWithConfig(middleware.RateLimitUpload, cfg.RateLimit),
				assetHandler.UploadHeroBanner,
			)
		}
	}
}
