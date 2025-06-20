package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type HeaderSanitizationConfig struct {
	RemoveServerHeaders  bool
	RemoveVersionHeaders bool
	RemoveDebugHeaders   bool
	CustomServerHeader   string
	HideFrameworkDetails bool
}

func NewHeaderSanitizationConfig() *HeaderSanitizationConfig {
	isProduction := os.Getenv("GIN_MODE") == "release"

	removeServerHeaders := getEnvBool("REMOVE_SERVER_HEADERS", isProduction)
	removeVersionHeaders := getEnvBool("REMOVE_VERSION_HEADERS", isProduction)
	removeDebugHeaders := getEnvBool("REMOVE_DEBUG_HEADERS", isProduction)
	customServerHeader := getEnvString("CUSTOM_SERVER_HEADER", "Portfolio-API")

	return &HeaderSanitizationConfig{
		RemoveServerHeaders:  removeServerHeaders,
		RemoveVersionHeaders: removeVersionHeaders,
		RemoveDebugHeaders:   removeDebugHeaders,
		CustomServerHeader:   customServerHeader,
		HideFrameworkDetails: isProduction,
	}
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}

func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func HeaderSanitizationMiddleware(config *HeaderSanitizationConfig) gin.HandlerFunc {
	if config == nil {
		config = NewHeaderSanitizationConfig()
	}

	return func(c *gin.Context) {
		c.Next()

		sanitizeResponseHeaders(c, config)
	}
}

func sanitizeResponseHeaders(c *gin.Context, config *HeaderSanitizationConfig) {
	headers := c.Writer.Header()

	sensitiveHeaders := []string{
		"Server",
		"X-Powered-By",
		"X-AspNet-Version",
		"X-AspNetMvc-Version",
		"X-Framework",
		"X-Runtime",
		"X-Version",
		"X-Served-By",
		"X-Generator",
	}

	if config.RemoveServerHeaders {
		for _, header := range sensitiveHeaders {
			headers.Del(header)
		}

		if config.CustomServerHeader != "" {
			headers.Set("Server", config.CustomServerHeader)
		}
	}

	if config.RemoveVersionHeaders {
		versionHeaders := []string{
			"X-API-Version",
			"X-App-Version",
			"X-Build-Version",
			"X-Release-Version",
		}
		for _, header := range versionHeaders {
			headers.Del(header)
		}
	}

	if config.RemoveDebugHeaders {
		debugHeaders := []string{
			"X-Debug",
			"X-Debug-Info",
			"X-Debug-Token",
			"X-Request-ID",
			"X-Trace-Id",
			"X-Response-Time",
			"X-Database-Queries",
			"X-Memory-Usage",
		}
		for _, header := range debugHeaders {
			headers.Del(header)
		}
	}

	if config.HideFrameworkDetails {
		frameworkHeaders := []string{
			"X-Gin-Mode",
			"X-Go-Version",
			"X-Golang-Version",
		}
		for _, header := range frameworkHeaders {
			headers.Del(header)
		}
	}

	sanitizeHeaderValues(headers)
}

func sanitizeHeaderValues(headers map[string][]string) {
	for headerName, values := range headers {
		for i, value := range values {
			if containsSensitiveInfo(value) {
				headers[headerName][i] = sanitizeValue(value)
			}
		}
	}
}

func containsSensitiveInfo(value string) bool {
	lowerValue := strings.ToLower(value)

	sensitivePatterns := []string{
		"password",
		"secret",
		"token",
		"key",
		"private",
		"internal",
		"debug",
		"dev",
		"localhost",
		"127.0.0.1",
		"development",
		"staging",
		"test",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerValue, pattern) {
			return true
		}
	}

	return false
}

func sanitizeValue(value string) string {
	return "[REDACTED]"
}

func NoServerHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		c.Writer.Header().Del("Server")
		c.Writer.Header().Del("X-Powered-By")

		for key := range c.Writer.Header() {
			if strings.HasPrefix(strings.ToLower(key), "x-") &&
				(strings.Contains(strings.ToLower(key), "server") ||
					strings.Contains(strings.ToLower(key), "powered")) {
				c.Writer.Header().Del(key)
			}
		}
	}
}

func DebugHeadersMiddleware() gin.HandlerFunc {
	isDevelopment := os.Getenv("GIN_MODE") != "release"

	return func(c *gin.Context) {
		if isDevelopment {
			c.Header("X-Environment", "development")
			c.Header("X-Debug-Mode", "enabled")
		}

		c.Next()

		if isDevelopment {
			if startTime, exists := c.Get("start_time"); exists {
				if start, ok := startTime.(int64); ok {
					c.Header("X-Response-Time", "calculated-response-time")
					_ = start
				}
			}
		}
	}
}
