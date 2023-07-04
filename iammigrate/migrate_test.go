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

package iammigrate

import (
	"bytes"
	"testing"
)

func TestFormatData(t *testing.T) {
	testCases := []struct {
		name        string
		data        []byte
		templateVar interface{}
		expected    []byte
		expectedErr bool
	}{
		{
			name: "ValidTemplate",
			data: []byte(`{{.Name}}, {{.Age}}`),
			templateVar: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			expected: []byte(`John, 30`),
		},
		{
			name: "EmptyTemplate",
			data: []byte(``),
			templateVar: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			expected: []byte(``),
		},
		{
			name:        "EmptyVar",
			data:        []byte(`{{.Name}}, {{.Age}}`),
			templateVar: nil,
			expected:    []byte(`{{.Name}}, {{.Age}}`),
		},
		{
			name: "InvalidTemplate",
			data: []byte(`{{.Name}, {{.Age}}`),
			templateVar: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			expected:    nil,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := FormatData(tc.data, tc.templateVar)
			if tc.expectedErr != (err != nil) {
				t.Errorf("%s unexpected error: got %v, want %v", tc.name, err, tc.expectedErr)
			}

			if !bytes.Equal(result, tc.expected) {
				t.Errorf("%s unexpected result: got %s, want %s", tc.name, result, tc.expected)
			}
		})
	}
}
