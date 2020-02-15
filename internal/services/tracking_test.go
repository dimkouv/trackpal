// +build unit_test

package services

import (
	"encoding/json"
	"testing"

	"github.com/dimkouv/trackpal/internal/models"

	"github.com/stretchr/testify/assert"

	"github.com/dimkouv/trackpal/internal/repository"
)

func TestTrackingService_GetDevicesAsJSON(t *testing.T) {
	repo := repository.NewTrackingRepositoryMock()
	server := NewTrackingService(repo)
	devices := make([]models.Device, 0)

	t.Run("if no devices exist a null response is received", func(t *testing.T) {
		b, err := server.GetDevicesAsJSON()
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

		b, err := server.GetDevicesAsJSON()
		assert.NoError(t, err)

		err = json.Unmarshal(b, &devices)
		assert.NoError(t, err)

		for i := range createdDevices {
			assert.Equal(t, createdDevices[i].Name, devices[i].Name)
		}
	})
}
