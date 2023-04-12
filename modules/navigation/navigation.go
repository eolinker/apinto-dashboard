package navigation

import (
	"context"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type INavigationService interface {
	Add(ctx context.Context, uuid string, title string, icon string) error
	Info(ctx context.Context, uuid string) (*navigation_model.Navigation, error)
	List(ctx context.Context) ([]*navigation_model.Navigation, error)
}
