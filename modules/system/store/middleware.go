package store

import (
	"github.com/eolinker/apinto-dashboard/modules/system/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type ISystemStore interface {
	store.IBaseStore[entry.System]
}

type SystemStore struct {
	*store.BaseStore[entry.System]
}

func newStore(db store.IDB) ISystemStore {
	return &SystemStore{BaseStore: store.CreateStore[entry.System](db)}
}
