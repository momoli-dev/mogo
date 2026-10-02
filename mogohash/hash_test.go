package mogohash_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/momoli-dev/mogo/mogohash"
)

func TestPassword(t *testing.T) {
	t.Parallel()
	text := "pa$$word"
	p := mogohash.NewHashedString(text)
	require.NotNil(t, p)

	require.True(t, p.Check(text))
	require.False(t, p.Check("wrong"))
}
