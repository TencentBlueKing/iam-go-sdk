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

	Describe("Contains", func() {
		DescribeTable("test string", func(expected bool, array []string, value string) {
			assert.Equal(GinkgoT(), expected, Contains(array, value))
		},
			Entry("empty", false, []string{}, ""),
			Entry("one string", true, []string{"1"}, "1"),
			Entry("three strings", true, []string{"1", "2", "3"}, "1"),
			Entry("not found", false, []string{"1", "2", "3"}, "4"),
		)

		DescribeTable("test int", func(expected bool, array []int, value int) {
			assert.Equal(GinkgoT(), expected, Contains(array, value))
		},
			Entry("empty", false, []int{}, 0),
			Entry("one int64", true, []int{1}, 1),
			Entry("three int64", true, []int{1, 2, 3}, 1),
			Entry("not found", false, []int{1, 2, 3}, 4),
		)
	})

})
