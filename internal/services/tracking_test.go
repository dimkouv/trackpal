// +build unit_test

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/dimkouv/trackpal/pkg/terror"

	"github.com/dimkouv/trackpal/internal/models"

	"github.com/stretchr/testify/assert"

	"github.com/dimkouv/trackpal/internal/repository"
)

func TestTrackingService_GetDevicesAsJSON(t *testing.T) {
	repo := repository.NewTrackingRepositoryMock()
	server := NewTrackingService(repo)
	devices := make([]models.Device, 0)

	t.Run("if no devices exist a null response is received", func(t *testing.T) {
		b, err := server.GetDevicesAsJSON(context.WithValue(context.Background(), "user", models.UserAccount{}))
		assert.NoError(t, err)
		err = json.Unmarshal(b, &devices)
		assert.NoError(t, err)
		assert.Empty(t, devices)
	})

	t.Run("if devices exists they should be received", func(t *testing.T) {
		createdDevices := []models.Device{
			{Name: "dev1"}, {Name: "dev2"}, {Name: "dev3"}, {Name: "dev4"},
		}
		for _, device := range createdDevices {
			_, err := repo.SaveNewDevice(device)
			assert.NoError(t, err)
		}

		b, err := server.GetDevicesAsJSON(context.WithValue(context.Background(), "user", models.UserAccount{}))
		assert.NoError(t, err)

		err = json.Unmarshal(b, &devices)
		assert.NoError(t, err)

		for i := range createdDevices {
			assert.Equal(t, createdDevices[i].Name, devices[i].Name)
		}
	})
}

func TestTrackingService_SaveDevice(t *testing.T) {
	repo := repository.NewTrackingRepositoryMock()
	server := NewTrackingService(repo)

	t.Run("should save and repond with the model succesfully", func(t *testing.T) {
		deviceAsJSON, err := json.Marshal(models.Device{Name: "my cool device"})
		assert.NoError(t, err)
		deviceAsJSON, err = server.SaveDevice(context.WithValue(context.Background(), "user", models.UserAccount{}), bytes.NewBufferString(string(deviceAsJSON)))
		assert.NoError(t, err)
		device := models.Device{}
		err = json.Unmarshal(deviceAsJSON, &device)
		assert.NoError(t, err)

		assert.Less(t, time.Now().UTC().Sub(device.CreatedAt).Seconds(), 2.0)
		assert.Equal(t, "my cool device", device.Name)
		assert.Greater(t, device.ID, int64(0))
		assert.Empty(t, device.UserID)
	})

	t.Run("should respond with ErrBodyRead", func(t *testing.T) {
		_, err := server.SaveDevice(context.Background(), IoReaderErrAlways{})
		terr := err.(terror.Terror)
		assert.Equal(t, ErrBodyRead, terr.Code())
	})

	t.Run("should respond with ErrBodyParse", func(t *testing.T) {
		cases := []string{
			"invalid json",
			"",
		}

		for _, c := range cases {
			_, err := server.SaveDevice(context.Background(), bytes.NewBufferString(c))
			terr := err.(terror.Terror)
			assert.Equal(t, ErrBodyParse, terr.Code())
		}
	})

	t.Run("should respond with ErrNotValid", func(t *testing.T) {
		cases := []string{
			`{"a": "b"}`,
			`{"name": ""}`,
			`{"name": "a"}`,
			`{"name": "aa"}`,
			`{"name": "aaa"}`,
		}

		for _, c := range cases {
			_, err := server.SaveDevice(context.Background(), bytes.NewBufferString(c))
			terr := err.(terror.Terror)
			assert.Equal(t, ErrNotValid, terr.Code())
		}
	})
}
