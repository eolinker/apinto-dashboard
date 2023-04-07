package upstream_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	api "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/discovery"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	upstream_model "github.com/eolinker/apinto-dashboard/modules/upstream/model"
	upstream_store "github.com/eolinker/apinto-dashboard/modules/upstream/store"
	upstream_dto "github.com/eolinker/apinto-dashboard/modules/upstream/upstream-dto"
	upstream_entry2 "github.com/eolinker/apinto-dashboard/modules/upstream/upstream-entry"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
	variable_model "github.com/eolinker/apinto-dashboard/modules/variable/variable-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type service struct {
	clusterService        cluster.IClusterService
	apintoClient          cluster.IApintoClient
	clusterNodeService    cluster.IClusterNodeService
	discoveryService      discovery.IDiscoveryService
	namespaceService      namespace.INamespaceService
	globalVariableService variable.IGlobalVariableService
	lockService           locker_service.IAsynLockService
	variableService       variable.IClusterVariableService
	apiService            api.IAPIService
	userInfoService       user.IUserInfoService
	serviceRuntimeStore   upstream_store.IServiceRuntimeStore
	serviceStore          upstream_store.IServiceStore
	serviceVersionStore   upstream_store.IServiceVersionStore
	serviceStatStore      upstream_store.IServiceStatStore
	quoteStore            quote_store.IQuoteStore
	historyStore          upstream_store.IServiceHistoryStore
}

func newServiceService() upstream.IService {
	s := &service{}
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.clusterNodeService)
	bean.Autowired(&s.serviceRuntimeStore)
	bean.Autowired(&s.serviceVersionStore)
	bean.Autowired(&s.serviceStatStore)
	bean.Autowired(&s.serviceStore)
	bean.Autowired(&s.variableService)
	bean.Autowired(&s.quoteStore)
	bean.Autowired(&s.globalVariableService)
	bean.Autowired(&s.namespaceService)
	bean.Autowired(&s.discoveryService)
	bean.Autowired(&s.lockService)
	bean.Autowired(&s.apiService)
	bean.Autowired(&s.historyStore)
	bean.Autowired(&s.userInfoService)
	return s
}

