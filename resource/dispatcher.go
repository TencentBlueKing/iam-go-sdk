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

package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Page the object for pagination
type Page struct {
	Offset int `json:"offset" binding:"required"`
	Limit  int `json:"limit" binding:"required"`
}

// Request the callback request body
type Request struct {
	Context context.Context        `json:"-"`
	Type    string                 `json:"type" binding:"required"`
	Method  string                 `json:"method" binding:"required"`
	Filter  map[string]interface{} `json:"filter" binding:"omitempty"`
	Page    Page                   `json:"page" binding:"omitempty"`
}

// Response the response body
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	// maybe a [] or {}
	Data interface{} `json:"data"`
}

// Dispatcher is the interface of dispatcher, for callback
type Dispatcher interface {
	RegisterProvider(_type string, provider Provider)
	GetProvider(_type string) (provider Provider, exist bool)
}

// NewDispatcher will create a dispatcher
func NewDispatcher() Dispatcher {
	return &dispatcher{
		providers: make(map[string]Provider, 8),
	}
}

type dispatcher struct {
	providers map[string]Provider
}

// RegisterProvider will register a provider
func (d *dispatcher) RegisterProvider(_type string, provider Provider) {
	d.providers[_type] = provider
}

// GetProvider get the provider by type
func (d *dispatcher) GetProvider(_type string) (provider Provider, exist bool) {
	provider, exist = d.providers[_type]
	return
}

// NewDispatchHandler will create a http handler for dispatcher
func NewDispatchHandler(d Dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := doDispatch(r, d)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

}

func doDispatch(r *http.Request, d Dispatcher) Response {
	// parse request.Body into req
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return Response{
			Code:    400,
			Message: "bad request, parse json fail",
		}
	}

	// get the provider via resourceType
	provider, exist := d.GetProvider(req.Type)
	if !exist {
		return Response{
			Code:    404,
			Message: fmt.Sprintf("type=%s not supported or the provider not registered", req.Type),
		}
	}

	// set context
	req.Context = r.Context()

	// dispatch the method
	switch req.Method {
	case "list_attr":
		return provider.ListAttr(req)
	case "list_attr_value":
		return provider.ListAttrValue(req)
	case "list_instance":
		return provider.ListInstance(req)
	case "fetch_instance_info":
		return provider.FetchInstanceInfo(req)
	case "list_instance_by_policy":
		return provider.ListInstanceByPolicy(req)
	case "search_instance":
		return provider.SearchInstance(req)
	default:
		return Response{
			Code:    404,
			Message: fmt.Sprintf("method=%s not supported", req.Method),
		}
	}
}
