/*
 * TencentBlueKing is pleased to support the open source community by making 蓝鲸智云PaaS平台社区版 (BlueKing PaaS
 * Community Edition) available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package iam

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"

	"github.com/TencentBlueKing/iam-go-sdk/expression"
)

// H is a shortcut for map[string]interface{}
type H map[string]interface{}

var validate = validator.New()

// Subject is the object of permission
type Subject struct {
	Type string `json:"type" binding:"required,oneof=user"`
	ID   string `json:"id" binding:"required"`
}

// NewSubject create a subject with type and id
func NewSubject(_type, id string) Subject {
	return Subject{
		Type: _type,
		ID:   id,
	}
}

// Action is the action of permission
type Action struct {
	ID string `json:"id" binding:"required"`
}

// NewAction create a action with id
func NewAction(id string) Action {
	return Action{
		ID: id,
	}
}

// ResourceNode is the mini unit of a resource
type ResourceNode struct {
	System    string                 `json:"system" binding:"required"`
	Type      string                 `json:"type" binding:"required"`
	ID        string                 `json:"id" binding:"required"`
	Attribute map[string]interface{} `json:"attribute" binding:"required"`
}

// NewResourceNode create a resrouce node
func NewResourceNode(system, _type, id string, attrs map[string]interface{}) ResourceNode {
	return ResourceNode{
		System:    system,
		Type:      _type,
		ID:        id,
		Attribute: attrs,
	}
}

// Resources means `one resource`
type Resources []ResourceNode

// Request is the policy query request body
type Request struct {
	System    string    `json:"system" binding:"required"`
	Subject   Subject   `json:"subject" binding:"required"`
	Action    Action    `json:"action" binding:"required"`
	Resources Resources `json:"resources" binding:"omitempty"`
}

// NewRequest create a new request for policy query
func NewRequest(system string, subject Subject, action Action, resources []ResourceNode) Request {
	return Request{
		System:    system,
		Subject:   subject,
		Action:    action,
		Resources: resources,
	}
}

// Validate will check if the request is valid
func (r *Request) Validate() error {
	return validate.Struct(r)
}

// GenObjectSet create an ObjectSet from the resources of request
func (r *Request) GenObjectSet() expression.ObjectSetInterface {
	return NewObjectSet(r.Resources)
}

// CacheKey make the unique key of a request
func (r *Request) CacheKey() (string, error) {
	b, err := jsoniter.Marshal(r)
	if err != nil {
		return "", err
	}

	h := md5.New()
	_, err = h.Write(b)
	if err != nil {
		return "", err
	}

	return "iam:" + hex.EncodeToString(h.Sum(nil)), nil
}

// NewObjectSet create an ObjectSet from resources
func NewObjectSet(resources Resources) expression.ObjectSetInterface {
	objSet := expression.NewObjectSet()

	if len(resources) == 0 {
		return objSet
	}

	for _, i := range resources {
		attrs := make(map[string]interface{}, len(i.Attribute)+1)
		attrs["id"] = i.ID

		for key, value := range i.Attribute {
			attrs[key] = value
		}

		objSet.Set(i.Type, attrs)
	}

	return objSet
}

// MultiActionRequest  is the request object for Multi Actions Request
type MultiActionRequest struct {
	System    string    `json:"system" binding:"required"`
	Subject   Subject   `json:"subject" binding:"required"`
	Actions   []Action  `json:"actions" binding:"required"`
	Resources Resources `json:"resources" binding:"omitempty"`
}

// NewMultiActionRequest create a request
func NewMultiActionRequest(
	system string,
	subject Subject,
	actions []Action,
	resources []ResourceNode,
) MultiActionRequest {
	return MultiActionRequest{
		System:    system,
		Subject:   subject,
		Actions:   actions,
		Resources: resources,
	}
}

// Validate will check if the request is valid
func (mar *MultiActionRequest) Validate() error {
	return validate.Struct(mar)
}

// ActionPolicy is the response struct
type ActionPolicy struct {
	Action    Action              `json:"action"`
	Condition expression.ExprCell `json:"condition"`
}

// ApplicationResourceNode  is the resourc node struct for application
type ApplicationResourceNode struct {
	Type string `json:"type" binding:"required"`
	ID   string `json:"id" binding:"required"`
}

// ApplicationResourceInstance is the resource instance for application
type ApplicationResourceInstance []ApplicationResourceNode

// ApplicationRelatedResourceType is the related resource type for application
type ApplicationRelatedResourceType struct {
	SystemID  string                        `json:"system_id"`
	Type      string                        `json:"type"`
	Instances []ApplicationResourceInstance `json:"instances"`
}

// Validate will check if the application related resource type is valid
func (arr *ApplicationRelatedResourceType) Validate() error {
	for i, ari := range arr.Instances {
		if len(ari) == 0 {
			return fmt.Errorf(
				"the RelatedResourceType.instances[%d] invalid: "+
					"ResourceInstance.resource_nodes should contain at least 1 ApplicationResourceNode",
				i)
		}

		for j, node := range ari {
			err := validate.Struct(node)
			if err != nil {
				return fmt.Errorf("the RelatedResourceType.instances[%d].ApplicationResourceNode[%d] invalid: %w", i, j, err)
			}
		}
	}

	return nil
}

// ApplicationAction is the action for application
type ApplicationAction struct {
	ID                   string                           `json:"id"`
	RelatedResourceTypes []ApplicationRelatedResourceType `json:"related_resource_types"`
}

// NewApplicationAction will create the application action
func NewApplicationAction(id string, rrt []ApplicationRelatedResourceType) ApplicationAction {
	return ApplicationAction{
		ID:                   id,
		RelatedResourceTypes: rrt,
	}
}

// Validate will check if the application action is valid
func (aa *ApplicationAction) Validate() error {
	for i, rrt := range aa.RelatedResourceTypes {
		err := rrt.Validate()
		if err != nil {
			return fmt.Errorf("the Action.related_resource_types[%d] invalid: %w", i, err)
		}
	}

	return nil
}

// Application is the application for permission
type Application struct {
	SystemID string              `json:"system_id"`
	Actions  []ApplicationAction `json:"actions"`
}

// Validate will check if the application is valid
func (a *Application) Validate() error {
	if len(a.Actions) == 0 {
		return errors.New("the Application.actions invalid: should contain at least 1 Action")
	}

	for i, action := range a.Actions {
		err := action.Validate()
		if err != nil {
			return fmt.Errorf("the Application.actions[%d] invalid: %w", i, err)
		}
	}

	return nil
}

// NewApplication will create the application
func NewApplication(system string, actions []ApplicationAction) Application {
	return Application{
		SystemID: system,
		Actions:  actions,
	}
}

// ApplicationResourceNodeWithName is the resourc node struct for application, which with the names of each field
type ApplicationResourceNodeWithName struct {
	Type     string `json:"type" binding:"required"`
	TypeName string `json:"type_name" binding:"required"`
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// ApplicationResourceInstanceWithName is the resource instance for application, which with the names of each field
type ApplicationResourceInstanceWithName []ApplicationResourceNodeWithName

// ApplicationRelatedResourceTypeWithName is the related resource type with names
type ApplicationRelatedResourceTypeWithName struct {
	SystemID   string                                `json:"system_id" binding:"required"`
	SystemName string                                `json:"system_name" binding:"required"`
	Type       string                                `json:"type" binding:"required"`
	TypeName   string                                `json:"type_name" binding:"required"`
	Instances  []ApplicationResourceInstanceWithName `json:"instances" binding:"required"`
}

// ApplicationActionForApply is the action for apply
type ApplicationActionForApply struct {
	ID                   string                                   `json:"id" binding:"required"`
	Name                 string                                   `json:"name" binding:"required"`
	RelatedResourceTypes []ApplicationRelatedResourceTypeWithName `json:"related_resource_types"`
}

// ApplicationActionListForApply is the action list for apply
type ApplicationActionListForApply struct {
	SystemID   string                      `json:"system_id" binding:"required"`
	SystemName string                      `json:"system_name" binding:"required"`
	Actions    []ApplicationActionForApply `json:"actions" binding:"required"`
}
