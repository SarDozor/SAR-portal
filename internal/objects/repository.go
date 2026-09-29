package objects

import (
	"context"
	"time"

	"/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func (r *Repository) ListObjects(ctx context.Context) ([]model.Object, error) {
	query := "SELECT id, name, country, latitude, longitude, raise_type, max_height_m, created_at, updated_at, status FROM objects"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var listObjects []model.Object
	for rows.Next() {
		var object model.Object
		err := rows.Scan(
			&object.ID, &object.Name, &object.Country, &object.Latitude,
			&object.Longitude, &object.RaiseType, &object.MaxHeightM,
			&object.CreatedAt, &object.UpdatedAt, &object.Status,
		)
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
	query := "SELECT id, name, country, latitude, longitude, raise_type, max_height_m, created_at, updated_at, status FROM objects WHERE id = $1"
	var object model.Object
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&object.ID, &object.Name, &object.Country, &object.Latitude,
		&object.Longitude, &object.RaiseType, &object.MaxHeightM,
		&object.CreatedAt, &object.UpdatedAt, &object.Status,
	)
	if err != nil {
		return model.Object{}, err
	}
	return object, nil
}

func (r *Repository) CreateObject(ctx context.Context, obj model.Object) (model.Object, error) {
	var id int
	var createdAt time.Time
	var updatedAt time.Time
	query := "INSERT INTO objects (name, country, latitude, longitude, raise_type, max_height_m, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at"
	err := r.pool.QueryRow(
		ctx, query, obj.Name, obj.Country, obj.Latitude,
		obj.Longitude, obj.RaiseType, obj.MaxHeightM, obj.Status,
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return model.Object{}, err
	}
	obj.ID = id
	obj.CreatedAt = createdAt
	obj.UpdatedAt = updatedAt
	return obj, nil
}

func (r *Repository) UpdateObject(ctx context.Context, id int, obj model.Object) (model.Object, error) {
	query := "UPDATE objects SET name = $1, country = $2, latitude = $3, longitude = $4, raise_type = $5, max_height_m = $6, status = $7, updated_at = now() WHERE id = $8 RETURNING id, name, country, latitude, longitude, raise_type, max_height_m, status, updated_at"
	var updatedObj model.Object
	err := r.pool.QueryRow(
		ctx, query, obj.Name, obj.Country, obj.Latitude, obj.Longitude,
		obj.RaiseType, obj.MaxHeightM, obj.Status, id,
	).Scan(
		&updatedObj.ID, &updatedObj.Name, &updatedObj.Country,
		&updatedObj.Latitude, &updatedObj.Longitude, &updatedObj.RaiseType,
		&updatedObj.MaxHeightM, &updatedObj.Status,
		&updatedObj.UpdatedAt,
	)
	if err != nil {
		return model.Object{}, err
	}
	return updatedObj, nil
}

func (r *Repository) DeleteObject(ctx context.Context, id int) (model.Object, error) {
	query := "DELETE FROM objects WHERE id = $1 RETURNING id, name, country, latitude, longitude, raise_type, max_height_m, status, created_at, updated_at"
	var deletedObject model.Object
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&deletedObject.ID, &deletedObject.Name, &deletedObject.Country,
		&deletedObject.Latitude, &deletedObject.Longitude, &deletedObject.RaiseType,
		&deletedObject.MaxHeightM, &deletedObject.Status,
		&deletedObject.CreatedAt, &deletedObject.UpdatedAt,
	)
	if err != nil {
		return model.Object{}, err
	}
	return deletedObject, nil
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
