package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// RegisterRoutes register all the routes that are declared in this package
func (ts TrackpalServer) RegisterRoutes() *mux.Router {
	ts.routes = []Route{
		{
			"GetDevices",
			"GET",
			"/devices",
			ts.getDevices,
		},
		{
			"CreateDevice",
			"POST",
			"/devices",
			ts.createDevice,
		},
		{
			"GetTrackRecords",
			"GET",
			"/devices/{deviceID:[0-9]+}/entries",
			ts.getTrackRecordsOfDevice,
		},
		{
			"NewTrackRecord",
			"POST",
			"/devices/{deviceID:[0-9]+}/entries",
			ts.addTrackRecordOfDevice,
		},
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range ts.routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// Listen starts listening for incoming requests
func (ts TrackpalServer) ListenAndServe(addr string, router *mux.Router) {
	logrus.Infof("Server running: addr=%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
