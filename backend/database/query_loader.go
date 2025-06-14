package database

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed queries/**/*.sql
var queryFiles embed.FS

type QueryLoader struct {
	queries map[string]string
}

func NewQueryLoader() (*QueryLoader, error) {
	loader := &QueryLoader{
		queries: make(map[string]string),
	}

	err := loader.loadQueries()
	if err != nil {
		return nil, fmt.Errorf("failed to load queries: %w", err)
	}

	return loader, nil
}

func (ql *QueryLoader) loadQueries() error {
	return fs.WalkDir(queryFiles, "queries", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		content, err := queryFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read query file %s: %w", path, err)
		}

		key := ql.generateQueryKey(path)
		ql.queries[key] = string(content)

		return nil
	})
}

func (ql *QueryLoader) generateQueryKey(path string) string {
	key := strings.TrimPrefix(path, "queries/")
	key = strings.TrimSuffix(key, ".sql")

	key = strings.ReplaceAll(key, "/", ".")

	return key
}

func (ql *QueryLoader) GetQuery(key string) (string, error) {
	query, exists := ql.queries[key]
	if !exists {
		return "", fmt.Errorf("query not found: %s", key)
	}
	return query, nil
}

func (ql *QueryLoader) ListQueries() []string {
	keys := make([]string, 0, len(ql.queries))
	for key := range ql.queries {
		keys = append(keys, key)
	}
	return keys
}

type AdminQueries struct {
	GetAdminByUsername   string
	GetAdminByID         string
	CreateAdmin          string
	UpdateAdminToken     string
	InvalidateAdminToken string
	GetAdminByToken      string
	CountAdmins          string
	CleanupExpiredTokens string
}

type LoginAttemptQueries struct {
	CreateLoginAttempt      string
	GetRecentLoginAttempts  string
	GetFailedLoginAttempts  string
	CleanupOldLoginAttempts string
}

var QueryKeys = struct {
	Admin        AdminQueries
	LoginAttempt LoginAttemptQueries
}{
	Admin: AdminQueries{
		GetAdminByUsername:   "admin.get_admin_by_username",
		GetAdminByID:         "admin.get_admin_by_id",
		CreateAdmin:          "admin.create_admin",
		UpdateAdminToken:     "admin.update_admin_token",
		InvalidateAdminToken: "admin.invalidate_admin_token",
		GetAdminByToken:      "admin.get_admin_by_token",
		CountAdmins:          "admin.count_admins",
		CleanupExpiredTokens: "admin.cleanup_expired_tokens",
	},
	LoginAttempt: LoginAttemptQueries{
		CreateLoginAttempt:      "login_attempts.create_login_attempt",
		GetRecentLoginAttempts:  "login_attempts.get_recent_login_attempts",
		GetFailedLoginAttempts:  "login_attempts.get_failed_login_attempts",
		CleanupOldLoginAttempts: "login_attempts.cleanup_old_login_attempts",
	},
}
