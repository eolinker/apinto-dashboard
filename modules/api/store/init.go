package api_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func InitStoreDB(db store.IDB) {
	api := NewAPIStore(db)
	apiStat := NewAPIStatStore(db)
	apiRuntime := NewApiRuntimeStore(db)
	apiVersion := newAPIVersionStore(db)
	apiHistory := newApiHistoryStore(db)

	bean.Injection(&api)
	bean.Injection(&apiStat)
	bean.Injection(&apiRuntime)
	bean.Injection(&apiVersion)
	bean.Injection(&apiHistory)
}
