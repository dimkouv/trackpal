package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (ts TrackpalServer) getDevices(w http.ResponseWriter, req *http.Request) {
	b, err := ts.trackingService.GetDevicesAsJSON(req.Context())
	ts.responseBasedOnErr(err, b, http.StatusOK, w)
}

func (ts TrackpalServer) createDevice(w http.ResponseWriter, req *http.Request) {
	b, err := ts.trackingService.SaveDevice(req.Context(), req.Body)
	ts.responseBasedOnErr(err, b, http.StatusCreated, w)
}

func (ts TrackpalServer) getTrackRecordsOfDevice(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b, err := ts.trackingService.GetAllTrackInputsOfDeviceAsJSON(req.Context(), vars)
	ts.responseBasedOnErr(err, b, http.StatusOK, w)
}

func (ts TrackpalServer) addTrackRecordOfDevice(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b, err := ts.trackingService.SaveTrackInput(req.Context(), vars, req.Body)
	ts.responseBasedOnErr(err, b, http.StatusCreated, w)
}

func (ts TrackpalServer) enableAlerting(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	err := ts.trackingService.EnableAlerting(req.Context(), vars, req.Body)
	ts.responseBasedOnErr(err, nil, http.StatusOK, w)
}
