package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/cryptoutils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

const repoName = "accounts_repository_postgres"

type AccountsRepositoryPostgres struct {
	db *sqlx.DB
}

func (repo AccountsRepositoryPostgres) SaveNewUser(ua models.UserAccount, password string) (*models.UserAccount, error) {
	if err := ua.Validate(); err != nil {
		return nil, err
	}

	passhash, err := cryptoutils.Argon2Hash(password)
	if err != nil {
		return nil, fmt.Errorf("unable to generate passhash: %v", err)
	}
	ua.Passhash = passhash
	ua.ActivationToken = xid.NewWithTime(time.Now()).String()
	ua.IsActive = false

	const sqlQuery = `insert into user_account(email, passhash, first_name, last_name, is_active, activation_token)` +
		` values ($1, $2, $3, $4, $5, $6) returning id`

	err = repo.db.QueryRow(sqlQuery, ua.Email, ua.Passhash, ua.FirstName, ua.LastName, ua.IsActive, ua.ActivationToken).
		Scan(&ua.ID)
	if err != nil {
		return nil, err
	}

	return &ua, nil
}

func (repo AccountsRepositoryPostgres) GetUserByEmailAndPassword(email, password string) (*models.UserAccount, error) {
	ua := models.UserAccount{}
	const sqlQuery = `select id, email, passhash, first_name, last_name, is_active, activation_token from` +
		` user_account where email=$1`
	err := repo.db.Get(&ua, sqlQuery, email)
	if err != nil {
		return nil, err
	}

	err = cryptoutils.Argon2Verify(password, ua.Passhash)
	switch err {
	case nil:
		return &ua, nil
	default:
		return nil, errors.New("user account not found")
	}
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
