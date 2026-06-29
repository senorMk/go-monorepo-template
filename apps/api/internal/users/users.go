// Package users holds the shared user response shape used by the auth and
// profile handlers, keeping the JSON contract identical across endpoints.
package users

import (
	"time"

	"github.com/GITHUB_USERNAME/APP_NAME/internal/db/sqlc"
)

// Response is the public JSON representation of a user.
type Response struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	DisplayName *string `json:"displayName"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

// Format converts a database user into its public representation.
func Format(u sqlc.User) Response {
	return Response{
		ID:          u.ID.String(),
		Email:       u.Email,
		DisplayName: u.DisplayName,
		CreatedAt:   ISOTime(u.CreatedAt.Time),
		UpdatedAt:   ISOTime(u.UpdatedAt.Time),
	}
}

// ISOTime formats a time as an ISO-8601 string with millisecond precision in
// UTC, matching JavaScript's Date.toISOString().
func ISOTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
