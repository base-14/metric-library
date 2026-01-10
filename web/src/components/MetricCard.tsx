'use client';

import { CanonicalMetric } from '@/types/api';

interface MetricCardProps {
  metric: CanonicalMetric;
  onClick?: (metric: CanonicalMetric) => void;
}

const instrumentTypeColors: Record<string, string> = {
  counter: 'bg-green-100 text-green-800',
  updowncounter: 'bg-yellow-100 text-yellow-800',
  gauge: 'bg-blue-100 text-blue-800',
  histogram: 'bg-purple-100 text-purple-800',
  summary: 'bg-pink-100 text-pink-800',
};

const componentTypeColors: Record<string, string> = {
  receiver: 'bg-indigo-100 text-indigo-800',
  processor: 'bg-orange-100 text-orange-800',
  exporter: 'bg-teal-100 text-teal-800',
  extension: 'bg-gray-100 text-gray-800',
  connector: 'bg-cyan-100 text-cyan-800',
};

export function MetricCard({ metric, onClick }: MetricCardProps) {
  return (
    <div
      className="p-4 bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer"
      onClick={() => onClick?.(metric)}
    >
      <div className="flex items-start justify-between mb-2">
        <h3 className="text-lg font-semibold text-gray-900 break-all">
          {metric.metric_name}
        </h3>
        <div className="flex gap-2 flex-shrink-0 ml-2">
          <span
            className={`px-2 py-1 text-xs font-medium rounded ${
              instrumentTypeColors[metric.instrument_type] || 'bg-gray-100 text-gray-800'
            }`}
          >
            {metric.instrument_type}
          </span>
        </div>
      </div>

      <p className="text-sm text-gray-600 mb-3 line-clamp-2">
        {metric.description || 'No description available'}
      </p>

      <div className="flex flex-wrap gap-2 text-xs">
        <span
          className={`px-2 py-1 rounded ${
            componentTypeColors[metric.component_type] || 'bg-gray-100 text-gray-800'
          }`}
        >
          {metric.component_type}
        </span>
        <span className="px-2 py-1 bg-gray-100 text-gray-700 rounded">
          {metric.component_name}
        </span>
        {metric.unit && (
          <span className="px-2 py-1 bg-gray-50 text-gray-600 rounded border border-gray-200">
            {metric.unit}
          </span>
        )}
      </div>
    </div>
  );
}
