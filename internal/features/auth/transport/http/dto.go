package auth_transport_http

import auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type RegisterResponse struct {
	User   auth_service.UserDTO   `json:"user"`
	Tokens auth_service.TokenPair `json:"tokens"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceID string `json:"device_id"`
}

type LoginResponse struct {
	User   auth_service.UserDTO   `json:"user"`
	Tokens auth_service.TokenPair `json:"tokens"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	Tokens auth_service.TokenPair `json:"tokens"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type MeResponse struct {
	User auth_service.UserDTO `json:"user"`
}
