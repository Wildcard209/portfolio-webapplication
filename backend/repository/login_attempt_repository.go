package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/database"
	"github.com/Wildcard209/portfolio-webapplication/models"
)

type LoginAttemptRepository struct {
	db          *sql.DB
	queryLoader *database.QueryLoader
}

func NewLoginAttemptRepository(db *sql.DB) *LoginAttemptRepository {
	queryLoader, err := database.NewQueryLoader()
	if err != nil {
		fmt.Printf("Warning: Failed to load queries: %v\n", err)
	}

	return &LoginAttemptRepository{
		db:          db,
		queryLoader: queryLoader,
	}
}

func (r *LoginAttemptRepository) CreateLoginAttempt(ipAddress, userAgent string, success bool, details *string) error {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.LoginAttempt.CreateLoginAttempt)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	_, err = r.db.Exec(query, ipAddress, userAgent, success, details)
	if err != nil {
		return fmt.Errorf("failed to create login attempt: %w", err)
	}

	return nil
}

func (r *LoginAttemptRepository) GetRecentLoginAttempts(ipAddress string, since time.Time) ([]models.LoginAttempt, error) {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.LoginAttempt.GetRecentLoginAttempts)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	rows, err := r.db.Query(query, ipAddress, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent login attempts: %w", err)
	}
	defer rows.Close()

	var attempts []models.LoginAttempt
	for rows.Next() {
		var attempt models.LoginAttempt
		err := rows.Scan(
			&attempt.ID,
			&attempt.IPAddress,
			&attempt.UserAgent,
			&attempt.Success,
			&attempt.AttemptAt,
			&attempt.Details,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan login attempt: %w", err)
		}
		attempts = append(attempts, attempt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate login attempts: %w", err)
	}

	return attempts, nil
}

func (r *LoginAttemptRepository) GetFailedLoginAttempts(ipAddress string, since time.Time) (int, error) {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.LoginAttempt.GetFailedLoginAttempts)
	if err != nil {
		return 0, fmt.Errorf("failed to get query: %w", err)
	}

	var count int
	err = r.db.QueryRow(query, ipAddress, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get failed login attempts count: %w", err)
	}

	return count, nil
}

func (r *LoginAttemptRepository) CleanupOldLoginAttempts(olderThan time.Time) error {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.LoginAttempt.CleanupOldLoginAttempts)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	result, err := r.db.Exec(query, olderThan)
	if err != nil {
		return fmt.Errorf("failed to cleanup old login attempts: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("Cleaned up %d old login attempts\n", rowsAffected)
	}

	return nil
}
