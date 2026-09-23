CREATE TABLE hyp3_jobs (
    id SERIAL PRIMARY KEY,
    object_id INTEGER NOT NULL REFERENCES objects(id),
    hyp3_job_id TEXT UNIQUE,
    job_type TEXT NOT NULL DEFAULT 'INSAR_GAMMA',
    granule_reference TEXT NOT NULL,
    granule_secondary TEXT NOT NULL,
    status_code TEXT NOT NULL DEFAULT 'PENDING',
    requested_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    result_url TEXT,
    error_message TEXT
);
