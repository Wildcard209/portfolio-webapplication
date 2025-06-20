package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/auth"
	"github.com/Wildcard209/portfolio-webapplication/database"
	"github.com/Wildcard209/portfolio-webapplication/repository"
)

type AdminService struct {
	db               *sql.DB
	authService      *auth.AuthService
	adminRepo        *repository.AdminRepository
	loginAttemptRepo *repository.LoginAttemptRepository
}

func NewAdminService(db *sql.DB, authService *auth.AuthService) *AdminService {
	return &AdminService{
		db:               db,
		authService:      authService,
		adminRepo:        repository.NewAdminRepository(db),
		loginAttemptRepo: repository.NewLoginAttemptRepository(db),
	}
}

func (s *AdminService) InitializeAdminSystem() error {
	log.Println("Initializing admin system...")

	if err := database.RunMigrations(s.db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	adminCount, err := s.adminRepo.CountAdmins()
	if err != nil {
		return fmt.Errorf("failed to count admin users: %w", err)
	}

	if adminCount == 0 {
		log.Println("No admin user found, creating default admin user...")
		if err := s.createDefaultAdmin(); err != nil {
			return fmt.Errorf("failed to create default admin: %w", err)
		}
	} else {
		log.Printf("Found %d admin user(s) in the system", adminCount)
	}

	if err := s.adminRepo.CleanupExpiredTokens(); err != nil {
		log.Printf("Warning: Failed to cleanup expired tokens: %v", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -30)
	if err := s.loginAttemptRepo.CleanupOldLoginAttempts(cutoffTime); err != nil {
		log.Printf("Warning: Failed to cleanup old login attempts: %v", err)
	}

	log.Println("Admin system initialized successfully")
	return nil
}

func (s *AdminService) createDefaultAdmin() error {
	username := os.Getenv("ADMIN_USER")
	password := os.Getenv("ADMIN_PASSWORD")

	if username == "" || password == "" {
		return fmt.Errorf("ADMIN_USER and ADMIN_PASSWORD environment variables are required")
	}

	hashedPassword, err := s.authService.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	admin, err := s.adminRepo.CreateAdminWithHashVersion(username, hashedPassword, "", 2)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("Default admin user created successfully with ID: %d, Username: %s", admin.ID, admin.Username)
	return nil
}

func (s *AdminService) StartMaintenanceTasks() {
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			s.runMaintenanceTasks()
		}
	}()

	log.Println("Maintenance tasks started")
}

func (s *AdminService) runMaintenanceTasks() {
	if err := s.adminRepo.CleanupExpiredTokens(); err != nil {
		log.Printf("Maintenance: Failed to cleanup expired tokens: %v", err)
	}

	cutoffTime := time.Now().AddDate(0, 0, -7)
	if err := s.loginAttemptRepo.CleanupOldLoginAttempts(cutoffTime); err != nil {
		log.Printf("Maintenance: Failed to cleanup old login attempts: %v", err)
	}
}

func (s *AdminService) GetRepositories() (*repository.AdminRepository, *repository.LoginAttemptRepository) {
	return s.adminRepo, s.loginAttemptRepo
}
