package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto/strategy-dto"
	"github.com/eolinker/apinto-dashboard/entry/strategy-entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	service2 "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"reflect"
	"sort"
	"strings"
	"time"
)

type IStrategyService[T any, K any] interface {
	GetList(ctx context.Context, namespaceId int, clusterName string) ([]*model.Strategy, error)
	GetInfo(ctx context.Context, namespaceId int, uuid string) (*model.StrategyInfoOutput[K], *model.ExtenderData, error)
	CreateStrategy(ctx context.Context, namespaceId int, operator int, clusterName string, input *strategy_dto.StrategyInfoInput[T]) error
	UpdateStrategy(ctx context.Context, namespaceId int, operator int, clusterName string, input *strategy_dto.StrategyInfoInput[T]) error
	DeleteStrategy(ctx context.Context, namespaceId, operator int, clusterName, uuid string) error
	RestoreStrategy(ctx context.Context, namespaceId, userID int, clusterName, uuid string) error
	UpdateStop(ctx context.Context, namespaceId, operator int, uuid, clusterName string, stop bool) error
	ToPublish(ctx context.Context, namespaceId int, clusterName string) ([]*model.StrategyToPublish[T], error)
	PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*model.StrategyPublishHistory, int, error)
	Publish(ctx context.Context, namespaceId, operator int, clusterName string, input *strategy_dto.StrategyPublish) error

	ChangePriority(ctx context.Context, namespaceId, userId int, clusterName string, maps map[string]int) error
	CheckInput(input *strategy_dto.StrategyInfoInput[T]) error

	checkPriorityReduplicative(ctx context.Context, clusterID, priority int, strategyType, uuid string) (int, error)
	getLatestStrategyVersion(ctx context.Context, strategyID int) (*strategy_entry.StrategyVersion, error)
	toApinto(name, desc string, isStop bool, priority int, filters []strategy_entry.StrategyFiltersConfig, conf T) map[string]interface{}
	encodeConfig(config *T) string
	decodeConfig(config string) *T
	IResetOnlineService
}

type strategyService[T any, K any] struct {
	apintoClient         IApintoClient
	clusterService       IClusterService
	strategyStore        store.IStrategyStore
	strategyStatStore    store.IStrategyStatStore
	publishVersionStore  store.IStrategyPublishVersionStore
	publishHistoryStore  store.IStrategyPublishHistoryStore
	strategyVersionStore store.IStrategyVersionStore
	strategyRuntimeStore store.IStrategyRuntimeStore
	strategyHistory      store.IStrategyHistoryStore
	lock                 IAsynLockService
	applicationService   IApplicationService
	apiService           service2.IAPIService
	service              upstream.IService
	userInfoService      IUserInfoService

	//strategyManager driver_manager.IStrategyDriverManager
	strategyHandler IStrategyHandler[T, K]
}

func NewStrategyService[T any, K any](handler IStrategyHandler[T, K], runtimeKind string) IStrategyService[T, K] {
	s := new(strategyService[T, K])

	bean.Autowired(&s.strategyStore)
	bean.Autowired(&s.strategyStatStore)
	bean.Autowired(&s.strategyVersionStore)
	bean.Autowired(&s.strategyHistory)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.applicationService)
	bean.Autowired(&s.service)
	bean.Autowired(&s.apiService)
	bean.Autowired(&s.lock)
	bean.Autowired(&s.userInfoService)

	s.strategyHandler = handler
	var db store.IDB
	bean.Autowired(&db)
	s.strategyRuntimeStore = store.NewStrategyRuntimeStore(db, runtimeKind)
	s.publishHistoryStore = store.NewStrategyPublishHistoryStore(db, "p_h_"+runtimeKind)
	s.publishVersionStore = store.NewStrategyPublishVersionStore(db, "p_v_"+runtimeKind)

	var strategyCommon IStrategyCommonService
	bean.Autowired(&strategyCommon)
	bean.AddInitializingBeanFunc(func() {
		strategyCommon.AddHandler(s)
	})
	return s
}

