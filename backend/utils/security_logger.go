package utils

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type SecurityLogger struct {
	logger *log.Logger
}

func NewSecurityLogger() *SecurityLogger {
	return &SecurityLogger{
		logger: log.Default(),
	}
}

func (sl *SecurityLogger) SanitizeLogMessage(message string) string {
	sanitized := message

	if strings.Contains(sanitized, "postgres://") {
		sanitized = strings.ReplaceAll(sanitized, "postgres://", "postgres://")
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

	patterns := []string{
		"password=",
		"pwd=",
		"secret=",
		"token=",
		"key=",
	}

	for _, pattern := range patterns {
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

	return sanitized
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
