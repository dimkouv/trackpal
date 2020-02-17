// +build integration_test

package repository

import (
	"testing"

	"github.com/rs/xid"

	"github.com/dimkouv/trackpal/internal/envlib"
	"github.com/dimkouv/trackpal/pkg/cryptoutils"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewAccountsRepositoryPostgres(t *testing.T) {
	_, err := NewAccountsRepositoryPostgres("invalid")
	assert.Error(t, err)

	repo, err := NewAccountsRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)
	assert.NotNil(t, repo.db)
}

func TestAccountsRepoPostgres_SaveNewUser(t *testing.T) {
	repo, err := NewAccountsRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

	ua := models.UserAccount{
		Email:     xid.New().String() + "myuser@email.com",
		FirstName: "my",
		LastName:  "user",
	}
	password := "my-password"

	t.Run("saving a new user should succeed", func(t *testing.T) {
		populatedUA, err := repo.SaveNewUser(ua, password)
		assert.NoError(t, err)
		assert.Equal(t, populatedUA.IsActive, false)
		assert.Greater(t, populatedUA.ID, int64(0))
	})

	t.Run("saving a user with an existing email should fail", func(t *testing.T) {
		_, err = repo.SaveNewUser(ua, password)
		assert.Error(t, err)
	})

	t.Run("saving an invalid user should fail", func(t *testing.T) {
		_, err = repo.SaveNewUser(models.UserAccount{}, password)
		assert.Error(t, err)
	})
}

func TestAccountsRepoPostgres_GetUserByEmailAndPassword(t *testing.T) {
	repo, err := NewAccountsRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

	ua := models.UserAccount{
		Email:     xid.New().String() + "@email.com",
		FirstName: "my",
		LastName:  "user",
	}
	password := "my-password"

	_, err = repo.SaveNewUser(ua, password)
	assert.NoError(t, err)

	t.Run("fetching an existing user with email and correct password should pass", func(t *testing.T) {
		populatedUA, err := repo.GetUserByEmailAndPassword(ua.Email, password)
		assert.NoError(t, err)
		assert.Equal(t, populatedUA.IsActive, false)
		assert.Greater(t, populatedUA.ID, int64(0))
	})

	t.Run("should fail if the password is not correct", func(t *testing.T) {
		_, err = repo.GetUserByEmailAndPassword(ua.Email, "invalid-password")
		assert.Error(t, err)
	})

	t.Run("should fail if the username is not correct", func(t *testing.T) {
		_, err = repo.GetUserByEmailAndPassword("a-non-existing@email.com", password)
		assert.Error(t, err)
	})
}

func TestAccountsRepositoryPostgres_UpdateUser(t *testing.T) {
	repo, err := NewAccountsRepositoryPostgres(envlib.GetEnvOrPanic(envlib.EnvPostgresDSN))
	assert.NoError(t, err)

	ua := models.UserAccount{
		Email:     xid.New().String() + "myuser@email.com",
		FirstName: "my",
		LastName:  "user",
	}
	password := "my-password"
	populatedUA, err := repo.SaveNewUser(ua, password)
	assert.NoError(t, err)

	updatedUA := models.UserAccount{
		Email:           xid.New().String() + "@me.com",
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
	assert.NoError(t, err)
	assert.True(t, updated)

	fetchedUA, err := repo.GetUserByEmailAndPassword(updatedUA.Email, password)
	assert.NoError(t, err)
	assert.Equal(t, updatedUA.IsActive, fetchedUA.IsActive)
	assert.Equal(t, updatedUA.Email, fetchedUA.Email)
	assert.Equal(t, updatedUA.FirstName, fetchedUA.FirstName)
	assert.Equal(t, updatedUA.LastName, fetchedUA.LastName)
	assert.Equal(t, updatedUA.ActivationToken, fetchedUA.ActivationToken)
}
