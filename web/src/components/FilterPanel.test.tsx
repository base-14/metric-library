import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { FilterPanel } from './FilterPanel';
import { FacetResponse } from '@/types/api';

const mockFacets: FacetResponse = {
  instrument_types: { histogram: 10, counter: 5, gauge: 3 },
  component_types: { receiver: 15, processor: 8, exporter: 5 },
  component_names: { httpreceiver: 10 },
  source_categories: { 'otel-collector-contrib': 20 },
  source_names: { 'metadata.yaml': 20 },
  confidence_levels: { high: 18, medium: 2 },
  units: { ms: 10, '1': 5 },
};

const mockFacetsWithGroupedSources: FacetResponse = {
  instrument_types: { histogram: 10 },
  component_types: { receiver: 15 },
  component_names: { httpreceiver: 10 },
  source_categories: { otel: 100 },
  source_names: {
    'otel-collector-contrib': 50,
    'otel-semconv': 30,
    'otel-python': 20,
    'prometheus-node': 40,
    'prometheus-redis': 15,
    'kubernetes-ksm': 25,
    'cloudwatch-ec2': 10,
    'openllmetry': 35,
    'openlit': 12,
    'custom-source': 8,
  },
  confidence_levels: { high: 100 },
  units: { ms: 50 },
};

describe('FilterPanel', () => {
  it('renders loading skeleton when facets are null', () => {
    render(
      <FilterPanel
        facets={null}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(document.querySelector('.animate-pulse')).toBeInTheDocument();
  });

  it('renders instrument type section', () => {
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByText('Instrument Type')).toBeInTheDocument();
    expect(screen.getByText('histogram')).toBeInTheDocument();
    expect(screen.getByText('counter')).toBeInTheDocument();
    expect(screen.getByText('gauge')).toBeInTheDocument();
  });

  it('renders component type section', () => {
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByText('Component Type')).toBeInTheDocument();
    expect(screen.getByText('receiver')).toBeInTheDocument();
    expect(screen.getByText('processor')).toBeInTheDocument();
  });

  it('renders source section', () => {
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByText('Source')).toBeInTheDocument();
    expect(screen.getByText('otel-collector-contrib')).toBeInTheDocument();
  });

  it('displays counts for each filter option', () => {
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    // httpreceiver and histogram both have count 10
    expect(screen.getAllByText('10')).toHaveLength(2);
    // counter and exporter both have 5
    expect(screen.getAllByText('5')).toHaveLength(2);
  });

  it('renders component section first', () => {
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByText('Component')).toBeInTheDocument();
    expect(screen.getByText('httpreceiver')).toBeInTheDocument();
  });

  it('calls onFilterChange when filter is clicked', () => {
    const onFilterChange = vi.fn();
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{}}
        onFilterChange={onFilterChange}
      />
    );

    fireEvent.click(screen.getByText('histogram'));
    expect(onFilterChange).toHaveBeenCalledWith('instrument_type', 'histogram');
  });

  it('clears filter when same filter is clicked again', () => {
    const onFilterChange = vi.fn();
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{ instrument_type: 'histogram' }}
        onFilterChange={onFilterChange}
      />
    );

    fireEvent.click(screen.getByText('histogram'));
    expect(onFilterChange).toHaveBeenCalledWith('instrument_type', undefined);
  });

  it('highlights selected filter', () => {
    render(
      <FilterPanel
        facets={mockFacets}
        selectedFilters={{ instrument_type: 'histogram' }}
        onFilterChange={() => {}}
      />
    );

    const histogramButton = screen.getByText('histogram').closest('button');
    expect(histogramButton).toHaveClass('bg-blue-100');
  });
});

