'use client';

import { FacetResponse } from '@/types/api';

interface FilterPanelProps {
  facets: FacetResponse | null;
  selectedFilters: {
    instrument_type?: string;
    component_type?: string;
    source_category?: string;
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
      <div className="p-4 bg-gray-50 rounded-lg animate-pulse">
        <div className="h-4 bg-gray-200 rounded w-24 mb-4"></div>
        <div className="space-y-2">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-6 bg-gray-200 rounded"></div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
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
        title="Source"
        items={facets.source_categories}
        selectedValue={selectedFilters.source_category}
        onSelect={(value) => handleFilterClick('source_category', value)}
      />
    </div>
  );
}

interface FilterSectionProps {
  title: string;
  items: Record<string, number>;
  selectedValue?: string;
  onSelect: (value: string) => void;
}

function FilterSection({ title, items, selectedValue, onSelect }: FilterSectionProps) {
  const sortedItems = Object.entries(items).sort((a, b) => b[1] - a[1]);

  if (sortedItems.length === 0) {
    return null;
  }

  return (
    <div>
      <h3 className="text-sm font-semibold text-gray-700 mb-2">{title}</h3>
      <div className="space-y-1">
        {sortedItems.map(([value, count]) => (
          <button
            key={value}
            onClick={() => onSelect(value)}
            className={`w-full flex items-center justify-between px-3 py-2 text-sm rounded-md transition-colors ${
              selectedValue === value
                ? 'bg-blue-100 text-blue-800'
                : 'hover:bg-gray-100 text-gray-700'
            }`}
          >
            <span className="capitalize">{value}</span>
            <span className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded">
              {count}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
