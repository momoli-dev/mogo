// Package mogogeo provides types and utilities for geographical data.
package mogogeo

import (
	"errors"
	"fmt"

	"github.com/twpayne/go-geom"
)

// DefaultSRID is the SRID of underlying [geom.Point].
const DefaultSRID = 4326

// Point is a geographical latitude and longitude backed by [geom.Point].
type Point struct {
	data *geom.Point

	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

// NewPoint creates a new [Point].
func NewPoint(lat, lng float64) *Point {
	p, err := geom.NewPoint(geom.XY).SetCoords([]float64{lng, lat})
	if err != nil {
		panic(fmt.Errorf("mogogeo failed to create new point: %w", err))
	}
	p.SetSRID(DefaultSRID)

	return &Point{
		data: p,
		Lat:  lat,
		Lng:  lng,
	}
}

// NewPointFromData creates a [Point] from a [geom.Point].
func NewPointFromData(p *geom.Point) (*Point, error) {
	if p.SRID() != DefaultSRID {
		return nil, errors.New("point is not in SRID DefaultSRID")
	}

	coords := p.Coords()
	if len(coords) < 2 {
		return nil, errors.New("point does not have enough coordinates")
	}

	// Lat and Lng are swapped here because we use SRID DefaultSRID which is lng followed by lat.
	return NewPoint(
		coords[1],
		coords[0],
	), nil
}

// MustNewPointFromData same as [NewPointFromData] but panics on any error.
func MustNewPointFromData(p *geom.Point) *Point {
	latlng, err := NewPointFromData(p)
	if err != nil {
		panic(err)
	}

	return latlng
}

// Data returns underlying [geom.Point].
func (l *Point) Data() *geom.Point {
	return l.data
}
