package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// sha256Hex returns the hex-encoded SHA-256 digest of s. Opaque tokens are
// stored hashed so a database leak does not expose usable tokens.
func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// newOpaqueToken returns a random, unguessable token (UUIDv4).
func newOpaqueToken() string {
	return uuid.NewString()
}

// timestamptz wraps a time.Time as a non-null pgtype.Timestamptz.
func timestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}
