package repository

import (
	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type TrackingRepositoryPostgres struct {
	db *sqlx.DB
}

func (t2 TrackingRepositoryPostgres) SaveNewTrackInput(t models.TrackInput) (*models.TrackInput, error) {
	panic("implement me")
}

func (t2 TrackingRepositoryPostgres) GetAllTrackInputsOfDevice(deviceID int64) ([]models.TrackInput, error) {
	panic("implement me")
}

func (t2 TrackingRepositoryPostgres) SaveNewDevice(d models.Device) (*models.Device, error) {
	panic("implement me")
}

func (t2 TrackingRepositoryPostgres) GetDevices() ([]models.Device, error) {
	panic("implement me")
}

func NewTrackingRepositoryPostgres(postgresDSN string) (*TrackingRepositoryPostgres, error) {
	db, err := sqlx.Connect("postgres", postgresDSN)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Error("unable to connect to postgres")
		return nil, err
	}
	logrus.Info("postgres connection success")

	return &TrackingRepositoryPostgres{db: db}, nil
}
