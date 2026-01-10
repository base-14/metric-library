-- migrate:up
CREATE VIRTUAL TABLE IF NOT EXISTS metrics_fts USING fts5(
    metric_name,
    description,
    component_name,
    source_name,
    content='metrics',
    content_rowid='rowid'
);

-- Triggers to keep FTS index in sync
CREATE TRIGGER IF NOT EXISTS metrics_ai AFTER INSERT ON metrics BEGIN
    INSERT INTO metrics_fts(rowid, metric_name, description, component_name, source_name)
    VALUES (NEW.rowid, NEW.metric_name, NEW.description, NEW.component_name, NEW.source_name);
END;

CREATE TRIGGER IF NOT EXISTS metrics_ad AFTER DELETE ON metrics BEGIN
    INSERT INTO metrics_fts(metrics_fts, rowid, metric_name, description, component_name, source_name)
    VALUES ('delete', OLD.rowid, OLD.metric_name, OLD.description, OLD.component_name, OLD.source_name);
END;

CREATE TRIGGER IF NOT EXISTS metrics_au AFTER UPDATE ON metrics BEGIN
    INSERT INTO metrics_fts(metrics_fts, rowid, metric_name, description, component_name, source_name)
    VALUES ('delete', OLD.rowid, OLD.metric_name, OLD.description, OLD.component_name, OLD.source_name);
    INSERT INTO metrics_fts(rowid, metric_name, description, component_name, source_name)
    VALUES (NEW.rowid, NEW.metric_name, NEW.description, NEW.component_name, NEW.source_name);
END;

-- migrate:down
DROP TRIGGER IF EXISTS metrics_au;
DROP TRIGGER IF EXISTS metrics_ad;
DROP TRIGGER IF EXISTS metrics_ai;
DROP TABLE IF EXISTS metrics_fts;
