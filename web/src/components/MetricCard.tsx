'use client';

import { CanonicalMetric } from '@/types/api';
import { badgeClass, componentTone, instrumentTone, semconvTone } from '@/lib/badges';

interface MetricCardProps {
  metric: CanonicalMetric;
  onClick?: (metric: CanonicalMetric) => void;
}

export function MetricCard({ metric, onClick }: MetricCardProps) {
  const semconv = semconvTone(metric.semconv_match);
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
          {semconv && (
            <span
              className={badgeClass(semconv.tone, semconv.outlined)}
              title={metric.semconv_match === 'exact' ? 'Matches semantic convention' : 'Prefix matches semantic convention'}
            >
              {semconv.outlined ? 'SemConv~' : 'SemConv'}
            </span>
          )}
          <span className={badgeClass(instrumentTone(metric.instrument_type))}>
            {metric.instrument_type}
          </span>
        </div>
      </div>

      <p className="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">
        {metric.description || 'No description available'}
      </p>

      <div className="flex flex-wrap gap-2 text-xs">
        <span className={badgeClass(componentTone(metric.component_type))}>
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
