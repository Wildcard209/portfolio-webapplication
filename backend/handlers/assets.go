package handlers

import (
	"net/http"

	"github.com/Wildcard209/portfolio-webapplication/models"
	"github.com/Wildcard209/portfolio-webapplication/services"
	"github.com/gin-gonic/gin"
)

type AssetHandler struct {
	assetService *services.AssetService
}

func NewAssetHandler(assetService *services.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

// GetHeroBanner handles GET requests for hero banner image
// @Summary Get hero banner image
// @Description Get the current hero banner image
// @Tags assets
// @Produce image/jpeg,image/png,image/gif,image/webp
// @Success 200 {file} file "Hero banner image"
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/assets/hero-banner [get]
func (h *AssetHandler) GetHeroBanner(c *gin.Context) {
	data, contentType, err := h.assetService.GetHeroBanner()
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "Hero banner not found",
			Message: "No hero banner image has been uploaded yet",
		})
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	c.Data(http.StatusOK, contentType, data)
}

// UploadHeroBanner handles POST requests to upload hero banner image
// @Summary Upload hero banner image
// @Description Upload a new hero banner image (requires authentication)
// @Tags assets
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param adminToken path string true "Admin Token"
// @Param file formData file true "Hero banner image file"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 413 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /{adminToken}/admin/assets/hero-banner [post]
func (h *AssetHandler) UploadHeroBanner(c *gin.Context) {
	// Get the uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "No file uploaded",
			Message: "Please select a file to upload",
		})
		return
	}
	defer file.Close()

	// Check file size (limit to 10MB)
	const maxFileSize = 10 << 20 // 10MB
	if header.Size > maxFileSize {
		c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{
			Error:   "File too large",
			Message: "File size must be less than 10MB",
		})
		return
	}

	// Check file type
	contentType := header.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid file type",
			Message: "Only JPEG, PNG, GIF, and WebP images are allowed",
		})
		return
	}

	// Upload the file
	err = h.assetService.UploadHeroBanner(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Upload failed",
			Message: "Failed to upload hero banner image",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Hero banner uploaded successfully",
	})
}

// GetAssetInfo returns information about available assets
// @Summary Get asset information
// @Description Get information about available assets
// @Tags assets
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/assets/info [get]
func (h *AssetHandler) GetAssetInfo(c *gin.Context) {
	hasHeroBanner := h.assetService.HasHeroBanner()

	c.JSON(http.StatusOK, gin.H{
		"hero_banner_available": hasHeroBanner,
		"hero_banner_url":       "/api/assets/hero-banner",
	})
}

func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}