func (s *service) GetServiceList(ctx context.Context, namespaceID int, searchName string, pageNum int, pageSize int) ([]*upstream_model.ServiceListItem, int, error) {
	var sl []*upstream_entry2.Service
	sl, total, err := s.serviceStore.GetListPage(ctx, namespaceID, searchName, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*upstream_model.ServiceListItem, 0, len(sl))
	for _, item := range sl {
		info, err := s.GetLatestServiceVersion(ctx, item.Id)
		if err != nil {
			return nil, 0, err
		}

		li := &upstream_model.ServiceListItem{
			Name:        item.Name,
			UUID:        item.UUID,
			Scheme:      info.Scheme,
			DiscoveryID: info.DiscoveryId,
			DriverName:  info.DriverName,
			Config:      info.FormatAddr,
			UpdateTime:  item.UpdateTime,
			IsDelete:    s.isDelete(ctx, item.Id),
		}

		list = append(list, li)
	}

	return list, total, nil
}
func (s *service) GetServiceListAll(ctx context.Context, namespaceID int) ([]*upstream_model.ServiceListItem, error) {
	sl, err := s.serviceStore.GetListAll(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	list := make([]*upstream_model.ServiceListItem, 0, len(sl))
	for _, item := range sl {

		li := &upstream_model.ServiceListItem{
			Name:       item.Name,
			UUID:       item.UUID,
			UpdateTime: item.UpdateTime,
		}

		list = append(list, li)
	}

	return list, nil
}
func (s *service) GetServiceListByNames(ctx context.Context, namespaceID int, names []string) ([]*upstream_model.ServiceListItem, error) {
	sl, err := s.serviceStore.GetByNames(ctx, namespaceID, names)
	if err != nil {
		return nil, err
	}

	list := make([]*upstream_model.ServiceListItem, 0, len(sl))
	for _, item := range sl {

		li := &upstream_model.ServiceListItem{
			Name:       item.Name,
			UUID:       item.UUID,
			UpdateTime: item.UpdateTime,
		}

		list = append(list, li)
	}

	return list, nil
}

func (s *service) isDelete(ctx context.Context, serviceId int) bool {
	quotedSet, err := s.quoteStore.GetTargetQuote(ctx, serviceId, quote_entry.QuoteTargetKindTypeService)
	if err != nil {
		return false
	}
	//有API绑定了上游服务
	if len(quotedSet[quote_entry.QuoteKindTypeAPI]) > 0 {
		return false
	}

	//在某个集群上线
	onlineCount, _ := s.serviceRuntimeStore.OnlineCount(ctx, serviceId)
	return onlineCount == 0
}

func (s *service) isDiscoveryCanDelete(ctx context.Context, namespaceID, serviceId int) (bool, error) {
	quotedSet, err := s.quoteStore.GetTargetQuote(ctx, serviceId, quote_entry.QuoteTargetKindTypeService)
	if err != nil {
		return false, err
	}

	for _, apiID := range quotedSet[quote_entry.QuoteKindTypeAPI] {
		name, err := s.apiService.GetAPINameByID(ctx, apiID)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("Service is in use by api %s. ", name)
	}

	clusters, err := s.clusterService.GetByNamespaceId(ctx, namespaceID)
	if err != nil {
		return false, err
	}
	for _, clusterInfo := range clusters {
		runtime, _ := s.serviceRuntimeStore.GetForCluster(ctx, serviceId, clusterInfo.Id)
		if runtime != nil && runtime.IsOnline {
			return false, fmt.Errorf("service is online in cluster %s. ", clusterInfo.Name)
		}
	}
	return true, nil

}

func (s *service) GetServiceInfo(ctx context.Context, namespaceID int, serviceName string) (*upstream_model.ServiceInfo, error) {
	serviceInfo, err := s.serviceStore.GetByName(ctx, namespaceID, serviceName)
	if err != nil {
		return nil, err
	}

	version, err := s.GetLatestServiceVersion(ctx, serviceInfo.Id)
	if err != nil {
		return nil, err
	}

	info := &upstream_model.ServiceInfo{
		ServiceVersion: version,
		Name:           serviceInfo.Name,
		Desc:           serviceInfo.Desc,
		ServiceId:      serviceInfo.Id,
		UUID:           serviceInfo.UUID,
	}

	return info, nil
}

func (s *service) CreateService(ctx context.Context, namespaceID, userId int, input *upstream_dto.ServiceInfo, variableList []string) (int, error) {
	input.Name = strings.ToLower(input.Name)
	//服务发现name查重
	_, err := s.serviceStore.GetByName(ctx, namespaceID, input.Name)
	if err != gorm.ErrRecordNotFound {
		if err == nil {
			return 0, fmt.Errorf("name %s already exist. ", input.Name)
		}
		return 0, err
	}

	if input.UUID == "" {
		input.UUID = uuid.New()
	}
	input.UUID = strings.ToLower(input.UUID)

	t := time.Now()
	serviceInfo := &upstream_entry2.Service{
		NamespaceId: namespaceID,
		UUID:        input.UUID,
		Name:        input.Name,
		Desc:        input.Desc,
		Operator:    userId,
		CreateTime:  t,
		UpdateTime:  t,
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.Name,
	})
	serviceID := 0
	return serviceID, s.serviceStore.Transaction(ctx, func(txCtx context.Context) error {

		if err := s.serviceStore.Save(txCtx, serviceInfo); err != nil {
			return err
		}
		serviceID = serviceInfo.Id
		//添加版本信息
		serviceVersionInfo := &upstream_entry2.ServiceVersion{
			ServiceId:   serviceInfo.Id,
			NamespaceID: namespaceID,
			ServiceVersionConfig: upstream_entry2.ServiceVersionConfig{
				DiscoveryId: input.DiscoveryID,
				DriverName:  input.DriverName,
				Scheme:      input.Scheme,
				Balance:     input.Balance,
				Timeout:     input.Timeout,
				FormatAddr:  input.FormatAddr,
				Config:      input.Config,
			},
			Operator:   userId,
			CreateTime: t,
		}

		if err := s.serviceVersionStore.Save(txCtx, serviceVersionInfo); err != nil {
			return err
		}
		//添加版本关联原表信息
		if err := s.serviceStatStore.Save(txCtx, &upstream_entry2.ServiceStat{
			ServiceId: serviceInfo.Id,
			VersionId: serviceVersionInfo.Id,
		}); err != nil {
			return err
		}
		//往引用表插入所引用服务发现
		if err := s.quoteStore.Set(txCtx, serviceInfo.Id, quote_entry.QuoteKindTypeService, quote_entry.QuoteTargetKindTypeDiscovery, input.DiscoveryID); err != nil {
			return err
		}

		//往引用表插入所使用的环境变量
		if len(variableList) > 0 {
			variables, err := s.globalVariableService.GetByKeys(ctx, namespaceID, variableList)
			if err != nil {
				return err
			}
			variableIds := common.SliceToSliceIds(variables, func(t *variable_model.GlobalVariable) int {
				return t.Id
			})

			err = s.quoteStore.Set(txCtx, serviceInfo.Id, quote_entry.QuoteKindTypeService, quote_entry.QuoteTargetKindTypeVariable, variableIds...)
			if err != nil {
				return err
			}
		}

		return s.historyStore.HistoryAdd(txCtx, namespaceID, serviceInfo.Id, &upstream_entry2.ServiceHistoryInfo{
			Service: *serviceInfo,
			Config:  serviceVersionInfo.ServiceVersionConfig,
		}, userId)
	})
}

