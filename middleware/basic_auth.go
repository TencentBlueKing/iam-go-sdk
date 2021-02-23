package middleware

import (
	"net/http"

	"github.com/TencentBlueKing/iam-go-sdk"
)

// NewIAMBasicAuth will create a middleware for http server, check the callback request
func NewIAMBasicAuth(i *iam.IAM) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				http.Error(w, "basic auth not provided or parseBasicAuth fail", http.StatusForbidden)
				return
			}

			err := i.IsBasicAuthAllowed(username, password)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
