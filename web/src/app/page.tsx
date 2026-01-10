'use client';

import { Suspense, useState, useEffect, useCallback, useMemo } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { SearchBar } from '@/components/SearchBar';
import { FilterPanel } from '@/components/FilterPanel';
import { MetricCard } from '@/components/MetricCard';
import { MetricDetail } from '@/components/MetricDetail';
import { ThemeToggle } from '@/components/ThemeToggle';
import { searchMetrics, getFacets, getMetric } from '@/lib/api';
import { CanonicalMetric, FacetResponse, SearchParams } from '@/types/api';

const FILTER_KEYS = [
  'q',
  'instrument_type',
  'component_type',
  'component_name',
  'source_category',
  'source_name',
  'semconv_match',
] as const;

function HomeContent() {
  const router = useRouter();
  const urlSearchParams = useSearchParams();

  const [metrics, setMetrics] = useState<CanonicalMetric[]>([]);
  const [facets, setFacets] = useState<FacetResponse | null>(null);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedMetric, setSelectedMetric] = useState<CanonicalMetric | null>(null);

  const searchParams = useMemo<SearchParams>(() => {
    const params: SearchParams = { limit: 20, offset: 0 };
    for (const key of FILTER_KEYS) {
      const value = urlSearchParams.get(key);
      if (value) {
        params[key] = value;
      }
    }
    return params;
  }, [urlSearchParams]);

  const updateUrl = useCallback(
    (newParams: SearchParams, metricId?: string | null) => {
      const params = new URLSearchParams();
      for (const key of FILTER_KEYS) {
        const value = newParams[key];
        if (value) {
          params.set(key, value);
        }
      }
      if (metricId) {
        params.set('metric', metricId);
      }
      const queryString = params.toString();
      router.replace(queryString ? `?${queryString}` : '/', { scroll: false });
    },
    [router]
  );

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

  const fetchFacets = useCallback(async (sourceName?: string) => {
    try {
      const response = await getFacets(sourceName);
      setFacets(response);
    } catch (err) {
      console.error('Failed to fetch facets:', err);
    }
  }, []);

  useEffect(() => {
    fetchMetrics();
  }, [fetchMetrics]);

  useEffect(() => {
    fetchFacets(searchParams.source_name);
  }, [fetchFacets, searchParams.source_name]);

  useEffect(() => {
    const metricId = urlSearchParams.get('metric');
    if (metricId && (!selectedMetric || selectedMetric.id !== metricId)) {
      getMetric(metricId)
        .then(setSelectedMetric)
        .catch(() => {
          updateUrl(searchParams, null);
        });
    } else if (!metricId && selectedMetric) {
      setSelectedMetric(null);
    }
  }, [urlSearchParams, selectedMetric, searchParams, updateUrl]);

  const handleSearch = useCallback(
    (query: string) => {
      const newParams = { ...searchParams, q: query || undefined, offset: 0 };
      updateUrl(newParams, urlSearchParams.get('metric'));
    },
    [searchParams, updateUrl, urlSearchParams]
  );

  const handleFilterChange = useCallback(
    (key: string, value: string | undefined) => {
      const newParams = { ...searchParams, [key]: value, offset: 0 };
      updateUrl(newParams, urlSearchParams.get('metric'));
    },
    [searchParams, updateUrl, urlSearchParams]
  );

  const handleMetricClick = useCallback(
    (metric: CanonicalMetric) => {
      setSelectedMetric(metric);
      updateUrl(searchParams, metric.id);
    },
    [searchParams, updateUrl]
  );

  const handleMetricClose = useCallback(() => {
    setSelectedMetric(null);
    updateUrl(searchParams, null);
  }, [searchParams, updateUrl]);

  const handleLoadMore = async () => {
    setLoadingMore(true);
    try {
      const response = await searchMetrics({
        ...searchParams,
        offset: metrics.length,
      });
      setMetrics(prev => [...prev, ...response.metrics]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load more metrics');
    } finally {
      setLoadingMore(false);
    }
  };

  const selectedFilters = {
    instrument_type: searchParams.instrument_type,
    component_type: searchParams.component_type,
    component_name: searchParams.component_name,
    source_category: searchParams.source_category,
    source_name: searchParams.source_name,
    semconv_match: searchParams.semconv_match,
  };

  const activeFilters = Object.entries(selectedFilters).filter(([, value]) => value);
  const hasActiveFilters = activeFilters.length > 0 || searchParams.q;

  const clearAllFilters = () => {
    updateUrl({ limit: 20, offset: 0 }, null);
  };

  const filterLabels: Record<string, string> = {
    instrument_type: 'Type',
    component_type: 'Component Type',
    component_name: 'Component',
    source_category: 'Source Category',
    source_name: 'Source',
    semconv_match: 'SemConv',
  };

  const semconvMatchLabels: Record<string, string> = {
    exact: 'Exact Match',
    prefix: 'Prefix Match',
    none: 'No Match',
  };

  return (
    <>
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
                {filterLabels[key]}: {key === 'semconv_match' ? semconvMatchLabels[value as string] || value : value}
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
                  onClick={handleMetricClick}
                />
              ))}
            </div>

            {!loading && metrics.length > 0 && metrics.length < total && (
              <div className="mt-6 text-center">
                <button
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                  className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loadingMore ? 'Loading...' : 'Load More'}
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
          onClose={handleMetricClose}
        />
      )}
    </>
  );
}

function LoadingFallback() {
  return (
    <main className="max-w-7xl mx-auto px-4 py-8">
      <div className="mb-6">
        <div className="h-12 bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
      </div>
      <div className="flex flex-col lg:flex-row gap-8">
        <aside className="w-full lg:w-64 flex-shrink-0">
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-32 bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
            ))}
          </div>
        </aside>
        <div className="flex-1">
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-24 bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    </main>
  );
}

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <header className="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
        <div className="max-w-7xl mx-auto px-4 py-6 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-white font-mono">metric-library</h1>
            <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
              Metric discovery platform
            </p>
          </div>
          <ThemeToggle />
        </div>
      </header>

      <Suspense fallback={<LoadingFallback />}>
        <HomeContent />
      </Suspense>
    </div>
  );
}
