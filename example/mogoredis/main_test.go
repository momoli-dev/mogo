package main_test

import (
	"log/slog"
	"testing"
	"uuid"

	"github.com/momoli-dev/mogo/mogoredis/mogoredistest"
)

func TestUsingRedis(t *testing.T) {
	t.Parallel()

	rdb := mogoredistest.New(t)
	key := "example:" + uuid.NewV7().String()
	mogoredistest.CleanupKeys(t, rdb, key)

	if err := rdb.Set(t.Context(), key, "bar", 0).Err(); err != nil {
		t.Fatal(err)
	}

	if err := rdb.Get(t.Context(), key).Err(); err != nil {
		t.Fatal(err)
	}

	slog.Info("Redis test passed")
}
