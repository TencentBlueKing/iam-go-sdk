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
	"encoding/json"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
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
			// if !isComparable {
			//	t.Error("object should be comparable for type " + currCase.cType)
			// }

			assert.Equal(GinkgoT(), resLess, compareLess)
			// if resLess != compareLess {
			//	t.Errorf("object less should be less than greater for type " + currCase.cType)
			// }

			resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
			assert.True(GinkgoT(), isComparable)
			// if !isComparable {
			//	t.Error("object are comparable for type " + currCase.cType)
			// }

			assert.Equal(GinkgoT(), resGreater, compareGreater)
			// if resGreater != compareGreater {
			//	t.Errorf("object greater should be greater than less for type " + currCase.cType)
			// }

			resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
			assert.True(GinkgoT(), isComparable)
			// if !isComparable {
			//	t.Error("object are comparable for type " + currCase.cType)
			// }

			assert.Equal(GinkgoT(), resEqual, compareEqual)
			// if resEqual != 0 {
			//	t.Errorf("objects should be equal for type " + currCase.cType)
			// }
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

	It("ValueEqual", func() {
		assert.True(GinkgoT(), ValueEqual("a", "a"))
		assert.True(GinkgoT(), ValueEqual(1, 1))
		assert.True(GinkgoT(), ValueEqual(float64(1), float64(1)))
		assert.True(GinkgoT(), ValueEqual(float64(1), 1))
		assert.True(GinkgoT(), ValueEqual(float64(1), json.Number("1")))
		assert.True(GinkgoT(), ValueEqual(json.Number("1"), json.Number("1")))

		assert.False(GinkgoT(), ValueEqual(1, 2))
		assert.False(GinkgoT(), ValueEqual("a", "b"))
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

	It("isNumberKind", func() {
		assert.True(GinkgoT(), isNumberKind(reflect.Uint))
		assert.True(GinkgoT(), isNumberKind(reflect.Uint8))
		assert.True(GinkgoT(), isNumberKind(reflect.Uint16))
		assert.True(GinkgoT(), isNumberKind(reflect.Uint32))
		assert.True(GinkgoT(), isNumberKind(reflect.Uint64))
		assert.True(GinkgoT(), isNumberKind(reflect.Int))
		assert.True(GinkgoT(), isNumberKind(reflect.Int8))
		assert.True(GinkgoT(), isNumberKind(reflect.Int16))
		assert.True(GinkgoT(), isNumberKind(reflect.Int32))
		assert.True(GinkgoT(), isNumberKind(reflect.Int64))
		assert.True(GinkgoT(), isNumberKind(reflect.Float32))
		assert.True(GinkgoT(), isNumberKind(reflect.Float64))

		assert.False(GinkgoT(), isNumberKind(reflect.String))
	})

	It("isFloatKind", func() {
		assert.True(GinkgoT(), isFloatKind(reflect.Float32))
		assert.True(GinkgoT(), isFloatKind(reflect.Float64))

		assert.False(GinkgoT(), isFloatKind(reflect.Uint))
		assert.False(GinkgoT(), isFloatKind(reflect.Uint8))
		assert.False(GinkgoT(), isFloatKind(reflect.Uint16))
		assert.False(GinkgoT(), isFloatKind(reflect.Uint32))
		assert.False(GinkgoT(), isFloatKind(reflect.Uint64))
		assert.False(GinkgoT(), isFloatKind(reflect.Int))
		assert.False(GinkgoT(), isFloatKind(reflect.Int8))
		assert.False(GinkgoT(), isFloatKind(reflect.Int16))
		assert.False(GinkgoT(), isFloatKind(reflect.Int32))
		assert.False(GinkgoT(), isFloatKind(reflect.Int64))
		assert.False(GinkgoT(), isFloatKind(reflect.String))
	})

	Context("toInt64", func() {
		DescribeTable("toInt64 cases", func(expected interface{}, input interface{}, willError bool) {
			v, err := toInt64(input)
			if willError {
				assert.Error(GinkgoT(), err)
				return
			}
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), expected, v)
		},
			Entry("int", int64(8), int(8), false),
			Entry("int8", int64(8), int8(8), false),
			Entry("int16", int64(8), int16(8), false),
			Entry("int32", int64(8), int32(8), false),
			Entry("int64", int64(8), int64(8), false),
			Entry("uint", int64(8), uint(8), false),
			Entry("uint8", int64(8), uint8(8), false),
			Entry("uint16", int64(8), uint16(8), false),
			Entry("uint32", int64(8), uint32(8), false),
			Entry("uint64", int64(8), uint64(8), false),

			Entry("float32", int64(0), float32(8), true),
			Entry("float64", int64(0), float64(8.31), true),
			Entry("string", float64(0), "8.31", true),
		)

	})

	Context("toFloat64", func() {
		DescribeTable("toFloat64 cases", func(expected interface{}, input interface{}, willError bool) {
			v, err := toFloat64(input)
			if willError {
				assert.Error(GinkgoT(), err)
				return
			}
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), expected, v)
		},
			Entry("int", float64(8), int(8), false),
			Entry("int8", float64(8), int8(8), false),
			Entry("int16", float64(8), int16(8), false),
			Entry("int32", float64(8), int32(8), false),
			Entry("int64", float64(8), int64(8), false),
			Entry("uint", float64(8), uint(8), false),
			Entry("uint8", float64(8), uint8(8), false),
			Entry("uint16", float64(8), uint16(8), false),
			Entry("uint32", float64(8), uint32(8), false),
			Entry("uint64", float64(8), uint64(8), false),
			Entry("float32", float64(8), float32(8), false),
			Entry("float64", float64(8.31), float64(8.31), false),
			Entry("string", float64(0), "8.31", true),
		)
	})

	Context("castJsonNumber", func() {
		It("type wrong", func() {
			_, _, err := castJsonNumber(123)
			assert.Error(GinkgoT(), err)
		})

		It("float64 ok", func() {
			v, t, err := castJsonNumber(json.Number("123.456"))
			assert.Equal(GinkgoT(), float64(123.456), v)
			assert.Equal(GinkgoT(), reflect.Float64, t)
			assert.NoError(GinkgoT(), err)
		})

		It("float64 error", func() {
			_, _, err := castJsonNumber(json.Number("123.4abc"))
			assert.Error(GinkgoT(), err)
		})

		It("int64 ok", func() {
			v, t, err := castJsonNumber(json.Number("123"))
			assert.Equal(GinkgoT(), int64(123), v)
			assert.Equal(GinkgoT(), reflect.Int64, t)
			assert.NoError(GinkgoT(), err)
		})

		It("int64 wrong", func() {
			_, _, err := castJsonNumber(json.Number("abc"))
			assert.Error(GinkgoT(), err)
		})
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
				{v1: nil, v2: float64(1)},
				{v1: "abc", v2: nil},
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

		It("compareTwoValues, no same type, int, ok", func() {
			for _, currCase := range []struct {
				v1           interface{}
				v2           interface{}
				compareTypes []CompareType
			}{
				{v1: 1, v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: int(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: int64(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: int32(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: int16(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: int8(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: uint(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: uint64(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: uint32(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: uint16(2), compareTypes: []CompareType{compareLess}},
				{v1: 1, v2: uint8(2), compareTypes: []CompareType{compareLess}},
				{v1: int(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: int64(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: int32(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: int16(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: int8(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: uint(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: uint64(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: uint32(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: uint16(1), v2: 2, compareTypes: []CompareType{compareLess}},
				{v1: uint8(1), v2: 2, compareTypes: []CompareType{compareLess}},
			} {
				compareResult := compareTwoValues(currCase.v1, currCase.v2, currCase.compareTypes)
				assert.True(GinkgoT(), compareResult)
			}

		})

		It("compareTwoValues, no same type, float, ok", func() {
			for _, currCase := range []struct {
				v1           interface{}
				v2           interface{}
				compareTypes []CompareType
			}{
				{v1: int(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: int64(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: int32(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: int16(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: int8(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: uint(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: uint64(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: uint32(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: uint16(1), v2: 2.1, compareTypes: []CompareType{compareLess}},
				{v1: uint8(1), v2: 2.1, compareTypes: []CompareType{compareLess}},

				{v1: int(1), v2: float32(2.1), compareTypes: []CompareType{compareLess}},
				{v1: int(1), v2: float64(2.1), compareTypes: []CompareType{compareLess}},
				{v1: float32(1.1), v2: float64(2.1), compareTypes: []CompareType{compareLess}},
			} {
				compareResult := compareTwoValues(currCase.v1, currCase.v2, currCase.compareTypes)
				assert.True(GinkgoT(), compareResult)
			}
		})

		It("compareTwoValues, json.Number, ok", func() {
			for _, currCase := range []struct {
				v1           interface{}
				v2           interface{}
				compareTypes []CompareType
			}{
				{v1: int(1), v2: json.Number("2"), compareTypes: []CompareType{compareLess}},
				{v1: int(1), v2: json.Number("2.1"), compareTypes: []CompareType{compareLess}},
				{v1: json.Number("1"), v2: json.Number("2.1"), compareTypes: []CompareType{compareLess}},
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
