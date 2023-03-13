package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/model"
)

type IStrategyHandler[T any, K any] interface {
	GetListLabel(conf *T) string
	GetType() string
	GetConfName() string
	GetBatchSettingName() string
	ToApintoConfig(conf T) interface{}
	FormatOut(ctx context.Context, namespaceID int, input *model.StrategyInfoOutput[T]) *model.StrategyInfoOutput[K]
	CheckInput(input *dto.StrategyInfoInput[T]) error
}
type FormatHandler[T any] struct {
}

func (f *FormatHandler[T]) FormatOut(ctx context.Context, namespaceID int, output *model.StrategyInfoOutput[T]) *model.StrategyInfoOutput[T] {
	return output
}
