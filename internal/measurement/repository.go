package measurement

import (
	"context"

	"github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const columns = "id, job_id, latitude, longitude, reference_method, displacement_mm, dam_coherence_mean, dam_coherence_valid_pixels, created_at"
const columnsWithJoin = "m.id, m.job_id, m.latitude, m.longitude, m.reference_method, m.displacement_mm, m.dam_coherence_mean, m.dam_coherence_valid_pixels, m.created_at"

type Repository struct {
	pool *pgxpool.Pool
}

func scanMeasurement(row pgx.Row) (model.Measurement, error) {
	var m model.Measurement
	err := row.Scan(
		&m.ID, &m.JobID, &m.Latitude, &m.Longitude,
		&m.ReferenceMethod, &m.DisplacementMM, &m.DamCoherenceMean,
		&m.DamCoherenceValidPixels, &m.CreatedAt,
	)
	return m, err
}

func scanObjectMeasurement(row pgx.Row) (model.ObjectMeasurement, error) {
	var m model.ObjectMeasurement
	err := row.Scan(
		&m.ID, &m.JobID, &m.Latitude, &m.Longitude,
		&m.ReferenceMethod, &m.DisplacementMM, &m.DamCoherenceMean,
		&m.DamCoherenceValidPixels, &m.CreatedAt,
		&m.ReferenceDate, &m.SecondaryDate,
	)
	return m, err
}

func (r *Repository) listMeasurement(ctx context.Context, query string, args ...any) ([]model.Measurement, error) {
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

func (r *Repository) listObjectMeasurement(ctx context.Context, query string, args ...any) ([]model.ObjectMeasurement, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listObjectMeasurements := []model.ObjectMeasurement{}
	for rows.Next() {
		objectMeasurement, err := scanObjectMeasurement(rows)
		if err != nil {
			return nil, err
		}
		listObjectMeasurements = append(listObjectMeasurements, objectMeasurement)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return listObjectMeasurements, nil
}

func (r *Repository) ListMeasurements(ctx context.Context) ([]model.Measurement, error) {
	return r.listMeasurement(ctx, "SELECT "+columns+" FROM measurements ORDER BY id")
}

func (r *Repository) ListMeasurementsByJob(ctx context.Context, jobID int) ([]model.Measurement, error) {
	return r.listMeasurement(ctx, "SELECT "+columns+" FROM measurements WHERE job_id = $1 ORDER BY id", jobID)
}

func (r *Repository) GetMeasurementByID(ctx context.Context, id int64) (model.Measurement, error) {
	query := "SELECT " + columns + " FROM measurements WHERE id = $1"
	return scanMeasurement(r.pool.QueryRow(ctx, query, id))
}

func (r *Repository) CreateMeasurement(ctx context.Context, measurement model.Measurement) (model.Measurement, error) {
	query := "INSERT INTO measurements (job_id, latitude, longitude, reference_method, displacement_mm, dam_coherence_mean, dam_coherence_valid_pixels) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING " + columns
	return scanMeasurement(r.pool.QueryRow(ctx, query,
		measurement.JobID, measurement.Latitude, measurement.Longitude,
		measurement.ReferenceMethod, measurement.DisplacementMM,
		measurement.DamCoherenceMean, measurement.DamCoherenceValidPixels,
	))
}

func (r *Repository) UpdateMeasurement(ctx context.Context, id int64, measurement model.Measurement) (model.Measurement, error) {
	query := "UPDATE measurements SET job_id = $1, latitude = $2, longitude = $3, reference_method = $4, displacement_mm = $5, dam_coherence_mean = $6, dam_coherence_valid_pixels = $7 WHERE id = $8 RETURNING " + columns
	return scanMeasurement(r.pool.QueryRow(ctx, query,
		measurement.JobID, measurement.Latitude, measurement.Longitude,
		measurement.ReferenceMethod, measurement.DisplacementMM,
		measurement.DamCoherenceMean, measurement.DamCoherenceValidPixels, id,
	))
}

func (r *Repository) DeleteMeasurement(ctx context.Context, id int64) (model.Measurement, error) {
	query := "DELETE FROM measurements WHERE id = $1 RETURNING " + columns
	return scanMeasurement(r.pool.QueryRow(ctx, query, id))
}

func (r *Repository) ListMeasurementsByObject(ctx context.Context, id int64) ([]model.ObjectMeasurement, error) {
	query := "SELECT " + columnsWithJoin + ", j.reference_date, j.secondary_date FROM measurements AS m JOIN insar_jobs AS j ON j.id = m.job_id WHERE j.object_id = $1 ORDER BY j.secondary_date, j.reference_date, m.reference_method, m.id"
	return r.listObjectMeasurement(ctx, query, id)
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
