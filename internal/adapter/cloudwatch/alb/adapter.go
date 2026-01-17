package alb

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
	return "cloudwatch-alb"
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
	return "https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-cloudwatch-metrics.html"
}

func (a *Adapter) Fetch(_ context.Context, _ adapter.FetchOptions) (*adapter.FetchResult, error) {
	return &adapter.FetchResult{
		Commit:    time.Now().Format("2006-01-02"),
		Timestamp: time.Now(),
	}, nil
}

func (a *Adapter) Extract(_ context.Context, _ *adapter.FetchResult) ([]*adapter.RawMetric, error) {
	return []*adapter.RawMetric{
		// Load Balancer Metrics
		{Name: "ActiveConnectionCount", Description: "Total number of concurrent TCP connections active from clients to the load balancer and from the load balancer to targets", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ClientTLSNegotiationErrorCount", Description: "Number of TLS connections initiated by the client that did not establish a session with the load balancer due to a TLS error", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "DroppedInvalidHeaderRequestCount", Description: "Number of requests where the load balancer removed HTTP headers with invalid header fields", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ForwardedInvalidHeaderRequestCount", Description: "Number of requests routed by the load balancer that had HTTP headers with invalid header fields", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "GrpcRequestCount", Description: "Number of gRPC requests processed over IPv4 and IPv6", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTP_Fixed_Response_Count", Description: "Number of fixed-response actions that were successful", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTP_Redirect_Count", Description: "Number of redirect actions that were successful", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTP_Redirect_Url_Limit_Exceeded_Count", Description: "Number of redirect actions that couldn't be completed because the URL exceeds 8K", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_3XX_Count", Description: "Number of HTTP 3XX redirection codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_4XX_Count", Description: "Number of HTTP 4XX client error codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_5XX_Count", Description: "Number of HTTP 5XX server error codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_500_Count", Description: "Number of HTTP 500 error codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_502_Count", Description: "Number of HTTP 502 error codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_503_Count", Description: "Number of HTTP 503 error codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_ELB_504_Count", Description: "Number of HTTP 504 error codes that originate from the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "IPv6ProcessedBytes", Description: "Total number of bytes processed by the load balancer over IPv6", Unit: "Bytes", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "IPv6RequestCount", Description: "Number of IPv6 requests received by the load balancer", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "NewConnectionCount", Description: "Total number of new TCP connections established from clients to the load balancer and from the load balancer to targets", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "NonStickyRequestCount", Description: "Number of requests where the load balancer chose a new target because it couldn't use an existing sticky session", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ProcessedBytes", Description: "Total number of bytes processed by the load balancer over IPv4 and IPv6", Unit: "Bytes", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "RejectedConnectionCount", Description: "Number of connections rejected because the load balancer had reached its maximum number of connections", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "RequestCount", Description: "Number of requests processed over IPv4 and IPv6", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "RuleEvaluations", Description: "Number of rules evaluated by the load balancer while processing requests", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		// LCU Metrics
		{Name: "ConsumedLCUs", Description: "Number of load balancer capacity units (LCU) used by your load balancer", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "PeakLCUs", Description: "Maximum number of load balancer capacity units (LCU) used at a given point in time", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		// Target Metrics
		{Name: "AnomalousHostCount", Description: "Number of hosts detected with anomalies", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HealthyHostCount", Description: "Number of targets that are considered healthy", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_Target_2XX_Count", Description: "Number of HTTP 2XX response codes generated by the targets", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_Target_3XX_Count", Description: "Number of HTTP 3XX response codes generated by the targets", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_Target_4XX_Count", Description: "Number of HTTP 4XX response codes generated by the targets", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HTTPCode_Target_5XX_Count", Description: "Number of HTTP 5XX response codes generated by the targets", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "MitigatedHostCount", Description: "Number of targets under mitigation", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "RequestCountPerTarget", Description: "Average request count per target in a target group", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "TargetConnectionErrorCount", Description: "Number of connections not successfully established between load balancer and target", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "TargetResponseTime", Description: "Time elapsed after request leaves load balancer until target starts to send response headers", Unit: "Seconds", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "TargetTLSNegotiationErrorCount", Description: "Number of TLS connections initiated by load balancer that did not establish a session with the target", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "UnHealthyHostCount", Description: "Number of targets that are considered unhealthy", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		// Target Group Health Metrics
		{Name: "HealthyStateDNS", Description: "Number of zones that meet the DNS healthy state requirements", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "HealthyStateRouting", Description: "Number of zones that meet the routing healthy state requirements", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "UnhealthyRoutingRequestCount", Description: "Number of requests routed using the routing failover action", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "UnhealthyStateDNS", Description: "Number of zones that do not meet the DNS healthy state requirements", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "UnhealthyStateRouting", Description: "Number of zones that do not meet the routing healthy state requirements", Unit: "Count", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		// Lambda Metrics
		{Name: "LambdaInternalError", Description: "Number of requests to Lambda function that failed due to internal load balancer or AWS Lambda issue", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "LambdaTargetProcessedBytes", Description: "Total number of bytes processed by load balancer for requests to and responses from Lambda function", Unit: "Bytes", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "LambdaUserError", Description: "Number of requests to Lambda function that failed due to issue with the Lambda function", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		// Auth Metrics
		{Name: "ELBAuthError", Description: "Number of user authentications that could not be completed due to misconfiguration or internal error", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ELBAuthFailure", Description: "Number of user authentications that could not be completed because IdP denied access", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ELBAuthLatency", Description: "Time elapsed to query the IdP for ID token and user info", Unit: "Milliseconds", InstrumentType: "gauge", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ELBAuthRefreshTokenSuccess", Description: "Number of times load balancer successfully refreshed user claims using refresh token", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ELBAuthSuccess", Description: "Number of authenticate actions that were successful", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
		{Name: "ELBAuthUserClaimsSizeExceeded", Description: "Number of times configured IdP returned user claims exceeding 11K bytes", Unit: "Count", InstrumentType: "counter", ComponentName: "ALB", ComponentType: "platform", SourceLocation: "AWS/ApplicationELB"},
	}, nil
}
