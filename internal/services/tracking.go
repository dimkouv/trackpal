package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/dimkouv/trackpal/internal/consts"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
)

type TrackingService struct {
	repo repository.TrackingRepository
}

func (service TrackingService) GetDevicesAsJSON(ctx context.Context) ([]byte, error) {
	ua := ctx.Value(consts.CtxUser).(models.UserAccount)

	results, err := service.repo.GetDevices(ua.ID)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("error getting devices")
		return nil, consts.ErrEnumInternal
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal results")
		return nil, consts.ErrEnumInternal
	}

	return b, nil
}

func (service TrackingService) SaveDevice(ctx context.Context, rc io.Reader) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return nil, consts.ErrEnumInvalidBody
	}

	d := models.Device{}
	err = json.Unmarshal(requestData, &d)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return nil, consts.ErrEnumInvalidBody
	}

	err = d.Validate()
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("the provided device is not valid")
		return nil, consts.ErrEnumInvalidData
	}

	ua := ctx.Value(consts.CtxUser).(models.UserAccount)
	d.UserID = ua.ID
	device, err := service.repo.SaveNewDevice(d)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to save new device")
		return nil, consts.ErrEnumInternal
	}

	b, err := json.Marshal(device)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal device")
		return nil, consts.ErrEnumInternal
	}

	return b, nil
}

func (service TrackingService) SaveTrackInput(
	ctx context.Context, vars map[string]string, rc io.Reader) ([]byte, error) {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return nil, consts.ErrEnumInvalidBody
	}

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.WithField(consts.LogFieldVars, vars).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse deviceID")
		return nil, consts.ErrEnumInvalidVars
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return nil, consts.ErrEnumNotFound
	}

	ua := ctx.Value(consts.CtxUser).(models.UserAccount)
	if device.UserID != ua.ID {
		logrus.WithField("user_id", ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to save track input")
		return nil, consts.ErrEnumNotFound
	}

	t := models.TrackInput{}
	err = json.Unmarshal(requestData, &t)
	if err != nil {
		logrus.WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return nil, consts.ErrEnumInvalidBody
	}
	t.DeviceID = int64(deviceID)

	ti, err := service.repo.SaveNewTrackInput(t)
	if err != nil {
		return nil, consts.ErrEnumInternal
	}

	go service.checkForAlert(ua, device, ti)

	b, err := json.Marshal(ti)
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal track input")
		return nil, consts.ErrEnumInternal
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
		return nil, consts.ErrEnumInvalidVars
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return nil, consts.ErrEnumNotFound
	}

	ua := ctx.Value(consts.CtxUser).(models.UserAccount)
	if device.UserID != ua.ID {
		logrus.
			WithField("user_id", ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to get track inputs of that device")
		return nil, consts.ErrEnumNotFound
	}

	results, err := service.repo.GetAllTrackInputsOfDevice(int64(deviceID))
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get track inputs of device")
		return nil, consts.ErrEnumInternal
	}

	b, err := json.Marshal(results)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to marshal track input")
		return nil, consts.ErrEnumInternal
	}

	return b, nil
}

// EnableAlerting enables alerting for a device. Enabling alerting locks the current position
// of the device (specified in the arguments), after that if the device is moved
// an alert is sent to the owner of the device.
func (service TrackingService) EnableAlerting(
	ctx context.Context, vars map[string]string, rc io.Reader) error {
	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return consts.ErrEnumInvalidBody
	}

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.WithField(consts.LogFieldVars, vars).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse deviceID")
		return consts.ErrEnumInvalidVars
	}

	device, err := service.repo.GetDeviceByID(int64(deviceID))
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).
			Errorf("unable to get device by id")
		return consts.ErrEnumNotFound
	}

	ua := ctx.Value(consts.CtxUser).(models.UserAccount)
	if device.UserID != ua.ID {
		logrus.WithField("user_id", ua.ID).
			WithField("device_owner", device.UserID).
			Errorf("unauthorized to enable alerting")
		return consts.ErrEnumNotFound
	}

	t := models.TrackInput{}
	err = json.Unmarshal(requestData, &t)
	if err != nil || t.Lat == 0 || t.Lng == 0 {
		logrus.WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return consts.ErrEnumInvalidBody
	}

	device.Lat, device.Lng = t.Lat, t.Lng
	device.AlertingEnabled = true
	err = service.repo.UpdateDevice(int64(deviceID), *device)
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).Errorf("unable to update device")
		return consts.ErrEnumInternal
	}

	return nil
}

// checkForAlert checks whether or not the user should be alerted based on the trackInput
func (service TrackingService) checkForAlert(
	ua models.UserAccount, device *models.Device, trackInput *models.TrackInput) {
	if !device.AlertingEnabled {
		return
	}

	if time.Since(device.LastAlertTimestamp).Minutes() < 5.0 {
		return
	}

	if trackInput.Location.HasMovedMoreThanM(device.Location, 3) {
		logrus.Infof("[alert] ua=%v device=%v", ua, device)
		// TODO: email/sms/...

		device.LastAlertTimestamp = time.Now().UTC()
		if err := service.repo.UpdateDevice(device.ID, *device); err != nil {
			logrus.WithField(consts.LogFieldErr, err).
				Errorf("unable to update device")
			return
		}
	}
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
