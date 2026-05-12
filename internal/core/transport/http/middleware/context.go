package core_http_middleware

import (
	"context"
	"fmt"
)

type contextKey string

const (
	contextKeyUserID    contextKey = "user_id"
	contextKeyRole      contextKey = "role"
	contextKeyTokenJTI  contextKey = "token_jti"
	contextKeyRequestID contextKey = "request_id"
)

const (
	RoleClient = "client"
	RoleLogist = "logist"
	RoleWorker = "worker"
	RoleAdmin  = "admin"
)

type CurrentUser struct {
	ID   int64
	Role string
	JTI  string
}

func WithUser(ctx context.Context, user CurrentUser) context.Context {
	ctx = context.WithValue(ctx, contextKeyUserID, user.ID)
	ctx = context.WithValue(ctx, contextKeyRole, user.Role)
	ctx = context.WithValue(ctx, contextKeyTokenJTI, user.JTI)

	return ctx
}

func CurrentUserFromContext(ctx context.Context) (CurrentUser, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return CurrentUser{}, err
	}

	role, err := RoleFromContext(ctx)
	if err != nil {
		return CurrentUser{}, err
	}

	jti, _ := TokenJTIFromContext(ctx)

	return CurrentUser{
		ID:   userID,
		Role: role,
		JTI:  jti,
	}, nil
}

func UserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(contextKeyUserID).(int64)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}

	return userID, nil
}

func RoleFromContext(ctx context.Context) (string, error) {
	role, ok := ctx.Value(contextKeyRole).(string)
	if !ok {
		return "", fmt.Errorf("role not found in context")
	}

	return role, nil
}

func TokenJTIFromContext(ctx context.Context) (string, error) {
	jti, ok := ctx.Value(contextKeyTokenJTI).(string)
	if !ok {
		return "", fmt.Errorf("token jti not found in context")
	}

	return jti, nil
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value(contextKeyRequestID).(string)
	if !ok {
		return "", fmt.Errorf("request id not found in context")
	}

	return requestID, nil
}
