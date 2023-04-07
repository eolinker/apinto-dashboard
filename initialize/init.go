package initialize

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {

	var db store.IDB
	bean.Autowired(&db)
	bean.AddInitializingBeanFunc(func() {
		initNavigation(db)
		initPlugins(db)
		initMiddleware(db)
	})

}