func (s *service) UpdateService(ctx context.Context, namespaceID, userId int, input *upstream_dto.ServiceInfo, variableList []string) error {

	serviceInfo, err := s.serviceStore.GetByName(ctx, namespaceID, input.Name)
	if err != nil {
		return err
	}

	isUpdateVersion := false

	currentVersion, err := s.GetLatestServiceVersion(ctx, serviceInfo.Id)
	if err != nil {
		return err
	}

	switch true {
	case currentVersion.Scheme != input.Scheme:
		isUpdateVersion = true
	case currentVersion.Balance != input.Balance:
		isUpdateVersion = true
	case currentVersion.Timeout != input.Timeout:
		isUpdateVersion = true
	case currentVersion.Config != input.Config:
		isUpdateVersion = true
	}

	oldServiceInfo := *serviceInfo

	serviceInfo.Desc = input.Desc
	serviceInfo.Operator = userId
	serviceInfo.UpdateTime = time.Now()

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.Name,
	})

	return s.serviceStore.Transaction(ctx, func(txCtx context.Context) error {

		//修改基础数据
		if _, err = s.serviceStore.Update(txCtx, serviceInfo); err != nil {
			return err
		}
		if isUpdateVersion {
			//修改配置信息 并更新一个版本
			serviceVersionInfo := &upstream_entry2.ServiceVersion{
				ServiceId: serviceInfo.Id,
				ServiceVersionConfig: upstream_entry2.ServiceVersionConfig{
					DiscoveryId: input.DiscoveryID,
					DriverName:  input.DriverName,
					Scheme:      input.Scheme,
					Balance:     input.Balance,
					Timeout:     input.Timeout,
					FormatAddr:  input.FormatAddr,
					Config:      input.Config,
				},
				NamespaceID: namespaceID,
				Operator:    userId,
				CreateTime:  time.Now(),
			}
			if err = s.serviceVersionStore.Save(txCtx, serviceVersionInfo); err != nil {
				return err
			}
			//添加版本关联原表信息
			stat := &upstream_entry2.ServiceStat{
				ServiceId: serviceInfo.Id,
				VersionId: serviceVersionInfo.Id,
			}
			if err = s.serviceStatStore.Save(txCtx, stat); err != nil {
				return err
			}

			if err = s.historyStore.HistoryEdit(txCtx, namespaceID, serviceInfo.Id, &upstream_entry2.ServiceHistoryInfo{
				Service: oldServiceInfo,
				Config:  currentVersion.ServiceVersionConfig,
			}, &upstream_entry2.ServiceHistoryInfo{
				Service: *serviceInfo,
				Config:  serviceVersionInfo.ServiceVersionConfig,
			}, userId); err != nil {
				return err
			}
		}

		//往引用表插入所引用服务发现
		if err = s.quoteStore.Set(txCtx, serviceInfo.Id, quote_entry.QuoteKindTypeService, quote_entry.QuoteTargetKindTypeDiscovery, input.DiscoveryID); err != nil {
			return err
		}

		//往引用表插入所使用的环境变量
		if len(variableList) > 0 {
			variables, err := s.globalVariableService.GetByKeys(ctx, namespaceID, variableList)

			variableIds := common.SliceToSliceIds(variables, func(t *variable_model.GlobalVariable) int {
				return t.Id
			})

			if err = s.quoteStore.Set(txCtx, serviceInfo.Id, quote_entry.QuoteKindTypeService, quote_entry.QuoteTargetKindTypeVariable, variableIds...); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *service) DeleteService(ctx context.Context, namespaceID, userId int, serviceName string) error {
	serviceInfo, err := s.serviceStore.GetByName(ctx, namespaceID, serviceName)
	if err != nil {
		return err
	}

	if err = s.lockService.Lock(locker_service.LockNameService, serviceInfo.Id); err != nil {
		return err
	}
	defer s.lockService.Unlock(locker_service.LockNameService, serviceInfo.Id)
	serviceInfo, err = s.serviceStore.GetByName(ctx, namespaceID, serviceName)
	if err != nil {
		return err
	}

	onlineCount, err := s.serviceRuntimeStore.OnlineCount(ctx, serviceInfo.Id)
	if err != nil {
		return err
	}
	if onlineCount > 0 {
		return errors.New("服务已上线不可删除")
	}

	_, err = s.isDiscoveryCanDelete(ctx, namespaceID, serviceInfo.Id)
	if err != nil {
		return err
	}
	version, err := s.GetLatestServiceVersion(ctx, serviceInfo.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: serviceInfo.UUID,
		Name: serviceName,
	})

	err = s.serviceStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = s.serviceStore.Delete(txCtx, serviceInfo.Id); err != nil {
			return err
		}

		delMap := make(map[string]interface{})
		delMap["`kind`"] = "service"
		delMap["`target`"] = serviceInfo.Id

		if _, err = s.serviceStatStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = s.serviceVersionStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = s.serviceRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if err = s.historyStore.HistoryDelete(txCtx, namespaceID, serviceInfo.Id, upstream_entry2.ServiceHistoryInfo{
			Service: *serviceInfo,
			Config:  version.ServiceVersionConfig,
		}, userId); err != nil {
			return err
		}

		return s.quoteStore.DelBySource(txCtx, serviceInfo.Id, quote_entry.QuoteKindTypeService)
	})
	if err != nil {
		return err
	}
	s.lockService.DeleteLock(locker_service.LockNameService, serviceInfo.Id)
	return nil
}

