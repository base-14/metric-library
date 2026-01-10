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

  it('calls onSearch when form is submitted', () => {
    const onSearch = vi.fn();
    render(<SearchBar onSearch={onSearch} />);

    const input = screen.getByPlaceholderText(
      'Search metrics by name or description...'
    );
    fireEvent.change(input, { target: { value: 'http.request' } });
    fireEvent.submit(input.closest('form')!);

    expect(onSearch).toHaveBeenCalledWith('http.request');
  });

  it('renders search button', () => {
    render(<SearchBar onSearch={() => {}} />);
    expect(screen.getByRole('button', { name: /search/i })).toBeInTheDocument();
  });
});
