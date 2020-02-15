package models

import (
	"errors"
	"fmt"
)

const minDeviceNameLength = 4

type Device struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func (d Device) IsEmpty() bool {
	return d == Device{}
}

func (d Device) Validate() error {
	if d.IsEmpty() {
		return errors.New("the device is empty")
	}

	if len(d.Name) < minDeviceNameLength {
		return fmt.Errorf("the device name should contain at least %d characters", minDeviceNameLength)
	}

	return nil
}
