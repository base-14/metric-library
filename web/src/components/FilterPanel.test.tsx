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
