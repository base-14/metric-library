package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/base14/otel-glossary/internal/domain"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) UpsertMetric(ctx context.Context, metric *domain.CanonicalMetric) error {
	metric.EnsureID()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := s.upsertMetricTx(ctx, tx, metric); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *SQLiteStore) upsertMetricTx(ctx context.Context, tx *sql.Tx, metric *domain.CanonicalMetric) error {
	query := `
		INSERT INTO metrics (
			id, metric_name, instrument_type, description, unit, enabled_by_default,
			component_type, component_name, source_category, source_name, source_location,
			extraction_method, source_confidence, repo, path, "commit", extracted_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(id) DO UPDATE SET
			metric_name = excluded.metric_name,
			instrument_type = excluded.instrument_type,
			description = excluded.description,
			unit = excluded.unit,
			enabled_by_default = excluded.enabled_by_default,
			component_type = excluded.component_type,
			component_name = excluded.component_name,
			source_category = excluded.source_category,
			source_name = excluded.source_name,
			source_location = excluded.source_location,
			extraction_method = excluded.extraction_method,
			source_confidence = excluded.source_confidence,
			repo = excluded.repo,
			path = excluded.path,
			"commit" = excluded."commit",
			extracted_at = excluded.extracted_at,
			updated_at = CURRENT_TIMESTAMP
	`

	enabledByDefault := 0
	if metric.EnabledByDefault {
		enabledByDefault = 1
	}

	_, err := tx.ExecContext(ctx, query,
		metric.ID, metric.MetricName, metric.InstrumentType, metric.Description, metric.Unit, enabledByDefault,
		metric.ComponentType, metric.ComponentName, metric.SourceCategory, metric.SourceName, metric.SourceLocation,
		metric.ExtractionMethod, metric.SourceConfidence, metric.Repo, metric.Path, metric.Commit, metric.ExtractedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert metric: %w", err)
	}

	// Delete existing attributes
	_, err = tx.ExecContext(ctx, "DELETE FROM metric_attributes WHERE metric_id = ?", metric.ID)
	if err != nil {
		return fmt.Errorf("failed to delete existing attributes: %w", err)
	}

	// Insert attributes
	for _, attr := range metric.Attributes {
		required := 0
		if attr.Required {
			required = 1
		}

		result, err := tx.ExecContext(ctx,
			"INSERT INTO metric_attributes (metric_id, attribute_name, attribute_type, description, required) VALUES (?, ?, ?, ?, ?)",
			metric.ID, attr.Name, attr.Type, attr.Description, required,
		)
		if err != nil {
			return fmt.Errorf("failed to insert attribute: %w", err)
		}

		if len(attr.Enum) > 0 {
			attrID, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get attribute id: %w", err)
			}

			for _, enumVal := range attr.Enum {
				_, err := tx.ExecContext(ctx,
					"INSERT INTO attribute_enum_values (attribute_id, enum_value) VALUES (?, ?)",
					attrID, enumVal,
				)
				if err != nil {
					return fmt.Errorf("failed to insert enum value: %w", err)
				}
			}
		}
	}

	return nil
}

func (s *SQLiteStore) UpsertMetrics(ctx context.Context, metrics []*domain.CanonicalMetric) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, metric := range metrics {
		metric.EnsureID()
		if err := s.upsertMetricTx(ctx, tx, metric); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *SQLiteStore) GetMetric(ctx context.Context, id string) (*domain.CanonicalMetric, error) {
	query := `
		SELECT id, metric_name, instrument_type, description, unit, enabled_by_default,
			component_type, component_name, source_category, source_name, source_location,
			extraction_method, source_confidence, repo, path, "commit", extracted_at
		FROM metrics WHERE id = ?
	`

	var metric domain.CanonicalMetric
	var enabledByDefault int
	var description, unit, sourceLocation, repo, path, commit sql.NullString

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&metric.ID, &metric.MetricName, &metric.InstrumentType, &description, &unit, &enabledByDefault,
		&metric.ComponentType, &metric.ComponentName, &metric.SourceCategory, &metric.SourceName, &sourceLocation,
		&metric.ExtractionMethod, &metric.SourceConfidence, &repo, &path, &commit, &metric.ExtractedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get metric: %w", err)
	}

	metric.Description = description.String
	metric.Unit = unit.String
	metric.SourceLocation = sourceLocation.String
	metric.Repo = repo.String
	metric.Path = path.String
	metric.Commit = commit.String
	metric.EnabledByDefault = enabledByDefault == 1

	// Get attributes
	attrs, err := s.getMetricAttributes(ctx, id)
	if err != nil {
		return nil, err
	}
	metric.Attributes = attrs

	return &metric, nil
}

