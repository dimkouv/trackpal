package server

import (
	"net/http"

	"github.com/dimkouv/trackpal/pkg/response"
	"github.com/dimkouv/trackpal/pkg/terror"
)

// responseBasedOnErr writes to w based on the error
func (ts TrackpalServer) responseBasedOnErr(err error, b []byte, succStatus int, w http.ResponseWriter) {
	terr, isTerr := err.(terror.Terror)

	switch {
	case err == nil:
		response.HTTP(w).Data(b).Status(succStatus).JSON()
	case !isTerr:
		response.HTTP(w).Status(http.StatusInternalServerError).TEXT()
	default:
		response.HTTP(w).Status(terr.Code()).Error(terr).JSON()
	}
}
