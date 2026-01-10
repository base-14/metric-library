package ksm

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type MetricDefinition struct {
	Name       string
	Help       string
	MetricType string
}

func ParseFile(filePath string) ([]MetricDefinition, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var defs []MetricDefinition

	ast.Inspect(file, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		if !isNewFamilyGeneratorWithStability(call) {
			return true
		}

		def := extractMetricDefinition(call)
		if def != nil {
			defs = append(defs, *def)
		}

		return true
	})

	return defs, nil
}

func isNewFamilyGeneratorWithStability(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	if sel.Sel.Name != "NewFamilyGeneratorWithStability" {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "generator"
}

func extractMetricDefinition(call *ast.CallExpr) *MetricDefinition {
	if len(call.Args) < 3 {
		return nil
	}

	name := extractStringLiteral(call.Args[0])
	if name == "" {
		return nil
	}

	help := extractStringLiteral(call.Args[1])

	metricType := extractMetricType(call.Args[2])

	return &MetricDefinition{
		Name:       name,
		Help:       help,
		MetricType: metricType,
	}
}

func extractStringLiteral(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, "\"")
}

func extractMetricType(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "gauge"
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "metric" {
		return "gauge"
	}

	switch sel.Sel.Name {
	case "Counter":
		return "counter"
	case "Gauge":
		return "gauge"
	case "Histogram":
		return "histogram"
	case "Summary":
		return "summary"
	default:
		return "gauge"
	}
}
