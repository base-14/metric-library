export interface Attribute {
  name: string;
  type: string;
  description: string;
  required: boolean;
  enum?: string[];
}

export interface CanonicalMetric {
  id: string;
  metric_name: string;
  instrument_type: string;
  description: string;
  unit: string;
  attributes: Attribute[];
  enabled_by_default: boolean;
  component_type: string;
  component_name: string;
  source_category: string;
  source_name: string;
  source_location: string;
  extraction_method: string;
  source_confidence: string;
  repo: string;
  path: string;
  commit: string;
  extracted_at: string;
}

export interface SearchResponse {
  metrics: CanonicalMetric[];
  total: number;
  limit: number;
  offset: number;
}

export interface FacetResponse {
  instrument_types: Record<string, number>;
  component_types: Record<string, number>;
  component_names: Record<string, number>;
  source_categories: Record<string, number>;
  source_names: Record<string, number>;
  confidence_levels: Record<string, number>;
  units: Record<string, number>;
}

export interface SearchParams {
  q?: string;
  instrument_type?: string;
  component_type?: string;
  component_name?: string;
  source_category?: string;
  source_name?: string;
  confidence?: string;
  limit?: number;
  offset?: number;
}
