package warn

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/warn/model"
	"time"
)

type IWarnStrategyService interface {
	CreateWarnStrategy(ctx context.Context, namespaceId, userId int, input *model.WarnStrategy) error
	UpdateWarnStrategy(ctx context.Context, namespaceId, userId int, input *model.WarnStrategy) error
	WarnStrategyListPage(ctx context.Context, namespaceId int, query *model.QueryWarnStrategyParam) ([]*model.WarnStrategy, int64, error)
	WarnStrategyAll(ctx context.Context, namespaceId, status int) ([]*model.WarnStrategy, error)
	WarnStrategyByUuid(ctx context.Context, namespaceId int, uuid string) (*model.WarnStrategy, error)
	UpdateWarnStrategyStatus(ctx context.Context, uuid string, isEnable bool) error
	DeleteWarnStrategy(ctx context.Context, uuid string) error
	DeleteWarnStrategyByPartitionId(ctx context.Context, namespaceId, partitionId int) error
}

type IWarnHistoryService interface {
	QueryList(ctx context.Context, namespaceId, partitionId, pageNum, pageSize int, startTime, endTime time.Time, name string) ([]*model.WarnHistoryInfo, int64, error)
	Create(ctx context.Context, namespaceId int, partitionId int, infos ...*model.WarnHistoryInfo) error
}
