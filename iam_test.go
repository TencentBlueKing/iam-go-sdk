package iam

import (
	"github.com/stretchr/testify/assert"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("iam", func() {

	Context("request.CacheKey", func() {

		It("ok", func() {
			request := NewRequest("system", NewSubject("type", "id"), NewAction("id"), []ResourceNode{
				NewResourceNode("system", "type", "id", map[string]interface{}{}),
			})

			key, err := request.CacheKey()

			assert.NoError(GinkgoT(), err)
			assert.Equal(
				GinkgoT(),
				key,
				"iam:9b9893808246fb3c7fbc0192a771d67e",
			)
		})
	})

	Context("iam.buildResourceID", func() {
		var iam = NewIAM("bk_paas", "bk_paas", "{app_secret}", "http://{iam_backend_addr}", "http://{paas_domain}")

		It("one node", func() {
			resources := []ResourceNode{NewResourceNode("system", "type", "id", map[string]interface{}{})}

			resourceID := iam.buildResourceID(resources)

			assert.Equal(GinkgoT(), resourceID, "id")
		})

		It("two nodes", func() {
			resources := []ResourceNode{
				NewResourceNode("system", "type", "id", map[string]interface{}{}),
				NewResourceNode("system", "type2", "id2", map[string]interface{}{}),
			}

			resourceID := iam.buildResourceID(resources)

			assert.Equal(GinkgoT(), resourceID, "type,id/type2,id2")
		})
	})
})
