package measurement

import (
	"context"

	"github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const columns = "id, job_id, latitude, longitude, displacement_mm, coherence, created_at"

type Repository struct {
	pool *pgxpool.Pool
}

func scanMeasurement(row pgx.Row) (model.Measurement, error) {
	var m model.Measurement
	err := row.Scan(
		&m.ID, &m.JobID, &m.Latitude, &m.Longitude,
		&m.DisplacementMM, &m.Coherence, &m.CreatedAt,
	)
	return m, err
}

func (r *Repository) list(ctx context.Context, query string, args ...any) ([]model.Measurement, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listMeasurements := []model.Measurement{}
	for rows.Next() {
		measurement, err := scanMeasurement(rows)
		if err != nil {
			return nil, err
		}
		listMeasurements = append(listMeasurements, measurement)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return listMeasurements, nil
}

func (r *Repository) ListMeasurements(ctx context.Context) ([]model.Measurement, error) {
	return r.list(ctx, "SELECT "+columns+" FROM measurements ORDER BY id")
}

func (r *Repository) ListMeasurementsByJob(ctx context.Context, jobID int) ([]model.Measurement, error) {
	return r.list(ctx, "SELECT "+columns+" FROM measurements WHERE job_id = $1 ORDER BY id", jobID)
}

func (r *Repository) GetMeasurementByID(ctx context.Context, id int64) (model.Measurement, error) {
	query := "SELECT " + columns + " FROM measurements WHERE id = $1"
	return scanMeasurement(r.pool.QueryRow(ctx, query, id))
}

func (r *Repository) CreateMeasurement(ctx context.Context, measurement model.Measurement) (model.Measurement, error) {
	query := "INSERT INTO measurements (job_id, latitude, longitude, displacement_mm, coherence) VALUES ($1, $2, $3, $4, $5) RETURNING " + columns
	return scanMeasurement(r.pool.QueryRow(ctx, query,
		measurement.JobID, measurement.Latitude, measurement.Longitude,
		measurement.DisplacementMM, measurement.Coherence,
	))
}

func (r *Repository) UpdateMeasurement(ctx context.Context, id int64, measurement model.Measurement) (model.Measurement, error) {
	query := "UPDATE measurements SET job_id = $1, latitude = $2, longitude = $3, displacement_mm = $4, coherence = $5 WHERE id = $6 RETURNING " + columns
	return scanMeasurement(r.pool.QueryRow(ctx, query,
		measurement.JobID, measurement.Latitude, measurement.Longitude,
		measurement.DisplacementMM, measurement.Coherence, id,
	))
}

func (r *Repository) DeleteMeasurement(ctx context.Context, id int64) (model.Measurement, error) {
	query := "DELETE FROM measurements WHERE id = $1 RETURNING " + columns
	return scanMeasurement(r.pool.QueryRow(ctx, query, id))
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
