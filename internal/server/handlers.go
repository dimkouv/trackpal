package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dimkouv/trackpal/internal/models"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (ts TrackpalServer) getDevices(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ts.trackingService.GetDevicesAsJSON()
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to get devices")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to write response")
	}
}

func (ts TrackpalServer) createDevice(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	requestData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.WithField("error", err).WithField("error", err).Errorf("unable to read request body")
		return
	}

	device := models.Device{}
	err = json.Unmarshal(requestData, &device)
	if err != nil {
		logrus.WithField("body", fmt.Sprintf("%s", requestData)).
			WithField("error", err).Errorf("unable to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return
	}

	b, err := ts.trackingService.SaveDevice(device)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to write response")
	}
}

func (ts TrackpalServer) getTrackRecordsOfDevice(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.WithField("vars", vars).WithField("error", err).Errorf("unable to parse device id")
		return
	}

	b, err := ts.trackingService.GetAllTrackInputsOfDeviceAsJSON(int64(deviceID))
	if err != nil {
		logrus.WithField("error", err).WithField("deviceID", deviceID).
			Errorf("unable to get track inputs of device")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to write response")
	}
}

func (ts TrackpalServer) addTrackRecordOfDevice(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)

	deviceID, err := strconv.Atoi(vars["deviceID"])
	if err != nil {
		logrus.WithField("vars", vars).WithField("error", err).Errorf("unable to parse device id")
		return
	}

	requestData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logrus.WithField("error", err).WithField("error", err).Errorf("unable to read request body")
		return
	}

	trackInput := models.TrackInput{}
	err = json.Unmarshal(requestData, &trackInput)
	if err != nil {
		logrus.WithField("body", fmt.Sprintf("%s", requestData)).
			WithField("error", err).Errorf("unable to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return
	}
	trackInput.DeviceID = int64(deviceID)

	b, err := ts.trackingService.SaveTrackInput(trackInput)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to parse request body")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"msg": "%s"}`, err.Error())))
		return
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to write response")
	}
}
