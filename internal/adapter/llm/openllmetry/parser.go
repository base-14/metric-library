package openllmetry

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type metricDef struct {
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
	// Positional argument patterns
	posStringPattern = regexp.MustCompile(`^\s*["']([^"']+)["']`)
	posConstPattern  = regexp.MustCompile(`^\s*(\w+(?:\.\w+)+)`)
	posUnitPattern   = regexp.MustCompile(`^\s*,\s*["']([^"']+)["']`)
	posDescPattern   = regexp.MustCompile(`^\s*,\s*["']([^"']+)["']\s*,\s*["']([^"']+)["']`)
)

func parseFileWithConstants(path string, constants map[string]string) ([]*metricDef, error) {
	cleanPath := filepath.Clean(path)
	content, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, err
	}

	return parseContentWithConstants(string(content), constants)
}

func parseContentWithConstants(content string, constants map[string]string) ([]*metricDef, error) {
	var metrics []*metricDef

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

		// Extract the arguments portion (content inside parentheses)
		parenStart := strings.Index(callContent, "(")
		argsContent := ""
		if parenStart >= 0 {
			argsContent = callContent[parenStart+1:]
		}

		var name, unit, description string

		// Try named argument pattern first: name="..."
		name = extractStringArg(callContent, namePattern)
		if name == "" {
			// Try named constant pattern: name=Meters.CONSTANT
			constMatch := nameConstPattern.FindStringSubmatch(callContent)
			if len(constMatch) > 1 {
				name = resolveConstantWithMap(constMatch[1], constants, content)
			}
		}

		// If still no name, try positional arguments
		if name == "" && argsContent != "" {
			name, unit, description = extractPositionalArgs(argsContent, constants, content)
		}

		if name == "" {
			continue
		}

		// If we didn't get unit/description from positional args, try named patterns
		if unit == "" {
			unit = extractStringArg(callContent, unitPattern)
		}
		if description == "" {
			description = extractStringArg(callContent, descPattern)
		}

		metrics = append(metrics, &metricDef{
			Name:           name,
			InstrumentType: instrumentType,
			Unit:           unit,
			Description:    description,
		})
	}

	return metrics, nil
}

func extractPositionalArgs(argsContent string, constants map[string]string, fileContent string) (name, unit, description string) {
	// Try to extract first positional argument as a constant (e.g., Meters.CONSTANT)
	constMatch := posConstPattern.FindStringSubmatch(argsContent)
	if len(constMatch) > 1 {
		name = resolveConstantWithMap(constMatch[1], constants, fileContent)
		if name != "" {
			// After the constant, look for unit and description
			afterConst := argsContent[len(constMatch[0]):]
			unitDescMatch := posDescPattern.FindStringSubmatch(afterConst)
			if len(unitDescMatch) >= 3 {
				unit = unitDescMatch[1]
				description = unitDescMatch[2]
			} else {
				unitMatch := posUnitPattern.FindStringSubmatch(afterConst)
				if len(unitMatch) > 1 {
					unit = unitMatch[1]
				}
			}
			return
		}
	}

	// Try to extract first positional argument as a string literal
	stringMatch := posStringPattern.FindStringSubmatch(argsContent)
	if len(stringMatch) > 1 {
		name = stringMatch[1]
		afterName := argsContent[len(stringMatch[0]):]
		unitDescMatch := posDescPattern.FindStringSubmatch(afterName)
		if len(unitDescMatch) >= 3 {
			unit = unitDescMatch[1]
			description = unitDescMatch[2]
		} else {
			unitMatch := posUnitPattern.FindStringSubmatch(afterName)
			if len(unitMatch) > 1 {
				unit = unitMatch[1]
			}
		}
	}

	return
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

func resolveConstantWithMap(constName string, constants map[string]string, content string) string {
	// Handle Meters.CONSTANT_NAME or GenAIMetrics.CONSTANT_NAME patterns
	parts := strings.Split(constName, ".")
	lookupName := parts[len(parts)-1]

	// First try the provided constants map
	if constants != nil {
		if val, ok := constants[lookupName]; ok {
			return val
		}
	}

	// Fallback: try to find in the same file
	pattern := regexp.MustCompile(lookupName + `\s*(?::\s*str)?\s*=\s*["']([^"']+)["']`)
	match := pattern.FindStringSubmatch(content)
	if len(match) > 1 {
		return match[1]
	}

	return ""
}
