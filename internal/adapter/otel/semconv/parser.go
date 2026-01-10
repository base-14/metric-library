package semconv

import (
	"os"

	"gopkg.in/yaml.v3"
)

type MetricDefinition struct {
	Name       string
	Brief      string
	Instrument string
	Unit       string
	Stability  string
	Attributes []AttributeRef
}

type AttributeRef struct {
	Ref              string
	RequirementLevel string
}

type metricsFile struct {
	Groups []group `yaml:"groups"`
}

type group struct {
	ID         string        `yaml:"id"`
	Type       string        `yaml:"type"`
	MetricName string        `yaml:"metric_name"`
	Brief      string        `yaml:"brief"`
	Instrument string        `yaml:"instrument"`
	Unit       string        `yaml:"unit"`
	Stability  string        `yaml:"stability"`
	Attributes []attrRefYAML `yaml:"attributes"`
}

type attrRefYAML struct {
	Ref              string      `yaml:"ref"`
	RequirementLevel interface{} `yaml:"requirement_level"`
}

func ParseFile(filePath string) ([]MetricDefinition, error) {
	data, err := os.ReadFile(filePath) //nolint:gosec // filePath is from controlled source
	if err != nil {
		return nil, err
	}

	var file metricsFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, err
	}

	var defs []MetricDefinition

	for _, g := range file.Groups {
		if g.Type != "metric" {
			continue
		}

		def := MetricDefinition{
			Name:       g.MetricName,
			Brief:      g.Brief,
			Instrument: g.Instrument,
			Unit:       g.Unit,
			Stability:  g.Stability,
		}

		for _, attr := range g.Attributes {
			ref := AttributeRef{
				Ref: attr.Ref,
			}

			switch v := attr.RequirementLevel.(type) {
			case string:
				ref.RequirementLevel = v
			case map[string]interface{}:
				for key := range v {
					ref.RequirementLevel = key
					break
				}
			}

			def.Attributes = append(def.Attributes, ref)
		}

		defs = append(defs, def)
	}

	return defs, nil
}
