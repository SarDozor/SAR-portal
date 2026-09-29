package objects

import "time"

type Object struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Country    string    `json:"country"`
	Latitude   *float64  `json:"latitude"`
	Longitude  *float64  `json:"longitude"`
	RaiseType  string    `json:"raise_type,omitempty"`
	MaxHeightM *float64  `json:"max_height_m"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Status     string    `json:"status"`
}

type InsarJob struct {
	ID               int        `json:"id"`
	ObjectId         int        `json:"object_id"`
	Hyp3JobId        string     `json:"hyp3_job_id"`
	JobType          string     `json:"job_type"`
	ReferenceGranule string     `json:"reference_granule"`
	SecondaryGranule string     `json:"secondary_granule"`
	ReferenceDate    time.Time  `json:"reference_date"`
	SecondaryDate    time.Time  `json:"secondary_date"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	FinishedAt       *time.Time `json:"finished_at"`
}

type Measurement struct {
	ID             int       `json:"id"`
	JobID          int       `jsob:"json_id"`
	Latitude       *float64  `json:"latitude"`
	Longitude      *float64  `json:"longitude"`
	DisplacementMM float64   `json:"displacement_mm"`
	Coherence      *float64  `json:"coherence"`
	CreatedAt      time.Time `json:"created_at"`
}
