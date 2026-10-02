// Package mogohttp provides opinionated helpers for serving HTTP servers.
package mogohttp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Config configures [Serve]. Start with [DefaultConfig] and
// override the values needed by the application.
type Config struct {
	Addr              string
	Handler           http.Handler
	Logger            *slog.Logger
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
}

// DefaultConfig returns the production-oriented server defaults.
//
//nolint:mnd // These are inherently not magic numbers.
func DefaultConfig(addr string, handler http.Handler) Config {
	return Config{
		Addr:              addr,
		Handler:           handler,
		Logger:            slog.Default(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		ShutdownTimeout:   10 * time.Second,
		MaxHeaderBytes:    1 * 1024 * 1024,
	}
}

// Serve serves an HTTP server with passed in [config] and graceful
// shut down support.
func Serve(ctx context.Context, config Config) error {
	s := &http.Server{
		Addr:              config.Addr,
		Handler:           config.Handler,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		MaxHeaderBytes:    config.MaxHeaderBytes,
	}

	errChan := make(chan error, 1)

	go func() {
		config.Logger.Info("Starting HTTP server", "addr", config.Addr)

		err := s.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		}

		errChan <- err
	}()

	select {
	case err := <-errChan:
		if err != nil {
			config.Logger.Error("Error encountered while serving HTTP server", "addr", config.Addr, "err", err)
			return fmt.Errorf("serve HTTP: %w", err)
		}

		return nil
	case <-ctx.Done():
		config.Logger.Info("Shutting down HTTP server", "addr", config.Addr, "timeout", config.ShutdownTimeout)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer cancel()

		//nolint:contextcheck // shutdownCtx uses new background context since ctx is already done.
		shutdownErr := s.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			shutdownErr = fmt.Errorf("shut down HTTP server: %w", shutdownErr)

			closeErr := s.Close()
			if closeErr != nil {
				closeErr = fmt.Errorf("force close HTTP server: %w", closeErr)
			}

			return errors.Join(shutdownErr, closeErr)
		}

		return nil
	}
}
