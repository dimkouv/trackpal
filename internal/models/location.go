package models

import "math"

// Location is a latitude, longitude pair of earth coordinates.
type Location struct {
	Lat float64 `json:"lat" db:"lat"`
	Lng float64 `json:"lng" db:"lng"`
}

func (loc Location) IsEmpty() bool {
	return loc == Location{}
}

// HaversineDistanceKM calculates the haversine distance between two locations in kilometers
func (loc Location) haversineDistanceKM(loc2 Location) float64 {
	degreesToRadians := func(d float64) float64 {
		return d * math.Pi / 180
	}
	const earthRadiusKm = 6371

	lat1, lng1 := degreesToRadians(loc.Lat), degreesToRadians(loc.Lng)
	lat2, lng2 := degreesToRadians(loc2.Lat), degreesToRadians(loc2.Lng)

	diffLat := lat2 - lat1
	diffLng := lng2 - lng1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(diffLng/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * earthRadiusKm
}

// HasMovedMoreThanM checks whether or not loc has moved more than 'm' meters
// from loc2
func (loc Location) HasMovedMoreThanM(loc2 Location, m int) bool {
	distMeters := loc.haversineDistanceKM(loc2) * 1000.0
	return distMeters > float64(m)
}
