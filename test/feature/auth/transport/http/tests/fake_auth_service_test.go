package auth_http_tests

import (
	"context"

	auth_service "github.com/pkpal-uhobp/fulfillment-app/internal/features/auth/service"
)

type fakeAuthService struct {
	registerFn          func(context.Context, auth_service.RegisterInput) (auth_service.UserDTO, auth_service.TokenPair, error)
	loginFn             func(context.Context, auth_service.LoginInput) (auth_service.UserDTO, auth_service.TokenPair, error)
	refreshFn           func(context.Context, auth_service.RefreshInput) (auth_service.TokenPair, error)
	logoutFn            func(context.Context, auth_service.LogoutInput) error
	getMeFn             func(context.Context, int64) (auth_service.UserDTO, error)
	verifyAccessTokenFn func(context.Context, string) (auth_service.AuthClaims, error)
}

func (f *fakeAuthService) Register(
	ctx context.Context,
	input auth_service.RegisterInput,
) (auth_service.UserDTO, auth_service.TokenPair, error) {
	if f.registerFn == nil {
		return auth_service.UserDTO{}, auth_service.TokenPair{}, nil
	}
	return f.registerFn(ctx, input)
}

func (f *fakeAuthService) Login(
	ctx context.Context,
	input auth_service.LoginInput,
) (auth_service.UserDTO, auth_service.TokenPair, error) {
	if f.loginFn == nil {
		return auth_service.UserDTO{}, auth_service.TokenPair{}, nil
	}
	return f.loginFn(ctx, input)
}

func (f *fakeAuthService) Refresh(
	ctx context.Context,
	input auth_service.RefreshInput,
) (auth_service.TokenPair, error) {
	if f.refreshFn == nil {
		return auth_service.TokenPair{}, nil
	}
	return f.refreshFn(ctx, input)
}

func (f *fakeAuthService) Logout(
	ctx context.Context,
	input auth_service.LogoutInput,
) error {
	if f.logoutFn == nil {
		return nil
	}
	return f.logoutFn(ctx, input)
}

func (f *fakeAuthService) GetMe(
	ctx context.Context,
	userID int64,
) (auth_service.UserDTO, error) {
	if f.getMeFn == nil {
		return auth_service.UserDTO{}, nil
	}
	return f.getMeFn(ctx, userID)
}

func (f *fakeAuthService) VerifyAccessToken(
	ctx context.Context,
	accessToken string,
) (auth_service.AuthClaims, error) {
	if f.verifyAccessTokenFn == nil {
		return auth_service.AuthClaims{}, nil
	}
	return f.verifyAccessTokenFn(ctx, accessToken)
}
