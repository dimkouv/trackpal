package repository

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/dimkouv/trackpal/internal/consts"
	"github.com/dimkouv/trackpal/internal/models"
)

type TrackingRepositoryPostgres struct {
	db *sqlx.DB
}

func (t TrackingRepositoryPostgres) GetDeviceByID(deviceID int64) (*models.Device, error) {
	var device models.Device

	const sqlQuery = `select id, name, created_at, user_id from device where id=$1`
	err := t.db.Get(&device, sqlQuery, deviceID)
	if err != nil {
		pqErr, isPqErr := err.(*pq.Error)
		if (isPqErr && pqErr.Code == consts.PQCodeForeignKeyViolation) || err == sql.ErrNoRows {
			return nil, ErrDeviceDoesNotExist
		}
		return nil, err
	}

	return &device, nil
}

func (t TrackingRepositoryPostgres) SaveNewTrackInput(trackInput models.TrackInput) (*models.TrackInput, error) {
	if err := trackInput.Validate(); err != nil {
		return nil, err
	}

	const sqlQuery = `insert into track_input(lat, lng, recorded_at, device_id) values ($1, $2, $3, $4) returning id`
	if err := t.db.QueryRow(
		sqlQuery,
		trackInput.Location.Lat,
		trackInput.Location.Lng,
		trackInput.RecordedAt.UTC(),
		trackInput.DeviceID,
	).Scan(&trackInput.ID); err != nil {
		pqErr := err.(*pq.Error)
		if pqErr.Code == consts.PQCodeForeignKeyViolation {
			return nil, ErrDeviceDoesNotExist
		}
		return nil, err
	}

	return &trackInput, nil
}

func (t TrackingRepositoryPostgres) GetAllTrackInputsOfDevice(deviceID int64) ([]models.TrackInput, error) {
	trackInputs := make([]models.TrackInput, 0)

	const sqlQueryDeviceExists = `select exists(select 1 from device where id=$1)`
	exists := false
	if err := t.db.QueryRow(sqlQueryDeviceExists, deviceID).Scan(&exists); err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrDeviceDoesNotExist
	}

	const sqlQuery = `select id, lat, lng, recorded_at, created_at from track_input where device_id=$1`
	if err := t.db.Select(&trackInputs, sqlQuery, deviceID); err != nil {
		return nil, err
	}

	return trackInputs, nil
}

func (t TrackingRepositoryPostgres) SaveNewDevice(d models.Device) (*models.Device, error) {
	if err := d.Validate(); err != nil {
		return nil, err
	}

	d.CreatedAt = time.Now().UTC().Truncate(time.Second)

	const sqlQuery = `insert into device(name, user_id, created_at) values ($1, $2, $3) returning id`
	err := t.db.QueryRow(sqlQuery, d.Name, d.UserID, d.CreatedAt).Scan(&d.ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (t TrackingRepositoryPostgres) GetDevices(userID int64) ([]models.Device, error) {
	trackInputs := make([]models.Device, 0)

	const sqlQuery = `select id, name, created_at from device where user_id=$1`
	if err := t.db.Select(&trackInputs, sqlQuery, userID); err != nil {
		return nil, err
	}

	return trackInputs, nil
}

func NewTrackingRepositoryPostgres(postgresDSN string) (*TrackingRepositoryPostgres, error) {
	db, err := sqlx.Connect("postgres", postgresDSN)

	logrus.Debugf("attempting postgres connection with dsn=%s", postgresDSN)
	if err != nil {
		logrus.
			WithField(consts.LogFieldErr, err).
			Error("unable to connect to postgres")
		return nil, err
	}
	logrus.
		WithField(consts.LogFieldRepo, "tracking_repository_postgres").
		Info("postgres connection success")

	return &TrackingRepositoryPostgres{db: db}, nil
}
