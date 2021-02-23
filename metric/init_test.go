package metric_test

import (
	. "github.com/onsi/ginkgo"

	"github.com/TencentBlueKing/iam-go-sdk/metric"
)

var _ = Describe("Init", func() {
	It("RegisterMetrics", func() {
		metric.RegisterMetrics()
	})
})
