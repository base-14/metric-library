import React, { useState, useMemo } from 'react';

// Sample data representing what we'd extract from OTel repos
const SAMPLE_METRICS = [
  {
    id: '1',
    name: 'system.cpu.utilization',
    type: 'gauge',
    description: 'CPU utilization as a percentage of total available CPU time',
    unit: '1',
    component_type: 'receiver',
    component_name: 'hostmetrics',
    enabled_by_default: true,
    attributes: [
      { name: 'cpu', type: 'string', description: 'CPU number starting at 0' },
      { name: 'state', type: 'string', description: 'CPU state', enum: ['idle', 'user', 'system', 'interrupt', 'nice', 'softirq', 'steal', 'wait'] }
    ]
  },
  {
    id: '2',
    name: 'system.memory.usage',
    type: 'sum',
    description: 'Memory usage in bytes',
    unit: 'By',
    component_type: 'receiver',
    component_name: 'hostmetrics',
    enabled_by_default: true,
    attributes: [
      { name: 'state', type: 'string', description: 'Memory state', enum: ['used', 'free', 'buffered', 'cached', 'slab_reclaimable', 'slab_unreclaimable'] }
    ]
  },
  {
    id: '3',
    name: 'system.disk.io',
    type: 'sum',
    description: 'Disk bytes transferred',
    unit: 'By',
    component_type: 'receiver',
    component_name: 'hostmetrics',
    enabled_by_default: true,
    attributes: [
      { name: 'device', type: 'string', description: 'Name of the disk device' },
      { name: 'direction', type: 'string', description: 'Direction of transfer', enum: ['read', 'write'] }
    ]
  },
  {
    id: '4',
    name: 'http.server.request.duration',
    type: 'histogram',
    description: 'Duration of HTTP server requests',
    unit: 's',
    component_type: 'receiver',
    component_name: 'httpcheck',
    enabled_by_default: true,
    attributes: [
      { name: 'http.request.method', type: 'string', description: 'HTTP request method' },
      { name: 'http.response.status_code', type: 'int', description: 'HTTP response status code' },
      { name: 'url.scheme', type: 'string', description: 'HTTP or HTTPS' }
    ]
  },
  {
    id: '5',
    name: 'kafka.consumer.lag',
    type: 'gauge',
    description: 'Current approximate lag of consumer group behind the head of partition',
    unit: '{messages}',
    component_type: 'receiver',
    component_name: 'kafkametrics',
    enabled_by_default: true,
    attributes: [
      { name: 'group', type: 'string', description: 'Consumer group name' },
      { name: 'topic', type: 'string', description: 'Topic name' },
      { name: 'partition', type: 'int', description: 'Partition number' }
    ]
  },
  {
    id: '6',
    name: 'kafka.brokers',
    type: 'gauge',
    description: 'Number of brokers in the cluster',
    unit: '{brokers}',
    component_type: 'receiver',
    component_name: 'kafkametrics',
    enabled_by_default: true,
    attributes: []
  },
  {
    id: '7',
    name: 'redis.commands.processed',
    type: 'sum',
    description: 'Total number of commands processed by the server',
    unit: '{commands}',
    component_type: 'receiver',
    component_name: 'redis',
    enabled_by_default: true,
    attributes: []
  },
  {
    id: '8',
    name: 'redis.memory.used',
    type: 'gauge',
    description: 'Total memory used by Redis',
    unit: 'By',
    component_type: 'receiver',
    component_name: 'redis',
    enabled_by_default: true,
    attributes: []
  },
  {
    id: '9',
    name: 'process.cpu.utilization',
    type: 'gauge',
    description: 'CPU utilization of the process',
    unit: '1',
    component_type: 'receiver',
    component_name: 'hostmetrics',
    enabled_by_default: false,
    attributes: []
  },
  {
    id: '10',
    name: 'db.client.connections.usage',
    type: 'gauge',
    description: 'Number of connections currently in use by the database client',
    unit: '{connections}',
    component_type: 'receiver',
    component_name: 'postgresql',
    enabled_by_default: true,
    attributes: [
      { name: 'state', type: 'string', description: 'Connection state', enum: ['idle', 'active'] },
      { name: 'database', type: 'string', description: 'Database name' }
    ]
  }
];

