# Todo

## In Progress: UI Improvements

- [ ] Search as you type with debounce (300ms)
  - `web/src/components/SearchBar.tsx`
  - `web/src/lib/debounce.ts` (new)

- [ ] Reorder filters (Component Name first, Instrument Type second)
  - `web/src/components/FilterPanel.tsx`

## Pending

- [ ] Manual E2E testing
- [ ] Helm chart for Kubernetes deployment

## Completed

- [x] Add CLI command to trigger metric extraction
- [x] Run extraction against otel-collector-contrib (1261 metrics)
- [x] Verify data loads in frontend
