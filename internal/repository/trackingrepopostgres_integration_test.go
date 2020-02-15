// +build integration_test

package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/dimkouv/trackpal/internal/envlib"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewTrackingRepositoryPostgres(t *testing.T) {
	_, err := NewTrackingRepositoryPostgres("invalid")
	assert.Error(t, err)

	repo, err := NewTrackingRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)
	assert.NotNil(t, repo.db)
}

func TestTrackingRepositoryPostgres_SaveNewDevice(t *testing.T) {
	repo, err := NewTrackingRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

	t.Run("trying to save an invalid device should raise an error", func(t *testing.T) {
		_, err := repo.SaveNewDevice(models.Device{
			ID:   0,
			Name: "",
		})
		assert.Error(t, err)
	})

	t.Run("saving a valid device should be successful", func(t *testing.T) {
		d, err := repo.SaveNewDevice(models.Device{
			ID:   -1, // the passed id is ignored
			Name: "my-device",
		})
		assert.NoError(t, err)
		assert.Greater(t, d.ID, int64(0))
	})
}

func TestTrackingRepositoryPostgres_GetDevices(t *testing.T) {
	repo, err := NewTrackingRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

	devices := []models.Device{
		{
			Name: fmt.Sprintf("%v", time.Now().UnixNano()),
		},
		{
			Name: fmt.Sprintf("%v", time.Now().UnixNano()),
		},
	}

	for _, device := range devices {
		_, err := repo.SaveNewDevice(device)
		assert.NoError(t, err)
	}

	fetchedDevices, err := repo.GetDevices()
	assert.NoError(t, err)
	found := 0
	for i := range devices {
		for j := range fetchedDevices {
			if fetchedDevices[j].Name == devices[i].Name {
				found++
			}
		}
	}
	assert.Equal(t, len(devices), found)
}

func TestTrackingRepositoryPostgres_SaveNewTrackInput(t *testing.T) {
	repo, err := NewTrackingRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

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

func TestTrackingRepositoryPostgres_GetAllTrackInputsOfDevice(t *testing.T) {
	repo, err := NewTrackingRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

	dev1, err := repo.SaveNewDevice(models.Device{Name: fmt.Sprintf("%v", time.Now().Unix())})
	assert.NoError(t, err)
	dev2, err := repo.SaveNewDevice(models.Device{Name: fmt.Sprintf("%v", time.Now().Unix())})
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
		t.Logf("%v", tis)
	})

	t.Run("an empty list of device inputs is fetched if none exists", func(t *testing.T) {
		tis, err := repo.GetAllTrackInputsOfDevice(dev1.ID)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(tis))
	})

	t.Run("if the specified device does not exist an error should be returned", func(t *testing.T) {
		_, err := repo.GetAllTrackInputsOfDevice(-1)
		assert.Equal(t, ErrDeviceDoesNotExist, err)
	})
}
