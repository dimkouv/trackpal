package services

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/sirupsen/logrus"

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

func (service TrackingService) SaveDevice(rc io.ReadCloser) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.WithField("error", err).WithField("error", err).Errorf("unable to read request body")
		return nil, err
	}

	d := models.Device{}
	err = json.Unmarshal(requestData, &d)
	if err != nil {
		logrus.WithField("body", fmt.Sprintf("%s", requestData)).
			WithField("error", err).Errorf("unable to parse request body")
		return nil, fmt.Errorf("unable to parse body: %v", err)
	}

	device, err := service.repo.SaveNewDevice(d)
	if err != nil {
		return nil, err
	}
	return json.Marshal(device)
}

func (service TrackingService) SaveTrackInput(vars map[string]string, rc io.ReadCloser) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.WithField("error", err).WithField("error", err).Errorf("unable to read request body")
		return nil, err
	}

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.WithField("vars", vars).WithField("error", err).Errorf("unable to parse device id")
		return nil, err
	}

	t := models.TrackInput{}
	err = json.Unmarshal(requestData, &t)
	if err != nil {
		logrus.WithField("body", fmt.Sprintf("%s", requestData)).
			WithField("error", err).Errorf("unable to parse request body")
		return nil, err
	}
	t.DeviceID = int64(deviceID)

	ti, err := service.repo.SaveNewTrackInput(t)
	if err != nil {
		return nil, err
	}
	return json.Marshal(ti)
}

func (service TrackingService) GetAllTrackInputsOfDeviceAsJSON(vars map[string]string) ([]byte, error) {
	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.WithField("vars", vars).WithField("error", err).Errorf("unable to parse device id")
		return nil, err
	}

	results, err := service.repo.GetAllTrackInputsOfDevice(int64(deviceID))
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
