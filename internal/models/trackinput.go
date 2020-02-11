package models

import (
	"errors"
	"time"
)

// TrackInput contains a location, the time that it was recorded and the id of the device that recorded it.
type TrackInput struct {
	Location  Location  `json:"location"`
	Timestamp time.Time `json:"timestamp"`

	DeviceID int64 `json:"-"`
}

func (t TrackInput) IsEmpty() bool {
	return t == TrackInput{}
}

func (t TrackInput) Validate() error {
	if t.IsEmpty() {
		return errors.New("the track input is empty")
	}

	if t.Location.IsEmpty() {
		return errors.New("the location of the track input is empty")
	}

	return nil
}
