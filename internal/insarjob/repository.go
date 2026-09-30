package insarjob

import (
	"context"

	"github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const columns = "id, object_id, hyp3_job_id, job_type, reference_granule, secondary_granule, reference_date, secondary_date, status, created_at, finished_at"

type Repository struct {
	pool *pgxpool.Pool
}

func scanInsarJob(row pgx.Row) (model.InsarJob, error) {
	var j model.InsarJob
	err := row.Scan(
		&j.ID, &j.ObjectID, &j.Hyp3JobID, &j.JobType,
		&j.ReferenceGranule, &j.SecondaryGranule,
		&j.ReferenceDate, &j.SecondaryDate,
		&j.Status, &j.CreatedAt, &j.FinishedAt,
	)
	return j, err
}

func (r *Repository) list(ctx context.Context, query string, args ...any) ([]model.InsarJob, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listInsarJobs := []model.InsarJob{}
	for rows.Next() {
		insarJob, err := scanInsarJob(rows)
		if err != nil {
			return nil, err
		}
		listInsarJobs = append(listInsarJobs, insarJob)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return listInsarJobs, nil
}

func (r *Repository) ListInsarJobs(ctx context.Context) ([]model.InsarJob, error) {
	return r.list(ctx, "SELECT "+columns+" FROM insar_jobs ORDER BY id")
}

func (r *Repository) ListInsarJobsByObject(ctx context.Context, objectID int) ([]model.InsarJob, error) {
	return r.list(ctx, "SELECT "+columns+" FROM insar_jobs WHERE object_id = $1 ORDER BY reference_date", objectID)
}

func (r *Repository) GetInsarJobByID(ctx context.Context, id int) (model.InsarJob, error) {
	query := "SELECT " + columns + " FROM insar_jobs WHERE id = $1"
	return scanInsarJob(r.pool.QueryRow(ctx, query, id))
}

func (r *Repository) CreateInsarJob(ctx context.Context, insarJob model.InsarJob) (model.InsarJob, error) {
	if insarJob.JobType == "" {
		insarJob.JobType = "INSAR_GAMMA"
	}
	if insarJob.Status == "" {
		insarJob.Status = "PENDING"
	}
	query := "INSERT INTO insar_jobs (object_id, hyp3_job_id, job_type, reference_granule, secondary_granule, reference_date, secondary_date, status, finished_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING " + columns
	return scanInsarJob(r.pool.QueryRow(ctx, query,
		insarJob.ObjectID, insarJob.Hyp3JobID, insarJob.JobType,
		insarJob.ReferenceGranule, insarJob.SecondaryGranule,
		insarJob.ReferenceDate, insarJob.SecondaryDate,
		insarJob.Status, insarJob.FinishedAt,
	))
}

func (r *Repository) UpdateInsarJob(ctx context.Context, id int, insarJob model.InsarJob) (model.InsarJob, error) {
	query := "UPDATE insar_jobs SET object_id = $1, hyp3_job_id = $2, job_type = $3, reference_granule = $4, secondary_granule = $5, reference_date = $6, secondary_date = $7, status = $8, finished_at = $9 WHERE id = $10 RETURNING " + columns
	return scanInsarJob(r.pool.QueryRow(ctx, query,
		insarJob.ObjectID, insarJob.Hyp3JobID, insarJob.JobType,
		insarJob.ReferenceGranule, insarJob.SecondaryGranule,
		insarJob.ReferenceDate, insarJob.SecondaryDate,
		insarJob.Status, insarJob.FinishedAt, id,
	))
}

func (r *Repository) DeleteInsarJob(ctx context.Context, id int) (model.InsarJob, error) {
	query := "DELETE FROM insar_jobs WHERE id = $1 RETURNING " + columns
	return scanInsarJob(r.pool.QueryRow(ctx, query, id))
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
