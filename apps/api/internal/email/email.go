// Package email is a stub email service. Swap the Send* methods for a real
// provider (Resend, SES, Postmark, ...) when wiring production email.
package email

import (
	"fmt"
	"log/slog"
)

// Service sends transactional emails. This stub just logs the links.
type Service struct {
	scheme string // deep-link scheme, e.g. "my-app"
	from   string
}

// New creates an email service stub.
func New(scheme, from string) *Service {
	return &Service{scheme: scheme, from: from}
}

// SendEmailConfirmation sends (logs) an email-confirmation link.
func (s *Service) SendEmailConfirmation(email, token string) {
	link := fmt.Sprintf("%s://auth/confirm?token=%s", s.scheme, token)
	slog.Info("[EMAIL] confirm email", "to", email, "from", s.from, "link", link)
}

// SendPasswordReset sends (logs) a password-reset link.
func (s *Service) SendPasswordReset(email, token string) {
	link := fmt.Sprintf("%s://auth/reset-password?token=%s", s.scheme, token)
	slog.Info("[EMAIL] password reset", "to", email, "from", s.from, "link", link)
}
