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
			Name:        "RegisterAccount",
			Method:      "POST",
			Pattern:     "/auth/register",
			HandlerFunc: ts.authRegister,
		},
		{
			Name:        "RefreshToken",
			Method:      "POST",
			Pattern:     "/auth/refresh",
			HandlerFunc: ts.withUser(ts.authRefresh),
		},
		{
			Name:        "RegisterAccount",
			Method:      "POST",
			Pattern:     "/auth/login",
			HandlerFunc: ts.authLogin,
		},
		{
			Name:        "GetDevices",
			Method:      "GET",
			Pattern:     "/devices",
			HandlerFunc: ts.withUser(ts.getDevices),
		},
		{
			Name:        "CreateDevice",
			Method:      "POST",
			Pattern:     "/devices",
			HandlerFunc: ts.withUser(ts.createDevice),
		},
		{
			Name:        "GetTrackRecords",
			Method:      "GET",
			Pattern:     "/devices/{deviceID:[0-9]+}/entries",
			HandlerFunc: ts.withUser(ts.getTrackRecordsOfDevice),
		},
		{
			Name:        "NewTrackRecord",
			Method:      "POST",
			Pattern:     "/devices/{deviceID:[0-9]+}/entries",
			HandlerFunc: ts.withUser(ts.addTrackRecordOfDevice),
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
func (ts TrackpalServer) ListenAndServe(addr string, router http.Handler) {
	logrus.Infof("Server running: addr=%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