func (s *service) GetServiceEnum(ctx context.Context, namespaceID int, searchName string) ([]string, error) {
	serviceList, _, err := s.serviceStore.GetListPage(ctx, namespaceID, searchName, 0, 0)
	if err != nil {
		return nil, err
	}

	serviceNameList := make([]string, 0, len(serviceList))
	for _, info := range serviceList {
		serviceNameList = append(serviceNameList, info.Name)
	}

	return serviceNameList, nil
}

// OnlineList 上线管理列表
func (s *service) OnlineList(ctx context.Context, namespaceId int, serviceName string) ([]*upstream_model.ServiceOnline, error) {
	serviceInfo, err := s.serviceStore.GetByName(ctx, namespaceId, serviceName)
	if err != nil {
		return nil, err
	}

	serviceId := serviceInfo.Id
	//获取工作空间下的所有集群
	clusters, err := s.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *cluster_model.Cluster) int {
		return t.Id
	})

	//获取当前服务下集群运行的版本
	runtimes, err := s.serviceRuntimeStore.GetByTarget(ctx, serviceId)
	if err != nil {
		return nil, err
	}
	runtimeMaps := common.SliceToMap(runtimes, func(t *upstream_entry2.ServiceRuntime) int {
		return t.ClusterId
	})

	//最新版本
	newVersion, err := s.GetLatestServiceVersion(ctx, serviceId)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(runtimes, func(t *upstream_entry2.ServiceRuntime) int {
		return t.Operator
	})

	userInfoMaps, _ := s.userInfoService.GetUserInfoMaps(ctx, userIds...)

	list := make([]*upstream_model.ServiceOnline, 0, len(clusters))
	for _, clusterInfo := range clusterMaps {
		serviceOnline := &upstream_model.ServiceOnline{
			ClusterID:   clusterInfo.Id,
			ClusterName: clusterInfo.Name,
			Env:         clusterInfo.Env,
			Status:      1, //默认为未上线状态
		}
		if runtime, ok := runtimeMaps[clusterInfo.Id]; ok {
			if runtime.IsOnline {
				serviceOnline.Status = 3
			} else {
				serviceOnline.Status = 2
			}

			if userInfo, uOk := userInfoMaps[runtime.Operator]; uOk {
				serviceOnline.Operator = userInfo.NickName
			}

			serviceOnline.UpdateTime = runtime.UpdateTime

			if serviceOnline.Status == 3 && runtime.VersionId != newVersion.Id {
				serviceOnline.Status = 4
			}
		}
		list = append(list, serviceOnline)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Status > list[j].Status
	})
	return list, nil
}

