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

package iam

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TencentBlueKing/gopkg/stringx"
	"github.com/golang-migrate/migrate/v4/source"
	jsoniter "github.com/json-iterator/go"
	"github.com/mitchellh/mapstructure"

	"github.com/TencentBlueKing/iam-go-sdk/cache"
	"github.com/TencentBlueKing/iam-go-sdk/client"
	"github.com/TencentBlueKing/iam-go-sdk/expression"
	"github.com/TencentBlueKing/iam-go-sdk/iammigrate"
	"github.com/TencentBlueKing/iam-go-sdk/logger"
	"github.com/golang-migrate/migrate/v4"

	// register file source
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// IAM is the instance of iam sdk
type IAM struct {
	appCode    string
	appSecret  string
	bkTenantID string

	client client.IAMBackendClient
}

type Option func(*IAM)

func WithBkTenantID(bkTenantID string) Option {
	return func(i *IAM) {
		i.bkTenantID = bkTenantID
	}
}

// NewIAM will create an IAM instance
func NewIAM(system, appCode, appSecret, bkAPIGatewayURL string, opts ...Option) *IAM {
	return NewAPIGatewayIAM(system, appCode, appSecret, bkAPIGatewayURL, opts...)
}

// NewAPIGatewayIAM will create an IAM instance, call all api through APIGateway
// if your TencentBlueking has a APIGateway, use this, recommend
func NewAPIGatewayIAM(system, appCode, appSecret, bkAPIGatewayURL string, opts ...Option) *IAM {
	c := &IAM{
		appCode:   appCode,
		appSecret: appSecret,
	}

	for _, opt := range opts {
		opt(c)
	}

	var clientOpts []client.Option
	if c.bkTenantID != "" {
		clientOpts = append(clientOpts, client.WithBkTenantID(c.bkTenantID))
	}
	c.client = client.NewIAMBackendClient(bkAPIGatewayURL, system, appCode, appSecret, clientOpts...)

	return c
}

// IsAllowed will check if the permission is allowed
func (i *IAM) IsAllowed(request Request) (allowed bool, err error) {
	logger.Debug("calling IAM.is_allowed(request)......")

	// 1. validate
	err = request.Validate()
	if err != nil {
		logger.Debugf("the request is invalid! err=%w", err)
		return
	}

	// 2. policy query
	logger.Debugf("the request: %v", request)
	data, err := i.client.V2PolicyQuery(request.System, request)
	if err != nil {
		logger.Errorf("do policy query fail! err=%w", err)
		return
	}
	logger.Debugf("the return policies: %#v", data)

	expr := expression.ExprCell{}
	err = mapstructure.Decode(data, &expr)
	if err != nil {
		logger.Errorf("decode policy query data to expr fail! err=%w", err)
		return
	}
	logger.Debugf("the expr: %#v", expr)

	// 3. make objSet
	objSet := request.GenObjectSet()

	// 4. eval
	evalBegin := time.Now()
	allowed = expr.Eval(objSet)
	logger.Debugf("the return expr: %s", expr.String())
	logger.Debugf("the return expr render: %s", expr.Render(objSet))
	logger.Debugf("the return expr eval: %v", allowed)
	logger.Debugf("the return expr eval took %s ms", time.Since(evalBegin)/time.Millisecond)

	return allowed, nil
}

// IsAllowedWithCache will check if the permission is allowed, will cache with ttl
func (i *IAM) IsAllowedWithCache(request Request, ttl time.Duration) (allowed bool, err error) {
	var k string
	k, err = request.CacheKey()
	if err != nil {
		return
	}

	value, found := cache.Get(k)
	if found {
		return value.(bool), nil
	}

	allowed, err = i.IsAllowed(request)
	if err != nil {
		return
	}

	cache.Set(k, allowed, ttl)
	return
}

// BatchIsAllowed will batch check the permission for resources lists
func (i *IAM) BatchIsAllowed(request Request, resourcesList []Resources) (result map[string]bool, err error) {
	// logger.debug("calling IAM.is_allowed(request)......")

	// 1. validate
	err = request.Validate()
	if err != nil {
		return
	}

	// 2. policy query without resources
	if len(request.Resources) != 0 {
		request.Resources = Resources{}
	}

	data, err := i.client.V2PolicyQuery(request.System, request)
	if err != nil {
		return
	}

	expr := expression.ExprCell{}
	err = mapstructure.Decode(data, &expr)
	if err != nil {
		return
	}

	result = make(map[string]bool, len(resourcesList))
	for _, resources := range resourcesList {
		// 3. make objSet
		objSet := NewObjectSet(resources)

		// 4. eval
		allowed := expr.Eval(objSet)
		result[i.buildResourceID(resources)] = allowed
	}

	return result, nil
}

func (i *IAM) buildResourceID(resources Resources) string {
	if len(resources) == 1 {
		return resources[0].ID
	}

	nodeIDs := make([]string, 0, len(resources))
	for _, node := range resources {
		nodeIDs = append(nodeIDs, fmt.Sprintf("%s,%s", node.Type, node.ID))
	}

	return strings.Join(nodeIDs, "/")
}

