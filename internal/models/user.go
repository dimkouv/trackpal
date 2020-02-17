package models

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

func (ua *UserAccount) FromJWT(jwt string) (*UserAccount, error) {
	return ua, errors.New("unimplemented")
}

func (ua *UserAccount) GetJWT(jwtSignBytes []byte) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      ua.Email,
		"first_name": ua.FirstName,
		"last_name":  ua.LastName,
		"createdAt":  time.Now(),
	})

	tokStr, err := tok.SignedString(jwtSignBytes)
	if err != nil {
		return "", err
	}

	return tokStr, nil
}

func NewUserAccount() *UserAccount {
	return &UserAccount{}
}
