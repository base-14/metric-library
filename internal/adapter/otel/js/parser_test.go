package js

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSemconvContent(t *testing.T) {
	content := `
/**
 * Total CPU seconds broken down by different states.
 *
 * @experimental This metric is experimental.
 */
export const METRIC_SYSTEM_CPU_TIME = 'system.cpu.time' as const;

/**
 * Reports memory in use by state.
 *
 * @experimental This metric is experimental.
 */
export const METRIC_SYSTEM_MEMORY_USAGE = 'system.memory.usage' as const;

/**
 * Event loop maximum delay.
 *
 * @note Value can be retrieved from histogram.max.
 *
 * @experimental This metric is experimental.
 */
export const METRIC_NODEJS_EVENTLOOP_DELAY_MAX = 'nodejs.eventloop.delay.max' as const;
`

	metrics, err := parseSemconvContent(content)
	require.NoError(t, err)
	assert.Len(t, metrics, 3)

	// Check first metric
	assert.Equal(t, "system.cpu.time", metrics[0].Name)
	assert.Contains(t, metrics[0].Description, "Total CPU seconds")
	assert.Equal(t, "counter", metrics[0].InstrumentType)

	// Check second metric
	assert.Equal(t, "system.memory.usage", metrics[1].Name)
	assert.Contains(t, metrics[1].Description, "memory in use")
	assert.Equal(t, "gauge", metrics[1].InstrumentType)

	// Check third metric
	assert.Equal(t, "nodejs.eventloop.delay.max", metrics[2].Name)
	assert.Contains(t, metrics[2].Description, "Event loop maximum delay")
	assert.Equal(t, "gauge", metrics[2].InstrumentType)
}

func TestParseSemconvContent_WithUnit(t *testing.T) {
	content := `
/**
 * Garbage collection duration.
 *
 * @experimental This metric is experimental.
 */
export const METRIC_V8JS_GC_DURATION = 'v8js.gc.duration' as const;
`

	metrics, err := parseSemconvContent(content)
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	assert.Equal(t, "v8js.gc.duration", metrics[0].Name)
	assert.Equal(t, "s", metrics[0].Unit) // Inferred from .duration
}

func TestParseInstrumentationContent(t *testing.T) {
	semconvContent := `
export const METRIC_SYSTEM_CPU_TIME = 'system.cpu.time' as const;
export const METRIC_SYSTEM_MEMORY_USAGE = 'system.memory.usage' as const;
`

	instrumentationContent := `
this._meter.createObservableCounter(METRIC_SYSTEM_CPU_TIME, {
  description: 'Cpu time in seconds',
  unit: 's',
});

this._meter.createObservableGauge(METRIC_SYSTEM_MEMORY_USAGE, {
  description: 'Memory usage in bytes',
});
`

	metrics, err := parseInstrumentationContent(instrumentationContent, semconvContent)
	require.NoError(t, err)
	assert.Len(t, metrics, 2)

	// Check first metric
	assert.Equal(t, "system.cpu.time", metrics[0].Name)
	assert.Equal(t, "Cpu time in seconds", metrics[0].Description)
	assert.Equal(t, "s", metrics[0].Unit)
	assert.Equal(t, "counter", metrics[0].InstrumentType)

	// Check second metric
	assert.Equal(t, "system.memory.usage", metrics[1].Name)
	assert.Equal(t, "Memory usage in bytes", metrics[1].Description)
	assert.Equal(t, "gauge", metrics[1].InstrumentType)
}

func TestParseInstrumentationContent_StringLiteral(t *testing.T) {
	content := `
this._meter.createHistogram('gen_ai.client.operation.duration', {
  description: 'GenAI operation duration',
  unit: 's',
});
`

	metrics, err := parseInstrumentationContent(content, "")
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	assert.Equal(t, "gen_ai.client.operation.duration", metrics[0].Name)
	assert.Equal(t, "GenAI operation duration", metrics[0].Description)
	assert.Equal(t, "s", metrics[0].Unit)
	assert.Equal(t, "histogram", metrics[0].InstrumentType)
}

func TestExtractJSDocDescription(t *testing.T) {
	tests := []struct {
		name     string
		jsdoc    string
		expected string
	}{
		{
			name: "simple description",
			jsdoc: `
 * Total CPU seconds broken down by different states.
 *
 * @experimental This metric is experimental.
`,
			expected: "Total CPU seconds broken down by different states.",
		},
		{
			name: "multiline description",
			jsdoc: `
 * Event loop maximum delay.
 *
 * @note Value can be retrieved from histogram.max.
 *
 * @experimental This metric is experimental.
`,
			expected: "Event loop maximum delay.",
		},
		{
			name: "description with note mixed in",
			jsdoc: `
 * The amount of physical memory in use.
 *
 * @experimental This metric is experimental.
`,
			expected: "The amount of physical memory in use.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSDocDescription(tt.jsdoc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInferInstrumentType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"system.cpu.time", "counter"},
		{"gen_ai.client.operation.duration", "counter"},
		{"system.memory.usage", "gauge"},
		{"system.cpu.utilization", "gauge"},
		{"nodejs.eventloop.delay.max", "gauge"},
		{"v8js.memory.heap.limit", "gauge"},
		{"system.network.errors", "counter"},
		{"system.network.io", "counter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferInstrumentType(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInferUnit(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"system.cpu.time", "s"},
		{"gen_ai.client.operation.duration", "s"},
		{"nodejs.eventloop.delay.max", "s"},
		{"system.memory.usage", "By"},
		{"system.network.io", "By"},
		{"system.cpu.utilization", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferUnit(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}
