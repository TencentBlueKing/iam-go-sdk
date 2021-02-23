package client_test

import (
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/onsi/ginkgo"

	"github.com/TencentBlueKing/iam-go-sdk/client"
)

var _ = Describe("Metric", func() {

	It("NewMetricCallback", func() {

		f := client.NewMetricCallback("test", time.Now())

		assert.NotNil(GinkgoT(), f)

	})

})
