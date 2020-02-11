package repository

import "github.com/dimkouv/trackpal/internal/models"

// TrackingRepository contains methods for CRUD operations on tracking inputs
type TrackingRepository interface {
	// SaveNewTrackInput stores a new TrackInput
	SaveNewTrackInput(t models.TrackInput) (*models.TrackInput, error)

	// GetAllTrackInputsOfDevice returns all the tracking inputs that were recorded from the target device
	GetAllTrackInputsOfDevice(deviceID int64) ([]models.TrackInput, error)

	// SaveNewDevice stores a new Device
	SaveNewDevice(d models.Device) (*models.Device, error)

	// GetDevices returns all the devices
	GetDevices() ([]models.Device, error)
}
