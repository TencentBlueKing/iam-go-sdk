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
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

// ESBClient is the interface of esb
type ESBClient interface {
	GetApplyURL(bkToken string, bkUsername string, body interface{}) (string, error)
	// NOTE: will remove soon, change all API to APIGateway
	//       so you should not add more interface here!!!!!!
}

type esbClient struct {
	Host string

	appCode   string
	appSecret string
}

// NewESBClient will create a esb client
func NewESBClient(host string, appCode string, appSecret string) ESBClient {
	return &esbClient{
		Host: host,

		appCode:   appCode,
		appSecret: appSecret,
	}
}

// ESBResponse is the struct of esb response
type ESBResponse struct {
	Code    int                    `json:"code"`
	Result  bool                   `json:"result"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// Error will check if response fail
func (r *ESBResponse) Error() error {
	if r.Code == 0 {
		return nil
	}

	return fmt.Errorf("response error[code=`%d`,  message=`%s`]", r.Code, r.Message)
}

func (c *esbClient) call(
	method Method,
	path string,
	data interface{},
	timeout int64,
	bkToken string,
	bkUsername string,
) (map[string]interface{}, error) {
	callTimeout := time.Duration(timeout) * time.Second
	if timeout == 0 {
		callTimeout = defaultTimeout
	}

	headers := map[string]string{}

	url := fmt.Sprintf("%s%s", c.Host, path)
	result := ESBResponse{}
	start := time.Now()
	callbackFunc := NewMetricCallback("ESB", start)

	request := gorequest.New().Timeout(callTimeout).Type("json")
	switch method {
	case POST:
		request = request.Post(url)
	case GET:
		request = request.Get(url)
	}

	// set headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// NOTE: set s.Data
	request.Data["bk_app_code"] = c.appCode
	request.Data["bk_app_secret"] = c.appSecret
	request.Data["bk_token"] = bkToken
	request.Data["bk_username"] = bkUsername

	// do request
	resp, _, errs := request.
		Send(data).
		EndStruct(&result, callbackFunc)

	logFailHTTPRequest(request, resp, errs, &result)

	if len(errs) != 0 {
		return nil, fmt.Errorf("gorequest errors=`%s`", errs)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gorequest statusCode is %d not 200", resp.StatusCode)
	}
	if result.Code != 0 {
		return nil, errors.New(result.Message)
	}

	return result.Data, nil
}

// GetApplyURL will get apply url from iam saas
func (c *esbClient) GetApplyURL(bkToken string, bkUsername string, body interface{}) (url string, err error) {
	path := "/api/c/compapi/v2/iam/application/"
	data, err := c.call(POST, path, body, 10, bkToken, bkUsername)
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
