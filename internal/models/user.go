package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/conf"
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

func (ua *UserAccount) FromJWT(tokenString string) (*UserAccount, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		logrus.Debugf("FromJWT parsing tokenString:%s", tokenString)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return conf.JWTSignBytes, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		logrus.WithField("claims", claims).Debug(claims)

		createdAtUTC, err := time.Parse("2006-01-02T15:04:05Z", claims["created_at"].(string))
		if err != nil {
			return nil, err
		}
		if createdAtUTC.Before(time.Now().UTC().Add(-5 * time.Minute)) {
			return nil, errors.New("token has expired")
		}

		return &UserAccount{
			Email:     claims["email"].(string),
			FirstName: claims["first_name"].(string),
			LastName:  claims["last_name"].(string),
		}, nil
	}

	return nil, errors.New("invalid token")

}

func (ua *UserAccount) GetJWT() (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      ua.Email,
		"first_name": ua.FirstName,
		"last_name":  ua.LastName,
		"created_at": time.Now().UTC(),
	})

	tokStr, err := tok.SignedString(conf.JWTSignBytes)
	if err != nil {
		return "", err
	}

	return tokStr, nil
}

func NewUserAccount() *UserAccount {
	return &UserAccount{}
}
