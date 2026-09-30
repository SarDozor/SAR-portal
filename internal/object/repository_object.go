package objects

import (
	"context"
	"time"

	objects "github.com/darkn3sssss/sar-dozor/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryObject struct {
	pool *pgxpool.Pool
}

func (r *RepositoryObject) ListObjects(ctx context.Context) ([]objects.Object, error) {
	query := "SELECT * FROM objects"
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var listObjects []objects.Object
	for rows.Next() {
		var object objects.Object
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

func (r *RepositoryObject) GetObjectByID(ctx context.Context, id int) (objects.Object, error) {
	query := "SELECT * FROM objects WHERE id = $1"
	var object objects.Object
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&object.ID, &object.Name, &object.Country, &object.Latitude,
		&object.Longitude, &object.RaiseType, &object.MaxHeightM,
		&object.CreatedAt, &object.UpdatedAt, &object.Status,
	)
	if err != nil {
		return objects.Object{}, err
	}
	return object, nil
}

func (r *RepositoryObject) CreateObject(ctx context.Context, obj objects.Object) (objects.Object, error) {
	var id int
	var createdAt time.Time
	var updatedAt time.Time
	query := "INSERT INTO objects (name, country, latitude, longitude, raise_type, max_height_m, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at"
	err := r.pool.QueryRow(
		ctx, query, obj.Name, obj.Country, obj.Latitude,
		obj.Longitude, obj.RaiseType, obj.MaxHeightM, obj.Status,
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return objects.Object{}, err
	}
	obj.ID = id
	obj.CreatedAt = createdAt
	obj.UpdatedAt = updatedAt
	return obj, nil
}

func (r *RepositoryObject) UpdateObject(ctx context.Context, id int, obj objects.Object) (objects.Object, error) {
	query := "UPDATE objects SET name = $1, country = $2, latitude = $3, longitude = $4, raise_type = $5, max_height_m = $6, status = $7, updated_at = now() WHERE id = $8 RETURNING id, name, country, latitude, longitude, raise_type, max_height_m, status, updated_at"
	var updatedObj objects.Object
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
		return objects.Object{}, err
	}
	return updatedObj, nil
}

func (r *RepositoryObject) DeleteObject(ctx context.Context, id int) (objects.Object, error) {
	query := "DELETE FROM objects WHERE id = $1 RETURNING *"
	var deletedObject objects.Object
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&deletedObject.ID, &deletedObject.Name, &deletedObject.Country,
		&deletedObject.Latitude, &deletedObject.Longitude, &deletedObject.RaiseType,
		&deletedObject.MaxHeightM, &deletedObject.Status,
		&deletedObject.CreatedAt, &deletedObject.UpdatedAt,
	)
	if err != nil {
		return objects.Object{}, err
	}
	return deletedObject, nil
}

func NewRepositoryObject(pool *pgxpool.Pool) *RepositoryObject {
	return &RepositoryObject{pool: pool}
}