// ResourceMultiActionsAllowed will check the permission of one-resource with multi-actions
func (i *IAM) ResourceMultiActionsAllowed(request MultiActionRequest) (result map[string]bool, err error) {
	// 1. validate
	err = request.Validate()
	if err != nil {
		return
	}

	// 2. batch action policy query
	logger.Debugf("the request: %v", request)
	data, err := i.client.V2PolicyQueryByActions(request.System, request)
	if err != nil {
		logger.Errorf("do policy query by actions fail! err=%w", err)
		return
	}
	logger.Debugf("the return policies of actions: %#v", data)

	result = make(map[string]bool, len(request.Actions))

	// 3. make objSet
	objSet := NewObjectSet(request.Resources)

	// 4. calculate perms
	var actionPolicies []ActionPolicy
	err = mapstructure.Decode(data, &actionPolicies)
	if err != nil {
		logger.Errorf("decode policy query by actions data to expr fail! err=%w", err)
		return
	}
	for _, actionPolicy := range actionPolicies {
		allowed := actionPolicy.Condition.Eval(objSet)
		result[actionPolicy.Action.ID] = allowed
	}
	return
}

// BatchResourceMultiActionsAllowed will check the permissions of batch-resource with multi-actions
func (i *IAM) BatchResourceMultiActionsAllowed(
	request MultiActionRequest,
	resourcesList []Resources,
) (results map[string]map[string]bool, err error) {
	// 1. validate
	err = request.Validate()
	if err != nil {
		return
	}

	// 2. policy query without resources
	if len(request.Resources) != 0 {
		request.Resources = Resources{}
	}

	// 3. batch action policy query
	logger.Debugf("the request: %v", request)
	data, err := i.client.V2PolicyQueryByActions(request.System, request)
	if err != nil {
		logger.Errorf("do policy query by actions fail! err=%w", err)
		return
	}
	logger.Debugf("the return policies of actions: %#v", data)

	results = make(map[string]map[string]bool, len(resourcesList))

	for _, resources := range resourcesList {
		result := make(map[string]bool, len(request.Actions))

		// 4. make objSet
		objSet := NewObjectSet(resources)

		// 5. calculate perms
		var actionPolicies []ActionPolicy
		err = mapstructure.Decode(data, &actionPolicies)
		if err != nil {
			logger.Errorf("decode policy query by actions data to expr fail! err=%w", err)
			return
		}
		for _, actionPolicy := range actionPolicies {
			allowed := actionPolicy.Condition.Eval(objSet)
			result[actionPolicy.Action.ID] = allowed
		}
		results[i.buildResourceID(resources)] = result
	}
	return
}

// GetToken will get the token of system
func (i *IAM) GetToken() (token string, err error) {
	return i.client.GetToken()
}

// IsBasicAuthAllowed will check basic auth of callback request
func (i *IAM) IsBasicAuthAllowed(username, password string) (err error) {
	if username != "bk_iam" {
		err = errors.New("username is not bk_iam")
		return
	}

	token, err := i.client.GetToken()
	if err != nil {
		err = fmt.Errorf("get system token fail: %w", err)
		return
	}

	if password != token {
		err = fmt.Errorf("password in basic_auth not equals to system token [password=%s***, token=%s***]",
			stringx.Truncate(password, 6), stringx.Truncate(token, 6))
		return
	}

	return nil
}

// GetApplyURL will generate the application URL
func (i *IAM) GetApplyURL(application Application) (url string, err error) {
	err = application.Validate()
	if err != nil {
		return
	}

	url, err = i.client.GetApplyURL(application)

	return
}

// GenPermissionApplyData will generate the apply data
func (i *IAM) GenPermissionApplyData(a ApplicationActionListForApply) (data H, err error) {
	j, err := jsoniter.Marshal(a)
	if err != nil {
		return
	}

	err = jsoniter.Unmarshal(j, &data)
	if err != nil {
		return
	}
	return
}

// Migrate migrates the iam model using the provided migrations directory.
//
// It takes the following parameters:
//   - db: a pointer to a sql.DB instance representing the database connection.
//   - migrateTable: the name of the table where migration information is stored.
//   - sourceDriver: the migrations source driver.
//   - timeout: the duration after which a migration statement times out.
//
// It returns an error if any error occurs during the migration process.
func (i *IAM) Migrate(db *sql.DB, sourceDriver source.Driver, migrateTable string, timeout time.Duration,
	templateVar interface{}) error {
	databaseInstance, err := iammigrate.WithInstance(db, &iammigrate.Config{
		MigrationsTable:  migrateTable,
		StatementTimeout: timeout,
		TemplateVar:      templateVar,
	}, i.client)
	if err != nil {
		return err
	}

	mig, err := migrate.NewWithInstance("migrations_source", sourceDriver, "bkiam_migrations", databaseInstance)
	if err != nil {
		return err
	}

	// remove dirty
	version, dirty, err := mig.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if dirty {
		// rollback to previous version
		if version >= 0 {
			version--
		}
		if err = mig.Force(int(version)); err != nil {
			return err
		}
	}

	// run migrations
	return mig.Up()
}

// TODO:
// - grant_resource_creator_actions
// - grant_batch_resource_creator_actions
// - grant_or_revoke_instance_permission
// - grant_or_revoke_path_permission
// - batch_grant_or_revoke_instance_permission
// - batch_grant_or_revoke_path_permission
// - query_polices_with_action_id
