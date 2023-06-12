package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type IPluginResources interface {
	store.IBaseStore[entry.PluginResources]
	GetByUUID(ctx context.Context, uuid string) (*entry.PluginResources, error)
}

type pluginResources struct {
	store.IBaseStore[entry.PluginResources]
}

func (c *pluginResources) GetByUUID(ctx context.Context, uuid string) (*entry.PluginResources, error) {
	return c.IBaseStore.First(ctx, map[string]interface{}{"uuid": uuid})
}

func newPluginResourcesStore(db store.IDB) IPluginResources {
	err := db.DB(context.Background()).AutoMigrate(&entry.PluginResources{})
	if err != nil {
		panic(err)
	}
	return &pluginResources{IBaseStore: store.CreateStore[entry.PluginResources](db)}
}
