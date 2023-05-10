package dynamic

import (
	"context"

	v2 "github.com/eolinker/apinto-dashboard/client/v2"

	dynamic_model "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-model"
)

type IDynamicService interface {
	Info(ctx context.Context, namespaceId int, profession string, name string) (*v2.WorkerInfo[dynamic_model.DynamicBasicInfo], error)
	List(ctx context.Context, namespaceId int, profession string, columns []string, drivers []string, keyword string, page int, pageSize int) ([]map[string]string, int, error)
	GetBySkill(ctx context.Context, namespaceId int, skill string) ([]*dynamic_model.DynamicBasicInfo, error)

	Online(ctx context.Context, namespaceId int, profession string, module string, name string, names []string, updater int, depends ...string) ([]string, []string, error)
	Offline(ctx context.Context, namespaceId int, profession string, module string, name string, names []string, updater int) ([]string, []string, error)

	ClusterStatuses(ctx context.Context, namespaceId int, profession string, names []string, drivers []string, keyword string, page int, pageSize int) (map[string]map[string]string, error)
	ClusterStatus(ctx context.Context, namespaceId int, profession string, name string) (*dynamic_model.DynamicBasicInfo, []*dynamic_model.DynamicCluster, error)
	ClusterStatusByClusterName(ctx context.Context, namespaceId int, profession, name string, clusterName string) (*dynamic_model.DynamicCluster, error)

	Create(ctx context.Context, namespaceId int, profession string, module string, skill string, title string, name string, driver string, description string, body string, updater int, depend ...string) error
	Save(ctx context.Context, namespaceId int, profession string, module string, title string, name string, description string, body string, updater int, depend ...string) error
	Delete(ctx context.Context, namespaceId int, profession string, module string, name string) error

	GetIDByName(ctx context.Context, namespaceId int, profession string, name string) (int, error)

	ListByNames(ctx context.Context, namespaceID int, profession string, names []string) ([]*dynamic_model.DynamicBasicInfo, error)
	ListByKeyword(ctx context.Context, namespaceID int, profession string, keyword string) ([]*dynamic_model.DynamicBasicInfo, error)

	Count(ctx context.Context, namespaceID int, profession string, addition map[string]interface{}) (int, error)
}
