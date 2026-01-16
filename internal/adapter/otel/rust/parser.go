package rust

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
	"u64_counter":                    "counter",
	"i64_counter":                    "counter",
	"f64_counter":                    "counter",
	"u64_up_down_counter":            "updowncounter",
	"i64_up_down_counter":            "updowncounter",
	"f64_up_down_counter":            "updowncounter",
	"u64_histogram":                  "histogram",
	"i64_histogram":                  "histogram",
	"f64_histogram":                  "histogram",
	"u64_gauge":                      "gauge",
	"i64_gauge":                      "gauge",
	"f64_gauge":                      "gauge",
	"u64_observable_counter":         "counter",
	"i64_observable_counter":         "counter",
	"f64_observable_counter":         "counter",
	"u64_observable_up_down_counter": "updowncounter",
	"i64_observable_up_down_counter": "updowncounter",
	"f64_observable_up_down_counter": "updowncounter",
	"u64_observable_gauge":           "gauge",
	"i64_observable_gauge":           "gauge",
	"f64_observable_gauge":           "gauge",
}

var (
	// Match meter.xxx_counter("name") or meter.xxx_histogram(CONSTANT_NAME)
	// Captures: 1=method name, 2=metric name (string or constant)
	meterCreatePattern = regexp.MustCompile(`(?:\w+\.)?meter\s*\.\s*([uif]64_(?:observable_)?(?:counter|up_down_counter|histogram|gauge))\s*\(\s*(?:Cow::from\()?(?:"([^"]+)"|([A-Z_][A-Z0-9_]*))\)?`)

	// Match .with_description("...")
	descriptionPattern = regexp.MustCompile(`\.with_description\s*\(\s*"([^"]+)"`)

	// Match .with_unit("...")
	unitPattern = regexp.MustCompile(`\.with_unit\s*\(\s*"([^"]+)"`)

	// Match const CONSTANT_NAME: &str = "value"
	constPattern = regexp.MustCompile(`const\s+([A-Z_][A-Z0-9_]*)\s*:\s*&str\s*=\s*"([^"]+)"`)
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

	// Find all meter.create* calls
	matches := meterCreatePattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		if len(match) < 8 {
			continue
		}

		methodName := content[match[2]:match[3]]

		// Group 2 is quoted string, Group 3 is constant name
		var metricName string
		if match[4] != -1 && match[5] != -1 {
			// Quoted string
			metricName = content[match[4]:match[5]]
		} else if match[6] != -1 && match[7] != -1 {
			// Constant name - resolve it
			constName := content[match[6]:match[7]]
			if resolved, ok := constants[constName]; ok {
				metricName = resolved
			} else {
				// Skip unresolved constants
				continue
			}
		} else {
			continue
		}

		instrumentType := methodToType[methodName]
		if instrumentType == "" {
			continue
		}

		// Find the end of the builder chain to extract unit/description
		// Use end of full match as starting point
		chainEnd := findChainEnd(content, match[1])
		if chainEnd == -1 {
			chainEnd = min(match[1]+500, len(content))
		}

		chainContent := content[match[0]:chainEnd]

		// Extract description
		description := extractStringFromPattern(chainContent, descriptionPattern)

		// Extract unit
		unit := extractStringFromPattern(chainContent, unitPattern)

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

func findChainEnd(content string, start int) int {
	maxEnd := min(start+500, len(content))

	// Look for .build() call
	for i := start; i < maxEnd; i++ {
		if i+7 <= len(content) && content[i:i+7] == ".build()" {
			return i + 7
		}
		// Stop at semicolon or new let/const statement
		if content[i] == ';' {
			return i
		}
		if i+4 <= len(content) && content[i:i+4] == "let " {
			return i
		}
		if i+6 <= len(content) && content[i:i+6] == "const " {
			return i
		}
	}

	return maxEnd
}

func extractStringFromPattern(content string, pattern *regexp.Regexp) string {
	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
