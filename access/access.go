package access

import (
	_ "embed"
	"fmt"
	"gopkg.in/yaml.v3"
)

type Access int

const (
	ClusterView Access = iota
	ClusterEdit
	VariableView
	VariableEdit
	ServiceView
	ServiceEdit
	DiscoveryView
	DiscoveryEdit
	ApplicationView
	ApplicationEdit
	ApiView
	ApiEdit
	StrategyTrafficView
	StrategyTrafficEdit
	StrategyFuseView
	StrategyFuseEdit
	StrategyVisitView
	StrategyVisitEdit
	StrategyCacheView
	StrategyCacheEdit
	StrategyGreyView
	StrategyGreyEdit
	AuditLogView
	ExtAPPView
	ExtAPPEdit
	PluginView
	PluginEdit
	PluginTemplateView
	PluginTemplateEdit
	lastId
	unknown
)

//go:embed config/access.yml
var accessData []byte

var (
	allAccess      []*GlobalAccess
	allModules     []*ModuleItem
	maxAccessDepth int
)

type AccessConfig struct {
	Modules []*AccessConfigItem `yaml:"modules"`
}

type AccessConfigItem struct {
	ID       int                 `yaml:"id"`                 //模块ID
	Title    string              `yaml:"title"`              //模块名称
	Module   string              `yaml:"module,omitempty"`   //模块路径
	Access   []string            `yaml:"access,omitempty"`   //模块权限标志
	Children []*AccessConfigItem `yaml:"children,omitempty"` //子模块
}

// GlobalAccess  用户中心所需要的权限
type GlobalAccess struct {
	ID       int             //模块ID
	Title    string          //模块名称
	Module   string          //模块路径
	Access   []*AccessItem   //模块权限标志
	Children []*GlobalAccess //子模块
	Parent   int             //父级模块ID
}

// ModuleItem 模块列表项，用于返回用户的权限列表
type ModuleItem struct {
	ID         int      //模块ID
	Title      string   //模块名称
	Router     string   //模块路径
	Access     []string //模块权限标志
	ModuleNeed []Access //模块出现所需要的AccessID
	Parent     int      //父级模块ID
}

type AccessItem struct {
	Key          string   `json:"key"`
	Title        string   `json:"title"`
	Dependencies []string `json:"dependencies"`
}

func initAccessConfig() {
	c := new(AccessConfig)
	err := yaml.Unmarshal(accessData, c)
	if err != nil {
		panic(err)
	}

	allAccess, _ = getAccessItems(c.Modules, 1, 0)
}

func GetGlobalAccessConfig() ([]*GlobalAccess, int) {
	return allAccess, maxAccessDepth
}

func GetAllModulesConfig() []*ModuleItem {
	return allModules
}

// getAccessItems 转化成前端所需要的权限列表
func getAccessItems(modules []*AccessConfigItem, depth, parent int) ([]*GlobalAccess, []Access) {
	items := make([]*GlobalAccess, len(modules))
	moduleNeedIds := make([]Access, 0)
	for i, module := range modules {
		item := &GlobalAccess{
			ID:     module.ID,
			Title:  module.Title,
			Module: module.Module,
			Parent: parent,
		}

		moduleItem := &ModuleItem{
			ID:         module.ID,
			Title:      module.Title,
			Router:     module.Module,
			ModuleNeed: nil,
			Parent:     parent,
		}

		if len(module.Access) > 0 {
			item.Access = make([]*AccessItem, 0, len(module.Access))
			accessIds := make([]Access, 0, len(module.Access))
			for _, key := range module.Access {
				access, has := accessParse[key]
				if !has {
					panic(fmt.Errorf("access %s:%w", key, ErrorAccessUnknown))
				}
				item.Access = append(item.Access, &AccessItem{
					Key:          key,
					Title:        access.Title(),
					Dependencies: dependenciesMap[key],
				})
				accessIds = append(accessIds, access)
			}

			moduleItem.Access = module.Access

			moduleItem.ModuleNeed = append(moduleItem.ModuleNeed, accessIds...)
			moduleNeedIds = append(moduleNeedIds, accessIds...)
		}
		items[i] = item

		allModules = append(allModules, moduleItem)
		if len(module.Children) > 0 {
			children, itemNeeds := getAccessItems(module.Children, depth+1, module.ID)
			moduleNeedIds = append(moduleNeedIds, moduleItem.ModuleNeed...)
			item.Children = children
			moduleItem.ModuleNeed = append(moduleItem.ModuleNeed, itemNeeds...)
		}
	}

	if depth > maxAccessDepth {
		maxAccessDepth = depth
	}
	return items, moduleNeedIds
}
