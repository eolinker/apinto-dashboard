package dynamic

import (
	"context"

	v2 "github.com/eolinker/apinto-dashboard/client/v2"

	dynamic_model "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-model"
)

type IDynamicService interface {
	Info(ctx context.Context, namespaceId int, profession string, name string) (*v2.WorkerInfo[dynamic_model.DynamicBasicInfo], error)
	List(ctx context.Context, namespaceId int, profession string, columns []string, keyword string, page int, pageSize int) ([]map[string]string, int, error)

	Online(ctx context.Context, namespaceId int, profession string, name string, names []string, updater int) ([]string, []string, error)
	Offline(ctx context.Context, namespaceId int, profession string, name string, names []string, updater int) ([]string, []string, error)

	ClusterStatuses(ctx context.Context, namespaceId int, profession string, names []string, keyword string, page int, pageSize int) (map[string]map[string]string, error)
	ClusterStatus(ctx context.Context, namespaceId int, profession string, name string) (*dynamic_model.DynamicBasicInfo, []*dynamic_model.DynamicCluster, error)

	Create(ctx context.Context, namespaceId int, module string, title string, name string, driver string, description string, body string, updater int) error
	Save(ctx context.Context, namespaceId int, module string, title string, name string, description string, body string, updater int) error
	Delete(ctx context.Context, namespaceId int, module string, name string) error
}
