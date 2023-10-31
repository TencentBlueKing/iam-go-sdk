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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/TencentBlueKing/gopkg/conv"
	"github.com/parnurzeal/gorequest"

	"github.com/TencentBlueKing/iam-go-sdk/logger"
	"github.com/TencentBlueKing/iam-go-sdk/util"
)

var _ IAMBackendClient = &iamBackendClient{}

const (
	bkIAMVersion = "1"
)

// Method is the type of http method
type Method string

var (
	// POST http post
	POST Method = "POST"
	// GET http get
	GET Method = "GET"
	// PUT http put
	PUT Method = "PUT"
	// DELETE http delete
	DELETE Method = "DELETE"
)

// IAMBackendBaseResponse is the struct of iam backend response
type IAMBackendBaseResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Error will check if the response with error
func (r *IAMBackendBaseResponse) Error() error {
	if r.Code == 0 {
		return nil
	}

	return fmt.Errorf("response error[code=`%d`,  message=`%s`]", r.Code, r.Message)
}

// String will return the detail text of the response
func (r *IAMBackendBaseResponse) String() string {
	return fmt.Sprintf("response[code=`%d`, message=`%s`, data=`%s`]", r.Code, r.Message, conv.BytesToString(r.Data))
}

// IAMBackendClient is the interface of iam backend client
type IAMBackendClient interface {
	Ping() error
	GetToken() (token string, err error)

	PolicyQuery(body interface{}) (map[string]interface{}, error)
	PolicyQueryByActions(body interface{}) ([]map[string]interface{}, error)

	V2PolicyQuery(system string, body interface{}) (data map[string]interface{}, err error)
	V2PolicyQueryByActions(system string, body interface{}) (data []map[string]interface{}, err error)
	V2PolicyAuth(system string, body interface{}) (data map[string]interface{}, err error)

	PolicyAuth(body interface{}) (data map[string]interface{}, err error)
	PolicyAuthByResources(body interface{}) (data map[string]interface{}, err error)
	PolicyAuthByActions(body interface{}) (data map[string]interface{}, err error)

	PolicyGet(policyID int64) (data map[string]interface{}, err error)
	PolicyList(body interface{}) (data map[string]interface{}, err error)
	PolicySubjects(policyIDs []int64) (data []map[string]interface{}, err error)

	GetApplyURL(body interface{}) (string, error)

	// Model
	ModelQuery(system string) (map[string]interface{}, error)
	AddSystem(body interface{}) error
	UpdateSystem(system string, body interface{}) error
	AddResourceType(system string, body interface{}) error
	UpdateResourceType(system, resourceTypeID string, body interface{}) error
	BatchDeleteResourceType(system string, resourceTypeIDs ...string) error
	AddInstanceSelection(system string, body interface{}) error
	UpdateInstanceSelection(system, instanceSelectionID string, body interface{}) error
	BatchDeleteInstanceSelection(system string, instanceSelectionIDs ...string) error
	AddAction(system string, body interface{}) error
	UpdateAction(system, actionID string, body interface{}) error
	BatchDeleteAction(system string, actionIDs ...string) error
	AddActionGroups(system string, body interface{}) error
	UpdateActionGroups(system string, body interface{}) error
	AddResourceCreatorActions(system string, body interface{}) error
	UpdateResourceCreatorActions(system string, body interface{}) error
	AddCommonActions(system string, body interface{}) error
	UpdateCommonActions(system string, body interface{}) error
	AddFeatureShieldRules(system string, body interface{}) error
	UpdateFeatureShieldRules(system string, body interface{}) error
}

type iamBackendClient struct {
	Host         string
	IsAPIGateway bool

	System    string
	appCode   string
	appSecret string

	isApiDebugEnabled bool
	isApiForceEnabled bool
}

