package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Utils", func() {

	Describe("Int64ArrayToString", func() {

		DescribeTable("cases", func(expected string, array []int64) {
			assert.Equal(GinkgoT(), expected, Int64ArrayToString(array, ","))
		},
			Entry("empty", "", []int64{}),
			Entry("one int", "1", []int64{1}),
			Entry("three ints", "1,2,3", []int64{1, 2, 3}),
		)
	})

})
