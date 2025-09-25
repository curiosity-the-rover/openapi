package v3_1

import "github.com/sv-tools/openapi/common"

type Webhooks = map[string]*common.RefOrSpec[common.Extendable[PathItem]]

func NewWebhooks() Webhooks {
	return make(Webhooks)
}
