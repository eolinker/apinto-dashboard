package service

// IsInnerPlugin 根据插件类型判断是否为内置插件
func IsInnerPlugin(pluginType int) bool {
	switch pluginType {
	case pluginTypeFrame, pluginTypeCore, pluginTypeInner:
		return true
	case pluginTypeNotInner:
		return false
	default:
		return false
	}
}

// IsPluginCanDisable 根据插件类型判断插件是否能停用
func IsPluginCanDisable(pluginType int) bool {
	//框架和核心插件不可以停用
	if pluginType == pluginTypeFrame || pluginType == pluginTypeCore {
		return false
	}
	return true
}
