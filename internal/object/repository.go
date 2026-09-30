package object

import (
	"context"

	"github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const columns = "id, name, country, latitude, longitude, status, raise_type, max_height_m, created_at, updated_at"

type Repository struct {
	pool *pgxpool.Pool
}

func scanObject(row pgx.Row) (model.Object, error) {
	var o model.Object
	err := row.Scan(
		&o.ID, &o.Name, &o.Country, &o.Latitude, &o.Longitude,
		&o.Status, &o.RaiseType, &o.MaxHeightM, &o.CreatedAt, &o.UpdatedAt,
	)
	return o, err
}

func (r *Repository) ListObjects(ctx context.Context) ([]model.Object, error) {
	query := "SELECT " + columns + " FROM objects ORDER BY id"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	listObjects := []model.Object{}
	for rows.Next() {
		object, err := scanObject(rows)
		if err != nil {
			return nil, err
		}
		listObjects = append(listObjects, object)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return listObjects, nil
}

func (r *Repository) GetObjectByID(ctx context.Context, id int) (model.Object, error) {
	query := "SELECT " + columns + " FROM objects WHERE id = $1"
	return scanObject(r.pool.QueryRow(ctx, query, id))
}

func (r *Repository) CreateObject(ctx context.Context, obj model.Object) (model.Object, error) {
	query := "INSERT INTO objects (name, country, latitude, longitude, status, raise_type, max_height_m) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING " + columns
	return scanObject(r.pool.QueryRow(
		ctx, query, obj.Name, obj.Country, obj.Latitude, obj.Longitude,
		obj.Status, obj.RaiseType, obj.MaxHeightM,
	))
}

func (r *Repository) UpdateObject(ctx context.Context, id int, obj model.Object) (model.Object, error) {
	query := "UPDATE objects SET name = $1, country = $2, latitude = $3, longitude = $4, status = $5, raise_type = $6, max_height_m = $7, updated_at = now() WHERE id = $8 RETURNING " + columns
	return scanObject(r.pool.QueryRow(
		ctx, query, obj.Name, obj.Country, obj.Latitude, obj.Longitude,
		obj.Status, obj.RaiseType, obj.MaxHeightM, id,
	))
}

func (r *Repository) DeleteObject(ctx context.Context, id int) (model.Object, error) {
	query := "DELETE FROM objects WHERE id = $1 RETURNING " + columns
	return scanObject(r.pool.QueryRow(ctx, query, id))
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
