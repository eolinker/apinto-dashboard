package navigation_service

import (
	"context"

	"github.com/eolinker/apinto-dashboard/modules/navigation"
	"github.com/eolinker/eosc/common/bean"

	"github.com/eolinker/eosc"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type navigationService struct {
	dataService   navigation.INavigationDataService
	navigationMap eosc.Untyped[string, *navigation_model.Navigation]
	navigations   []*navigation_model.Navigation
}

func newNavigationService() *navigationService {
	n := &navigationService{}
	bean.Autowired(&n.navigationMap)
	return n
}

func (n *navigationService) initData() {
	navigations := n.dataService.GetNavigationData()
	ns := make([]*navigation_model.Navigation, 0, len(navigations))
	nsMap := eosc.BuildUntyped[string, *navigation_model.Navigation]()
	for _, nv := range navigations {
		_, ok := nsMap.Get(nv.Uuid)
		if !ok {
			v := &navigation_model.Navigation{
				Uuid:  nv.Uuid,
				Title: nv.Title,
				Icon:  nv.Icon,
			}
			nsMap.Set(nv.Uuid, v)
			ns = append(ns, v)
		}
	}
	n.navigations = ns
	n.navigationMap = nsMap
}

func (n *navigationService) Info(ctx context.Context, uuid string) (*navigation_model.Navigation, bool) {
	if n.navigationMap == nil {
		n.initData()
	}
	return n.navigationMap.Get(uuid)
}

func (n *navigationService) List(ctx context.Context) []*navigation_model.Navigation {
	if n.navigations == nil {
		n.initData()
	}
	return n.navigations
}
