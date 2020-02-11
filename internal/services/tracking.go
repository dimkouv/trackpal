package services

import (
	"encoding/json"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
)

type TrackingService struct {
	repo repository.TrackingRepository
}

func (service TrackingService) SaveTrackInput(t models.TrackInput) {
	service.repo.SaveNewTrackInput(t)
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
