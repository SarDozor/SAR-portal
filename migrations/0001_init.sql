BEGIN;

CREATE TABLE objects (
    id           INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name         TEXT NOT NULL,
    country      TEXT NOT NULL,
    latitude     DOUBLE PRECISION CHECK (latitude  BETWEEN -90  AND 90),
    longitude    DOUBLE PRECISION CHECK (longitude BETWEEN -180 AND 180),
    status       TEXT,
    raise_type   TEXT,
    max_height_m DOUBLE PRECISION CHECK (max_height_m >= 0),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX objects_country_idx ON objects (country);
CREATE INDEX objects_status_idx  ON objects (status);


CREATE TABLE insar_jobs (
    id                INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    object_id         INTEGER NOT NULL REFERENCES objects (id) ON DELETE CASCADE,
    pair_name         TEXT NOT NULL,
    hyp3_job_id       TEXT UNIQUE,     
    job_type          TEXT NOT NULL DEFAULT 'INSAR_GAMMA',
    reference_granule TEXT NOT NULL,   
    secondary_granule TEXT NOT NULL,
    reference_date    DATE NOT NULL,
    secondary_date    DATE NOT NULL,
    status            TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'RUNNING', 'SUCCEEDED', 'FAILED')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at       TIMESTAMPTZ,

    CHECK (secondary_date > reference_date),
    UNIQUE (object_id, pair_name)
);

CREATE INDEX insar_jobs_object_id_idx ON insar_jobs (object_id);

CREATE TABLE measurements (
    id                          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    job_id                      INTEGER NOT NULL REFERENCES insar_jobs (id) ON DELETE CASCADE,
    latitude                    DOUBLE PRECISION NOT NULL CHECK (latitude  BETWEEN -90  AND 90),
    longitude                   DOUBLE PRECISION NOT NULL CHECK (longitude BETWEEN -180 AND 180),
    reference_method            TEXT NOT NULL,
    displacement_mm             DOUBLE PRECISION NOT NULL,
    dam_coherence_mean          DOUBLE PRECISION NOT NULL CHECK(dam_coherence_mean BETWEEN 0 AND 1),
    dam_coherence_valid_pixels  INTEGER NOT NULL CHECK(dam_coherence_valid_pixels >= 0),
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (job_id, reference_method)
);

CREATE INDEX measurements_job_id_idx ON measurements (job_id);

COMMIT;
