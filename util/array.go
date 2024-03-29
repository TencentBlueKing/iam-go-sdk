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
	"strconv"
	"strings"
)

// Int64ArrayToString will covert inter array to string with separator `,`
func Int64ArrayToString(input []int64, sep string) string {
	b := make([]string, len(input))
	for i, v := range input {
		b[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(b, ",")
}

// Contains checks if an element exists in a given slice.
//
// The function takes a slice of elements and a target value. It returns a boolean
// value indicating whether the target value is present in the slice or not.
func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
