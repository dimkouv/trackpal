package models

import "errors"

type Device struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (d Device) IsEmpty() bool {
	return d == Device{}
}

func (d Device) Validate() error {
	if d.IsEmpty() {
		return errors.New("the device is empty")
	}

	if d.Name == "" {
		return errors.New("the device does not have a name")
	}

	return nil
}
