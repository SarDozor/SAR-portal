package objects

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func (r *Repository) ListObjects(ctx context.Context) ([]Object, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, country, latitude, longitude, dam_type, height_m, construction_year, created_at FROM objects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var listObjects []Object
	for rows.Next() {
		var object Object
		err := rows.Scan(&object.ID, &object.Name, &object.Country, &object.Latitude, &object.Longitude, &object.DamType, &object.HeightM, &object.ConstructionYear, &object.CreatedAt)
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

func (r *Repository) CreateObject(ctx context.Context, obj Object) (Object, error) {
	var id int
	var createdAt time.Time
	query := "INSERT INTO objects (name, country, latitude, longitude, dam_type, height_m, construction_year) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at"
	err := r.pool.QueryRow(ctx, query, obj.Name, obj.Country, obj.Latitude, obj.Longitude, obj.DamType, obj.HeightM, obj.ConstructionYear).Scan(&id, &createdAt)
	if err != nil {
		return Object{}, err
	}
	obj.ID = id
	obj.CreatedAt = createdAt
	return obj, nil
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}
