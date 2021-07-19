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

package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	serviceName = "iam"
)

var (
	// ClientRequestDuration 依赖 api 响应时间分布
	ClientRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "client_request_duration_milliseconds",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": serviceName},
		Buckets:     []float64{20, 50, 100, 200, 500, 1000, 2000, 5000},
	},
		[]string{"method", "path", "status", "component"},
	)
)

// RegisterMetrics will register the mtrics
func RegisterMetrics() {
	// Register the summary and the histogram with Prometheus's default registry.
	prometheus.MustRegister(ClientRequestDuration)
}
