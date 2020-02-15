// +build unit_test

package repository

import (
	"testing"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountsRepoMock(t *testing.T) {
	repo := NewAccountsRepoMock()
	assert.Empty(t, repo.userAccount)
	assert.NotNil(t, repo.userAccount)
}

func TestAccountsRepoMock_SaveNewUser(t *testing.T) {
	repo := NewAccountsRepoMock()

	ua := models.UserAccount{
		Email:     "myuser@email.com",
		FirstName: "my",
		LastName:  "user",
	}
	password := "my-password"

	populatedUA, err := repo.SaveNewUser(ua, password)
	assert.NoError(t, err)
	assert.Equal(t, populatedUA.IsActive, false)
}

func TestAccountsRepoMock_GetUserByEmailAndPassword(t *testing.T) {
	repo := NewAccountsRepoMock()

	ua := models.UserAccount{
		Email:     "myuser@email.com",
		FirstName: "my",
		LastName:  "user",
	}
	password := "my-password"

	_, err := repo.SaveNewUser(ua, password)
	assert.NoError(t, err)

	populatedUA, err := repo.GetUserByEmailAndPassword(ua.Email, password)
	assert.NoError(t, err)
	assert.Equal(t, populatedUA.IsActive, false)

	_, err = repo.GetUserByEmailAndPassword(ua.Email, "invalid-password")
	assert.Error(t, err)
}
