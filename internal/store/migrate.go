package store

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func (s *SQLiteStore) Migrate(ctx context.Context) error {
	if err := s.createMigrationsTable(ctx); err != nil {
		return err
	}

	applied, err := s.getAppliedMigrations(ctx)
	if err != nil {
		return err
	}

	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	for _, file := range files {
		name := strings.TrimSuffix(file, ".sql")
		if applied[name] {
			continue
		}

		content, err := migrationsFS.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		upSQL := extractUpMigration(string(content))
		if upSQL == "" {
			continue
		}

		if err := s.runMigration(ctx, name, upSQL); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", name, err)
		}
	}

	return nil
}

func (s *SQLiteStore) createMigrationsTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func (s *SQLiteStore) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}
	return applied, rows.Err()
}

func (s *SQLiteStore) runMigration(ctx context.Context, name, sql string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	statements := splitStatements(sql)
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("statement failed: %w\nSQL: %s", err, stmt)
		}
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", name); err != nil {
		return err
	}

	return tx.Commit()
}

func extractUpMigration(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inUp := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- migrate:up") {
			inUp = true
			continue
		}
		if strings.HasPrefix(trimmed, "-- migrate:down") {
			break
		}
		if inUp {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func splitStatements(sql string) []string {
	var statements []string
	var current strings.Builder
	inTrigger := false

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		upperTrimmed := strings.ToUpper(trimmed)

		if strings.HasPrefix(upperTrimmed, "CREATE TRIGGER") {
			inTrigger = true
		}

		current.WriteString(line)
		current.WriteString("\n")

		if inTrigger && strings.HasSuffix(trimmed, "END;") {
			statements = append(statements, current.String())
			current.Reset()
			inTrigger = false
		} else if !inTrigger && strings.HasSuffix(trimmed, ";") {
			statements = append(statements, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func NewSQLiteStoreWithMigrations(dbPath string) (*SQLiteStore, error) {
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		return nil, err
	}

	if err := store.Migrate(context.Background()); err != nil {
		_ = store.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

var _ interface{ Migrate(context.Context) error } = (*SQLiteStore)(nil)
