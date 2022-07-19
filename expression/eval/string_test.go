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

	It("StringContains", func() {
		assert.False(GinkgoT(), StringContains("", 1))
		assert.True(GinkgoT(), StringContains("hello", "he"))
		assert.True(GinkgoT(), StringContains("hello", "el"))
		assert.True(GinkgoT(), StringContains("hello", "lo"))
		assert.False(GinkgoT(), StringContains("hello", "abc"))
	})
})
