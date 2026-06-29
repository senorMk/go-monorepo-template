package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/GITHUB_USERNAME/APP_NAME/internal/httpx"
)

type ctxKey int

const userCtxKey ctxKey = iota

// User is the authenticated principal extracted from a valid access token.
type User struct {
	ID    string
	Email string
}

// Middleware returns net/http middleware that requires a valid Bearer access
// token and stores the resulting User in the request context.
func Middleware(jwt *JWTManager) func(http.Handler) http.Handler {
	const prefix = "Bearer "
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, prefix) {
				httpx.WriteError(w, httpx.Unauthorized("Missing or expired token"))
				return
			}
			claims, err := jwt.Parse(strings.TrimPrefix(header, prefix))
			if err != nil {
				httpx.WriteError(w, httpx.Unauthorized("Missing or expired token"))
				return
			}
			ctx := context.WithValue(r.Context(), userCtxKey, User{ID: claims.Subject, Email: claims.Email})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CurrentUser returns the authenticated user from the request context.
func CurrentUser(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(userCtxKey).(User)
	return u, ok
}
