package objects

import "time"

type Object struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Country          string    `json:"country"`
	Latitude         *float64  `json:"latitude"`
	Longitude        *float64  `json:"longitude"`
	DamType          string    `json:"dam_type"`
	HeightM          *float64  `json:"height_m"`
	ConstructionYear *int      `json:"construction_year"`
	CreatedAt        time.Time `json:"created_at"`
}