describe('SourceFilterSection', () => {
  it('renders search input', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByPlaceholderText('Search sources...')).toBeInTheDocument();
  });

  it('renders grouped sources with group labels', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByText('OpenTelemetry')).toBeInTheDocument();
    expect(screen.getByText('Prometheus')).toBeInTheDocument();
    expect(screen.getByText('Kubernetes')).toBeInTheDocument();
    expect(screen.getByText('CloudWatch')).toBeInTheDocument();
    expect(screen.getByText('OpenLLMetry')).toBeInTheDocument();
    expect(screen.getByText('OpenLIT')).toBeInTheDocument();
  });

  it('expands group when clicking on group header', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    expect(screen.queryByText('otel-collector-contrib')).not.toBeInTheDocument();

    fireEvent.click(screen.getByText('OpenTelemetry'));

    expect(screen.getByText('otel-collector-contrib')).toBeInTheDocument();
    expect(screen.getByText('otel-semconv')).toBeInTheDocument();
    expect(screen.getByText('otel-python')).toBeInTheDocument();
  });

  it('collapses group when clicking expanded group header', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    fireEvent.click(screen.getByText('OpenTelemetry'));
    expect(screen.getByText('otel-collector-contrib')).toBeInTheDocument();

    fireEvent.click(screen.getByText('OpenTelemetry'));
    expect(screen.queryByText('otel-collector-contrib')).not.toBeInTheDocument();
  });

  it('filters sources when typing in search', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    const searchInput = screen.getByPlaceholderText('Search sources...');
    fireEvent.change(searchInput, { target: { value: 'prom' } });

    expect(screen.getByText('prometheus-node')).toBeInTheDocument();
    expect(screen.getByText('prometheus-redis')).toBeInTheDocument();
    expect(screen.queryByText('OpenTelemetry')).not.toBeInTheDocument();
  });

  it('expands groups automatically when searching', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    expect(screen.queryByText('otel-collector-contrib')).not.toBeInTheDocument();

    const searchInput = screen.getByPlaceholderText('Search sources...');
    fireEvent.change(searchInput, { target: { value: 'otel' } });

    expect(screen.getByText('otel-collector-contrib')).toBeInTheDocument();
  });

  it('shows ungrouped sources directly', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    expect(screen.getByText('custom-source')).toBeInTheDocument();
  });

  it('calls onFilterChange when clicking source in expanded group', () => {
    const onFilterChange = vi.fn();
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={onFilterChange}
      />
    );

    fireEvent.click(screen.getByText('OpenTelemetry'));
    fireEvent.click(screen.getByText('otel-collector-contrib'));

    expect(onFilterChange).toHaveBeenCalledWith('source_name', 'otel-collector-contrib');
  });

  it('displays group total count', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    const otelGroup = screen.getByText('OpenTelemetry').closest('button');
    const countSpan = otelGroup?.querySelectorAll('span.text-xs')[1];
    expect(countSpan).toHaveTextContent('100');
  });

  it('highlights group header when a source in the group is selected', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithGroupedSources}
        selectedFilters={{ source_name: 'otel-collector-contrib' }}
        onFilterChange={() => {}}
      />
    );

    const otelGroupButton = screen.getByText('OpenTelemetry').closest('button');
    expect(otelGroupButton).toHaveClass('bg-blue-50');
  });
});

describe('SearchableFilterSection (Component)', () => {
  const mockFacetsWithManyComponents: FacetResponse = {
    instrument_types: { histogram: 10 },
    component_types: { receiver: 15 },
    component_names: {
      httpreceiver: 50,
      kafkareceiver: 30,
      prometheusreceiver: 25,
      jaegerreceiver: 20,
      otlpreceiver: 15,
    },
    source_categories: { otel: 100 },
    source_names: { 'otel-collector-contrib': 100 },
    confidence_levels: { high: 100 },
    units: { ms: 50 },
  };

  it('renders search input for components', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithManyComponents}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );
    expect(screen.getByPlaceholderText('Search components...')).toBeInTheDocument();
  });

  it('filters components when typing in search', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithManyComponents}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    expect(screen.getByText('httpreceiver')).toBeInTheDocument();
    expect(screen.getByText('kafkareceiver')).toBeInTheDocument();

    const searchInput = screen.getByPlaceholderText('Search components...');
    fireEvent.change(searchInput, { target: { value: 'kafka' } });

    expect(screen.getByText('kafkareceiver')).toBeInTheDocument();
    expect(screen.queryByText('httpreceiver')).not.toBeInTheDocument();
  });

  it('calls onFilterChange when clicking filtered component', () => {
    const onFilterChange = vi.fn();
    render(
      <FilterPanel
        facets={mockFacetsWithManyComponents}
        selectedFilters={{}}
        onFilterChange={onFilterChange}
      />
    );

    const searchInput = screen.getByPlaceholderText('Search components...');
    fireEvent.change(searchInput, { target: { value: 'prometheus' } });
    fireEvent.click(screen.getByText('prometheusreceiver'));

    expect(onFilterChange).toHaveBeenCalledWith('component_name', 'prometheusreceiver');
  });

  it('shows all components when search is cleared', () => {
    render(
      <FilterPanel
        facets={mockFacetsWithManyComponents}
        selectedFilters={{}}
        onFilterChange={() => {}}
      />
    );

    const searchInput = screen.getByPlaceholderText('Search components...');
    fireEvent.change(searchInput, { target: { value: 'kafka' } });
    expect(screen.queryByText('httpreceiver')).not.toBeInTheDocument();

    fireEvent.change(searchInput, { target: { value: '' } });
    expect(screen.getByText('httpreceiver')).toBeInTheDocument();
    expect(screen.getByText('kafkareceiver')).toBeInTheDocument();
  });
});
