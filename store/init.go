package store

import (
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

func InitStoreDB(idb *gorm.DB) {
	var db IDB = &myDB{db: idb}
	bean.Injection(&db)
	runHandler(db)

}
