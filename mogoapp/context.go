// Package mogoapp provides heavily opinionated helpers for patterns encountered while building Go applications.
package mogoapp

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// NewAppCtx returns a completely new context that is intended to be used as the root context for an application. It
// registers signal handling.
func NewAppCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	return ctx, cancel
}
