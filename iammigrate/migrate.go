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
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	"github.com/TencentBlueKing/iam-go-sdk/client"
	"github.com/TencentBlueKing/iam-go-sdk/util"
	"github.com/mitchellh/mapstructure"
)

// DoMigate is a Go function that performs the migration.
//
// It takes three parameters: ctx, a context.Context object representing the execution context,
// client, a *client.IAMBackendClient object representing the IAM backend client, and
// data, a []byte object representing the data to be migrated.
// The function does not return anything.
func DoMigate(ctx context.Context, cli client.IAMBackendClient, data []byte, templateVar interface{},
	version int) error {
	// ping
	if err := cli.Ping(); err != nil {
		return fmt.Errorf("iam service is not available: %s", err.Error())
	}

	// format migration file, fill template variables
	data, err := FormatData(data, templateVar)
	if err != nil {
		return fmt.Errorf("format data error: %s", err.Error())
	}

	var migrations Migrations
	if err = json.Unmarshal(data, &migrations); err != nil {
		return err
	}

	if !migrations.Enabled {
		return nil
	}

	// get current model
	models, err := queryAllModels(cli, migrations.SystemID)
	if err != nil && version != 0 {
		return fmt.Errorf("query all models fail, %s", err.Error())
	}

	// do migrate
	for _, v := range migrations.Operations {
		if v.Operation == "" || v.Data == nil {
			continue
		}
		var opData []byte
		var err error
		switch d := v.Data.(type) {
		case map[string]interface{}, []map[string]interface{}, []interface{}:
			opData, err = json.Marshal(d)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("operation %s data is invalid", v.Operation)
		}
		if err = migrateFuncs[v.Operation](ctx, migrations.SystemID, cli, opData, models); err != nil {
			return fmt.Errorf("do migrate [%s] fail, %s", v.Operation, err.Error())
		}
	}
	return nil
}

// FormatData formats the given data using the provided templateVar.
//
// It takes in a byte slice of data and an interface{} templateVar.
// It returns a byte slice and an error.
func FormatData(data []byte, templateVar interface{}) ([]byte, error) {
	if templateVar == nil {
		return data, nil
	}
	tmpl, err := template.New("migrations").Parse(string(data))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := io.MultiWriter(&buf)
	err = tmpl.Execute(writer, templateVar)
	return buf.Bytes(), err
}

func queryAllModels(cli client.IAMBackendClient, systemID string) (ModelIDs, error) {
	var models ModelIDs
	var model Model
	result, err := cli.ModelQuery(systemID)
	if err != nil {
		return models, err
	}
	err = mapstructure.Decode(result, &model)
	if err != nil {
		return models, err
	}

	if model.BaseInfo.ID != "" {
		models.SystemIDs = append(models.SystemIDs, model.BaseInfo.ID)
	}
	for _, v := range model.ResourceTypes {
		models.ResourceTypeIDs = append(models.ResourceTypeIDs, v.ID)
	}
	for _, v := range model.Actions {
		models.ActionIDs = append(models.ActionIDs, v.ID)
	}
	for _, v := range model.InstanceSelections {
		models.InstanceSelectionIDs = append(models.InstanceSelectionIDs, v.ID)
	}
	return models, nil
}

