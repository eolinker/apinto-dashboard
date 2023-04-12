package navigation

import (
	"context"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type INavigationService interface {
	Info(ctx context.Context, uuid string) (*navigation_model.Navigation, bool)
	List(ctx context.Context) ([]*navigation_model.Navigation, error)
}

type INavigationDataService interface {
	GetNavigationData() []*navigation_model.Navigation
}
