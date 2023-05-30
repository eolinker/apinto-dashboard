package application_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		application := newApplicationStore(db)
		applicationStat := newApplicationStatStore(db)
		applicationVersion := newApplicationVersionStore(db)
		applicationAuth := newApplicationAuthStore(db)
		applicationHistory := newApplicationHistoryStore(db)
		appPublishHistory := newAppPublishHistoryStore(db)

		bean.Injection(&application)
		bean.Injection(&applicationVersion)
		bean.Injection(&applicationStat)
		bean.Injection(&applicationAuth)
		bean.Injection(&applicationHistory)
		bean.Injection(&appPublishHistory)
	})
}
