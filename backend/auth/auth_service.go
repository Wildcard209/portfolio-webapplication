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
	jwtSecret          []byte
	tokenExpiry        time.Duration
	refreshTokenExpiry time.Duration
	pepper             string
}

func NewAuthService(jwtSecret string, tokenExpiry time.Duration) *AuthService {
	pepper := os.Getenv("PASSWORD_PEPPER")
	if pepper == "" {
		pepper = "default-pepper-change-in-production"
	}

	return &AuthService{
		jwtSecret:          []byte(jwtSecret),
		tokenExpiry:        tokenExpiry,
		refreshTokenExpiry: 7 * 24 * time.Hour, // 7 days
		pepper:             pepper,
	}
}

type CustomClaims struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

func (s *AuthService) GenerateSalt() (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func (s *AuthService) HashPasswordWithSalt(password, salt string) (string, error) {
	// Use SHA-256 to hash the combined password+salt+pepper to ensure it's always under bcrypt's 72-byte limit
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hasher.Write([]byte(salt))
	hasher.Write([]byte(s.pepper))

	hashedInput := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(hashedInput), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func (s *AuthService) HashPassword(password string) (string, error) {
	salt, err := s.GenerateSalt()
	if err != nil {
		return "", err
	}
	return s.HashPasswordWithSalt(password, salt)
}

func (s *AuthService) VerifyPassword(hashedPassword, password, salt string) error {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hasher.Write([]byte(salt))
	hasher.Write([]byte(s.pepper))

	hashedInput := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(hashedInput))
}

func (s *AuthService) GenerateTokenPair(userID int, username string) (*TokenPair, error) {
	accessExpirationTime := time.Now().Add(s.tokenExpiry)
	refreshExpirationTime := time.Now().Add(s.refreshTokenExpiry)

	// Generate access token
	accessClaims := &CustomClaims{
		UserID:    userID,
		Username:  username,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Generate refresh token
	refreshClaims := &CustomClaims{
		UserID:    userID,
		Username:  username,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:      accessTokenString,
		RefreshToken:     refreshTokenString,
		AccessExpiresAt:  accessExpirationTime,
		RefreshExpiresAt: refreshExpirationTime,
	}, nil
}

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

func (s *AuthService) ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "access" {
		return nil, errors.New("invalid token type: expected access token")
	}

	return claims, nil
}

func (s *AuthService) ValidateRefreshToken(tokenString string) (*CustomClaims, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("invalid token type: expected refresh token")
	}

	return claims, nil
}

func (s *AuthService) RefreshTokenPair(refreshTokenString string) (*TokenPair, error) {
	claims, err := s.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	return s.GenerateTokenPair(claims.UserID, claims.Username)
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
