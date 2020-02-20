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

func (repo *AccountsRepoMock) ActivateUserAccount(email, token string) error {
	xidToken, err := xid.FromString(token)
	if err != nil {
		return err
	}
	if xidToken.Time().Before(time.Now().UTC().Add(-10 * time.Minute)) {
		return errors.New("your token has expired")
	}

	for i := range repo.userAccount {
		if repo.userAccount[i].Email != email {
			continue
		}

		actualToken, err := xid.FromString(repo.userAccount[i].ActivationToken)
		if err != nil {
			return errors.New("internal error: the token of the user is not valid")
		}

		if actualToken.Compare(xidToken) == 0 {
			repo.userAccount[i].IsActive = true
			return nil
		}
	}

	return errors.New("account with the provided email address not found")
}

func (repo *AccountsRepoMock) UpdateUser(userID int64, input *UpdateUserInput) (bool, error) {
	if input == nil {
		return false, nil
	}

	for i := range repo.userAccount {
		if repo.userAccount[i].ID != userID {
			continue
		}

		switch {
		case input.ActivationToken != nil:
			repo.userAccount[i].ActivationToken = *input.ActivationToken
			fallthrough
		case input.IsActive != nil:
			repo.userAccount[i].IsActive = *input.IsActive
			fallthrough
		case input.LastName != nil:
			repo.userAccount[i].LastName = *input.LastName
			fallthrough
		case input.FirstName != nil:
			repo.userAccount[i].FirstName = *input.FirstName
			fallthrough
		case input.Email != nil:
			repo.userAccount[i].Email = *input.Email
			fallthrough
		case input.Password != nil:
			passhash, err := cryptoutils.Argon2Hash(*input.Password)
			if err != nil {
				return false, fmt.Errorf("unable to generate passhash: %v", err)
			}
			repo.userAccount[i].Passhash = passhash
		}

		return true, nil
	}

	return false, errors.New("user account not found")
}

func (repo *AccountsRepoMock) SaveNewUser(ua models.UserAccount, password string) (*models.UserAccount, error) {
	passhash, err := cryptoutils.Argon2Hash(password)
	if err != nil {
		return nil, fmt.Errorf("unable to generate passhash: %v", err)
	}
	ua.ID = int64(len(repo.userAccount) + 1)
	ua.Passhash = passhash
	ua.ActivationToken = xid.NewWithTime(time.Now().UTC()).String()
	ua.IsActive = false

	for _, u := range repo.userAccount {
		if u.Email == ua.Email {
			return nil, ErrAccountExists
		}
	}

	repo.userAccount = append(repo.userAccount, ua)
	return &ua, nil
}

func (repo *AccountsRepoMock) GetUserByEmailAndPassword(email, password string) (*models.UserAccount, error) {
	for _, ua := range repo.userAccount {
		if ua.Email == email {
			if !ua.IsActive {
				return nil, errors.New("account is not active")
			}
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
