package redis

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/base14/otel-glossary/internal/adapter"
	"github.com/base14/otel-glossary/internal/domain"
	"github.com/base14/otel-glossary/internal/fetcher"
)

const repoURL = "https://github.com/oliver006/redis_exporter"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "prometheus-redis"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourcePrometheus
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceDerived
}

func (a *Adapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionAST
}

func (a *Adapter) RepoURL() string {
	return repoURL
}

func (a *Adapter) Fetch(ctx context.Context, opts adapter.FetchOptions) (*adapter.FetchResult, error) {
	fetchOpts := fetcher.FetchOptions{
		RepoURL: repoURL,
		Commit:  opts.Commit,
		Shallow: true,
		Depth:   1,
		Force:   opts.Force,
	}

	result, err := a.fetcher.Fetch(ctx, fetchOpts)
	if err != nil {
		return nil, err
	}

	return &adapter.FetchResult{
		RepoPath:  result.RepoPath,
		Commit:    result.Commit,
		Timestamp: result.Timestamp,
	}, nil
}

func (a *Adapter) Extract(ctx context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	exporterDir := filepath.Join(result.RepoPath, "exporter")

	entries, err := os.ReadDir(exporterDir)
	if err != nil {
		return nil, err
	}

	var metrics []*adapter.RawMetric

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		if strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}

		filePath := filepath.Join(exporterDir, entry.Name())
		fileMetrics, err := parseRedisFile(filePath)
		if err != nil {
			continue
		}

		for _, m := range fileMetrics {
			rawMetric := &adapter.RawMetric{
				Name:             "redis_" + m.Name,
				Description:      m.Description,
				InstrumentType:   inferInstrumentType(m.Name, m.MetricType),
				Attributes:       labelsToAttributes(m.Labels),
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    "redis",
				SourceLocation:   filePath,
				Path:             filePath,
			}
			metrics = append(metrics, rawMetric)
		}
	}

	return metrics, nil
}

type redisMetric struct {
	Name        string
	Description string
	Labels      []string
	MetricType  string // "gauge", "counter", or ""
}

func parseRedisFile(filePath string) ([]redisMetric, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var metrics []redisMetric

	ast.Inspect(node, func(n ast.Node) bool {
		// Look for key-value expressions in composite literals (struct field assignments)
		// Pattern: metricMapGauges: map[string]string{...}
		if kv, ok := n.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				metricType := ""
				switch ident.Name {
				case "metricMapGauges":
					metricType = "gauge"
				case "metricMapCounters":
					metricType = "counter"
				default:
					return true
				}

				if compLit, ok := kv.Value.(*ast.CompositeLit); ok {
					for _, elt := range compLit.Elts {
						if mapKV, ok := elt.(*ast.KeyValueExpr); ok {
							if keyLit, ok := mapKV.Key.(*ast.BasicLit); ok && keyLit.Kind == token.STRING {
								metricName := strings.Trim(keyLit.Value, `"`)
								m := redisMetric{
									Name:        metricName,
									Description: "Redis " + metricName,
									MetricType:  metricType,
								}
								metrics = append(metrics, m)
							}
						}
					}
				}
			}
		}

		// Look for range statements over inline map literals for metricDescriptions
		// Pattern: for k, desc := range map[string]struct{txt string; lbls []string}{...}
		if rangeStmt, ok := n.(*ast.RangeStmt); ok {
			if compLit, ok := rangeStmt.X.(*ast.CompositeLit); ok {
				if mapType, ok := compLit.Type.(*ast.MapType); ok {
					// Check if this looks like the metricDescriptions pattern
					if structType, ok := mapType.Value.(*ast.StructType); ok {
						if hasDescriptionFields(structType) {
							for _, elt := range compLit.Elts {
								if mapKV, ok := elt.(*ast.KeyValueExpr); ok {
									if keyLit, ok := mapKV.Key.(*ast.BasicLit); ok && keyLit.Kind == token.STRING {
										metricName := strings.Trim(keyLit.Value, `"`)
										m := parseMetricDescription(metricName, mapKV.Value)
										if m.Name != "" {
											metrics = append(metrics, m)
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// Also look for variable declarations (for tests and potential global vars)
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				varName := name.Name

				metricType := ""
				switch varName {
				case "metricMapGauges":
					metricType = "gauge"
				case "metricMapCounters":
					metricType = "counter"
				case "metricDescriptions":
					metricType = ""
				default:
					continue
				}

				if i >= len(valueSpec.Values) {
					continue
				}

				compLit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					continue
				}

				for _, elt := range compLit.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}

					keyLit, ok := kv.Key.(*ast.BasicLit)
					if !ok || keyLit.Kind != token.STRING {
						continue
					}
					metricName := strings.Trim(keyLit.Value, `"`)

					if varName == "metricDescriptions" {
						m := parseMetricDescription(metricName, kv.Value)
						if m.Name != "" {
							metrics = append(metrics, m)
						}
					} else {
						m := redisMetric{
							Name:        metricName,
							Description: "Redis " + metricName,
							MetricType:  metricType,
						}
						metrics = append(metrics, m)
					}
				}
			}
		}

		return true
	})

	return metrics, nil
}

func hasDescriptionFields(st *ast.StructType) bool {
	hasTxt := false
	hasLbls := false
	for _, field := range st.Fields.List {
		for _, name := range field.Names {
			if name.Name == "txt" {
				hasTxt = true
			}
			if name.Name == "lbls" {
				hasLbls = true
			}
		}
	}
	return hasTxt && hasLbls
}

func parseMetricDescription(name string, value ast.Expr) redisMetric {
	m := redisMetric{Name: name}

	compLit, ok := value.(*ast.CompositeLit)
	if !ok {
		return m
	}

	for _, elt := range compLit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		keyIdent, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch keyIdent.Name {
		case "txt":
			if lit, ok := kv.Value.(*ast.BasicLit); ok && lit.Kind == token.STRING {
				m.Description = strings.Trim(lit.Value, "`\"")
			}
		case "lbls":
			if comp, ok := kv.Value.(*ast.CompositeLit); ok {
				for _, lblElt := range comp.Elts {
					if lit, ok := lblElt.(*ast.BasicLit); ok && lit.Kind == token.STRING {
						m.Labels = append(m.Labels, strings.Trim(lit.Value, `"`))
					}
				}
			}
		}
	}

	return m
}

func inferInstrumentType(metricName, metricType string) string {
	if metricType != "" {
		return metricType
	}

	switch {
	case strings.HasSuffix(metricName, "_total"):
		return string(domain.InstrumentCounter)
	case strings.HasSuffix(metricName, "_bucket"):
		return string(domain.InstrumentHistogram)
	case strings.HasSuffix(metricName, "_seconds") && !strings.Contains(metricName, "duration"):
		return string(domain.InstrumentGauge)
	default:
		return string(domain.InstrumentGauge)
	}
}

func labelsToAttributes(labels []string) []domain.Attribute {
	attrs := make([]domain.Attribute, 0, len(labels))
	for _, label := range labels {
		attrs = append(attrs, domain.Attribute{
			Name: label,
			Type: "string",
		})
	}
	return attrs
}
