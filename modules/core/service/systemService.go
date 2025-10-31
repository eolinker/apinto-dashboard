package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/pm3"

	navigation_service "github.com/eolinker/apinto-dashboard/modules/navigation"

	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/apinto-dashboard/modules/core/model"
)

var _ core.ISystemService = (*navigation)(nil)

type navigation struct {
	navigationService   navigation_service.INavigationService
	modulePluginService mpm3.IModuleService
	frontendService     mpm3.IFrontendService
}

func (n *navigation) PluginConfig(ctx context.Context) ([]pm3.PFrontend, error) {
	return n.frontendService.GetEnable(ctx), nil

}

func newNavigationService() core.ISystemService {
	n := &navigation{}
	bean.Autowired(&n.navigationService)
	bean.Autowired(&n.frontendService)
	bean.Autowired(&n.modulePluginService)
	return n
}

func (n *navigation) Navigations(ctx context.Context) ([]*model.Navigation, error) {
	modules := n.modulePluginService.GetEnable(ctx)

	moduleMap := make(map[string][]*model.Module)
	moduleSort := make([]string, 0, len(modules))
	for _, m := range modules {
		if _, ok := moduleMap[m.Navigation]; !ok {
			moduleMap[m.Navigation] = make([]*model.Module, 0, 4)
			moduleSort = append(moduleSort, m.Navigation)
		}
		//access[m.Name] = "edit"
		moduleMap[m.Navigation] = append(moduleMap[m.Navigation], &model.Module{
			Name:  m.Name,
			Title: m.CName,
			Path:  m.Router,
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
		if len(ms) == 1 && l.Quick {
			defaultModule = ms[0].Name
		}

		//若导航下没有模块, 则不显示
		if len(ms) == 0 {
			continue
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
	return navigations, nil
}
