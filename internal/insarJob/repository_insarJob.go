package insarjob

import (
	"context"
	"time"

	objects "github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInsarJob struct {
	pool *pgxpool.Pool
}

func (r *RepositoryInsarJob) ListInsarJob(ctx context.Context) ([]objects.InsarJob, error) {
	query := "SELECT * FROM insar_jobs"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var listInsarJob []objects.InsarJob
	for rows.Next() {
		var insarJob objects.InsarJob
		if err = rows.Scan(
			&insarJob.ID, &insarJob.ObjectID, &insarJob.Hyp3JobID,
			&insarJob.JobType, &insarJob.ReferenceGranule, &insarJob.SecondaryGranule,
			&insarJob.ReferenceDate, &insarJob.SecondaryDate, &insarJob.Status,
			&insarJob.CreatedAt, &insarJob.FinishedAt); err != nil {
			return nil, err
		}
		listInsarJob = append(listInsarJob, insarJob)
	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return listInsarJob, nil
}

func (r *RepositoryInsarJob) GetInsarJobByID(ctx context.Context, id int) (objects.InsarJob, error) {
	var insarJob objects.InsarJob
	query := "SELECT * FROM insar_jobs WHERE id = $1"
	err := r.pool.QueryRow(ctx, query, id).Scan(&insarJob.ID, &insarJob.ObjectID, &insarJob.Hyp3JobID,
		&insarJob.JobType, &insarJob.ReferenceGranule, &insarJob.SecondaryGranule,
		&insarJob.ReferenceDate, &insarJob.SecondaryDate, &insarJob.Status,
		&insarJob.CreatedAt, &insarJob.FinishedAt,
	)
	if err != nil {
		return objects.InsarJob{}, err
	}
	return insarJob, nil
}

func (r *RepositoryInsarJob) CreateInsarJob(ctx context.Context, insarJob objects.InsarJob) (objects.InsarJob, error) {
	var id int
	var objectID int
	var hyp3JobID string
	var createdAt time.Time
	query := "INSERT INTO insar_jobs (job_type, reference_granule, secondary_granule, reference_date, secondary_date, status, finished_at) VALUSE ($1 $2 $3 $4 $5 $6 $7) RETURNING id, object_id, hyp3_job_id, created_at"
	err := r.pool.QueryRow(ctx, query, insarJob.JobType,
		insarJob.ReferenceGranule, insarJob.SecondaryGranule,
		insarJob.ReferenceDate, insarJob.SecondaryDate,
		insarJob.Status, insarJob.FinishedAt,
	).Scan(&id, &objectID, &hyp3JobID, &createdAt)
	if err != nil {
		return objects.InsarJob{}, err
	}
	insarJob.ID = id
	insarJob.ObjectID = objectID
	insarJob.Hyp3JobID = hyp3JobID
	insarJob.CreatedAt = createdAt
	return insarJob, nil
}

func (r *RepositoryInsarJob) UpdateInsarJob(ctx context.Context, id int, insarJob objects.InsarJob) (objects.InsarJob, error) {
	var updatedInsarJob objects.InsarJob
	query := "UPDATE insar_jobs SET job_type = $1, reference_granule = $2, secondary_granule = $3, reference_date = $4, secondary_date = $5, status = $6, finished_at = $7 WHERE id = $8 RETURNING id, object_id, hyp3_job_id, job_type, reference_granule, secondary_granule, reference_date, secondary_date, status, finished_at"
	err := r.pool.QueryRow(ctx, query,
		insarJob.JobType,
		insarJob.ReferenceGranule, insarJob.SecondaryGranule,
		insarJob.ReferenceDate, insarJob.SecondaryDate,
		insarJob.Status, insarJob.FinishedAt, id,
	).Scan(
		&insarJob.ID, &insarJob.ObjectID, &insarJob.Hyp3JobID,
		&insarJob.JobType, &insarJob.ReferenceGranule, &insarJob.SecondaryGranule,
		&insarJob.ReferenceDate, &insarJob.SecondaryDate, &insarJob.Status,
		&insarJob.FinishedAt,
	)
	if err != nil {
		return objects.InsarJob{}, err
	}
	return updatedInsarJob, nil
}

func (r *RepositoryInsarJob) DeleteInsarJob(ctx context.Context, id int) (objects.InsarJob, error) {
	var deletedInsarJob objects.InsarJob
	query := "DELETE FROM insar_job WHERE id = $1 RETURNING *"
	err := r.pool.QueryRow(ctx, query, id).Scan(&deletedInsarJob.ID,
		&deletedInsarJob.ObjectID, &deletedInsarJob.Hyp3JobID,
		&deletedInsarJob.JobType, &deletedInsarJob.ReferenceGranule,
		&deletedInsarJob.SecondaryGranule, &deletedInsarJob.ReferenceDate,
		&deletedInsarJob.SecondaryDate, &deletedInsarJob.Status,
		&deletedInsarJob.CreatedAt, &deletedInsarJob.FinishedAt,
	)
	if err != nil {
		return objects.InsarJob{}, err
	}
	return deletedInsarJob, nil
}

func NewRepositoryInsarJob(pool *pgxpool.Pool) *RepositoryInsarJob {
	return &RepositoryInsarJob{pool: pool}
}
