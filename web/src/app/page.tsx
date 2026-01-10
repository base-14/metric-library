'use client';

import { useState, useEffect, useCallback } from 'react';
import { SearchBar } from '@/components/SearchBar';
import { FilterPanel } from '@/components/FilterPanel';
import { MetricCard } from '@/components/MetricCard';
import { MetricDetail } from '@/components/MetricDetail';
import { ThemeToggle } from '@/components/ThemeToggle';
import { searchMetrics, getFacets } from '@/lib/api';
import { CanonicalMetric, FacetResponse, SearchParams } from '@/types/api';

export default function Home() {
  const [metrics, setMetrics] = useState<CanonicalMetric[]>([]);
  const [facets, setFacets] = useState<FacetResponse | null>(null);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedMetric, setSelectedMetric] = useState<CanonicalMetric | null>(null);
  const [searchParams, setSearchParams] = useState<SearchParams>({
    limit: 20,
    offset: 0,
  });

  const fetchMetrics = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await searchMetrics(searchParams);
      setMetrics(response.metrics);
      setTotal(response.total);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch metrics');
    } finally {
      setLoading(false);
    }
  }, [searchParams]);

  const fetchFacets = useCallback(async () => {
    try {
      const response = await getFacets();
      setFacets(response);
    } catch (err) {
      console.error('Failed to fetch facets:', err);
    }
  }, []);

  useEffect(() => {
    fetchMetrics();
  }, [fetchMetrics]);

  useEffect(() => {
    fetchFacets();
  }, [fetchFacets]);

  const handleSearch = useCallback((query: string) => {
    setSearchParams((prev) => ({
      ...prev,
      q: query || undefined,
      offset: 0,
    }));
  }, []);

  const handleFilterChange = (key: string, value: string | undefined) => {
    setSearchParams((prev) => ({
      ...prev,
      [key]: value,
      offset: 0,
    }));
  };

  const handleLoadMore = () => {
    setSearchParams((prev) => ({
      ...prev,
      offset: (prev.offset || 0) + (prev.limit || 20),
    }));
  };

  const selectedFilters = {
    instrument_type: searchParams.instrument_type,
    component_type: searchParams.component_type,
    component_name: searchParams.component_name,
    source_category: searchParams.source_category,
  };

  const activeFilters = Object.entries(selectedFilters).filter(([, value]) => value);
  const hasActiveFilters = activeFilters.length > 0 || searchParams.q;

  const clearAllFilters = () => {
    setSearchParams({ limit: 20, offset: 0 });
  };

  const filterLabels: Record<string, string> = {
    instrument_type: 'Type',
    component_type: 'Component Type',
    component_name: 'Component',
    source_category: 'Source',
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 py-6 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white">OTel Glossary</h1>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Discover and explore OpenTelemetry metrics
            </p>
          </div>
          <ThemeToggle />
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="mb-6">
          <SearchBar onSearch={handleSearch} initialValue={searchParams.q} />
        </div>

        {hasActiveFilters && (
          <div className="mb-6 flex flex-wrap items-center gap-2">
            <span className="text-sm text-gray-500 dark:text-gray-400">Active filters:</span>
            {searchParams.q && (
              <span className="inline-flex items-center gap-1 px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-sm rounded">
                Search: {searchParams.q}
                <button
                  onClick={() => handleSearch('')}
                  className="ml-1 hover:text-blue-600 dark:hover:text-blue-400"
                  aria-label="Clear search"
                >
                  <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            )}
            {activeFilters.map(([key, value]) => (
              <span key={key} className="inline-flex items-center gap-1 px-2 py-1 bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 text-sm rounded">
                {filterLabels[key]}: {value}
                <button
                  onClick={() => handleFilterChange(key, undefined)}
                  className="ml-1 hover:text-blue-600 dark:hover:text-blue-400"
                  aria-label={`Clear ${filterLabels[key]} filter`}
                >
                  <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            ))}
            <button
              onClick={clearAllFilters}
              className="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 underline"
            >
              Clear all
            </button>
          </div>
        )}

        <div className="flex flex-col lg:flex-row gap-8">
          <aside className="w-full lg:w-64 flex-shrink-0">
            <FilterPanel
              facets={facets}
              selectedFilters={selectedFilters}
              onFilterChange={handleFilterChange}
            />
          </aside>

          <div className="flex-1">
            {error && (
              <div className="p-4 mb-4 bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 rounded-lg text-red-700 dark:text-red-400">
                {error}
              </div>
            )}

            <div className="mb-4 text-sm text-gray-600 dark:text-gray-400">
              {loading ? 'Loading...' : `${total} metrics found`}
            </div>

            <div className="space-y-4">
              {metrics.map((metric) => (
                <MetricCard
                  key={metric.id}
                  metric={metric}
                  onClick={setSelectedMetric}
                />
              ))}
            </div>

            {!loading && metrics.length > 0 && metrics.length < total && (
              <div className="mt-6 text-center">
                <button
                  onClick={handleLoadMore}
                  className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Load More
                </button>
              </div>
            )}

            {!loading && metrics.length === 0 && (
              <div className="text-center py-12 text-gray-500 dark:text-gray-400">
                No metrics found. Try adjusting your search or filters.
              </div>
            )}
          </div>
        </div>
      </main>

      {selectedMetric && (
        <MetricDetail
          metric={selectedMetric}
          onClose={() => setSelectedMetric(null)}
        />
      )}
    </div>
  );
}
