package server

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/conf"
)

func (ts TrackpalServer) authRoutes() []Route {
	return []Route{
		{
			Pattern:     "/auth/register",
			HandlerFunc: ts.authRegister,
			Method:      "POST",
			Name:        "authRegister",
		},
		{
			Pattern:     "/auth/activate",
			HandlerFunc: ts.authActivate,
			Method:      "POST",
			Name:        "authActivate",
		},
		{
			Pattern:     "/auth/refresh",
			HandlerFunc: ts.withUser(ts.authRefresh),
			Method:      "POST",
			Name:        "authRefresh",
		},
		{
			Pattern:     "/auth/login",
			HandlerFunc: ts.authLogin,
			Method:      "POST",
			Name:        "authLogin",
		},
	}
}

func (ts TrackpalServer) trackingRoutes() []Route {
	return []Route{
		{
			Pattern:     "/tracking/devices",
			HandlerFunc: ts.withUser(ts.getDevices),
			Method:      "GET",
			Name:        "getTrackingDevices",
		},
		{
			Pattern:     "/tracking/devices",
			HandlerFunc: ts.withUser(ts.createDevice),
			Method:      "POST",
			Name:        "postTrackingDevice",
		},
		{
			Pattern:     "/tracking/devices/{deviceID:[0-9]+}/records",
			HandlerFunc: ts.withUser(ts.getTrackRecordsOfDevice),
			Method:      "GET",
			Name:        "getTrackingDeviceRecords",
		},
		{
			Pattern:     "/tracking/devices/{deviceID:[0-9]+}/records",
			HandlerFunc: ts.withUser(ts.addTrackRecordOfDevice),
			Method:      "POST",
			Name:        "postTrackingDeviceRecord",
		},
	}
}

func (ts TrackpalServer) alertingRoutes() []Route {
	return []Route{
		{
			Pattern:     "/tracking/devices/{deviceID:[0-9]+}/alerting/enable",
			HandlerFunc: ts.withUser(ts.enableAlerting),
			Method:      "POST",
			Name:        "postAlertingEnable",
		},
	}
}

// RegisterRoutes register all the routes that are declared in this package
func (ts TrackpalServer) RegisterRoutes() *mux.Router {
	ts.routes = []Route{}
	ts.routes = append(ts.routes, ts.authRoutes()...)
	ts.routes = append(ts.routes, ts.trackingRoutes()...)
	ts.routes = append(ts.routes, ts.alertingRoutes()...)

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

// ListenAndServe starts listening for incoming requests
func (ts TrackpalServer) ListenAndServe(addr string, router http.Handler) {
	logrus.Infof("Server running: addr=%s", addr)
	if err := http.ListenAndServe(addr, handlers.CORS(
		handlers.AllowedOrigins([]string{conf.AccessControlAllowOrigin}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept", "Origin", "Authorization"}),
	)(router)); err != nil {
		panic(err)
	}
}
