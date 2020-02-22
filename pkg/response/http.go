package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type errResponseData struct {
	Error string `json:"error"`
}

func (errResp errResponseData) json() []byte {
	b, _ := json.Marshal(errResp)
	return b
}

type Response struct {
	data   []byte
	status int
	w      http.ResponseWriter
	err    error
}

func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}

func (r *Response) Data(data []byte) *Response {
	r.data = data
	return r
}

func (r *Response) Error(err error) *Response {
	r.err = err
	return r
}

func (r *Response) ErrorStr(errStr string) *Response {
	r.err = errors.New(errStr)
	return r
}

func (r *Response) JSON() {
	r.w.Header().Set("Content-Type", "application/json")
	r.compile()
}

func (r *Response) TEXT() {
	r.w.Header().Set("Content-Type", "text/plain")
	r.compile()
}

func (r *Response) compile() {
	r.w.WriteHeader(r.status)

	if r.err != nil {
		r.data = errResponseData{Error: r.err.Error()}.json()
	}

	if _, err := r.w.Write(r.data); err != nil {
		logrus.Errorf("unable to write response data: %v", err)
	}
}

func HTTP(w http.ResponseWriter) *Response {
	return &Response{w: w}
}
