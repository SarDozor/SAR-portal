BEGIN;

-- objects

ALTER TABLE objects DROP COLUMN construction_year;

-- dam_type -> raise_type
ALTER TABLE objects RENAME COLUMN dam_type TO raise_type;
ALTER TABLE objects ALTER COLUMN raise_type DROP NOT NULL;

-- height_m -> max_height_m
ALTER TABLE objects RENAME COLUMN height_m TO max_height_m;

-- статус хвостохранилища
ALTER TABLE objects ADD COLUMN status TEXT;

ALTER TABLE objects ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now();

ALTER TABLE objects
    ADD CONSTRAINT objects_latitude_range  CHECK (latitude  BETWEEN -90  AND 90),
    ADD CONSTRAINT objects_longitude_range CHECK (longitude BETWEEN -180 AND 180),
    ADD CONSTRAINT objects_height_positive CHECK (max_height_m >= 0);

CREATE INDEX objects_country_idx ON objects (country);
CREATE INDEX objects_status_idx  ON objects (status);

-- insar_jobs

CREATE TABLE insar_jobs (
    id                INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    object_id         INTEGER NOT NULL REFERENCES objects (id) ON DELETE CASCADE,
    hyp3_job_id       TEXT UNIQUE,                -- id задачи на стороне NASA
    job_type          TEXT NOT NULL DEFAULT 'INSAR_GAMMA',
    reference_granule TEXT NOT NULL,              -- первый снимок
    secondary_granule TEXT NOT NULL,              -- второй снимок
    reference_date    DATE NOT NULL,
    secondary_date    DATE NOT NULL,
    status            TEXT NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'RUNNING', 'SUCCEEDED', 'FAILED')),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at       TIMESTAMPTZ,

    CHECK (secondary_date > reference_date)
);

CREATE INDEX insar_jobs_object_id_idx ON insar_jobs (object_id);

-- measurements

CREATE TABLE measurements (
    id              BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    job_id          INTEGER NOT NULL REFERENCES insar_jobs (id) ON DELETE CASCADE,
    latitude        DOUBLE PRECISION CHECK (latitude  BETWEEN -90  AND 90),
    longitude       DOUBLE PRECISION CHECK (longitude BETWEEN -180 AND 180),
    displacement_mm DOUBLE PRECISION NOT NULL,    -- вдоль линии визирования, + к спутнику
    coherence       DOUBLE PRECISION CHECK (coherence BETWEEN 0 AND 1),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX measurements_job_id_idx ON measurements (job_id);

COMMIT;
