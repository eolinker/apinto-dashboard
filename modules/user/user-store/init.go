package user_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		userInfo := newUserInfoStore(db)
		userRole := newUserRoleStore(db)
		bean.Injection(&userInfo)

		bean.Injection(&userRole)
	})
}
