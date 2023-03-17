package discovery_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		discovery := newDiscoveryStore(db)
		discoveryVersionStore := newDiscoveryVersionStore(db)
		discoveryStatStore := newDiscoveryStatStore(db)
		discoveryRuntime := newDiscoveryRuntimeStore(db)
		discoveryHistory := newDiscoveryHistoryStore(db)
		bean.Injection(&discovery)
		bean.Injection(&discoveryVersionStore)
		bean.Injection(&discoveryStatStore)
		bean.Injection(&discoveryRuntime)
		bean.Injection(&discoveryHistory)
	})
}
