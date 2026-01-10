-- migrate:up
CREATE TABLE IF NOT EXISTS metric_attributes (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    metric_id       TEXT NOT NULL REFERENCES metrics(id) ON DELETE CASCADE,
    attribute_name  TEXT NOT NULL,
    attribute_type  TEXT,
    description     TEXT,
    required        INTEGER DEFAULT 0,
    UNIQUE(metric_id, attribute_name)
);

CREATE TABLE IF NOT EXISTS attribute_enum_values (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    attribute_id    INTEGER NOT NULL REFERENCES metric_attributes(id) ON DELETE CASCADE,
    enum_value      TEXT NOT NULL,
    UNIQUE(attribute_id, enum_value)
);

CREATE INDEX IF NOT EXISTS idx_metric_attributes_metric_id ON metric_attributes(metric_id);
CREATE INDEX IF NOT EXISTS idx_metric_attributes_name ON metric_attributes(attribute_name);
CREATE INDEX IF NOT EXISTS idx_attribute_enum_values_attr_id ON attribute_enum_values(attribute_id);

-- migrate:down
DROP INDEX IF EXISTS idx_attribute_enum_values_attr_id;
DROP INDEX IF EXISTS idx_metric_attributes_name;
DROP INDEX IF EXISTS idx_metric_attributes_metric_id;
DROP TABLE IF EXISTS attribute_enum_values;
DROP TABLE IF EXISTS metric_attributes;
