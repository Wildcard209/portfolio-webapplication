package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"github.com/minio/minio-go/v7"
)

type AssetService struct {
	minioClient *minio.Client
	bucketName  string
}

func NewAssetService(minioClient *minio.Client) *AssetService {
	service := &AssetService{
		minioClient: minioClient,
		bucketName:  "portfolio-assets",
	}

	// Ensure bucket exists
	if err := service.ensureBucketExists(); err != nil {
		log.Printf("Warning: Failed to ensure bucket exists: %v", err)
	}

	return service
}

func (s *AssetService) ensureBucketExists() error {
	ctx := context.Background()

	exists, err := s.minioClient.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = s.minioClient.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Created bucket: %s", s.bucketName)
	}

	return nil
}

// GetHeroBanner gets the current hero banner image
func (s *AssetService) GetHeroBanner() ([]byte, string, error) {
	ctx := context.Background()
	objectName := "hero-banner"

	// Try to get the object
	object, err := s.minioClient.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get hero banner: %w", err)
	}
	defer object.Close()

	// Read the object data
	data, err := io.ReadAll(object)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read hero banner data: %w", err)
	}

	// Get object info to determine content type
	objectInfo, err := object.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object info: %w", err)
	}

	contentType := objectInfo.ContentType
	if contentType == "" {
		// Try to determine content type from object name or data
		contentType = "application/octet-stream"
	}

	return data, contentType, nil
}

// UploadHeroBanner uploads a new hero banner image
func (s *AssetService) UploadHeroBanner(file multipart.File, header *multipart.FileHeader) error {
	ctx := context.Background()
	objectName := "hero-banner"

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Try to determine from filename
		filename := strings.ToLower(header.Filename)
		if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
			contentType = "image/jpeg"
		} else if strings.HasSuffix(filename, ".png") {
			contentType = "image/png"
		} else if strings.HasSuffix(filename, ".gif") {
			contentType = "image/gif"
		} else if strings.HasSuffix(filename, ".webp") {
			contentType = "image/webp"
		} else {
			contentType = "application/octet-stream"
		}
	}

	// Upload the file
	_, err := s.minioClient.PutObject(ctx, s.bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload hero banner: %w", err)
	}

	log.Printf("Successfully uploaded hero banner: %s (size: %d bytes, type: %s)",
		header.Filename, header.Size, contentType)
	return nil
}

// HasHeroBanner checks if a hero banner exists
func (s *AssetService) HasHeroBanner() bool {
	ctx := context.Background()
	objectName := "hero-banner"

	_, err := s.minioClient.StatObject(ctx, s.bucketName, objectName, minio.StatObjectOptions{})
	return err == nil
}

// DeleteHeroBanner deletes the current hero banner
func (s *AssetService) DeleteHeroBanner() error {
	ctx := context.Background()
	objectName := "hero-banner"

	err := s.minioClient.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete hero banner: %w", err)
	}

	log.Println("Successfully deleted hero banner")
	return nil
}

// GetAssetURL returns the URL for an asset (for direct access)
func (s *AssetService) GetAssetURL(objectName string) (string, error) {
	// For now, we'll serve assets through our API endpoint
	// In production, you might want to use presigned URLs
	return fmt.Sprintf("/api/assets/%s", objectName), nil
}
