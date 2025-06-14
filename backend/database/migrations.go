package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migration struct {
	Version  string
	Filename string
	SQL      string
}

func GetMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to read migration directory: %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version := strings.TrimSuffix(entry.Name(), ".sql")

		sqlContent, err := migrationFiles.ReadFile(filepath.Join("migrations", entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version:  version,
			Filename: entry.Name(),
			SQL:      string(sqlContent),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func RunMigrations(db *sql.DB) error {
	migrationTableSQL := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) UNIQUE NOT NULL,
			filename VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := db.Exec(migrationTableSQL); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations, err := GetMigrations()
	if err != nil {
		return fmt.Errorf("failed to get migrations: %w", err)
	}

	for _, migration := range migrations {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = $1", migration.Version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", migration.Version, err)
		}

		if count > 0 {
			log.Printf("Migration %s already applied, skipping", migration.Version)
			continue
		}

		log.Printf("Applying migration: %s (%s)", migration.Version, migration.Filename)
		if _, err := db.Exec(migration.SQL); err != nil {
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		if _, err := db.Exec("INSERT INTO migrations (version, filename) VALUES ($1, $2)", migration.Version, migration.Filename); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		log.Printf("Migration %s applied successfully", migration.Version)
	}

	return nil
}
