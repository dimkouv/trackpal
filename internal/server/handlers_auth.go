package server

import "net/http"

func (ts TrackpalServer) authRegister(w http.ResponseWriter, req *http.Request) {
	err := ts.userService.CreateUserAccount(req.Context(), req.Body)
	ts.responseBasedOnErr(err, nil, http.StatusCreated, w)
}

func (ts TrackpalServer) authActivate(w http.ResponseWriter, req *http.Request) {
	err := ts.userService.ActivateUserAccount(req.Context(), req.Body)
	ts.responseBasedOnErr(err, nil, http.StatusOK, w)
}

func (ts TrackpalServer) authLogin(w http.ResponseWriter, req *http.Request) {
	b, err := ts.userService.GetJWTFromEmailAndPassword(req.Context(), req.Body)
	ts.responseBasedOnErr(err, b, http.StatusOK, w)
}

func (ts TrackpalServer) authRefresh(w http.ResponseWriter, req *http.Request) {
	b, err := ts.userService.RefreshJWT(req.Context())
	ts.responseBasedOnErr(err, b, http.StatusOK, w)
}