func (s *strategyService[T, K]) GetList(ctx context.Context, namespaceId int, clusterName string) ([]*model.Strategy, error) {
	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	strategyType := s.strategyHandler.GetType()

	strategies, err := s.strategyStore.GetListByType(ctx, cluster.Id, strategyType)
	if err != nil {
		return nil, err
	}

	_, publishMaps, err := s.getRuntimePublishMaps(ctx, cluster.Id)

	userIds := common.SliceToSliceIds(strategies, func(t *strategy_entry.Strategy) int {
		return t.Operator
	})

	userInfoMaps, _ := s.userInfoService.GetUserInfoMaps(ctx, userIds...)

	resList := make([]*model.Strategy, 0, len(strategies))
	for _, strategy := range strategies {

		version, err := s.getLatestStrategyVersion(ctx, strategy.Id)
		if err != nil {
			return nil, err
		}

		status := enum.StrategyOnlineStatusNOTGOONLINE //默认未上线
		if publish, ok := publishMaps[strategy.Id]; ok {
			status = enum.StrategyOnlineStatusGOONLINE //在发布记录中表示已上线
			if publish.StrategyVersion.Id != version.Id {
				status = enum.StrategyOnlineStatusTOUPDATE //有变更为待更新
			}
			if publish.StrategyVersion.IsStop != strategy.IsStop {
				status = enum.StrategyOnlineStatusTOUPDATE
			}
			if publish.StrategyVersion.Priority != strategy.Priority {
				status = enum.StrategyOnlineStatusTOUPDATE
			}
		}
		//上线了 并且被删除了 为待删除状态
		if status != enum.StrategyOnlineStatusNOTGOONLINE && strategy.IsDelete {
			status = enum.StrategyOnlineStatusTODELETE
		}

		filters, err := s.getFiltersStr(ctx, namespaceId, version)
		if err != nil {
			return nil, err
		}

		operatorName := ""
		if userInfo, ok := userInfoMaps[strategy.Operator]; ok {
			operatorName = userInfo.NickName
		}

		resList = append(resList, &model.Strategy{
			Strategy:    strategy,
			Version:     version,
			Filters:     filters,
			Conf:        s.strategyHandler.GetListLabel(s.decodeConfig(version.StrategyConfigInfo.Config)),
			OperatorStr: operatorName,
			Status:      status,
		})

	}

	return resList, nil
}

func (s *strategyService[T, K]) GetInfo(ctx context.Context, namespaceId int, uuid string) (*model.StrategyInfoOutput[K], *model.ExtenderData, error) {
	strategyInfo, err := s.strategyStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, nil, err
	}
	version, err := s.getLatestStrategyVersion(ctx, strategyInfo.Id)
	if err != nil {
		return nil, nil, err
	}

	//拼装筛选条件所需要的信息
	extenderData := &model.ExtenderData{
		Api:         make(map[string]*model.RemoteApis),
		Service:     make(map[string]*model.RemoteServices),
		Application: make(map[string]*model.RemoteApplications),
	}
	filters := make([]*model.FilterOutput, 0, len(version.Filters))
	for _, f := range version.Filters {
		filter := &model.FilterOutput{
			Name:   f.Name,
			Values: f.Values,
		}
		if len(f.Values) == 0 {
			continue
		}
		switch filter.Name {
		case enum.FilterApplication:
			filter.Type = enum.FilterTypeRemote
			filter.Title = "应用"
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = "所有应用"
			} else {
				apps, err := s.applicationService.AppListByUUIDS(ctx, namespaceId, f.Values)
				if err != nil {
					return nil, nil, err
				}
				if len(apps) == 0 {
					continue
				}
				labels := make([]string, len(apps))
				for i, app := range apps {
					extenderData.Application[app.IdStr] = &model.RemoteApplications{
						Name: app.Name,
						Uuid: app.IdStr,
						Desc: app.Desc,
					}
					labels[i] = app.Name
				}
				filter.Label = strings.Join(labels, ",")
			}
		case enum.FilterApi:
			filter.Type = enum.FilterTypeRemote
			filter.Title = "API"
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = "所有API"
			} else {
				apis, err := s.apiService.GetAPIRemoteByUUIDS(ctx, namespaceId, f.Values)
				if err != nil {
					return nil, nil, err
				}
				if len(apis) == 0 {
					continue
				}
				labels := make([]string, len(apis))
				for i, api := range apis {
					extenderData.Api[api.Uuid] = api
					labels[i] = api.Name
				}
				filter.Label = strings.Join(labels, ",")
			}
		case enum.FilterPath:
			filter.Type = enum.FilterTypePattern
			filter.Title = "API路径"
			filter.Label = f.Values[0]
		case enum.FilterService:
			filter.Type = enum.FilterTypeRemote
			filter.Title = "上游服务"
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = "所有上游服务"
			} else {
				services, err := s.service.GetServiceRemoteByNames(ctx, namespaceId, f.Values)
				if err != nil {
					return nil, nil, err
				}
				if len(services) == 0 {
					continue
				}
				labels := make([]string, len(services))
				for i, sv := range services {
					extenderData.Service[sv.Uuid] = sv
					labels[i] = sv.Name
				}
				filter.Label = strings.Join(labels, ",")
			}
		case enum.FilterMethod:
			filter.Type = enum.FilterTypeStatic
			filter.Title = "API请求方式"
			if f.Values[0] == enum.HttpALL {
				filter.Label = "所有API请求方式"
			} else {
				filter.Label = strings.Join(f.Values, ",")
			}
		case enum.FilterIP:
			filter.Type = enum.FilterTypePattern
			filter.Title = "IP"
			filter.Label = strings.Join(f.Values, ",")
		default: //KEY(应用)
			filter.Type = enum.FilterTypeStatic
			filter.Title = fmt.Sprintf("%s(应用)", common.GetFilterAppKey(filter.Name))
			if f.Values[0] == enum.FilterValuesALL {
				filter.Label = fmt.Sprintf("%s(应用)所有值", common.GetFilterAppKey(filter.Name))
			} else {
				filter.Label = strings.Join(f.Values, ",")
			}
		}

		filters = append(filters, filter)
	}
	input := &model.StrategyInfoOutput[T]{
		Strategy: strategyInfo,
		Filters:  filters,
	}

	input.Config = s.decodeConfig(version.StrategyConfigInfo.StrategyVersionConfig.Config)

	return s.strategyHandler.FormatOut(ctx, namespaceId, input), extenderData, nil
}