// OnlineService  上线服务
func (s *service) OnlineService(ctx context.Context, namespaceId, operator int, serviceName, clusterName string) (*frontend_model.Router, error) {
	serviceInfo, err := s.serviceStore.GetByName(ctx, namespaceId, serviceName)
	if err != nil {
		return nil, err
	}

	serviceId := serviceInfo.Id
	t := time.Now()

	//获取当前集群信息
	clusterInfo, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}
	clusterId := clusterInfo.Id

	if err = s.lockService.Lock(locker_service.LockNameService, serviceInfo.Id); err != nil {
		return nil, err
	}
	defer s.lockService.Unlock(locker_service.LockNameService, serviceInfo.Id)

	//拿到锁后需要重新获取下信息
	serviceInfo, err = s.serviceStore.GetByName(ctx, namespaceId, serviceName)
	if err != nil {
		return nil, err
	}

	namespaceInfo, err := s.namespaceService.GetById(namespaceId)
	if err != nil {
		return nil, err
	}

	//获取当前服务的版本
	version, err := s.GetLatestServiceVersion(ctx, serviceId)
	if err != nil {
		return nil, err
	}

	var discoveryName = ""
	var discoveryId = version.ServiceVersionConfig.DiscoveryId
	if discoveryId == 0 { //表示静态服务发现
		router := &frontend_model.Router{
			Name:   frontend_model.RouterNameClusterVariable,
			Params: make(map[string]string),
		}
		router.Params["cluster_name"] = clusterName
		//服务引用的环境变量
		quoteMaps, err := s.quoteStore.GetSourceQuote(ctx, serviceId, quote_entry.QuoteKindTypeService)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if variableIds, ok := quoteMaps[quote_entry.QuoteTargetKindTypeVariable]; ok {
			//通用错误处理文档，当数据库找不到数据时返回
			var errMsg error
			globalVariable, err := s.globalVariableService.GetById(ctx, variableIds[0])
			if err != nil {
				return nil, err
			}
			errMsg = errors.New(fmt.Sprintf("${%s}未上线到{%s}，上线/更新失败", globalVariable.Key, clusterName))

			//获取集群正在运行的环境变量版本
			variableVersion, err := s.variableService.GetPublishVersion(ctx, clusterInfo.Id)
			if err != nil || variableVersion == nil {
				return router, errMsg
			}

			//已发布的环境变量
			toMap := common.SliceToMap(variableVersion.ClusterVariable, func(t *variable_entry.ClusterVariable) int {
				return t.VariableId
			})

			for _, variableId := range variableIds {
				if _, ok := toMap[variableId]; !ok {
					//当前环境变量未发布
					globalVariable, err = s.globalVariableService.GetById(ctx, variableId)
					if err != nil {
						return nil, err
					}
					return router, errors.New(fmt.Sprintf("${%s}未上线到{%s}，上线/更新失败", globalVariable.Key, clusterName))
				}
			}
		}
	} else {

		//动态服务的话 需要判断这个动态服务有没有上线
		discoveryName, err = s.discoveryService.GetDiscoveryName(ctx, discoveryId)
		if err != nil {
			return nil, err
		}
		if !s.discoveryService.IsOnline(ctx, clusterId, discoveryId) {
			router := &frontend_model.Router{
				Name:   frontend_model.RouterNameDiscoveryOnline,
				Params: make(map[string]string),
			}
			router.Params["discovery_name"] = discoveryName
			return router, errors.New(fmt.Sprintf("绑定的%s未上线到%s", discoveryName, clusterName))
		}
	}

	//获取当前运行的版本
	runtime, err := s.serviceRuntimeStore.GetForCluster(ctx, serviceId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	client, err := s.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return nil, err
	}

	_, _, driverInfo, err := s.discoveryService.GetServiceDiscoveryDriverByID(ctx, discoveryId)
	if err != nil {
		return nil, err
	}

	//判断是否是更新
	isApintoUpdate := false
	if runtime == nil {
		runtime = &upstream_entry2.ServiceRuntime{
			NamespaceId: namespaceId,
			ServiceId:   serviceId,
			ClusterId:   clusterId,
			VersionId:   version.Id,
			IsOnline:    true,
			Operator:    operator,
			CreateTime:  t,
			UpdateTime:  t,
		}
	} else {
		if runtime.IsOnline {
			isApintoUpdate = true
		}
		runtime.IsOnline = true
		runtime.UpdateTime = t
		runtime.VersionId = version.Id
		runtime.Operator = operator
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        serviceInfo.UUID,
		Name:        serviceName,
		ClusterId:   clusterId,
		ClusterName: clusterName,
		PublishType: 1,
	})

	return nil, s.serviceRuntimeStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = s.serviceRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}

		config := version.ServiceVersionConfig
		serviceConfig := driverInfo.ToApinto(strings.ToLower(serviceInfo.Name), namespaceInfo.Name, serviceInfo.Desc, config.Scheme, config.Balance, strings.ToLower(discoveryName), "http", config.Timeout, []byte(config.Config))
		if isApintoUpdate {
			return client.ForService().Update(strings.ToLower(serviceInfo.Name), *serviceConfig)
		}

		return client.ForService().Create(*serviceConfig)
	})
}

