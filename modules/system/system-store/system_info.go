package system_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/system/system-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

var (
	_ ISystemInfoStore = (*systemInfoStore)(nil)
)

type ISystemInfoStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
}

type systemInfoStore struct {
	store.IBaseStore[system_entry.SystemInfo]
}

func (s *systemInfoStore) Get(ctx context.Context, key string) ([]byte, error) {

	ent, err := s.IBaseStore.First(ctx, map[string]any{"key": key})
	if err != nil {
		return nil, err
	}
	return ent.Value, nil
}

func (s *systemInfoStore) Set(ctx context.Context, key string, value []byte) error {
	err := s.IBaseStore.Save(ctx, &system_entry.SystemInfo{
		Id:    0,
		Key:   key,
		Value: value,
	})
	return err
}

func newSystemInfoStore(db store.IDB) ISystemInfoStore {
	return &systemInfoStore{IBaseStore: store.CreateStore[system_entry.SystemInfo](db)}
}
