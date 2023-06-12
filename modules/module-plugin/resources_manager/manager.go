package resources_manager

import (
	"github.com/eolinker/eosc"
)

var _innerPluginManager *innerPluginManager

func init() {
	_innerPluginManager = newEmbedPluginManager()
}

type IPluginResources interface {
	ICon() ([]byte, bool)
	RM() ([]byte, bool)
	ReadMe(name string) ([]byte, bool)
	Resources(path string) ([]byte, bool)
}

type innerPluginManager struct {
	externPluginResources eosc.Untyped[string, IPluginResources]
	embedPluginResources  eosc.Untyped[string, IPluginResources]
}

func newEmbedPluginManager() *innerPluginManager {
	return &innerPluginManager{
		externPluginResources: eosc.BuildUntyped[string, IPluginResources](),
		embedPluginResources:  eosc.BuildUntyped[string, IPluginResources](),
	}
}

func StoreExternPluginResources(id string, resource IPluginResources) {
	_innerPluginManager.externPluginResources.Set(id, resource)
}

func StoreEmbedPluginResources(id string, resource IPluginResources) {
	_innerPluginManager.embedPluginResources.Set(id, resource)
}

func GetExternPluginResources(id string) (IPluginResources, bool) {
	return _innerPluginManager.externPluginResources.Get(id)
}

func GetEmbedPluginResources(id string) (IPluginResources, bool) {
	return _innerPluginManager.embedPluginResources.Get(id)
}
