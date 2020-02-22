package server

import (
	"net/http"

	"github.com/dimkouv/trackpal/pkg/terror"

	"github.com/dimkouv/trackpal/pkg/response"

	"github.com/gorilla/mux"
)

func (ts TrackpalServer) getDevices(w http.ResponseWriter, req *http.Request) {
	b, err := ts.trackingService.GetDevicesAsJSON(req.Context())
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) createDevice(w http.ResponseWriter, req *http.Request) {
	b, err := ts.trackingService.SaveDevice(req.Context(), req.Body)
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(http.StatusCreated).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) getTrackRecordsOfDevice(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b, err := ts.trackingService.GetAllTrackInputsOfDeviceAsJSON(req.Context(), vars)
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) addTrackRecordOfDevice(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	b, err := ts.trackingService.SaveTrackInput(req.Context(), vars, req.Body)
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(http.StatusCreated).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) authRegister(w http.ResponseWriter, req *http.Request) {
	err := ts.userService.CreateUserAccount(req.Context(), req.Body)
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Status(http.StatusCreated).TEXT()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) authActivate(w http.ResponseWriter, req *http.Request) {
	err := ts.userService.ActivateUserAccount(req.Context(), req.Body)
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Status(http.StatusOK).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) authLogin(w http.ResponseWriter, req *http.Request) {
	b, err := ts.userService.GetJWTFromEmailAndPassword(req.Context(), req.Body)
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}

func (ts TrackpalServer) authRefresh(w http.ResponseWriter, req *http.Request) {
	b, err := ts.userService.RefreshJWT(req.Context())
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(http.StatusOK).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}
