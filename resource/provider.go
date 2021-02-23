package resource

// Provider is the interface for provider
type Provider interface {
	ListAttr(req Request) Response
	ListAttrValue(req Request) Response
	ListInstance(req Request) Response
	FetchInstanceInfo(req Request) Response
	ListInstanceByPolicy(req Request) Response
	SearchInstance(req Request) Response
}
