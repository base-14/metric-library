'use client';

import { useState, useEffect } from 'react';
import { useDebounce } from '@/hooks/useDebounce';

interface SearchBarProps {
  onSearch: (query: string) => void;
  initialValue?: string;
}

export function SearchBar({ onSearch, initialValue = '' }: SearchBarProps) {
  const [query, setQuery] = useState(initialValue);
  const debouncedQuery = useDebounce(query, 300);

  useEffect(() => {
    onSearch(debouncedQuery);
  }, [debouncedQuery, onSearch]);

  return (
    <div className="w-full">
      <div className="relative">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search metrics by name or description..."
          className="w-full px-4 py-3 pl-12 pr-4 text-gray-900 bg-white border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <svg
          className="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
      </div>
    </div>
  );
}
