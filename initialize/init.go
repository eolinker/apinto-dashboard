package initialize

import "github.com/eolinker/apinto-dashboard/store"

func init() {
	store.RegisterStore(func(db store.IDB) {
		InitNavigation(db)
		InitPlugins(db)
		InitMiddleware(db)
	})
}
