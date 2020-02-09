// +build unit_test

package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dimkouv/trackpal/internal/models"
)

func TestTrackingRepositoryMock_SaveNewTrackInput(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	trackingInputs := []models.TrackInput{
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-1 * time.Hour),
			DeviceID:  1,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-50 * time.Minute),
			DeviceID:  1,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-2 * time.Hour),
			DeviceID:  1,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-1 * time.Hour),
			DeviceID:  2,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-2 * time.Hour),
			DeviceID:  2,
		},
	}

	for i := range trackingInputs {
		repo.SaveNewTrackInput(trackingInputs[i])
	}
}

func TestTrackingRepositoryMock_GetAllTrackInputsOfDevice(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	trackingInputs := []models.TrackInput{
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-1 * time.Hour),
			DeviceID:  1,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-50 * time.Minute),
			DeviceID:  1,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-2 * time.Hour),
			DeviceID:  1,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-1 * time.Hour),
			DeviceID:  2,
		},
		{
			Location:  models.Location{Lat: 12.123, Lng: 12.321},
			Timestamp: time.Now().Add(-2 * time.Hour),
			DeviceID:  2,
		},
	}

	for i := range trackingInputs {
		repo.SaveNewTrackInput(trackingInputs[i])
	}

	var results []models.TrackInput

	t.Run("the records should be fetched correctly", func(t *testing.T) {
		results = repo.GetAllTrackInputsOfDevice(1)
		assert.Len(t, results, 3)
		results = append(results, repo.GetAllTrackInputsOfDevice(2)...)
		assert.Len(t, results, 5)
		results = append(results, repo.GetAllTrackInputsOfDevice(3)...)
		assert.Len(t, results, 5)
	})

	t.Run("validate that the results are coming in timestamp ascending order", func(t *testing.T) {
		assert.Equal(t, results[0], trackingInputs[2])
		assert.Equal(t, results[1], trackingInputs[0])
		assert.Equal(t, results[2], trackingInputs[1])
		assert.Equal(t, results[3], trackingInputs[4])
		assert.Equal(t, results[4], trackingInputs[3])
	})
}
