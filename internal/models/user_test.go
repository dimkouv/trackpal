// +build unit_test

package models

import (
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestUserAccount_IsEmpty(t *testing.T) {
	assert.Equal(t, true, UserAccount{}.IsEmpty())
	assert.Equal(t, false, UserAccount{ID: 123}.IsEmpty())
}

func TestUserAccount_Validate(t *testing.T) {
	assert.Error(t, UserAccount{}.Validate())
	assert.NoError(t, UserAccount{ID: 123}.Validate())
}

func TestNewUserAccount(t *testing.T) {
	assert.Equal(t, &UserAccount{}, NewUserAccount())
}

func TestUserAccount_GetJWT(t *testing.T) {
	ua := UserAccount{
		ID:              123,
		Email:           "ken@plan9.com",
		Passhash:        xid.New().String(),
		FirstName:       "Ken",
		LastName:        "Thompson",
		IsActive:        true,
		ActivationToken: xid.New().String(),
	}

	jwtToken, err := ua.GetJWT(time.Now().UTC().Add(time.Minute))
	assert.NoError(t, err)
	assert.Greater(t, len(jwtToken), 128)
}

func TestUserAccount_FromJWT(t *testing.T) {
	t.Run("fetch user from valid jwt", func(t *testing.T) {
		ua := UserAccount{
			ID:              123,
			Email:           "ken@plan9.com",
			Passhash:        xid.New().String(),
			FirstName:       "Ken",
			LastName:        "Thompson",
			IsActive:        true,
			ActivationToken: xid.New().String(),
		}

		jwtToken, err := ua.GetJWT(time.Now().UTC().Add(2 * time.Minute))
		assert.NoError(t, err)
		ua2, err := NewUserAccount().FromJWT(jwtToken)
		assert.NoError(t, err)
		assert.Equal(t, ua.ID, ua2.ID)
		assert.Equal(t, ua.Email, ua2.Email)
		assert.Equal(t, ua.FirstName, ua2.FirstName)
		assert.Equal(t, ua.LastName, ua2.LastName)
	})

	t.Run("fail if jwt not valid", func(t *testing.T) {
		_, err := NewUserAccount().FromJWT("123")
		assert.Error(t, err)
	})

	t.Run("fail if jwt expired", func(t *testing.T) {
		ua := UserAccount{
			ID:              123,
			Email:           "ken@plan9.com",
			Passhash:        xid.New().String(),
			FirstName:       "Ken",
			LastName:        "Thompson",
			IsActive:        true,
			ActivationToken: xid.New().String(),
		}

		jwtToken, err := ua.GetJWT(time.Now().UTC())
		assert.NoError(t, err)
		_, err = NewUserAccount().FromJWT(jwtToken)
		assert.Equal(t, ErrJWTTokenExpired, err)
	})
}
