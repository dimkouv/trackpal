package server

import (
	"context"
	"net/http"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/response"
)

func (ts *TrackpalServer) withUser(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		ua, err := models.NewUserAccount().FromJWT(req.Header.Get(consts.Authorization))
		if err != nil {
			response.HTTP(w).Error(err).Status(http.StatusUnauthorized).JSON()
			return
		}

		ctx := context.WithValue(req.Context(), "user", *ua)
		next.ServeHTTP(w, req.WithContext(ctx))
	}

	return f
}
