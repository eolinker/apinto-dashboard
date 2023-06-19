/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package plugin_group

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
