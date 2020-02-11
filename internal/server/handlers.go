package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

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
		return
	}

	_, err = w.Write(b)
	if err != nil {
		logrus.WithField("error", err).Errorf("unable to write response")
	}
}

func (ts TrackpalServer) addTrackRecordOfDevice(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write([]byte("hello friend"))
}
