package client

import (
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/TencentBlueKing/iam-go-sdk/metric"
)

// CallbackFunc is the func object of http callback
type CallbackFunc func(response gorequest.Response, v interface{}, body []byte, errs []error)

// NewMetricCallback will record the http request data into metrics
func NewMetricCallback(system string, start time.Time) CallbackFunc {
	return func(response gorequest.Response, v interface{}, body []byte, errs []error) {
		duration := time.Since(start)

		metric.ClientRequestDuration.With(prometheus.Labels{
			"method":    response.Request.Method,
			"path":      response.Request.URL.Path,
			"status":    strconv.Itoa(response.StatusCode),
			"component": system,
		}).Observe(float64(duration / time.Millisecond))
	}
}
