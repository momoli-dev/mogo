// Package main is an example of using mogohttp.
package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/momoli-dev/mogo/mogohttp"
)

func main() {
	// Create and configure a mux.
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		_, err := io.WriteString(w, "pong")
		if err != nil {
			panic(err)
		}
	})

	// Start with default config.
	cfg := mogohttp.DefaultConfig(":8080", mux)

	// Modify the config as you wish.
	cfg.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Context is respected by the server.
	ctx := context.Background()

	// Start the HTTP server.
	if err := mogohttp.Serve(ctx, cfg); err != nil {
		panic(err)
	}
}
