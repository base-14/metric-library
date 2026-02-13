package cockroachdb

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
	"github.com/base-14/metric-library/internal/fetcher"
)

const repoURL = "https://github.com/cockroachdb/cockroach"

type Adapter struct {
	fetcher *fetcher.GitFetcher
}

func NewAdapter(cacheDir string) *Adapter {
	return &Adapter{
		fetcher: fetcher.NewGitFetcher(cacheDir),
	}
}

func (a *Adapter) Name() string {
	return "prometheus-cockroachdb"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourcePrometheus
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceAuthoritative
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

func (a *Adapter) Extract(_ context.Context, result *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	pkgDir := filepath.Join(result.RepoPath, "pkg")

	seen := make(map[string]bool)
	var metrics []*adapter.RawMetric

	err := filepath.Walk(pkgDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		defs, parseErr := parseMetricMetadata(filepath.Base(path), readFile(path))
		if parseErr != nil {
			return nil
		}

		for _, def := range defs {
			if def.Name == "" || seen[def.Name] {
				continue
			}
			seen[def.Name] = true

			metrics = append(metrics, &adapter.RawMetric{
				Name:             def.Name,
				Description:      def.Help,
				Unit:             mapUnit(def.Unit),
				InstrumentType:   inferType(def.Name),
				EnabledByDefault: true,
				ComponentType:    string(domain.ComponentPlatform),
				ComponentName:    "cockroachdb",
				SourceLocation:   path,
				Path:             path,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path) //nolint:gosec // Reading Go source files from cloned repos is intentional
	if err != nil {
		return nil
	}
	return data
}

type metadataDef struct {
	Name string
	Help string
	Unit string
}

func parseMetricMetadata(filename string, src []byte) ([]metadataDef, error) {
	if src == nil {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, err
	}

	var defs []metadataDef

	ast.Inspect(f, func(n ast.Node) bool {
		comp, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		if !isMetricMetadataType(comp) {
			return true
		}

		def := extractMetadataFields(comp)
		if def.Name != "" {
			defs = append(defs, def)
		}

		return true
	})

	return defs, nil
}

func isMetricMetadataType(comp *ast.CompositeLit) bool {
	sel, ok := comp.Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "metric" && sel.Sel.Name == "Metadata"
}

func extractMetadataFields(comp *ast.CompositeLit) metadataDef {
	var def metadataDef

	for _, elt := range comp.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "Name":
			def.Name = extractStringValue(kv.Value)
		case "Help":
			def.Help = extractStringValue(kv.Value)
		case "Unit":
			def.Unit = extractUnitValue(kv.Value)
		}
	}

	return def
}

func extractStringValue(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, `"`)
}

func extractUnitValue(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	return sel.Sel.Name
}

func mapUnit(unit string) string {
	switch unit {
	case "Unit_BYTES":
		return "bytes"
	case "Unit_COUNT":
		return "count"
	case "Unit_NANOSECONDS", "Unit_TIMESTAMP_NS":
		return "nanoseconds"
	case "Unit_SECONDS", "Unit_TIMESTAMP_SEC":
		return "seconds"
	case "Unit_PERCENT":
		return "percent"
	default:
		return ""
	}
}

func inferType(name string) string {
	switch {
	case strings.HasSuffix(name, ".count"):
		return string(domain.InstrumentCounter)
	case strings.Contains(name, ".latency"):
		return string(domain.InstrumentHistogram)
	default:
		return string(domain.InstrumentGauge)
	}
}
