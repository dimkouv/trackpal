// +build integration_test

package repository

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/xid"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/envlib"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/stretchr/testify/assert"
)

var (
	ua   *models.UserAccount
	repo *TrackingRepositoryPostgres
)

func TestMain(m *testing.M) {
	postgresDSN := envlib.GetEnvOrPanic(envlib.EnvPostgresDSN)

	postgresRepo, err := NewTrackingRepositoryPostgres(postgresDSN)
	if err != nil {
		logrus.Fatalf("unable to create repo: %v", err)
	}
	repo = postgresRepo

	uaRepo, err := NewAccountsRepositoryPostgres(postgresDSN)
	if err != nil {
		logrus.Fatalf("unable to create repo: %v", err)
	}

	user, err := uaRepo.SaveNewUser(models.UserAccount{
		Email: xid.New().String() + "@trackpal.com",
	}, "my-password")
	if err != nil {
		logrus.Fatalf("unable to create user: %v", err)
	}
	ua = user

	os.Exit(m.Run())
}

func TestNewTrackingRepositoryPostgres(t *testing.T) {
	_, err := NewTrackingRepositoryPostgres("invalid")
	assert.Error(t, err)

	repo, err := NewTrackingRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)
	assert.NotNil(t, repo.db)
}

func TestTrackingRepositoryPostgres_SaveNewDevice(t *testing.T) {
	t.Run("trying to save an invalid device should raise an error", func(t *testing.T) {
		_, err := repo.SaveNewDevice(models.Device{
			ID:     0,
			Name:   "",
			UserID: ua.ID,
		})
		assert.Error(t, err)
	})

	t.Run("saving a valid device should be successful", func(t *testing.T) {
		d, err := repo.SaveNewDevice(models.Device{
			ID:     -1, // the passed id is ignored
			Name:   "my-device",
			UserID: ua.ID,
		})
		if !assert.NoError(t, err) {
			fmt.Println(">>>", err)
		}
		assert.Greater(t, d.ID, int64(0))
	})
}

func TestTrackingRepositoryPostgres_GetDevices(t *testing.T) {
	devices := make([]models.Device, 0)
	for i := 0; i < 5; i++ {
		device, err := repo.SaveNewDevice(models.Device{
			Name:   xid.New().String(),
			UserID: ua.ID,
		})
		assert.NoError(t, err)
		devices = append(devices, *device)
	}

	t.Run("should fetch all the devices of the user", func(t *testing.T) {
		fetchedDevices, err := repo.GetDevices(ua.ID)
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
	})

	t.Run("should respond with an empty list if user does not exist", func(t *testing.T) {
		fetchedDevices, err := repo.GetDevices(-1)
		assert.NoError(t, err)
		assert.Empty(t, fetchedDevices)
	})
}

func TestTrackingRepositoryPostgres_SaveNewTrackInput(t *testing.T) {
	devices := []models.Device{
		{Name: "dev1"}, {Name: "dev2"}, {Name: "dev3"}, {Name: "dev4"},
	}
	for i, device := range devices {
		device.UserID = ua.ID
		dev, err := repo.SaveNewDevice(device)
		assert.NoError(t, err)
		devices[i] = *dev
	}

	t.Run("saving a track input for a non existing device should respond with error", func(t *testing.T) {
		_, err := repo.SaveNewTrackInput(models.TrackInput{
			Location: models.Location{
				Lat: 12.123,
				Lng: 32.123,
			},
			RecordedAt: time.Now(),
			DeviceID:   -1,
		})
		assert.Equal(t, ErrDeviceDoesNotExist, err)
	})

	t.Run("saving a valid track input for an existing device should succeed", func(t *testing.T) {
		_, err := repo.SaveNewTrackInput(models.TrackInput{
			Location: models.Location{
				Lat: 12.123,
				Lng: 21.123,
			},
			RecordedAt: time.Now(),
			DeviceID:   devices[0].ID,
		})
		assert.NoError(t, err)
	})
}

func TestTrackingRepositoryPostgres_GetAllTrackInputsOfDevice(t *testing.T) {
	dev1, err := repo.SaveNewDevice(models.Device{
		Name:   fmt.Sprintf("%v", time.Now().Unix()),
		UserID: ua.ID,
	})
	assert.NoError(t, err)
	dev2, err := repo.SaveNewDevice(models.Device{
		Name:   fmt.Sprintf("%v", time.Now().Unix()),
		UserID: ua.ID,
	})
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

func TestTrackingRepositoryPostgres_GetDeviceByID(t *testing.T) {
	devices := make([]models.Device, 0)

	for i := 0; i < 5; i++ {
		device, err := repo.SaveNewDevice(models.Device{
			Name:   xid.New().String(),
			UserID: ua.ID,
		})
		assert.NoError(t, err)
		devices = append(devices, *device)
	}

	t.Run("it should respond with the device if it exists", func(t *testing.T) {
		for i := range devices {
			fetchedDevice, err := repo.GetDeviceByID(devices[i].ID)
			assert.NoError(t, err)
			assert.Equal(t, devices[i].ID, fetchedDevice.ID)
			assert.Equal(t, devices[i].Name, fetchedDevice.Name)
		}
	})

	t.Run("should respond with error if device does not exist", func(t *testing.T) {
		_, err := repo.GetDeviceByID(-1234)
		assert.Equal(t, ErrDeviceDoesNotExist, err)
	})
}

func TestTrackingRepositoryPostgres_UpdateDevice(t *testing.T) {
	dev, err := repo.SaveNewDevice(models.Device{Name: "my device", UserID: ua.ID})
	assert.NoError(t, err)

	t.Run("should update the device details", func(t *testing.T) {
		err = repo.UpdateDevice(dev.ID, models.Device{
			Name:            "my new device name",
			AlertingEnabled: true,
		})
		assert.NoError(t, err)
		dev2, err := repo.GetDeviceByID(dev.ID)
		assert.NoError(t, err)
		assert.Equal(t, "my new device name", dev2.Name)
		assert.True(t, dev2.AlertingEnabled)
	})

	t.Run("should not overwrite the device id", func(t *testing.T) {
		err = repo.UpdateDevice(dev.ID, models.Device{ID: 123, Name: "my new name"})
		assert.NoError(t, err)
		dev2, err := repo.GetDeviceByID(dev.ID)
		if !assert.NoError(t, err) {
			fmt.Println(">>>", err)
		}
		assert.Equal(t, dev.ID, dev2.ID)
	})

}