func (s *strategyService[T, K]) CreateStrategy(ctx context.Context, namespaceId int, operator int, clusterName string, input *strategy_dto.StrategyInfoInput[T]) error {
	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	err = s.lock.Lock(LockNameStrategy, cluster.Id)
	if err != nil {
		return err
	}
	defer s.lock.Unlock(LockNameStrategy, cluster.Id)

	strategyType := s.strategyHandler.GetType()

	priority, err := s.checkPriorityReduplicative(ctx, cluster.Id, input.Priority, strategyType, "")
	if err != nil {
		return err
	}
	input.Priority = priority

	if input.Uuid == "" {
		input.Uuid = uuid.New()
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        input.Uuid,
		Name:        input.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return s.strategyStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		strategyInfo := &strategy_entry.Strategy{
			UUID:        input.Uuid,
			NamespaceId: namespaceId,
			ClusterId:   cluster.Id,
			Type:        strategyType,
			Name:        input.Name,
			Desc:        input.Desc,
			Priority:    input.Priority,
			Operator:    operator,
			CreateTime:  t,
			UpdateTime:  t,
		}
		if err := s.strategyStore.Save(txCtx, strategyInfo); err != nil {
			return err
		}

		//添加版本信息
		strategyVersionInfo := &strategy_entry.StrategyVersion{
			StrategyConfigInfo: strategy_entry.StrategyConfigInfo{
				Priority:              input.Priority,
				Type:                  strategyType,
				StrategyVersionConfig: strategy_entry.StrategyVersionConfig{},
			},
			StrategyId:  strategyInfo.Id,
			NamespaceId: namespaceId,
			Operator:    operator,
			CreateTime:  t,
		}

		filters := make([]strategy_entry.StrategyFiltersConfig, 0, len(input.Filters))
		for _, f := range input.Filters {
			filter := strategy_entry.StrategyFiltersConfig{
				Name:   f.Name,
				Values: f.Values,
			}
			filters = append(filters, filter)
		}
		strategyVersionInfo.Filters = filters

		strategyVersionInfo.StrategyVersionConfig = strategy_entry.StrategyVersionConfig{
			Config: s.encodeConfig(input.Config),
		}

		if err = s.strategyVersionStore.Save(txCtx, strategyVersionInfo); err != nil {
			return err
		}

		if err = s.strategyHistory.HistoryAdd(txCtx, namespaceId, strategyInfo.Id, &strategy_entry.StrategyHistoryInfo{
			Strategy: *strategyInfo,
			Config:   strategyVersionInfo.StrategyConfigInfo,
		}, operator); err != nil {
			return err
		}

		stat := &strategy_entry.StrategyStat{
			StrategyId: strategyInfo.Id,
			VersionId:  strategyVersionInfo.Id,
		}

		//添加版本关联原表信息
		return s.strategyStatStore.Save(txCtx, stat)
	})
}

