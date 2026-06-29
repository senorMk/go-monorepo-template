// Package config loads and validates runtime configuration from the
// environment. It mirrors the variables in .env.example.
package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// minJWTSecretLen is the smallest acceptable JWT signing secret (256 bits).
const minJWTSecretLen = 32

// Config holds all runtime configuration, populated from environment variables.
type Config struct {
	Env  string `env:"APP_ENV" envDefault:"development"`
	Port int    `env:"PORT" envDefault:"3000"`

	DatabaseURL string `env:"DATABASE_URL,required"`

	JWTSecret string `env:"JWT_SECRET,required"`
	// JWTRefreshSecret is reserved: refresh tokens are opaque (random, hashed at
	// rest), so no signing secret is used today. Kept for forward-compatibility.
	JWTRefreshSecret    string        `env:"JWT_REFRESH_SECRET" envDefault:""`
	JWTAccessExpiresIn  time.Duration `env:"JWT_ACCESS_EXPIRES_IN" envDefault:"15m"`
	JWTRefreshExpiresIn time.Duration `env:"JWT_REFRESH_EXPIRES_IN" envDefault:"720h"`

	CORSOrigins []string `env:"CORS_ORIGINS" envSeparator:"," envDefault:"http://localhost:4000"`

	EmailFrom string `env:"EMAIL_FROM" envDefault:"info@APP_DOMAIN"`

	// AppScheme is the deep-link scheme used in confirmation / reset emails.
	AppScheme string `env:"APP_SCHEME" envDefault:"APP_NAME"`
}

// IsProduction reports whether the API is running in production mode.
func (c Config) IsProduction() bool { return c.Env == "production" }

// Load reads configuration from the environment, applying defaults and
// validating that required values are present.
func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	if len(cfg.JWTSecret) < minJWTSecretLen {
		return Config{}, errors.New("JWT_SECRET must be at least 32 bytes (use a random 256-bit string)")
	}
	return cfg, nil
}
