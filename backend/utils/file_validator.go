package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

// FileSignature represents a file type with its magic number signature
type FileSignature struct {
	MimeType  string
	Extension string
	Signature []byte
}

// Common file signatures for validation
var allowedFileSignatures = []FileSignature{
	{MimeType: "image/jpeg", Extension: ".jpg", Signature: []byte{0xFF, 0xD8, 0xFF}},
	{MimeType: "image/jpeg", Extension: ".jpeg", Signature: []byte{0xFF, 0xD8, 0xFF}},
	{MimeType: "image/png", Extension: ".png", Signature: []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}},
	{MimeType: "image/gif", Extension: ".gif", Signature: []byte{0x47, 0x49, 0x46, 0x38}},
	{MimeType: "image/webp", Extension: ".webp", Signature: []byte{0x52, 0x49, 0x46, 0x46}}, // RIFF header, WebP has additional checks
}

// FileValidator provides comprehensive file validation
type FileValidator struct {
	maxFileSize     int64
	maxFilenameLen  int
	allowedTypes    []string
	filenamePattern *regexp.Regexp
}

// NewFileValidator creates a new file validator with specified limits
func NewFileValidator(maxFileSize int64, maxFilenameLen int, allowedTypes []string) *FileValidator {
	// Pattern to allow only safe filename characters
	safeFilenamePattern := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

	return &FileValidator{
		maxFileSize:     maxFileSize,
		maxFilenameLen:  maxFilenameLen,
		allowedTypes:    allowedTypes,
		filenamePattern: safeFilenamePattern,
	}
}

func (fv *FileValidator) ValidateFile(file multipart.File, header *multipart.FileHeader) error {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	if header.Size > fv.maxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", header.Size, fv.maxFileSize)
	}

	if header.Size == 0 {
		return errors.New("file is empty")
	}

	if err := fv.validateFilename(header.Filename); err != nil {
		return fmt.Errorf("invalid filename: %w", err)
	}

	contentType, err := fv.detectContentType(file)
	if err != nil {
		return fmt.Errorf("failed to detect file type: %w", err)
	}

	if !fv.isContentTypeAllowed(contentType) {
		return fmt.Errorf("file type %s is not allowed", contentType)
	}

	headerContentType := header.Header.Get("Content-Type")
	if headerContentType != "" && !fv.isContentTypeMatching(headerContentType, contentType) {
		return fmt.Errorf("header content type %s does not match actual content type %s", headerContentType, contentType)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to reset file pointer after validation: %w", err)
	}

	return nil
}

func (fv *FileValidator) validateFilename(filename string) error {
	if filename == "" {
		return errors.New("filename is empty")
	}

	if len(filename) > fv.maxFilenameLen {
		return fmt.Errorf("filename length %d exceeds maximum allowed length of %d", len(filename), fv.maxFilenameLen)
	}

	dangerous := []string{"..", "/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\x00"}
	for _, char := range dangerous {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains dangerous character: %s", char)
		}
	}

	if !fv.filenamePattern.MatchString(filename) {
		return errors.New("filename contains invalid characters (only alphanumeric, dots, underscores, and hyphens allowed)")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return errors.New("filename must have an extension")
	}

	return nil
}

func (fv *FileValidator) detectContentType(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file header: %w", err)
	}

	buffer = buffer[:n]

	for _, sig := range allowedFileSignatures {
		if len(buffer) >= len(sig.Signature) {
			if bytes.HasPrefix(buffer, sig.Signature) {
				if sig.MimeType == "image/webp" {
					if len(buffer) >= 12 && string(buffer[8:12]) == "WEBP" {
						return sig.MimeType, nil
					}
				} else {
					return sig.MimeType, nil
				}
			}
		}
	}

	return "", errors.New("unrecognized or unsupported file type")
}

func (fv *FileValidator) isContentTypeAllowed(contentType string) bool {
	for _, allowed := range fv.allowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

func (fv *FileValidator) isContentTypeMatching(headerType, actualType string) bool {
	headerType = strings.ToLower(strings.TrimSpace(headerType))
	actualType = strings.ToLower(strings.TrimSpace(actualType))

	if headerType == actualType {
		return true
	}

	jpegTypes := []string{"image/jpeg", "image/jpg"}
	if contains(jpegTypes, headerType) && contains(jpegTypes, actualType) {
		return true
	}

	return false
}

func SanitizeFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	safeChars := regexp.MustCompile(`[^a-zA-Z0-9._-]`)
	safeName := safeChars.ReplaceAllString(nameWithoutExt, "_")

	multipleUnderscores := regexp.MustCompile(`_{2,}`)
	safeName = multipleUnderscores.ReplaceAllString(safeName, "_")

	safeName = strings.Trim(safeName, "_")

	if safeName == "" {
		safeName = "unnamed"
	}

	maxNameLen := 100 - len(ext)
	if len(safeName) > maxNameLen {
		safeName = safeName[:maxNameLen]
	}

	return safeName + strings.ToLower(ext)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
