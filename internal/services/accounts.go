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
	"github.com/dimkouv/trackpal/pkg/terror"
)

type UserAccountService struct {
	repo repository.UserAccountRepository
}

func (s *UserAccountService) CreateUserAccount(ctx context.Context, rc io.Reader) error {
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
		return terror.New(ErrBodyRead, err.Error())
	}

	uaReq := createUserAccountInput{}
	if err = json.Unmarshal(requestData, &uaReq); err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")

		return terror.New(ErrBodyParse, err.Error())
	}

	ua := models.UserAccount{
		Email:     uaReq.Email,
		FirstName: uaReq.FirstName,
		LastName:  uaReq.LastName,
	}
	if err = ua.Validate(); err != nil {
		return err
	}

	createdUA, err := s.repo.SaveNewUser(ua, uaReq.Password)
	if err != nil {
		logrus.WithField(consts.LogFieldErr, err).Errorf("unable to save user")
		return terror.New(ErrPlain, err.Error())
	}

	logrus.Infof("user account created: activationToken=%v", createdUA.ActivationToken)
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

		return nil, terror.New(ErrBodyRead, err.Error())
	}

	uaReq := emailPasswordInput{}
	err = json.Unmarshal(requestData, &uaReq)
	if err != nil {
		logrus.
			WithField(consts.LogFieldBody, fmt.Sprintf("%s", requestData)).
			WithField(consts.LogFieldErr, err).
			Errorf("unable to parse request body")

		return nil, terror.New(ErrBodyParse, err.Error())
	}

	ua, err := s.repo.GetUserByEmailAndPassword(uaReq.Email, uaReq.Password)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to fetch user by email and password")

		return nil, terror.New(ErrBodyParse, err.Error())
	}

	tokenString, err := ua.GetJWT()
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get jwt")

		return nil, terror.New(ErrBodyParse, err.Error())
	}

	return []byte(tokenString), nil
}

func (s *UserAccountService) RefreshJWT(ctx context.Context) ([]byte, error) {
	ua := ctx.Value("user").(models.UserAccount)

	tokenString, err := ua.GetJWT()
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Errorf("unable to get jwt")

		return nil, terror.New(ErrBodyParse, err.Error())
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
