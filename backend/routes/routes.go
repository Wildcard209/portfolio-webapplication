package routes

import (
	"github.com/Wildcard209/portfolio-webapplication/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine) {
	// API routes group
	api := r.Group("/api")
	{
		api.GET("/test", handlers.Hello)
		// Configure Swagger to use the correct URL for API calls
		url := ginSwagger.URL("/api/swagger/doc.json") // The url pointing to API definition
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
}
