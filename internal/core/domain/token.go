package domain

import (
	"time"

	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

func (t TokenType) String() string {
	return string(t)
}

func (t TokenType) IsValid() bool {
	switch t {
	case TokenTypeAccess, TokenTypeRefresh:
		return true
	default:
		return false
	}
}

type IssuedToken struct {
	ID        int64
	UserID    int64
	JTI       uuid.UUID
	TokenType TokenType
	DeviceID  uuid.UUID
	Revoked   bool
	ExpiresAt time.Time
}
