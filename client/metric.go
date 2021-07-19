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

package client

import (
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/TencentBlueKing/iam-go-sdk/metric"
)

// CallbackFunc is the func object of http callback
type CallbackFunc func(response gorequest.Response, v interface{}, body []byte, errs []error)

// NewMetricCallback will record the http request data into metrics
func NewMetricCallback(system string, start time.Time) CallbackFunc {
	return func(response gorequest.Response, v interface{}, body []byte, errs []error) {
		duration := time.Since(start)

		metric.ClientRequestDuration.With(prometheus.Labels{
			"method":    response.Request.Method,
			"path":      response.Request.URL.Path,
			"status":    strconv.Itoa(response.StatusCode),
			"component": system,
		}).Observe(float64(duration / time.Millisecond))
	}
}
