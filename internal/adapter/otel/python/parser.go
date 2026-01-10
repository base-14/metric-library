package python

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
	"create_counter":                    "counter",
	"create_up_down_counter":            "updowncounter",
	"create_histogram":                  "histogram",
	"create_gauge":                      "gauge",
	"create_observable_counter":         "counter",
	"create_observable_up_down_counter": "updowncounter",
	"create_observable_gauge":           "gauge",
}

var (
	meterCreatePattern = regexp.MustCompile(`(?:meter|self\._meter|self\.meter)\.(create_\w+)\s*\(`)
	namePattern        = regexp.MustCompile(`name\s*=\s*["']([^"']+)["']`)
	nameConstPattern   = regexp.MustCompile(`name\s*=\s*(\w+(?:\.\w+)*)`)
	unitPattern        = regexp.MustCompile(`unit\s*=\s*["']([^"']+)["']`)
	descPattern        = regexp.MustCompile(`description\s*=\s*["']([^"']+)["']`)
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

		name := extractStringArg(callContent, namePattern)
		if name == "" {
			constMatch := nameConstPattern.FindStringSubmatch(callContent)
			if len(constMatch) > 1 {
				name = resolveConstant(constMatch[1], content)
			}
		}

		if name == "" {
			continue
		}

		unit := extractStringArg(callContent, unitPattern)
		description := extractStringArg(callContent, descPattern)

		metrics = append(metrics, &MetricDef{
			Name:           name,
			InstrumentType: instrumentType,
			Unit:           unit,
			Description:    description,
		})
	}

	return metrics, nil
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

		if !inString && (c == '"' || c == '\'') {
			if i+2 < len(content) && content[i:i+3] == string(c)+string(c)+string(c) {
				inString = true
				stringChar = c
				i += 2
				continue
			}
			inString = true
			stringChar = c
			continue
		}

		if inString && c == stringChar {
			if i+2 < len(content) && content[i:i+3] == string(c)+string(c)+string(c) {
				inString = false
				i += 2
				continue
			}
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

func extractStringArg(content string, pattern *regexp.Regexp) string {
	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func resolveConstant(constName string, content string) string {
	parts := strings.Split(constName, ".")
	lookupName := parts[len(parts)-1]

	pattern := regexp.MustCompile(lookupName + `\s*(?::\s*str)?\s*=\s*["']([^"']+)["']`)
	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}

	return ""
}
