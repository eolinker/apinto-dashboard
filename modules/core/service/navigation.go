package service

import (
	"context"

	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"

	"github.com/eolinker/apinto-dashboard/modules/navigation"

	"github.com/eolinker/apinto-dashboard/modules/core/model"
)

type INavigationService interface {
	List(ctx context.Context) ([]*model.Navigation, error)
}

type navigationService struct {
	navigationService   navigation.INavigationService
	modulePluginService module_plugin.IModulePlugin
}

func (n *navigationService) List(ctx context.Context) ([]*model.Navigation, error) {
	list, err := n.navigationService.List(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]int, 0, len(list))
	for _, l := range list {
		ids = append(ids, l.ID)
	}
	moduleMap, err := n.modulePluginService.GetModulesByNavigations(ctx, ids)
	if err != nil {
		return nil, err
	}
	navigations := make([]*model.Navigation, 0, len(list))
	for _, l := range list {
		modules := make([]*model.Module, 0)
		defaultModule := ""
		if vs, ok := moduleMap[l.ID]; ok {
			for i, v := range vs {
				if i == 0 {
					defaultModule = v.Name
				}
				modules = append(modules, &model.Module{
					Name:  v.Name,
					Title: v.Title,
					Type:  v.Type,
					Path:  v.Path,
				})
			}
		}

		navigations = append(navigations, &model.Navigation{
			Title:    l.Title,
			Icon:     l.Icon,
			IconType: l.IconType,
			Modules:  modules,
			Default:  defaultModule,
		})
	}
	return navigations, nil
}
