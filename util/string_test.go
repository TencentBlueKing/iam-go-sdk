package util

import (
	"math/rand"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"
)

func rawBytesToStr(b []byte) string {
	return string(b)
}

func rawStrToBytes(s string) []byte {
	return []byte(s)
}

const letterBytesForTest = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandStringBytesMaskImprSrcSB will generate a n-chars string
func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytesForTest) {
			sb.WriteByte(letterBytesForTest[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

var _ = Describe("Utils", func() {

	Describe("TruncateString", func() {
		var s = "helloworld"

		DescribeTable("TruncateString cases", func(expected string, truncatedSize int) {
			assert.Equal(GinkgoT(), expected, TruncateString(s, truncatedSize))
		},
			Entry("truncated size less than real size", "he", 2),
			Entry("truncated size equals to real size", s, 10),
			Entry("truncated size greater than real size", s, 20),
		)
	})

	Describe("StringToBytes", func() {
		It("a normal string", func() {
			b := StringToBytes("abc")
			assert.Equal(GinkgoT(), []byte("abc"), b)
		})

		It("random generate", func() {
			for i := 0; i < 100; i++ {
				s := RandStringBytesMaskImprSrcSB(64)

				assert.Equal(GinkgoT(), rawStrToBytes(s), StringToBytes(s))
			}
		})
	})

	Describe("BytesToString", func() {
		It("a normal bytes", func() {
			s := BytesToString([]byte("abc"))
			assert.Equal(GinkgoT(), "abc", s)
		})

		It("random generate", func() {
			data := make([]byte, 1024)
			for i := 0; i < 100; i++ {
				rand.Read(data)
				assert.Equal(GinkgoT(), rawBytesToStr(data), BytesToString(data))
			}
		})
	})

})
