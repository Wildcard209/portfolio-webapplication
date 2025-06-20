package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/config"
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
		var tokenString string
		var err error

		accessToken, cookieErr := c.Cookie("access_token")
		if cookieErr == nil && accessToken != "" {
			tokenString = accessToken
		} else {
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

		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			refreshToken, refreshErr := c.Cookie("refresh_token")
			if refreshErr != nil || refreshToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			refreshClaims, refreshValidErr := authService.ValidateRefreshToken(refreshToken)
			if refreshValidErr != nil {
				isHttps := os.Getenv("HTTPS_MODE") == "true"
				c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
				c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired, please login again"})
				c.Abort()
				return
			}

			admin, adminErr := adminRepo.GetAdminByToken(refreshToken)
			if adminErr != nil || admin == nil {
				isHttps := os.Getenv("HTTPS_MODE") == "true"
				c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
				c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Session has been revoked"})
				c.Abort()
				return
			}

			tokenPair, tokenErr := authService.GenerateTokenPair(refreshClaims.UserID, refreshClaims.Username)
			if tokenErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh session"})
				c.Abort()
				return
			}

			if updateErr := adminRepo.UpdateAdminToken(refreshClaims.UserID, tokenPair.RefreshToken, tokenPair.RefreshExpiresAt); updateErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session"})
				c.Abort()
				return
			}

			isHttps := os.Getenv("HTTPS_MODE") == "true"
			c.SetCookie(
				"access_token",
				tokenPair.AccessToken,
				int(tokenPair.AccessExpiresAt.Sub(time.Now()).Seconds()),
				"/",
				"",
				isHttps,
				true,
			)

			c.SetCookie(
				"refresh_token",
				tokenPair.RefreshToken,
				int(tokenPair.RefreshExpiresAt.Sub(time.Now()).Seconds()),
				"/",
				"",
				isHttps,
				true,
			)

			claims = &auth.CustomClaims{
				UserID:    refreshClaims.UserID,
				Username:  refreshClaims.Username,
				TokenType: "access",
			}

			c.Set("userID", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("admin", admin)

			c.Next()
			return
		}

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

func SecurityHeadersMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.SecurityHeaders.Enabled {
			c.Next()
			return
		}

		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")

		if cfg.SecurityHeaders.HTTPSMode {
			hstsValue := fmt.Sprintf("max-age=%d; includeSubDomains", cfg.SecurityHeaders.HSTSMaxAge)
			c.Header("Strict-Transport-Security", hstsValue)
		}

		cspPolicy := generateCSPPolicy(cfg.SecurityHeaders.CSPMode)
		c.Header("Content-Security-Policy", cspPolicy)

		c.Next()
	}
}

func generateCSPPolicy(mode string) string {
	reportURI := " report-uri /api/csp-report"

	switch mode {
	case "production":
		return "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' https:; " +
			"connect-src 'self'; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"frame-ancestors 'none';" +
			reportURI
	case "development":
		return "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https: http:; " +
			"font-src 'self' https: http:; " +
			"connect-src 'self' ws: wss:; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"frame-ancestors 'none';" +
			reportURI
	default:
		return "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https: http:; " +
			"font-src 'self' https: http:; " +
			"connect-src 'self' ws: wss:; " +
			"media-src 'self'; " +
			"object-src 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'; " +
			"frame-ancestors 'none';" +
			reportURI
	}
}

func CORSMiddleware() gin.HandlerFunc {
	allowedOrigins := getAllowedOrigins()

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		isOriginAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				isOriginAllowed = true
				break
			}
		}

		if isOriginAllowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		} else {
			fmt.Printf("SECURITY WARNING: Blocked CORS request from unauthorized origin: %s\n", origin)
		}

		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Max-Age", "86400")

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

func getAllowedOrigins() []string {
	defaultOrigins := []string{
		"http://localhost",
	}

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

func RequestBodySizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
			c.Next()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)

		c.Next()
	})
}

func FileUploadSizeLimitMiddleware(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		contentType := c.GetHeader("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		}

		c.Next()
	})
}

func GetRequestBodySizeLimit() int64 {
	defaultSize := int64(1 << 20)

	if sizeStr := os.Getenv("MAX_REQUEST_BODY_SIZE"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil && size > 0 {
			return size
		}
	}

	return defaultSize
}

func GetFileUploadSizeLimit() int64 {
	defaultSize := int64(10 << 20)

	if sizeStr := os.Getenv("MAX_FILE_SIZE"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil && size > 0 {
			return size
		}
	}

	return defaultSize
}
