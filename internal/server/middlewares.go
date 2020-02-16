package server

import (
	"net/http"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/response"
)

func (ts TrackpalServer) withUser(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		ua, err := models.NewUserAccount().FromJWT(req.Header.Get(consts.Authorization))
		if err != nil {
			response.HTTP(w).Error(err).Status(http.StatusUnauthorized).JSON()
			return
		}
		ts.trackingService.SetUser(*ua)

		next.ServeHTTP(w, req)
	}

	return f
}