// MigrateFunc is a function that performs the migration.
type MigrateFunc func(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error

var migrateFuncs = map[OperatorType]MigrateFunc{
	"upsert_system":                   UpsertSystem,
	"add_system":                      AddSystem,
	"update_system":                   UpdateSystem,
	"upsert_resource_type":            UpsertResourceType,
	"add_resource_type":               AddResourceType,
	"update_resource_type":            UpdateResourceType,
	"delete_resource_type":            DeleteResourceType,
	"upsert_instance_selection":       UpsertInstanceSelection,
	"add_instance_selection":          AddInstanceSelection,
	"update_instance_selection":       UpdateInstanceSelection,
	"delete_instance_selection":       DeleteInstanceSelection,
	"upsert_action":                   UpsertAction,
	"add_action":                      AddAction,
	"update_action":                   UpdateAction,
	"delete_action":                   DeleteAction,
	"upsert_action_groups":            UpsertActionGroups,
	"upsert_resource_creator_actions": UpsertResourceCreatorActions,
	"upsert_common_actions":           UpsertCommonActions,
	"upsert_feature_shield_rules":     UpsertFeatureShieldRules,
}

// UpsertSystem is a function that performs some operation in Go.
func UpsertSystem(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	if util.Contains(models.SystemIDs, systemID) {
		return UpdateSystem(ctx, systemID, cli, data, models)
	}
	return AddSystem(ctx, systemID, cli, data, models)
}

// AddSystem adds a system to the IAM backend
func AddSystem(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var system System
	if err := json.Unmarshal(data, &system); err != nil {
		return err
	}

	if system.ID == "" {
		return fmt.Errorf("system id is empty")
	}
	return cli.AddSystem(system)
}

// UpdateSystem updates a system in the IAM backend
func UpdateSystem(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var system System
	if err := json.Unmarshal(data, &system); err != nil {
		return err
	}

	if system.ID == "" {
		return fmt.Errorf("system id is empty")
	}
	return cli.UpdateSystem(systemID, system)
}

// UpsertResourceType is a function that performs some operation in Go.
func UpsertResourceType(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var rt ResourceType
	if err := json.Unmarshal(data, &rt); err != nil {
		return err
	}

	if rt.ID == "" {
		return fmt.Errorf("resource type id is empty")
	}
	if util.Contains(models.ResourceTypeIDs, rt.ID) {
		return cli.UpdateResourceType(systemID, rt.ID, rt)
	}
	var rts []ResourceType
	rts = append(rts, rt)
	return cli.AddResourceType(systemID, rts)
}

// AddResourceType adds a resource type to the IAM backend
func AddResourceType(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var rt ResourceType
	if err := json.Unmarshal(data, &rt); err != nil {
		return err
	}

	if rt.ID == "" {
		return fmt.Errorf("resource type id is empty")
	}
	var rts []ResourceType
	rts = append(rts, rt)
	return cli.AddResourceType(systemID, rts)
}

// UpdateResourceType updates a resource type in the IAM backend
func UpdateResourceType(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var rt ResourceType
	if err := json.Unmarshal(data, &rt); err != nil {
		return err
	}

	if rt.ID == "" {
		return fmt.Errorf("resource type id is empty")
	}
	return cli.UpdateResourceType(systemID, rt.ID, rt)
}

// DeleteResourceType deletes the resource type for a given system ID using the IAMBackendClient.
func DeleteResourceType(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var id ID
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	return cli.BatchDeleteResourceType(systemID, id.ID)
}

// UpsertInstanceSelection is a function that performs some operation in Go.
func UpsertInstanceSelection(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var is InstanceSelection
	if err := json.Unmarshal(data, &is); err != nil {
		return err
	}

	if is.ID == "" {
		return fmt.Errorf("instance selection id is empty")
	}
	if util.Contains(models.InstanceSelectionIDs, is.ID) {
		return cli.UpdateInstanceSelection(systemID, is.ID, is)
	}
	var iss []InstanceSelection
	iss = append(iss, is)
	return cli.AddInstanceSelection(systemID, iss)
}

// AddInstanceSelection adds an instance selection to the IAM backend
func AddInstanceSelection(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var is InstanceSelection
	if err := json.Unmarshal(data, &is); err != nil {
		return err
	}

	if is.ID == "" {
		return fmt.Errorf("instance selection id is empty")
	}
	var iss []InstanceSelection
	iss = append(iss, is)
	return cli.AddInstanceSelection(systemID, iss)
}

// UpdateInstanceSelection updates an instance selection in the IAM backend
func UpdateInstanceSelection(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var is InstanceSelection
	if err := json.Unmarshal(data, &is); err != nil {
		return err
	}

	if is.ID == "" {
		return fmt.Errorf("instance selection id is empty")
	}
	return cli.UpdateInstanceSelection(systemID, is.ID, is)
}

// DeleteInstanceSelection deletes the instance selection for a given system ID using the IAMBackendClient.
func DeleteInstanceSelection(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var id ID
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	return cli.BatchDeleteInstanceSelection(systemID, id.ID)
}

// UpsertAction is a function that performs some operation in Go.
func UpsertAction(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var action Action
	if err := json.Unmarshal(data, &action); err != nil {
		return err
	}

	if action.ID == "" {
		return fmt.Errorf("action id is empty")
	}
	if util.Contains(models.ActionIDs, action.ID) {
		return cli.UpdateAction(systemID, action.ID, action)
	}
	var actions []Action
	actions = append(actions, action)
	return cli.AddAction(systemID, actions)
}

// AddAction adds an action to the specified system
func AddAction(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var action Action
	if err := json.Unmarshal(data, &action); err != nil {
		return err
	}

	if action.ID == "" {
		return fmt.Errorf("action id is empty")
	}
	var actions []Action
	actions = append(actions, action)
	return cli.AddAction(systemID, actions)
}

// UpdateAction updates the action for a given system ID using the IAMBackendClient.
func UpdateAction(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var action Action
	if err := json.Unmarshal(data, &action); err != nil {
		return err
	}

	if action.ID == "" {
		return fmt.Errorf("action id is empty")
	}
	return cli.UpdateAction(systemID, action.ID, action)
}

// DeleteAction deletes the action for a given system ID using the IAMBackendClient.
func DeleteAction(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var id ID
	if err := json.Unmarshal(data, &id); err != nil {
		return err
	}
	return cli.BatchDeleteAction(systemID, id.ID)
}

// UpsertActionGroups updates the action groups for a given system ID using the IAMBackendClient.
func UpsertActionGroups(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var updateData interface{}
	if err := json.Unmarshal(data, &updateData); err != nil {
		return err
	}
	return cli.UpdateActionGroups(systemID, updateData)
}

// UpsertResourceCreatorActions updates the resource creator actions for a given system ID using the IAMBackendClient.
func UpsertResourceCreatorActions(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var updateData interface{}
	if err := json.Unmarshal(data, &updateData); err != nil {
		return err
	}
	return cli.UpdateResourceCreatorActions(systemID, updateData)
}

// UpsertCommonActions updates the common actions in the IAM backend for the given systemID.
func UpsertCommonActions(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var updateData interface{}
	if err := json.Unmarshal(data, &updateData); err != nil {
		return err
	}
	return cli.UpdateCommonActions(systemID, updateData)
}

// UpsertFeatureShieldRules updates the feature shield rules for a system.
func UpsertFeatureShieldRules(ctx context.Context, systemID string, cli client.IAMBackendClient, data []byte,
	models ModelIDs) error {
	var updateData interface{}
	if err := json.Unmarshal(data, &updateData); err != nil {
		return err
	}
	return cli.UpdateFeatureShieldRules(systemID, updateData)
}
