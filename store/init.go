package store

import (
	"github.com/eolinker/eosc/common/bean"
)

func InitStoreDB(db IDB) {

	bean.Injection(&db)
	runHandler(db)

}
