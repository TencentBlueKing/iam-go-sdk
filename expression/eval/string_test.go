package eval

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("String", func() {

	It("convertArgsToString", func() {
		_, _, ok := convertArgsToString(1, "")
		assert.False(GinkgoT(), ok)
		_, _, ok = convertArgsToString("", 1)
		assert.False(GinkgoT(), ok)
		_, _, ok = convertArgsToString("", "")
		assert.True(GinkgoT(), ok)
	})

	It("StartsWith", func() {
		assert.False(GinkgoT(), StartsWith("", 1))
		assert.True(GinkgoT(), StartsWith("hello", "he"))
	})

	It("NotStartsWith", func() {
		assert.False(GinkgoT(), NotStartsWith("", 1))
		assert.False(GinkgoT(), NotStartsWith("hello", "he"))
		assert.True(GinkgoT(), NotStartsWith("hello", "abc"))
	})

	It("EndsWith", func() {
		assert.False(GinkgoT(), EndsWith("", 1))
		assert.True(GinkgoT(), EndsWith("hello", "lo"))
	})

	It("NotEndsWith", func() {
		assert.False(GinkgoT(), NotEndsWith("", 1))
		assert.False(GinkgoT(), NotEndsWith("hello", "lo"))
		assert.True(GinkgoT(), NotEndsWith("hello", "abc"))
	})

})
