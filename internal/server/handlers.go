package server

import (
	"net/http"

	"github.com/dimkouv/trackpal/pkg/response"

	"github.com/gorilla/mux"
)

func (ts TrackpalServer) getDevices(w http.ResponseWriter, req *http.Request) {
	b, err := ts.trackingService.GetDevicesAsJSON()

	switch err {
	case nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	default:
		response.HTTP(w).Error(err).Status(http.StatusBadRequest).JSON()
	}
}

func (ts TrackpalServer) createDevice(w http.ResponseWriter, req *http.Request) {
	b, err := ts.trackingService.SaveDevice(req.Body)

	switch err {
	case nil:
		response.HTTP(w).Data(b).Status(http.StatusCreated).JSON()
	default:
		response.HTTP(w).Error(err).Status(http.StatusBadRequest).JSON()
	}
}

func (ts TrackpalServer) getTrackRecordsOfDevice(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b, err := ts.trackingService.GetAllTrackInputsOfDeviceAsJSON(vars)

	switch err {
	case nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	default:
		response.HTTP(w).Error(err).Status(http.StatusBadRequest).JSON()
	}
}

func (ts TrackpalServer) addTrackRecordOfDevice(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b, err := ts.trackingService.SaveTrackInput(vars, req.Body)

	switch err {
	case nil:
		response.HTTP(w).Data(b).Status(http.StatusCreated).JSON()
	default:
		response.HTTP(w).Error(err).Status(http.StatusBadRequest).JSON()
	}
}

func (ts TrackpalServer) authRegister(w http.ResponseWriter, req *http.Request) {
	err := ts.userService.CreateUserAccount(req.Body)

	switch err {
	case nil:
		response.HTTP(w).Status(http.StatusCreated).JSON()
	default:
		response.HTTP(w).Error(err).Status(http.StatusBadRequest).JSON()
	}
}

func (ts TrackpalServer) authLogin(w http.ResponseWriter, req *http.Request) {
	b, err := ts.userService.GetJWTFromEmailAndPassword(req.Body)

	switch err {
	case nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	default:
		response.HTTP(w).Error(err).Status(http.StatusBadRequest).JSON()
	}
}
