import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { SearchBar } from './SearchBar';

describe('SearchBar', () => {
  it('renders search input with placeholder', () => {
    render(<SearchBar onSearch={() => {}} />);
    expect(
      screen.getByPlaceholderText('Search metrics by name or description...')
    ).toBeInTheDocument();
  });

  it('renders with initial value', () => {
    render(<SearchBar onSearch={() => {}} initialValue="test query" />);
    expect(screen.getByDisplayValue('test query')).toBeInTheDocument();
  });

  it('updates input value when typing', () => {
    const onSearch = vi.fn();
    render(<SearchBar onSearch={onSearch} />);

    const input = screen.getByPlaceholderText(
      'Search metrics by name or description...'
    );
    fireEvent.change(input, { target: { value: 'http.request' } });

    expect(input).toHaveValue('http.request');
  });

  it('renders search icon', () => {
    render(<SearchBar onSearch={() => {}} />);
    const svg = document.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });
});
