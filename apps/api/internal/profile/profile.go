// Package profile exposes the authenticated user's own profile.
package profile

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/GITHUB_USERNAME/APP_NAME/internal/auth"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/db/sqlc"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/httpx"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/users"
)

// optionalString distinguishes an omitted JSON key (Present == false) from an
// explicit value or null, so PATCH can skip fields that were not sent — the
// same "undefined skips, null clears" semantics as the Prisma original.
type optionalString struct {
	Value   *string
	Present bool
}

func (o *optionalString) UnmarshalJSON(b []byte) error {
	o.Present = true
	if string(b) == "null" {
		o.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	o.Value = &s
	return nil
}

type updateProfileRequest struct {
	DisplayName optionalString `json:"displayName" validate:"-"`
}

// Service holds the profile dependencies.
type Service struct {
	q *sqlc.Queries
}

// NewService builds a profile Service.
func NewService(q *sqlc.Queries) *Service { return &Service{q: q} }

// GetMe returns the user's own profile.
func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (users.Response, error) {
	u, err := s.q.GetUserByID(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return users.Response{}, httpx.NotFound("User not found")
	} else if err != nil {
		return users.Response{}, err
	}
	return users.Format(u), nil
}

// UpdateMe updates the user's display name.
func (s *Service) UpdateMe(ctx context.Context, userID uuid.UUID, displayName *string) (users.Response, error) {
	u, err := s.q.UpdateUserDisplayName(ctx, sqlc.UpdateUserDisplayNameParams{
		ID:          userID,
		DisplayName: displayName,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return users.Response{}, httpx.NotFound("User not found")
	} else if err != nil {
		return users.Response{}, err
	}
	return users.Format(u), nil
}

// Handler exposes the profile HTTP endpoints.
type Handler struct {
	svc *Service
}

// NewHandler builds a profile Handler.
func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

// RegisterRoutes mounts the /v1/profile routes behind the auth middleware.
func (h *Handler) RegisterRoutes(r chi.Router, authMW func(http.Handler) http.Handler) {
	r.Route("/v1/profile", func(r chi.Router) {
		r.Use(authMW)
		r.Get("/me", h.getMe)
		r.Patch("/me", h.updateMe)
	})
}

func (h *Handler) currentUserID(r *http.Request) (uuid.UUID, error) {
	user, ok := auth.CurrentUser(r.Context())
	if !ok {
		return uuid.Nil, httpx.Unauthorized("Missing or expired token")
	}
	id, err := uuid.Parse(user.ID)
	if err != nil {
		return uuid.Nil, httpx.Unauthorized("Invalid token subject")
	}
	return id, nil
}

func (h *Handler) getMe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.currentUserID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	res, err := h.svc.GetMe(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) updateMe(w http.ResponseWriter, r *http.Request) {
	userID, err := h.currentUserID(r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	var req updateProfileRequest
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.WriteError(w, err)
		return
	}

	// PATCH: a field that was not sent leaves the stored value unchanged.
	if !req.DisplayName.Present {
		res, err := h.svc.GetMe(r.Context(), userID)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, res)
		return
	}

	res, err := h.svc.UpdateMe(r.Context(), userID, req.DisplayName.Value)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, res)
}