func (s *service) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	runtimes, err := s.serviceRuntimeStore.GetByCluster(ctx, clusterId)
	if err != nil {
		log.Errorf("service-ResetOnline-getRuntimes clusterId=%d err=%s", clusterId, err.Error())
		return
	}
	client, err := s.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		log.Errorf("service-ResetOnline-getClient clusterId=%d err=%s", clusterId, err.Error())
		return
	}
	namespaceInfo, err := s.namespaceService.GetById(namespaceId)
	if err != nil {
		return
	}

	for _, runtime := range runtimes {
		if !runtime.IsOnline {
			continue
		}

		discoveryName := ""
		version, err := s.serviceVersionStore.Get(ctx, runtime.VersionId)
		if err != nil {
			log.Errorf("service-ResetOnline-getVersion versionId=%d, clusterId=%d, err=%s", runtime.VersionId, clusterId, err.Error())
			continue
		}
		discoveryId := version.DiscoveryId

		if version.DiscoveryId > 0 {
			discoveryName, err = s.discoveryService.GetDiscoveryName(ctx, discoveryId)
			if err != nil {
				log.Errorf("service-ResetOnline-getDiscoveryName discoveryId=%d, clusterId=%d, err=%s", discoveryId, clusterId, err.Error())
				continue
			}
		}

		_, _, driverInfo, err := s.discoveryService.GetServiceDiscoveryDriverByID(ctx, discoveryId)
		if err != nil {
			log.Errorf("service-ResetOnline-getDriverInfo discoveryId=%d clusterId=%d err=%s", discoveryId, clusterId, err.Error())
			continue
		}

		serviceInfo, err := s.serviceStore.Get(ctx, runtime.ServiceId)
		if err != nil {
			log.Errorf("service-ResetOnline-getService serviceId=%d clusterId=%d err=%s", runtime.ServiceId, clusterId, err.Error())
			continue
		}

		versionConfig := version.ServiceVersionConfig
		serviceConfig := driverInfo.ToApinto(strings.ToLower(serviceInfo.Name), namespaceInfo.Name, serviceInfo.Desc, versionConfig.Scheme, versionConfig.Balance, strings.ToLower(discoveryName), "http", versionConfig.Timeout, []byte(versionConfig.Config))
		if err = client.ForService().Create(*serviceConfig); err != nil {
			log.Errorf("service-ResetOnline-apintoCreate serviceConfig=%v clusterId=%d err=%s", serviceConfig, clusterId, err.Error())
			continue
		}
	}
}

