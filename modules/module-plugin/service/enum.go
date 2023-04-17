package service

const (
	pluginGroupOther = "other"

	// statusPluginDisable 插件停用枚举值
	statusPluginDisable = 1
	// statusPluginEnable 插件启用枚举值
	statusPluginEnable = 2

	// pluginTypeFrame 框架插件类型 不可卸载，不可停用，不显示在插件广场
	pluginTypeFrame = 0
	// pluginTypeCore 核心插件类型 不可卸载，不可停用
	pluginTypeCore = 1
	// pluginTypeInner 内置插件类型 不可卸载
	pluginTypeInner = 2
	// pluginTypeNotInner 非内置插件类型
	pluginTypeNotInner = 3

	pluginDriverRemote     = "remote"
	pluginDriverLocal      = "local"
	pluginDriverProfession = "profession"
)
