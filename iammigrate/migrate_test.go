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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	"github.com/stretchr/testify/assert"
)

var _ = Describe("Migrations", func() {

	Describe("FormatData", func() {

		DescribeTable("ValidTemplate", func(expected []byte, data []byte, tempVar interface{}) {
			got, err := FormatData(data, tempVar)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), expected, got)
		},
			Entry("ValidTemplate", []byte("John, 30"), []byte(`{{.Name}}, {{.Age}}`), struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			}),
		)

		DescribeTable("EmptyTemplate", func(expected []byte, data []byte, tempVar interface{}) {
			got, err := FormatData(data, tempVar)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), expected, got)
		},
			Entry("EmptyTemplate", []byte(nil), []byte(``), struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			}),
		)

		DescribeTable("EmptyVar", func(expected []byte, data []byte, tempVar interface{}) {
			got, err := FormatData(data, tempVar)
			assert.NoError(GinkgoT(), err)
			assert.Equal(GinkgoT(), expected, got)
		},
			Entry("EmptyVar", []byte("{{.Name}}, {{.Age}}"), []byte(`{{.Name}}, {{.Age}}`), nil),
		)

		DescribeTable("InvalidTemplate", func(expected []byte, data []byte, tempVar interface{}) {
			got, err := FormatData(data, tempVar)
			assert.Error(GinkgoT(), err)
			assert.Equal(GinkgoT(), expected, got)
		},
			Entry("InvalidTemplate", nil, []byte(`{{.Name}, {{.Age}}`), struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			}),
		)

	})

})
