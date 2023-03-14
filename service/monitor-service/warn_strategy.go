package monitor_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/entry/monitor-entry"
	"github.com/eolinker/apinto-dashboard/entry/quote-entry"
	"github.com/eolinker/apinto-dashboard/model/audit-model"
	"github.com/eolinker/apinto-dashboard/model/monitor-model"
	"github.com/eolinker/apinto-dashboard/model/notice-model"
	"github.com/eolinker/apinto-dashboard/service/notice-service"
	"github.com/eolinker/apinto-dashboard/service/user-service"
	"github.com/eolinker/apinto-dashboard/store/monitor-store"
	"github.com/eolinker/apinto-dashboard/store/quote-store"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
	"time"
)

type IWarnStrategyService interface {
	CreateWarnStrategy(ctx context.Context, namespaceId, userId int, input *monitor_model.WarnStrategy) error
	UpdateWarnStrategy(ctx context.Context, namespaceId, userId int, input *monitor_model.WarnStrategy) error
	WarnStrategyListPage(ctx context.Context, namespaceId int, query *monitor_model.QueryWarnStrategyParam) ([]*monitor_model.WarnStrategy, int64, error)
	WarnStrategyAll(ctx context.Context, namespaceId, status int) ([]*monitor_model.WarnStrategy, error)
	WarnStrategyByUuid(ctx context.Context, namespaceId int, uuid string) (*monitor_model.WarnStrategy, error)
	UpdateWarnStrategyStatus(ctx context.Context, uuid string, isEnable bool) error
	DeleteWarnStrategy(ctx context.Context, uuid string) error
	DeleteWarnStrategyByPartitionId(ctx context.Context, namespaceId, partitionId int) error
}

type warnStrategyService struct {
	warnStrategyStore    monitor_store.IWarnStrategyIStore
	noticeChannelService notice_service.INoticeChannelService
	monitorService       IMonitorService
	userService          user_service.IUserInfoService
	quoteStore           quote_store.IQuoteStore
}

func newWarnStrategyService() IWarnStrategyService {
	w := new(warnStrategyService)
	bean.Autowired(&w.warnStrategyStore)
	bean.Autowired(&w.noticeChannelService)
	bean.Autowired(&w.quoteStore)
	bean.Autowired(&w.userService)
	bean.Autowired(&w.monitorService)
	return w
}

func (w *warnStrategyService) UpdateWarnStrategyStatus(ctx context.Context, uuid string, isEnable bool) error {
	warnStrategy, err := w.warnStrategyStore.GetByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	return w.warnStrategyStore.UpdateIsEnable(ctx, warnStrategy.Id, isEnable)
}

func (w *warnStrategyService) DeleteWarnStrategy(ctx context.Context, uuid string) error {
	warnStrategy, err := w.warnStrategyStore.GetByUuid(ctx, uuid)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuid,
		Name: warnStrategy.Title,
	})

	return w.warnStrategyStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = w.warnStrategyStore.Delete(txCtx, warnStrategy.Id); err != nil {
			return err
		}
		//删除被绑定的渠道信息
		return w.quoteStore.DelBySource(txCtx, warnStrategy.Id, quote_entry.QuoteKindTypeWarnStrategy)
	})
}

func (w *warnStrategyService) DeleteWarnStrategyByPartitionId(ctx context.Context, namespaceId int, partitionId int) error {

	warnStrategy, err := w.warnStrategyStore.GetByPartitionId(ctx, namespaceId, partitionId)
	if err != nil {
		return err
	}

	for _, strategy := range warnStrategy {
		if _, err = w.warnStrategyStore.Delete(ctx, strategy.Id); err != nil {
			return err
		}
		//删除被绑定的渠道信息
		if err = w.quoteStore.DelBySource(ctx, strategy.Id, quote_entry.QuoteKindTypeWarnStrategy); err != nil {
			return err
		}
	}

	return nil
}

