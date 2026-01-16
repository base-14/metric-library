package golang

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

var methodToType = map[string]string{
	"Int64Counter":                   "counter",
	"Int64UpDownCounter":             "updowncounter",
	"Int64Histogram":                 "histogram",
	"Int64Gauge":                     "gauge",
	"Int64ObservableCounter":         "counter",
	"Int64ObservableUpDownCounter":   "updowncounter",
	"Int64ObservableGauge":           "gauge",
	"Float64Counter":                 "counter",
	"Float64UpDownCounter":           "updowncounter",
	"Float64Histogram":               "histogram",
	"Float64Gauge":                   "gauge",
	"Float64ObservableCounter":       "counter",
	"Float64ObservableUpDownCounter": "updowncounter",
	"Float64ObservableGauge":         "gauge",
}

var (
	// Match meter.Int64Counter("name", ...) or meter.Float64Histogram("name", ...) etc.
	// Captures: 1=method name, 2=metric name
	meterCreatePattern = regexp.MustCompile(`(?:\w+\.)?meter\s*\.\s*((?:Int64|Float64)(?:Observable)?(?:Counter|UpDownCounter|Histogram|Gauge))\s*\(\s*\n?\s*"([^"]+)"`)

	// Match metric.WithDescription("...")
	descriptionPattern = regexp.MustCompile(`metric\.WithDescription\s*\(\s*"([^"]+)"`)

	// Match metric.WithUnit("...")
	unitPattern = regexp.MustCompile(`metric\.WithUnit\s*\(\s*"([^"]+)"`)
)

func ParseFile(path string) ([]*MetricDef, error) {
	cleanPath := filepath.Clean(path)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	return parseContent(string(content))
}

func parseContent(content string) ([]*MetricDef, error) {
	var metrics []*MetricDef

	// Find all meter.Create* calls
	matches := meterCreatePattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		methodName := content[match[2]:match[3]]
		metricName := content[match[4]:match[5]]

		instrumentType := methodToType[methodName]
		if instrumentType == "" {
			continue
		}

		// Find the end of the method call to extract unit/description
		callEnd := findCallEnd(content, match[5])
		if callEnd == -1 {
			callEnd = min(match[5]+500, len(content))
		}

		callContent := content[match[0]:callEnd]

		// Extract description
		description := extractStringFromPattern(callContent, descriptionPattern)

		// Extract unit
		unit := extractStringFromPattern(callContent, unitPattern)

		metrics = append(metrics, &MetricDef{
			Name:           metricName,
			InstrumentType: instrumentType,
			Unit:           unit,
			Description:    description,
		})
	}

	return metrics, nil
}

func findCallEnd(content string, start int) int {
	depth := 1
	inString := false
	escaped := false

	for i := start; i < len(content); i++ {
		c := content[i]

		if escaped {
			escaped = false
			continue
		}

		if c == '\\' && inString {
			escaped = true
			continue
		}

		if c == '"' && !inString {
			inString = true
			continue
		}

		if c == '"' && inString {
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

func extractStringFromPattern(content string, pattern *regexp.Regexp) string {
	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
