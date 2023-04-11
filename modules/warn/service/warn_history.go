package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/warn"
	"github.com/eolinker/apinto-dashboard/modules/warn/store"
	"github.com/eolinker/apinto-dashboard/modules/warn/warn-entry"
	"github.com/eolinker/apinto-dashboard/modules/warn/warn-model"
	"github.com/eolinker/eosc/common/bean"
	"time"
)

type warnHistoryService struct {
	warnHistoryStore store.IWarnHistoryIStore
}

func newWarnHistoryService() warn.IWarnHistoryService {
	w := new(warnHistoryService)
	bean.Autowired(&w.warnHistoryStore)
	return w
}

func (w *warnHistoryService) Create(ctx context.Context, namespaceId int, partitionId int, infos ...*warn_model.WarnHistoryInfo) error {
	if len(infos) == 0 {
		return nil
	}
	list := make([]*warn_entry.WarnHistory, 0, len(infos))
	for _, info := range infos {
		history := &warn_entry.WarnHistory{
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

func (w *warnHistoryService) QueryList(ctx context.Context, _, partitionId, pageNum, pageSize int, startTime, endTime time.Time, name string) ([]*warn_model.WarnHistoryInfo, int64, error) {
	list, count, err := w.warnHistoryStore.GetPage(ctx, partitionId, name, pageNum, pageSize, startTime, endTime)
	if err != nil {
		return nil, 0, err
	}

	warnHistory := make([]*warn_model.WarnHistoryInfo, 0, len(list))
	for _, history := range list {
		warnHistory = append(warnHistory, &warn_model.WarnHistoryInfo{
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