func (s *SQLiteStore) getMetricAttributes(ctx context.Context, metricID string) ([]domain.Attribute, error) {
	query := `
		SELECT id, attribute_name, attribute_type, description, required
		FROM metric_attributes WHERE metric_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, metricID)
	if err != nil {
		return nil, fmt.Errorf("failed to query attributes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var attrs []domain.Attribute
	for rows.Next() {
		var attr domain.Attribute
		var attrID int64
		var required int
		var attrType, description sql.NullString

		if err := rows.Scan(&attrID, &attr.Name, &attrType, &description, &required); err != nil {
			return nil, fmt.Errorf("failed to scan attribute: %w", err)
		}

		attr.Type = attrType.String
		attr.Description = description.String
		attr.Required = required == 1

		// Get enum values
		enumRows, err := s.db.QueryContext(ctx, "SELECT enum_value FROM attribute_enum_values WHERE attribute_id = ?", attrID)
		if err != nil {
			return nil, fmt.Errorf("failed to query enum values: %w", err)
		}

		for enumRows.Next() {
			var enumVal string
			if err := enumRows.Scan(&enumVal); err != nil {
				_ = enumRows.Close()
				return nil, fmt.Errorf("failed to scan enum value: %w", err)
			}
			attr.Enum = append(attr.Enum, enumVal)
		}
		_ = enumRows.Close()

		attrs = append(attrs, attr)
	}

	return attrs, rows.Err()
}

func (s *SQLiteStore) DeleteMetric(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM metrics WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete metric: %w", err)
	}
	return nil
}

func (s *SQLiteStore) DeleteMetricsBySource(ctx context.Context, sourceName string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM metrics WHERE source_name = ?", sourceName)
	if err != nil {
		return fmt.Errorf("failed to delete metrics by source: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Search(ctx context.Context, query SearchQuery) (*SearchResult, error) {
	start := time.Now()

	var conditions []string
	var args []interface{}

	// Substring search using LIKE on metric_name and description only
	if query.Text != "" {
		searchPattern := "%" + query.Text + "%"
		conditions = append(conditions, "(m.metric_name LIKE ? OR m.description LIKE ?)")
		args = append(args, searchPattern, searchPattern)
	}

	// Filters
	if len(query.InstrumentTypes) > 0 {
		placeholders := make([]string, len(query.InstrumentTypes))
		for i, t := range query.InstrumentTypes {
			placeholders[i] = "?"
			args = append(args, t)
		}
		conditions = append(conditions, fmt.Sprintf("m.instrument_type IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.ComponentTypes) > 0 {
		placeholders := make([]string, len(query.ComponentTypes))
		for i, t := range query.ComponentTypes {
			placeholders[i] = "?"
			args = append(args, t)
		}
		conditions = append(conditions, fmt.Sprintf("m.component_type IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.ComponentNames) > 0 {
		placeholders := make([]string, len(query.ComponentNames))
		for i, n := range query.ComponentNames {
			placeholders[i] = "?"
			args = append(args, n)
		}
		conditions = append(conditions, fmt.Sprintf("m.component_name IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.SourceCategories) > 0 {
		placeholders := make([]string, len(query.SourceCategories))
		for i, c := range query.SourceCategories {
			placeholders[i] = "?"
			args = append(args, c)
		}
		conditions = append(conditions, fmt.Sprintf("m.source_category IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.SourceNames) > 0 {
		placeholders := make([]string, len(query.SourceNames))
		for i, n := range query.SourceNames {
			placeholders[i] = "?"
			args = append(args, n)
		}
		conditions = append(conditions, fmt.Sprintf("m.source_name IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.ConfidenceLevels) > 0 {
		placeholders := make([]string, len(query.ConfidenceLevels))
		for i, c := range query.ConfidenceLevels {
			placeholders[i] = "?"
			args = append(args, c)
		}
		conditions = append(conditions, fmt.Sprintf("m.source_confidence IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.Units) > 0 {
		placeholders := make([]string, len(query.Units))
		for i, u := range query.Units {
			placeholders[i] = "?"
			args = append(args, u)
		}
		conditions = append(conditions, fmt.Sprintf("m.unit IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.AttributeNames) > 0 {
		placeholders := make([]string, len(query.AttributeNames))
		for i, n := range query.AttributeNames {
			placeholders[i] = "?"
			args = append(args, n)
		}
		conditions = append(conditions, fmt.Sprintf("m.id IN (SELECT metric_id FROM metric_attributes WHERE attribute_name IN (%s))", strings.Join(placeholders, ",")))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM metrics m %s", whereClause) //nolint:gosec // SQL injection not possible - whereClause uses parameterized queries
	var total int
	if err := s.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count metrics: %w", err)
	}

	// Get metrics
	limit := query.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := query.Offset
	if offset < 0 {
		offset = 0
	}

	// Order by: metric_name matches first, then description matches, then alphabetically
	orderClause := "ORDER BY m.metric_name"
	if query.Text != "" {
		searchPattern := "%" + query.Text + "%"
		orderClause = `ORDER BY
			CASE WHEN m.metric_name LIKE ? THEN 0 ELSE 1 END,
			CASE WHEN m.description LIKE ? THEN 0 ELSE 1 END,
			m.metric_name`
		args = append(args, searchPattern, searchPattern)
	}

	//nolint:gosec // SQL injection not possible - whereClause and orderClause use parameterized queries
	selectQuery := fmt.Sprintf(`
		SELECT id, metric_name, instrument_type, description, unit, enabled_by_default,
			component_type, component_name, source_category, source_name, source_location,
			extraction_method, source_confidence, repo, path, "commit", extracted_at
		FROM metrics m %s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)

	args = append(args, limit, offset)

	rows, err := s.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var metrics []*domain.CanonicalMetric
	for rows.Next() {
		var metric domain.CanonicalMetric
		var enabledByDefault int
		var description, unit, sourceLocation, repo, path, commit sql.NullString

		if err := rows.Scan(
			&metric.ID, &metric.MetricName, &metric.InstrumentType, &description, &unit, &enabledByDefault,
			&metric.ComponentType, &metric.ComponentName, &metric.SourceCategory, &metric.SourceName, &sourceLocation,
			&metric.ExtractionMethod, &metric.SourceConfidence, &repo, &path, &commit, &metric.ExtractedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %w", err)
		}

		metric.Description = description.String
		metric.Unit = unit.String
		metric.SourceLocation = sourceLocation.String
		metric.Repo = repo.String
		metric.Path = path.String
		metric.Commit = commit.String
		metric.EnabledByDefault = enabledByDefault == 1

		// Get attributes (could be optimized with a join)
		attrs, err := s.getMetricAttributes(ctx, metric.ID)
		if err != nil {
			return nil, err
		}
		metric.Attributes = attrs

		metrics = append(metrics, &metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate metrics: %w", err)
	}

	return &SearchResult{
		Metrics: metrics,
		Total:   total,
		Took:    time.Since(start),
	}, nil
}

