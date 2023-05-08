package dynamic_store

import (
	"context"
	"fmt"
	"strings"

	dynamic_entry "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IDynamicStore interface {
	store.IBaseStore[dynamic_entry.Dynamic]
	ListPageByKeyword(ctx context.Context, params map[string]interface{}, drivers []string, keyword string, page int, pageSize int) ([]*dynamic_entry.Dynamic, int, error)
	ListByKeyword(ctx context.Context, params map[string]interface{}, names []string, keyword string) ([]*dynamic_entry.Dynamic, error)
	Count(ctx context.Context, params map[string]interface{}) (int, error)
}

type dynamicStore struct {
	*store.BaseStore[dynamic_entry.Dynamic]
}

func (d *dynamicStore) Count(ctx context.Context, params map[string]interface{}) (int, error) {
	var count int64
	err := d.DB(ctx).Where(params).Count(&count).Error
	return int(count), err
}

func (d *dynamicStore) ListPageByKeyword(ctx context.Context, params map[string]interface{}, drivers []string, keyword string, page int, pageSize int) ([]*dynamic_entry.Dynamic, int, error) {
	isInit := false
	builder := strings.Builder{}
	args := make([]interface{}, 0, len(params))
	for key, value := range params {
		if isInit {
			builder.WriteString(" AND ")
		}
		builder.WriteString(fmt.Sprintf("`%s` = ? ", key))
		args = append(args, value)
		isInit = true
	}

	if keyword != "" {
		builder.WriteString(" AND `title` LIKE ?")
		args = append(args, "%"+keyword+"%")
	}
	if len(drivers) > 0 {
		builder.WriteString(fmt.Sprintf(" AND `driver` IN (\"%s\")", strings.Join(drivers, "\",\"")))
	}

	return d.ListPage(ctx, builder.String(), page, pageSize, args, "update_time")
}

func (d *dynamicStore) ListByKeyword(ctx context.Context, params map[string]interface{}, names []string, keyword string) ([]*dynamic_entry.Dynamic, error) {
	isInit := false
	builder := strings.Builder{}
	args := make([]interface{}, 0, len(params))
	for key, value := range params {
		if isInit {
			builder.WriteString(" AND ")
		}
		builder.WriteString(fmt.Sprintf("`%s` = ? ", key))
		args = append(args, value)
		isInit = true
	}

	if keyword != "" {
		builder.WriteString(" AND `title` LIKE ?")
		args = append(args, "%"+keyword+"%")
	}
	if len(names) > 0 {
		builder.WriteString(fmt.Sprintf(" AND `name` IN (\"%s\")", strings.Join(names, "\",\"")))
	}

	return d.ListQuery(ctx, builder.String(), args, "update_time")
}

func newDynamicStore(db store.IDB) IDynamicStore {
	return &dynamicStore{BaseStore: store.CreateStore[dynamic_entry.Dynamic](db)}
}
