package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/pkg/cryptoutils"
)

type AccountsRepositoryPostgres struct {
	db *sqlx.DB
}

func (repo AccountsRepositoryPostgres) ActivateUserAccount(email, token string) error {
	xidToken, err := xid.FromString(token)
	if err != nil {
		return err
	}
	if xidToken.Time().Before(time.Now().UTC().Add(-10 * time.Minute)) {
		return ErrTokenExpired
	}

	q := `update user_account set is_active=true where email=$1 and activation_token=$2`
	res, err := repo.db.Exec(q, email, token)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	} else if n == 0 {
		return ErrTokenNotFound
	}

	return nil
}

func (repo AccountsRepositoryPostgres) UpdateUser(userID int64, input *UpdateUserInput) (bool, error) {
	if input == nil {
		return false, nil
	}

	queryUpdateFields := make([]string, 0)
	args := make([]interface{}, 0)

	if input.ActivationToken != nil {
		args = append(args, *input.ActivationToken)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("activation_token=$%d", len(args)))
	}

	if input.IsActive != nil {
		args = append(args, *input.IsActive)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("is_active=$%d", len(args)))
	}

	if input.LastName != nil {
		args = append(args, *input.LastName)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("last_name=$%d", len(args)))
	}

	if input.FirstName != nil {
		args = append(args, *input.FirstName)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("first_name=$%d", len(args)))
	}

	if input.Email != nil {
		args = append(args, *input.Email)
		queryUpdateFields = append(queryUpdateFields, fmt.Sprintf("email=$%d", len(args)))
	}

	if input.Password != nil {
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
	passhash, err := cryptoutils.Argon2Hash(password)
	if err != nil {
		return nil, fmt.Errorf("unable to generate passhash: %v", err)
	}
	ua.Passhash = passhash
	ua.ActivationToken = xid.NewWithTime(time.Now().UTC()).String()
	ua.IsActive = false

	const sqlQuery = `insert into user_account(email, passhash, first_name, last_name, is_active, activation_token)` +
		` values ($1, $2, $3, $4, $5, $6) returning id`

	err = repo.db.QueryRow(sqlQuery, ua.Email, ua.Passhash, ua.FirstName, ua.LastName, ua.IsActive, ua.ActivationToken).
		Scan(&ua.ID)
	if err != nil {
		pqErr, isPqErr := err.(*pq.Error)
		if isPqErr && pqErr.Code == consts.PQCodeUniqueKeyViolation {
			return nil, ErrAccountExists
		}
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
		if err == sql.ErrNoRows {
			return nil, ErrUserAccountNotFound
		}
		return nil, err
	}

	err = cryptoutils.Argon2Verify(password, ua.Passhash)
	switch err {
	case nil:
		return &ua, nil
	default:
		return nil, ErrUserAccountNotFound
	}
}

func NewAccountsRepositoryPostgres(postgresDSN string) (*AccountsRepositoryPostgres, error) {
	repoName := "accounts_repository_postgres"
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
