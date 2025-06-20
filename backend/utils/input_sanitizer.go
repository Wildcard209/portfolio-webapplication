package utils

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"unicode"
)

type InputSanitizer struct {
	maxStringLength int
	allowedChars    *regexp.Regexp
}

func NewInputSanitizer(maxStringLength int) *InputSanitizer {
	safeCharsPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.@#$%&*+=!?()[\]{}|;:,<>/"'\\]+$`)

	return &InputSanitizer{
		maxStringLength: maxStringLength,
		allowedChars:    safeCharsPattern,
	}
}

func (is *InputSanitizer) SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	sanitized := strings.TrimSpace(input)

	sanitized = removeControlCharacters(sanitized)

	sanitized = html.EscapeString(sanitized)

	if len(sanitized) > is.maxStringLength {
		sanitized = sanitized[:is.maxStringLength]
	}

	return sanitized
}

func (is *InputSanitizer) ValidateString(input string, fieldName string, minLength, maxLength int) error {
	if input == "" && minLength > 0 {
		return NewValidationError(fieldName, "is required")
	}

	if len(input) < minLength {
		return NewValidationError(fieldName, "must be at least %d characters long", minLength)
	}

	if len(input) > maxLength {
		return NewValidationError(fieldName, "must not exceed %d characters", maxLength)
	}

	if containsDangerousPatterns(input) {
		return NewValidationError(fieldName, "contains potentially dangerous content")
	}

	return nil
}

func (is *InputSanitizer) SanitizeUsername(username string) string {
	if username == "" {
		return ""
	}

	sanitized := strings.ToLower(strings.TrimSpace(username))

	usernamePattern := regexp.MustCompile(`[^a-z0-9_-]`)
	sanitized = usernamePattern.ReplaceAllString(sanitized, "")

	if len(sanitized) > 50 {
		sanitized = sanitized[:50]
	}

	return sanitized
}

func (is *InputSanitizer) ValidateUsername(username string) error {
	if username == "" {
		return NewValidationError("username", "is required")
	}

	if len(username) < 3 {
		return NewValidationError("username", "must be at least 3 characters long")
	}

	if len(username) > 50 {
		return NewValidationError("username", "must not exceed 50 characters")
	}

	usernamePattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernamePattern.MatchString(username) {
		return NewValidationError("username", "can only contain letters, numbers, underscores, and hyphens")
	}

	if !unicode.IsLetter(rune(username[0])) && !unicode.IsDigit(rune(username[0])) {
		return NewValidationError("username", "must start with a letter or number")
	}

	return nil
}

func (is *InputSanitizer) ValidatePassword(password string) error {
	if password == "" {
		return NewValidationError("password", "is required")
	}

	if len(password) < 8 {
		return NewValidationError("password", "must be at least 8 characters long")
	}

	if len(password) > 128 {
		return NewValidationError("password", "must not exceed 128 characters")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper {
		return NewValidationError("password", "must contain at least one uppercase letter")
	}

	if !hasLower {
		return NewValidationError("password", "must contain at least one lowercase letter")
	}

	if !hasDigit {
		return NewValidationError("password", "must contain at least one digit")
	}

	if !hasSpecial {
		return NewValidationError("password", "must contain at least one special character")
	}

	return nil
}

func removeControlCharacters(input string) string {
	var result strings.Builder
	for _, r := range input {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

func containsDangerousPatterns(input string) bool {
	lowerInput := strings.ToLower(input)

	sqlPatterns := []string{
		"union select", "drop table", "delete from", "insert into",
		"update set", "create table", "alter table", "exec(",
		"execute(", "sp_", "xp_", "--;", "/*", "*/",
	}

	xssPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:",
		"onload=", "onerror=", "onclick=", "onmouseover=",
		"eval(", "expression(", "url(javascript",
	}

	for _, pattern := range sqlPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	for _, pattern := range xssPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

type ValidationError struct {
	Field   string
	Message string
}

func (ve *ValidationError) Error() string {
	return ve.Field + " " + ve.Message
}

func NewValidationError(field, message string, args ...interface{}) *ValidationError {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
