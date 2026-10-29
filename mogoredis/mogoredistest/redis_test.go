package mogoredistest_test

import (
	"testing"
	"uuid"

	"github.com/momoli-dev/mogo/mogoredis/mogoredistest"
	"github.com/stretchr/testify/require"
)

//nolint:tparallel // The parent verifies deletion after the sequential subtest cleans up.
func TestCleanupKeys(t *testing.T) {
	t.Parallel()
	client := mogoredistest.New(t)
	key := "cleanup:" + uuid.NewV7().String()
	other := "cleanup:" + uuid.NewV7().String()
	mogoredistest.CleanupKeys(t, client, key, other)
	require.NoError(t, client.Set(t.Context(), other, "keep", 0).Err())

	t.Run("register", func(t *testing.T) {
		// Keep this subtest sequential so its cleanup finishes before the assertions below.
		mogoredistest.CleanupKeys(t, client, key)
		require.NoError(t, client.Set(t.Context(), key, "remove", 0).Err())
	})

	count, err := client.Exists(t.Context(), key).Result()
	require.NoError(t, err)
	require.Zero(t, count)
	value, err := client.Get(t.Context(), other).Result()
	require.NoError(t, err)
	require.Equal(t, "keep", value)
}
