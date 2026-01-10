'use client';

import { CanonicalMetric } from '@/types/api';

interface MetricDetailProps {
  metric: CanonicalMetric;
  onClose: () => void;
}

export function MetricDetail({ metric, onClose }: MetricDetailProps) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] overflow-hidden">
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <h2 className="text-xl font-semibold text-gray-900 break-all">
            {metric.metric_name}
          </h2>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-full transition-colors"
            aria-label="Close"
          >
            <svg
              className="w-5 h-5 text-gray-500"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        <div className="p-6 overflow-y-auto max-h-[calc(90vh-80px)]">
          <div className="space-y-6">
            <section>
              <h3 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2">
                Description
              </h3>
              <p className="text-gray-700">
                {metric.description || 'No description available'}
              </p>
            </section>

            <section>
              <h3 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2">
                Details
              </h3>
              <dl className="grid grid-cols-2 gap-4">
                <DetailItem label="Instrument Type" value={metric.instrument_type} />
                <DetailItem label="Unit" value={metric.unit || 'N/A'} />
                <DetailItem label="Component Type" value={metric.component_type} />
                <DetailItem label="Component Name" value={metric.component_name} />
                <DetailItem label="Source Category" value={metric.source_category} />
                <DetailItem label="Source Name" value={metric.source_name} />
                <DetailItem
                  label="Enabled by Default"
                  value={metric.enabled_by_default ? 'Yes' : 'No'}
                />
                <DetailItem label="Confidence" value={metric.source_confidence} />
              </dl>
            </section>

            {metric.attributes && metric.attributes.length > 0 && (
              <section>
                <h3 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2">
                  Attributes
                </h3>
                <div className="space-y-3">
                  {metric.attributes.map((attr) => (
                    <div
                      key={attr.name}
                      className="p-3 bg-gray-50 rounded-lg border border-gray-200"
                    >
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-mono text-sm font-semibold text-gray-900">
                          {attr.name}
                        </span>
                        <span className="px-2 py-0.5 text-xs bg-gray-200 text-gray-700 rounded">
                          {attr.type}
                        </span>
                        {attr.required && (
                          <span className="px-2 py-0.5 text-xs bg-red-100 text-red-700 rounded">
                            required
                          </span>
                        )}
                      </div>
                      {attr.description && (
                        <p className="text-sm text-gray-600">{attr.description}</p>
                      )}
                      {attr.enum && attr.enum.length > 0 && (
                        <div className="mt-2">
                          <span className="text-xs text-gray-500">Allowed values: </span>
                          <span className="text-xs font-mono text-gray-700">
                            {attr.enum.join(', ')}
                          </span>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </section>
            )}

            <section>
              <h3 className="text-sm font-semibold text-gray-500 uppercase tracking-wide mb-2">
                Source Information
              </h3>
              <dl className="grid grid-cols-1 gap-2 text-sm">
                <div>
                  <dt className="text-gray-500">Repository</dt>
                  <dd className="font-mono text-gray-700 break-all">{metric.repo}</dd>
                </div>
                <div>
                  <dt className="text-gray-500">Path</dt>
                  <dd className="font-mono text-gray-700 break-all">{metric.path}</dd>
                </div>
                <div>
                  <dt className="text-gray-500">Commit</dt>
                  <dd className="font-mono text-gray-700">{metric.commit}</dd>
                </div>
                <div>
                  <dt className="text-gray-500">Extracted At</dt>
                  <dd className="text-gray-700">
                    {new Date(metric.extracted_at).toLocaleString()}
                  </dd>
                </div>
              </dl>
            </section>
          </div>
        </div>
      </div>
    </div>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-sm text-gray-500">{label}</dt>
      <dd className="text-gray-900 capitalize">{value}</dd>
    </div>
  );
}
