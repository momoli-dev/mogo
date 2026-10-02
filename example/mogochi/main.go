// Package main is an example of using mogochi.
package main

import (
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/momoli-dev/mogo/mogochi"
)

func main() {
	// Start with a default config.
	cfg := mogochi.DefaultConfig()

	// Modify the config as you wish.
	cfg.AllowedOrigins = []string{"https://example.com"}
	cfg.AllowCredentials = true

	// Create a new mux with the config.
	mux := mogochi.NewChi(cfg)

	// Use the mux.
	mux.Get("/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, err := io.WriteString(w, "pong")
		if err != nil {
			panic(err)
		}
	})

	// Serve the mux.
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
		//nolint:mnd // 5.
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("server listening",
		"addr", srv.Addr,
		"url", "http://localhost:8080/ping",
	)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
