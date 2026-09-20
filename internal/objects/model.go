package objects

import "time"

type Object struct {
	ID               int
	Name             string
	Country          string
	Latitude         *float64
	Longitude        *float64
	DamType          string
	HeightM          *float64
	ConstructionYear *int
	CreatedAt        time.Time
}
