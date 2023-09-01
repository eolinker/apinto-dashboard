package initialize

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed plugin_group.yml
	pluginGroupsContent []byte

	groupList []*PluginGroupsItem
)

type PluginGroupsItem struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

func init() {
	// 初始化插件分组列表
	groupList = make([]*PluginGroupsItem, 0)
	err := yaml.Unmarshal(pluginGroupsContent, &groupList)
	if err != nil {
		panic(err)
	}
}

func GetModulePluginGroups() []*PluginGroupsItem {
	return groupList
}
