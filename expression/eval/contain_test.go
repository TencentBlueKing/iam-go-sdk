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
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Contain", func() {

	It("includeElement", func() {
		list1 := []string{"Foo", "Bar"}
		list2 := []int{1, 2}
		simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

		ok, found := includeElement("Hello World", "World")
		assert.True(GinkgoT(), ok)
		assert.True(GinkgoT(), found)

		ok, found = includeElement(list1, "Foo")
		assert.True(GinkgoT(), ok)
		assert.True(GinkgoT(), found)

		ok, found = includeElement(list1, "Bar")
		assert.True(GinkgoT(), ok)
		assert.True(GinkgoT(), found)

		ok, found = includeElement(list2, 1)
		assert.True(GinkgoT(), ok)
		assert.True(GinkgoT(), found)

		ok, found = includeElement(list2, 2)
		assert.True(GinkgoT(), ok)
		assert.True(GinkgoT(), found)

		ok, found = includeElement(list1, "Foo!")
		assert.True(GinkgoT(), ok)
		assert.False(GinkgoT(), found)

		ok, found = includeElement(list2, 3)
		assert.True(GinkgoT(), ok)
		assert.False(GinkgoT(), found)

		ok, found = includeElement(list2, "1")
		assert.True(GinkgoT(), ok)
		assert.False(GinkgoT(), found)

		ok, found = includeElement(simpleMap, "Foo")
		assert.True(GinkgoT(), ok)
		assert.True(GinkgoT(), found)

		ok, found = includeElement(simpleMap, "Bar")
		assert.True(GinkgoT(), ok)
		assert.False(GinkgoT(), found)

		ok, found = includeElement(1433, "1")
		assert.False(GinkgoT(), ok)
		assert.False(GinkgoT(), found)
	})

	It("Contains and NotContains", func() {
		// A is a temp struct for testing
		type A struct {
			Name, Value string
		}
		list := []string{"Foo", "Bar"}

		complexList := []*A{
			{"b", "c"},
			{"d", "e"},
			{"g", "h"},
			{"j", "k"},
		}
		simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

		cases := []struct {
			expected interface{}
			actual   interface{}
			result   bool
		}{
			{"Hello World", "Hello", true},
			{"Hello World", "Salut", false},
			{list, "Bar", true},
			{list, "Salut", false},
			{complexList, &A{"g", "h"}, true},
			{complexList, &A{"g", "e"}, false},
			{simpleMap, "Foo", true},
			{simpleMap, "Bar", false},
		}

		for _, c := range cases {
			res := Contains(c.expected, c.actual)
			assert.Equal(GinkgoT(), c.result, res)
		}

		for _, c := range cases {
			res := NotContains(c.expected, c.actual)
			assert.Equal(GinkgoT(), !c.result, res)
		}
	})

	It("In and NotIn", func() {
		// A is a temp struct for testing
		type A struct {
			Name, Value string
		}
		list := []string{"Foo", "Bar"}

		complexList := []*A{
			{"b", "c"},
			{"d", "e"},
			{"g", "h"},
			{"j", "k"},
		}
		simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

		cases := []struct {
			expected interface{}
			actual   interface{}
			result   bool
		}{
			{"Hello World", "Hello", true},
			{"Hello World", "Salut", false},
			{list, "Bar", true},
			{list, "Salut", false},
			{complexList, &A{"g", "h"}, true},
			{complexList, &A{"g", "e"}, false},
			{simpleMap, "Foo", true},
			{simpleMap, "Bar", false},
		}
		// `in` is the reverse of contains, so, same data, but reverse expected and actual

		for _, c := range cases {
			res := In(c.actual, c.expected)
			assert.Equal(GinkgoT(), c.result, res)
		}

		for _, c := range cases {
			res := NotIn(c.actual, c.expected)
			assert.Equal(GinkgoT(), !c.result, res)
		}
	})

})
