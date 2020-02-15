package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocation_IsEmpty(t *testing.T) {
	testCases := []struct {
		loc           Location
		shouldBeEmpty bool
	}{
		{
			loc: Location{
				Lat: 0,
				Lng: 0,
			},
			shouldBeEmpty: true,
		},
		{
			loc:           Location{},
			shouldBeEmpty: true,
		},
		{
			loc: Location{
				Lat: 123,
				Lng: 321,
			},
			shouldBeEmpty: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.shouldBeEmpty, tc.loc.IsEmpty())
	}
}
