package mogotest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// RequireEqualTime is a helper for asserting that two times are equal, at millisecond precision.
func RequireEqualTime(t *testing.T, expected, actual time.Time) {
	t.Helper()

	expected = expected.Truncate(time.Millisecond)
	actual = actual.Truncate(time.Millisecond)

	require.Equal(t, expected, actual)
}
