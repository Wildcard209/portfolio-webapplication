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
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get admin by username: %w", err)
	}

	return admin, nil
}

func (r *AdminRepository) GetAdminByID(id int) (*models.Admin, error) {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.GetAdminByID)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	admin := &models.Admin{}
	err = r.db.QueryRow(query, id).Scan(
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
		return nil, fmt.Errorf("failed to get admin by ID: %w", err)
	}

	return admin, nil
}

func (r *AdminRepository) CreateAdmin(username, passwordHash, passwordSalt string) (*models.Admin, error) {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.CreateAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	admin := &models.Admin{}
	err = r.db.QueryRow(query, username, passwordHash, passwordSalt).Scan(
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
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.UpdateAdminToken)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

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
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.InvalidateAdminToken)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to invalidate admin token: %w", err)
	}

	return nil
}

func (r *AdminRepository) GetAdminByToken(token string) (*models.Admin, error) {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.GetAdminByToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get query: %w", err)
	}

	admin := &models.Admin{}
	err = r.db.QueryRow(query, token).Scan(
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
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.CountAdmins)
	if err != nil {
		return 0, fmt.Errorf("failed to get query: %w", err)
	}

	var count int
	err = r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admins: %w", err)
	}
	return count, nil
}

func (r *AdminRepository) CleanupExpiredTokens() error {
	query, err := r.queryLoader.GetQuery(database.QueryKeys.Admin.CleanupExpiredTokens)
	if err != nil {
		return fmt.Errorf("failed to get query: %w", err)
	}

	_, err = r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	return nil
}
