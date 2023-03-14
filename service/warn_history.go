package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/monitor-entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"time"
)

type IWarnHistoryService interface {
	QueryList(ctx context.Context, namespaceId, partitionId, pageNum, pageSize int, startTime, endTime time.Time, name string) ([]*model.WarnHistoryInfo, int64, error)
	Create(ctx context.Context, namespaceId int, partitionId int, infos ...*model.WarnHistoryInfo) error
}

type warnHistoryService struct {
	warnHistoryStore store.IWarnHistoryIStore
}

func newWarnHistoryService() IWarnHistoryService {
	w := new(warnHistoryService)
	bean.Autowired(&w.warnHistoryStore)
	return w
}

func (w *warnHistoryService) Create(ctx context.Context, namespaceId int, partitionId int, infos ...*model.WarnHistoryInfo) error {
	if len(infos) == 0 {
		return nil
	}
	list := make([]*monitor_entry.WarnHistory, 0, len(infos))
	for _, info := range infos {
		history := &monitor_entry.WarnHistory{
			NamespaceID:   namespaceId,
			PartitionId:   partitionId,
			StrategyTitle: info.StrategyTitle,
			ErrMsg:        info.ErrMsg,
			Status:        info.Status,
			Dimension:     info.Dimension,
			Quota:         info.Quota,
			Target:        info.Target,
			Content:       info.Content,
			CreateTime:    info.CreateTime,
		}
		list = append(list, history)
	}

	return w.warnHistoryStore.Insert(ctx, list...)
}

func (w *warnHistoryService) QueryList(ctx context.Context, _, partitionId, pageNum, pageSize int, startTime, endTime time.Time, name string) ([]*model.WarnHistoryInfo, int64, error) {
	list, count, err := w.warnHistoryStore.GetPage(ctx, partitionId, name, pageNum, pageSize, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	warnHistory := make([]*model.WarnHistoryInfo, 0, len(list))
	for _, history := range list {
		warnHistory = append(warnHistory, &model.WarnHistoryInfo{
			StrategyTitle: history.StrategyTitle,
			Dimension:     history.Dimension,
			Quota:         history.Quota,
			Target:        history.Target,
			Content:       history.Content,
			ErrMsg:        history.ErrMsg,
			Status:        history.Status,
			CreateTime:    history.CreateTime,
		})
	}

	return warnHistory, count, nil
}
