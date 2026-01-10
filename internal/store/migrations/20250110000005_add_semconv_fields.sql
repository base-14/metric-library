-- migrate:up
ALTER TABLE metrics ADD COLUMN semconv_match TEXT DEFAULT '';
ALTER TABLE metrics ADD COLUMN semconv_name TEXT DEFAULT '';
ALTER TABLE metrics ADD COLUMN semconv_stability TEXT DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_metrics_semconv_match ON metrics(semconv_match);

-- migrate:down
DROP INDEX IF EXISTS idx_metrics_semconv_match;
-- SQLite doesn't support DROP COLUMN, so we leave the columns
