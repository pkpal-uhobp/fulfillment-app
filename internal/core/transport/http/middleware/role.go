package core_http_middleware

import (
	"net/http"
	"strings"
)

func RequireRoles(verifier AccessTokenVerifier) func(roles ...string) Middleware {
	return func(roles ...string) Middleware {
		allowedRoles := normalizeRoles(roles...)

		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if len(allowedRoles) == 0 {
					next.ServeHTTP(w, r)
					return
				}

				user, err := CurrentUserFromContext(r.Context())
				if err != nil {
					claims, ok := verifyRequestAccessToken(w, r, verifier)
					if !ok {
						return
					}

					user = CurrentUser{
						ID:   claims.UserID,
						Role: claims.Role,
						JTI:  claims.JTI,
					}

					r = r.WithContext(WithUser(r.Context(), user))
				}

				userRole := strings.ToLower(strings.TrimSpace(user.Role))

				if _, ok := allowedRoles[userRole]; !ok {
					writeMiddlewareError(
						w,
						http.StatusForbidden,
						"access denied",
					)
					return
				}

				next.ServeHTTP(w, r)
			})
		}
	}
}

func normalizeRoles(roles ...string) map[string]struct{} {
	normalized := make(map[string]struct{}, len(roles))

	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role == "" {
			continue
		}

		normalized[role] = struct{}{}
	}

	return normalized
}
