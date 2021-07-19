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

package iam

import (
	"github.com/stretchr/testify/assert"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("iam", func() {

	Context("request.CacheKey", func() {

		It("ok", func() {
			request := NewRequest("system", NewSubject("type", "id"), NewAction("id"), []ResourceNode{
				NewResourceNode("system", "type", "id", map[string]interface{}{}),
			})

			key, err := request.CacheKey()

			assert.NoError(GinkgoT(), err)
			assert.Equal(
				GinkgoT(),
				key,
				"iam:9b9893808246fb3c7fbc0192a771d67e",
			)
		})
	})

	Context("iam.buildResourceID", func() {
		var iam = NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")

		It("one node", func() {
			resources := []ResourceNode{NewResourceNode("system", "type", "id", map[string]interface{}{})}

			resourceID := iam.buildResourceID(resources)

			assert.Equal(GinkgoT(), resourceID, "id")
		})

		It("two nodes", func() {
			resources := []ResourceNode{
				NewResourceNode("system", "type", "id", map[string]interface{}{}),
				NewResourceNode("system", "type2", "id2", map[string]interface{}{}),
			}

			resourceID := iam.buildResourceID(resources)

			assert.Equal(GinkgoT(), resourceID, "type,id/type2,id2")
		})
	})
})
