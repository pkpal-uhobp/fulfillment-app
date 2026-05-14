package auth_transport_http

import auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10,max=16"`
}

type RegisterResponse struct {
	User   auth_service.UserDTO   `json:"user"`
	Tokens auth_service.TokenPair `json:"tokens"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	DeviceID string `json:"device_id,omitempty" validate:"omitempty,uuid"`
}

type LoginResponse struct {
	User   auth_service.UserDTO   `json:"user"`
	Tokens auth_service.TokenPair `json:"tokens"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	Tokens auth_service.TokenPair `json:"tokens"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type MeResponse struct {
	User auth_service.UserDTO `json:"user"`
}
