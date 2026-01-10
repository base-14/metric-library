'use client';

import { useState, useEffect } from 'react';
import { CanonicalMetric } from '@/types/api';

interface MetricDetailProps {
  metric: CanonicalMetric;
  onClose: () => void;
}

export function MetricDetail({ metric, onClose }: MetricDetailProps) {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [onClose]);

  const copyToClipboard = async () => {
    await navigator.clipboard.writeText(metric.metric_name);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const getGitHubUrl = () => {
    if (!metric.repo || !metric.source_location) return null;
    const repoPath = metric.source_location.split(metric.repo.replace('https://github.com/', '').split('/').slice(0, 2).join('/'))[1];
    if (!repoPath) return null;
    return `${metric.repo}/blob/${metric.commit}${repoPath}`;
  };

  const githubUrl = getGitHubUrl();

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] overflow-hidden">
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center gap-2 min-w-0 flex-1">
            <h2 className="text-xl font-semibold font-mono text-gray-900 dark:text-white break-all">
              {metric.metric_name}
            </h2>
            <button
              onClick={copyToClipboard}
              className="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-700 rounded transition-colors flex-shrink-0"
              aria-label="Copy metric name"
              title="Copy metric name"
            >
              {copied ? (
                <svg className="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              ) : (
                <svg className="w-4 h-4 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              )}
            </button>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-full transition-colors"
            aria-label="Close"
          >
            <svg
              className="w-5 h-5 text-gray-500 dark:text-gray-400"
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
              <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">
                Description
              </h3>
              <p className="text-gray-700 dark:text-gray-300">
                {metric.description || 'No description available'}
              </p>
            </section>

            {metric.semconv_match && (metric.semconv_match === 'exact' || metric.semconv_match === 'prefix') && (
              <section>
                <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">
                  Semantic Conventions
                </h3>
                <div className={`p-3 rounded-lg border ${
                  metric.semconv_match === 'exact'
                    ? 'bg-emerald-50 dark:bg-emerald-900/20 border-emerald-200 dark:border-emerald-800'
                    : 'bg-amber-50 dark:bg-amber-900/20 border-amber-200 dark:border-amber-800'
                }`}>
                  <div className="flex items-center gap-2 mb-1">
                    <span className={`px-2 py-0.5 text-xs font-medium rounded ${
                      metric.semconv_match === 'exact'
                        ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-900 dark:text-emerald-200'
                        : 'bg-amber-100 text-amber-800 dark:bg-amber-900 dark:text-amber-200'
                    }`}>
                      {metric.semconv_match === 'exact' ? 'Exact Match' : 'Prefix Match'}
                    </span>
                    {metric.semconv_stability && (
                      <span className="px-2 py-0.5 text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded">
                        {metric.semconv_stability}
                      </span>
                    )}
                  </div>
                  {metric.semconv_name && (
                    <p className="text-sm font-mono text-gray-700 dark:text-gray-300 mt-1">
                      {metric.semconv_name}
                    </p>
                  )}
                  <p className="text-xs text-gray-600 dark:text-gray-400 mt-2">
                    {metric.semconv_match === 'exact'
                      ? 'This metric exactly matches the OpenTelemetry semantic conventions.'
                      : 'This metric has a name that starts with a semantic convention metric.'}
                  </p>
                </div>
              </section>
            )}

            <section>
              <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">
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
                <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">
                  Attributes
                </h3>
                <div className="space-y-3">
                  {metric.attributes.map((attr) => (
                    <div
                      key={attr.name}
                      className="p-3 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700"
                    >
                      <div className="flex items-center gap-2 mb-1">
                        <span className="font-mono text-sm font-semibold text-gray-900 dark:text-white">
                          {attr.name}
                        </span>
                        <span className="px-2 py-0.5 text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded">
                          {attr.type}
                        </span>
                        {attr.required && (
                          <span className="px-2 py-0.5 text-xs bg-red-100 dark:bg-red-900 text-red-700 dark:text-red-300 rounded">
                            required
                          </span>
                        )}
                      </div>
                      {attr.description && (
                        <p className="text-sm text-gray-600 dark:text-gray-400">{attr.description}</p>
                      )}
                      {attr.enum && attr.enum.length > 0 && (
                        <div className="mt-2">
                          <span className="text-xs text-gray-500 dark:text-gray-400">Allowed values: </span>
                          <span className="text-xs font-mono text-gray-700 dark:text-gray-300">
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
              <h3 className="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide mb-2">
                Source Information
              </h3>
              <dl className="grid grid-cols-1 gap-2 text-sm">
                <div>
                  <dt className="text-gray-500 dark:text-gray-400">Repository</dt>
                  <dd className="font-mono text-gray-700 dark:text-gray-300 break-all">{metric.repo}</dd>
                </div>
                <div>
                  <dt className="text-gray-500 dark:text-gray-400">Commit</dt>
                  <dd className="font-mono text-gray-700 dark:text-gray-300">{metric.commit?.slice(0, 12)}</dd>
                </div>
                <div>
                  <dt className="text-gray-500 dark:text-gray-400">Extracted At</dt>
                  <dd className="text-gray-700 dark:text-gray-300">
                    {new Date(metric.extracted_at).toLocaleString()}
                  </dd>
                </div>
                {githubUrl && (
                  <div className="mt-2">
                    <a
                      href={githubUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-flex items-center gap-1.5 px-3 py-1.5 bg-gray-900 dark:bg-gray-700 text-white text-sm rounded hover:bg-gray-700 dark:hover:bg-gray-600 transition-colors"
                    >
                      <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                        <path fillRule="evenodd" clipRule="evenodd" d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.604-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.138 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z" />
                      </svg>
                      View on GitHub
                    </a>
                  </div>
                )}
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
      <dt className="text-sm text-gray-500 dark:text-gray-400">{label}</dt>
      <dd className="text-gray-900 dark:text-white">{value}</dd>
    </div>
  );
}
