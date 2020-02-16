package services

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/dimkouv/trackpal/internal/consts"

	"github.com/dimkouv/trackpal/pkg/terror"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
)

const (
	ErrPlain = iota
	ErrBodyRead
	ErrBodyParse
	ErrVarNotFound
	ErrMarshal
	ErrPermissionDenied
	ErrNotValid
)

type TrackingService struct {
	ua   models.UserAccount
	repo repository.TrackingRepository
}

func (service *TrackingService) SetUser(ua models.UserAccount) {
	service.ua = ua
}

func (service *TrackingService) GetUser() models.UserAccount {
	return service.ua
}

func (service TrackingService) GetDevicesAsJSON() ([]byte, error) {
	results, err := service.repo.GetDevices(service.ua.ID)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("error getting devices")
		return nil, terror.New(ErrPlain, err.Error())
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal results")
		return nil, terror.New(ErrMarshal, err.Error())
	}

	return b, nil
}

func (service TrackingService) SaveDevice(rc io.Reader) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return nil, terror.New(ErrBodyRead, err.Error())
	}

	d := models.Device{}
	err = json.Unmarshal(requestData, &d)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")

		return nil, terror.New(ErrBodyParse, err.Error())
	}

	if err := d.Validate(); err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("the provided device is not valid")
		return nil, terror.New(ErrNotValid, err.Error())
	}

	d.UserID = service.ua.ID
	device, err := service.repo.SaveNewDevice(d)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to save new device")
		return nil, terror.New(ErrPlain, err.Error())
	}

	b, err := json.Marshal(device)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal device")
		return nil, terror.New(ErrMarshal, err.Error())
	}

	return b, nil
}

func (service TrackingService) SaveTrackInput(vars map[string]string, rc io.Reader) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return nil, terror.New(ErrBodyRead, err.Error())
	}

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.
			WithField(consts.LogFieldVars, vars).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse deviceID")
		return nil, terror.New(ErrVarNotFound, err.Error())
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return nil, terror.New(ErrPlain, err.Error())
	}

	if device.UserID != service.ua.ID {
		logrus.
			WithField("user_id", service.ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to save track input")
		return nil, terror.New(ErrPermissionDenied, "device owner is not the claimed one")
	}

	t := models.TrackInput{}
	err = json.Unmarshal(requestData, &t)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return nil, terror.New(ErrBodyParse, err.Error())
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
		return nil, terror.New(ErrMarshal, err.Error())
	}

	return b, nil
}

func (service TrackingService) GetAllTrackInputsOfDeviceAsJSON(vars map[string]string) ([]byte, error) {
	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.
			WithField(consts.LogFieldVars, vars).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse deviceID")
		return nil, terror.New(ErrVarNotFound, err.Error())
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return nil, terror.New(ErrPlain, err.Error())
	}

	if device.UserID != service.ua.ID {
		logrus.
			WithField("user_id", service.ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to get track inputs of that device")
		return nil, terror.New(ErrPermissionDenied, "device owner is not the claimed one")
	}

	results, err := service.repo.GetAllTrackInputsOfDevice(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get track inputs of device")
		return nil, terror.New(ErrPlain, err.Error())
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal track input")
		return nil, terror.New(ErrMarshal, err.Error())
	}

	return b, nil
}

// NewTrackingService receives a repository and returns a tracking service
func NewTrackingService(repo repository.TrackingRepository) TrackingService {
	return TrackingService{repo: repo}
}

// NewTrackingServicePostgres returns a tracking service with a postgres repository
func NewTrackingServicePostgres(postgresDSN string) (*TrackingService, error) {
	repo, err := repository.NewTrackingRepositoryPostgres(postgresDSN)
	if err != nil {
		return nil, err
	}

	return &TrackingService{repo: repo}, nil
}
