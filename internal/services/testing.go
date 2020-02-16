package services

import "errors"

// IoReaderErrAlways is a io.Reader impl. that always responds with an error
type IoReaderErrAlways struct{}

// Read returns 0 and an error
func (e IoReaderErrAlways) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader error")
}
