package enum

import (
	"fmt"
)

type PluginStateType int

const (
	PluginStateTypeNone = iota
	PluginStateTypeDisable
	PluginStateTypeEnable
	PluginStateTypeGlobal
	PluginStateTypeAll
)

const (
	PluginStateNameDisable = "DISABLE"
	PluginStateNameEnable  = "ENABLE"
	PluginStateNameGlobal  = "GLOBAL"
)

var (
	pluginStateNames = map[PluginStateType]string{
		PluginStateTypeNone:    "",
		PluginStateTypeDisable: PluginStateNameDisable, //禁用
		PluginStateTypeEnable:  PluginStateNameEnable,  //启用
		PluginStateTypeGlobal:  PluginStateNameGlobal,  //全局启用
	}
	pluginStateIndex = map[string]PluginStateType{
		PluginStateNameDisable: PluginStateTypeDisable,
		PluginStateNameEnable:  PluginStateTypeEnable,
		PluginStateNameGlobal:  PluginStateTypeGlobal,
	}
)

func (p *PluginStateType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint("\"", p.String(), "\"")), nil
}

func (p *PluginStateType) String() string {
	if *p >= PluginStateTypeAll {
		return "unknown"
	}
	return pluginStateNames[*p]
}

func GetPluginState(stateName string) (int, bool) {
	state, exist := pluginStateIndex[stateName]
	if !exist {
		return 0, false
	}
	return int(state), true
}
