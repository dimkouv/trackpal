package models

import "time"

// TrackInput contains a location, the time that it was recorded and the id of the device that recorded it.
type TrackInput struct {
	Location  Location  `json:"location"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  int64     `json:"deviceID"`
}
