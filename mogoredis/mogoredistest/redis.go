// Package mogoredistest provides clients for local Redis testing.
package mogoredistest

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// New creates a Redis client for tests using local Redis database 15.
func New(t *testing.T) *redis.Client {
	t.Helper()

	const (
		redisDatabaseNumber = 15
	)

	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB:   redisDatabaseNumber,
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	require.NoError(t, client.Ping(ctx).Err())

	t.Cleanup(func() {
		closeErr := client.Close()
		require.NoError(t, closeErr)
	})

	return client
}
