package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	serviceName = "iam"
)

var (
	// ClientRequestDuration 依赖 api 响应时间分布
	ClientRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "client_request_duration_milliseconds",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     []float64{20, 50, 100, 200, 500, 1000, 2000, 5000},
	},
		[]string{"method", "path", "status", "component"},
	)
)

// RegisterMetrics will register the mtrics
func RegisterMetrics() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(ClientRequestDuration)
}
