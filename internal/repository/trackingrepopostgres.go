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

	const sqlQuery = `select
						id, 
						name,
						created_at,
						user_id,
						alerting_enabled,
						coalesce(lat, 0) "lat",
						coalesce(lng, 0) "lng",
						coalesce(last_alert_timestamp, to_timestamp(0)) "last_alert_timestamp" 
					  from device where id=$1`
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

	const sqlQuery = `insert into 
    					device(name, user_id, created_at, alerting_enabled, lat, lng)
						values ($1, $2, $3, $4, $5, $6) returning id`
	err := t.db.QueryRow(sqlQuery, d.Name, d.UserID, d.CreatedAt, d.AlertingEnabled, d.Lat, d.Lng).Scan(&d.ID)
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (t TrackingRepositoryPostgres) GetDevices(userID int64) ([]models.Device, error) {
	trackInputs := make([]models.Device, 0)

	const sqlQuery = `select
						id, 
						name,
						created_at,
						user_id,
						alerting_enabled,
						coalesce(lat, 0) "lat",
						coalesce(lng, 0) "lng",
						coalesce(last_alert_timestamp, to_timestamp(0)) "last_alert_timestamp"
					  from device where user_id=$1`
	if err := t.db.Select(&trackInputs, sqlQuery, userID); err != nil {
		return nil, err
	}

	return trackInputs, nil
}

// UpdateDevice updates an existing device
func (t TrackingRepositoryPostgres) UpdateDevice(deviceID int64, device models.Device) error {
	const sqlQuery = `update device set
						name=$1,
                  		alerting_enabled=$2,
                  		lat=$3,
                  		lng=$4,
                  		last_alert_timestamp=$5
					  where id=$6`

	res, err := t.db.Exec(sqlQuery,
		device.Name, device.AlertingEnabled, device.Lat, device.Lng, device.LastAlertTimestamp, deviceID)
	if err != nil {
		return err
	}

	raf, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if raf == 0 {
		return ErrDeviceDoesNotExist
	}

	return nil
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
