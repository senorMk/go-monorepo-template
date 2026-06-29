// Package server wires the HTTP router: middleware stack, routes and docs.
package server

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/unrolled/secure"

	"github.com/GITHUB_USERNAME/APP_NAME/internal/auth"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/config"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/db"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/docs"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/email"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/httpx"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/profile"
	"github.com/GITHUB_USERNAME/APP_NAME/internal/users"
)

// New builds the application's HTTP handler.
func New(cfg config.Config, database *db.DB) http.Handler {
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessExpiresIn)
	emailSvc := email.New(cfg.AppScheme, cfg.EmailFrom)
	authMW := auth.Middleware(jwtManager)

	authHandler := auth.NewHandler(auth.NewService(database.Q, database.Pool, jwtManager, emailSvc, cfg.JWTRefreshExpiresIn))
	profileHandler := profile.NewHandler(profile.NewService(database.Q))

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(requestLogger)
	r.Use(recoverer)
	r.Use(maxBodyBytes(1 << 20)) // 1 MiB request body cap
	r.Use(secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		ReferrerPolicy:     "no-referrer",
	}).Handler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(middleware.Compress(5))

	r.Get("/health", health)

	authHandler.RegisterRoutes(r, authMW)
	profileHandler.RegisterRoutes(r, authMW)

	// Swagger UI is only mounted outside production.
	if !cfg.IsProduction() {
		docs.Mount(r)
	}

	return r
}

func health(w http.ResponseWriter, _ *http.Request) {
	httpx.WriteJSON(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"timestamp": users.ISOTime(time.Now()),
	})
}

// maxBodyBytes caps the size of request bodies to guard against
// memory-exhaustion. Decoding an oversized body yields a 400 via httpx.
func maxBodyBytes(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

// recoverer turns a panic into the standard JSON error envelope (500) instead
// of chi's default plain-text response, keeping the error contract consistent.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler {
					panic(rec)
				}
				slog.Error("panic recovered", "err", rec, "path", r.URL.Path)
				httpx.WriteError(w, errors.New("internal server error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requestLogger logs each request as a structured slog line.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		next.ServeHTTP(ww, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.Status(),
			"bytes", ww.BytesWritten(),
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", middleware.GetReqID(r.Context()),
		)
	})
}
