package inner_plugin_manager

import (
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/bean"
	"sync"
)

var _innerPluginManager *innerPluginManager

func init() {
	_innerPluginManager = newEmbedPluginManager()
}

type innerPluginManager struct {
	data                eosc.Untyped[string, IInnerPlugin]
	modulePluginService module_plugin.IModulePluginService
	sync.Mutex
}

func newEmbedPluginManager() *innerPluginManager {
	manager := &innerPluginManager{
		data: eosc.BuildUntyped[string, IInnerPlugin](),
	}
	bean.Autowired(&manager.modulePluginService)
	return manager
}

func (e *innerPluginManager) getInnerPlugin(pluginID string) (IInnerPlugin, bool) {
	return _innerPluginManager.data.Get(pluginID)
}

func (e *innerPluginManager) setInnerPlugin(pluginID string, plugin IInnerPlugin) {
	e.Lock()
	defer e.Unlock()
	e.data.Set(pluginID, plugin)
	//TODO 检查数据库中有没有该内置插件，没有则插入, 有则判断版本有没有变化
}

func (e *innerPluginManager) enableInnerPlugin(pluginID string) error {
	e.Lock()
	defer e.Unlock()
	//TODO
	//return e.modulePluginService.EnablePlugin()
	return nil
}
