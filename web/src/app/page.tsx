'use client';

import { useState, useEffect, useCallback } from 'react';
import { SearchBar } from '@/components/SearchBar';
import { FilterPanel } from '@/components/FilterPanel';
import { MetricCard } from '@/components/MetricCard';
import { MetricDetail } from '@/components/MetricDetail';
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

  const handleSearch = (query: string) => {
    setSearchParams((prev) => ({
      ...prev,
      q: query || undefined,
      offset: 0,
    }));
  };

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
    source_category: searchParams.source_category,
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <h1 className="text-2xl font-bold text-gray-900">OTel Glossary</h1>
          <p className="text-sm text-gray-600 mt-1">
            Discover and explore OpenTelemetry metrics
          </p>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-8">
        <div className="mb-8">
          <SearchBar onSearch={handleSearch} initialValue={searchParams.q} />
        </div>

        <div className="flex gap-8">
          <aside className="w-64 flex-shrink-0">
            <FilterPanel
              facets={facets}
              selectedFilters={selectedFilters}
              onFilterChange={handleFilterChange}
            />
          </aside>

          <div className="flex-1">
            {error && (
              <div className="p-4 mb-4 bg-red-50 border border-red-200 rounded-lg text-red-700">
                {error}
              </div>
            )}

            <div className="mb-4 text-sm text-gray-600">
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
              <div className="text-center py-12 text-gray-500">
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
