package main_test

import (
	"log/slog"
	"testing"

	"github.com/momoli-dev/mogo/mogoredis/mogoredistest"
)

func TestUsingRedis(t *testing.T) {
	t.Parallel()

	rdb := mogoredistest.New(t)

	if err := rdb.Set(t.Context(), "foo", "bar", 0).Err(); err != nil {
		t.Fatal(err)
	}

	if err := rdb.Get(t.Context(), "foo").Err(); err != nil {
		t.Fatal(err)
	}

	slog.Info("Redis test passed")
}
