-- migrate:up
CREATE TABLE IF NOT EXISTS extraction_runs (
    id              TEXT PRIMARY KEY,
    adapter_name    TEXT NOT NULL,
    "commit"        TEXT,
    started_at      TIMESTAMP NOT NULL,
    completed_at    TIMESTAMP,
    metrics_count   INTEGER DEFAULT 0,
    status          TEXT NOT NULL DEFAULT 'running',
    error_message   TEXT
);

CREATE INDEX IF NOT EXISTS idx_extraction_runs_adapter ON extraction_runs(adapter_name);
CREATE INDEX IF NOT EXISTS idx_extraction_runs_status ON extraction_runs(status);
CREATE INDEX IF NOT EXISTS idx_extraction_runs_started_at ON extraction_runs(started_at);

-- migrate:down
DROP INDEX IF EXISTS idx_extraction_runs_started_at;
DROP INDEX IF EXISTS idx_extraction_runs_status;
DROP INDEX IF EXISTS idx_extraction_runs_adapter;
DROP TABLE IF EXISTS extraction_runs;
