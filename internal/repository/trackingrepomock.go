package repository

import (
	"sort"

	"github.com/dimkouv/trackpal/internal/models"
)

type TrackingRepositoryMock struct {
	devices     []models.Device
	trackInputs []models.TrackInput

	trackInputsDeviceIDX map[int64][]models.TrackInput
}

func (repo *TrackingRepositoryMock) SaveNewDevice(d models.Device) (*models.Device, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	deviceID := int64(len(repo.devices) + 1)
	d.ID = deviceID
	repo.devices = append(repo.devices, d)
	repo.trackInputsDeviceIDX[deviceID] = make([]models.TrackInput, 0)
	return &d, nil
}

func (repo *TrackingRepositoryMock) GetDevices() ([]models.Device, error) {
	return repo.devices, nil
}

func (repo *TrackingRepositoryMock) SaveNewTrackInput(trackInput models.TrackInput) (*models.TrackInput, error) {
	repo.trackInputs = append(repo.trackInputs, trackInput)
	_, exists := repo.trackInputsDeviceIDX[trackInput.DeviceID]
	if !exists {
		return nil, ErrDeviceDoesNotExist
	}

	if err := trackInput.Validate(); err != nil {
		return nil, err
	}

	repo.trackInputsDeviceIDX[trackInput.DeviceID] = append(repo.trackInputsDeviceIDX[trackInput.DeviceID], trackInput)
	return &trackInput, nil
}

func (repo *TrackingRepositoryMock) GetAllTrackInputsOfDevice(deviceID int64) ([]models.TrackInput, error) {
	trackInputs, exists := repo.trackInputsDeviceIDX[deviceID]
	if !exists {
		return nil, ErrDeviceDoesNotExist
	}

	sort.Slice(trackInputs, func(i, j int) bool {
		return trackInputs[i].RecordedAt.Before(trackInputs[j].RecordedAt)
	})

	return trackInputs, nil
}

// NewTrackingRepositoryMock returns a new instance of a mock tracking repository
func NewTrackingRepositoryMock() *TrackingRepositoryMock {
	return &TrackingRepositoryMock{
		trackInputs:          make([]models.TrackInput, 0),
		trackInputsDeviceIDX: make(map[int64][]models.TrackInput),
	}
}
