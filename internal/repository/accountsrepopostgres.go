package repository

import (
	"errors"
	"fmt"
	"strings"
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

func (repo AccountsRepositoryPostgres) UpdateUser(userID int64, input *UpdateUserInput) (bool, error) {
	if input == nil {
		return false, nil
	}

	queryUpdateFields := make([]string, 0)
	args := make([]interface{}, 0)
	switch {
	case input.ActivationToken != nil:
		args = append(args, *input.ActivationToken)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("activation_token=$%d", len(args)))
		fallthrough
	case input.IsActive != nil:
		args = append(args, *input.IsActive)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("is_active=$%d", len(args)))
		fallthrough
	case input.LastName != nil:
		args = append(args, *input.LastName)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("last_name=$%d", len(args)))
		fallthrough
	case input.FirstName != nil:
		args = append(args, *input.FirstName)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("first_name=$%d", len(args)))
		fallthrough
	case input.Email != nil:
		args = append(args, *input.Email)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("email=$%d", len(args)))
		fallthrough
	case input.Password != nil:
		passhash, err := cryptoutils.Argon2Hash(*input.Password)
		if err != nil {
			return false, fmt.Errorf("unable to generate passhash: %v", err)
		}
		args = append(args, passhash)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("passhash=$%d", len(args)))
	}

	args = append(args, userID)
	q := fmt.Sprintf("update user_account set %s where id=$%d",
		strings.Join(queryUpdateFields, ", "), len(args))

	res, err := repo.db.Exec(q, args...)
	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if n == 0 {
		return false, nil
	}

	return true, nil
}

func (repo AccountsRepositoryPostgres) SaveNewUser(
	ua models.UserAccount, password string) (*models.UserAccount, error) {
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
