package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/dimkouv/trackpal/internal/consts"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
)

type TrackingService struct {
	repo repository.TrackingRepository
}

func (service TrackingService) GetDevicesAsJSON(ctx context.Context) ([]byte, error) {
	ua := ctx.Value("user").(models.UserAccount)

	results, err := service.repo.GetDevices(ua.ID)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("error getting devices")
		return nil, err
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal results")
		return nil, err
	}

	return b, nil
}

func (service TrackingService) SaveDevice(ctx context.Context, rc io.Reader) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return nil, err
	}

	d := models.Device{}
	err = json.Unmarshal(requestData, &d)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")

		return nil, err
	}

	err = d.Validate()
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("the provided device is not valid")
		return nil, err
	}

	ua := ctx.Value("user").(models.UserAccount)
	d.UserID = ua.ID
	device, err := service.repo.SaveNewDevice(d)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to save new device")
		return nil, err
	}

	b, err := json.Marshal(device)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal device")
		return nil, err
	}

	return b, nil
}

func (service TrackingService) SaveTrackInput(
	ctx context.Context, vars map[string]string, rc io.Reader) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return nil, err
	}

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.
			WithField(consts.LogFieldVars, vars).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse deviceID")
		return nil, err
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return nil, err
	}

	ua := ctx.Value("user").(models.UserAccount)
	if device.UserID != ua.ID {
		logrus.
			WithField("user_id", ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to save track input")
		return nil, err
	}

	t := models.TrackInput{}
	err = json.Unmarshal(requestData, &t)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return nil, err
	}
	t.DeviceID = int64(deviceID)

	ti, err := service.repo.SaveNewTrackInput(t)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(ti)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal track input")
		return nil, err
	}

	return b, nil
}

func (service TrackingService) GetAllTrackInputsOfDeviceAsJSON(
	ctx context.Context, vars map[string]string) ([]byte, error) {
	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.
			WithField(consts.LogFieldVars, vars).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse deviceID")
		return nil, err
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return nil, err
	}

	ua := ctx.Value("user").(models.UserAccount)
	if device.UserID != ua.ID {
		logrus.
			WithField("user_id", ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to get track inputs of that device")
		return nil, err
	}

	results, err := service.repo.GetAllTrackInputsOfDevice(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get track inputs of device")
		return nil, err
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal track input")
		return nil, err
	}

	return b, nil
}

// NewTrackingService receives a repository and returns a tracking service
func NewTrackingService(repo repository.TrackingRepository) *TrackingService {
	return &TrackingService{repo: repo}
}

// NewTrackingServicePostgres returns a tracking service with a postgres repository
func NewTrackingServicePostgres(postgresDSN string) (*TrackingService, error) {
	repo, err := repository.NewTrackingRepositoryPostgres(postgresDSN)
	if err != nil {
		return nil, err
	}

	return &TrackingService{repo: repo}, nil
}
