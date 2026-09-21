CREATE TABLE objects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    dam_type TEXT NOT NULL,
    height_m DOUBLE PRECISION,
    construction_year INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);