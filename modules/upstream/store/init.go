package upstream_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func InitStoreHandler(db store.IDB) {
	service := newServiceStore(db)
	serviceVersionStore := newServiceVersionStore(db)
	serviceStatStore := newServiceStatStore(db)
	serviceRuntime := newServiceRuntimeStore(db)
	serviceHistory := newServiceHistoryStore(db)

	bean.Injection(&serviceRuntime)

	bean.Injection(&service)

	bean.Injection(&serviceRuntime)
	bean.Injection(&serviceVersionStore)
	bean.Injection(&serviceHistory)
	bean.Injection(&serviceStatStore)
}