func (s *strategyService[T, K]) UpdateStrategy(ctx context.Context, namespaceId int, operator int, clusterName string, input *strategy_dto.StrategyInfoInput[T]) error {

	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	err = s.lock.Lock(LockNameStrategy, cluster.Id)
	if err != nil {
		return err
	}
	defer s.lock.Unlock(LockNameStrategy, cluster.Id)
	strategyInfo, err := s.strategyStore.GetByUUID(ctx, input.Uuid)
	if err != nil {
		return err
	}

	strategyType := s.strategyHandler.GetType()

	priority, err := s.checkPriorityReduplicative(ctx, cluster.Id, input.Priority, strategyType, strategyInfo.UUID)
	if err != nil {
		return err
	}

	currentVersion, err := s.getLatestStrategyVersion(ctx, strategyInfo.Id)
	if err != nil {
		return err
	}

	t := time.Now()

	oldStrategyInfo := *strategyInfo

	strategyInfo.Name = input.Name
	strategyInfo.Desc = input.Desc
	strategyInfo.Priority = priority
	strategyInfo.Operator = operator
	strategyInfo.UpdateTime = t

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        input.Uuid,
		Name:        input.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return s.strategyStore.Transaction(ctx, func(txCtx context.Context) error {
		//修改基础数据
		if _, err = s.strategyStore.Update(txCtx, strategyInfo); err != nil {
			return err
		}

		latestVersionConfig := &strategy_entry.StrategyVersion{
			StrategyConfigInfo: strategy_entry.StrategyConfigInfo{
				Priority: strategyInfo.Priority,
				IsStop:   strategyInfo.IsStop,
				Type:     strategyType,
			},
			StrategyId:  strategyInfo.Id,
			NamespaceId: namespaceId,
			Operator:    operator,
			CreateTime:  t,
		}

		filters := make([]strategy_entry.StrategyFiltersConfig, 0, len(input.Filters))
		for _, f := range input.Filters {
			filter := strategy_entry.StrategyFiltersConfig{
				Name:   f.Name,
				Values: f.Values,
			}
			filters = append(filters, filter)
		}
		latestVersionConfig.Filters = filters
		latestVersionConfig.StrategyVersionConfig = strategy_entry.StrategyVersionConfig{
			Config: s.encodeConfig(input.Config),
		}

		//判断配置信息是否有更新
		if s.isStrategyVersionConfChange(latestVersionConfig, currentVersion) {
			if err = s.strategyVersionStore.Save(txCtx, latestVersionConfig); err != nil {
				return err
			}
			//添加版本关联原表信息
			stat := &strategy_entry.StrategyStat{
				StrategyId: strategyInfo.Id,
				VersionId:  latestVersionConfig.Id,
			}
			if err = s.strategyStatStore.Save(txCtx, stat); err != nil {
				return err
			}
		}

		return s.strategyHistory.HistoryEdit(txCtx, namespaceId, strategyInfo.Id, &strategy_entry.StrategyHistoryInfo{
			Strategy: oldStrategyInfo,
			Config:   currentVersion.StrategyConfigInfo,
		}, &strategy_entry.StrategyHistoryInfo{
			Strategy: *strategyInfo,
			Config:   latestVersionConfig.StrategyConfigInfo,
		}, operator)

	})
}

func (s *strategyService[T, K]) DeleteStrategy(ctx context.Context, namespaceId, operator int, clusterName, uuid string) error {
	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	strategyInfo, err := s.strategyStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.lock.Lock(LockNameStrategy, cluster.Id)
	if err != nil {
		return err
	}
	defer s.lock.Unlock(LockNameStrategy, cluster.Id)

	strategyInfo, err = s.strategyStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	//若已软删除
	if strategyInfo.IsDelete {
		return errors.New("strategy is already delete. ")
	}

	//判断strategy是否已发布，若已上线或待更新，则软删；若未上线则硬删
	isOnline := false
	_, publishMap, err := s.getRuntimePublishMaps(ctx, strategyInfo.ClusterId)
	if _, ok := publishMap[strategyInfo.Id]; ok {
		isOnline = true
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        uuid,
		Name:        strategyInfo.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	t := time.Now()
	err = s.strategyStore.Transaction(ctx, func(txCtx context.Context) error {
		if isOnline {
			strategyInfo.IsDelete = true
			strategyInfo.Operator = operator
			strategyInfo.UpdateTime = t
			_, err = s.strategyStore.Update(txCtx, strategyInfo)
			return err
		}

		if _, err = s.strategyStore.Delete(txCtx, strategyInfo.Id); err != nil {
			return err
		}
		version, err := s.getLatestStrategyVersion(ctx, strategyInfo.Id)
		if err != nil {
			return err
		}
		delMap := make(map[string]interface{})
		delMap["`kind`"] = "strategy"
		delMap["`target`"] = strategyInfo.Id
		if _, err = s.strategyStatStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = s.strategyVersionStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		return s.strategyHistory.HistoryDelete(txCtx, namespaceId, strategyInfo.Id, &strategy_entry.StrategyHistoryInfo{
			Strategy: *strategyInfo,
			Config:   version.StrategyConfigInfo,
		}, operator)
	})
	if err != nil {
		return err
	}

	s.lock.DeleteLock(LockNameStrategy, cluster.Id)
	return nil
}

func (s *strategyService[T, K]) RestoreStrategy(ctx context.Context, namespaceId, userID int, clusterName, uuid string) error {
	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	strategyInfo, err := s.strategyStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = s.lock.Lock(LockNameStrategy, cluster.Id)
	if err != nil {
		return err
	}
	defer s.lock.Unlock(LockNameStrategy, cluster.Id)

	strategyInfo, err = s.strategyStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	//若未软删除
	if !strategyInfo.IsDelete {
		return errors.New("strategy is already restore. ")
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        uuid,
		Name:        strategyInfo.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return s.strategyStore.Transaction(ctx, func(txCtx context.Context) error {
		strategyInfo.IsDelete = false
		strategyInfo.Operator = userID
		strategyInfo.UpdateTime = time.Now()
		return s.strategyStore.Save(txCtx, strategyInfo)
	})
}

func (s *strategyService[T, K]) ToPublish(ctx context.Context, namespaceId int, clusterName string) ([]*model.StrategyToPublish[T], error) {

	strategies, err := s.GetList(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	resList := make([]*model.StrategyToPublish[T], 0)
	for _, strategy := range strategies {
		//已上线的状态过滤掉
		if strategy.Status == enum.StrategyOnlineStatusGOONLINE {
			continue
		}
		resList = append(resList, &model.StrategyToPublish[T]{
			Status:          int(strategy.Status),
			Strategy:        strategy.Strategy,
			StrategyVersion: strategy.Version,
		})
	}

	sort.Slice(resList, func(i, j int) bool {
		return resList[i].Strategy.Priority < resList[j].Strategy.Priority
	})
	return resList, nil
}

func (s *strategyService[T, K]) PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*model.StrategyPublishHistory, int, error) {
	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, 0, common.ClusterNotExist
	}

	list, count, err := s.publishHistoryStore.GetByClusterPage(ctx, pageNum, pageSize, cluster.Id)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]*model.StrategyPublishHistory, 0, len(list))

	for _, history := range list {

		details := make([]*model.StrategyPublishHistoryDetails, 0)
		for _, publish := range history.Publish {
			details = append(details, &model.StrategyPublishHistoryDetails{
				Name:     publish.Strategy.Name,
				Priority: publish.Strategy.Priority,
				Status:   publish.Status,
				OptTime:  publish.Strategy.UpdateTime,
			})
		}

		resp = append(resp, &model.StrategyPublishHistory{
			Id:         history.Id,
			Name:       history.VersionName,
			OptType:    history.OptType,
			Operator:   "",
			CreateTime: history.CreateTime,
			Details:    details,
		})
	}

	return resp, count, nil
}

