package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := "db" // Docker Compose service name for PostgreSQL
	dbPort := "5432"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Database connection error: %v", err)
		}
	}(db)

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	} else {
		log.Println("Successfully connected to PostgreSQL!")
	}

	minioUser := os.Getenv("MINIO_ROOT_USER")
	minioPassword := os.Getenv("MINIO_ROOT_PASSWORD")
	minioEndpoint := "localhost:9000" // MinIO runs locally in this setup

	_, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioPassword, ""),
		Secure: false, // Use true if MinIO supports HTTPS
	})
	if err != nil {
		log.Fatalf("Failed to connect to MinIO: %v", err)
	}

	log.Println("Successfully connected to MinIO!")
	log.Printf("MinIO User: %s\n", minioUser)

	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Write JSON response
		response := `{"message": "Hello from Go backend 2!"}`

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(response))
		if err != nil {
			return
		}
	})

	port := ":8080"
	log.Printf("Server is running on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
