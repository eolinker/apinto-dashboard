package strategy

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
)

type IStrategyService[T any, K any] interface {
	GetList(ctx context.Context, namespaceId int, clusterName string) ([]*strategy_model.Strategy, error)
	GetInfo(ctx context.Context, namespaceId int, uuid string) (*strategy_model.StrategyInfoOutput[K], error)
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
}
type IStrategyRemoteOptionHandle interface {
	Get(namespaceId int, keyword, groupUUID string, pageNum, pageSize int)
}
type IStrategyCommonService interface {
	GetFilterOptions(ctx context.Context, namespaceId int) ([]*strategy_model.FilterOptionsItem, error)
	GetFilterRemote(ctx context.Context, namespaceId int, targetType, keyword, groupUUID string, pageNum, pageSize int) (*strategy_model.FilterRemoteOutput, int, error)

	GetFilterLabel(ctx context.Context, namespaceId int, name string, value []string) (string, string, string)
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
