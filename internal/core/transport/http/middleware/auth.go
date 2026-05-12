package core_http_middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type AccessTokenClaims struct {
	UserID int64
	Role   string
	JTI    string
}

type AccessTokenVerifier interface {
	VerifyAccessToken(ctx context.Context, token string) (*AccessTokenClaims, error)
}

func Auth(verifier AccessTokenVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := verifyRequestAccessToken(w, r, verifier)
			if !ok {
				return
			}

			ctx := WithUser(r.Context(), CurrentUser{
				ID:   claims.UserID,
				Role: claims.Role,
				JTI:  claims.JTI,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyRequestAccessToken(
	w http.ResponseWriter,
	r *http.Request,
	verifier AccessTokenVerifier,
) (*AccessTokenClaims, bool) {
	token, ok := bearerTokenFromRequest(w, r)
	if !ok {
		return nil, false
	}

	claims, err := verifier.VerifyAccessToken(r.Context(), token)
	if err != nil {
		writeMiddlewareError(
			w,
			http.StatusUnauthorized,
			"invalid access token",
		)
		return nil, false
	}

	if claims == nil {
		writeMiddlewareError(
			w,
			http.StatusUnauthorized,
			"invalid access token claims",
		)
		return nil, false
	}

	if claims.UserID <= 0 || strings.TrimSpace(claims.Role) == "" {
		writeMiddlewareError(
			w,
			http.StatusUnauthorized,
			"invalid access token claims",
		)
		return nil, false
	}

	return claims, true
}

func bearerTokenFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeMiddlewareError(
			w,
			http.StatusUnauthorized,
			"missing authorization header",
		)
		return "", false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		writeMiddlewareError(
			w,
			http.StatusUnauthorized,
			"invalid authorization header",
		)
		return "", false
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		writeMiddlewareError(
			w,
			http.StatusUnauthorized,
			"empty bearer token",
		)
		return "", false
	}

	return token, true
}

func writeMiddlewareError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}
