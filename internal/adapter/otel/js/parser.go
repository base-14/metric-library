package js

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MetricDef struct {
	Name           string
	InstrumentType string
	Unit           string
	Description    string
}

var (
	// Match METRIC_* exports with their values
	// export const METRIC_SYSTEM_CPU_TIME = 'system.cpu.time' as const;
	metricExportPattern = regexp.MustCompile(`export\s+const\s+(METRIC_\w+)\s*=\s*['"]([^'"]+)['"]\s*as\s+const`)

	// Match meter.create* calls for additional metrics
	// this._meter.createObservableCounter(METRIC_NAME, { description: '...', unit: '...' })
	// this._meter.createHistogram('metric.name', { ... })
	meterCreatePattern = regexp.MustCompile(`(?:this\._?meter|meter)\.(create\w+)\s*\(`)

	// Match descriptions in options objects
	descriptionPattern = regexp.MustCompile(`description:\s*['"]([^'"]+)['"]`)
	unitPattern        = regexp.MustCompile(`unit:\s*['"]([^'"]+)['"]`)
)

var methodToType = map[string]string{
	"createCounter":                 "counter",
	"createUpDownCounter":           "updowncounter",
	"createHistogram":               "histogram",
	"createGauge":                   "gauge",
	"createObservableCounter":       "counter",
	"createObservableUpDownCounter": "updowncounter",
	"createObservableGauge":         "gauge",
}

func ParseSemconvFile(path string) ([]*MetricDef, error) {
	cleanPath := filepath.Clean(path)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	return parseSemconvContent(string(content))
}

func parseSemconvContent(content string) ([]*MetricDef, error) {
	var metrics []*MetricDef

	// Find all METRIC_* exports with their positions
	exportMatches := metricExportPattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range exportMatches {
		if len(match) < 6 {
			continue
		}

		// match[0]:match[1] is the full match
		// match[2]:match[3] is the constant name (METRIC_*)
		// match[4]:match[5] is the metric value
		exportStart := match[0]
		constName := content[match[2]:match[3]]
		metricName := content[match[4]:match[5]]

		// Look backward from export to find the immediately preceding JSDoc block
		description := findPrecedingJSDoc(content, exportStart)

		// Infer instrument type from metric name suffix
		instrumentType := inferInstrumentType(metricName)

		// Try to extract unit from metric name (if ends in .time, .duration, etc.)
		unit := inferUnit(metricName)

		metrics = append(metrics, &MetricDef{
			Name:           metricName,
			InstrumentType: instrumentType,
			Unit:           unit,
			Description:    description,
		})

		_ = constName // unused but kept for debugging
	}

	return metrics, nil
}

func findPrecedingJSDoc(content string, exportStart int) string {
	// Look backward from exportStart to find the end of a JSDoc block (*/)
	searchArea := content[:exportStart]

	// Find the last JSDoc block before this export
	// We need to find */ and then look backward for /**
	lastBlockEnd := strings.LastIndex(searchArea, "*/")
	if lastBlockEnd == -1 {
		return ""
	}

	// Check if there's only whitespace between */ and the export
	between := strings.TrimSpace(content[lastBlockEnd+2 : exportStart])
	if between != "" {
		// There's something between the JSDoc and the export, not the right block
		return ""
	}

	// Find the matching /** for this */
	blockStart := strings.LastIndex(searchArea[:lastBlockEnd], "/**")
	if blockStart == -1 {
		return ""
	}

	// Extract the JSDoc content (between /** and */)
	jsdocContent := content[blockStart+3 : lastBlockEnd]

	return extractJSDocDescription(jsdocContent)
}

func ParseInstrumentationFile(path string) ([]*MetricDef, error) {
	cleanPath := filepath.Clean(path)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	return parseInstrumentationContent(string(content), "")
}

