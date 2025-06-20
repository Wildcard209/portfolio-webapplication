package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/repository"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimitMiddleware(rate limiter.Rate) gin.HandlerFunc {
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)

	return mgin.NewMiddleware(rateLimiter)
}

func AuthMiddleware(authService *auth.AuthService, adminRepo *repository.AdminRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from cookie first, then fallback to header for backward compatibility
		var tokenString string
		var err error

		// Check for access token in cookie
		accessToken, cookieErr := c.Cookie("access_token")
		if cookieErr == nil && accessToken != "" {
			tokenString = accessToken
		} else {
			// Fallback to Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
				c.Abort()
				return
			}

			tokenString, err = authService.ExtractTokenFromHeader(authHeader)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
				c.Abort()
				return
			}
		}

		// Validate access token
		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			// If access token is invalid/expired, try to refresh using refresh token
			refreshToken, refreshErr := c.Cookie("refresh_token")
			if refreshErr != nil || refreshToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			// Validate refresh token
			refreshClaims, refreshValidErr := authService.ValidateRefreshToken(refreshToken)
			if refreshValidErr != nil {
				// Clear invalid cookies
				isHttps := os.Getenv("HTTPS_MODE") == "true"
				c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
				c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired, please login again"})
				c.Abort()
				return
			}

			// Check if refresh token is still valid in database
			admin, adminErr := adminRepo.GetAdminByToken(refreshToken)
			if adminErr != nil || admin == nil {
				// Clear invalid cookies
				isHttps := os.Getenv("HTTPS_MODE") == "true"
				c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
				c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session has been revoked"})
				c.Abort()
				return
			}

			// Generate new token pair
			tokenPair, tokenErr := authService.GenerateTokenPair(refreshClaims.UserID, refreshClaims.Username)
			if tokenErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh session"})
				c.Abort()
				return
			}

			// Update refresh token in database
			if updateErr := adminRepo.UpdateAdminToken(refreshClaims.UserID, tokenPair.RefreshToken, tokenPair.RefreshExpiresAt); updateErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
				c.Abort()
				return
			}

			// Set new HTTP-only cookies
			isHttps := os.Getenv("HTTPS_MODE") == "true"
			c.SetCookie(
				"access_token",
				tokenPair.AccessToken,
				int(tokenPair.AccessExpiresAt.Sub(time.Now()).Seconds()),
				"/",
				"",
				isHttps, // Secure flag - true for HTTPS, false for HTTP
				true,    // HTTP-only
			)

			c.SetCookie(
				"refresh_token",
				tokenPair.RefreshToken,
				int(tokenPair.RefreshExpiresAt.Sub(time.Now()).Seconds()),
				"/",
				"",
				isHttps, // Secure flag - true for HTTPS, false for HTTP
				true,    // HTTP-only
			)

			// Use the new access token claims
			claims = &auth.CustomClaims{
				UserID:    refreshClaims.UserID,
				Username:  refreshClaims.Username,
				TokenType: "access",
			}

			// Set admin context from database
			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("admin", admin)

			c.Next()
			return
		}

		// For valid access tokens, verify against database
		admin, err := adminRepo.GetAdminByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify token"})
			c.Abort()
			return
		}

		if admin == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("admin", admin)

		c.Next()
	}
}

func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	// Define allowed origins from environment variables with secure defaults
	allowedOrigins := getAllowedOrigins()

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is in our allowed list
		isOriginAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isOriginAllowed = true
				break
			}
		}

		// Only set CORS headers if origin is allowed
		if isOriginAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		} else {
			// Log suspicious requests for monitoring
			fmt.Printf("SECURITY WARNING: Blocked CORS request from unauthorized origin: %s\n", origin)
		}

		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Max-Age", "86400") // Cache preflight for 24 hours

		if c.Request.Method == "OPTIONS" {
			if isOriginAllowed {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}

		c.Next()
	}
}

// getAllowedOrigins returns the list of allowed origins from environment variables
func getAllowedOrigins() []string {
	// Default allowed origins for development and production
	defaultOrigins := []string{
		"http://localhost:3000", // Local development
		"http://localhost",      // Local with nginx
	}

	// Get additional allowed origins from environment variable
	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	if envOrigins != "" {
		// Split by comma and trim whitespace
		additionalOrigins := strings.Split(envOrigins, ",")
		for i, origin := range additionalOrigins {
			additionalOrigins[i] = strings.TrimSpace(origin)
		}
		defaultOrigins = append(defaultOrigins, additionalOrigins...)
	}

	return defaultOrigins
}

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func ValidateContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
