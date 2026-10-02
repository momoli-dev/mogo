package main_test

import (
	"testing"
	"time"

	"github.com/momoli-dev/mogo/mogotest"
	"github.com/stretchr/testify/require"
)

type greeting struct {
	Message string `json:"message"`
}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()

	// Decode JSON into a typed value. Invalid JSON fails the test automatically.
	got := mogotest.JSONUnmarshal[greeting](t, []byte(`{"message":"Hello, World!"}`))
	require.Equal(t, "Hello, World!", got.Message)
}

func TestContextWithTimeout(t *testing.T) {
	t.Parallel()

	// Simple context with timeout that is canceled at the end of the test.
	ctx := mogotest.CtxWithTimeout(t, 5*time.Second)
	require.NotNil(t, ctx)
}
