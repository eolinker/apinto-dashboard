package namespace_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		namespace := newNamespaceStore(db)

		bean.Injection(&namespace)
	})
}
