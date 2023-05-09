package audit

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"time"
)

type IAuditLogService interface {
	GetLogsList(ctx context.Context, namespaceID, operateType int, kind, keyword string, start, end int64, pageNum, pageSize int) ([]*audit_model.LogListItem, int, error)
	GetLogDetail(ctx context.Context, logID int) ([]*audit_model.LogDetailArg, error)
	Log(namespace int, userId int, operate int, kind string, url, object, ip, userAgent, body, err string, start, end time.Time)
}
