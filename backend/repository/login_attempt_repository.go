package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/models"
)

type LoginAttemptRepository struct {
	db *sql.DB
}

func NewLoginAttemptRepository(db *sql.DB) *LoginAttemptRepository {
	return &LoginAttemptRepository{db: db}
}

func (r *LoginAttemptRepository) CreateLoginAttempt(ipAddress, userAgent string, success bool, details *string) error {
	query := `
		INSERT INTO login_attempts (ip_address, user_agent, success, attempt_at, details)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, $4)
	`

	_, err := r.db.Exec(query, ipAddress, userAgent, success, details)
	if err != nil {
		return fmt.Errorf("failed to create login attempt: %w", err)
	}

	return nil
}

func (r *LoginAttemptRepository) GetRecentLoginAttempts(ipAddress string, since time.Time) ([]models.LoginAttempt, error) {
	query := `
		SELECT id, ip_address, user_agent, success, attempt_at, details
		FROM login_attempts 
		WHERE ip_address = $1 AND attempt_at >= $2
		ORDER BY attempt_at DESC
	`

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
	query := `
		SELECT COUNT(*) 
		FROM login_attempts 
		WHERE ip_address = $1 AND success = FALSE AND attempt_at >= $2
	`

	var count int
	err := r.db.QueryRow(query, ipAddress, since).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get failed login attempts count: %w", err)
	}

	return count, nil
}

func (r *LoginAttemptRepository) CleanupOldLoginAttempts(olderThan time.Time) error {
	query := `DELETE FROM login_attempts WHERE attempt_at < $1`

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
