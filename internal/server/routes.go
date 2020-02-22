package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

// RegisterRoutes register all the routes that are declared in this package
func (ts TrackpalServer) RegisterRoutes() *mux.Router {
	ts.routes = []Route{}
	ts.routes = append(ts.routes, ts.authRoutes()...)
	ts.routes = append(ts.routes, ts.trackingRoutes()...)

	router := mux.NewRouter().StrictSlash(true)
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.WriteHeader(http.StatusNoContent)
	})

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
	if err := http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
