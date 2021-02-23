// +build !codeanalysis

package main

import (
	"net/http"

	"github.com/TencentBlueKing/iam-go-sdk/resource"
)

// DummyProvider is an example of provider
type DummyProvider struct {
}

// ListAttr implements the list_attr
func (d DummyProvider) ListAttr(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// ListAttrValue implements the list_attr_value
func (d DummyProvider) ListAttrValue(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// ListInstance implements the list_instance
func (d DummyProvider) ListInstance(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// FetchInstanceInf implements the fetch_instance_info
func (d DummyProvider) FetchInstanceInfo(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// ListInstanceByPolicy implements the list_instance_by_policy
func (d DummyProvider) ListInstanceByPolicy(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

// SearchInstance implements the search_instance
func (d DummyProvider) SearchInstance(req resource.Request) resource.Response {
	return resource.Response{
		Code:    0,
		Message: req.Method,
	}
}

func main() {
	d := resource.NewDispatcher()
	dummyProvider := DummyProvider{}

	// type=dummy will use the dummyProvider
	d.RegisterProvider("dummy", dummyProvider)

	handler := resource.NewDispatchHandler(d)

	// we just register only one api to handle all resource types - callback api
	http.HandleFunc("/api/v1/resource", handler)

}
