/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云PaaS平台社区版 (BlueKing PaaS
 * Community Edition) available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
