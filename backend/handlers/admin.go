package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/models"
	"github.com/Wildcard209/portfolio-webapplication/repository"
	"github.com/Wildcard209/portfolio-webapplication/utils"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	authService         *auth.AuthService
	adminRepo           *repository.AdminRepository
	loginAttemptRepo    *repository.LoginAttemptRepository
	inputSanitizer      *utils.InputSanitizer
	errorHandler        *utils.ErrorHandler
	maxFailedAttempts   int
	lockoutDuration     time.Duration
	failedAttemptWindow time.Duration
}

func NewAdminHandler(
	authService *auth.AuthService,
	adminRepo *repository.AdminRepository,
	loginAttemptRepo *repository.LoginAttemptRepository,
) *AdminHandler {
	return &AdminHandler{
		authService:         authService,
		adminRepo:           adminRepo,
		loginAttemptRepo:    loginAttemptRepo,
		inputSanitizer:      utils.NewInputSanitizer(1000),
		errorHandler:        utils.NewErrorHandler(),
		maxFailedAttempts:   5,
		lockoutDuration:     15 * time.Minute,
		failedAttemptWindow: 5 * time.Minute,
	}
}

// Login handles admin login
// @Summary Admin login
// @Description Authenticate admin user and return JWT token
// @Tags admin
// @Accept json
// @Produce json
// @Param adminToken path string true "Admin Token"
// @Param loginRequest body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 429 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /{adminToken}/admin/login [post]
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Invalid request format: %v", err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	if err := h.inputSanitizer.ValidateUsername(req.Username); err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Invalid username: %v", err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: "Username " + err.(*utils.ValidationError).Message,
		})
		return
	}

	if err := h.inputSanitizer.ValidateString(req.Password, "password", 1, 128); err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Invalid password: %v", err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: "Password " + err.(*utils.ValidationError).Message,
		})
		return
	}

	req.Username = h.inputSanitizer.SanitizeUsername(req.Username)

	clientIP := c.ClientIP()

	failedAttempts, err := h.loginAttemptRepo.GetFailedLoginAttempts(
		clientIP,
		time.Now().Add(-h.failedAttemptWindow),
	)
	if err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Failed to check login attempts: %v", err))
		h.errorHandler.HandleError(c, err, "Failed to process login request", utils.ErrorLevelError)
		return
	}

	if failedAttempts >= h.maxFailedAttempts {
		h.logLoginAttempt(c, false, fmt.Sprintf("IP locked out due to %d failed attempts", failedAttempts))
		h.errorHandler.HandleRateLimitError(c, fmt.Sprintf("IP address locked out for %v due to too many failed login attempts", h.lockoutDuration))
		return
	}

	admin, err := h.adminRepo.GetAdminByUsername(req.Username)
	if err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Database error: %v", err))
		h.errorHandler.HandleError(c, err, "Failed to process login request", utils.ErrorLevelError)
		return
	}

	if admin == nil {
		h.logLoginAttempt(c, false, "User not found")
		h.errorHandler.HandleAuthError(c, fmt.Errorf("user not found"), "Username or password is incorrect")
		return
	}

	if err := h.authService.VerifyPasswordWithHashVersion(admin.PasswordHash, req.Password, admin.HashVersion, admin.PasswordSalt); err != nil {
		h.logLoginAttempt(c, false, "Invalid password: "+err.Error())
		if err.Error() == "legacy password format no longer supported - please reset your password" {
			h.errorHandler.HandleAuthError(c, err, "Your password format needs to be updated. Please reset your password.")
		} else {
			h.errorHandler.HandleAuthError(c, err, "Username or password is incorrect")
		}
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(admin.ID, admin.Username)
	if err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Failed to generate token: %v", err))
		h.errorHandler.HandleError(c, err, "Failed to generate authentication token", utils.ErrorLevelError)
		return
	}

	if err := h.adminRepo.UpdateAdminToken(admin.ID, tokenPair.RefreshToken, tokenPair.RefreshExpiresAt); err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Failed to update token: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to complete login process",
		})
		return
	}

	isHttps := os.Getenv("HTTPS_MODE") == "true"

	c.SetCookie(
		"access_token",
		tokenPair.AccessToken,
		int(tokenPair.AccessExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"",
		isHttps,
		true,
	)

	c.SetCookie(
		"refresh_token",
		tokenPair.RefreshToken,
		int(tokenPair.RefreshExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"",
		isHttps,
		true,
	)

	h.logLoginAttempt(c, true, "Login successful")

	response := models.LoginResponse{
		Token:     "",
		ExpiresAt: tokenPair.AccessExpiresAt,
		User: models.AdminUser{
			ID:        admin.ID,
			Username:  admin.Username,
			LastLogin: admin.LastLogin,
		},
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles admin logout
// @Summary Admin logout
// @Description Invalidate current admin session
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param adminToken path string true "Admin Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /{adminToken}/admin/logout [post]
func (h *AdminHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	adminID := userID.(int)

	if err := h.adminRepo.InvalidateAdminToken(adminID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to logout",
		})
		return
	}

	isHttps := os.Getenv("HTTPS_MODE") == "true"
	c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
	c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags admin
// @Produce json
// @Param adminToken path string true "Admin Token"
// @Success 200 {object} models.SuccessResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /{adminToken}/admin/refresh [post]
func (h *AdminHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Refresh token not found",
		})
		return
	}

	claims, err := h.authService.ValidateRefreshToken(refreshToken)
	if err != nil {
		isHttps := os.Getenv("HTTPS_MODE") == "true"
		c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
		c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Invalid refresh token",
		})
		return
	}

	admin, err := h.adminRepo.GetAdminByToken(refreshToken)
	if err != nil || admin == nil {
		isHttps := os.Getenv("HTTPS_MODE") == "true"
		c.SetCookie("access_token", "", -1, "/", "", isHttps, true)
		c.SetCookie("refresh_token", "", -1, "/", "", isHttps, true)

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Unauthorized",
			Message: "Refresh token has been revoked",
		})
		return
	}

	tokenPair, err := h.authService.GenerateTokenPair(claims.UserID, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to generate new tokens",
		})
		return
	}

	if err := h.adminRepo.UpdateAdminToken(claims.UserID, tokenPair.RefreshToken, tokenPair.RefreshExpiresAt); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to update tokens",
		})
		return
	}

	isHttps := os.Getenv("HTTPS_MODE") == "true"
	c.SetCookie(
		"access_token",
		tokenPair.AccessToken,
		int(tokenPair.AccessExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"",
		isHttps,
		true,
	)

	c.SetCookie(
		"refresh_token",
		tokenPair.RefreshToken,
		int(tokenPair.RefreshExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"",
		isHttps,
		true,
	)

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Tokens refreshed successfully",
	})
}

func (h *AdminHandler) logLoginAttempt(c *gin.Context, success bool, details string) {
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	detailsPtr := &details
	if details == "" {
		detailsPtr = nil
	}

	go func() {
		if err := h.loginAttemptRepo.CreateLoginAttempt(clientIP, userAgent, success, detailsPtr); err != nil {
			fmt.Printf("Failed to log login attempt: %v\n", err)
		}
	}()
}
