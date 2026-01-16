package dotnet

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
	"CreateCounter":                 "counter",
	"CreateUpDownCounter":           "updowncounter",
	"CreateHistogram":               "histogram",
	"CreateGauge":                   "gauge",
	"CreateObservableCounter":       "counter",
	"CreateObservableUpDownCounter": "updowncounter",
	"CreateObservableGauge":         "gauge",
}

var (
	// Match Meter.CreateXxx<T>("name" or variations like MeterInstance, _meter, this, etc.
	// Captures: 1=method name, 2=metric name (handles multiline)
	meterCreatePattern = regexp.MustCompile(`(?:(?:\w*)[mM]eter(?:Instance)?!?|this)\s*\.\s*(Create(?:Observable)?(?:Counter|UpDownCounter|Histogram|Gauge))(?:<[^>]+>)?\s*\(\s*(?:name:\s*)?(?:\r?\n\s*)?"([^"]+)"`)

	// Match description parameter: description: "..."
	descriptionPattern = regexp.MustCompile(`description:\s*"([^"]+)"`)

	// Match unit parameter: unit: "..."
	unitPattern = regexp.MustCompile(`unit:\s*"([^"]+)"`)

	// Match constant definitions for metric names
	constPattern = regexp.MustCompile(`(?:const|static\s+readonly)\s+string\s+(\w+)\s*=\s*"([^"]+)"`)
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

	// First, build a map of constants
	constants := extractConstants(content)

	// Find all meter.CreateXxx calls
	matches := meterCreatePattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		methodName := content[match[2]:match[3]]
		metricName := content[match[4]:match[5]]

		// Resolve constant reference if needed
		if resolved, ok := constants[metricName]; ok {
			metricName = resolved
		}

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

func extractConstants(content string) map[string]string {
	constants := make(map[string]string)

	matches := constPattern.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) >= 3 {
			constants[match[1]] = match[2]
		}
	}

	return constants
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
			case ';':
				// Statement end without matching parens
				return i
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
