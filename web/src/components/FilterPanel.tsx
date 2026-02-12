'use client';

import { useState, useMemo } from 'react';
import { FacetResponse } from '@/types/api';

const semconvLabels: Record<string, string> = {
  exact: 'Exact Match',
  prefix: 'Prefix Match',
  none: 'No Match',
};

const sourceGroupConfig: { prefix: string; label: string }[] = [
  { prefix: 'otel-', label: 'OpenTelemetry' },
  { prefix: 'prometheus-', label: 'Prometheus' },
  { prefix: 'kubernetes-', label: 'Kubernetes' },
  { prefix: 'cloudwatch-', label: 'CloudWatch' },
  { prefix: 'gcp-', label: 'GCP' },
  { prefix: 'azure-', label: 'Azure' },
  { prefix: 'openllmetry', label: 'OpenLLMetry' },
  { prefix: 'openlit', label: 'OpenLIT' },
];

interface FilterPanelProps {
  facets: FacetResponse | null;
  selectedFilters: {
    instrument_type?: string;
    component_type?: string;
    component_name?: string;
    source_category?: string;
    source_name?: string;
    semconv_match?: string;
  };
  onFilterChange: (key: string, value: string | undefined) => void;
}

export function FilterPanel({ facets, selectedFilters, onFilterChange }: FilterPanelProps) {
  const handleFilterClick = (key: string, value: string) => {
    if (selectedFilters[key as keyof typeof selectedFilters] === value) {
      onFilterChange(key, undefined);
    } else {
      onFilterChange(key, value);
    }
  };

  if (!facets) {
    return (
      <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg animate-pulse">
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24 mb-4"></div>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-6 bg-gray-200 dark:bg-gray-700 rounded"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <SourceFilterSection
        items={facets.source_names}
        selectedValue={selectedFilters.source_name}
        onSelect={(value) => handleFilterClick('source_name', value)}
      />

      <SearchableFilterSection
        title="Component"
        placeholder="Search components..."
        items={facets.component_names}
        selectedValue={selectedFilters.component_name}
        onSelect={(value) => handleFilterClick('component_name', value)}
      />

      <FilterSection
        title="Instrument Type"
        items={facets.instrument_types}
        selectedValue={selectedFilters.instrument_type}
        onSelect={(value) => handleFilterClick('instrument_type', value)}
      />

      <FilterSection
        title="Component Type"
        items={facets.component_types}
        selectedValue={selectedFilters.component_type}
        onSelect={(value) => handleFilterClick('component_type', value)}
      />

      <FilterSection
        title="Source Category"
        items={facets.source_categories}
        selectedValue={selectedFilters.source_category}
        onSelect={(value) => handleFilterClick('source_category', value)}
      />

      {facets.semconv_matches && Object.keys(facets.semconv_matches).length > 0 && (
        <FilterSection
          title="Semantic Convention"
          items={facets.semconv_matches}
          selectedValue={selectedFilters.semconv_match}
          onSelect={(value) => handleFilterClick('semconv_match', value)}
          labelMap={semconvLabels}
        />
      )}
    </div>
  );
}

interface FilterSectionProps {
  title: string;
  items: Record<string, number>;
  selectedValue?: string;
  onSelect: (value: string) => void;
  labelMap?: Record<string, string>;
}

function FilterSection({ title, items, selectedValue, onSelect, labelMap }: FilterSectionProps) {
  const sortedItems = Object.entries(items).sort((a, b) => b[1] - a[1]);

  if (sortedItems.length === 0) {
    return null;
  }

  return (
    <div>
      <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">{title}</h3>
      <div className="space-y-1">
        {sortedItems.map(([value, count]) => (
          <button
            key={value}
            onClick={() => onSelect(value)}
            className={`w-full flex items-center justify-between px-3 py-2 text-sm rounded-md transition-colors ${
              selectedValue === value
                ? 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'
                : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
            }`}
          >
            <span>{labelMap?.[value] || value}</span>
            <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">
              {count}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}

interface SearchableFilterSectionProps {
  title: string;
  placeholder: string;
  items: Record<string, number>;
  selectedValue?: string;
  onSelect: (value: string) => void;
}

function SearchableFilterSection({ title, placeholder, items, selectedValue, onSelect }: SearchableFilterSectionProps) {
  const [searchTerm, setSearchTerm] = useState('');

  const filteredItems = useMemo(() => {
    const sorted = Object.entries(items).sort((a, b) => b[1] - a[1]);
    if (!searchTerm) return sorted;
    const term = searchTerm.toLowerCase();
    return sorted.filter(([name]) => name.toLowerCase().includes(term));
  }, [items, searchTerm]);

  if (Object.keys(items).length === 0) {
    return null;
  }

  return (
    <div>
      <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">{title}</h3>
      <input
        type="text"
        placeholder={placeholder}
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        className="w-full px-3 py-2 mb-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
      />
      <div className="space-y-1">
        {filteredItems.map(([value, count]) => (
          <button
            key={value}
            onClick={() => onSelect(value)}
            className={`w-full flex items-center justify-between px-3 py-2 text-sm rounded-md transition-colors ${
              selectedValue === value
                ? 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'
                : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
            }`}
          >
            <span>{value}</span>
            <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">
              {count}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}

interface SourceFilterSectionProps {
  items: Record<string, number>;
  selectedValue?: string;
  onSelect: (value: string) => void;
}

function SourceFilterSection({ items, selectedValue, onSelect }: SourceFilterSectionProps) {
  const [searchTerm, setSearchTerm] = useState('');
  const [expandedGroups, setExpandedGroups] = useState<Set<string>>(new Set());

  const groupedSources = useMemo(() => {
    const groups: Record<string, [string, number][]> = {};
    const other: [string, number][] = [];

    Object.entries(items).forEach(([name, count]) => {
      const groupConfig = sourceGroupConfig.find((g) => name.startsWith(g.prefix));
      if (groupConfig) {
        if (!groups[groupConfig.label]) {
          groups[groupConfig.label] = [];
        }
        groups[groupConfig.label].push([name, count]);
      } else {
        other.push([name, count]);
      }
    });

    Object.values(groups).forEach((group) => group.sort((a, b) => b[1] - a[1]));
    other.sort((a, b) => b[1] - a[1]);

    return { groups, other };
  }, [items]);

  const filteredGroups = useMemo(() => {
    if (!searchTerm) return groupedSources;

    const term = searchTerm.toLowerCase();
    const filtered: Record<string, [string, number][]> = {};
    const filteredOther: [string, number][] = [];

    Object.entries(groupedSources.groups).forEach(([label, sources]) => {
      const matches = sources.filter(([name]) => name.toLowerCase().includes(term));
      if (matches.length > 0) {
        filtered[label] = matches;
      }
    });

    groupedSources.other.forEach(([name, count]) => {
      if (name.toLowerCase().includes(term)) {
        filteredOther.push([name, count]);
      }
    });

    return { groups: filtered, other: filteredOther };
  }, [groupedSources, searchTerm]);

  const toggleGroup = (label: string) => {
    setExpandedGroups((prev) => {
      const next = new Set(prev);
      if (next.has(label)) {
        next.delete(label);
      } else {
        next.add(label);
      }
      return next;
    });
  };

  const getGroupTotal = (sources: [string, number][]) => {
    return sources.reduce((sum, [, count]) => sum + count, 0);
  };

  if (Object.keys(items).length === 0) {
    return null;
  }

  const orderedGroupLabels = sourceGroupConfig
    .map((g) => g.label)
    .filter((label) => filteredGroups.groups[label]);

  return (
    <div>
      <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">Source</h3>
      <input
        type="text"
        placeholder="Search sources..."
        value={searchTerm}
        onChange={(e) => setSearchTerm(e.target.value)}
        className="w-full px-3 py-2 mb-2 text-sm border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
      />
      <div className="space-y-1">
        {orderedGroupLabels.map((label) => {
          const sources = filteredGroups.groups[label];
          const isExpanded = expandedGroups.has(label) || searchTerm.length > 0;
          const hasSelection = sources.some(([name]) => name === selectedValue);

          return (
            <div key={label}>
              <button
                onClick={() => toggleGroup(label)}
                className={`w-full flex items-center justify-between px-3 py-2 text-sm rounded-md transition-colors ${
                  hasSelection
                    ? 'bg-blue-50 dark:bg-blue-950 text-blue-700 dark:text-blue-300'
                    : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
                }`}
              >
                <span className="flex items-center gap-2">
                  <span className="text-xs text-gray-400">{isExpanded ? '▼' : '▶'}</span>
                  <span className="font-medium">{label}</span>
                </span>
                <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">
                  {getGroupTotal(sources)}
                </span>
              </button>
              {isExpanded && (
                <div className="ml-4 mt-1 space-y-1">
                  {sources.map(([name, count]) => (
                    <button
                      key={name}
                      onClick={() => onSelect(name)}
                      className={`w-full flex items-center justify-between px-3 py-1.5 text-sm rounded-md transition-colors ${
                        selectedValue === name
                          ? 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'
                          : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-600 dark:text-gray-400'
                      }`}
                    >
                      <span className="truncate">{name}</span>
                      <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded ml-2">
                        {count}
                      </span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          );
        })}
        {filteredGroups.other.length > 0 && (
          <>
            {filteredGroups.other.map(([name, count]) => (
              <button
                key={name}
                onClick={() => onSelect(name)}
                className={`w-full flex items-center justify-between px-3 py-2 text-sm rounded-md transition-colors ${
                  selectedValue === name
                    ? 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200'
                    : 'hover:bg-gray-100 dark:hover:bg-gray-700 text-gray-700 dark:text-gray-300'
                }`}
              >
                <span>{name}</span>
                <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">
                  {count}
                </span>
              </button>
            ))}
          </>
        )}
      </div>
    </div>
  );
}
