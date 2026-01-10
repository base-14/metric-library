-- migrate:up
CREATE TABLE IF NOT EXISTS metrics (
    id                  TEXT PRIMARY KEY,
    metric_name         TEXT NOT NULL,
    instrument_type     TEXT NOT NULL,
    description         TEXT,
    unit                TEXT,
    enabled_by_default  INTEGER DEFAULT 1,

    -- Component info
    component_type      TEXT NOT NULL,
    component_name      TEXT NOT NULL,

    -- Source info
    source_category     TEXT NOT NULL,
    source_name         TEXT NOT NULL,
    source_location     TEXT,

    -- Provenance
    extraction_method   TEXT NOT NULL,
    source_confidence   TEXT NOT NULL,
    repo                TEXT,
    path                TEXT,
    "commit"            TEXT,
    extracted_at        TIMESTAMP NOT NULL,

    -- Metadata
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_metrics_metric_name ON metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_metrics_instrument_type ON metrics(instrument_type);
CREATE INDEX IF NOT EXISTS idx_metrics_component_type ON metrics(component_type);
CREATE INDEX IF NOT EXISTS idx_metrics_component_name ON metrics(component_name);
CREATE INDEX IF NOT EXISTS idx_metrics_source_category ON metrics(source_category);
CREATE INDEX IF NOT EXISTS idx_metrics_source_name ON metrics(source_name);
CREATE INDEX IF NOT EXISTS idx_metrics_source_confidence ON metrics(source_confidence);

-- migrate:down
DROP INDEX IF EXISTS idx_metrics_source_confidence;
DROP INDEX IF EXISTS idx_metrics_source_name;
DROP INDEX IF EXISTS idx_metrics_source_category;
DROP INDEX IF EXISTS idx_metrics_component_name;
DROP INDEX IF EXISTS idx_metrics_component_type;
DROP INDEX IF EXISTS idx_metrics_instrument_type;
DROP INDEX IF EXISTS idx_metrics_metric_name;
DROP TABLE IF EXISTS metrics;