func (s *strategyService[T, K]) UpdateStop(ctx context.Context, namespaceId, operator int, uuid, clusterName string, stop bool) error {
	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	strategy, err := s.strategyStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	version, err := s.getLatestStrategyVersion(ctx, strategy.Id)
	if err != nil {
		return err
	}

	oldValue := &strategy_entry.StrategyHistoryInfo{
		Strategy: *strategy,
		Config:   version.StrategyConfigInfo,
	}

	t := time.Now()
	strategy.UpdateTime = t
	strategy.Operator = operator
	strategy.IsStop = stop

	newValue := &strategy_entry.StrategyHistoryInfo{
		Strategy: *strategy,
		Config:   version.StrategyConfigInfo,
	}

	enableOperate := 1
	if stop {
		enableOperate = 2
	}
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:          uuid,
		Name:          strategy.Name,
		ClusterId:     cluster.Id,
		ClusterName:   clusterName,
		EnableOperate: enableOperate,
	})

	return s.strategyStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = s.strategyHistory.HistoryEdit(txCtx, namespaceId, strategy.Id, oldValue, newValue, operator); err != nil {
			return err
		}
		return s.strategyStore.Save(txCtx, strategy)

	})

}

func (s *strategyService[T, K]) ResetOnline(ctx context.Context, _, clusterId int) {

	client, err := s.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return
	}

	//查询上个版本发布的策略信息
	runtime, publishMaps, err := s.getRuntimePublishMaps(ctx, clusterId)
	if err != nil {
		log.Errorf("strategyService-ResetOnline-GetRuntime clusterId=%d,err=%s", clusterId, err.Error())
		return
	}
	if len(publishMaps) == 0 {
		return
	}
	strategyInfos := make([]interface{}, 0)

	for _, publish := range publishMaps {
		strategy := publish.Strategy
		conf := s.decodeConfig(publish.StrategyVersion.StrategyConfigInfo.Config)
		strategyInfo := s.toApinto(strategy.UUID, strategy.Desc, strategy.IsStop, strategy.Priority, publish.StrategyVersion.Filters, *conf)
		strategyInfos = append(strategyInfos, &strategyInfo)
	}
	if err = client.ForStrategy().Batch(s.strategyHandler.GetBatchSettingName(), strategyInfos); err != nil {
		log.Errorf("strategyService-ResetOnline-Batch clusterId=%d,runtimeId=%d,versionId=%d,err=%s", clusterId, runtime.Id, runtime.VersionId, err.Error())
		return
	}

}

