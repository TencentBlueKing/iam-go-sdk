//go:build !codeanalysis
// +build !codeanalysis

/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 权限中心 Go SDK(iam-go-sdk) available.
 * Copyright (C) 2017-2021 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/TencentBlueKing/iam-go-sdk"
	"github.com/TencentBlueKing/iam-go-sdk/logger"
)

func main() {
	// create a logger
	log := &logrus.Logger{
		Out:          os.Stderr,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.DebugLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	// do set logger
	logger.SetLogger(log)

	req := iam.NewRequest(
		"bk_paas",
		iam.NewSubject("user", "admin"),
		iam.NewAction("access_developer_center"),
		[]iam.ResourceNode{},
	)

	i := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{bk_iam_apigateway_url}")
	// if your TencentBlueking has a APIGateway, use NewAPIGatewayIAM, the url suffix is /stage/(for testing) and /prod/(for production)
	// i := iam.NewAPIGatewayIAM("bk_paas", "bk_paas", "{app_secret}", "http://bk-iam.{APIGATEWAY_DOMAIN}/stage/")

	// support multiple tenants, you can use NewIAM with bk tenant id option
	// i := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{bk_iam_apigateway_url}", iam.WithBkTenantID("my_tenant_id"))

	allowed, err := i.IsAllowed(req)
	fmt.Println("isAllowed:", allowed, err)

	// check 3 times but only call iam backend once
	allowed, err = i.IsAllowedWithCache(req, 10*time.Second)
	allowed, err = i.IsAllowedWithCache(req, 10*time.Second)
	i2 := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{bk_iam_apigateway_url}")
	allowed, err = i2.IsAllowedWithCache(req, 10*time.Second)
	fmt.Println("isAllowedWithCache:", allowed, err)

	multiReq := iam.NewMultiActionRequest(
		"bk_sops",
		iam.NewSubject("user", "admin"),
		[]iam.Action{
			iam.NewAction("task_delete"),
			iam.NewAction("task_edit"),
			iam.NewAction("task_view"),
		},
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "1", map[string]interface{}{"iam_resource_owner": "admin"}),
		},
	)
	i3 := iam.NewIAM("bk_sops", "bk_sops", "{app_secret}", "http://{bk_iam_apigateway_url}")
	result, err := i3.ResourceMultiActionsAllowed(multiReq)
	fmt.Println("ResourceMultiActionsAllowed: ", result, err)

	multiReq.Resources = iam.Resources{}
	resourcesList := []iam.Resources{
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "1", map[string]interface{}{"iam_resource_owner": "admin"}),
		},
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "2", map[string]interface{}{"iam_resource_owner": "admin2"}),
		},
		[]iam.ResourceNode{
			iam.NewResourceNode("bk_sops", "task", "3", map[string]interface{}{"iam_resource_owner": "admin3"}),
		},
	}
	results, err := i3.BatchResourceMultiActionsAllowed(multiReq, resourcesList)
	fmt.Println("BatchResourceMultiActionsAllowed: ", results, err)

	actions := []iam.ApplicationAction{}
	application := iam.NewApplication("bk_paas", actions)

	url, err := i.GetApplyURL(application)
	fmt.Println("GetApplyURL:", url, err)

	err = i.IsBasicAuthAllowed("bk_iam", "3b223guhmzwlq7oto417b4g41rqnboip")
	fmt.Println("IsBasicAuthAllowed:", err)

	token, err := i.GetToken()
	fmt.Println("GetToken:", token, err)

}