// NewIAMBackendClient will create a iam backend client
func NewIAMBackendClient(host string, isAPIGateway bool, system string, appCode string, appSecret string) IAMBackendClient {
	host = strings.TrimRight(host, "/")
	return &iamBackendClient{
		Host:         host,
		IsAPIGateway: isAPIGateway,

		System:    system,
		appCode:   appCode,
		appSecret: appSecret,

		// will add ?debug=true in url, for debug api/policy, show the details
		isApiDebugEnabled: os.Getenv("IAM_API_DEBUG") == "true" || os.Getenv("BKAPP_IAM_API_DEBUG") == "true",
		// will add ?force=true in url, for api/policy run without cache(all data from database)
		isApiForceEnabled: os.Getenv("IAM_API_FORCE") == "true" || os.Getenv("BKAPP_IAM_API_FORCE") == "true",
	}
}

func (c *iamBackendClient) call(
	method Method, path string,
	data interface{},
	timeout int64,
	responseData interface{},
) error {
	callTimeout := time.Duration(timeout) * time.Second
	if timeout == 0 {
		callTimeout = defaultTimeout
	}

	headers := map[string]string{
		"X-Bk-IAM-Version": bkIAMVersion,
	}

	if c.IsAPIGateway {
		auth, err := json.Marshal(map[string]string{
			"bk_app_code":   c.appCode,
			"bk_app_secret": c.appSecret,
		})
		if err != nil {
			return fmt.Errorf("generate apigateway call header fail. err=`%s`", err)
		}

		headers["X-Bkapi-Authorization"] = conv.BytesToString(auth)
	} else {
		headers["X-BK-APP-CODE"] = c.appCode
		headers["X-BK-APP-SECRET"] = c.appSecret
	}

	url := fmt.Sprintf("%s%s", c.Host, path)
	start := time.Now()
	callbackFunc := NewMetricCallback("IAMBackend", start)

	logger.Debugf("do http request: method=`%s`, url=`%s`, data=`%s`", method, url, data)

	// request := gorequest.New().Timeout(callTimeout).Post(url).Type("json")
	request := gorequest.New().Timeout(callTimeout).Type("json")
	switch method {
	case POST:
		request = request.Post(url).Send(data)
	case GET:
		request = request.Get(url).Query(data)
	case PUT:
		request = request.Put(url).Send(data)
	case DELETE:
		request = request.Delete(url).Send(data)
	}

	if c.isApiDebugEnabled {
		request.QueryData.Add("debug", "true")
	}
	if c.isApiForceEnabled {
		request.QueryData.Add("force", "true")
	}

	// set headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// do request
	baseResult := IAMBackendBaseResponse{}
	resp, respBody, errs := request.
		EndStruct(&baseResult, callbackFunc)

	if len(errs) != 0 {
		logFailHTTPRequest(request, resp, respBody, errs, &baseResult)
		return fmt.Errorf("gorequest errors=`%s`", errs)
	}

	body := ""
	duration := time.Since(start)
	if respBody != nil {
		body = conv.BytesToString(respBody)
	}

	logger.Debugf("http request result: %+v", baseResult.String())
	logger.Debugf("http request took %v ms", float64(duration/time.Millisecond))
	logger.Debugf("http response: status_code=%s, body=%+v", resp.StatusCode, body)

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("gorequest statusCode is %d not 200", resp.StatusCode)
		if baseResult.Message != "" {
			err = fmt.Errorf("%w. response body.code: %d, message:%s", err, baseResult.Code, baseResult.Message)
		}

		return err
	}

	if baseResult.Code != 0 {
		return fmt.Errorf("response body.code: %d, message:%s", baseResult.Code, baseResult.Message)
	}

	err := json.Unmarshal(baseResult.Data, responseData)
	if err != nil {
		return fmt.Errorf("http request response body data not valid: %w, data=`%v`", err, baseResult.Data)
	}
	return nil
}

func (c *iamBackendClient) callWithReturnMapData(
	method Method, path string,
	data interface{},
	timeout int64,
) (map[string]interface{}, error) {
	var responseData map[string]interface{}
	err := c.call(method, path, data, timeout, &responseData)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return responseData, nil
}

func (c *iamBackendClient) callWithReturnSliceMapData(
	method Method, path string,
	data interface{},
	timeout int64,
) ([]map[string]interface{}, error) {
	var responseData []map[string]interface{}
	err := c.call(method, path, data, timeout, &responseData)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	return responseData, nil
}

