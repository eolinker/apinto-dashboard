package audit_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	auditLog := newAuditLogService()
	bean.Injection(&auditLog)
}
