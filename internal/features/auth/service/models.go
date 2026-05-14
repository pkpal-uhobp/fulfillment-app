package auth_service

import "time"

type UserDTO struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`
	Role     string `json:"role"`
}

type TokenPair struct {
	AccessToken      string    `json:"access_token"`
	RefreshToken     string    `json:"refresh_token"`
	AccessExpiresAt  time.Time `json:"access_expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	DeviceID         string    `json:"device_id"`
}

type RegisterInput struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type LoginInput struct {
	Email    string
	Password string
	DeviceID string
}

type RefreshInput struct {
	RefreshToken string
}

type LogoutInput struct {
	RefreshToken string
}

type AuthClaims struct {
	UserID   int64
	Role     string
	JTI      string
	DeviceID string
}
