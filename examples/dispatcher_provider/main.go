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
// +build !codeanalysis

package main

import (
	"net/http"

	"github.com/TencentBlueKing/iam-go-sdk/resource"
)

// DummyProvider is an example of provider
type DummyProvider struct {
}

// ListAttr implements the list_attr
func (d DummyProvider) ListAttr(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// ListAttrValue implements the list_attr_value
func (d DummyProvider) ListAttrValue(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// ListInstance implements the list_instance
func (d DummyProvider) ListInstance(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// FetchInstanceInf implements the fetch_instance_info
func (d DummyProvider) FetchInstanceInfo(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// ListInstanceByPolicy implements the list_instance_by_policy
func (d DummyProvider) ListInstanceByPolicy(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// SearchInstance implements the search_instance
func (d DummyProvider) SearchInstance(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

func main() {
	d := resource.NewDispatcher()
	dummyProvider := DummyProvider{}

	// type=dummy will use the dummyProvider
	d.RegisterProvider("dummy", dummyProvider)

	handler := resource.NewDispatchHandler(d)

	// we just register only one api to handle all resource types - callback api
	http.HandleFunc("/api/v1/resource", handler)

}