// Ping will check the iam backend service is ping-able
func (c *iamBackendClient) Ping() (err error) {
	url := fmt.Sprintf("%s%s", c.Host, "/ping")

	resp, _, errs := gorequest.New().Timeout(defaultTimeout).Get(url).EndBytes()
	if len(errs) != 0 {
		return fmt.Errorf("ping fail! errs=%v", errs)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ping fail! status_code=%d", resp.StatusCode)
	}
	return nil
}

// GetToken will get the token of system, use for callback requests basic auth
func (c *iamBackendClient) GetToken() (token string, err error) {
	path := fmt.Sprintf("/api/v1/model/systems/%s/token", c.System)
	data, err := c.callWithReturnMapData(GET, path, map[string]interface{}{}, 10)
	if err != nil {
		return "", err
	}
	tokenI, ok := data["token"]
	if !ok {
		return "", errors.New("no token in response body")
	}
	token, ok = tokenI.(string)
	if !ok {
		return "", errors.New("token is not a valid string")
	}
	return token, nil
}

// PolicyQuery will do policy query
func (c *iamBackendClient) PolicyQuery(body interface{}) (data map[string]interface{}, err error) {
	path := "/api/v1/policy/query"
	data, err = c.callWithReturnMapData(POST, path, body, 10)
	return
}

// V2PolicyQuery will do policy query
func (c *iamBackendClient) V2PolicyQuery(system string, body interface{}) (data map[string]interface{}, err error) {
	path := "/api/v2/policy/systems/" + system + "/query/"
	data, err = c.callWithReturnMapData(POST, path, body, 10)
	return
}

// PolicyQueryByActions will do policy query by actions
func (c *iamBackendClient) PolicyQueryByActions(body interface{}) (data []map[string]interface{}, err error) {
	path := "/api/v1/policy/query_by_actions"
	data, err = c.callWithReturnSliceMapData(POST, path, body, 10)
	return
}

// V2PolicyQueryByActions will do policy query by actions
func (c *iamBackendClient) V2PolicyQueryByActions(system string, body interface{}) (data []map[string]interface{}, err error) {
	path := "/api/v2/policy/systems/" + system + "/query_by_actions/"
	data, err = c.callWithReturnSliceMapData(POST, path, body, 10)
	return
}

// PolicyAuth will do policy auth
func (c *iamBackendClient) PolicyAuth(body interface{}) (data map[string]interface{}, err error) {
	path := "/api/v1/policy/auth"
	data, err = c.callWithReturnMapData(POST, path, body, 10)
	return
}

// V2PolicyAuth will do policy auth
func (c *iamBackendClient) V2PolicyAuth(system string, body interface{}) (data map[string]interface{}, err error) {
	path := "/api/v2/policy/systems/" + system + "/auth/"
	data, err = c.callWithReturnMapData(POST, path, body, 10)
	return
}

// PolicyAuthByResources will do policy auth by resources
func (c *iamBackendClient) PolicyAuthByResources(body interface{}) (data map[string]interface{}, err error) {
	path := "/api/v1/policy/auth_by_resources"
	data, err = c.callWithReturnMapData(POST, path, body, 10)
	return
}

// PolicyAuthByActions will do policy auth by actions
func (c *iamBackendClient) PolicyAuthByActions(body interface{}) (data map[string]interface{}, err error) {
	path := "/api/v1/policy/auth_by_actions"
	data, err = c.callWithReturnMapData(POST, path, body, 10)
	return
}

// PolicyGet will get the policy detail by id
func (c *iamBackendClient) PolicyGet(policyID int64) (data map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/systems/%s/policies/%d", c.System, policyID)
	data, err = c.callWithReturnMapData(GET, path, map[string]interface{}{}, 10)
	return
}

// PolicyList will list all the policy
func (c *iamBackendClient) PolicyList(body interface{}) (data map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/systems/%s/policies", c.System)
	data, err = c.callWithReturnMapData(GET, path, body, 10)
	return
}

