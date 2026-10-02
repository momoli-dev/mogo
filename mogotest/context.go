// Package mogotest provides helpers for testing.
package mogotest

import (
	"context"
	"testing"
	"time"
)

// CtxWithTimeout creates a context canceled at the timeout or test cleanup.
func CtxWithTimeout(t *testing.T, timeout time.Duration) context.Context {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)

	return ctx
}
