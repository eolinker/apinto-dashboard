package navigation

import (
	"context"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type INavigationService interface {
	Add(ctx context.Context, uuid string, name string, icon string) error
	Save(ctx context.Context, uuid string, name string, icon string) error
	Delete(ctx context.Context, uuid string) error
	List(ctx context.Context) ([]*navigation_model.NavigationBasicInfo, error)
	Info(ctx context.Context, uuid string) (*navigation_model.Navigation, error)
	GetUUIDByID(ctx context.Context, id int) (string, error)
	// Sort 更新排序
	Sort(ctx context.Context, uuids []string) error
}