func (s *SQLiteStore) GetFacetCounts(ctx context.Context) (*FacetCounts, error) {
	facets := &FacetCounts{
		InstrumentTypes:  make(map[domain.InstrumentType]int),
		ComponentTypes:   make(map[domain.ComponentType]int),
		ComponentNames:   make(map[string]int),
		SourceCategories: make(map[domain.SourceCategory]int),
		SourceNames:      make(map[string]int),
		ConfidenceLevels: make(map[domain.ConfidenceLevel]int),
		Units:            make(map[string]int),
	}

	queries := []struct {
		query  string
		target interface{}
	}{
		{"SELECT instrument_type, COUNT(*) FROM metrics GROUP BY instrument_type", &facets.InstrumentTypes},
		{"SELECT component_type, COUNT(*) FROM metrics GROUP BY component_type", &facets.ComponentTypes},
		{"SELECT component_name, COUNT(*) FROM metrics GROUP BY component_name", &facets.ComponentNames},
		{"SELECT source_category, COUNT(*) FROM metrics GROUP BY source_category", &facets.SourceCategories},
		{"SELECT source_name, COUNT(*) FROM metrics GROUP BY source_name", &facets.SourceNames},
		{"SELECT source_confidence, COUNT(*) FROM metrics GROUP BY source_confidence", &facets.ConfidenceLevels},
		{"SELECT unit, COUNT(*) FROM metrics WHERE unit IS NOT NULL AND unit != '' GROUP BY unit", &facets.Units},
	}

	for _, q := range queries {
		rows, err := s.db.QueryContext(ctx, q.query)
		if err != nil {
			return nil, fmt.Errorf("failed to query facets: %w", err)
		}

		for rows.Next() {
			var key string
			var count int
			if err := rows.Scan(&key, &count); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("failed to scan facet: %w", err)
			}

			switch target := q.target.(type) {
			case *map[domain.InstrumentType]int:
				(*target)[domain.InstrumentType(key)] = count
			case *map[domain.ComponentType]int:
				(*target)[domain.ComponentType(key)] = count
			case *map[string]int:
				(*target)[key] = count
			case *map[domain.SourceCategory]int:
				(*target)[domain.SourceCategory(key)] = count
			case *map[domain.ConfidenceLevel]int:
				(*target)[domain.ConfidenceLevel(key)] = count
			}
		}
		_ = rows.Close()
	}

	return facets, nil
}

