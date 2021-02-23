package client

import (
	"net/http"

	"github.com/parnurzeal/gorequest"
	"moul.io/http2curl"

	"github.com/TencentBlueKing/iam-go-sdk/logger"
)

// AsCurlCommand returns a string representing the runnable `curl' command
// version of the request.
func AsCurlCommand(request *gorequest.SuperAgent) (string, error) {
	req, err := request.MakeRequest()
	if err != nil {
		return "", err
	}

	// 脱敏, 去掉-H 中 Authorization
	req.Header.Del("Authorization")

	cmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return "", err
	}
	return cmd.String(), nil
}

func logFailHTTPRequest(request *gorequest.SuperAgent, response gorequest.Response, errs []error, data responseBody) {

	dump, err := AsCurlCommand(request)
	if err != nil {
		logger.Errorf("component request AsCurlCommand fail, error: %v", err)
	}

	status := -1
	requestID := ""
	if response != nil {
		status = response.StatusCode
		requestID = response.Header.Get("X-Request-Id")
	}

	responseBodyError := data.Error()

	if len(errs) != 0 || status != http.StatusOK || responseBodyError != nil {
		message := "-"
		if responseBodyError != nil {
			message = responseBodyError.Error()
		}
		logger.Errorf("[http request fail] %s! status=`%d`, errs=`%v`, request_id=`%s`, request=`%s`",
			message, status, errs, requestID, dump)
	} else {
		logger.Debugf("[http request] %s! status=`%d`, errs=`%v`, request_id=`%s`, request=`%s`",
			status, errs, requestID, dump)
	}
}
