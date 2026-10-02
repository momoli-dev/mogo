package mogogeo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/momoli-dev/mogo/mogogeo"
)

func TestLatLngsAreWithinDistance(t *testing.T) {
	t.Parallel()
	a := mogogeo.NewPoint(1, 2)
	c := mogogeo.NewPoint(5, 6)

	require.False(t, mogogeo.PointsAreWithinDistance(a, c, 100))
}
