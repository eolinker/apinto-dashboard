package application_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		application := newApplicationStore(db)
		applicationRuntime := newApplicationRuntimeStore(db)
		applicationStat := newApplicationStatStore(db)
		applicationVersion := newApplicationVersionStore(db)
		applicationAuth := newApplicationAuthStore(db)
		applicationAuthVersion := newApplicationAuthVersionStore(db)
		applicationAuthStat := newApplicationAuthStatStore(db)
		applicationAuthRuntimeStore := newApplicationAuthRuntimeStore(db)
		applicationAuthPublish := newApplicationAuthPublishStore(db)
		applicationHistory := newApplicationHistoryStore(db)
		applicationAuthHistory := newApplicationAuthHistoryStore(db)
		bean.Injection(&application)
		bean.Injection(&applicationRuntime)
		bean.Injection(&applicationVersion)
		bean.Injection(&applicationStat)
		bean.Injection(&applicationAuth)
		bean.Injection(&applicationAuthVersion)
		bean.Injection(&applicationAuthStat)
		bean.Injection(&applicationAuthRuntimeStore)
		bean.Injection(&applicationAuthPublish)
		bean.Injection(&applicationHistory)
		bean.Injection(&applicationAuthHistory)
	})
}
