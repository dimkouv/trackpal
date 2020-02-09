package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (ts TrackpalServer) info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"hello": "friend"}`))

	if err != nil {
		logrus.Errorf("err=%v", err)
	}
}
