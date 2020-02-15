// +build unit_test

package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrackInput_IsEmpty(t *testing.T) {
	testCases := []struct {
		ti            TrackInput
		shouldBeEmpty bool
	}{
		{
			TrackInput{
				Location:  Location{},
				Timestamp: time.Time{},
				DeviceID:  0,
			},
			true,
		},
		{
			TrackInput{
				Location: Location{
					Lat: 123,
					Lng: 0,
				},
				Timestamp: time.Time{},
				DeviceID:  0,
			},
			false,
		},
		{
			TrackInput{
				Location:  Location{},
				Timestamp: time.Time{},
				DeviceID:  1,
			},
			false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.shouldBeEmpty, tc.ti.IsEmpty())
	}
}

func TestTrackInput_Validate(t *testing.T) {
	testCases := []struct {
		ti            TrackInput
		shouldBeValid bool
	}{
		{
			ti: TrackInput{
				Location:  Location{},
				Timestamp: time.Time{},
				DeviceID:  0,
			},
			shouldBeValid: false,
		},
		{
			ti: TrackInput{
				Location: Location{
					Lat: 0,
					Lng: 0,
				},
				Timestamp: time.Time{},
				DeviceID:  123,
			},
			shouldBeValid: false,
		},
		{
			ti: TrackInput{
				Location: Location{
					Lat: 123,
					Lng: 0,
				},
				Timestamp: time.Time{},
				DeviceID:  123,
			},
			shouldBeValid: true,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.shouldBeValid, tc.ti.Validate() == nil) {
			t.Logf("%v", tc)
		}
	}
}