// PolicySubjects will query the subject of each policy
func (c *iamBackendClient) PolicySubjects(policyIDs []int64) (data []map[string]interface{}, err error) {
	path := fmt.Sprintf("/api/v1/systems/%s/policies/-/subjects", c.System)

	body := map[string]interface{}{
		"ids": util.Int64ArrayToString(policyIDs, ","),
	}
	data, err = c.callWithReturnSliceMapData(GET, path, body, 10)
	return
}

// GetApplyURL will get apply url from iam saas
func (c *iamBackendClient) GetApplyURL(body interface{}) (url string, err error) {
	path := "/api/v1/open/application/"
	data, err := c.callWithReturnMapData(POST, path, body, 10)
	if err != nil {
		return "", err
	}

	urlI, ok := data["url"]
	if !ok {
		return "", errors.New("no url in response body")
	}
	url, ok = urlI.(string)
	if !ok {
		return "", errors.New("url is not a valid string")
	}
	return url, nil
}

// ModelQuery performs a model query using the specified system.
//
// system: the name of the system.
// Returns a map[string]interface{} and an error.
func (c *iamBackendClient) ModelQuery(system string) (map[string]interface{}, error) {
	if system == "" {
		system = c.System
	}
	path := fmt.Sprintf("/api/v1/model/systems/%s/query", system)
	return c.callWithReturnMapData(GET, path, map[string]interface{}{}, 10)
}

