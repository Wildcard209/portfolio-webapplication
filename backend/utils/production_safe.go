package utils

import (
	"os"
	"strings"
)

type ProductionSafeResponse struct {
	isProduction bool
}

func NewProductionSafeResponse() *ProductionSafeResponse {
	return &ProductionSafeResponse{
		isProduction: os.Getenv("GIN_MODE") == "release",
	}
}

func (psr *ProductionSafeResponse) SafeErrorMessage(originalError, userFriendlyMessage string) string {
	if psr.isProduction {
		if psr.containsSensitiveInfo(originalError) {
			return userFriendlyMessage
		}
		return userFriendlyMessage
	}

	if originalError != "" {
		return originalError
	}
	return userFriendlyMessage
}

func (psr *ProductionSafeResponse) containsSensitiveInfo(message string) bool {
	lowerMessage := strings.ToLower(message)

	sensitiveIndicators := []string{
		"database", "sql", "postgres", "mysql", "mongodb",
		"connection", "timeout", "refused", "network",
		"file", "path", "directory", "permission",
		"token", "secret", "key", "password", "auth",
		"panic", "runtime", "stack", "trace",
		"minio", "s3", "bucket",
		"nil pointer", "index out of range",
		"unmarshaling", "parsing", "decode",
		"localhost", "127.0.0.1", "internal",
	}

	for _, indicator := range sensitiveIndicators {
		if strings.Contains(lowerMessage, indicator) {
			return true
		}
	}

	return false
}

func (psr *ProductionSafeResponse) GetGenericErrorMessage(errorType string) string {
	genericMessages := map[string]string{
		"database":       "A data processing error occurred",
		"authentication": "Authentication failed",
		"authorization":  "Access denied",
		"validation":     "Invalid input provided",
		"file_upload":    "File upload failed",
		"file_download":  "File retrieval failed",
		"network":        "A network error occurred",
		"internal":       "An internal error occurred",
		"rate_limit":     "Too many requests. Please try again later",
		"not_found":      "The requested resource was not found",
		"conflict":       "A conflict occurred with the current state",
		"timeout":        "The request timed out",
		"maintenance":    "Service temporarily unavailable",
	}

	if message, exists := genericMessages[errorType]; exists {
		return message
	}

	return "An error occurred while processing your request"
}

func (psr *ProductionSafeResponse) SanitizeUserInput(input string) string {
	if psr.isProduction {
		if len(input) > 100 {
			return input[:100] + "...[truncated]"
		}

		dangerousPatterns := []string{
			"<script", "</script>", "javascript:", "onclick=", "onerror=",
			"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
			"drop table", "delete from", "insert into", "update set",
		}

		sanitized := strings.ToLower(input)
		for _, pattern := range dangerousPatterns {
			if strings.Contains(sanitized, pattern) {
				return "[POTENTIALLY_DANGEROUS_INPUT_REDACTED]"
			}
		}
	}

	return input
}

func (psr *ProductionSafeResponse) IsProductionMode() bool {
	return psr.isProduction
}

func (psr *ProductionSafeResponse) GetSafeHeaders() map[string]string {
	headers := make(map[string]string)

	if psr.isProduction {
		headers["X-Environment"] = "production"
		headers["X-Content-Type-Options"] = "nosniff"
		headers["X-Frame-Options"] = "DENY"
	} else {
		headers["X-Environment"] = "development"
		headers["X-Debug-Mode"] = "enabled"
		headers["X-Content-Type-Options"] = "nosniff"
		headers["X-Frame-Options"] = "DENY"
	}

	return headers
}
