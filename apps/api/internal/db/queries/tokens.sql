-- Refresh tokens -----------------------------------------------------------

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = now() WHERE id = $1;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked_at = now()
WHERE user_id = $1 AND revoked_at IS NULL;

-- Password reset tokens -----------------------------------------------------

-- name: CreatePasswordResetToken :exec
INSERT INTO password_reset_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetPasswordResetTokenByHash :one
SELECT * FROM password_reset_tokens WHERE token_hash = $1;

-- name: MarkPasswordResetTokenUsed :exec
UPDATE password_reset_tokens SET used_at = now() WHERE id = $1;

-- Email confirmation tokens -------------------------------------------------

-- name: CreateEmailConfirmationToken :exec
INSERT INTO email_confirmation_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3);

-- name: GetEmailConfirmationTokenByHash :one
SELECT * FROM email_confirmation_tokens WHERE token_hash = $1;

-- name: MarkEmailConfirmationTokenUsed :exec
UPDATE email_confirmation_tokens SET used_at = now() WHERE id = $1;
