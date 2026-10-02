package mogogeo_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/momoli-dev/mogo/mogogeo"
	"github.com/momoli-dev/mogo/mogotest"
)

func TestNewPoint(t *testing.T) {
	t.Parallel()
	p := mogogeo.NewPoint(1.0, 2.0)
	mogotest.RequireEqualFloat64(t, 1.0, p.Lat)
	mogotest.RequireEqualFloat64(t, 2.0, p.Lng)
}

func TestNewPointFromPoint(t *testing.T) {
	t.Parallel()
	data := mogogeo.NewPoint(1.0, 2.0).Data()
	p, err := mogogeo.NewPointFromData(data)
	require.NoError(t, err)
	mogotest.RequireEqualFloat64(t, 1.0, p.Lat)
	mogotest.RequireEqualFloat64(t, 2.0, p.Lng)
}

func TestNewPointFromPointError_WrongSRID(t *testing.T) {
	t.Parallel()
	p := mogogeo.NewPoint(1, 2).Data()
	p.SetSRID(1111)
	_, err := mogogeo.NewPointFromData(p)
	require.Error(t, err)
}

func TestMustNewPointFromPoint(t *testing.T) {
	t.Parallel()
	data := mogogeo.NewPoint(1.0, 2.0).Data()
	p := mogogeo.MustNewPointFromData(data)
	mogotest.RequireEqualFloat64(t, 1.0, p.Lat)
	mogotest.RequireEqualFloat64(t, 2.0, p.Lng)
}

func TestMustNewPointFromPointError_WrongSRID(t *testing.T) {
	t.Parallel()
	data := mogogeo.NewPoint(1, 2).Data()
	data.SetSRID(1111)
	require.Panics(t, func() {
		mogogeo.MustNewPointFromData(data)
	})
}

func TestMustNewPointFromPointError_NoCoords(t *testing.T) {
	t.Parallel()
	data := mogogeo.NewPoint(1, 2).Data()
	_, err := data.SetCoords([]float64{})
	require.Error(t, err)
	require.Panics(t, func() {
		mogogeo.MustNewPointFromData(data)
	})
}

func TestMustNewPointFromPointError_NotEnoughCoords(t *testing.T) {
	t.Parallel()
	data := mogogeo.NewPoint(1, 2).Data()
	_, err := data.SetCoords([]float64{1})
	require.Error(t, err)
	require.Panics(t, func() {
		mogogeo.MustNewPointFromData(data)
	})
}

func TestPoint(t *testing.T) {
	t.Parallel()
	p := mogogeo.NewPoint(1.0, 2.0)
	coords := p.Data().Coords()
	mogotest.RequireEqualFloat64(t, coords[0], 2.0)
	mogotest.RequireEqualFloat64(t, coords[1], 1.0)
}
