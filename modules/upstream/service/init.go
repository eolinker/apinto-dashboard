package upstream_service

import (
	upstream_store "github.com/eolinker/apinto-dashboard/modules/upstream/store"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(upstream_store.InitStoreHandler)

	upstreamService := newServiceService()
	bean.Injection(&upstreamService)
}
