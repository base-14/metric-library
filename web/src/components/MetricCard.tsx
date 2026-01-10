'use client';

import { CanonicalMetric } from '@/types/api';

interface MetricCardProps {
  metric: CanonicalMetric;
  onClick?: (metric: CanonicalMetric) => void;
}

const instrumentTypeColors: Record<string, string> = {
  counter: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200',
  updowncounter: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200',
  gauge: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200',
  histogram: 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200',
  summary: 'bg-pink-100 text-pink-800 dark:bg-pink-900 dark:text-pink-200',
};

const componentTypeColors: Record<string, string> = {
  receiver: 'bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-200',
  processor: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-200',
  exporter: 'bg-teal-100 text-teal-800 dark:bg-teal-900 dark:text-teal-200',
  extension: 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200',
  connector: 'bg-cyan-100 text-cyan-800 dark:bg-cyan-900 dark:text-cyan-200',
};

const semconvMatchColors: Record<string, string> = {
  exact: 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-200',
  prefix: 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200',
};

const semconvMatchLabels: Record<string, string> = {
  exact: 'SemConv',
  prefix: 'SemConv~',
};

export function MetricCard({ metric, onClick }: MetricCardProps) {
  return (
    <div
      className="p-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer"
      onClick={() => onClick?.(metric)}
    >
      <div className="flex items-start justify-between mb-2">
        <h3 className="text-lg font-semibold font-mono text-gray-900 dark:text-white break-all">
          {metric.metric_name}
        </h3>
        <div className="flex gap-2 flex-shrink-0 ml-2">
          {metric.semconv_match && (metric.semconv_match === 'exact' || metric.semconv_match === 'prefix') && (
            <span
              className={`px-2 py-1 text-xs font-medium rounded ${semconvMatchColors[metric.semconv_match]}`}
              title={metric.semconv_match === 'exact' ? 'Matches semantic convention' : 'Prefix matches semantic convention'}
            >
              {semconvMatchLabels[metric.semconv_match]}
            </span>
          )}
          <span
            className={`px-2 py-1 text-xs font-medium rounded ${
              instrumentTypeColors[metric.instrument_type] || 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
            }`}
          >
            {metric.instrument_type}
          </span>
        </div>
      </div>

      <p className="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">
        {metric.description || 'No description available'}
      </p>

      <div className="flex flex-wrap gap-2 text-xs">
        <span
          className={`px-2 py-1 rounded ${
            componentTypeColors[metric.component_type] || 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
          }`}
        >
          {metric.component_type}
        </span>
        <span className="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded">
          {metric.component_name}
        </span>
        {metric.unit && (
          <span className="px-2 py-1 bg-gray-50 dark:bg-gray-900 text-gray-600 dark:text-gray-400 rounded border border-gray-200 dark:border-gray-600">
            {metric.unit}
          </span>
        )}
      </div>
    </div>
  );
}
