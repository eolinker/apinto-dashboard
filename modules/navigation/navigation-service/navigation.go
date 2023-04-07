package navigation_service

import (
	"context"
	"sort"

	"github.com/eolinker/eosc/common/bean"

	navigation_entry "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-entry"

	"gorm.io/gorm"

	navigation_store "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-store"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type navigationService struct {
	navigationStore navigation_store.INavigationStore
}

func (n *navigationService) GetUUIDByID(ctx context.Context, id int) (string, error) {
	v, err := n.navigationStore.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return v.Uuid, nil
}

func (n *navigationService) GetIDByUUID(ctx context.Context, uuid string) (int, error) {
	v, err := n.navigationStore.First(ctx, map[string]interface{}{
		"uuid": uuid,
	})
	if err != nil {
		return 0, err
	}
	return v.Id, nil
}

func newNavigationService() *navigationService {
	c := &navigationService{}
	bean.Autowired(&c.navigationStore)
	return c
}

func (n *navigationService) Add(ctx context.Context, uuid string, name string, icon string) error {
	index, err := n.navigationStore.MaxSort(ctx)
	if err != nil {
		return err
	}
	return n.navigationStore.Insert(ctx, &navigation_entry.Navigation{
		Uuid:  uuid,
		Title: name,
		Icon:  icon,
		Sort:  index + 1,
	})
}

func (n *navigationService) Save(ctx context.Context, uuid string, name string, icon string) error {

	return n.navigationStore.Save(ctx, &navigation_entry.Navigation{
		Uuid:  uuid,
		Title: name,
		Icon:  icon,
	})
}

func (n *navigationService) Delete(ctx context.Context, uuid string) error {
	_, err := n.navigationStore.DeleteWhere(ctx, map[string]interface{}{
		"uuid": uuid,
	})
	return err
}

func (n *navigationService) List(ctx context.Context) ([]*navigation_model.NavigationBasicInfo, error) {
	list, err := n.navigationStore.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	navigations := make([]*navigation_model.NavigationBasicInfo, 0)
	for _, l := range list {
		_, err = n.navigationStore.First(ctx, map[string]interface{}{
			"navigation_uuid": l.Uuid,
		})
		canDel := false
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
			canDel = true
		}
		navigations = append(navigations, &navigation_model.NavigationBasicInfo{
			Uuid:      l.Uuid,
			Title:     l.Title,
			Icon:      l.Icon,
			CanDelete: canDel,
			Sort:      l.Sort,
		})
	}
	sort.Sort(navigation_model.Navigations(navigations))
	return navigations, nil
}

func (n *navigationService) Info(ctx context.Context, uuid string) (*navigation_model.Navigation, error) {
	info, err := n.navigationStore.First(ctx, map[string]interface{}{
		"uuid": uuid,
	})
	if err != nil {
		return nil, err
	}

	pluginService.Get(ctx)
	//connModules, err := n.navigationModuleStore.List(ctx, map[string]interface{}{
	//	"navigation_uuid": info.Uuid,
	//})
	//
	//modules := make([]*navigation_model.Module, 0, len(connModules))
	//for _, m := range connModules {
	//	// TODO:获取模块标题信息
	//	modules = append(modules, &navigation_model.Module{
	//		Id:    m.ModuleId,
	//		Title: m.ModuleId,
	//	})
	//}
	//
	return nil, nil
}

func (n *navigationService) Sort(ctx context.Context, uuids []string) error {
	return n.navigationStore.Transaction(ctx, func(txCtx context.Context) error {
		return n.navigationStore.SortByUUIDs(txCtx, uuids)
	})
}