func (s *SQLiteStore) GetFilteredFacetCounts(ctx context.Context, query FacetQuery) (*FacetCounts, error) {
	facets := &FacetCounts{
		InstrumentTypes:  make(map[domain.InstrumentType]int),
		ComponentTypes:   make(map[domain.ComponentType]int),
		ComponentNames:   make(map[string]int),
		SourceCategories: make(map[domain.SourceCategory]int),
		SourceNames:      make(map[string]int),
		ConfidenceLevels: make(map[domain.ConfidenceLevel]int),
		Units:            make(map[string]int),
	}

	whereClause := ""
	var args []interface{}
	if query.SourceName != "" {
		whereClause = " WHERE source_name = ?"
		args = append(args, query.SourceName)
	}

	queries := []struct {
		query  string
		target interface{}
	}{
		{"SELECT instrument_type, COUNT(*) FROM metrics" + whereClause + " GROUP BY instrument_type", &facets.InstrumentTypes},
		{"SELECT component_type, COUNT(*) FROM metrics" + whereClause + " GROUP BY component_type", &facets.ComponentTypes},
		{"SELECT component_name, COUNT(*) FROM metrics" + whereClause + " GROUP BY component_name", &facets.ComponentNames},
		{"SELECT source_category, COUNT(*) FROM metrics" + whereClause + " GROUP BY source_category", &facets.SourceCategories},
		{"SELECT source_name, COUNT(*) FROM metrics GROUP BY source_name", &facets.SourceNames}, // Always show all sources
		{"SELECT source_confidence, COUNT(*) FROM metrics" + whereClause + " GROUP BY source_confidence", &facets.ConfidenceLevels},
		{"SELECT unit, COUNT(*) FROM metrics WHERE unit IS NOT NULL AND unit != ''" + strings.Replace(whereClause, "WHERE", "AND", 1) + " GROUP BY unit", &facets.Units},
	}

	for _, q := range queries {
		var rows *sql.Rows
		var err error
		if strings.Contains(q.query, "?") {
			rows, err = s.db.QueryContext(ctx, q.query, args...)
		} else {
			rows, err = s.db.QueryContext(ctx, q.query)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to query facets: %w", err)
		}

		for rows.Next() {
			var key string
			var count int
			if err := rows.Scan(&key, &count); err != nil {
				_ = rows.Close()
				return nil, fmt.Errorf("failed to scan facet: %w", err)
			}

			switch target := q.target.(type) {
			case *map[domain.InstrumentType]int:
				(*target)[domain.InstrumentType(key)] = count
			case *map[domain.ComponentType]int:
				(*target)[domain.ComponentType(key)] = count
			case *map[string]int:
				(*target)[key] = count
			case *map[domain.SourceCategory]int:
				(*target)[domain.SourceCategory(key)] = count
			case *map[domain.ConfidenceLevel]int:
				(*target)[domain.ConfidenceLevel(key)] = count
			}
		}
		_ = rows.Close()
	}

	return facets, nil
}

