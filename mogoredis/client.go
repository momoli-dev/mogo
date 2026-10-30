// Package mogoredis provides helpers for working with Redis.
package mogoredis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// NewClient connects to Redis using a connection URL and verifies the connection.
// For example, connStr may be redis://localhost:6379/0.
// The caller must close the returned client.
func NewClient(ctx context.Context, connStr string) (*redis.Client, error) {
	opts, err := redis.ParseURL(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to redis.ParseURL: %w", err)
	}

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		pingErr := fmt.Errorf("failed to client.Ping: %w", err)
		if closeErr := client.Close(); closeErr != nil {
			return nil, errors.Join(pingErr, fmt.Errorf("failed to client.Close: %w", closeErr))
		}
		return nil, pingErr
	}

	return client, nil
}
