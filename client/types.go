package client

import (
	"time"
)

const (
	defaultTimeout = 5 * time.Second
)

// responseBody is the interface for fail http response log
type responseBody interface {
	Error() error
}