func (s *SQLiteStore) CreateExtractionRun(ctx context.Context, run *ExtractionRun) error {
	query := `
		INSERT INTO extraction_runs (id, adapter_name, "commit", started_at, status)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := s.db.ExecContext(ctx, query, run.ID, run.AdapterName, run.Commit, run.StartedAt, run.Status)
	if err != nil {
		return fmt.Errorf("failed to create extraction run: %w", err)
	}
	return nil
}

func (s *SQLiteStore) UpdateExtractionRun(ctx context.Context, run *ExtractionRun) error {
	query := `
		UPDATE extraction_runs
		SET completed_at = ?, metrics_count = ?, status = ?, error_message = ?
		WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query, run.CompletedAt, run.MetricsCount, run.Status, run.ErrorMessage, run.ID)
	if err != nil {
		return fmt.Errorf("failed to update extraction run: %w", err)
	}
	return nil
}

func (s *SQLiteStore) GetExtractionRun(ctx context.Context, id string) (*ExtractionRun, error) {
	query := `
		SELECT id, adapter_name, "commit", started_at, completed_at, metrics_count, status, error_message
		FROM extraction_runs WHERE id = ?
	`

	var run ExtractionRun
	var completedAt sql.NullTime
	var commit, errorMessage sql.NullString

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&run.ID, &run.AdapterName, &commit, &run.StartedAt, &completedAt, &run.MetricsCount, &run.Status, &errorMessage,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get extraction run: %w", err)
	}

	run.Commit = commit.String
	run.ErrorMessage = errorMessage.String
	if completedAt.Valid {
		run.CompletedAt = &completedAt.Time
	}

	return &run, nil
}

func (s *SQLiteStore) GetLatestExtractionRun(ctx context.Context, adapterName string) (*ExtractionRun, error) {
	query := `
		SELECT id, adapter_name, "commit", started_at, completed_at, metrics_count, status, error_message
		FROM extraction_runs WHERE adapter_name = ?
		ORDER BY started_at DESC LIMIT 1
	`

	var run ExtractionRun
	var completedAt sql.NullTime
	var commit, errorMessage sql.NullString

	err := s.db.QueryRowContext(ctx, query, adapterName).Scan(
		&run.ID, &run.AdapterName, &commit, &run.StartedAt, &completedAt, &run.MetricsCount, &run.Status, &errorMessage,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest extraction run: %w", err)
	}

	run.Commit = commit.String
	run.ErrorMessage = errorMessage.String
	if completedAt.Valid {
		run.CompletedAt = &completedAt.Time
	}

	return &run, nil
}

