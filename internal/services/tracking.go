package services

import (
	"encoding/json"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
)

type TrackingService struct {
	repo repository.TrackingRepository
}

func (service TrackingService) GetDevicesAsJSON() ([]byte, error) {
	results, err := service.repo.GetDevices()
	if err != nil {
		return nil, err
	}
	return json.Marshal(results)
}

func (service TrackingService) SaveDevice(d models.Device) ([]byte, error) {
	device, err := service.repo.SaveNewDevice(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(device)
}

func (service TrackingService) SaveTrackInput(t models.TrackInput) ([]byte, error) {
	ti, err := service.repo.SaveNewTrackInput(t)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ti)
}

func (service TrackingService) GetAllTrackInputsOfDeviceAsJSON(deviceID int64) ([]byte, error) {
	results, err := service.repo.GetAllTrackInputsOfDevice(deviceID)
	if err != nil {
		return nil, err
	}
	return json.Marshal(results)
}

// NewTrackingService receives a repository and returns a tracking service
func NewTrackingService(repo repository.TrackingRepository) TrackingService {
	return TrackingService{repo: repo}
}

// NewTrackingService returns a tracking service with a mock repository
func NewTrackingServiceMock() TrackingService {
	return TrackingService{repo: repository.NewTrackingRepositoryMock()}
}