const TYPE_COLORS = {
  gauge: { bg: 'bg-blue-100', text: 'text-blue-800', border: 'border-blue-200' },
  sum: { bg: 'bg-green-100', text: 'text-green-800', border: 'border-green-200' },
  counter: { bg: 'bg-green-100', text: 'text-green-800', border: 'border-green-200' },
  histogram: { bg: 'bg-purple-100', text: 'text-purple-800', border: 'border-purple-200' },
};

const COMPONENT_TYPES = ['receiver', 'processor', 'exporter', 'extension', 'connector'];
const METRIC_TYPES = ['gauge', 'sum', 'histogram'];

// Get unique component names from data
const getComponentNames = (metrics) => [...new Set(metrics.map(m => m.component_name))].sort();

function TypeBadge({ type }) {
  const colors = TYPE_COLORS[type] || TYPE_COLORS.gauge;
  return (
    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${colors.bg} ${colors.text}`}>
      {type}
    </span>
  );
}

function CopyButton({ text }) {
  const [copied, setCopied] = useState(false);
  
  const handleCopy = async (e) => {
    e.stopPropagation();
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };
  
  return (
    <button
      onClick={handleCopy}
      className="p-1 rounded hover:bg-gray-200 transition-colors"
      title="Copy metric name"
    >
      {copied ? (
        <svg className="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
        </svg>
      ) : (
        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      )}
    </button>
  );
}

function MetricCard({ metric, onClick, isSelected }) {
  return (
    <div 
      className={`p-4 border rounded-lg cursor-pointer transition-all hover:border-gray-400 hover:shadow-sm ${isSelected ? 'border-blue-500 bg-blue-50' : 'border-gray-200 bg-white'}`}
      onClick={() => onClick(metric)}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <code className="text-sm font-semibold text-gray-900 truncate">{metric.name}</code>
            <CopyButton text={metric.name} />
          </div>
          <p className="text-sm text-gray-600 line-clamp-2">{metric.description}</p>
        </div>
        <TypeBadge type={metric.type} />
      </div>
      <div className="flex items-center gap-2 mt-3 text-xs text-gray-500">
        <span className="inline-flex items-center px-2 py-0.5 rounded bg-gray-100">
          {metric.component_name}
        </span>
        {metric.attributes.length > 0 && (
          <span>{metric.attributes.length} attribute{metric.attributes.length !== 1 ? 's' : ''}</span>
        )}
        {metric.unit && metric.unit !== '1' && (
          <span>Unit: {metric.unit}</span>
        )}
      </div>
    </div>
  );
}

function MetricDetail({ metric, onClose }) {
  if (!metric) return null;
  
  return (
    <div className="bg-white border-l border-gray-200 h-full overflow-y-auto">
      <div className="sticky top-0 bg-white border-b border-gray-200 p-4 flex items-center justify-between">
        <h2 className="font-semibold text-gray-900">Metric Details</h2>
        <button onClick={onClose} className="p-1 rounded hover:bg-gray-100">
          <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <div className="p-4 space-y-6">
        {/* Name */}
        <div>
          <div className="flex items-center gap-2">
            <code className="text-lg font-semibold text-gray-900 break-all">{metric.name}</code>
            <CopyButton text={metric.name} />
          </div>
        </div>
        
        {/* Type & Unit */}
        <div className="flex gap-4">
          <div>
            <label className="text-xs font-medium text-gray-500 uppercase tracking-wide">Type</label>
            <div className="mt-1">
              <TypeBadge type={metric.type} />
            </div>
          </div>
          {metric.unit && (
            <div>
              <label className="text-xs font-medium text-gray-500 uppercase tracking-wide">Unit</label>
              <p className="mt-1 text-sm font-mono text-gray-900">{metric.unit}</p>
            </div>
          )}
          <div>
            <label className="text-xs font-medium text-gray-500 uppercase tracking-wide">Enabled</label>
            <p className="mt-1 text-sm text-gray-900">{metric.enabled_by_default ? 'Yes' : 'No'}</p>
          </div>
        </div>
        
        {/* Description */}
        <div>
          <label className="text-xs font-medium text-gray-500 uppercase tracking-wide">Description</label>
          <p className="mt-1 text-sm text-gray-700">{metric.description}</p>
        </div>
        
        {/* Source */}
        <div>
          <label className="text-xs font-medium text-gray-500 uppercase tracking-wide">Source</label>
          <div className="mt-1 flex items-center gap-2">
            <span className="inline-flex items-center px-2 py-1 rounded bg-gray-100 text-sm">
              {metric.component_type}/{metric.component_name}
            </span>
            <a href="#" className="text-sm text-blue-600 hover:underline">View on GitHub →</a>
          </div>
        </div>
        
        {/* Attributes */}
        {metric.attributes.length > 0 && (
          <div>
            <label className="text-xs font-medium text-gray-500 uppercase tracking-wide">Attributes</label>
            <div className="mt-2 space-y-2">
              {metric.attributes.map((attr, idx) => (
                <div key={idx} className="p-3 bg-gray-50 rounded-lg">
                  <div className="flex items-center gap-2">
                    <code className="text-sm font-medium text-gray-900">{attr.name}</code>
                    <span className="text-xs text-gray-500">({attr.type})</span>
                  </div>
                  <p className="text-sm text-gray-600 mt-1">{attr.description}</p>
                  {attr.enum && (
                    <div className="mt-2 flex flex-wrap gap-1">
                      {attr.enum.map((val, i) => (
                        <code key={i} className="text-xs px-1.5 py-0.5 bg-white border border-gray-200 rounded">{val}</code>
                      ))}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function FilterSection({ title, options, selected, onChange, counts = {} }) {
  return (
    <div className="mb-4">
      <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wide mb-2">{title}</h3>
      <div className="space-y-1">
        {options.map(option => (
          <label key={option} className="flex items-center gap-2 cursor-pointer py-1 px-2 rounded hover:bg-gray-100">
            <input
              type="checkbox"
              checked={selected.includes(option)}
              onChange={(e) => {
                if (e.target.checked) {
                  onChange([...selected, option]);
                } else {
                  onChange(selected.filter(s => s !== option));
                }
              }}
              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
            />
            <span className="text-sm text-gray-700 flex-1">{option}</span>
            {counts[option] !== undefined && (
              <span className="text-xs text-gray-400">{counts[option]}</span>
            )}
          </label>
        ))}
      </div>
    </div>
  );
}

export default function OTelGlossary() {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedMetric, setSelectedMetric] = useState(null);
  const [selectedTypes, setSelectedTypes] = useState([]);
  const [selectedComponents, setSelectedComponents] = useState([]);
  const [showFilters, setShowFilters] = useState(true);
  
  const componentNames = useMemo(() => getComponentNames(SAMPLE_METRICS), []);
  
  // Compute counts for filters
  const typeCounts = useMemo(() => {
    return METRIC_TYPES.reduce((acc, type) => {
      acc[type] = SAMPLE_METRICS.filter(m => m.type === type).length;
      return acc;
    }, {});
  }, []);
  
  const componentCounts = useMemo(() => {
    return componentNames.reduce((acc, name) => {
      acc[name] = SAMPLE_METRICS.filter(m => m.component_name === name).length;
      return acc;
    }, {});
  }, [componentNames]);
  
  // Filter and search metrics
  const filteredMetrics = useMemo(() => {
    return SAMPLE_METRICS.filter(metric => {
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        const matchesName = metric.name.toLowerCase().includes(query);
        const matchesDesc = metric.description.toLowerCase().includes(query);
        const matchesComponent = metric.component_name.toLowerCase().includes(query);
        if (!matchesName && !matchesDesc && !matchesComponent) return false;
      }
      
      // Type filter
      if (selectedTypes.length > 0 && !selectedTypes.includes(metric.type)) {
        return false;
      }
      
      // Component filter
      if (selectedComponents.length > 0 && !selectedComponents.includes(metric.component_name)) {
        return false;
      }
      
      return true;
    });
  }, [searchQuery, selectedTypes, selectedComponents]);
  
  const clearFilters = () => {
    setSelectedTypes([]);
    setSelectedComponents([]);
  };
  
  const hasActiveFilters = selectedTypes.length > 0 || selectedComponents.length > 0;
  
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex items-center gap-3 mb-4">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <svg className="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <h1 className="text-2xl font-bold text-gray-900 font-mono">metric-library</h1>
              <p className="text-sm text-gray-500">Metric discovery platform</p>
            </div>
          </div>
          
          {/* Search Bar */}
          <div className="relative max-w-2xl">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              placeholder="Search metrics by name, description, or component..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="block w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg bg-white text-gray-900 placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg className="h-5 w-5 text-gray-400 hover:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            )}
          </div>
          
          {/* Quick Stats */}
          <div className="flex items-center gap-4 mt-4 text-sm text-gray-500">
            <span>{SAMPLE_METRICS.length} total metrics</span>
            <span>•</span>
            <span>{componentNames.length} components</span>
            <span>•</span>
            <span>Updated Jan 10, 2025</span>
          </div>
        </div>
      </header>
      
      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 py-6">
        <div className="flex gap-6">
          {/* Sidebar Filters */}
          {showFilters && (
            <aside className="w-56 flex-shrink-0">
              <div className="sticky top-6">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="font-semibold text-gray-900">Filters</h2>
                  {hasActiveFilters && (
                    <button onClick={clearFilters} className="text-xs text-blue-600 hover:underline">
                      Clear all
                    </button>
                  )}
                </div>
                
                <FilterSection
                  title="Metric Type"
                  options={METRIC_TYPES}
                  selected={selectedTypes}
                  onChange={setSelectedTypes}
                  counts={typeCounts}
                />
                
                <FilterSection
                  title="Component"
                  options={componentNames}
                  selected={selectedComponents}
                  onChange={setSelectedComponents}
                  counts={componentCounts}
                />
              </div>
            </aside>
          )}
          
          {/* Results */}
          <main className="flex-1 min-w-0">
            {/* Results Header */}
            <div className="flex items-center justify-between mb-4">
              <p className="text-sm text-gray-600">
                {filteredMetrics.length} metric{filteredMetrics.length !== 1 ? 's' : ''} found
                {searchQuery && <span> for "<strong>{searchQuery}</strong>"</span>}
              </p>
              <button
                onClick={() => setShowFilters(!showFilters)}
                className="text-sm text-gray-500 hover:text-gray-700 flex items-center gap-1"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
                </svg>
                {showFilters ? 'Hide' : 'Show'} filters
              </button>
            </div>
            
            {/* Active Filters */}
            {hasActiveFilters && (
              <div className="flex flex-wrap gap-2 mb-4">
                {selectedTypes.map(type => (
                  <span key={type} className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-blue-100 text-blue-800 text-sm">
                    {type}
                    <button onClick={() => setSelectedTypes(selectedTypes.filter(t => t !== type))}>
                      <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </span>
                ))}
                {selectedComponents.map(comp => (
                  <span key={comp} className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-gray-200 text-gray-800 text-sm">
                    {comp}
                    <button onClick={() => setSelectedComponents(selectedComponents.filter(c => c !== comp))}>
                      <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </span>
                ))}
              </div>
            )}
            
            {/* Metrics Grid */}
            <div className={`grid gap-3 ${selectedMetric ? 'grid-cols-1' : 'grid-cols-1 lg:grid-cols-2'}`}>
              {filteredMetrics.map(metric => (
                <MetricCard
                  key={metric.id}
                  metric={metric}
                  onClick={setSelectedMetric}
                  isSelected={selectedMetric?.id === metric.id}
                />
              ))}
            </div>
            
            {filteredMetrics.length === 0 && (
              <div className="text-center py-12">
                <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <h3 className="mt-2 text-sm font-medium text-gray-900">No metrics found</h3>
                <p className="mt-1 text-sm text-gray-500">Try adjusting your search or filters</p>
                <button
                  onClick={() => { setSearchQuery(''); clearFilters(); }}
                  className="mt-4 text-sm text-blue-600 hover:underline"
                >
                  Clear all filters
                </button>
              </div>
            )}
          </main>
          
          {/* Detail Panel */}
          {selectedMetric && (
            <aside className="w-96 flex-shrink-0">
              <div className="sticky top-6">
                <MetricDetail metric={selectedMetric} onClose={() => setSelectedMetric(null)} />
              </div>
            </aside>
          )}
        </div>
      </div>
      
      {/* Footer */}
      <footer className="border-t border-gray-200 bg-white mt-12">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex items-center justify-between text-sm text-gray-500">
            <p>Data sourced from <a href="https://github.com/open-telemetry/opentelemetry-collector-contrib" className="text-blue-600 hover:underline">opentelemetry-collector-contrib</a></p>
            <div className="flex gap-4">
              <a href="#" className="hover:text-gray-700">GitHub</a>
              <a href="#" className="hover:text-gray-700">API</a>
              <a href="#" className="hover:text-gray-700">About</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}