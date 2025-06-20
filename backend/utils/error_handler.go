package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/Wildcard209/portfolio-webapplication/models"
	"github.com/gin-gonic/gin"
)

type ErrorHandler struct {
	isProduction bool
	logger       *SecurityLogger
}

type ErrorLevel int

const (
	ErrorLevelInfo ErrorLevel = iota
	ErrorLevelWarning
	ErrorLevelError
	ErrorLevelCritical
)

func NewErrorHandler() *ErrorHandler {
	isProduction := os.Getenv("GIN_MODE") == "release"
	return &ErrorHandler{
		isProduction: isProduction,
		logger:       NewSecurityLogger(),
	}
}

func (eh *ErrorHandler) HandleError(c *gin.Context, err error, userMessage string, level ErrorLevel) {
	eh.logErrorWithContext(c, err, level)

	if eh.isProduction {
		eh.respondWithSanitizedError(c, userMessage, http.StatusInternalServerError)
	} else {
		eh.respondWithDetailedError(c, err, userMessage, http.StatusInternalServerError)
	}
}

func (eh *ErrorHandler) HandleAuthError(c *gin.Context, err error, userMessage string) {
	eh.logErrorWithContext(c, err, ErrorLevelWarning)

	eh.respondWithSanitizedError(c, userMessage, http.StatusUnauthorized)
}

func (eh *ErrorHandler) HandleValidationError(c *gin.Context, err error, userMessage string) {
	eh.logErrorWithContext(c, err, ErrorLevelInfo)

	if eh.isProduction {
		eh.respondWithSanitizedError(c, userMessage, http.StatusBadRequest)
	} else {
		eh.respondWithDetailedError(c, err, userMessage, http.StatusBadRequest)
	}
}

func (eh *ErrorHandler) HandleNotFoundError(c *gin.Context, resource string) {
	err := fmt.Errorf("resource not found: %s", resource)
	eh.logErrorWithContext(c, err, ErrorLevelInfo)

	eh.respondWithSanitizedError(c, "Resource not found", http.StatusNotFound)
}

func (eh *ErrorHandler) HandleRateLimitError(c *gin.Context, message string) {
	err := fmt.Errorf("rate limit exceeded for IP: %s", c.ClientIP())
	eh.logErrorWithContext(c, err, ErrorLevelWarning)

	eh.respondWithSanitizedError(c, message, http.StatusTooManyRequests)
}

func (eh *ErrorHandler) logErrorWithContext(c *gin.Context, err error, level ErrorLevel) {
	if err == nil {
		return
	}

	context := map[string]interface{}{
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"ip":         c.ClientIP(),
		"user_agent": c.GetHeader("User-Agent"),
		"error":      err.Error(),
		"level":      eh.levelToString(level),
	}

	if !eh.isProduction {
		context["stack"] = eh.getStackTrace()
	}

	if userID, exists := c.Get("user_id"); exists {
		context["user_id"] = userID
	}

	eh.logger.LogSecurityEvent("APPLICATION_ERROR", context)
}

func (eh *ErrorHandler) respondWithSanitizedError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func (eh *ErrorHandler) respondWithDetailedError(c *gin.Context, err error, message string, statusCode int) {
	response := gin.H{
		"error":   http.StatusText(statusCode),
		"message": message,
	}

	if err != nil {
		response["details"] = err.Error()
		response["stack"] = eh.getStackTrace()
	}

	c.JSON(statusCode, response)
}

func (eh *ErrorHandler) getStackTrace() []string {
	var stack []string
	for i := 2; i < 10; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, "portfolio-webapplication") {
			stack = append(stack, fmt.Sprintf("%s:%d", file, line))
		}
	}
	return stack
}

func (eh *ErrorHandler) levelToString(level ErrorLevel) string {
	switch level {
	case ErrorLevelInfo:
		return "INFO"
	case ErrorLevelWarning:
		return "WARNING"
	case ErrorLevelError:
		return "ERROR"
	case ErrorLevelCritical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func (eh *ErrorHandler) SanitizeErrorMessage(message string) string {
	if eh.isProduction {
		sensitivePatterns := []string{
			"password",
			"token",
			"secret",
			"key",
			"connection failed",
			"database error",
			"sql:",
			"postgres:",
			"minio:",
			"panic:",
		}

		lowerMessage := strings.ToLower(message)
		for _, pattern := range sensitivePatterns {
			if strings.Contains(lowerMessage, pattern) {
				return "An error occurred while processing your request"
			}
		}
	}
	return message
}

func (eh *ErrorHandler) LogCriticalError(component string, err error, context map[string]interface{}) {
	if context == nil {
		context = make(map[string]interface{})
	}

	context["component"] = component
	context["error"] = err.Error()
	context["level"] = "CRITICAL"

	eh.logger.LogSecurityEvent("CRITICAL_ERROR", context)

	log.Printf("CRITICAL ERROR in %s: %v", component, err)
}
