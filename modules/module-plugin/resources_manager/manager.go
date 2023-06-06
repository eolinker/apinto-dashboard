package resources_manager

import (
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/eosc"
)

var _innerPluginManager *innerPluginManager

func init() {
	_innerPluginManager = newEmbedPluginManager()
}

type innerPluginManager struct {
	externPluginResources eosc.Untyped[string, *model.PluginResources]
	embedPluginResources  eosc.Untyped[string, *model.EmbedPluginResources]
}

func newEmbedPluginManager() *innerPluginManager {
	return &innerPluginManager{
		externPluginResources: eosc.BuildUntyped[string, *model.PluginResources](),
		embedPluginResources:  eosc.BuildUntyped[string, *model.EmbedPluginResources](),
	}
}

func StoreExternPluginResources(id string, resource *model.PluginResources) {
	_innerPluginManager.externPluginResources.Set(id, resource)
}

func StoreEmbedPluginResources(id string, resource *model.EmbedPluginResources) {
	_innerPluginManager.embedPluginResources.Set(id, resource)
}

func GetExternPluginResources(id string) (*model.PluginResources, bool) {
	return _innerPluginManager.externPluginResources.Get(id)
}

func GetEmbedPluginResources(id string) (*model.EmbedPluginResources, bool) {
	return _innerPluginManager.embedPluginResources.Get(id)
}
