package warn

import (
	"context"
	warn_model "github.com/eolinker/apinto-dashboard/modules/warn/warn-model"
	"time"
)

type IWarnStrategyService interface {
	CreateWarnStrategy(ctx context.Context, namespaceId, userId int, input *warn_model.WarnStrategy) error
	UpdateWarnStrategy(ctx context.Context, namespaceId, userId int, input *warn_model.WarnStrategy) error
	WarnStrategyListPage(ctx context.Context, namespaceId int, query *warn_model.QueryWarnStrategyParam) ([]*warn_model.WarnStrategy, int64, error)
	WarnStrategyAll(ctx context.Context, namespaceId, status int) ([]*warn_model.WarnStrategy, error)
	WarnStrategyByUuid(ctx context.Context, namespaceId int, uuid string) (*warn_model.WarnStrategy, error)
	UpdateWarnStrategyStatus(ctx context.Context, uuid string, isEnable bool) error
	DeleteWarnStrategy(ctx context.Context, uuid string) error
	DeleteWarnStrategyByPartitionId(ctx context.Context, namespaceId, partitionId int) error
}

type IWarnHistoryService interface {
	QueryList(ctx context.Context, namespaceId, partitionId, pageNum, pageSize int, startTime, endTime time.Time, name string) ([]*warn_model.WarnHistoryInfo, int64, error)
	Create(ctx context.Context, namespaceId int, partitionId int, infos ...*warn_model.WarnHistoryInfo) error
}
