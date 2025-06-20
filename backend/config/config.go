package config

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	DB          *sql.DB
	MinioClient *minio.Client
	Port        string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type SanitizedError struct {
	Operation string
	Cause     string
}

func (e *SanitizedError) Error() string {
	return fmt.Sprintf("database operation failed: %s - %s", e.Operation, e.Cause)
}

func sanitizeError(operation string, err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errMsg = strings.ReplaceAll(errMsg, "password=", "password=***")

	if strings.Contains(errMsg, "postgres://") {
		parts := strings.Split(errMsg, "postgres://")
		if len(parts) > 1 {
			connStr := parts[1]
			if idx := strings.IndexAny(connStr, " \t\n\"'"); idx != -1 {
				connStr = connStr[:idx]
			}
			sanitized := "postgres://***:***@***:****/***"
			errMsg = strings.ReplaceAll(errMsg, "postgres://"+connStr, sanitized)
		}
	}

	return &SanitizedError{
		Operation: operation,
		Cause:     errMsg,
	}
}

func NewConfig() (*Config, error) {
	var err error

	config := &Config{
		Port: getEnv("PORT", "8080"),
	}

	if os.Getenv("TEST_MODE") == "true" {
		log.Println("Running in test mode - skipping database connections")
		return config, nil
	}

	config.DB, err = initDB()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}

	config.MinioClient, err = initMinio()
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO: %v", err)
	}

	return config, nil
}

func initDB() (*sql.DB, error) {
	dbConfig := DatabaseConfig{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
		Host:     getEnv("POSTGRES_HOST", "db"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
	}

	if dbConfig.User == "" || dbConfig.Password == "" || dbConfig.Database == "" {
		return nil, sanitizeError("configuration validation",
			fmt.Errorf("missing required database configuration: POSTGRES_USER, POSTGRES_PASSWORD, and POSTGRES_DB must be set"))
	}

	params := url.Values{}
	params.Set("sslmode", "disable")
	params.Set("application_name", "portfolio-webapp")

	actualDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		url.QueryEscape(dbConfig.User),
		url.QueryEscape(dbConfig.Password),
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
		params.Encode())

	db, err := sql.Open("pgx", actualDSN)
	if err != nil {
		return nil, sanitizeError("database connection", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, sanitizeError("database ping", err)
	}

	log.Printf("Successfully connected to PostgreSQL database: %s@%s:%s/%s",
		"[USER_REDACTED]", dbConfig.Host, dbConfig.Port, dbConfig.Database)

	return db, nil
}

func initMinio() (*minio.Client, error) {
	minioConfig := MinioConfig{
		AccessKey: os.Getenv("MINIO_ROOT_USER"),
		SecretKey: os.Getenv("MINIO_ROOT_PASSWORD"),
		Endpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
		UseSSL:    false,
	}

	if minioConfig.AccessKey == "" || minioConfig.SecretKey == "" {
		return nil, sanitizeError("minio configuration validation",
			fmt.Errorf("missing required MinIO configuration: MINIO_ROOT_USER and MINIO_ROOT_PASSWORD must be set"))
	}

	minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Secure: minioConfig.UseSSL,
	})
	if err != nil {
		return nil, sanitizeError("minio client creation", err)
	}

	log.Printf("Successfully connected to MinIO storage: [USER_REDACTED]@%s", minioConfig.Endpoint)

	return minioClient, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *Config) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