// OfflineService 下线服务
func (s *service) OfflineService(ctx context.Context, namespaceId, operator int, serviceName, clusterName string) error {

	serviceInfo, err := s.serviceStore.GetByName(ctx, namespaceId, serviceName)
	if err != nil {
		return err
	}

	serviceId := serviceInfo.Id

	//获取当前集群信息
	clusterInfo, err := s.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	quote, err := s.quoteStore.GetTargetQuote(ctx, serviceInfo.Id, quote_entry.QuoteTargetKindTypeService)
	if err != nil {
		return err
	}

	for kindType, ids := range quote {
		switch kindType {
		case quote_entry.QuoteKindTypeAPI:
			for _, apiId := range ids {
				if s.apiService.IsAPIOnline(ctx, clusterInfo.Id, apiId) {
					name, err := s.apiService.GetAPINameByID(ctx, apiId)
					if err != nil {
						return err
					}
					return errors.New(fmt.Sprintf("当前服务被{%s}API所引用", name))
				}
			}
		}
	}

	if err = s.lockService.Lock(locker_service.LockNameService, serviceInfo.Id); err != nil {
		return err
	}
	defer s.lockService.Unlock(locker_service.LockNameService, serviceInfo.Id)

	//拿到锁后需要重新获取下信息
	serviceInfo, err = s.serviceStore.GetByName(ctx, namespaceId, serviceName)
	if err != nil {
		return err
	}

	//获取当前的版本
	runtime, err := s.serviceRuntimeStore.GetForCluster(ctx, serviceId, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	t := time.Now()

	if runtime == nil {
		return errors.New("invalid version")
	}
	if !runtime.IsOnline {
		return errors.New("已下线不可重复下线")
	}
	runtime.IsOnline = false
	runtime.UpdateTime = t
	runtime.Operator = operator

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        serviceInfo.UUID,
		Name:        serviceName,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
		PublishType: 2,
	})

	return s.serviceStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = s.serviceRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}
		client, err := s.apintoClient.GetClient(ctx, clusterInfo.Id)
		if err != nil {
			return err
		}
		return common.CheckWorkerNotExist(client.ForService().Delete(strings.ToLower(serviceInfo.Name)))
	})
}

func (s *service) GetServiceIDByName(ctx context.Context, namespaceId int, serviceName string) (int, error) {
	info, err := s.serviceStore.GetByName(ctx, namespaceId, serviceName)
	if err != nil {
		return 0, err
	}

	return info.Id, nil
}

func (s *service) GetLatestServiceVersion(ctx context.Context, serviceID int) (*upstream_entry2.ServiceVersion, error) {
	stat, err := s.serviceStatStore.Get(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	version, err := s.serviceVersionStore.Get(ctx, stat.VersionId)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (s *service) GetServiceSchemaInfo(ctx context.Context, serviceID int) (*upstream_entry2.Service, error) {
	return s.serviceStore.Get(ctx, serviceID)
}

func (s *service) IsOnline(ctx context.Context, clusterId, serviceId int) bool {
	runtime, err := s.serviceRuntimeStore.GetForCluster(ctx, serviceId, clusterId)
	if err != nil {
		return false
	}
	return runtime.IsOnline
}

func (s *service) GetServiceRemoteOptions(ctx context.Context, namespaceID, pageNum, pageSize int, keyword string) ([]*strategy_model.RemoteServices, int, error) {
	sl, total, err := s.serviceStore.GetListPage(ctx, namespaceID, keyword, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*strategy_model.RemoteServices, 0, len(sl))
	for _, item := range sl {
		info, err := s.GetServiceInfo(ctx, namespaceID, item.Name)
		if err != nil {
			return nil, 0, err
		}

		li := &strategy_model.RemoteServices{
			Uuid:   item.Name, //TODO
			Name:   item.Name,
			Scheme: info.Scheme,
			Desc:   item.Desc,
		}

		list = append(list, li)
	}

	return list, total, nil
}

func (s *service) GetServiceRemoteByNames(ctx context.Context, namespaceID int, names []string) ([]*strategy_model.RemoteServices, error) {
	sl, err := s.serviceStore.GetByNames(ctx, namespaceID, names)
	if err != nil {
		return nil, err
	}

	list := make([]*strategy_model.RemoteServices, 0, len(sl))
	for _, item := range sl {
		info, err := s.GetServiceInfo(ctx, namespaceID, item.Name)
		if err != nil {
			return nil, err
		}

		li := &strategy_model.RemoteServices{
			Uuid:   item.Name, //TODO
			Name:   item.Name,
			Scheme: info.Scheme,
			Desc:   item.Desc,
		}

		list = append(list, li)
	}

	return list, nil
}
