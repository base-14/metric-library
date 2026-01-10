import { CanonicalMetric, FacetResponse, SearchParams, SearchResponse } from '@/types/api';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export async function searchMetrics(params: SearchParams = {}): Promise<SearchResponse> {
  const searchParams = new URLSearchParams();

  if (params.q) searchParams.set('q', params.q);
  if (params.instrument_type) searchParams.set('instrument_type', params.instrument_type);
  if (params.component_type) searchParams.set('component_type', params.component_type);
  if (params.component_name) searchParams.set('component_name', params.component_name);
  if (params.source_category) searchParams.set('source_category', params.source_category);
  if (params.source_name) searchParams.set('source_name', params.source_name);
  if (params.confidence) searchParams.set('confidence', params.confidence);
  if (params.limit) searchParams.set('limit', params.limit.toString());
  if (params.offset) searchParams.set('offset', params.offset.toString());

  const url = `${API_BASE}/api/metrics?${searchParams.toString()}`;
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`Search failed: ${response.statusText}`);
  }

  return response.json();
}

export async function getMetric(id: string): Promise<CanonicalMetric> {
  const response = await fetch(`${API_BASE}/api/metrics/${id}`);

  if (!response.ok) {
    throw new Error(`Failed to get metric: ${response.statusText}`);
  }

  return response.json();
}

export async function getFacets(sourceName?: string): Promise<FacetResponse> {
  const params = new URLSearchParams();
  if (sourceName) params.set('source_name', sourceName);

  const url = params.toString() ? `${API_BASE}/api/facets?${params.toString()}` : `${API_BASE}/api/facets`;
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error(`Failed to get facets: ${response.statusText}`);
  }

  return response.json();
}
