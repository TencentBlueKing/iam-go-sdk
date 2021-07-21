## Usage

### create an iam instance first

```go
    import "github.com/TencentBlueKing/iam-go-sdk"
	i := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")
```

如果是使用 APIGateway 的方式

  
```go
    import "github.com/TencentBlueKing/iam-go-sdk"
	// if your TencentBlueking has a APIGateway, use NewAPIGatewayIAM, the url suffix is /stage/(for testing) and /prod/(for production)
	i := iam.NewAPIGatewayIAM("bk_paas", "bk_paas", "{app_secret}", "http://bk-iam.{APIGATEWAY_DOMAIN}/stage/")
```


### IsAllowed

```go
	req := iam.NewRequest(
		"bk_paas",
		iam.NewSubject("user", "admin"),
		iam.NewAction("access_developer_center"),
		[]iam.ResourceNode{},
	)

	allowed, err := i.IsAllowed(req)
	fmt.Println("isAllowed:", allowed, err)
```

### IsAllowedWithCache

```go
	// check 3 times but only call iam backend once
	allowed, err = i.IsAllowedWithCache(req, 10*time.Second)
	allowed, err = i.IsAllowedWithCache(req, 10*time.Second)
	i2 := iam.NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")
	allowed, err = i2.IsAllowedWithCache(req, 10*time.Second)
	fmt.Println("isAllowedWithCache:", allowed, err)
```


### BatchIsAllowed

```go
	req := iam.NewRequest(
		"bk_paas",
		iam.NewSubject("user", "admin"),
		iam.NewAction("develop_app"),
        // NOTE: here is empty resource
		[]iam.ResourceNode{},
	)

	resourcesList := []iam.Resources{
		[]iam.ResourceNode{
			{
				System: "bk_paas",
				Type:   "app",
				ID:     "test1",
			},
		},
	}

    result, err := i.BatchIsAllowed(req, resourcesList)
	fmt.Println("BatchIsAllowed:", result, err)

```

### ResourceMultiActionsAllowed

```go
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
	result, err := i.ResourceMultiActionsAllowed(multiReq)
	fmt.Println("ResourceMultiActionsAllowed: ", result, err)

```

### BatchResourceMultiActionsAllowed

```go
	multiReq := iam.NewMultiActionRequest(
		"bk_sops",
		iam.NewSubject("user", "admin"),
		[]iam.Action{
			iam.NewAction("task_delete"),
			iam.NewAction("task_edit"),
			iam.NewAction("task_view"),
		},
        // NOTE: here is empty resource
		[]iam.ResourceNode{},
	)
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
	results, err := i.BatchResourceMultiActionsAllowed(multiReq, resourcesList)
	fmt.Println("BatchResourceMultiActionsAllowed: ", results, err)

```


### GetApplyURL

```go
	actions := []iam.ApplicationAction{
		iam.NewApplicationAction("access_developer_center", []iam.ApplicationRelatedResourceType{}),
	}
	application := iam.NewApplication("bk_paas", actions)

	url, err := i.GetApplyURL(application, "", "admin")
	fmt.Println("GetApplyURL:", url, err)
```

### GenPermissionApplyData

```go
    d := ApplicationActionListForApply{}
    data, err := i.GenPermissionApplyData(d)
```

### IsBasicAuthAllowed

```go
    // get and parse the basic auth from http header
	err := c.IsBasicAuthAllowed("bk_iam", "theToken")
	fmt.Println("IsBasicAuthAllowed:", err)
```

### GetToken


```go
    token, err := i.GetToken()
    fmt.Println("GetToken:", token, err)
```

### enable prometheus metrics

```go
    import "github.com/TencentBlueKing/iam-go-sdk/metric"
    metric.RegisterMetrics()
```

### Implement resource callback api via dispatcher/provider interface

see examples/dispatcher_provider/main.go

```go
func main() {
	d := resource.NewDispatcher()
	dummyProvider := DummyProvider{}

	// type=dummy will use the dummyProvider
	d.RegisterProvider("dummy", dummyProvider)

	handler := resource.NewDispatchHandler(d)

	// we just register only one api to handle all resource types - callback api
	http.HandleFunc("/api/v1/resource", handler)
}
```
