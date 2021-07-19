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

/* NOTE: copied from https://github.com/stretchr/testify/assert/assertion_compare_test.go and modified
 *  The original versions of the files are MIT licensed
 */

package eval

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Compare", func() {

	It("compare", func() {
		type customInt int
		type customInt8 int8
		type customInt16 int16
		type customInt32 int32
		type customInt64 int64
		type customUInt uint
		type customUInt8 uint8
		type customUInt16 uint16
		type customUInt32 uint32
		type customUInt64 uint64
		type customFloat32 float32
		type customFloat64 float64
		type customString string
		for _, currCase := range []struct {
			less    interface{}
			greater interface{}
			cType   string
		}{
			{less: customString("a"), greater: customString("b"), cType: "string"},
			{less: "a", greater: "b", cType: "string"},
			{less: customInt(1), greater: customInt(2), cType: "int"},
			{less: int(1), greater: int(2), cType: "int"},
			{less: customInt8(1), greater: customInt8(2), cType: "int8"},
			{less: int8(1), greater: int8(2), cType: "int8"},
			{less: customInt16(1), greater: customInt16(2), cType: "int16"},
			{less: int16(1), greater: int16(2), cType: "int16"},
			{less: customInt32(1), greater: customInt32(2), cType: "int32"},
			{less: int32(1), greater: int32(2), cType: "int32"},
			{less: customInt64(1), greater: customInt64(2), cType: "int64"},
			{less: int64(1), greater: int64(2), cType: "int64"},
			{less: customUInt(1), greater: customUInt(2), cType: "uint"},
			{less: uint8(1), greater: uint8(2), cType: "uint8"},
			{less: customUInt8(1), greater: customUInt8(2), cType: "uint8"},
			{less: uint16(1), greater: uint16(2), cType: "uint16"},
			{less: customUInt16(1), greater: customUInt16(2), cType: "uint16"},
			{less: uint32(1), greater: uint32(2), cType: "uint32"},
			{less: customUInt32(1), greater: customUInt32(2), cType: "uint32"},
			{less: uint64(1), greater: uint64(2), cType: "uint64"},
			{less: customUInt64(1), greater: customUInt64(2), cType: "uint64"},
			{less: float32(1.23), greater: float32(2.34), cType: "float32"},
			{less: customFloat32(1.23), greater: customFloat32(2.23), cType: "float32"},
			{less: float64(1.23), greater: float64(2.34), cType: "float64"},
			{less: customFloat64(1.23), greater: customFloat64(2.34), cType: "float64"},
		} {
			resLess, isComparable := compare(currCase.less, currCase.greater, reflect.ValueOf(currCase.less).Kind())

			assert.True(GinkgoT(), isComparable)
			//if !isComparable {
			//	t.Error("object should be comparable for type " + currCase.cType)
			//}

			assert.Equal(GinkgoT(), resLess, compareLess)
			//if resLess != compareLess {
			//	t.Errorf("object less should be less than greater for type " + currCase.cType)
			//}

			resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
			assert.True(GinkgoT(), isComparable)
			//if !isComparable {
			//	t.Error("object are comparable for type " + currCase.cType)
			//}

			assert.Equal(GinkgoT(), resGreater, compareGreater)
			//if resGreater != compareGreater {
			//	t.Errorf("object greater should be greater than less for type " + currCase.cType)
			//}

			resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
			assert.True(GinkgoT(), isComparable)
			//if !isComparable {
			//	t.Error("object are comparable for type " + currCase.cType)
			//}

			assert.Equal(GinkgoT(), resEqual, compareEqual)
			//if resEqual != 0 {
			//	t.Errorf("objects should be equal for type " + currCase.cType)
			//}
		}

	})

	It("Greater", func() {
		assert.True(GinkgoT(), Greater(2, 1))
		assert.False(GinkgoT(), Greater(1, 1))
		assert.False(GinkgoT(), Greater(1, 2))

		// Check error report
		for _, currCase := range []struct {
			less    interface{}
			greater interface{}
			msg     string
		}{
			{less: "a", greater: "b", msg: `"a" is not greater than "b"`},
			{less: int(1), greater: int(2), msg: `"1" is not greater than "2"`},
			{less: int8(1), greater: int8(2), msg: `"1" is not greater than "2"`},
			{less: int16(1), greater: int16(2), msg: `"1" is not greater than "2"`},
			{less: int32(1), greater: int32(2), msg: `"1" is not greater than "2"`},
			{less: int64(1), greater: int64(2), msg: `"1" is not greater than "2"`},
			{less: uint8(1), greater: uint8(2), msg: `"1" is not greater than "2"`},
			{less: uint16(1), greater: uint16(2), msg: `"1" is not greater than "2"`},
			{less: uint32(1), greater: uint32(2), msg: `"1" is not greater than "2"`},
			{less: uint64(1), greater: uint64(2), msg: `"1" is not greater than "2"`},
			{less: float32(1.23), greater: float32(2.34), msg: `"1.23" is not greater than "2.34"`},
			{less: float64(1.23), greater: float64(2.34), msg: `"1.23" is not greater than "2.34"`},
		} {
			assert.False(GinkgoT(), Greater(currCase.less, currCase.greater))
		}

	})

	It("GreaterOrEqual", func() {
		assert.True(GinkgoT(), GreaterOrEqual(2, 1))
		assert.True(GinkgoT(), GreaterOrEqual(1, 1))
		assert.False(GinkgoT(), GreaterOrEqual(1, 2))

		// Check error report
		for _, currCase := range []struct {
			less    interface{}
			greater interface{}
			msg     string
		}{
			{less: "a", greater: "b", msg: `"a" is not greater than or equal to "b"`},
			{less: int(1), greater: int(2), msg: `"1" is not greater than or equal to "2"`},
			{less: int8(1), greater: int8(2), msg: `"1" is not greater than or equal to "2"`},
			{less: int16(1), greater: int16(2), msg: `"1" is not greater than or equal to "2"`},
			{less: int32(1), greater: int32(2), msg: `"1" is not greater than or equal to "2"`},
			{less: int64(1), greater: int64(2), msg: `"1" is not greater than or equal to "2"`},
			{less: uint8(1), greater: uint8(2), msg: `"1" is not greater than or equal to "2"`},
			{less: uint16(1), greater: uint16(2), msg: `"1" is not greater than or equal to "2"`},
			{less: uint32(1), greater: uint32(2), msg: `"1" is not greater than or equal to "2"`},
			{less: uint64(1), greater: uint64(2), msg: `"1" is not greater than or equal to "2"`},
			{less: float32(1.23), greater: float32(2.34), msg: `"1.23" is not greater than or equal to "2.34"`},
			{less: float64(1.23), greater: float64(2.34), msg: `"1.23" is not greater than or equal to "2.34"`},
		} {
			assert.False(GinkgoT(), GreaterOrEqual(currCase.less, currCase.greater))
		}
	})

	It("Less", func() {
		assert.True(GinkgoT(), Less(1, 2))
		assert.False(GinkgoT(), Less(1, 1))
		assert.False(GinkgoT(), Less(2, 1))

		// Check error report
		for _, currCase := range []struct {
			less    interface{}
			greater interface{}
			msg     string
		}{
			{less: "a", greater: "b", msg: `"b" is not less than "a"`},
			{less: int(1), greater: int(2), msg: `"2" is not less than "1"`},
			{less: int8(1), greater: int8(2), msg: `"2" is not less than "1"`},
			{less: int16(1), greater: int16(2), msg: `"2" is not less than "1"`},
			{less: int32(1), greater: int32(2), msg: `"2" is not less than "1"`},
			{less: int64(1), greater: int64(2), msg: `"2" is not less than "1"`},
			{less: uint8(1), greater: uint8(2), msg: `"2" is not less than "1"`},
			{less: uint16(1), greater: uint16(2), msg: `"2" is not less than "1"`},
			{less: uint32(1), greater: uint32(2), msg: `"2" is not less than "1"`},
			{less: uint64(1), greater: uint64(2), msg: `"2" is not less than "1"`},
			{less: float32(1.23), greater: float32(2.34), msg: `"2.34" is not less than "1.23"`},
			{less: float64(1.23), greater: float64(2.34), msg: `"2.34" is not less than "1.23"`},
		} {
			assert.False(GinkgoT(), Less(currCase.greater, currCase.less))
		}

	})

	It("LessOrEqual", func() {
		assert.True(GinkgoT(), LessOrEqual(1, 2))

		assert.True(GinkgoT(), LessOrEqual(1, 1))
		assert.False(GinkgoT(), LessOrEqual(2, 1))

		// Check error report
		for _, currCase := range []struct {
			less    interface{}
			greater interface{}
			msg     string
		}{
			{less: "a", greater: "b", msg: `"b" is not less than or equal to "a"`},
			{less: int(1), greater: int(2), msg: `"2" is not less than or equal to "1"`},
			{less: int8(1), greater: int8(2), msg: `"2" is not less than or equal to "1"`},
			{less: int16(1), greater: int16(2), msg: `"2" is not less than or equal to "1"`},
			{less: int32(1), greater: int32(2), msg: `"2" is not less than or equal to "1"`},
			{less: int64(1), greater: int64(2), msg: `"2" is not less than or equal to "1"`},
			{less: uint8(1), greater: uint8(2), msg: `"2" is not less than or equal to "1"`},
			{less: uint16(1), greater: uint16(2), msg: `"2" is not less than or equal to "1"`},
			{less: uint32(1), greater: uint32(2), msg: `"2" is not less than or equal to "1"`},
			{less: uint64(1), greater: uint64(2), msg: `"2" is not less than or equal to "1"`},
			{less: float32(1.23), greater: float32(2.34), msg: `"2.34" is not less than or equal to "1.23"`},
			{less: float64(1.23), greater: float64(2.34), msg: `"2.34" is not less than or equal to "1.23"`},
		} {
			assert.False(GinkgoT(), LessOrEqual(currCase.greater, currCase.less))
		}
	})

	Context("compareTwoValues", func() {

		It("compareTwoValuesDifferentValuesTypes", func() {
			for _, currCase := range []struct {
				v1            interface{}
				v2            interface{}
				compareResult bool
			}{
				{v1: 123, v2: "abc"},
				{v1: "abc", v2: 123456},
				{v1: float64(12), v2: "123"},
				{v1: "float(12)", v2: float64(1)},
			} {
				compareResult := compareTwoValues(currCase.v1, currCase.v2, []CompareType{compareLess, compareEqual, compareGreater})
				assert.False(GinkgoT(), compareResult)
			}
		})
		It("compareTwoValuesNotComparableValues", func() {
			// CompareStruct is a temp struct for testing
			type CompareStruct struct {
			}

			for _, currCase := range []struct {
				v1 interface{}
				v2 interface{}
			}{
				{v1: CompareStruct{}, v2: CompareStruct{}},
				{v1: map[string]int{}, v2: map[string]int{}},
				{v1: make([]int, 5), v2: make([]int, 5)},
			} {
				compareResult := compareTwoValues(currCase.v1, currCase.v2, []CompareType{compareLess, compareEqual, compareGreater})
				assert.False(GinkgoT(), compareResult)
			}

		})
		It("compareTwoValuesCorrectCompareResult", func() {
			for _, currCase := range []struct {
				v1           interface{}
				v2           interface{}
				compareTypes []CompareType
			}{
				{v1: 1, v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: 2, compareTypes: []CompareType{compareLess, compareEqual}},
				{v1: 2, v2: 2, compareTypes: []CompareType{compareGreater, compareEqual}},
				{v1: 2, v2: 2, compareTypes: []CompareType{compareEqual}},
				{v1: 2, v2: 1, compareTypes: []CompareType{compareEqual, compareGreater}},
				{v1: 2, v2: 1, compareTypes: []CompareType{compareGreater}},
			} {
				compareResult := compareTwoValues(currCase.v1, currCase.v2, currCase.compareTypes)
				assert.True(GinkgoT(), compareResult)
			}

		})

	})

	It("containsValue", func() {

		for _, currCase := range []struct {
			values []CompareType
			value  CompareType
			result bool
		}{
			{values: []CompareType{compareGreater}, value: compareGreater, result: true},
			{values: []CompareType{compareGreater, compareLess}, value: compareGreater, result: true},
			{values: []CompareType{compareGreater, compareLess}, value: compareLess, result: true},
			{values: []CompareType{compareGreater, compareLess}, value: compareEqual, result: false},
		} {
			compareResult := containsValue(currCase.values, currCase.value)
			assert.Equal(GinkgoT(), currCase.result, compareResult)
		}

	})

})
