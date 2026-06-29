package auth

import "github.com/GITHUB_USERNAME/APP_NAME/internal/users"

// Request payloads ----------------------------------------------------------

type signupRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8"`
	DisplayName *string `json:"displayName" validate:"omitempty"`
}

type signinRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type resetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type confirmEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// Response payloads ---------------------------------------------------------

// AuthResult is returned by signup and signin. On signup the token fields are
// null and RequiresEmailConfirmation is true.
type AuthResult struct {
	User                      *users.Response `json:"user"`
	AccessToken               *string         `json:"accessToken"`
	RefreshToken              *string         `json:"refreshToken"`
	ExpiresAt                 *string         `json:"expiresAt"`
	RequiresEmailConfirmation bool            `json:"requiresEmailConfirmation"`
}

// SessionResult is returned by confirm-email: a user plus a fresh token pair.
type SessionResult struct {
	User         *users.Response `json:"user"`
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"refreshToken"`
	ExpiresAt    string          `json:"expiresAt"`
}

// RefreshResult is returned by refresh: a rotated token pair, no user.
type RefreshResult struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    string `json:"expiresAt"`
}
