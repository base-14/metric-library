package java

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

var builderToType = map[string]string{
	"counterBuilder":       "counter",
	"histogramBuilder":     "histogram",
	"gaugeBuilder":         "gauge",
	"upDownCounterBuilder": "updowncounter",
}

var (
	// Match meter.xxxBuilder("name") patterns
	meterBuilderPattern = regexp.MustCompile(`meter\s*\.\s*(counter|histogram|gauge|upDownCounter)Builder\s*\(\s*"([^"]+)"`)
	// Match .setDescription("...") in method chain
	descriptionPattern = regexp.MustCompile(`\.setDescription\s*\(\s*"([^"]+)"`)
	// Match .setUnit("...") in method chain
	unitPattern = regexp.MustCompile(`\.setUnit\s*\(\s*"([^"]+)"`)
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

	// Find all meter builder calls
	matches := meterBuilderPattern.FindAllStringSubmatchIndex(content, -1)

	for _, match := range matches {
		if len(match) < 6 {
			continue
		}

		builderType := content[match[2]:match[3]]
		metricName := content[match[4]:match[5]]

		instrumentType := builderToType[builderType+"Builder"]
		if instrumentType == "" {
			continue
		}

		// Find the end of the method chain (look for .build() or semicolon)
		chainStart := match[0]
		chainEnd := findChainEnd(content, match[5])
		if chainEnd == -1 {
			chainEnd = len(content)
		}

		chainContent := content[chainStart:chainEnd]

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

func findChainEnd(content string, start int) int {
	// Look for .build(), .buildObserver(), .buildWithCallback(), or semicolon
	buildPattern := regexp.MustCompile(`\.build(?:Observer|WithCallback)?\s*\(`)

	// Search from start position
	remaining := content[start:]

	// Find the next .build*() call
	buildMatch := buildPattern.FindStringIndex(remaining)
	if buildMatch != nil {
		// Find the closing paren after build
		parenStart := start + buildMatch[1] - 1
		parenEnd := findMatchingParen(content, parenStart)
		if parenEnd != -1 {
			return parenEnd
		}
	}

	// Fallback: find semicolon
	for i := start; i < len(content); i++ {
		if content[i] == ';' {
			return i
		}
		// Stop at method or class boundary
		if i+6 < len(content) && (content[i:i+6] == "public" || content[i:i+7] == "private") {
			return i
		}
	}

	return -1
}

func findMatchingParen(content string, start int) int {
	if start >= len(content) || content[start] != '(' {
		return -1
	}

	depth := 1
	inString := false
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
