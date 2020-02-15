// +build unit_test

package repository

import (
	"testing"
	"time"

	"github.com/dimkouv/trackpal/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestNewTrackingRepositoryMock(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	assert.Equal(t, 0, len(repo.devices))
	assert.Equal(t, 0, len(repo.trackInputs))
	assert.Equal(t, 0, len(repo.trackInputsDeviceIDX))
}

func TestTrackingRepositoryMock_SaveNewDevice(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	t.Run("trying to save an invalid device should raise an error", func(t *testing.T) {
		_, err := repo.SaveNewDevice(models.Device{
			ID:   0,
			Name: "",
		})
		assert.Error(t, err)
	})

	t.Run("saving a valid device should be successful", func(t *testing.T) {
		d, err := repo.SaveNewDevice(models.Device{
			ID:   123, // the passed id is ignored
			Name: "my-device",
		})
		assert.NoError(t, err)
		assert.Equal(t, int64(1), d.ID)
	})
}

func TestTrackingRepositoryMock_GetDevices(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	devices := []models.Device{
		{Name: "dev1"}, {Name: "dev2"}, {Name: "dev3"}, {Name: "dev4"},
	}

	for _, device := range devices {
		_, err := repo.SaveNewDevice(device)
		assert.NoError(t, err)
	}

	fetchedDevices, err := repo.GetDevices()
	assert.NoError(t, err)
	for i := range fetchedDevices {
		assert.Equal(t, devices[i].Name, fetchedDevices[i].Name)
	}
}

func TestTrackingRepositoryMock_SaveNewTrackInput(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	devices := []models.Device{
		{Name: "dev1"}, {Name: "dev2"}, {Name: "dev3"}, {Name: "dev4"},
	}
	for _, device := range devices {
		_, err := repo.SaveNewDevice(device)
		assert.NoError(t, err)
	}

	t.Run("saving a track input for a non existing device should respond with error", func(t *testing.T) {
		_, err := repo.SaveNewTrackInput(models.TrackInput{
			Location: models.Location{
				Lat: 12.123,
				Lng: 32.123,
			},
			RecordedAt: time.Now(),
			DeviceID:   123,
		})
		assert.Equal(t, ErrDeviceDoesNotExist, err)
	})

	t.Run("saving an invalid track input for an existing device should respond with error", func(t *testing.T) {
		_, err := repo.SaveNewTrackInput(models.TrackInput{
			Location:   models.Location{},
			RecordedAt: time.Now(),
			DeviceID:   2,
		})
		assert.Error(t, err)
		assert.NotEqual(t, ErrDeviceDoesNotExist, err)
	})

	t.Run("saving a valid track input for an existing device should succeed", func(t *testing.T) {
		_, err := repo.SaveNewTrackInput(models.TrackInput{
			Location: models.Location{
				Lat: 12.123,
				Lng: 21.123,
			},
			RecordedAt: time.Now(),
			DeviceID:   2,
		})
		assert.NoError(t, err)
	})
}

func TestTrackingRepositoryMock_GetAllTrackInputsOfDevice(t *testing.T) {
	repo := NewTrackingRepositoryMock()

	dev1, err := repo.SaveNewDevice(models.Device{Name: "dev1"})
	assert.NoError(t, err)
	dev2, err := repo.SaveNewDevice(models.Device{Name: "dev2"})
	assert.NoError(t, err)

	trackInputs := make([]models.TrackInput, 0)
	for i := 0; i < 5; i++ {
		ti := models.TrackInput{
			Location: models.Location{
				Lat: 12.123 + float64(i),
				Lng: 21.123 + float64(i),
			},
			RecordedAt: time.Now().Add(time.Duration(i) * time.Minute),
			DeviceID:   dev2.ID,
		}

		trackInputs = append(trackInputs, ti)
		_, err := repo.SaveNewTrackInput(ti)
		assert.NoError(t, err)
	}

	t.Run("all the device inputs should be fetched", func(t *testing.T) {
		tis, err := repo.GetAllTrackInputsOfDevice(dev2.ID)
		assert.NoError(t, err)
		for i := range trackInputs {
			assert.Equal(t, trackInputs[i].DeviceID, tis[i].DeviceID)
			assert.Equal(t, trackInputs[i].Location, tis[i].Location)
			assert.Equal(t, trackInputs[i].RecordedAt, tis[i].RecordedAt)
		}
	})

	t.Run("an empty list of device inputs is fetched if none exists", func(t *testing.T) {
		tis, err := repo.GetAllTrackInputsOfDevice(dev1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(tis))
	})

	t.Run("if the specified device does not exist an error should be returned", func(t *testing.T) {
		_, err := repo.GetAllTrackInputsOfDevice(123)
		assert.Equal(t, ErrDeviceDoesNotExist, err)
	})
}
