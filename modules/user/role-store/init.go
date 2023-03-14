package role_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		role := newRoleStore(db)
		roleAccess := newRoleAccessStore(db)
		roleAccessLog := newRoleAccessLogStore(db)
		bean.Injection(&role)
		bean.Injection(&roleAccess)
		bean.Injection(&roleAccessLog)
	})
}
