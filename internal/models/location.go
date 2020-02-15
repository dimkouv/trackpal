package models

// Location is a latitude, longitude pair of earth coordinates.
type Location struct {
	Lat float64 `json:"lat" db:"lat"`
	Lng float64 `json:"lng" db:"lng"`
}

func (loc Location) IsEmpty() bool {
	return loc == Location{}
}
