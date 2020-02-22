package server

import (
	"context"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/response"
)

func (ts *TrackpalServer) withUser(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		ua, err := models.NewUserAccount().FromJWT(req.Header.Get(consts.Authorization))

		switch {
		case err == nil:
			ctx := context.WithValue(req.Context(), consts.CtxUser, *ua)
			next.ServeHTTP(w, req.WithContext(ctx))
		case err == models.ErrJWTTokenExpired:
			response.HTTP(w).Status(http.StatusUnauthorized).JSON()
		default:
			logrus.WithField(consts.LogFieldErr, err).Error("unable to load user from jwt")
			response.HTTP(w).Status(http.StatusUnauthorized).JSON()
		}
	}

	return f
}