// RunMigrations creates the database schema
func (s *SQLiteStore) RunMigrations(ctx context.Context) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS metrics (
			id                  TEXT PRIMARY KEY,
			metric_name         TEXT NOT NULL,
			instrument_type     TEXT NOT NULL,
			description         TEXT,
			unit                TEXT,
			enabled_by_default  INTEGER DEFAULT 1,
			component_type      TEXT NOT NULL,
			component_name      TEXT NOT NULL,
			source_category     TEXT NOT NULL,
			source_name         TEXT NOT NULL,
			source_location     TEXT,
			extraction_method   TEXT NOT NULL,
			source_confidence   TEXT NOT NULL,
			repo                TEXT,
			path                TEXT,
			"commit"            TEXT,
			extracted_at        TIMESTAMP NOT NULL,
			created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_metric_name ON metrics(metric_name)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_instrument_type ON metrics(instrument_type)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_component_type ON metrics(component_type)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_component_name ON metrics(component_name)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_source_category ON metrics(source_category)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_source_name ON metrics(source_name)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_source_confidence ON metrics(source_confidence)`,
		`CREATE TABLE IF NOT EXISTS metric_attributes (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			metric_id       TEXT NOT NULL REFERENCES metrics(id) ON DELETE CASCADE,
			attribute_name  TEXT NOT NULL,
			attribute_type  TEXT,
			description     TEXT,
			required        INTEGER DEFAULT 0,
			UNIQUE(metric_id, attribute_name)
		)`,
		`CREATE TABLE IF NOT EXISTS attribute_enum_values (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			attribute_id    INTEGER NOT NULL REFERENCES metric_attributes(id) ON DELETE CASCADE,
			enum_value      TEXT NOT NULL,
			UNIQUE(attribute_id, enum_value)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_metric_attributes_metric_id ON metric_attributes(metric_id)`,
		`CREATE INDEX IF NOT EXISTS idx_metric_attributes_name ON metric_attributes(attribute_name)`,
		`CREATE INDEX IF NOT EXISTS idx_attribute_enum_values_attr_id ON attribute_enum_values(attribute_id)`,
		`CREATE VIRTUAL TABLE IF NOT EXISTS metrics_fts USING fts5(
			metric_name,
			description,
			component_name,
			source_name,
			content='metrics',
			content_rowid='rowid'
		)`,
		`CREATE TABLE IF NOT EXISTS extraction_runs (
			id              TEXT PRIMARY KEY,
			adapter_name    TEXT NOT NULL,
			"commit"        TEXT,
			started_at      TIMESTAMP NOT NULL,
			completed_at    TIMESTAMP,
			metrics_count   INTEGER DEFAULT 0,
			status          TEXT NOT NULL DEFAULT 'running',
			error_message   TEXT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_extraction_runs_adapter ON extraction_runs(adapter_name)`,
		`CREATE INDEX IF NOT EXISTS idx_extraction_runs_status ON extraction_runs(status)`,
		`CREATE INDEX IF NOT EXISTS idx_extraction_runs_started_at ON extraction_runs(started_at)`,
	}

	for _, migration := range migrations {
		if _, err := s.db.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}
	}

	// Create triggers for FTS sync (need to check if they exist first)
	triggers := []struct {
		name string
		sql  string
	}{
		{
			"metrics_ai",
			`CREATE TRIGGER metrics_ai AFTER INSERT ON metrics BEGIN
				INSERT INTO metrics_fts(rowid, metric_name, description, component_name, source_name)
				VALUES (NEW.rowid, NEW.metric_name, NEW.description, NEW.component_name, NEW.source_name);
			END`,
		},
		{
			"metrics_ad",
			`CREATE TRIGGER metrics_ad AFTER DELETE ON metrics BEGIN
				INSERT INTO metrics_fts(metrics_fts, rowid, metric_name, description, component_name, source_name)
				VALUES ('delete', OLD.rowid, OLD.metric_name, OLD.description, OLD.component_name, OLD.source_name);
			END`,
		},
		{
			"metrics_au",
			`CREATE TRIGGER metrics_au AFTER UPDATE ON metrics BEGIN
				INSERT INTO metrics_fts(metrics_fts, rowid, metric_name, description, component_name, source_name)
				VALUES ('delete', OLD.rowid, OLD.metric_name, OLD.description, OLD.component_name, OLD.source_name);
				INSERT INTO metrics_fts(rowid, metric_name, description, component_name, source_name)
				VALUES (NEW.rowid, NEW.metric_name, NEW.description, NEW.component_name, NEW.source_name);
			END`,
		},
	}

	for _, trigger := range triggers {
		var exists int
		err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM sqlite_master WHERE type='trigger' AND name=?", trigger.name).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check trigger existence: %w", err)
		}
		if exists == 0 {
			if _, err := s.db.ExecContext(ctx, trigger.sql); err != nil {
				return fmt.Errorf("failed to create trigger %s: %w", trigger.name, err)
			}
		}
	}

	return nil
}
