package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDevice_IsEmpty(t *testing.T) {
	testCases := []struct {
		device        Device
		shouldBeEmpty bool
	}{
		{
			device: Device{
				ID:   0,
				Name: "",
			},
			shouldBeEmpty: true,
		},
		{
			device: Device{
				ID:   1,
				Name: "",
			},
			shouldBeEmpty: false,
		},
		{
			device: Device{
				ID:   0,
				Name: " ",
			},
			shouldBeEmpty: false,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.shouldBeEmpty, tc.device.IsEmpty())
	}
}

func TestDevice_Validate(t *testing.T) {
	testCases := []struct {
		device        Device
		shouldBeValid bool
	}{
		{
			device: Device{
				ID:   0,
				Name: "",
			},
			shouldBeValid: false,
		},
		{
			device: Device{
				ID:   123,
				Name: "",
			},
			shouldBeValid: false,
		},
		{
			device: Device{
				ID:   123,
				Name: "123",
			},
			shouldBeValid: false,
		},
		{
			device: Device{
				ID:   123,
				Name: "dev123",
			},
			shouldBeValid: true,
		},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.shouldBeValid, tc.device.Validate() == nil)
	}
}
