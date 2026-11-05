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
      className="p-4 bg-primary border border-weak rounded-[4px] hover:bg-hover transition-colors cursor-pointer"
      onClick={() => onClick?.(metric)}
    >
      <div className="flex items-start justify-between mb-2">
        <h3 className="text-sm font-medium font-mono text-primary break-all">
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

      <p className="text-sm text-secondary mb-3 line-clamp-2">
        {metric.description || 'No description available'}
      </p>

      <div className="flex flex-wrap gap-2 text-xs">
        <span className={badgeClass(componentTone(metric.component_type))}>
          {metric.component_type}
        </span>
        <span className="px-2 py-1 bg-secondary text-primary rounded">
          {metric.component_name}
        </span>
        {metric.unit && (
          <span className="px-2 py-1 bg-secondary text-secondary rounded border border-weak">
            {metric.unit}
          </span>
        )}
      </div>
    </div>
  );
}
