package mogotest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// RequireEqualFloat64 is a helper for asserting that two floats are equal.
func RequireEqualFloat64(t *testing.T, expected, actual float64) {
	t.Helper()
	const epsilon = 0.001
	require.InEpsilon(t, expected, actual, epsilon)
}
