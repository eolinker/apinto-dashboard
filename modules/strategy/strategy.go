package strategy

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
)

type IStrategyService[T any, K any] interface {
	GetList(ctx context.Context, namespaceId int, clusterName string) ([]*strategy_model.Strategy, error)
	GetInfo(ctx context.Context, namespaceId int, uuid string) (*strategy_model.StrategyInfoOutput[K], *strategy_model.ExtenderData, error)
	CreateStrategy(ctx context.Context, namespaceId int, operator int, clusterName string, input *strategy_dto.StrategyInfoInput[T]) error
	UpdateStrategy(ctx context.Context, namespaceId int, operator int, clusterName string, input *strategy_dto.StrategyInfoInput[T]) error
	DeleteStrategy(ctx context.Context, namespaceId, operator int, clusterName, uuid string) error
	RestoreStrategy(ctx context.Context, namespaceId, userID int, clusterName, uuid string) error
	UpdateStop(ctx context.Context, namespaceId, operator int, uuid, clusterName string, stop bool) error
	ToPublish(ctx context.Context, namespaceId int, clusterName string) ([]*strategy_model.StrategyToPublish[T], error)
	PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*strategy_model.StrategyPublishHistory, int, error)
	Publish(ctx context.Context, namespaceId, operator int, clusterName string, input *strategy_dto.StrategyPublish) error

	ChangePriority(ctx context.Context, namespaceId, userId int, clusterName string, maps map[string]int) error
	CheckInput(input *strategy_dto.StrategyInfoInput[T]) error

	//checkPriorityReduplicative(ctx context.Context, clusterID, priority int, strategyType, uuid string) (int, error)
	//getLatestStrategyVersion(ctx context.Context, strategyID int) (*strategy_entry.StrategyVersion, error)
	//toApinto(name, desc string, isStop bool, priority int, filters []strategy_entry.StrategyFiltersConfig, conf T) map[string]interface{}
	//encodeConfig(config *T) string
	//decodeConfig(config string) *T
	//online.IResetOnlineService
}

type IStrategyCommonService interface {
	GetFilterOptions(ctx context.Context, namespaceId int) ([]*strategy_model.FilterOptionsItem, error)
	GetFilterRemote(ctx context.Context, namespaceId int, targetType, keyword, groupUUID string, pageNum, pageSize int) (*strategy_model.FilterRemoteOutput, int, error)
	GetMetricsOptions() ([]*strategy_model.MetricsOptionsItem, error)
	//AddHandler(onlineService online.IResetOnlineService)
	//online.IResetOnlineService
}

type IStrategyHandler[T any, K any] interface {
	GetListLabel(conf *T) string
	GetType() string
	GetConfName() string
	GetBatchSettingName() string
	ToApintoConfig(conf T) interface{}
	FormatOut(ctx context.Context, namespaceID int, input *strategy_model.StrategyInfoOutput[T]) *strategy_model.StrategyInfoOutput[K]
	CheckInput(input *strategy_dto.StrategyInfoInput[T]) error
}

type FormatHandler[T any] struct {
}

func (f *FormatHandler[T]) FormatOut(ctx context.Context, namespaceID int, output *strategy_model.StrategyInfoOutput[T]) *strategy_model.StrategyInfoOutput[T] {
	return output
}
