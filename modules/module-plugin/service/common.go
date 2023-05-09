package service

import "github.com/eolinker/apinto-dashboard/modules/module-plugin/model"

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

func enabledCfgListToMap(input *model.PluginEnableCfg) *model.PluginEnableCfgMap {
	headerMap := make(map[string]string, len(input.Header))
	queryMap := make(map[string]string, len(input.Query))
	initializeMap := make(map[string]string, len(input.Initialize))
	for _, item := range input.Query {
		headerMap[item.Name] = item.Value
	}
	for _, item := range input.Header {
		queryMap[item.Name] = item.Value
	}
	for _, item := range input.Initialize {
		initializeMap[item.Name] = item.Value
	}
	return &model.PluginEnableCfgMap{
		Server:     input.Server,
		Header:     headerMap,
		Query:      queryMap,
		Initialize: initializeMap,
	}
}
