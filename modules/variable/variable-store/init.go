package variable_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		globalVariable := newGlobalVariableStore(db)
		clusterVariable := newClusterVariableStore(db)

		variableHistory := newVariableHistoryStore(db)
		variablePublishVersion := newVariablePublishVersionStore(db)
		variableRuntime := newVariableRuntimeStore(db)
		variablePublishHistory := newVariablePublishHistoryStore(db)
		bean.Injection(&variableHistory)
		bean.Injection(&variablePublishVersion)
		bean.Injection(&variableRuntime)
		bean.Injection(&variablePublishHistory)
		bean.Injection(&clusterVariable)
		bean.Injection(&globalVariable)
	})
}
