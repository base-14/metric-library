package apigateway

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
	return "cloudwatch-apigateway"
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
	return "https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-metrics-and-dimensions.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		{Name: "4XXError", Description: "The number of client-side errors captured in a given period", Unit: "Count", InstrumentType: "counter", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
		{Name: "5XXError", Description: "The number of server-side errors captured in a given period", Unit: "Count", InstrumentType: "counter", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
		{Name: "CacheHitCount", Description: "The number of requests served from the API cache in a given period", Unit: "Count", InstrumentType: "counter", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
		{Name: "CacheMissCount", Description: "The number of requests served from the backend in a given period when API caching is enabled", Unit: "Count", InstrumentType: "counter", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
		{Name: "Count", Description: "The total number of API requests in a given period", Unit: "Count", InstrumentType: "counter", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
		{Name: "IntegrationLatency", Description: "The time between when API Gateway relays a request to the backend and when it receives a response from the backend", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
		{Name: "Latency", Description: "The time between when API Gateway receives a request from a client and when it returns a response to the client", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "APIGateway", ComponentType: "platform", SourceLocation: "AWS/ApiGateway"},
	}, nil
}
