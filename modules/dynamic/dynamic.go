package dynamic

import (
	"context"

	dynamic_model "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-model"
)

type IDynamicService interface {
	Info(ctx context.Context, namespaceId int, module string, name string) (*dynamic_model.DynamicBasicInfo, map[string]interface{}, error)
	List(ctx context.Context, namespaceId int, module string, columns []string, keyword string, page int, pageSize int) ([]map[string]string, error)

	Online(ctx context.Context, namespaceId int, module string, profession string, name string, names []string, updater int) ([]string, []string, error)
	Offline(ctx context.Context, namespaceId int, module string, profession string, name string, names []string, updater int) ([]string, []string, error)

	ClusterStatuses(ctx context.Context, namespaceId int, module string, profession string, names []string, keyword string, page int, pageSize int) (map[string]map[string]string, error)
	ClusterStatus(ctx context.Context, namespaceId int, module string, profession string, name string) (*dynamic_model.DynamicBasicInfo, []*dynamic_model.DynamicCluster, error)

	Create(ctx context.Context, namespaceId int, module string, title string, name string, driver string, description string, body string, updater int) error
	Save(ctx context.Context, namespaceId int, module string, title string, name string, description string, body string, updater int) error
	Delete(ctx context.Context, namespaceId int, module string, name string) error
}
