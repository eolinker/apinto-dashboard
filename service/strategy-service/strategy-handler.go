package strategy_service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/dto/strategy-dto"
	"github.com/eolinker/apinto-dashboard/model/openapi-model"
)

type IStrategyHandler[T any, K any] interface {
	GetListLabel(conf *T) string
	GetType() string
	GetConfName() string
	GetBatchSettingName() string
	ToApintoConfig(conf T) interface{}
	FormatOut(ctx context.Context, namespaceID int, input *openapi_model.StrategyInfoOutput[T]) *openapi_model.StrategyInfoOutput[K]
	CheckInput(input *strategy_dto.StrategyInfoInput[T]) error
}
type FormatHandler[T any] struct {
}

func (f *FormatHandler[T]) FormatOut(ctx context.Context, namespaceID int, output *openapi_model.StrategyInfoOutput[T]) *openapi_model.StrategyInfoOutput[T] {
	return output
}
