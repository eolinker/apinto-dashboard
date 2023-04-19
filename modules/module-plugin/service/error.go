package service

import "errors"

var (
	ErrModulePluginNotFound    = errors.New("插件不存在")
	ErrModulePluginInstalled   = errors.New("插件已安装")
	ErrModulePluginHasDisabled = errors.New("插件已停用")
	ErrModulePluginHasEnabled  = errors.New("插件已启用")
)
