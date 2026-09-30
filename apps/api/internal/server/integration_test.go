package server_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/GITHUB_USERNAME/APP_NAME/internal/config"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/db"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/server"
)

// Run with TEST_DATABASE_URL pointing to a disposable PostgreSQL database.
// CI supplies one through its Postgres service.
func TestAuthWithPostgres(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run the Postgres integration test")
	}
	ctx := context.Background()
	if err := db.Migrate(dsn); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if err := db.Migrate(dsn); err != nil {
		t.Fatalf("reapply migrations: %v", err)
	}
	database, err := db.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connect to Postgres: %v", err)
	}
	t.Cleanup(database.Close)

	cfg := config.Config{
		Env: "test", JWTSecret: "integration-test-secret-at-least-32-bytes",
		JWTAccessExpiresIn: time.Minute, JWTRefreshExpiresIn: time.Hour,
		AppScheme: "test-app", EmailFrom: "test@example.com",
	}
	handler := server.New(cfg, database)
	email := "integration-" + uuid.NewString() + "@example.com"
	password := "correct-password"

	signup := requestJSON(t, handler, "/v1/auth/signup", map[string]string{
		"email": email, "password": password,
	})
	if signup.Code != http.StatusCreated {
		t.Fatalf("signup status = %d, body = %s", signup.Code, signup.Body.String())
	}
	var created struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		RequiresEmailConfirmation bool `json:"requiresEmailConfirmation"`
	}
	decodeJSON(t, signup, &created)
	if !created.RequiresEmailConfirmation || created.User.ID == "" {
		t.Fatalf("signup did not create an unconfirmed user: %+v", created)
	}
	userID, err := uuid.Parse(created.User.ID)
	if err != nil {
		t.Fatalf("parse user ID: %v", err)
	}
	t.Cleanup(func() {
		if _, err := database.Pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userID); err != nil {
			t.Errorf("remove test user: %v", err)
		}
	})

	signin := requestJSON(t, handler, "/v1/auth/signin", map[string]string{
		"email": email, "password": password,
	})
	if signin.Code != http.StatusForbidden {
		t.Fatalf("unconfirmed signin status = %d, body = %s", signin.Code, signin.Body.String())
	}

	// The email stub logs the real link. Seed a known token to exercise the
	// confirmation endpoint without depending on log output.
	confirmationToken := uuid.NewString()
	hash := sha256.Sum256([]byte(confirmationToken))
	_, err = database.Pool.Exec(ctx,
		"INSERT INTO email_confirmation_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)",
		userID, hex.EncodeToString(hash[:]), time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("seed confirmation token: %v", err)
	}
	confirmed := requestJSON(t, handler, "/v1/auth/confirm-email", map[string]string{"token": confirmationToken})
	if confirmed.Code != http.StatusOK {
		t.Fatalf("confirm status = %d, body = %s", confirmed.Code, confirmed.Body.String())
	}

	signin = requestJSON(t, handler, "/v1/auth/signin", map[string]string{
		"email": email, "password": password,
	})
	if signin.Code != http.StatusOK {
		t.Fatalf("confirmed signin status = %d, body = %s", signin.Code, signin.Body.String())
	}
	var session struct {
		RefreshToken string `json:"refreshToken"`
	}
	decodeJSON(t, signin, &session)
	if session.RefreshToken == "" {
		t.Fatal("signin returned no refresh token")
	}
	refreshed := requestJSON(t, handler, "/v1/auth/refresh", map[string]string{
		"refreshToken": session.RefreshToken,
	})
	if refreshed.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, body = %s", refreshed.Code, refreshed.Body.String())
	}
	var next struct {
		RefreshToken string `json:"refreshToken"`
	}
	decodeJSON(t, refreshed, &next)
	if next.RefreshToken == "" || next.RefreshToken == session.RefreshToken {
		t.Fatal("refresh did not rotate the token")
	}
	if replay := requestJSON(t, handler, "/v1/auth/refresh", map[string]string{"refreshToken": session.RefreshToken}); replay.Code != http.StatusUnauthorized {
		t.Fatalf("replayed token status = %d, body = %s", replay.Code, replay.Body.String())
	}
	if revoked := requestJSON(t, handler, "/v1/auth/refresh", map[string]string{"refreshToken": next.RefreshToken}); revoked.Code != http.StatusUnauthorized {
		t.Fatalf("token family status = %d, body = %s", revoked.Code, revoked.Body.String())
	}
}

func requestJSON(t *testing.T, handler http.Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode response %s: %v", response.Body.String(), err)
	}
}
