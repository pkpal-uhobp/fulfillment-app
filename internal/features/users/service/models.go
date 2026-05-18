package users_service

type UserDTO struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone,omitempty"`
	Role      string `json:"role"`
	IsActive  bool   `json:"is_active"`
	IsBlocked bool   `json:"is_blocked"`
}

type UserFilter struct {
	Role      string
	IsActive  *bool
	IsBlocked *bool
	Search    string
	Page      int
	Limit     int
}

type CreateUserInput struct {
	Email    string
	Password string
	FullName string
	Phone    string
	Role     string
}

type PatchUserInput struct {
	Email     *string
	Password  *string
	FullName  *string
	Phone     *string
	Role      *string
	IsActive  *bool
	IsBlocked *bool
}

type BlockUserInput struct {
	Reason *string
}
