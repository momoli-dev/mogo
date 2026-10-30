package mogotest

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// RequireEqualFloat64 requires actual to be within a relative error of 0.001 of expected.
func RequireEqualFloat64(t *testing.T, expected, actual float64) {
	t.Helper()
	const epsilon = 0.001
	require.InEpsilon(t, expected, actual, epsilon)
}
