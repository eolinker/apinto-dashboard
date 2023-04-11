package navigation_store

import (
	"context"

	navigation_entry "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type INavigationStore interface {
	store.IBaseStore[navigation_entry.Navigation]
	SortByUUIDs(ctx context.Context, uuids []string) error
	MaxSort(ctx context.Context) (int, error)
}

type navigationStore struct {
	*store.BaseStore[navigation_entry.Navigation]
}

func (n *navigationStore) SortByUUIDs(ctx context.Context, uuids []string) error {
	sql := "UPDATE `navigation` SET `sort` = ? WHERE uuid = ?"
	for i, uuid := range uuids {
		err := n.DB(ctx).Exec(sql, i+1, uuid).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *navigationStore) MaxSort(ctx context.Context) (int, error) {
	var maxSort = 0
	err := n.DB(ctx).Table(`navigation`).Select("IFNULL(max(sort), 0)").Scan(&maxSort).Error
	return maxSort, err
}

func newNavigationStore(db store.IDB) INavigationStore {
	return &navigationStore{BaseStore: store.CreateStore[navigation_entry.Navigation](db)}
}
