'use client';

import { Suspense, useState, useEffect, useCallback, useMemo, useRef } from 'react';
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

  const filterValues = useMemo(() => {
    const values: Record<string, string> = {};
    for (const key of FILTER_KEYS) {
      const value = urlSearchParams.get(key);
      if (value) {
        values[key] = value;
      }
    }
    return values;
  }, [urlSearchParams]);

  const filterKey = JSON.stringify(filterValues);

  const searchParams = useMemo<SearchParams>(() => {
    return { limit: 20, offset: 0, ...filterValues };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterKey]);

  const updateMetricInUrl = useCallback((metricId: string | null) => {
    const params = new URLSearchParams(window.location.search);
    if (metricId) {
      params.set('metric', metricId);
    } else {
      params.delete('metric');
    }
    const queryString = params.toString();
    const newUrl = queryString ? `?${queryString}` : window.location.pathname;
    window.history.replaceState(null, '', newUrl);
  }, []);

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

  const selectedMetricIdRef = useRef<string | null>(null);

  useEffect(() => {
    selectedMetricIdRef.current = selectedMetric?.id ?? null;
  }, [selectedMetric]);

  useEffect(() => {
    const metricId = urlSearchParams.get('metric');
    if (metricId && metricId !== selectedMetricIdRef.current) {
      getMetric(metricId)
        .then(setSelectedMetric)
        .catch(() => {
          updateMetricInUrl(null);
        });
    }
    // Only react to URL changes, not state changes
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [urlSearchParams]);

  const handleSearch = useCallback(
    (query: string) => {
      const newParams = { ...searchParams, q: query || undefined, offset: 0 };
      updateUrl(newParams, selectedMetricIdRef.current);
    },
    [searchParams, updateUrl]
  );

  const handleFilterChange = useCallback(
    (key: string, value: string | undefined) => {
      const newParams = { ...searchParams, [key]: value, offset: 0 };
      updateUrl(newParams, selectedMetricIdRef.current);
    },
    [searchParams, updateUrl]
  );

  const handleMetricClick = useCallback(
    (metric: CanonicalMetric) => {
      setSelectedMetric(metric);
      updateMetricInUrl(metric.id);
    },
    [updateMetricInUrl]
  );

  const handleMetricClose = useCallback(() => {
    setSelectedMetric(null);
    updateMetricInUrl(null);
  }, [updateMetricInUrl]);

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
            <span className="text-sm text-secondary">Active filters:</span>
            {searchParams.q && (
              <span className="inline-flex items-center gap-1 px-2 py-1 bg-hue-blue/15 text-hue-blue text-sm rounded">
                Search: {searchParams.q}
                <button
                  onClick={() => handleSearch('')}
                  className="ml-1 hover:text-hue-blue"
                  aria-label="Clear search"
                >
                  <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            )}
            {activeFilters.map(([key, value]) => (
              <span key={key} className="inline-flex items-center gap-1 px-2 py-1 bg-hue-blue/15 text-hue-blue text-sm rounded">
                {filterLabels[key]}: {key === 'semconv_match' ? semconvMatchLabels[value as string] || value : value}
                <button
                  onClick={() => handleFilterChange(key, undefined)}
                  className="ml-1 hover:text-hue-blue"
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
              className="text-sm text-secondary hover:text-primary underline"
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
              <div className="p-4 mb-4 bg-hue-red/10 border border-hue-red/40 rounded-[4px] text-hue-red">
                {error}
              </div>
            )}

            <div className="mb-4 text-sm text-secondary">
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
                  className="px-6 py-2 bg-accent text-white rounded-[4px] hover:bg-accent-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {loadingMore ? 'Loading...' : 'Load More'}
                </button>
              </div>
            )}

            {!loading && metrics.length === 0 && (
              <div className="text-center py-12 text-secondary">
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
        <div className="h-12 bg-secondary rounded-[4px] animate-pulse" />
      </div>
      <div className="flex flex-col lg:flex-row gap-8">
        <aside className="w-full lg:w-64 flex-shrink-0">
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-32 bg-secondary rounded-[4px] animate-pulse" />
            ))}
          </div>
        </aside>
        <div className="flex-1">
          <div className="space-y-4">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-24 bg-secondary rounded-[4px] animate-pulse" />
            ))}
          </div>
        </div>
      </div>
    </main>
  );
}

export default function Home() {
  return (
    <div className="min-h-screen bg-canvas">
      <header className="bg-primary border-b border-weak">
        <div className="max-w-7xl mx-auto px-4 py-6 flex items-center justify-between">
          <div>
            <h1 className="text-lg font-semibold text-primary font-mono">metric-registry</h1>
            <p className="text-sm text-secondary mt-1">
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
