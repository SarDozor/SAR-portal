package measurement

import (
	"context"
	"time"

	objects "github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryMeasurement struct {
	pool *pgxpool.Pool
}

func (r *RepositoryMeasurement) ListMeasurement(ctx context.Context) ([]objects.Measurement, error) {
	var listMeasurement []objects.Measurement
	query := "SELECT * FROM measurements"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var measurement objects.Measurement
		if err = rows.Scan(
			&measurement.ID, &measurement.JobID,
			&measurement.Latitude, &measurement.Longitude,
			&measurement.DisplacementMM, &measurement.Coherence, &measurement.CreatedAt,
		); err != nil {
			return nil, err
		}
		listMeasurement = append(listMeasurement, measurement)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return listMeasurement, nil
}

func (r *RepositoryMeasurement) GetMeasurementByID(ctx context.Context, id int) (objects.Measurement, error) {
	var measurement objects.Measurement
	query := "SELECT * FROM measurements WHERE id = $1"
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&measurement.ID, &measurement.JobID,
		&measurement.Latitude, &measurement.Longitude,
		&measurement.DisplacementMM, &measurement.Coherence, &measurement.CreatedAt,
	); err != nil {
		return objects.Measurement{}, err
	}
	return measurement, nil
}

func (r *RepositoryMeasurement) CreateMeasurement(ctx context.Context, measurement objects.Measurement) (objects.Measurement, error) {
	var id int
	var jobID int
	var createdAt time.Time
	query := "INSERT INTO measurements (latitude, longitude, displacement_mm, coherence) VALUES ($1, $2, $3, $4) RETURNING id, job_id, created_at"
	if err := r.pool.QueryRow(ctx, query,
		measurement.Latitude, measurement.Longitude, measurement.DisplacementMM, measurement.Coherence,
	).Scan(
		&id, &jobID, &createdAt,
	); err != nil {
		return objects.Measurement{}, err
	}
	measurement.ID = id
	measurement.JobID = jobID
	measurement.CreatedAt = createdAt
	return measurement, nil
}

func (r *RepositoryMeasurement) UpdateMeasurement(ctx context.Context, id int, measurement objects.Measurement) (objects.Measurement, error) {
	var updatedMeasurement objects.Measurement
	query := "UPDATE measurements SET latitude = $1, longitude = $2, displacement_mm = $3, coherence = $4 WHERE id = $5 RETURNING id, job_id"
	err := r.pool.QueryRow(ctx, query,
		measurement.Latitude, measurement.Longitude, measurement.DisplacementMM, measurement.Coherence, id,
	).Scan(
		&measurement.ID, &measurement.JobID,
		&measurement.Latitude, &measurement.Longitude,
		&measurement.DisplacementMM, &measurement.Coherence,
	)
	if err != nil {
		return objects.Measurement{}, err
	}
	return updatedMeasurement, nil
}

func (r *RepositoryMeasurement) DeleteMeasurement(ctx context.Context, id int) (objects.Measurement, error) {
	var deletedMeasurement objects.Measurement
	query := "DELETE FROM measurements WHERE id = $1 RETURNING *"
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&deletedMeasurement.ID, &deletedMeasurement.JobID, &deletedMeasurement.Latitude,
		&deletedMeasurement.Longitude, &deletedMeasurement.DisplacementMM,
		&deletedMeasurement.Coherence, &deletedMeasurement.CreatedAt,
	)
	if err != nil {
		return objects.Measurement{}, err
	}
	return deletedMeasurement, nil
}

func NewRepositoryMeasurement(pool *pgxpool.Pool) *RepositoryMeasurement {
	return &RepositoryMeasurement{pool: pool}
}
