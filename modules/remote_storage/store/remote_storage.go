package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/remote_storage/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IRemoteStorage interface {
	Get(ctx context.Context, module, key string) (*entry.RemoteKeyObject, error)
	Save(ctx context.Context, object *entry.RemoteKeyObject) error
}
type remoteStorage struct {
	store.IBaseStore[entry.RemoteKeyObject]
}

func (c *remoteStorage) List(ctx context.Context, secret string) ([]*entry.RemoteKeyObject, error) {
	list, err := c.IBaseStore.List(ctx, map[string]interface{}{"secret": secret})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (c *remoteStorage) Get(ctx context.Context, module, key string) (*entry.RemoteKeyObject, error) {
	return c.IBaseStore.First(ctx, map[string]interface{}{"key": key, "module": module})
}

func newRemoteStorageStore(db store.IDB) IRemoteStorage {
	err := db.DB(context.Background()).AutoMigrate(&entry.RemoteKeyObject{})
	if err != nil {
		panic(err)
	}
	base := store.CreateStore[entry.RemoteKeyObject](db)

	return &remoteStorage{IBaseStore: base}
}
