package users_transport_http

import users_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/users/service"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10,max=16"`
	Role     string `json:"role" validate:"required"`
}

type PatchUserRequest struct {
	Email     *string `json:"email,omitempty" validate:"omitempty,email"`
	Password  *string `json:"password,omitempty" validate:"omitempty,min=6"`
	FullName  *string `json:"full_name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone     *string `json:"phone,omitempty" validate:"omitempty,min=10,max=16"`
	Role      *string `json:"role,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	IsBlocked *bool   `json:"is_blocked,omitempty"`
}

type BlockUserRequest struct {
	Reason *string `json:"reason,omitempty"`
}

type UsersResponse struct {
	Users []users_service.UserDTO `json:"users"`
}

type UserResponse struct {
	User users_service.UserDTO `json:"user"`
}