// WarnStrategyAll
// status -1 查所有状态的 0只查禁用 1只查启用状态
func (w *warnStrategyService) WarnStrategyAll(ctx context.Context, namespaceId, status int) ([]*monitor_model.WarnStrategy, error) {
	strategies, err := w.warnStrategyStore.GetAll(ctx, namespaceId, status)
	if err != nil {
		return nil, err
	}

	result := make([]*monitor_model.WarnStrategy, 0, len(strategies))
	for _, warnStrategy := range strategies {

		warnStrategyConfig := new(monitor_model.WarnStrategyConfig)

		if err = json.Unmarshal([]byte(warnStrategy.Config), warnStrategyConfig); err != nil {
			return nil, err
		}

		item, err := w.monitorService.PartitionById(ctx, warnStrategy.PartitionId)
		if err != nil {
			return nil, err
		}

		result = append(result, &monitor_model.WarnStrategy{
			PartitionId:        warnStrategy.PartitionId,
			NamespaceId:        warnStrategy.NamespaceID,
			Uuid:               warnStrategy.UUID,
			Title:              warnStrategy.Title,
			Desc:               warnStrategy.Desc,
			IsEnable:           warnStrategy.IsEnable,
			Dimension:          warnStrategy.Dimension,
			PartitionUUID:      item.Id,
			Quota:              monitor_model.QuotaType(warnStrategy.Quota),
			Every:              warnStrategy.Every,
			WarnStrategyConfig: warnStrategyConfig,
			CreateTime:         warnStrategy.CreateTime,
			UpdateTime:         warnStrategy.UpdateTime,
		})
	}

	return result, nil
}
func (w *warnStrategyService) WarnStrategyListPage(ctx context.Context, namespaceId int, query *monitor_model.QueryWarnStrategyParam) ([]*monitor_model.WarnStrategy, int64, error) {
	list, count, err := w.warnStrategyStore.GetPage(ctx, namespaceId, query.PartitionId, query.StrategyName, query.Dimension, query.Status, query.PageNum, query.PageSize)
	if err != nil {
		return nil, 0, err
	}

	userIds := common.SliceToSliceIds(list, func(t *monitor_entry.WarnStrategy) int {
		return t.Operator
	})

	userInfoMaps, err := w.userService.GetUserInfoMaps(ctx, userIds...)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*monitor_model.WarnStrategy, 0, len(list))
	for _, warnStrategy := range list {
		warnStrategyConfig := new(monitor_model.WarnStrategyConfig)

		if err = json.Unmarshal([]byte(warnStrategy.Config), warnStrategyConfig); err != nil {
			return nil, 0, err
		}

		operator := ""
		if userInfo, ok := userInfoMaps[warnStrategy.Operator]; ok {
			operator = userInfo.NickName
		}

		result = append(result, &monitor_model.WarnStrategy{
			Uuid:               warnStrategy.UUID,
			Title:              warnStrategy.Title,
			Desc:               warnStrategy.Desc,
			IsEnable:           warnStrategy.IsEnable,
			Dimension:          warnStrategy.Dimension,
			Quota:              monitor_model.QuotaType(warnStrategy.Quota),
			Every:              warnStrategy.Every,
			WarnStrategyConfig: warnStrategyConfig,
			Operator:           operator,
			CreateTime:         warnStrategy.CreateTime,
			UpdateTime:         warnStrategy.UpdateTime,
		})
	}

	return result, count, nil
}

func (w *warnStrategyService) WarnStrategyByUuid(ctx context.Context, _ int, uuid string) (*monitor_model.WarnStrategy, error) {
	warnStrategy, err := w.warnStrategyStore.GetByUuid(ctx, uuid)
	if err != nil {
		return nil, err
	}

	warnStrategyConfig := new(monitor_model.WarnStrategyConfig)

	if err = json.Unmarshal([]byte(warnStrategy.Config), warnStrategyConfig); err != nil {
		return nil, err
	}

	result := &monitor_model.WarnStrategy{
		Uuid:               warnStrategy.UUID,
		Title:              warnStrategy.Title,
		Desc:               warnStrategy.Desc,
		IsEnable:           warnStrategy.IsEnable,
		Dimension:          warnStrategy.Dimension,
		Quota:              monitor_model.QuotaType(warnStrategy.Quota),
		Every:              warnStrategy.Every,
		WarnStrategyConfig: warnStrategyConfig,
	}

	return result, nil
}

