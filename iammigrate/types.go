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

// OperatorType migration operator type
type OperatorType string

// Model all permissions model
type Model struct {
	BaseInfo           System `json:"base_info" mapstructure:"base_info"`
	ResourceTypes      []ID   `json:"resource_types" mapstructure:"resource_types"`
	Actions            []ID   `json:"actions" mapstructure:"actions"`
	InstanceSelections []ID   `json:"instance_selections" mapstructure:"instance_selections"`
}

// ModelIDs model ids
type ModelIDs struct {
	SystemIDs            []string `json:"system_ids"`
	ResourceTypeIDs      []string `json:"resource_type_ids"`
	ActionIDs            []string `json:"action_ids"`
	InstanceSelectionIDs []string `json:"instance_selection_ids"`
}

// ID model id
type ID struct {
	ID string `json:"id"`
}

// Migrations migration
type Migrations struct {
	SystemID   string      `json:"system_id"`
	Enabled    bool        `json:"enabled"`
	Operations []Operation `json:"operations"`
}

// Operation operation
type Operation struct {
	Operation OperatorType `json:"operation"`
	Data      interface{}  `json:"data"`
}

// System system
type System struct {
	ID             string         `json:"id"`
	Name           string         `json:"name"`
	NameEN         string         `json:"name_en"`
	Description    string         `json:"description"`
	DescriptionEN  string         `json:"description_en"`
	Clients        string         `json:"clients"`
	ProviderConfig ProviderConfig `json:"provider_config"`
}

// ProviderConfig provider config
type ProviderConfig struct {
	Host    string `json:"host,omitempty"`
	Path    string `json:"path,omitempty"`
	Auth    string `json:"auth,omitempty"` // basic
	Healthz string `json:"healthz,omitempty"`
}

// Action action
type Action struct {
	ID                   string                `json:"id"`
	Name                 string                `json:"name"`
	NameEN               string                `json:"name_en"`
	Description          string                `json:"description"`
	DescriptionEN        string                `json:"description_en"`
	Type                 string                `json:"type"`
	Version              int                   `json:"version"`
	RelatedResourceTypes []RelatedResourceType `json:"related_resource_types"`
	RelatedActions       []string              `json:"related_actions"`
}

// RelatedResourceType related resource type
type RelatedResourceType struct {
	SystemID                  string                     `json:"system_id"`
	ID                        string                     `json:"id"`
	NameAlias                 string                     `json:"name_alias"`
	NameAliasEN               string                     `json:"name_alias_en"`
	SelectionMode             string                     `json:"selection_mode"`
	RelatedInstanceSelections []RelatedInstanceSelection `json:"related_instance_selections"`
}

// RelatedInstanceSelection related instance selection
type RelatedInstanceSelection struct {
	SystemID      string `json:"system_id"`
	ID            string `json:"id"`
	IgnoreIAMPath bool   `json:"ignore_iam_path"`
}

// ResourceType resource type
type ResourceType struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	NameEN         string              `json:"name_en"`
	Description    string              `json:"description"`
	DescriptionEN  string              `json:"description_en"`
	Sensitivity    int                 `json:"sensitivity"`
	Parents        []ResourceTypeChain `json:"parents"`
	ProviderConfig ProviderConfig      `json:"provider_config"`
	Version        int                 `json:"version"`
}

// ResourceTypeChain resource type chain
type ResourceTypeChain struct {
	ID       string `json:"id"`
	SystemID string `json:"system_id"`
}

// InstanceSelection instance selection
type InstanceSelection struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	NameEN            string              `json:"name_en"`
	IsDynamic         bool                `json:"is_dynamic"`
	ResourceTypeChain []ResourceTypeChain `json:"resource_type_chain"`
}
