package astparser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type MetricDef struct {
	Name   string
	Help   string
	Labels []string
}

func ParseSource(filename string, src []byte) ([]MetricDef, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return extractMetrics(f)
}

func ParseFile(path string) ([]MetricDef, error) {
	src, err := os.ReadFile(path) //nolint:gosec // Reading Go source files from cloned repos is intentional
	if err != nil {
		return nil, err
	}
	return ParseSource(filepath.Base(path), src)
}

func extractMetrics(f *ast.File) ([]MetricDef, error) {
	var metrics []MetricDef
	constants := extractConstants(f)

	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isNewDescCall(call) {
			return true
		}

		if len(call.Args) < 2 {
			return true
		}

		name := extractMetricName(call.Args[0], constants)
		help := extractStringLiteral(call.Args[1])

		var labels []string
		if len(call.Args) >= 3 {
			labels = extractLabels(call.Args[2])
		}

		if name != "" {
			metrics = append(metrics, MetricDef{
				Name:   name,
				Help:   help,
				Labels: labels,
			})
		}

		return true
	})

	return metrics, nil
}

func extractConstants(f *ast.File) map[string]string {
	constants := make(map[string]string)

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if i < len(valueSpec.Values) {
					if lit, ok := valueSpec.Values[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
						constants[name.Name] = strings.Trim(lit.Value, `"`)
					}
				}
			}
		}
	}

	return constants
}

func isNewDescCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "prometheus" && sel.Sel.Name == "NewDesc"
}

func extractMetricName(arg ast.Expr, constants map[string]string) string {
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, `"`)
	}

	if call, ok := arg.(*ast.CallExpr); ok {
		if isBuildFQNameCall(call) {
			return buildFQName(call, constants)
		}
	}

	return ""
}

func isBuildFQNameCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "prometheus" && sel.Sel.Name == "BuildFQName"
}

func buildFQName(call *ast.CallExpr, constants map[string]string) string {
	if len(call.Args) != 3 {
		return ""
	}

	namespace := resolveStringArg(call.Args[0], constants)
	subsystem := resolveStringArg(call.Args[1], constants)
	name := resolveStringArg(call.Args[2], constants)

	parts := []string{}
	if namespace != "" {
		parts = append(parts, namespace)
	}
	if subsystem != "" {
		parts = append(parts, subsystem)
	}
	if name != "" {
		parts = append(parts, name)
	}

	return strings.Join(parts, "_")
}

func resolveStringArg(arg ast.Expr, constants map[string]string) string {
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, `"`)
	}

	if ident, ok := arg.(*ast.Ident); ok {
		if val, found := constants[ident.Name]; found {
			return val
		}
	}

	return ""
}

func extractStringLiteral(arg ast.Expr) string {
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		return strings.Trim(lit.Value, `"`)
	}
	return ""
}

func extractLabels(arg ast.Expr) []string {
	comp, ok := arg.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	var labels []string
	for _, elt := range comp.Elts {
		if lit, ok := elt.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			labels = append(labels, strings.Trim(lit.Value, `"`))
		}
	}

	return labels
}
