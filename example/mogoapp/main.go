// Package main is an example of using mogoapp.
package main

import (
	"log/slog"

	"github.com/momoli-dev/mogo/mogoapp"
)

func main() {
	slog.Info("Starting application")

	// Create a root context with signal handling.
	ctx, cancel := mogoapp.NewAppCtx()
	defer cancel()
	<-ctx.Done()

	slog.Info("Application stopped")
}
