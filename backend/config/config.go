package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Config holds all configuration for our application
type Config struct {
	DB          *sql.DB
	MinioClient *minio.Client
	Port        string
}

// NewConfig creates a new configuration instance
func NewConfig() (*Config, error) {
	// Load environment variables
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Don't fail here, as environment variables might be set externally
	}

	config := &Config{
		Port: getEnv("PORT", "8080"),
	}

	// Check if we're in test mode (skip database connections)
	if os.Getenv("TEST_MODE") == "true" {
		log.Println("Running in test mode - skipping database connections")
		return config, nil
	}

	// Initialize database connection
	config.DB, err = initDB()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		// Continue without database for now
	}

	// Initialize MinIO client
	config.MinioClient, err = initMinio()
	if err != nil {
		log.Printf("Warning: Failed to initialize MinIO: %v", err)
		// Continue without MinIO for now
	}

	return config, nil
}

func initDB() (*sql.DB, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := getEnv("POSTGRES_HOST", "db") // Docker Compose service name for PostgreSQL
	dbPort := getEnv("POSTGRES_PORT", "5432")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL!")
	return db, nil
}

func initMinio() (*minio.Client, error) {
	minioUser := os.Getenv("MINIO_ROOT_USER")
	minioPassword := os.Getenv("MINIO_ROOT_PASSWORD")
	minioEndpoint := getEnv("MINIO_ENDPOINT", "localhost:9000")

	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioPassword, ""),
		Secure: false, // Use true if MinIO supports HTTPS
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	log.Println("Successfully connected to MinIO!")
	log.Printf("MinIO User: %s", minioUser)
	return minioClient, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Close closes all connections
func (c *Config) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