func (s *strategyService[T, K]) Publish(ctx context.Context, namespaceId, operator int, clusterName string, input *strategy_dto.StrategyPublish) error {

	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}
	clusterId := cluster.Id

	publishHistory, err := s.publishHistoryStore.GetByVersionName(ctx, input.VersionName, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if publishHistory != nil && publishHistory.Id > 0 {
		return errors.New("该版本名称已存在")
	}

	//解析ToPublish接口返回的数据
	bytes, err := common.Base64Decode(input.Source)
	if err != nil {
		return err
	}

	publishes := make([]*model.StrategyToPublish[T], 0)

	if err = json.Unmarshal(bytes, &publishes); err != nil {
		return err
	}

	if len(publishes) == 0 {
		return errors.New("当前没有发布任何策略")
	}

	apintoClient, err := s.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	if err = s.lock.Lock(LockNameStrategy, clusterId); err != nil {
		return err
	}
	defer s.lock.Unlock(LockNameStrategy, clusterId)

	//查询上个版本发布的策略信息
	runtime, publishMaps, err := s.getRuntimePublishMaps(ctx, cluster.Id)
	if err != nil {
		return err
	}

	strategyPublishHistoryInfo := make([]*strategy_entry.StrategyPublishConfigInfo, 0)
	//需要物理删除的策略
	deleteStrategyMaps := common.Map[int, *model.StrategyToPublish[T]]{}

	for _, publish := range publishes {

		strategy, err := s.strategyStore.Get(ctx, publish.Strategy.Id)
		if err != nil && err == gorm.ErrRecordNotFound {
			return errors.New("其他人修改数据导致数据错乱，请重新刷新页面")
		}

		publish.StrategyVersion.Priority = strategy.Priority
		publish.StrategyVersion.IsStop = strategy.IsStop

		publishMaps[publish.Strategy.Id] = &strategy_entry.StrategyPublishConfigInfo{
			Strategy:        *publish.Strategy,
			StrategyVersion: *publish.StrategyVersion,
			Status:          int(publish.Status),
		}

		strategyPublishHistoryInfo = append(strategyPublishHistoryInfo, publishMaps[publish.Strategy.Id])
		if publish.Status == enum.StrategyOnlineStatusTODELETE {
			deleteStrategyMaps[publish.Strategy.Id] = publish
			delete(publishMaps, publish.Strategy.Id)
		}
	}

	//全量发布的数据
	strategyVersionPublishConfig := make([]*strategy_entry.StrategyPublishConfigInfo, 0)

	for _, info := range publishMaps {
		strategyVersionPublishConfig = append(strategyVersionPublishConfig, info)
	}
	sort.Slice(strategyVersionPublishConfig, func(i, j int) bool {
		return strategyVersionPublishConfig[i].Strategy.Priority < strategyVersionPublishConfig[j].Strategy.Priority
	})

	t := time.Now()
	if runtime == nil { //整体发布
		runtime = &strategy_entry.StrategyRuntime{
			ClusterId:   clusterId,
			NamespaceId: namespaceId,
			IsOnline:    true,
			Operator:    operator,
			CreateTime:  t,
			UpdateTime:  t,
		}
	} else {
		runtime.Operator = operator
		runtime.UpdateTime = t
	}

	//发布版本信息
	publishVersion := &strategy_entry.StrategyPublishVersion{
		ClusterId:   clusterId,
		NamespaceId: namespaceId,
		Publish:     strategyVersionPublishConfig,
		Operator:    operator,
		CreateTime:  t,
	}

	publishUUIDS := make([]string, 0)
	publishNames := make([]string, 0)
	for _, publish := range publishVersion.Publish {
		publishUUIDS = append(publishUUIDS, publish.Strategy.UUID)
		publishNames = append(publishNames, publish.Strategy.Name)
	}
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        strings.Join(publishUUIDS, ","),
		Name:        strings.Join(publishNames, ","),
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return s.publishVersionStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = s.publishVersionStore.Save(txCtx, publishVersion); err != nil {
			return err
		}

		runtime.VersionId = publishVersion.Id

		if err = s.strategyRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}

		//发布记录
		publishHistory = &strategy_entry.StrategyPublishHistory{
			VersionName: input.VersionName,
			ClusterId:   clusterId,
			NamespaceId: namespaceId,
			Desc:        input.Desc,
			VersionId:   publishVersion.Id,
			Publish:     strategyPublishHistoryInfo,
			OptType:     1,
			Operator:    operator,
			CreateTime:  t,
		}

		if err = s.publishHistoryStore.Insert(txCtx, publishHistory); err != nil {
			return err
		}

		for strategyId, info := range deleteStrategyMaps {

			if _, err = s.strategyStore.Delete(txCtx, strategyId); err != nil {
				return err
			}

			delMap := make(map[string]interface{})
			delMap["`kind`"] = "strategy"
			delMap["`target`"] = strategyId
			if _, err = s.strategyStatStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}

			if _, err = s.strategyVersionStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}

			//删除策略操作记录

			if err = s.strategyHistory.HistoryDelete(txCtx, namespaceId, strategyId, &strategy_entry.StrategyHistoryInfo{
				Strategy: *info.Strategy,
				Config:   info.StrategyVersion.StrategyConfigInfo,
			}, operator); err != nil {
				return err
			}
		}

		strategyInfos := make([]interface{}, 0)
		for _, publish := range publishVersion.Publish {
			strategy := publish.Strategy
			conf := s.decodeConfig(publish.StrategyVersion.StrategyConfigInfo.Config)
			strategyInfo := s.toApinto(strategy.UUID, strategy.Desc, strategy.IsStop, strategy.Priority, publish.StrategyVersion.Filters, *conf)
			strategyInfos = append(strategyInfos, &strategyInfo)
		}

		return apintoClient.ForStrategy().Batch(s.strategyHandler.GetBatchSettingName(), strategyInfos)
	})

}

