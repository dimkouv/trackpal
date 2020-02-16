package repository

import (
	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

const repoName = "accounts_repository_postgres"

type AccountsRepositoryPostgres struct {
	db *sqlx.DB
}

func (a AccountsRepositoryPostgres) SaveNewUser(ua models.UserAccount, password string) (*models.UserAccount, error) {
	panic("implement me")
}

func (a AccountsRepositoryPostgres) GetUserByEmailAndPassword(email, password string) (*models.UserAccount, error) {
	panic("implement me")
}

func NewAccountsRepositoryPostgres(postgresDSN string) (*AccountsRepositoryPostgres, error) {
	db, err := sqlx.Connect("postgres", postgresDSN)

	logrus.
		WithField(consts.LogFieldRepo, repoName).
		Debugf("attempting postgres connection with dsn=%s", postgresDSN)

	if err != nil {
		logrus.
			WithField(consts.LogFieldRepo, repoName).
			WithField(consts.LogFieldErr, err).
			Error("unable to connect to postgres")
		return nil, err
	}

	logrus.
		WithField(consts.LogFieldRepo, repoName).
		Info("postgres connection success")

	return &AccountsRepositoryPostgres{db: db}, nil
}
