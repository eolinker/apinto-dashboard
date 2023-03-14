package audit_store

import (
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	store.RegisterStore(func(db store.IDB) {
		auditLog := newAuditLogStore(db)
		bean.Injection(&auditLog)
	})
}
