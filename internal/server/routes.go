package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// RegisterRoutes register all the routes that are declared in this package
func (ts TrackpalServer) RegisterRoutes() {
	routes := map[string]http.HandlerFunc{
		"/info": ts.info,
	}

	for routeName, handler := range routes {
		logrus.Debugf("Register route <%s>", routeName)
		http.HandleFunc(routeName, handler)
	}
}

// Listen starts listening for incoming requests
func (ts TrackpalServer) Listen(addr string) {
	logrus.Infof("Server running: addr=%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
