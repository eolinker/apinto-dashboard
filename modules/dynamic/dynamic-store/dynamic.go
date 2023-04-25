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
	ListPageByKeyword(ctx context.Context, params map[string]interface{}, keyword string, page int, pageSize int) ([]*dynamic_entry.Dynamic, error)
}

type dynamicStore struct {
	*store.BaseStore[dynamic_entry.Dynamic]
}

func (d *dynamicStore) ListPageByKeyword(ctx context.Context, params map[string]interface{}, keyword string, page int, pageSize int) ([]*dynamic_entry.Dynamic, error) {
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
	result, _, err := d.ListPage(ctx, builder.String(), page, pageSize, args, "update_time")
	return result, err
}

func newDynamicStore(db store.IDB) IDynamicStore {
	return &dynamicStore{BaseStore: store.CreateStore[dynamic_entry.Dynamic](db)}
}
