package navigation_service

import (
	"context"
	"sort"

	"github.com/eolinker/eosc"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type navigationService struct {
	navigations eosc.Untyped[string, *navigation_model.Navigation]
	index       int
}

func newNavigationService() *navigationService {
	return &navigationService{navigations: eosc.BuildUntyped[string, *navigation_model.Navigation]()}
}

func (n *navigationService) Add(ctx context.Context, uuid string, title string, icon string) error {
	_, has := n.navigations.Get(uuid)
	if has {
		return nil
	}
	n.navigations.Set(uuid, &navigation_model.Navigation{
		Uuid:  uuid,
		Title: title,
		Icon:  icon,
		Index: n.index,
	})
	n.index++
	return nil
}

func (n *navigationService) Info(ctx context.Context, uuid string) (*navigation_model.Navigation, bool) {
	return n.navigations.Get(uuid)
}

func (n *navigationService) List(ctx context.Context) ([]*navigation_model.Navigation, error) {
	navigations := n.navigations.List()
	sort.Sort(navigation_model.Navigations(navigations))
	return navigations, nil
}
