package cadvisor

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
	Labels     []string
}

func ParseFile(filePath string) ([]MetricDefinition, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var defs []MetricDefinition

	ast.Inspect(file, func(n ast.Node) bool {
		lit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		if !isMetricComposite(lit) {
			return true
		}

		def := extractMetricDefinition(lit)
		if def != nil {
			defs = append(defs, *def)
		}

		return true
	})

	return defs, nil
}

func isMetricComposite(lit *ast.CompositeLit) bool {
	if lit.Type == nil {
		return hasMetricFields(lit)
	}

	ident, ok := lit.Type.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "containerMetric" || ident.Name == "machineMetric"
}

func hasMetricFields(lit *ast.CompositeLit) bool {
	hasName := false
	hasValueType := false

	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		if key.Name == "name" {
			hasName = true
		}
		if key.Name == "valueType" {
			hasValueType = true
		}
	}

	return hasName && hasValueType
}

func extractMetricDefinition(lit *ast.CompositeLit) *MetricDefinition {
	def := &MetricDefinition{
		MetricType: "gauge",
	}

	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}

		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}

		switch key.Name {
		case "name":
			def.Name = extractStringLiteral(kv.Value)
		case "help":
			def.Help = extractStringLiteral(kv.Value)
		case "valueType":
			def.MetricType = extractValueType(kv.Value)
		case "extraLabels":
			def.Labels = extractStringSlice(kv.Value)
		}
	}

	if def.Name == "" {
		return nil
	}

	return def
}

func extractStringLiteral(expr ast.Expr) string {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return ""
	}
	return strings.Trim(lit.Value, "\"")
}

func extractValueType(expr ast.Expr) string {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return "gauge"
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok || ident.Name != "prometheus" {
		return "gauge"
	}

	switch sel.Sel.Name {
	case "CounterValue":
		return "counter"
	case "GaugeValue":
		return "gauge"
	case "UntypedValue":
		return "gauge"
	default:
		return "gauge"
	}
}

func extractStringSlice(expr ast.Expr) []string {
	comp, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	var result []string
	for _, elt := range comp.Elts {
		if s := extractStringLiteral(elt); s != "" {
			result = append(result, s)
		}
	}
	return result
}
