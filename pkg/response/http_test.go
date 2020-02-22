// +build unit_test

package response

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {
	w := &httptest.ResponseRecorder{
		Code:    0,
		Body:    nil,
		Flushed: false,
	}

	resp := HTTP(w)
	assert.Equal(t, w, resp.w)
	assert.Empty(t, resp.data)
}

func TestResponse_Data(t *testing.T) {
	w := &httptest.ResponseRecorder{}
	b := []byte("hello")
	resp := HTTP(w).Data(b)
	assert.Equal(t, resp.data, b)
}

func TestResponse_Status(t *testing.T) {
	w := &httptest.ResponseRecorder{}
	resp := HTTP(w).Status(http.StatusBadRequest)
	assert.Equal(t, resp.status, http.StatusBadRequest)
}

func TestResponse_JSON(t *testing.T) {
	w := &httptest.ResponseRecorder{}
	resp := HTTP(w)
	resp.JSON()
	assert.Equal(t, resp.w.Header().Get("content-type"), "application/json")
}

func TestResponse_TEXT(t *testing.T) {
	w := &httptest.ResponseRecorder{}
	resp := HTTP(w)
	resp.TEXT()
	assert.Equal(t, resp.w.Header().Get("content-type"), "text/plain")
}

func TestResponse_Error(t *testing.T) {
	w := &httptest.ResponseRecorder{}
	err := errors.New("just an error")
	resp := HTTP(w).Data([]byte("hello")).Error(err).Data([]byte("hello"))
	resp.JSON()
	assert.Equal(t, `{"error":"`+err.Error()+`"}`, string(resp.data))
}
