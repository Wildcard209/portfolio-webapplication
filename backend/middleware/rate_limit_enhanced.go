package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/config"
	"github.com/Wildcard209/portfolio-webapplication/utils"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type RateLimitType string

const (
	RateLimitLogin   RateLimitType = "login"
	RateLimitRefresh RateLimitType = "refresh"
	RateLimitUpload  RateLimitType = "upload"
	RateLimitAPI     RateLimitType = "api"
	RateLimitPublic  RateLimitType = "public"
	RateLimitAdmin   RateLimitType = "admin"
)

func RateLimitMiddlewareWithConfig(rateLimitType RateLimitType, rateLimitConfig *config.EnhancedRateLimitConfig) gin.HandlerFunc {
	var rateLimit config.RateLimit

	switch rateLimitType {
	case RateLimitLogin:
		rateLimit = rateLimitConfig.Login
	case RateLimitRefresh:
		rateLimit = rateLimitConfig.Refresh
	case RateLimitUpload:
		rateLimit = rateLimitConfig.Upload
	case RateLimitAPI:
		rateLimit = rateLimitConfig.API
	case RateLimitPublic:
		rateLimit = rateLimitConfig.Public
	case RateLimitAdmin:
		rateLimit = rateLimitConfig.Admin
	default:
		rateLimit = rateLimitConfig.API
	}

	store := memory.NewStore()
	rateLimiter := limiter.New(store, rateLimit.ToLimiterRate())

	return func(c *gin.Context) {
		addRateLimitHeaders(c, rateLimit)

		logRateLimitAttempt(c, string(rateLimitType), rateLimit)

		middleware := mgin.NewMiddleware(rateLimiter)
		middleware(c)
	}
}

func addRateLimitHeaders(c *gin.Context, rateLimit config.RateLimit) {
	c.Header("X-RateLimit-Limit", strconv.Itoa(rateLimit.Requests))

	now := time.Now()
	resetTime := now.Truncate(rateLimit.Period).Add(rateLimit.Period)
	c.Header("X-RateLimit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
}

func logRateLimitAttempt(c *gin.Context, rateLimitType string, rateLimit config.RateLimit) {
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	endpoint := c.Request.URL.Path
	method := c.Request.Method

	// Log the rate limit attempt
	securityLogger := utils.NewSecurityLogger()
	securityLogger.LogSecurityEvent("rate_limit_attempt", map[string]interface{}{
		"type":       rateLimitType,
		"client_ip":  clientIP,
		"user_agent": userAgent,
		"endpoint":   endpoint,
		"method":     method,
		"limit":      rateLimit.Requests,
		"period":     rateLimit.Period.String(),
		"timestamp":  time.Now().UTC(),
	})
}

func RateLimitViolationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() == http.StatusTooManyRequests {
			logRateLimitViolation(c)
		}
	}
}

func logRateLimitViolation(c *gin.Context) {
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	endpoint := c.Request.URL.Path
	method := c.Request.Method

	securityLogger := utils.NewSecurityLogger()
	securityLogger.LogSecurityEvent("rate_limit_violation", map[string]interface{}{
		"client_ip":   clientIP,
		"user_agent":  userAgent,
		"endpoint":    endpoint,
		"method":      method,
		"timestamp":   time.Now().UTC(),
		"severity":    "medium",
		"description": fmt.Sprintf("Rate limit exceeded for %s %s from IP %s", method, endpoint, clientIP),
	})
}

func GetRateLimitForEndpoint(endpoint string, rateLimitConfig *config.EnhancedRateLimitConfig) config.RateLimit {
	switch endpoint {
	case "/api/admin/login":
		return rateLimitConfig.Login
	case "/api/admin/refresh":
		return rateLimitConfig.Refresh
	case "/api/admin/assets/hero-banner":
		return rateLimitConfig.Upload
	case "/api/admin/logout":
		return rateLimitConfig.Admin
	default:
		if endpoint == "/api/hello" || endpoint == "/api/assets/hero-banner" || endpoint == "/api/assets/info" {
			return rateLimitConfig.Public
		}
		return rateLimitConfig.API
	}
}
