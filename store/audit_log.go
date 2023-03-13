package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry"
	"time"
)

var (
	_ IAuditLogStore = (*auditLogStore)(nil)
)

type IAuditLogStore interface {
	IBaseStore[entry.AuditLog]
	GetLogsByCondition(ctx context.Context, namespaceID, operateType int, kind, keyword string, start, end int64, pageNum, pageSize int) ([]*entry.AuditLog, int, error)
}

type auditLogStore struct {
	*baseStore[entry.AuditLog]
}

func newAuditLogStore(db IDB) IAuditLogStore {
	return &auditLogStore{baseStore: createStore[entry.AuditLog](db)}
}

func (c *auditLogStore) GetLogsByCondition(ctx context.Context, namespaceID, operateType int, kind, keyword string, start, end int64, pageNum, pageSize int) ([]*entry.AuditLog, int, error) {
	logs := make([]*entry.AuditLog, 0)
	db := c.DB(ctx).Where("`namespace` = ?", namespaceID)
	count := int64(0)
	if operateType > 0 {
		db = db.Where("`operate` = ?", operateType)
	}
	if kind != "" {
		db = db.Where("`kind` = ?", kind)
	}
	//TODO
	if keyword != "" {
		db = db.Where("`username` like ? or `ip` like ? ", "%"+keyword+"%", "%"+keyword+"%")
	}
	if start != 0 {
		startTime := time.Unix(start, 0).UTC()
		db = db.Where("`start_time` >= ? ", startTime) //start_time为操作时间
	}
	if end != 0 {
		endTime := time.Unix(end, 0).UTC()
		db = db.Where("`start_time` <= ? ", endTime) //start_time为操作时间
	}

	err := db.Model(logs).Count(&count).Order("start_time DESC").Limit(pageSize).Offset(entry.PageIndex(pageNum, pageSize)).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, int(count), nil
}
