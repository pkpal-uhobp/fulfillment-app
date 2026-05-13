package domain

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	FullName     string
	Phone        *string
	Role         Role
	IsActive     bool
	IsBlocked    bool
}

func (u User) IsAvailable() bool {
	return u.IsActive && !u.IsBlocked
}
