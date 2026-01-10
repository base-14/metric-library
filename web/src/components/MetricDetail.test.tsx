import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { MetricDetail } from './MetricDetail';
import { CanonicalMetric } from '@/types/api';

const mockMetric: CanonicalMetric = {
  id: 'test-id',
  metric_name: 'http.server.request.duration',
  instrument_type: 'histogram',
  description: 'Duration of HTTP server requests',
  unit: 'ms',
  attributes: [
    {
      name: 'http.method',
      type: 'string',
      description: 'HTTP method',
      required: true,
      enum: ['GET', 'POST', 'PUT', 'DELETE'],
    },
    {
      name: 'http.status_code',
      type: 'int',
      description: 'HTTP status code',
      required: false,
    },
  ],
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

describe('MetricDetail', () => {
  it('renders metric name in header', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
  });

  it('renders description section', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('Description')).toBeInTheDocument();
    expect(
      screen.getByText('Duration of HTTP server requests')
    ).toBeInTheDocument();
  });

  it('renders detail fields', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('Instrument Type')).toBeInTheDocument();
    expect(screen.getByText('histogram')).toBeInTheDocument();
    expect(screen.getByText('Component Type')).toBeInTheDocument();
    expect(screen.getByText('receiver')).toBeInTheDocument();
  });

  it('renders attributes section', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('Attributes')).toBeInTheDocument();
    expect(screen.getByText('http.method')).toBeInTheDocument();
    expect(screen.getByText('http.status_code')).toBeInTheDocument();
  });

  it('shows required badge for required attributes', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('required')).toBeInTheDocument();
  });

  it('shows enum values when present', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('GET, POST, PUT, DELETE')).toBeInTheDocument();
  });

  it('renders source information', () => {
    render(<MetricDetail metric={mockMetric} onClose={() => {}} />);
    expect(screen.getByText('Source Information')).toBeInTheDocument();
    expect(
      screen.getByText(
        'https://github.com/open-telemetry/opentelemetry-collector-contrib'
      )
    ).toBeInTheDocument();
    expect(screen.getByText('abc123')).toBeInTheDocument();
  });

  it('calls onClose when close button is clicked', () => {
    const onClose = vi.fn();
    render(<MetricDetail metric={mockMetric} onClose={onClose} />);

    fireEvent.click(screen.getByLabelText('Close'));
    expect(onClose).toHaveBeenCalled();
  });

  it('does not render attributes section when no attributes', () => {
    const metricWithoutAttrs = { ...mockMetric, attributes: [] };
    render(<MetricDetail metric={metricWithoutAttrs} onClose={() => {}} />);
    expect(screen.queryByText('Attributes')).not.toBeInTheDocument();
  });
});
