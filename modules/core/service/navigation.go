package service

import (
	"context"

	navigation_service "github.com/eolinker/apinto-dashboard/modules/navigation"

	"github.com/eolinker/eosc/common/bean"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"

	"github.com/eolinker/apinto-dashboard/modules/core/model"
)

type INavigationService interface {
	List(ctx context.Context) ([]*model.Navigation, map[string]string, error)
}

type navigation struct {
	navigationService   navigation_service.INavigationService
	modulePluginService module_plugin.IModulePlugin
}

func newNavigationService() INavigationService {
	n := &navigation{}
	bean.Autowired(&n.navigationService)
	bean.Autowired(&n.modulePluginService)
	return n
}

func (n *navigation) List(ctx context.Context) ([]*model.Navigation, map[string]string, error) {
	modules, err := n.modulePluginService.GetNavigationModules(ctx)
	if err != nil {
		return nil, nil, err
	}

	moduleMap := make(map[string][]*model.Module)
	access := make(map[string]string)
	moduleSort := make([]string, 0, len(modules))
	for _, m := range modules {
		if _, ok := moduleMap[m.Navigation]; !ok {
			moduleMap[m.Navigation] = make([]*model.Module, 0, 4)
			moduleSort = append(moduleSort, m.Navigation)
		}
		access[m.Name] = "edit"
		moduleMap[m.Navigation] = append(moduleMap[m.Navigation], &model.Module{
			Name:  m.Name,
			Title: m.Title,
			Type:  m.Type,
			Path:  m.Path,
		})
	}
	list := n.navigationService.List()

	navigations := make([]*model.Navigation, 0, len(list))
	for _, l := range list {
		ms := make([]*model.Module, 0)
		if v, ok := moduleMap[l.Uuid]; ok {
			ms = v
			delete(moduleMap, l.Uuid)
		}
		defaultModule := ""
		if len(ms) == 1 && ms[0].Name == l.Name {
			defaultModule = ms[0].Name
		}

		navigations = append(navigations, &model.Navigation{
			Title:   l.Title,
			Icon:    l.Icon,
			Modules: ms,
			Default: defaultModule,
		})
	}
	for _, m := range moduleSort {
		if vs, ok := moduleMap[m]; ok {
			for _, v := range vs {
				navigations = append(navigations, &model.Navigation{
					Title:   v.Title,
					Modules: []*model.Module{v},
					Default: v.Name,
				})
			}

		}
	}
	return navigations, access, nil
}
