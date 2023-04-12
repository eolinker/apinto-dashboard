package service

import (
	"context"

	"github.com/eolinker/eosc/common/bean"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"

	"github.com/eolinker/apinto-dashboard/modules/core/model"
)

type INavigationService interface {
	List(ctx context.Context) ([]*model.Navigation, map[string]string, error)
}

type navigation struct {
	//navigationService   navigation_service.INavigationService
	modulePluginService module_plugin.IModulePlugin
}

func newNavigationService() INavigationService {
	n := &navigation{}
	//bean.Autowired(&n.navigationService)
	bean.Autowired(&n.modulePluginService)
	return n
}

func (n *navigation) List(ctx context.Context) ([]*model.Navigation, map[string]string, error) {
	//return nil, nil, nil
	list := n.navigationService.List()

	//ids := make([]int, 0, len(list))
	//for _, l := range list {
	//	ids = append(ids, l.ID)
	//}
	//moduleMap, err := n.modulePluginService.GetModulesByNavigations(ctx, ids)
	//if err != nil {
	//	return nil, nil, err
	//}
	//access := make(map[string]string)
	//navigations := make([]*model.Navigation, 0, len(list))
	//for _, l := range list {
	//	modules := make([]*model.Module, 0)
	//	defaultModule := ""
	//	if vs, ok := moduleMap[l.ID]; ok {
	//
	//		for i, v := range vs {
	//			if i == 0 && v.Title == l.Title {
	//				defaultModule = v.Name
	//			}
	//			access[v.Name] = "edit"
	//			modules = append(modules, &model.Module{
	//				Name:  v.Name,
	//				Title: v.Title,
	//				Type:  v.Type,
	//				Path:  v.Path,
	//			})
	//		}
	//		if len(vs) > 1 {
	//			defaultModule = ""
	//		}
	//	}
	//
	//	navigations = append(navigations, &model.Navigation{
	//		Title:    l.Title,
	//		Icon:     l.Icon,
	//		IconType: "css",
	//		Modules:  modules,
	//		Default:  defaultModule,
	//	})
	//}
	//return navigations, access, nil
}
