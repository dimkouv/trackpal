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
			"GetTrackRecords",
			"GET",
			"/track/{deviceID:[0-9]+}",
			ts.getTrackRecordsOfDevice,
		},
		{
			"NewTrackRecord",
			"POST",
			"/track/{deviceID:[0-9]+}",
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
