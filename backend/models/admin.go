package models

import (
	"database/sql/driver"
	"time"
)

type Admin struct {
	ID              int       `json:"id" db:"id"`
	Username        string    `json:"username" db:"username"`
	PasswordHash    string    `json:"-" db:"password_hash"`
	PasswordSalt    string    `json:"-" db:"password_salt"`
	LastLogin       NullTime  `json:"last_login" db:"last_login"`
	CurrentToken    *string   `json:"-" db:"current_token"`
	TokenExpiration NullTime  `json:"-" db:"token_expiration"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type LoginAttempt struct {
	ID        int       `json:"id" db:"id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Success   bool      `json:"success" db:"success"`
	AttemptAt time.Time `json:"attempt_at" db:"attempt_at"`
	Details   *string   `json:"details,omitempty" db:"details"`
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Time = value.(time.Time)
	nt.Valid = true
	return nil
}

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"password"`
}

type LoginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" example:"2023-12-31T23:59:59Z"`
	User      AdminUser `json:"user"`
}

type AdminUser struct {
	ID        int      `json:"id" example:"1"`
	Username  string   `json:"username" example:"admin"`
	LastLogin NullTime `json:"last_login"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid credentials"`
	Message string `json:"message,omitempty" example:"Username or password is incorrect"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}
