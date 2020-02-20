package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/dimkouv/trackpal/internal/repository"
)

type UserAccountService struct {
	repo repository.UserAccountRepository
}

func (s *UserAccountService) CreateUserAccount(_ context.Context, rc io.Reader) error {
	type createUserAccountInput struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return consts.ErrEnumInvalidBody
	}

	uaReq := createUserAccountInput{}
	if err = json.Unmarshal(requestData, &uaReq); err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return consts.ErrEnumInvalidBody
	}

	ua := models.UserAccount{
		Email:     uaReq.Email,
		FirstName: uaReq.FirstName,
		LastName:  uaReq.LastName,
	}

	if !consts.RgxEmail.MatchString(uaReq.Email) {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")
		return consts.ErrEnumInvalidEmail
	}

	if len(uaReq.Password) < 6 {
		logrus.Errorf("the provided password '%s' is not secure", uaReq.Password)
		return consts.ErrEnumInsecurePassword
	}

	createdUA, err := s.repo.SaveNewUser(ua, uaReq.Password)
	switch {
	case err == repository.ErrAccountExists:
		logrus.Errorf("target account already exists '%s'", uaReq.Email)
		return consts.ErrEnumAccountExists
	case err != nil:
		logrus.WithField(consts.LogFieldErr, err).Errorf("unable to save user")
		return err
	}

	// TODO: send it with email, sms, etc...
	logrus.Infof("user account created: activationToken=%v", createdUA.ActivationToken)
	return nil
}

func (s *UserAccountService) ActivateUserAccount(ctx context.Context, rc io.Reader) error {
	type activateUserAccountInput struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")
		return consts.ErrEnumInvalidBody
	}

	activationReq := activateUserAccountInput{}
	if err = json.Unmarshal(requestData, &activationReq); err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")

		return consts.ErrEnumInvalidBody
	}

	if err := s.repo.ActivateUserAccount(activationReq.Email, activationReq.Token); err != nil {
		logrus.WithField(consts.LogFieldErr, err).Errorf("account activation failed")
		return err
	}

	logrus.Debug("user account activated, token invalidated")
	return nil
}

func (s *UserAccountService) GetJWTFromEmailAndPassword(ctx context.Context, rc io.Reader) ([]byte, error) {
	type emailPasswordInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	requestData, err := ioutil.ReadAll(rc)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to read request body")

		return nil, consts.ErrEnumInvalidBody
	}

	uaReq := emailPasswordInput{}
	err = json.Unmarshal(requestData, &uaReq)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")

		return nil, consts.ErrEnumInvalidBody
	}

	ua, err := s.repo.GetUserByEmailAndPassword(uaReq.Email, uaReq.Password)
	switch {
	case err == repository.ErrUserAccountNotFound:
		logrus.Errorf("target account not found '%s'", uaReq.Email)
		return nil, consts.ErrEnumNotFound
	case err != nil:
		logrus.WithField(consts.LogFieldErr, err).Errorf("unable to save user")
		return nil, err
	case !ua.IsActive:
		logrus.Errorf("target account not active '%s'", uaReq.Email)
		return nil, consts.ErrEnumNotActivated
	}

	tokenString, err := ua.GetJWT()
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get jwt")

		return nil, consts.ErrInternal
	}

	return []byte(tokenString), nil
}

func (s *UserAccountService) RefreshJWT(ctx context.Context) ([]byte, error) {
	ua, exists := ctx.Value("user").(models.UserAccount)
	if !exists {
		return nil, consts.ErrEnumUnauthorized
	}

	tokenString, err := ua.GetJWT()
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get jwt")
		return nil, err
	}

	return []byte(tokenString), nil
}

// NewUserAccountService receives a repository and returns a user account service
func NewUserAccountService(repo repository.UserAccountRepository) *UserAccountService {
	return &UserAccountService{repo: repo}
}

// NewUserAccountServicePostgres returns a user account service with a postgres repository
func NewUserAccountServicePostgres(postgresDSN string) (*UserAccountService, error) {
	repo, err := repository.NewAccountsRepositoryPostgres(postgresDSN)
	if err != nil {
		return nil, err
	}

	return &UserAccountService{repo: repo}, nil
}