func parseInstrumentationContent(content string, semconvContent string) ([]*MetricDef, error) {
	var metrics []*MetricDef

	// Build constant map from semconv content if provided
	constants := make(map[string]string)
	if semconvContent != "" {
		matches := metricExportPattern.FindAllStringSubmatch(semconvContent, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				constants[match[1]] = match[2]
			}
		}
	}

	// Also look for local constant definitions
	localMatches := metricExportPattern.FindAllStringSubmatch(content, -1)
	for _, match := range localMatches {
		if len(match) >= 3 {
			constants[match[1]] = match[2]
		}
	}

	matches := meterCreatePattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		methodName := content[match[2]:match[3]]
		instrumentType, ok := methodToType[methodName]
		if !ok {
			continue
		}

		callStart := match[0]
		callEnd := findMatchingParen(content, match[1]-1)
		if callEnd == -1 {
			continue
		}

		callContent := content[callStart:callEnd]

		// Extract metric name (constant reference or string literal)
		name := extractMetricName(callContent, constants)
		if name == "" {
			continue
		}

		// Extract description and unit from options object
		description := ""
		if descMatch := descriptionPattern.FindStringSubmatch(callContent); len(descMatch) > 1 {
			description = descMatch[1]
		}

		unit := ""
		if unitMatch := unitPattern.FindStringSubmatch(callContent); len(unitMatch) > 1 {
			unit = unitMatch[1]
		}

		metrics = append(metrics, &MetricDef{
			Name:           name,
			InstrumentType: instrumentType,
			Unit:           unit,
			Description:    description,
		})
	}

	return metrics, nil
}

func extractMetricName(callContent string, constants map[string]string) string {
	// Look for METRIC_* constant reference
	constRefPattern := regexp.MustCompile(`\(\s*(METRIC_\w+)`)
	if match := constRefPattern.FindStringSubmatch(callContent); len(match) > 1 {
		if resolved, ok := constants[match[1]]; ok {
			return resolved
		}
	}

	// Look for string literal
	stringPattern := regexp.MustCompile(`\(\s*['"]([^'"]+)['"]`)
	if match := stringPattern.FindStringSubmatch(callContent); len(match) > 1 {
		return match[1]
	}

	return ""
}

func extractJSDocDescription(jsdoc string) string {
	lines := strings.Split(jsdoc, "\n")
	var descParts []string

	for _, line := range lines {
		// Remove leading asterisks and whitespace
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)

		// Skip annotation lines (@example, @note, @experimental, etc.)
		if strings.HasPrefix(line, "@") {
			continue
		}

		// Skip empty lines
		if line == "" {
			continue
		}

		descParts = append(descParts, line)
	}

	return strings.Join(descParts, " ")
}

func inferInstrumentType(metricName string) string {
	// Look at metric name suffixes to infer type
	switch {
	case strings.HasSuffix(metricName, ".time") || strings.HasSuffix(metricName, ".duration"):
		return "counter" // Usually cumulative time
	case strings.HasSuffix(metricName, ".count") || strings.HasSuffix(metricName, ".total"):
		return "counter"
	case strings.HasSuffix(metricName, ".errors"):
		return "counter"
	case strings.HasSuffix(metricName, ".io"):
		return "counter"
	case strings.HasSuffix(metricName, ".usage") || strings.HasSuffix(metricName, ".used"):
		return "gauge"
	case strings.HasSuffix(metricName, ".utilization"):
		return "gauge"
	case strings.HasSuffix(metricName, ".limit"):
		return "gauge"
	case strings.HasSuffix(metricName, ".size"):
		return "gauge"
	case strings.Contains(metricName, ".delay."):
		return "gauge"
	default:
		return "gauge"
	}
}

func inferUnit(metricName string) string {
	switch {
	case strings.HasSuffix(metricName, ".time") || strings.Contains(metricName, ".duration"):
		return "s"
	case strings.HasSuffix(metricName, ".usage") && strings.Contains(metricName, "memory"):
		return "By"
	case strings.HasSuffix(metricName, ".io"):
		return "By"
	case strings.Contains(metricName, ".delay."):
		return "s"
	default:
		return ""
	}
}

func findMatchingParen(content string, start int) int {
	depth := 1
	inString := false
	stringChar := byte(0)
	escaped := false

	for i := start + 1; i < len(content); i++ {
		c := content[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' && inString {
			escaped = true
			continue
		}

		if !inString && (c == '"' || c == '\'' || c == '`') {
			inString = true
			stringChar = c
			continue
		}

		if inString && c == stringChar {
			inString = false
			continue
		}

		if !inString {
			switch c {
			case '(':
				depth++
			case ')':
				depth--
				if depth == 0 {
					return i + 1
				}
			}
		}
	}

	return -1
}
