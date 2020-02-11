package repository

import (
	"errors"
	"sort"

	"github.com/dimkouv/trackpal/internal/models"
)

type TrackingRepositoryMock struct {
	trackInputs          []models.TrackInput
	trackInputsDeviceIDX map[int64][]models.TrackInput
}

var (
	ErrDeviceNotExists = errors.New("device does not exist")
)

func (repo *TrackingRepositoryMock) SaveNewTrackInput(trackInput models.TrackInput) {
	repo.trackInputs = append(repo.trackInputs, trackInput)

	if _, exists := repo.trackInputsDeviceIDX[trackInput.DeviceID]; !exists {
		repo.trackInputsDeviceIDX[trackInput.DeviceID] = make([]models.TrackInput, 0)
	}
	repo.trackInputsDeviceIDX[trackInput.DeviceID] = append(repo.trackInputsDeviceIDX[trackInput.DeviceID], trackInput)
}

func (repo *TrackingRepositoryMock) GetAllTrackInputsOfDevice(deviceID int64) ([]models.TrackInput, error) {
	trackInputs, exists := repo.trackInputsDeviceIDX[deviceID]
	if !exists {
		return nil, ErrDeviceNotExists
	}

	sort.Slice(trackInputs, func(i, j int) bool {
		return trackInputs[i].Timestamp.Before(trackInputs[j].Timestamp)
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
