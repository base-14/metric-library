import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, Mock } from 'vitest';
import Home from './page';
import { ThemeProvider } from '@/components/ThemeProvider';
import * as api from '@/lib/api';
import { CanonicalMetric } from '@/types/api';

function renderWithProviders(ui: React.ReactElement) {
  return render(<ThemeProvider>{ui}</ThemeProvider>);
}

const mockRouterReplace = vi.fn();
const mockSearchParams = new URLSearchParams();

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    replace: mockRouterReplace,
  }),
  useSearchParams: () => mockSearchParams,
}));

vi.mock('@/lib/api', () => ({
  searchMetrics: vi.fn(),
  getMetric: vi.fn(),
  getFacets: vi.fn(),
}));

const mockMetric: CanonicalMetric = {
  id: 'test-metric-id',
  metric_name: 'http.server.request.duration',
  instrument_type: 'histogram',
  description: 'Duration of HTTP server requests',
  unit: 'ms',
  attributes: [],
  enabled_by_default: true,
  component_type: 'receiver',
  component_name: 'httpreceiver',
  source_category: 'otel-collector-contrib',
  source_name: 'metadata.yaml',
  source_location: '/receiver/httpreceiver',
  extraction_method: 'yaml',
  source_confidence: 'high',
  repo: 'https://github.com/open-telemetry/opentelemetry-collector-contrib',
  path: '/receiver/httpreceiver/metadata.yaml',
  commit: 'abc123',
  extracted_at: '2024-01-01T00:00:00Z',
};

const mockMetric2: CanonicalMetric = {
  ...mockMetric,
  id: 'test-metric-id-2',
  metric_name: 'http.client.request.duration',
};

describe('Home Page - Metric Detail Interactions', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockSearchParams.delete('metric');
    mockSearchParams.delete('q');

    (api.searchMetrics as Mock).mockResolvedValue({
      metrics: [mockMetric, mockMetric2],
      total: 2,
    });
    (api.getFacets as Mock).mockResolvedValue({
      instrument_types: [],
      component_types: [],
      component_names: [],
      source_categories: [],
      source_names: [],
    });
    (api.getMetric as Mock).mockResolvedValue(mockMetric);

    Object.defineProperty(window, 'location', {
      value: { search: '', pathname: '/' },
      writable: true,
    });
    window.history.replaceState = vi.fn();

    const localStorageMock = {
      getItem: vi.fn().mockReturnValue(null),
      setItem: vi.fn(),
      removeItem: vi.fn(),
      clear: vi.fn(),
    };
    Object.defineProperty(window, 'localStorage', {
      value: localStorageMock,
      writable: true,
    });

    Object.defineProperty(window, 'matchMedia', {
      value: vi.fn().mockImplementation((query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
      })),
      writable: true,
    });
  });

  it('renders metric list on load', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
      expect(screen.getByText('http.client.request.duration')).toBeInTheDocument();
    });

    expect(api.searchMetrics).toHaveBeenCalledTimes(1);
  });

  it('opens metric detail when clicking a metric card', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByText('Details')).toBeInTheDocument();
    });
  });

  it('does not call getMetric API when opening from list (data already available)', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });

    expect(api.getMetric).not.toHaveBeenCalled();
  });

  it('updates URL when opening metric detail', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(window.history.replaceState).toHaveBeenCalled();
      const lastCall = (window.history.replaceState as Mock).mock.calls.slice(-1)[0];
      expect(lastCall[2]).toContain('metric=test-metric-id');
    });
  });

  it('closes metric detail when clicking close button', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByLabelText('Close'));

    await waitFor(() => {
      expect(screen.queryByText('Description')).not.toBeInTheDocument();
    });
  });

  it('closes metric detail when pressing ESC key', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });

    fireEvent.keyDown(window, { key: 'Escape' });

    await waitFor(() => {
      expect(screen.queryByText('Description')).not.toBeInTheDocument();
    });
  });

  it('updates URL when closing metric detail', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });

    (window.history.replaceState as Mock).mockClear();
    fireEvent.click(screen.getByLabelText('Close'));

    await waitFor(() => {
      expect(window.history.replaceState).toHaveBeenCalled();
      const lastCall = (window.history.replaceState as Mock).mock.calls.slice(-1)[0];
      expect(lastCall[2]).not.toContain('metric=');
    });
  });

  it('does not refetch metrics list when opening metric detail', async () => {
    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(screen.getByText('http.server.request.duration')).toBeInTheDocument();
    });

    const initialCallCount = (api.searchMetrics as Mock).mock.calls.length;

    fireEvent.click(screen.getByText('http.server.request.duration'));

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
    });

    expect(api.searchMetrics).toHaveBeenCalledTimes(initialCallCount);
  });

  it('fetches metric from API when loading page with metric in URL', async () => {
    mockSearchParams.set('metric', 'test-metric-id');

    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(api.getMetric).toHaveBeenCalledWith('test-metric-id');
    });
  });

  it('shows metric detail when loading page with metric in URL', async () => {
    mockSearchParams.set('metric', 'test-metric-id');

    renderWithProviders(<Home />);

    await waitFor(() => {
      expect(api.getMetric).toHaveBeenCalledWith('test-metric-id');
    });

    await waitFor(() => {
      expect(screen.getByText('Description')).toBeInTheDocument();
      expect(screen.getByText('Details')).toBeInTheDocument();
      expect(screen.getByText('Source Information')).toBeInTheDocument();
    });
  });
});
