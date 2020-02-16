package services

import (
	"github.com/dimkouv/trackpal/internal/repository"
)

type UserAccountService struct {
	repo repository.UserAccountRepository
}

// NewUserAccountService receives a repository and returns a user account service
func NewUserAccountService(repo repository.UserAccountRepository) UserAccountService {
	return UserAccountService{repo: repo}
}

// NewUserAccountServicePostgres returns a user account service with a postgres repository
func NewUserAccountServicePostgres(postgresDSN string) (*UserAccountService, error) {
	repo, err := repository.NewAccountsRepositoryPostgres(postgresDSN)
	if err != nil {
		return nil, err
	}

	return &UserAccountService{repo: repo}, nil
}
