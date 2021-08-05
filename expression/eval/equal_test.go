/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云-权限中心Go SDK(iam-go-sdk) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

/* MIT License
 * Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/* NOTE: copied from https://github.com/stretchr/testify/assert/assertions_test.go and modified
 *  The original versions of the files are MIT licensed
 */

package eval

import (
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

// AssertionTesterConformingObject is an object that conforms to the AssertionTesterInterface interface
type AssertionTesterConformingObject struct {
}

func (a *AssertionTesterConformingObject) TestMethod() {
}

var _ = Describe("Operator", func() {

	It("Equal", func() {
		type myType string

		var m map[string]interface{}

		cases := []struct {
			expected interface{}
			actual   interface{}
			result   bool
			remark   string
		}{
			{"Hello World", "Hello World", true, ""},
			{123, 123, true, ""},
			{123.5, 123.5, true, ""},
			{[]byte("Hello World"), []byte("Hello World"), true, ""},
			{nil, nil, true, ""},
			{int32(123), int32(123), true, ""},
			{uint64(123), uint64(123), true, ""},
			{myType("1"), myType("1"), true, ""},
			{&struct{}{}, &struct{}{}, true, "pointer equality is based on equality of underlying value"},

			// Not expected to be equal
			{m["bar"], "something", false, ""},
			{myType("1"), myType("2"), false, ""},

			// A case that might be confusing, especially with numeric literals
			{10, uint(10), false, ""},
		}

		for _, c := range cases {
			res := Equal(c.expected, c.actual)
			assert.Equal(GinkgoT(), c.result, res)
		}

	})

	It("NotEqual", func() {
		cases := []struct {
			expected interface{}
			actual   interface{}
			result   bool
		}{
			// cases that are expected not to match
			{"Hello World", "Hello World!", true},
			{123, 1234, true},
			{123.5, 123.55, true},
			{[]byte("Hello World"), []byte("Hello World!"), true},
			{nil, new(AssertionTesterConformingObject), true},

			// cases that are expected to match
			{nil, nil, false},
			{"Hello World", "Hello World", false},
			{123, 123, false},
			{123.5, 123.5, false},
			{[]byte("Hello World"), []byte("Hello World"), false},
			{new(AssertionTesterConformingObject), new(AssertionTesterConformingObject), false},
			{&struct{}{}, &struct{}{}, false},
			// {func() int { return 23 }, func() int { return 24 }, false},
			// A case that might be confusing, especially with numeric literals
			{int(10), uint(10), true},
		}

		for _, c := range cases {
			res := NotEqual(c.expected, c.actual)
			assert.Equal(GinkgoT(), c.result, res)
		}
	})

	It("validateEqualArgs", func() {
		// assert.NotNil(GinkgoT(), validateEqualArgs(func() {}, func() {}))
		assert.Nil(GinkgoT(), validateEqualArgs(nil, nil))
	})

	It("isFunction", func() {
		assert.False(GinkgoT(), isFunction(nil))

		assert.False(GinkgoT(), isFunction(1))
		assert.True(GinkgoT(), isFunction(func() {}))
	})

	It("ObjectsAreEqual", func() {
		cases := []struct {
			expected interface{}
			actual   interface{}
			result   bool
		}{
			// cases that are expected to be equal
			{"Hello World", "Hello World", true},
			{123, 123, true},
			{123.5, 123.5, true},
			{[]byte("Hello World"), []byte("Hello World"), true},
			{nil, nil, true},

			// cases that are expected not to be equal
			{map[int]int{5: 10}, map[int]int{10: 20}, false},
			{'x', "x", false},
			{"x", 'x', false},
			{0, 0.1, false},
			{0.1, 0, false},
			{time.Now, time.Now, false},
			{func() {}, func() {}, false},
			{uint32(10), int32(10), false},
		}

		for _, c := range cases {
			res := ObjectsAreEqual(c.expected, c.actual)
			assert.Equal(GinkgoT(), c.result, res)
		}
		// Cases where type differ but values are equal
		//if !ObjectsAreEqualValues(uint32(10), int32(10)) {
		//	t.Error("ObjectsAreEqualValues should return true")
		//}
		//if ObjectsAreEqualValues(0, nil) {
		//	t.Fail()
		//}
		//if ObjectsAreEqualValues(nil, 0) {
		//	t.Fail()
		//}

	})

})
