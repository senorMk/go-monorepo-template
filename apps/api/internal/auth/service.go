// Package auth implements signup, signin, token refresh, email confirmation
// and password reset. Opaque tokens are stored as SHA-256 hashes; passwords
// are hashed with bcrypt; access tokens are short-lived JWTs.
package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/GITHUB_USERNAME/APP_NAME/internal/db/sqlc"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/email"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/httpx"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/users"
)

const (
	emailConfirmationTTL = 24 * time.Hour
	passwordResetTTL     = 1 * time.Hour
)

// Service holds the dependencies for authentication operations.
type Service struct {
	q          *sqlc.Queries
	pool       *pgxpool.Pool
	jwt        *JWTManager
	email      *email.Service
	refreshTTL time.Duration
}

// NewService builds an auth Service.
func NewService(q *sqlc.Queries, pool *pgxpool.Pool, jwt *JWTManager, email *email.Service, refreshTTL time.Duration) *Service {
	return &Service{q: q, pool: pool, jwt: jwt, email: email, refreshTTL: refreshTTL}
}

// withTx runs fn inside a database transaction, committing on success and
// rolling back on any error, so multi-statement flows are atomic.
func (s *Service) withTx(ctx context.Context, fn func(q *sqlc.Queries) error) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if err := fn(s.q.WithTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// issueRefreshToken creates and stores a new opaque refresh token, returning
// the plaintext token to hand back to the client.
func issueRefreshToken(ctx context.Context, q *sqlc.Queries, userID uuid.UUID, ttl time.Duration) (string, error) {
	plain := newOpaqueToken()
	err := q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: sha256Hex(plain),
		ExpiresAt: timestamptz(time.Now().Add(ttl)),
	})
	if err != nil {
		return "", err
	}
	return plain, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation
}

// Signup registers a new user and sends an email-confirmation link. It does
// not return tokens: the account must be confirmed before signin.
func (s *Service) Signup(ctx context.Context, req signupRequest) (AuthResult, error) {
	if _, err := s.q.GetUserByEmail(ctx, req.Email); err == nil {
		return AuthResult{}, httpx.Conflict("Email already registered")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return AuthResult{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, err
	}

	var user sqlc.User
	var confirmToken string
	err = s.withTx(ctx, func(q *sqlc.Queries) error {
		var err error
		user, err = q.CreateUser(ctx, sqlc.CreateUserParams{
			Email:        req.Email,
			PasswordHash: string(hash),
			DisplayName:  req.DisplayName,
		})
		if err != nil {
			return err
		}
		confirmToken = newOpaqueToken()
		return q.CreateEmailConfirmationToken(ctx, sqlc.CreateEmailConfirmationTokenParams{
			UserID:    user.ID,
			TokenHash: sha256Hex(confirmToken),
			ExpiresAt: timestamptz(time.Now().Add(emailConfirmationTTL)),
		})
	})
	if err != nil {
		// A concurrent signup with the same email loses the unique-constraint race.
		if isUniqueViolation(err) {
			return AuthResult{}, httpx.Conflict("Email already registered")
		}
		return AuthResult{}, err
	}

	s.email.SendEmailConfirmation(user.Email, confirmToken)

	resp := users.Format(user)
	return AuthResult{User: &resp, RequiresEmailConfirmation: true}, nil
}

// Signin authenticates a confirmed user and issues an access/refresh pair.
func (s *Service) Signin(ctx context.Context, req signinRequest) (AuthResult, error) {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return AuthResult{}, httpx.Unauthorized("Invalid credentials")
	} else if err != nil {
		return AuthResult{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return AuthResult{}, httpx.Unauthorized("Invalid credentials")
	}
	if !user.EmailVerifiedAt.Valid {
		return AuthResult{}, httpx.Forbidden("Email not confirmed")
	}

	access, expiresAt, err := s.jwt.Issue(user.ID.String(), user.Email)
	if err != nil {
		return AuthResult{}, err
	}
	refresh, err := issueRefreshToken(ctx, s.q, user.ID, s.refreshTTL)
	if err != nil {
		return AuthResult{}, err
	}

	resp := users.Format(user)
	iso := users.ISOTime(expiresAt)
	return AuthResult{
		User:                      &resp,
		AccessToken:               &access,
		RefreshToken:              &refresh,
		ExpiresAt:                 &iso,
		RequiresEmailConfirmation: false,
	}, nil
}

// Signout revokes all of the user's active refresh tokens.
func (s *Service) Signout(ctx context.Context, userID uuid.UUID) error {
	return s.q.RevokeAllUserRefreshTokens(ctx, userID)
}

// Refresh rotates a refresh token, returning a new access/refresh pair. If a
// previously-revoked token is replayed, the entire token family is revoked
// (reuse detection).
func (s *Service) Refresh(ctx context.Context, refreshToken string) (RefreshResult, error) {
	stored, err := s.q.GetRefreshTokenByHash(ctx, sha256Hex(refreshToken))
	if errors.Is(err, pgx.ErrNoRows) {
		return RefreshResult{}, httpx.Unauthorized("Refresh token expired or revoked")
	} else if err != nil {
		return RefreshResult{}, err
	}

	// Replaying an already-revoked token signals theft: kill the whole family.
	if stored.RevokedAt.Valid {
		if err := s.q.RevokeAllUserRefreshTokens(ctx, stored.UserID); err != nil {
			slog.Error("revoke token family after reuse", "err", err, "user_id", stored.UserID)
		}
		return RefreshResult{}, httpx.Unauthorized("Refresh token expired or revoked")
	}
	if stored.ExpiresAt.Time.Before(time.Now()) {
		return RefreshResult{}, httpx.Unauthorized("Refresh token expired or revoked")
	}

	var user sqlc.User
	var newRefresh string
	err = s.withTx(ctx, func(q *sqlc.Queries) error {
		if err := q.RevokeRefreshToken(ctx, stored.ID); err != nil {
			return err
		}
		u, err := q.GetUserByID(ctx, stored.UserID)
		if err != nil {
			return err
		}
		user = u
		newRefresh, err = issueRefreshToken(ctx, q, stored.UserID, s.refreshTTL)
		return err
	})
	if err != nil {
		return RefreshResult{}, err
	}

	access, expiresAt, err := s.jwt.Issue(user.ID.String(), user.Email)
	if err != nil {
		return RefreshResult{}, err
	}
	return RefreshResult{
		AccessToken:  access,
		RefreshToken: newRefresh,
		ExpiresAt:    users.ISOTime(expiresAt),
	}, nil
}

// ForgotPassword issues a password-reset token. It never reveals whether the
// email exists.
func (s *Service) ForgotPassword(ctx context.Context, emailAddr string) error {
	user, err := s.q.GetUserByEmail(ctx, emailAddr)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	} else if err != nil {
		return err
	}

	token := newOpaqueToken()
	if err := s.q.CreatePasswordResetToken(ctx, sqlc.CreatePasswordResetTokenParams{
		UserID:    user.ID,
		TokenHash: sha256Hex(token),
		ExpiresAt: timestamptz(time.Now().Add(passwordResetTTL)),
	}); err != nil {
		return err
	}
	s.email.SendPasswordReset(emailAddr, token)
	return nil
}