func (s *strategyService[T, K]) ChangePriority(ctx context.Context, namespaceId, userId int, clusterName string, maps map[string]int) error {

	cluster, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	if err = s.lock.Lock(LockNameStrategy, cluster.Id); err != nil {
		return err
	}
	defer s.lock.Unlock(LockNameStrategy, cluster.Id)

	strategyType := s.strategyHandler.GetType()

	strategies, err := s.strategyStore.GetListByType(ctx, cluster.Id, strategyType)
	if err != nil {
		return err
	}

	strategyMaps := common.SliceToMap(strategies, func(t *strategy_entry.Strategy) string {
		return t.UUID
	})

	historyList := make([]*strategy_entry.StrategyHistory, 0)
	for key, priority := range maps {

		if strategy, ok := strategyMaps[key]; ok && strategy.Priority != priority {
			oldStrategy := *strategy
			strategy.Priority = priority
			version, err := s.getLatestStrategyVersion(ctx, strategy.Id)
			if err != nil {
				return err
			}

			history := &strategy_entry.StrategyHistory{

				OldValue: strategy_entry.StrategyHistoryInfo{
					Strategy: oldStrategy,
					Config:   version.StrategyConfigInfo,
				},
				NewValue: strategy_entry.StrategyHistoryInfo{
					Strategy: *strategy,
					Config:   version.StrategyConfigInfo,
				},
			}

			history.NewValue.Strategy.Priority = priority

			historyList = append(historyList, history)
		}
	}

	//没有变更记录 说明没有更改过优先级
	if len(historyList) == 0 {
		return nil
	}

	return s.strategyStore.Transaction(ctx, func(txCtx context.Context) error {

		for _, h := range historyList {
			if err = s.strategyHistory.HistoryEdit(txCtx, namespaceId, h.StrategyId, &h.OldValue, &h.NewValue, userId); err != nil {
				return err
			}
		}

		return s.strategyStore.UpdatePriority(txCtx, maps)
	})

}

// getRuntimePublishMaps 获取当前集群已经发布的策略信息
func (s *strategyService[T, K]) getRuntimePublishMaps(ctx context.Context, clusterId int) (*strategy_entry.StrategyRuntime, map[int]*strategy_entry.StrategyPublishConfigInfo, error) {
	var runtime *strategy_entry.StrategyRuntime
	var err error

	runtime, err = s.strategyRuntimeStore.GetForCluster(ctx, clusterId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, nil, err
	}

	//查询当前发布的版本
	publishMaps := common.Map[int, *strategy_entry.StrategyPublishConfigInfo]{}
	if runtime != nil {
		publishVersion, err := s.publishVersionStore.Get(ctx, runtime.VersionId)
		if err != nil {
			return nil, nil, err
		}
		publishMaps = common.SliceToMap(publishVersion.Publish, func(t *strategy_entry.StrategyPublishConfigInfo) int {
			return t.Strategy.Id
		})
	}

	return runtime, publishMaps, nil
}

func (s *strategyService[T, K]) getLatestStrategyVersion(ctx context.Context, strategyID int) (*strategy_entry.StrategyVersion, error) {
	stat, err := s.strategyStatStore.Get(ctx, strategyID)
	if err != nil {
		return nil, err
	}
	return s.strategyVersionStore.Get(ctx, stat.VersionId)
}

func (s *strategyService[T, K]) checkPriorityReduplicative(ctx context.Context, clusterID, priority int, strategyType, uuid string) (int, error) {
	if priority < 1 {
		// 策略表中最大priority + 1
		maxPriority, err := s.strategyStore.GetMaxPriority(ctx, clusterID, strategyType)
		if err != nil {
			return 0, err
		}

		return maxPriority + 1, nil
	}

	strategy, err := s.strategyStore.GetByPriority(ctx, clusterID, priority, strategyType)
	if err != nil {
		return 0, err
	}
	if strategy != nil && strategy.Id > 0 && strategy.UUID != uuid {
		return 0, errors.New("priority is reduplicative. ")
	}

	return priority, nil
}

