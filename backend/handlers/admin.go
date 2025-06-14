package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/models"
	"github.com/Wildcard209/portfolio-webapplication/repository"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

type AdminHandler struct {
	authService         *auth.AuthService
	adminRepo           *repository.AdminRepository
	loginAttemptRepo    *repository.LoginAttemptRepository
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

	clientIP := c.ClientIP()

	failedAttempts, err := h.loginAttemptRepo.GetFailedLoginAttempts(
		clientIP,
		time.Now().Add(-h.failedAttemptWindow),
	)
	if err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Failed to check login attempts: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to process login request",
		})
		return
	}

	if failedAttempts >= h.maxFailedAttempts {
		h.logLoginAttempt(c, false, fmt.Sprintf("IP locked out due to %d failed attempts", failedAttempts))
		c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
			Error:   "Too many failed attempts",
			Message: fmt.Sprintf("IP address locked out for %v due to too many failed login attempts", h.lockoutDuration),
		})
		return
	}

	admin, err := h.adminRepo.GetAdminByUsername(req.Username)
	if err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Database error: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to process login request",
		})
		return
	}

	if admin == nil {
		h.logLoginAttempt(c, false, "User not found")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Username or password is incorrect",
		})
		return
	}

	if err := h.authService.VerifyPassword(admin.PasswordHash, req.Password, admin.PasswordSalt); err != nil {
		h.logLoginAttempt(c, false, "Invalid password")
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "Invalid credentials",
			Message: "Username or password is incorrect",
		})
		return
	}

	token, expiresAt, err := h.authService.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Failed to generate token: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to generate authentication token",
		})
		return
	}

	if err := h.adminRepo.UpdateAdminToken(admin.ID, token, expiresAt); err != nil {
		h.logLoginAttempt(c, false, fmt.Sprintf("Failed to update token: %v", err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal server error",
			Message: "Failed to complete login process",
		})
		return
	}

	h.logLoginAttempt(c, true, "Login successful")

	response := models.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
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

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Successfully logged out",
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

func GetRateLimiterForLogin() limiter.Rate {
	return limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}
}
