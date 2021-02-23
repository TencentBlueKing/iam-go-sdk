package metric_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMetric(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metric Suite")
}