func (s *strategyService[T, K]) isStrategyVersionConfChange(latest *strategy_entry.StrategyVersion, current *strategy_entry.StrategyVersion) bool {
	if latest.Priority != current.Priority {
		return true
	}

	if !reflect.DeepEqual(latest.Filters, current.Filters) {
		return true
	}

	return !reflect.DeepEqual(latest.StrategyVersionConfig, current.StrategyVersionConfig)
}

func (s *strategyService[T, K]) getFiltersStr(ctx context.Context, namespaceId int, version *strategy_entry.StrategyVersion) (string, error) {
	//filters := ""
	filterList := make([]string, 0)
	for _, filter := range version.Filters {
		switch filter.Name {
		case enum.FilterApplication:
			if len(filter.Values) > 0 {
				if filter.Values[0] == enum.FilterValuesALL {
					filters := "[应用：全部应用]"
					filterList = append(filterList, filters)
				} else {
					applications, err := s.applicationService.AppListByUUIDS(ctx, namespaceId, filter.Values)
					if err != nil {
						return "", err
					}
					filters := "[应用："
					for i, application := range applications {
						if i == len(applications)-1 {
							filters += application.Name + "]"
						} else {
							filters += application.Name + ","
						}
					}
					filterList = append(filterList, filters)
				}
			}
		case enum.FilterApi:
			if len(filter.Values) > 0 {
				if filter.Values[0] == enum.FilterValuesALL {
					filters := "[API：全部API]"
					filterList = append(filterList, filters)
				} else {
					apis, err := s.apiService.GetAPIRemoteByUUIDS(ctx, namespaceId, filter.Values)
					if err != nil {
						return "", err
					}
					filters := "[API："
					for i, api := range apis {
						if i == len(apis)-1 {
							filters += api.Name + "]"
						} else {
							filters += api.Name + ","
						}
					}
					filterList = append(filterList, filters)
				}
			}
		case enum.FilterService:
			if len(filter.Values) > 0 {
				if filter.Values[0] == enum.FilterValuesALL {
					filters := "[上游服务：全部上游服务]"
					filterList = append(filterList, filters)
				} else {
					services, err := s.service.GetServiceRemoteByNames(ctx, namespaceId, filter.Values)
					if err != nil {
						return "", err
					}
					filters := "[上游服务："
					for i, val := range services {
						if i == len(services)-1 {
							filters += val.Name + "]"
						} else {
							filters += val.Name + ","
						}
					}
					filterList = append(filterList, filters)
				}
			}
		case enum.FilterMethod:
			if len(filter.Values) > 0 {
				if filter.Values[0] == enum.FilterValuesALL {
					filters := "[API请求方式：全部请求方式]"
					filterList = append(filterList, filters)
				} else {
					filters := fmt.Sprintf("[API请求方式：%s]", strings.Join(filter.Values, ","))
					filterList = append(filterList, filters)
				}
			}
		case enum.FilterPath:
			if len(filter.Values) > 0 {
				filters := fmt.Sprintf("[API路径：%s]", filter.Values[0])
				filterList = append(filterList, filters)
			}
		case enum.FilterIP:
			if len(filter.Values) > 0 {
				filters := fmt.Sprintf("[IP：%s]", strings.Join(filter.Values, ","))
				filterList = append(filterList, filters)
			}
		default:
			//appKey
			key := common.GetFilterAppKey(filter.Name)
			if filter.Name != key {
				if len(filter.Values) > 0 {
					if filter.Values[0] == enum.FilterValuesALL {
						filters := fmt.Sprintf("[%s(应用)：全部数据]", key)
						filterList = append(filterList, filters)
					} else {
						filters := fmt.Sprintf("[%s(应用)：", key)
						for i, val := range filter.Values {
							if i == len(filter.Values)-1 {
								filters += val + "]"
							} else {
								filters += val + ","
							}
						}
						filterList = append(filterList, filters)
					}
				}
			}
		}
	}
	return strings.Join(filterList, ";"), nil
}

func (s *strategyService[T, K]) CheckInput(input *strategy_dto.StrategyInfoInput[T]) error {
	return s.strategyHandler.CheckInput(input)
}

func (s *strategyService[T, K]) toApinto(name, desc string, isStop bool, priority int, filters []strategy_entry.StrategyFiltersConfig, conf T) map[string]interface{} {
	limitingFilters := make(map[string][]string)

	for _, filter := range filters {
		limitingFilters[filter.Name] = filter.Values
	}

	confName := s.strategyHandler.GetConfName()
	return map[string]interface{}{
		"name":        name,
		"stop":        isStop,
		"description": desc,
		"priority":    priority,
		"filters":     limitingFilters,
		confName:      s.strategyHandler.ToApintoConfig(conf),
	}
}

func (s *strategyService[T, K]) encodeConfig(config *T) string {
	data, _ := json.Marshal(config)
	return string(data)
}

func (s *strategyService[T, K]) decodeConfig(config string) *T {
	data := new(T)
	_ = json.Unmarshal([]byte(config), data)

	return data
}
