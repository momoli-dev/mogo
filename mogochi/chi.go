// Package mogochi provides opinionated helpers for using Chi as a production ready HTTP Mux.
package mogochi

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/httprate"
)

const (
	requestTimeout           = 30 * time.Second
	corsMaxAge               = 600
	maxRequestCountPerMinute = 250
)

// Config configures [NewChi].
type Config struct {
	AllowedOrigins      []string
	Logger              *slog.Logger
	RequestTimeout      time.Duration
	RateLimit           int
	RateWindow          time.Duration
	RateLimitMiddleware func(http.Handler) http.Handler
	CORSMaxAge          int
	AllowCredentials    bool
}

// DefaultConfig returns opinionated, production oriented [Config].
func DefaultConfig() Config {
	return Config{
		AllowedOrigins:      []string{"https://*"},
		Logger:              slog.Default(),
		RequestTimeout:      requestTimeout,
		RateLimit:           maxRequestCountPerMinute,
		RateWindow:          time.Minute,
		RateLimitMiddleware: nil,
		CORSMaxAge:          corsMaxAge,
		AllowCredentials:    false,
	}
}

// NewChi creates a new mux with chi library using passed in config.
func NewChi(p Config) *chi.Mux {
	r := chi.NewMux()

	r.Use(
		chimiddleware.ClientIPFromXFF(),
		httplog.RequestLogger(
			p.Logger,
			&httplog.Options{
				Level:         slog.LevelWarn,
				Schema:        httplog.SchemaECS,
				RecoverPanics: true,
				Skip: func(_ *http.Request, respStatus int) bool {
					return respStatus == http.StatusNotFound
				},
			},
		),
	)
	if p.RateLimitMiddleware != nil {
		r.Use(p.RateLimitMiddleware)
	} else if p.RateLimit > 0 && p.RateWindow > 0 {
		r.Use(
			httprate.Limit(
				p.RateLimit,
				p.RateWindow,
				httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
					ip := chimiddleware.GetClientIP(r.Context())
					if ip == "" {
						return "", errors.New("missing client IP")
					}
					return ip, nil
				},
				),
			),
		)
	}
	if p.RequestTimeout > 0 {
		r.Use(chimiddleware.Timeout(p.RequestTimeout))
	}
	r.Use(
		cors.Handler(
			//nolint:exhaustruct
			cors.Options{
				AllowedOrigins: p.AllowedOrigins,
				AllowedMethods: []string{
					http.MethodGet,
					http.MethodPost,
					http.MethodPut,
					http.MethodDelete,
					http.MethodPatch,
					http.MethodOptions,
					http.MethodHead,
				},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: p.AllowCredentials,
				MaxAge:           p.CORSMaxAge,
				Debug:            false,
			},
		),
	)

	return r
}
