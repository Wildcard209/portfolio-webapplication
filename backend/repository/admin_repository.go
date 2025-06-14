package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Wildcard209/portfolio-webapplication/database"
	"github.com/Wildcard209/portfolio-webapplication/models"
)

type AdminRepository struct {
	db          *sql.DB
	queryLoader *database.QueryLoader
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	queryLoader, err := database.NewQueryLoader()
	if err != nil {
		// Log error but don't fail - fallback to inline queries if needed
		fmt.Printf("Warning: Failed to load queries: %v\n", err)
	}

	return &AdminRepository{
		db:          db,
		queryLoader: queryLoader,
	}
}

func (r *AdminRepository) GetAdminByUsername(username string) (*models.Admin, error) {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.GetAdminByUsername)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	admin := &models.Admin{}
	err = r.db.QueryRow(query, username).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.PasswordSalt,
		&admin.LastLogin,
		&admin.CurrentToken,
		&admin.TokenExpiration,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No admin found
		}
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}

	return admin, nil
}

func (r *AdminRepository) GetAdminByID(id int) (*models.Admin, error) {
	query := `
		SELECT id, username, password_hash, password_salt, last_login, current_token, token_expiration, created_at, updated_at
		FROM admins 
		WHERE id = $1
	`

	admin := &models.Admin{}
	err := r.db.QueryRow(query, id).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.PasswordSalt,
		&admin.LastLogin,
		&admin.CurrentToken,
		&admin.TokenExpiration,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No admin found
		}
		return nil, fmt.Errorf("failed to get admin by ID: %w", err)
	}

	return admin, nil
}

func (r *AdminRepository) CreateAdmin(username, passwordHash, passwordSalt string) (*models.Admin, error) {
	query := `
		INSERT INTO admins (username, password_hash, password_salt, created_at, updated_at) 
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, username, password_hash, password_salt, last_login, current_token, token_expiration, created_at, updated_at
	`

	admin := &models.Admin{}
	err := r.db.QueryRow(query, username, passwordHash, passwordSalt).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.PasswordSalt,
		&admin.LastLogin,
		&admin.CurrentToken,
		&admin.TokenExpiration,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create admin: %w", err)
	}

	return admin, nil
}

func (r *AdminRepository) UpdateAdminToken(id int, token string, expiration time.Time) error {
	query := `
		UPDATE admins 
		SET current_token = $1, token_expiration = $2, last_login = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	result, err := r.db.Exec(query, token, expiration, id)
	if err != nil {
		return fmt.Errorf("failed to update admin token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no admin found with ID %d", id)
	}

	return nil
}

func (r *AdminRepository) InvalidateAdminToken(id int) error {
	query := `
		UPDATE admins 
		SET current_token = NULL, token_expiration = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to invalidate admin token: %w", err)
	}

	return nil
}

func (r *AdminRepository) GetAdminByToken(token string) (*models.Admin, error) {
	query := `
		SELECT id, username, password_hash, password_salt, last_login, current_token, token_expiration, created_at, updated_at
		FROM admins 
		WHERE current_token = $1 AND token_expiration > CURRENT_TIMESTAMP
	`

	admin := &models.Admin{}
	err := r.db.QueryRow(query, token).Scan(
		&admin.ID,
		&admin.Username,
		&admin.PasswordHash,
		&admin.PasswordSalt,
		&admin.LastLogin,
		&admin.CurrentToken,
		&admin.TokenExpiration,
		&admin.CreatedAt,
		&admin.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get admin by token: %w", err)
	}

	return admin, nil
}

func (r *AdminRepository) CountAdmins() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM admins").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admins: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) CleanupExpiredTokens() error {
	query := `
		UPDATE admins 
		SET current_token = NULL, token_expiration = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE token_expiration < CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}
