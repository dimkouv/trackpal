package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"

	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/cryptoutils"
)

type AccountsRepoMock struct {
	userAccount []models.UserAccount
}

func (repo *AccountsRepoMock) SaveNewUser(ua models.UserAccount, password string) (*models.UserAccount, error) {
	passhash, err := cryptoutils.Argon2Hash(password)
	if err != nil {
		return nil, fmt.Errorf("unable to generate passhash: %v", err)
	}
	ua.ID = int64(len(repo.userAccount) + 1)
	ua.Passhash = passhash
	ua.ActivationToken = xid.NewWithTime(time.Now()).String()
	ua.IsActive = false

	repo.userAccount = append(repo.userAccount, ua)
	return &ua, nil
}

func (repo *AccountsRepoMock) GetUserByEmailAndPassword(email, password string) (*models.UserAccount, error) {
	for _, ua := range repo.userAccount {
		if ua.Email == email {
			if err := cryptoutils.Argon2Verify(password, ua.Passhash); err != nil {
				break
			}
			return &ua, nil
		}
	}

	return nil, errors.New("user account not found")
}

func NewAccountsRepoMock() *AccountsRepoMock {
	return &AccountsRepoMock{
		userAccount: make([]models.UserAccount, 0),
	}
}
