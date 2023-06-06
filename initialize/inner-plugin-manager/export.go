package inner_plugin_manager

import "embed"

// RegisterInnerPluginFromEmbed 保存来源为内嵌目录的内置插件
func RegisterInnerPluginFromEmbed(pluginID string, source embed.FS) error {
	//TODO 从来源解析出实现了IInnerPlugin的资源

	_innerPluginManager.setInnerPlugin(pluginID, &embedPlugin{})
	return nil
}

func EnablePlugin(pluginID string) error {
	return _innerPluginManager.enableInnerPlugin(pluginID)
}

func GetInnerPlugin(pluginID string) (IInnerPlugin, bool) {
	return _innerPluginManager.getInnerPlugin(pluginID)
}
