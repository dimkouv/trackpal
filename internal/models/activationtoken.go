package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dimkouv/trackpal/pkg/cryptoutils"
)

// activationToken contains the fields that are stored in user's activationToken field
type activationToken struct {
	pin      string
	issuedAt time.Time
}

// parse receives an activationToken's string representation and parses it into a struct
func (at activationToken) parse(s string) (activationToken, error) {
	parts := strings.Split(s, "\t")
	if len(parts) != 2 {
		return at, errors.New("the provided string is not a valid activation token")
	}

	pin := parts[0]
	if len(pin) == 0 {
		return at, errors.New("invalid pin length")
	}

	dt := new(time.Time)
	if err := dt.UnmarshalText([]byte(parts[1])); err != nil {
		return at, fmt.Errorf("invalid activationToken format. %v", err)
	}

	return activationToken{pin: pin, issuedAt: *dt}, nil
}

// string converts an activation token to it's string representation
func (at activationToken) string() string {
	dt, err := at.issuedAt.MarshalText()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s\t%s", at.pin, dt)
}

// validate receives a pin and checks whether or not it is valid
func (at activationToken) validate(pin string) error {
	if time.Since(at.issuedAt) > 1*time.Hour {
		return errors.New("the activation token has expired")
	}

	if at.pin != pin {
		return errors.New("the pin is not valid")
	}

	return nil
}

// NewActivationToken returns a new activation token with a random pin
func newActivationToken(tokenLength int) activationToken {
	return activationToken{
		pin: cryptoutils.RandomString(
			tokenLength,
			[]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I',
				'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'},
		),
		issuedAt: time.Now(),
	}
}
