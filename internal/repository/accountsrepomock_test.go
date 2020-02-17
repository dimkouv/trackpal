// +build unit_test

package repository

import (
	"testing"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/cryptoutils"

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

func TestAccountsRepoMock_UpdateUser(t *testing.T) {
	repo := NewAccountsRepoMock()
	ua := models.UserAccount{}
	password := "my-password"
	populatedUA, err := repo.SaveNewUser(ua, password)
	assert.NoError(t, err)

	updatedUA := models.UserAccount{
		Email:           "me@me.com",
		FirstName:       "ken",
		LastName:        "thompson",
		IsActive:        true,
		ActivationToken: cryptoutils.RandomString(6, []rune{'0', '1', '2'}),
	}
	password = "new-password"

	updated, err := repo.UpdateUser(populatedUA.ID, &UpdateUserInput{
		Email:           &updatedUA.Email,
		Password:        &password,
		FirstName:       &updatedUA.FirstName,
		LastName:        &updatedUA.LastName,
		IsActive:        &updatedUA.IsActive,
		ActivationToken: &updatedUA.ActivationToken,
	})

	assert.True(t, updated)
	assert.NoError(t, err)

	fetchedUA, err := repo.GetUserByEmailAndPassword(updatedUA.Email, password)
	assert.NoError(t, err)
	assert.Equal(t, updatedUA.IsActive, fetchedUA.IsActive)
	assert.Equal(t, updatedUA.Email, fetchedUA.Email)
	assert.Equal(t, updatedUA.FirstName, fetchedUA.FirstName)
	assert.Equal(t, updatedUA.LastName, fetchedUA.LastName)
	assert.Equal(t, updatedUA.ActivationToken, fetchedUA.ActivationToken)
}
