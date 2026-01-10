import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { MetricCard } from './MetricCard';
import { CanonicalMetric } from '@/types/api';

const mockMetric: CanonicalMetric = {
  id: 'test-id',
  metric_name: 'http.server.request.duration',
  instrument_type: 'histogram',
  description: 'Duration of HTTP server requests',
  unit: 'ms',
  attributes: [],
  enabled_by_default: true,
  component_type: 'receiver',
  component_name: 'httpreceiver',
  source_category: 'otel-collector-contrib',
  source_name: 'metadata.yaml',
  source_location: '/receiver/httpreceiver',
  extraction_method: 'yaml',
  source_confidence: 'high',
  repo: 'https://github.com/open-telemetry/opentelemetry-collector-contrib',
  path: '/receiver/httpreceiver/metadata.yaml',
  commit: 'abc123',
  extracted_at: '2024-01-01T00:00:00Z',
};

describe('MetricCard', () => {
  it('renders metric name', () => {
    render(<MetricCard metric={mockMetric} />);
    expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
  });

  it('renders metric description', () => {
    render(<MetricCard metric={mockMetric} />);
    expect(
      screen.getByText('Duration of HTTP server requests')
    ).toBeInTheDocument();
  });

  it('renders instrument type badge', () => {
    render(<MetricCard metric={mockMetric} />);
    expect(screen.getByText('histogram')).toBeInTheDocument();
  });

  it('renders component type and name', () => {
    render(<MetricCard metric={mockMetric} />);
    expect(screen.getByText('receiver')).toBeInTheDocument();
    expect(screen.getByText('httpreceiver')).toBeInTheDocument();
  });

  it('renders unit when present', () => {
    render(<MetricCard metric={mockMetric} />);
    expect(screen.getByText('ms')).toBeInTheDocument();
  });

  it('does not render unit when absent', () => {
    const metricWithoutUnit = { ...mockMetric, unit: '' };
    render(<MetricCard metric={metricWithoutUnit} />);
    expect(screen.queryByText('ms')).not.toBeInTheDocument();
  });

  it('calls onClick when clicked', () => {
    const onClick = vi.fn();
    render(<MetricCard metric={mockMetric} onClick={onClick} />);

    fireEvent.click(screen.getByText('http.server.request.duration'));
    expect(onClick).toHaveBeenCalledWith(mockMetric);
  });

  it('shows fallback description when none provided', () => {
    const metricWithoutDesc = { ...mockMetric, description: '' };
    render(<MetricCard metric={metricWithoutDesc} />);
    expect(screen.getByText('No description available')).toBeInTheDocument();
  });
});
