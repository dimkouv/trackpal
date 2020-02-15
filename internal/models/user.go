package models

import "errors"

type UserAccount struct {
	ID              int64  `json:"id" db:"id"`
	Email           string `json:"email" db:"email"`
	Passhash        string `json:"-" db:"passhash"`
	FirstName       string `json:"firstName" db:"first_name"`
	LastName        string `json:"lastName" db:"last_name"`
	IsActive        bool   `json:"-" db:"is_active"`
	ActivationToken string `json:"-" db:"activation_token"`
}

func (ua UserAccount) IsEmpty() bool {
	return ua == UserAccount{}
}

func (ua UserAccount) Validate() error {
	if ua.IsEmpty() {
		return errors.New("the user account is empty")
	}

	return nil
}
