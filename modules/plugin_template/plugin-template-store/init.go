package plugin_template_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		runtimeStore := newPluginTemplateRuntimeStore(db)
		historyStore := newPluginTemplateHistoryStore(db)
		templateStore := newPluginTemplateStore(db)
		statStore := newPluginTemplateStatStore(db)
		versionStore := newPluginTemplateVersionStore(db)
		publishHistory := newPluginTemplatePublishHistoryStore(db)

		bean.Injection(&runtimeStore)
		bean.Injection(&historyStore)
		bean.Injection(&templateStore)
		bean.Injection(&statStore)
		bean.Injection(&versionStore)
		bean.Injection(&publishHistory)
	})
}