// AddSystem is a function that adds a system to the IAM backend.
//
// It takes a parameter called body which represents the system to be added.
// It returns an error if the operation fails.
func (c *iamBackendClient) AddSystem(body interface{}) error {
	path := "/api/v1/model/systems"
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateSystem updates the specified system in the IAM backend.
//
// system: The name of the system to be updated.
// body: The updated data for the system.
// error: An error if the update operation fails.
func (c *iamBackendClient) UpdateSystem(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s", system)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// AddResourceType description of the Go function.
//
// Adds a resource type to the specified system.
//
// Parameters:
// - system: The system to add the resource type to.
// - body: The body of the resource type to add.
//
// Returns:
// - error: An error if the operation fails.
func (c *iamBackendClient) AddResourceType(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/resource-types", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateResourceType updates a resource type in the IAM backend.
//
// Parameters:
//   - system: the system ID
//   - resourceTypeID: the ID of the resource type
//   - body: the body of the request
//
// Returns:
//   - error: an error if the update fails
func (c *iamBackendClient) UpdateResourceType(system, resourceTypeID string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/resource-types/%s", system, resourceTypeID)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// BatchDeleteResourceType deletes multiple resource types in the IAM backend.
//
// system: The system from which the resource types will be deleted.
// resourceTypeIDs: The IDs of the resource types to be deleted.
// error: Returns an error if the deletion fails.
func (c *iamBackendClient) BatchDeleteResourceType(system string, resourceTypeIDs ...string) error {
	if len(resourceTypeIDs) == 0 {
		return nil
	}
	path := fmt.Sprintf("/api/v1/model/systems/%s/resource-types", system)
	body := []map[string]string{}
	for _, v := range resourceTypeIDs {
		body = append(body, map[string]string{
			"id": v,
		})
	}
	_, err := c.callWithReturnMapData(DELETE, path, body, 10)
	return err
}

// AddInstanceSelection description of the Go function.
//
// Adds an instance selection for a given system.
//
// Parameters:
// - system: The name of the system.
// - body: The instance selection data.
//
// Returns:
// - error: An error if the operation fails.
func (c *iamBackendClient) AddInstanceSelection(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/instance-selections", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateInstanceSelection updates an instance selection for a system.
//
// system: the name of the system.
// instanceSelectionID: the ID of the instance selection.
// body: the data to update the instance selection with.
// Returns an error if the update fails.
func (c *iamBackendClient) UpdateInstanceSelection(system, instanceSelectionID string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/instance-selections/%s", system, instanceSelectionID)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// BatchDeleteInstanceSelection deletes multiple instance selections for a given system.
//
// Parameters:
// - system: the name of the system
// - instanceSelectionIDs: the IDs of the instance selections to delete
//
// Return type:
// - error: returns an error if the deletion fails
func (c *iamBackendClient) BatchDeleteInstanceSelection(system string, instanceSelectionIDs ...string) error {
	if len(instanceSelectionIDs) == 0 {
		return nil
	}
	path := fmt.Sprintf("/api/v1/model/systems/%s/instance-selections", system)
	body := []map[string]string{}
	for _, v := range instanceSelectionIDs {
		body = append(body, map[string]string{
			"id": v,
		})
	}
	_, err := c.callWithReturnMapData(DELETE, path, body, 10)
	return err
}

// AddAction adds an action to the specified system.
//
// system: the name of the system.
// body: the data to be sent in the request body.
// error: an error, if any, encountered during the process.
func (c *iamBackendClient) AddAction(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/actions", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateAction updates an action in the IAM backend.
//
// system: the system name.
// actionID: the ID of the action.
// body: the updated action data.
// Returns an error if the update fails.
func (c *iamBackendClient) UpdateAction(system, actionID string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/actions/%s", system, actionID)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// BatchDeleteAction deletes a batch of actions for a given system.
//
// It takes in the system string and a variadic parameter actionIDs of type string.
// It returns an error.
func (c *iamBackendClient) BatchDeleteAction(system string, actionIDs ...string) error {
	if len(actionIDs) == 0 {
		return nil
	}
	path := fmt.Sprintf("/api/v1/model/systems/%s/actions", system)
	body := []map[string]string{}
	for _, v := range actionIDs {
		body = append(body, map[string]string{
			"id": v,
		})
	}
	_, err := c.callWithReturnMapData(DELETE, path, body, 10)
	return err
}

// AddActionGroups is a function that adds action groups to a system.
//
// It takes two parameters:
// - system: a string representing the system to which the action groups will be added.
// - body: an interface{} containing the data for the action groups.
//
// It returns an error.
func (c *iamBackendClient) AddActionGroups(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/action_groups", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateActionGroups updates the action groups for a given system.
//
// It takes the following parameter(s):
// - system: a string representing the system to update the action groups for.
// - body: an interface{} representing the data to be sent in the request body.
//
// It returns an error indicating any issues encountered during the update process.
func (c *iamBackendClient) UpdateActionGroups(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/action_groups", system)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// AddResourceCreatorActions is a function that adds resource creator actions to the system.
//
// It takes in the following parameter:
// - system: a string that represents the system.
// - body: an interface{} that represents the body of the request.
//
// It returns an error.
func (c *iamBackendClient) AddResourceCreatorActions(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/resource_creator_actions", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateResourceCreatorActions updates the resource creator actions for a given system.
//
// system: the name of the system for which to update the resource creator actions.
// body: the data representing the updated resource creator actions.
// Return type: error.
func (c *iamBackendClient) UpdateResourceCreatorActions(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/resource_creator_actions", system)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// AddCommonActions adds common actions to the IAM backend client.
//
// system: The system to add common actions to.
// body: The data to be sent in the request body.
// error: An error that occurred during the function execution, if any.
func (c *iamBackendClient) AddCommonActions(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/common_actions", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateCommonActions updates the common actions for a given system.
//
// system: the name of the system to update.
// body: the new common actions configuration.
// error: an error if the update fails.
func (c *iamBackendClient) UpdateCommonActions(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/common_actions", system)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}

// AddFeatureShieldRules adds feature shield rules for a system.
//
// system: the system to add feature shield rules for.
// body: the data containing the feature shield rules.
// Returns an error if there was a problem adding the rules.
func (c *iamBackendClient) AddFeatureShieldRules(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/feature_shield_rules", system)
	_, err := c.callWithReturnMapData(POST, path, body, 10)
	return err
}

// UpdateFeatureShieldRules updates the feature shield rules for a given system.
//
// Parameters:
//   - system: the name of the system to update the feature shield rules for.
//   - body: the data containing the updated feature shield rules.
//
// Return type:
//   - error: an error if the update fails.
func (c *iamBackendClient) UpdateFeatureShieldRules(system string, body interface{}) error {
	path := fmt.Sprintf("/api/v1/model/systems/%s/configs/feature_shield_rules", system)
	_, err := c.callWithReturnMapData(PUT, path, body, 10)
	return err
}
