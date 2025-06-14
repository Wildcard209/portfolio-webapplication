package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret   []byte
	tokenExpiry time.Duration
	pepper      string
}

func NewAuthService(jwtSecret string, tokenExpiry time.Duration) *AuthService {
	// Get pepper from environment or use default for development
	pepper := os.Getenv("PASSWORD_PEPPER")
	if pepper == "" {
		pepper = "default-pepper-change-in-production"
	}

	return &AuthService{
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
		pepper:      pepper,
	}
}

type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateSalt generates a random salt for password hashing
func (s *AuthService) GenerateSalt() (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPasswordWithSalt hashes a password with salt and pepper using a secure method
func (s *AuthService) HashPasswordWithSalt(password, salt string) (string, error) {
	// Use SHA-256 to hash the combined password+salt+pepper to ensure it's always under bcrypt's 72-byte limit
	// This is a common pattern for handling long passwords/salts with bcrypt
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hasher.Write([]byte(salt))
	hasher.Write([]byte(s.pepper))

	// Get the SHA-256 hash and encode it as base64 for consistent length
	hashedInput := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Now use bcrypt on the pre-hashed input (which is always 44 bytes in base64)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(hashedInput), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// HashPassword generates salt and hashes password (for backward compatibility)
func (s *AuthService) HashPassword(password string) (string, error) {
	salt, err := s.GenerateSalt()
	if err != nil {
		return "", err
	}
	return s.HashPasswordWithSalt(password, salt)
}

// VerifyPassword verifies a password against its hash using salt and pepper
func (s *AuthService) VerifyPassword(hashedPassword, password, salt string) error {
	// Use the same SHA-256 pre-hashing method as in HashPasswordWithSalt
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hasher.Write([]byte(salt))
	hasher.Write([]byte(s.pepper))

	// Get the SHA-256 hash and encode it as base64
	hashedInput := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	// Compare with bcrypt
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(hashedInput))
}

// GenerateToken generates a JWT token for the user
func (s *AuthService) GenerateToken(userID int, username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(s.tokenExpiry)

	claims := &CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expirationTime, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			return nil, errors.New("token has expired")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) RefreshToken(oldTokenString string) (string, time.Time, error) {
	claims, err := s.ValidateToken(oldTokenString)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid token for refresh: %w", err)
	}

	return s.GenerateToken(claims.UserID, claims.Username)
}

func (s *AuthService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("invalid authorization header format")
	}

	return authHeader[len(bearerPrefix):], nil
}