// ResetPassword consumes a valid reset token and sets a new password.
func (s *Service) ResetPassword(ctx context.Context, req resetPasswordRequest) error {
	stored, err := s.q.GetPasswordResetTokenByHash(ctx, sha256Hex(req.Token))
	if errors.Is(err, pgx.ErrNoRows) {
		return httpx.Unprocessable("Token expired or already used")
	} else if err != nil {
		return err
	}
	if stored.UsedAt.Valid || stored.ExpiresAt.Time.Before(time.Now()) {
		return httpx.Unprocessable("Token expired or already used")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.withTx(ctx, func(q *sqlc.Queries) error {
		if err := q.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
			ID:           stored.UserID,
			PasswordHash: string(hash),
		}); err != nil {
			return err
		}
		return q.MarkPasswordResetTokenUsed(ctx, stored.ID)
	})
}

// ConfirmEmail consumes a valid confirmation token, marks the user verified
// and issues a first token pair.
func (s *Service) ConfirmEmail(ctx context.Context, token string) (SessionResult, error) {
	stored, err := s.q.GetEmailConfirmationTokenByHash(ctx, sha256Hex(token))
	if errors.Is(err, pgx.ErrNoRows) {
		return SessionResult{}, httpx.Unprocessable("Token expired or already used")
	} else if err != nil {
		return SessionResult{}, err
	}
	if stored.UsedAt.Valid || stored.ExpiresAt.Time.Before(time.Now()) {
		return SessionResult{}, httpx.Unprocessable("Token expired or already used")
	}

	var user sqlc.User
	var refresh string
	err = s.withTx(ctx, func(q *sqlc.Queries) error {
		u, err := q.SetUserEmailVerified(ctx, stored.UserID)
		if err != nil {
			return err
		}
		user = u
		if err := q.MarkEmailConfirmationTokenUsed(ctx, stored.ID); err != nil {
			return err
		}
		refresh, err = issueRefreshToken(ctx, q, stored.UserID, s.refreshTTL)
		return err
	})
	if err != nil {
		return SessionResult{}, err
	}

	access, expiresAt, err := s.jwt.Issue(user.ID.String(), user.Email)
	if err != nil {
		return SessionResult{}, err
	}

	resp := users.Format(user)
	return SessionResult{
		User:         &resp,
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresAt:    users.ISOTime(expiresAt),
	}, nil
}
