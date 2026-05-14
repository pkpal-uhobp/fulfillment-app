package domain

import (
	"fmt"
	"regexp"
	"strings"

	core_errors "github.com/pkpal-uhobp/fulfillment-app/internal/core/errors"
)

var (
	emailRegexp = regexp.MustCompile(`^[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,63}$`)
	phoneRegexp = regexp.MustCompile(`^[+]?[0-9]{10,15}$`)
)

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

func NewUser(
	email string,
	passwordHash string,
	fullName string,
	phone *string,
	role Role,
) (User, error) {
	user := User{
		Email:        normalizeEmail(email),
		PasswordHash: strings.TrimSpace(passwordHash),
		FullName:     strings.TrimSpace(fullName),
		Phone:        normalizePhone(phone),
		Role:         role,
		IsActive:     true,
		IsBlocked:    false,
	}

	if err := user.Validate(); err != nil {
		return User{}, err
	}

	return user, nil
}

func (u User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("%w: email is required", core_errors.ErrInvalidArgument)
	}

	if !emailRegexp.MatchString(strings.ToUpper(u.Email)) {
		return fmt.Errorf("%w: invalid email", core_errors.ErrInvalidArgument)
	}

	if u.PasswordHash == "" {
		return fmt.Errorf("%w: password hash is required", core_errors.ErrInvalidArgument)
	}

	if u.FullName == "" {
		return fmt.Errorf("%w: full name is required", core_errors.ErrInvalidArgument)
	}

	if u.Phone != nil && !phoneRegexp.MatchString(*u.Phone) {
		return fmt.Errorf("%w: invalid phone", core_errors.ErrInvalidArgument)
	}

	if !u.Role.IsValid() {
		return fmt.Errorf("%w: invalid user role", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u User) IsAvailable() bool {
	return u.IsActive && !u.IsBlocked
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizePhone(phone *string) *string {
	if phone == nil {
		return nil
	}

	value := strings.TrimSpace(*phone)
	if value == "" {
		return nil
	}

	return &value
}
