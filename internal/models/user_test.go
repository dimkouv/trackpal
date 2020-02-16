// +build unit_test

package models

import (
	"testing"

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
