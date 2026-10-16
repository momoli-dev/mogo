// Package mogoredistest provides clients for local Redis testing.
package mogoredistest

import (
	"context"
	"testing"
	"time"

	"github.com/momoli-dev/mogo/mogoredis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// New creates a Redis client for tests using local Redis database 15.
func New(t *testing.T) *redis.Client {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := mogoredis.NewClient(ctx, "redis://localhost:6379/15")
	require.NoError(t, err)

	t.Cleanup(func() {
		closeErr := client.Close()
		require.NoError(t, closeErr)
	})

	return client
}
