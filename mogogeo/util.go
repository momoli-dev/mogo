package mogogeo

import "math"

// PointsAreWithinDistance returns true if the distance between a and b is less than dMeters.
//
//nolint:mnd // This is maths.
func PointsAreWithinDistance(a, b *Point, dMeters float64) bool {
	// Approximate degrees to meters
	const metersPerDegLat = 111320.0
	metersPerDegLng := 111320.0 * math.Cos(a.Lat*math.Pi/180)

	dLat := (a.Lat - b.Lat) * metersPerDegLat
	dLng := (a.Lng - b.Lng) * metersPerDegLng

	return dLat*dLat+dLng*dLng <= dMeters*dMeters
}
