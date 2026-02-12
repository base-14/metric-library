package loadbalancing

import (
	"context"
	"time"

	"github.com/base-14/metric-library/internal/adapter"
	"github.com/base-14/metric-library/internal/domain"
)

type Adapter struct{}

func NewAdapter(_ string) *Adapter {
	return &Adapter{}
}

func (a *Adapter) Name() string {
	return "gcp-loadbalancing"
}

func (a *Adapter) SourceCategory() domain.SourceCategory {
	return domain.SourceCloud
}

func (a *Adapter) Confidence() domain.ConfidenceLevel {
	return domain.ConfidenceDocumented
}

func (a *Adapter) ExtractionMethod() domain.ExtractionMethod {
	return domain.ExtractionScrape
}

func (a *Adapter) RepoURL() string {
	return "https://cloud.google.com/monitoring/api/metrics_gcp#gcp-loadbalancing"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	var metrics []*adapter.RawMetric
	metrics = append(metrics, httpMetrics()...)
	metrics = append(metrics, tcpSslMetrics()...)
	return metrics, nil
}

func httpMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "loadbalancing.googleapis.com/https/request_count", Description: "Number of requests served by the HTTP(S) load balancer", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/request_bytes_count", Description: "Number of bytes sent as requests from clients to the HTTP(S) load balancer", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/response_bytes_count", Description: "Number of bytes sent as responses from the HTTP(S) load balancer to clients", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/total_latencies", Description: "Distribution of latency calculated from when the request was received by the load balancer proxy to when the proxy received ACK from the client on the last response byte", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/backend_latencies", Description: "Distribution of latency calculated from when the request was sent by the load balancer proxy to the backend until the proxy received from the backend the last byte of response", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/backend_request_count", Description: "Number of requests sent from the HTTP(S) load balancer to the backends", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/backend_request_bytes_count", Description: "Number of bytes sent as requests from the HTTP(S) load balancer to the backends", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/backend_response_bytes_count", Description: "Number of bytes sent as responses from the backends to the HTTP(S) load balancer", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/external_regional/total_latencies", Description: "Distribution of latency for regional external HTTP(S) load balancer", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/external_regional/backend_latencies", Description: "Distribution of backend latency for regional external HTTP(S) load balancer", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/external_regional/request_count", Description: "Number of requests served by the regional external HTTP(S) load balancer", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/external_regional/request_bytes_count", Description: "Number of bytes sent as requests from clients to the regional external HTTP(S) load balancer", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/external_regional/response_bytes_count", Description: "Number of bytes sent as responses from the regional external HTTP(S) load balancer to clients", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/internal/total_latencies", Description: "Distribution of latency for internal HTTP(S) load balancer", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/internal/backend_latencies", Description: "Distribution of backend latency for internal HTTP(S) load balancer", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/internal/request_count", Description: "Number of requests served by the internal HTTP(S) load balancer", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/internal/request_bytes_count", Description: "Number of bytes sent as requests from clients to the internal HTTP(S) load balancer", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/https/internal/response_bytes_count", Description: "Number of bytes sent as responses from the internal HTTP(S) load balancer to clients", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
	}
}

func tcpSslMetrics() []*adapter.RawMetric {
	return []*adapter.RawMetric{
		{Name: "loadbalancing.googleapis.com/tcp_ssl_proxy/open_connections", Description: "Number of connections that are open at the current moment", Unit: "1", InstrumentType: "gauge", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/tcp_ssl_proxy/new_connections", Description: "Number of connections that were created (client successfully connected to backend)", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/tcp_ssl_proxy/closed_connections", Description: "Number of connections that were terminated", Unit: "1", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/tcp_ssl_proxy/ingress_bytes_count", Description: "Number of bytes sent from client to backend using the proxy", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/tcp_ssl_proxy/egress_bytes_count", Description: "Number of bytes sent from backend to client using the proxy", Unit: "By", InstrumentType: "counter", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
		{Name: "loadbalancing.googleapis.com/tcp_ssl_proxy/frontend_tcp_rtt", Description: "Distribution of smoothed RTT measured for each connection between client and the proxy", Unit: "ms", InstrumentType: "histogram", ComponentName: "Cloud Load Balancing", ComponentType: "platform", SourceLocation: "loadbalancing.googleapis.com"},
	}
}
