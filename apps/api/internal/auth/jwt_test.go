package auth

import (
	"testing"
	"time"
)

func TestJWTRoundTrip(t *testing.T) {
	m := NewJWTManager("a-test-secret", time.Minute)

	token, expiresAt, err := m.Issue("user-1", "alice@example.com")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	if !expiresAt.After(time.Now()) {
		t.Fatal("expiry should be in the future")
	}

	claims, err := m.Parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Fatalf("subject = %q, want user-1", claims.Subject)
	}
	if claims.Email != "alice@example.com" {
		t.Fatalf("email = %q, want alice@example.com", claims.Email)
	}
}

func TestJWTRejectsWrongSecret(t *testing.T) {
	token, _, _ := NewJWTManager("right", time.Minute).Issue("u", "e")
	if _, err := NewJWTManager("wrong", time.Minute).Parse(token); err == nil {
		t.Fatal("expected error parsing token signed with a different secret")
	}
}

func TestJWTRejectsExpired(t *testing.T) {
	token, _, _ := NewJWTManager("secret", -time.Minute).Issue("u", "e")
	if _, err := NewJWTManager("secret", time.Minute).Parse(token); err == nil {
		t.Fatal("expected error parsing an expired token")
	}
}