func (w *warnStrategyService) CreateWarnStrategy(ctx context.Context, namespaceId, userId int, input *monitor_model.WarnStrategy) error {

	warnStrategy, err := w.warnStrategyStore.GetByTitle(ctx, namespaceId, input.PartitionId, input.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if warnStrategy != nil {
		return errors.New("策略名称重复")
	}

	t := time.Now()
	config, _ := json.Marshal(input.WarnStrategyConfig)

	strategy := &monitor_entry.WarnStrategy{
		NamespaceID: namespaceId,
		PartitionId: input.PartitionId,
		Title:       input.Title,
		UUID:        input.Uuid,
		Desc:        input.Desc,
		IsEnable:    input.IsEnable,
		Dimension:   input.Dimension,
		Quota:       string(input.Quota),
		Every:       input.Every,
		Config:      string(config),
		Operator:    userId,
		CreateTime:  t,
		UpdateTime:  t,
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.Uuid,
		Name: input.Title,
	})

	channels := make([]*notice_model.NoticeChannel, 0)
	//关联策略绑定的通知渠道信息
	for _, rule := range input.WarnStrategyConfig.Rule {
		for _, uuid := range rule.ChannelUuids {
			if uuid != "" {
				channel, err := w.noticeChannelService.NoticeChannelByName(ctx, namespaceId, uuid)
				if err != nil {
					return errors.New(fmt.Sprintf("%s 通知渠道不存在", uuid))
				}
				channels = append(channels, channel)
			}
		}
	}

	return w.warnStrategyStore.Transaction(ctx, func(txCtx context.Context) error {

		if err := w.warnStrategyStore.Insert(txCtx, strategy); err != nil {
			return err
		}

		if len(channels) > 0 {
			//往引用表插入所引用的通知渠道
			quoteMap := make(map[quote_entry.QuoteTargetKindType][]int)
			for _, channel := range channels {
				quoteMap[quote_entry.QuoteTargetKindTypeNoticeChannel] = append(quoteMap[quote_entry.QuoteTargetKindTypeNoticeChannel], channel.Id)
			}

			if err := w.quoteStore.Set(txCtx, strategy.Id, quote_entry.QuoteKindTypeWarnStrategy, quoteMap); err != nil {
				return err
			}
		}

		return nil
	})

}

func (w *warnStrategyService) UpdateWarnStrategy(ctx context.Context, namespaceId, userId int, input *monitor_model.WarnStrategy) error {

	warnStrategy, err := w.warnStrategyStore.GetByTitle(ctx, namespaceId, input.PartitionId, input.Title)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if warnStrategy != nil && warnStrategy.UUID != input.Uuid {
		return errors.New("策略名称重复")
	}

	warnStrategy, err = w.warnStrategyStore.GetByUuid(ctx, input.Uuid)
	if err != nil {
		return err
	}

	t := time.Now()
	config, _ := json.Marshal(input.WarnStrategyConfig)

	strategy := &monitor_entry.WarnStrategy{
		Id:          warnStrategy.Id,
		NamespaceID: namespaceId,
		PartitionId: warnStrategy.PartitionId,
		Title:       input.Title,
		UUID:        input.Uuid,
		Desc:        input.Desc,
		IsEnable:    input.IsEnable,
		Dimension:   input.Dimension,
		Quota:       string(input.Quota),
		Every:       input.Every,
		Config:      string(config),
		Operator:    userId,
		UpdateTime:  t,
	}

	channels := make([]*notice_model.NoticeChannel, 0)
	//关联策略绑定的通知渠道信息
	for _, rule := range input.WarnStrategyConfig.Rule {
		for _, uuid := range rule.ChannelUuids {
			if uuid != "" {
				channel, err := w.noticeChannelService.NoticeChannelByName(ctx, namespaceId, uuid)
				if err != nil {
					return errors.New(fmt.Sprintf("%s 通知渠道不存在", uuid))
				}
				channels = append(channels, channel)
			}
		}
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.Uuid,
		Name: input.Title,
	})

	return w.warnStrategyStore.Transaction(ctx, func(txCtx context.Context) error {

		if err := w.warnStrategyStore.Save(txCtx, strategy); err != nil {
			return err
		}

		if len(channels) > 0 {
			//往引用表插入所引用的通知渠道
			quoteMap := make(map[quote_entry.QuoteTargetKindType][]int)
			for _, channel := range channels {
				quoteMap[quote_entry.QuoteTargetKindTypeNoticeChannel] = append(quoteMap[quote_entry.QuoteTargetKindTypeNoticeChannel], channel.Id)
			}

			if err := w.quoteStore.Set(txCtx, strategy.Id, quote_entry.QuoteKindTypeWarnStrategy, quoteMap); err != nil {
				return err
			}
		}

		return nil
	})
}
