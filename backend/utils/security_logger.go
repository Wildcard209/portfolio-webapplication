package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type SecurityLogger struct {
	logger       *log.Logger
	isProduction bool
}

func NewSecurityLogger() *SecurityLogger {
	return &SecurityLogger{
		logger:       log.Default(),
		isProduction: os.Getenv("GIN_MODE") == "release",
	}
}

func (sl *SecurityLogger) SanitizeLogMessage(message string) string {
	if sl.isProduction && strings.TrimSpace(message) == "" {
		return message
	}

	sanitized := message

	sensitivePatterns := map[string]string{
		`postgres://[^:]+:[^@]+@`: "postgres://***:***@",
		`mysql://[^:]+:[^@]+@`:    "mysql://***:***@",
		`mongodb://[^:]+:[^@]+@`:  "mongodb://***:***@",

		`[Aa]pi[_-]?[Kk]ey[=:\s]+[^\s\n]+`: "api_key=***",
		`[Tt]oken[=:\s]+[^\s\n]+`:          "token=***",
		`[Aa]uthorization[=:\s]+[^\s\n]+`:  "authorization=***",
		`[Bb]earer\s+[^\s\n]+`:             "Bearer ***",

		`[Pp]assword[=:\s]+[^\s\n]+`: "password=***",
		`[Ss]ecret[=:\s]+[^\s\n]+`:   "secret=***",
		`[Kk]ey[=:\s]+[^\s\n]+`:      "key=***",

		`eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]*`: "JWT_TOKEN_REDACTED",

		`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`: func() string {
			if sl.isProduction {
				return "IP_REDACTED"
			}
			return "$0"
		}(),

		`/[a-zA-Z0-9_./\-]+`: func() string {
			if sl.isProduction {
				return "/PATH_REDACTED"
			}
			return "$0"
		}(),
	}

	for pattern, replacement := range sensitivePatterns {
		if strings.Contains(replacement, "$0") && !sl.isProduction {
			continue
		}

		re := regexp.MustCompile(pattern)
		sanitized = re.ReplaceAllString(sanitized, replacement)
	}

	if strings.Contains(sanitized, "postgres://") {
		lines := strings.Split(sanitized, "\n")
		for i, line := range lines {
			if strings.Contains(line, "postgres://") && strings.Contains(line, "@") {
				start := strings.Index(line, "postgres://")
				if start != -1 {
					remaining := line[start+11:]
					if credEnd := strings.Index(remaining, "@"); credEnd != -1 {
						before := line[:start+11]
						after := line[start+11+credEnd:]
						lines[i] = before + "***:***" + after
					}
				}
			}
		}
		sanitized = strings.Join(lines, "\n")
	}

	paramPatterns := []string{
		"password=", "pwd=", "secret=", "token=", "key=", "auth=", "authorization=",
	}

	for _, pattern := range paramPatterns {
		if strings.Contains(strings.ToLower(sanitized), pattern) {
			words := strings.Fields(sanitized)
			for j, word := range words {
				if strings.Contains(strings.ToLower(word), pattern) {
					if eqIndex := strings.Index(word, "="); eqIndex != -1 {
						words[j] = word[:eqIndex+1] + "***"
					}
				}
			}
			sanitized = strings.Join(words, " ")
		}
	}

	if sl.isProduction {
		sanitized = sl.sanitizeProductionSensitiveInfo(sanitized)
	}

	return sanitized
}

func (sl *SecurityLogger) sanitizeProductionSensitiveInfo(message string) string {
	if strings.Contains(message, "panic:") || strings.Contains(message, "runtime") {
		return "Internal error occurred - details suppressed in production"
	}

	productionSensitiveTerms := []string{
		"database error", "connection failed", "sql error", "query failed",
		"minio error", "file system error", "internal server error",
		"panic", "runtime error", "nil pointer", "index out of range",
	}

	lowerMessage := strings.ToLower(message)
	for _, term := range productionSensitiveTerms {
		if strings.Contains(lowerMessage, term) {
			return "An error occurred while processing the request"
		}
	}

	return message
}

func (sl *SecurityLogger) LogSecure(message string) {
	sanitized := sl.SanitizeLogMessage(message)
	sl.logger.Printf("[SECURE] %s", sanitized)
}

func (sl *SecurityLogger) LogSecureError(operation string, err error) {
	if err != nil {
		sanitized := sl.SanitizeLogMessage(err.Error())
		sl.logger.Printf("[SECURE_ERROR] %s failed: %s", operation, sanitized)
	}
}

func (sl *SecurityLogger) LogSecureInfo(message string) {
	sanitized := sl.SanitizeLogMessage(message)
	sl.logger.Printf("[SECURE_INFO] %s: %s", time.Now().UTC().Format(time.RFC3339), sanitized)
}

func (sl *SecurityLogger) LogSecureWarning(message string) {
	sanitized := sl.SanitizeLogMessage(message)
	sl.logger.Printf("[SECURE_WARNING] %s", sanitized)
}

func (sl *SecurityLogger) LogSecurityEvent(event string, details map[string]interface{}) {
	timestamp := time.Now().UTC().Format(time.RFC3339)

	message := "SECURITY_EVENT: " + event + " at " + timestamp

	for key, value := range details {
		if valueStr, ok := value.(string); ok {
			sanitizedValue := sl.SanitizeLogMessage(valueStr)
			message += " " + key + "=" + sanitizedValue
		} else {
			message += " " + key + "=" + sl.SanitizeLogMessage(fmt.Sprintf("%v", value))
		}
	}

	sl.logger.Printf("[SECURITY_AUDIT] %s", message)
}

func (sl *SecurityLogger) LogApplicationError(event string, err error, context map[string]interface{}) {
	if context == nil {
		context = make(map[string]interface{})
	}

	context["error"] = err.Error()
	context["timestamp"] = time.Now().UTC().Format(time.RFC3339)

	if sl.isProduction {
		context["error"] = sl.sanitizeProductionSensitiveInfo(err.Error())
	}

	sl.LogSecurityEvent(event, context)
}

func (sl *SecurityLogger) LogInformationDisclosureAttempt(clientIP, userAgent, endpoint, reason string) {
	context := map[string]interface{}{
		"client_ip":  clientIP,
		"user_agent": userAgent,
		"endpoint":   endpoint,
		"reason":     reason,
		"severity":   "HIGH",
	}

	sl.LogSecurityEvent("INFORMATION_DISCLOSURE_ATTEMPT", context)
}

func (sl *SecurityLogger) LogSensitiveDataAccess(userID interface{}, endpoint, method, clientIP string) {
	context := map[string]interface{}{
		"user_id":   userID,
		"endpoint":  endpoint,
		"method":    method,
		"client_ip": clientIP,
		"severity":  "MEDIUM",
	}

	sl.LogSecurityEvent("SENSITIVE_DATA_ACCESS", context)
}

func (sl *SecurityLogger) LogProductionError(component string, operation string, errorCode string) {
	if sl.isProduction {
		context := map[string]interface{}{
			"component":   component,
			"operation":   operation,
			"error_code":  errorCode,
			"environment": "production",
		}
		sl.LogSecurityEvent("PRODUCTION_ERROR", context)
	}
}
