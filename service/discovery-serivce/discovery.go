package discovery_serivce

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	driver_manager "github.com/eolinker/apinto-dashboard/driver-manager"
	"github.com/eolinker/apinto-dashboard/driver-manager/driver"
	"github.com/eolinker/apinto-dashboard/dto/discover-dto"
	"github.com/eolinker/apinto-dashboard/entry/discovery-entry"
	"github.com/eolinker/apinto-dashboard/entry/quote-entry"
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
	"github.com/eolinker/apinto-dashboard/model/audit-model"
	"github.com/eolinker/apinto-dashboard/model/cluster-model"
	"github.com/eolinker/apinto-dashboard/model/discovery-model"
	"github.com/eolinker/apinto-dashboard/model/frontend-model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"

	"github.com/eolinker/apinto-dashboard/service/cluster-service"
	"github.com/eolinker/apinto-dashboard/service/locker-service"
	"github.com/eolinker/apinto-dashboard/service/namespace-service"
	"github.com/eolinker/apinto-dashboard/service/user-service"
	"github.com/eolinker/apinto-dashboard/service/variable-service"
	"github.com/eolinker/apinto-dashboard/store/discovery-store"
	"github.com/eolinker/apinto-dashboard/store/quote-store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type IDiscoveryService interface {
	GetDiscoveryList(ctx context.Context, namespaceID int, searchName string) ([]*discovery_model.DiscoveryListItem, error)
	GetDiscoveryVersionInfo(ctx context.Context, namespaceID int, discoveryName string) (*discovery_model.DiscoveryInfo, error)
	CreateDiscovery(ctx context.Context, namespaceID int, userID int, input *discover_dto.DiscoveryInfoProxy) error
	UpdateDiscovery(ctx context.Context, namespaceID int, userID int, input *discover_dto.DiscoveryInfoProxy) error
	DeleteDiscovery(ctx context.Context, namespaceId, userId int, discoveryName string) error

	OnlineList(ctx context.Context, namespaceId int, discoveryName string) ([]*discovery_model.DiscoveryOnline, error)
	OnlineDiscovery(ctx context.Context, namespaceId, operator int, discoveryName, clusterName string) (*frontend_model.Router, error)
	OfflineDiscovery(ctx context.Context, namespaceId, operator int, discoveryName, clusterName string) error

	GetDiscoveryName(ctx context.Context, discoveryID int) (string, error)
	GetDiscoveryInfoByID(ctx context.Context, discoveryID int) (*discovery_model.DiscoveryListItem, error)
	GetDiscoveryID(ctx context.Context, namespaceID int, discoveryName string) (int, error)
	GetDiscoveryEnum(ctx context.Context, namespaceID int) ([]*discovery_model.DiscoveryEnum, error)
	GetDriversRender() []*driver_manager.DriverInfo
	GetLatestDiscoveryVersion(ctx context.Context, discoveryID int) (*discovery_model.DiscoveryVersion, error)
	IsOnline(ctx context.Context, clusterId, discoveryId int) bool

	//通过服务名获取配置上游服务时所需要的discoveryDriver
	GetServiceDiscoveryDriver(ctx context.Context, namespaceID int, discoveryName string) (int, string, driver.IServiceDriver, error)
	GetServiceDiscoveryDriverByID(ctx context.Context, discoveryID int) (string, string, driver.IServiceDriver, error)
	cluster_service.IResetOnlineService
}

type discoveryService struct {
	discoveryStore        discovery_store.IDiscoveryStore
	discoveryHistoryStore discovery_store.IDiscoveryHistoryStore
	discoveryStatStore    discovery_store.IDiscoveryStatStore
	discoveryVersionStore discovery_store.IDiscoveryVersionStore
	discoveryRuntimeStore discovery_store.IDiscoveryRuntimeStore
	quoteStore            quote_store.IQuoteStore

	clusterNodeService    cluster_service.IClusterNodeService
	apintoClient          cluster_service.IApintoClient
	userInfoService       user_service.IUserInfoService
	namespaceService      namespace_service.INamespaceService
	globalVariableService variable_service.IGlobalVariableService
	clusterService        cluster_service.IClusterService
	variableService       variable_service.IClusterVariableService
	service               upstream.IService

	lockService      locker_service.IAsynLockService
	discoveryManager driver_manager.IDiscoveryDriverManager
	staticDriver     driver.IServiceDriver
}

func newDiscoveryService() IDiscoveryService {
	s := &discoveryService{}
	bean.Autowired(&s.discoveryStore)
	bean.Autowired(&s.discoveryStatStore)
	bean.Autowired(&s.discoveryVersionStore)
	bean.Autowired(&s.discoveryRuntimeStore)
	bean.Autowired(&s.discoveryHistoryStore)

	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.globalVariableService)

	bean.Autowired(&s.service)
	bean.Autowired(&s.userInfoService)
	bean.Autowired(&s.clusterNodeService)
	bean.Autowired(&s.discoveryManager)
	bean.Autowired(&s.staticDriver)
	bean.Autowired(&s.quoteStore)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.namespaceService)
	bean.Autowired(&s.variableService)
	bean.Autowired(&s.lockService)
	return s
}

func (d *discoveryService) GetDiscoveryList(ctx context.Context, namespaceID int, searchName string) ([]*discovery_model.DiscoveryListItem, error) {
	dl, err := d.discoveryStore.GetList(ctx, namespaceID, searchName)
	if err != nil {
		return nil, err
	}

	list := make([]*discovery_model.DiscoveryListItem, 0, len(dl))
	for _, discovery := range dl {
		item := &discovery_model.DiscoveryListItem{
			Name:       discovery.Name,
			UUID:       discovery.UUID,
			Driver:     discovery.Driver,
			Desc:       discovery.Desc,
			UpdateTime: discovery.UpdateTime,
			IsDelete:   false,
		}

		isDelete, _ := d.isDiscoveryCanDelete(ctx, namespaceID, discovery.Id)

		item.IsDelete = isDelete
		list = append(list, item)
	}

	return list, nil
}

func (d *discoveryService) GetDiscoveryVersionInfo(ctx context.Context, namespaceID int, discoveryName string) (*discovery_model.DiscoveryInfo, error) {
	discovery, err := d.discoveryStore.GetByName(ctx, namespaceID, discoveryName)
	if err != nil {
		return nil, err
	}
	version, err := d.GetLatestDiscoveryVersion(ctx, discovery.Id)
	if err != nil {
		return nil, err
	}

	discoveryDriver := d.discoveryManager.GetDriver(discovery.Driver)

	info := &discovery_model.DiscoveryInfo{
		Name:   discoveryName,
		UUID:   discovery.UUID,
		Driver: discovery.Driver,
		Desc:   discovery.Desc,
		Config: discoveryDriver.FormatConfig([]byte(version.Config)),
		Render: discoveryDriver.Render(),
	}

	return info, nil
}

func (d *discoveryService) CreateDiscovery(ctx context.Context, namespaceID int, userID int, input *discover_dto.DiscoveryInfoProxy) error {
	discoveryDriver := d.discoveryManager.GetDriver(input.Driver)
	if discoveryDriver == nil {
		return fmt.Errorf("Driver %s doesn't exit. ", input.Driver)
	}

	//服务发现name查重
	_, err := d.discoveryStore.GetByName(ctx, namespaceID, input.Name)
	if err != gorm.ErrRecordNotFound {
		if err == nil {
			return fmt.Errorf("name %s already exit. ", input.Name)
		}
		return err
	}

	newConf, _, variableList, err := discoveryDriver.CheckInput(input.Config)
	if err != nil {
		return err
	}
	input.Config = newConf

	if input.UUID == "" {
		input.UUID = uuid.New()
	}

	input.UUID = strings.ToLower(input.UUID)

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.Name,
	})

	return d.discoveryStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		discoveryInfo := &discovery_entry.Discovery{
			NamespaceId: namespaceID,
			Name:        input.Name,
			UUID:        input.UUID,
			Driver:      input.Driver,
			Desc:        input.Desc,
			Operator:    userID,
			CreateTime:  t,
			UpdateTime:  t,
		}
		if err = d.discoveryStore.Save(txCtx, discoveryInfo); err != nil {
			return err
		}
		//添加版本信息
		discoveryVersionInfo := &discovery_entry.DiscoveryVersion{
			DiscoveryID: discoveryInfo.Id,
			NamespaceID: namespaceID,
			DiscoveryVersionConfig: discovery_entry.DiscoveryVersionConfig{
				Config: input.Config.String(),
			},
			Operator:   userID,
			CreateTime: t,
		}

		if err = d.discoveryVersionStore.Save(txCtx, discoveryVersionInfo); err != nil {
			return err
		}
		stat := &discovery_entry.DiscoveryStat{
			DiscoveryID: discoveryInfo.Id,
			VersionID:   discoveryVersionInfo.Id,
		}
		//添加版本关联原表信息
		if err = d.discoveryStatStore.Save(txCtx, stat); err != nil {
			return err
		}
		//创建所有引用的环境变量
		if len(variableList) > 0 {
			variables, err := d.globalVariableService.GetByKeys(ctx, namespaceID, variableList)
			if err != nil {
				return err
			}
			quoteMap := make(map[quote_entry.QuoteTargetKindType][]int)
			for _, variable := range variables {
				quoteMap[quote_entry.QuoteTargetKindTypeVariable] = append(quoteMap[quote_entry.QuoteTargetKindTypeVariable], variable.Id)
			}

			if err = d.quoteStore.Set(txCtx, discoveryInfo.Id, quote_entry.QuoteKindTypeDiscovery, quoteMap); err != nil {
				return err
			}
		}

		return d.discoveryHistoryStore.HistoryAdd(txCtx, namespaceID, discoveryInfo.Id, &discovery_entry.DiscoveryHistoryInfo{
			Discovery: *discoveryInfo,
			Config: discovery_entry.DiscoveryVersionConfig{
				Config: input.Config.String(),
			},
		}, userID)
	})

}

func (d *discoveryService) UpdateDiscovery(ctx context.Context, namespaceID int, userID int, input *discover_dto.DiscoveryInfoProxy) error {
	discoveryInfo, err := d.discoveryStore.GetByName(ctx, namespaceID, input.Name)
	if err != nil {
		return err
	}

	discoveryDriver := d.discoveryManager.GetDriver(discoveryInfo.Driver)
	if discoveryDriver == nil {
		return fmt.Errorf("Driver %s doesn't exit. ", discoveryInfo.Driver)
	}
	newConf, _, variableList, err := discoveryDriver.CheckInput(input.Config)
	if err != nil {
		return err
	}
	input.Config = newConf

	latestVersion, err := d.GetLatestDiscoveryVersion(ctx, discoveryInfo.Id)
	if err != nil {
		return err
	}

	oldDiscovery := *discoveryInfo
	t := time.Now()

	discoveryInfo.Desc = input.Desc
	discoveryInfo.Operator = userID
	discoveryInfo.UpdateTime = t

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.Name,
	})

	return d.discoveryStore.Transaction(ctx, func(txCtx context.Context) error {

		//修改基础数据
		if _, err = d.discoveryStore.Update(txCtx, discoveryInfo); err != nil {
			return err
		}

		//若配置更新了才往version插入新数据
		if discoveryDriver.CheckConfIsChange([]byte(latestVersion.Config), input.Config) {
			discoveryVersionInfo := &discovery_entry.DiscoveryVersion{
				DiscoveryID: discoveryInfo.Id,
				NamespaceID: namespaceID,
				DiscoveryVersionConfig: discovery_entry.DiscoveryVersionConfig{
					Config: input.Config.String(),
				},
				Operator:   userID,
				CreateTime: t,
			}
			if err = d.discoveryVersionStore.Save(txCtx, discoveryVersionInfo); err != nil {
				return err
			}
			//添加版本关联原表信息
			stat := &discovery_entry.DiscoveryStat{
				DiscoveryID: discoveryInfo.Id,
				VersionID:   discoveryVersionInfo.Id,
			}
			if err = d.discoveryStatStore.Save(txCtx, stat); err != nil {
				return err
			}

			//更新引用， 获取新的引用变量ID
			targetMaps := make(map[quote_entry.QuoteTargetKindType][]int)
			variableIDList := make([]int, 0)
			variables, err := d.globalVariableService.GetByKeys(ctx, namespaceID, variableList)
			if err != nil {
				return err
			}
			for _, variable := range variables {
				variableIDList = append(variableIDList, variable.Id)
			}
			targetMaps[quote_entry.QuoteTargetKindTypeVariable] = variableIDList
			if err = d.quoteStore.Set(txCtx, discoveryInfo.Id, quote_entry.QuoteKindTypeDiscovery, targetMaps); err != nil {
				return err
			}
		}

		//操作记录
		return d.discoveryHistoryStore.HistoryEdit(txCtx, namespaceID, discoveryInfo.Id, &discovery_entry.DiscoveryHistoryInfo{
			Discovery: oldDiscovery,
			Config: discovery_entry.DiscoveryVersionConfig{
				Config: latestVersion.Config,
			},
		}, &discovery_entry.DiscoveryHistoryInfo{
			Discovery: *discoveryInfo,
			Config: discovery_entry.DiscoveryVersionConfig{
				Config: input.Config.String(),
			},
		}, userID)
	})

}

func (d *discoveryService) DeleteDiscovery(ctx context.Context, namespaceID, userId int, discoveryName string) error {
	discoveryInfo, err := d.discoveryStore.GetByName(ctx, namespaceID, discoveryName)
	if err != nil {
		return err
	}

	if err = d.lockService.Lock(locker_service.LockNameDiscovery, discoveryInfo.Id); err != nil {
		return err
	}
	defer d.lockService.Unlock(locker_service.LockNameDiscovery, discoveryInfo.Id)

	discoveryInfo, err = d.discoveryStore.GetByName(ctx, namespaceID, discoveryName)
	if err != nil {
		return err
	}

	_, err = d.isDiscoveryCanDelete(ctx, namespaceID, discoveryInfo.Id)
	if err != nil {
		return err
	}

	version, err := d.GetLatestDiscoveryVersion(ctx, discoveryInfo.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: discoveryInfo.UUID,
		Name: discoveryName,
	})

	err = d.discoveryStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = d.discoveryStore.Delete(txCtx, discoveryInfo.Id); err != nil {
			return err
		}

		if err = d.discoveryHistoryStore.HistoryDelete(txCtx, namespaceID, discoveryInfo.Id, &discovery_entry.DiscoveryHistoryInfo{
			Discovery: *discoveryInfo,
			Config:    discovery_entry.DiscoveryVersionConfig{Config: version.Config},
		}, userId); err != nil {
			return err
		}

		d.lockService.DeleteLock(locker_service.LockNameDiscovery, discoveryInfo.Id)
		delMap := make(map[string]interface{})
		delMap["`kind`"] = "discovery"
		delMap["`target`"] = discoveryInfo.Id

		if _, err = d.discoveryStatStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = d.discoveryVersionStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = d.discoveryRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		return d.quoteStore.DelBySource(txCtx, discoveryInfo.Id, quote_entry.QuoteKindTypeDiscovery)
	})
	if err != nil {
		return err
	}
	d.lockService.DeleteLock(locker_service.LockNameDiscovery, discoveryInfo.Id)
	return nil
}

func (d *discoveryService) OnlineList(ctx context.Context, namespaceId int, discoveryName string) ([]*discovery_model.DiscoveryOnline, error) {
	discoveryInfo, err := d.discoveryStore.GetByName(ctx, namespaceId, discoveryName)
	if err != nil {
		return nil, err
	}

	//获取工作空间下的所有集群
	clusters, err := d.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *cluster_model.Cluster) int {
		return t.Id
	})

	//获取当前服务发现下集群运行的版本
	runtimes, err := d.discoveryRuntimeStore.GetByTarget(ctx, discoveryInfo.Id)
	if err != nil {
		return nil, err
	}
	runtimeMaps := common.SliceToMap(runtimes, func(t *discovery_entry.DiscoveryRuntime) int {
		return t.ClusterId
	})

	list := make([]*discovery_model.DiscoveryOnline, 0, len(clusters))

	latestVersion, err := d.GetLatestDiscoveryVersion(ctx, discoveryInfo.Id)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(runtimes, func(t *discovery_entry.DiscoveryRuntime) int {
		return t.Operator
	})

	userInfoMaps, _ := d.userInfoService.GetUserInfoMaps(ctx, userIds...)

	for _, cluster := range clusterMaps {
		discoveryOnline := &discovery_model.DiscoveryOnline{
			ClusterName: cluster.Name,
			Env:         cluster.Env,
			Status:      1, //默认为未上线状态
		}
		if runtime, ok := runtimeMaps[cluster.Id]; ok {
			discoveryOnline.UpdateTime = runtime.UpdateTime
			if runtime.IsOnline {
				discoveryOnline.Status = 3
			} else {
				discoveryOnline.Status = 2
			}

			if userInfo, uOk := userInfoMaps[runtime.Operator]; uOk {
				discoveryOnline.Operator = userInfo.NickName
			}

			//已上线需要对比是否更新过 服务发现信息
			if discoveryOnline.Status == 3 && runtime.VersionID != latestVersion.Id {
				discoveryOnline.Status = 4
			}
		}

		list = append(list, discoveryOnline)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Status > list[j].Status
	})
	return list, nil
}

func (d *discoveryService) OnlineDiscovery(ctx context.Context, namespaceId, operator int, discoveryName, clusterName string) (*frontend_model.Router, error) {
	discoveryInfo, err := d.discoveryStore.GetByName(ctx, namespaceId, discoveryName)
	if err != nil {
		return nil, err
	}

	discoveryID := discoveryInfo.Id
	t := time.Now()

	//获取当前集群信息
	cluster, err := d.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	if err = d.lockService.Lock(locker_service.LockNameDiscovery, discoveryID); err != nil {
		return nil, err
	}
	defer d.lockService.Unlock(locker_service.LockNameDiscovery, discoveryID)

	discoveryInfo, err = d.discoveryStore.GetByName(ctx, namespaceId, discoveryName)
	if err != nil {
		return nil, err
	}

	latestVersion, err := d.GetLatestDiscoveryVersion(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	router := &frontend_model.Router{
		Name:   frontend_model.RouterNameClusterVariable,
		Params: make(map[string]string),
	}
	router.Params["cluster_name"] = clusterName
	//服务引用的环境变量
	quoteMaps, err := d.quoteStore.GetSourceQuote(ctx, discoveryID, quote_entry.QuoteKindTypeDiscovery)
	if err != nil {
		return nil, err
	}
	variableIds := quoteMaps[quote_entry.QuoteTargetKindTypeVariable]
	if len(variableIds) > 0 {
		//获取集群正在运行的环境变量版本
		variablePublishVersion, err := d.variableService.GetPublishVersion(ctx, cluster.Id)
		if err != nil || variablePublishVersion == nil {
			globalVariable, err := d.globalVariableService.GetById(ctx, variableIds[0])
			if err != nil {
				return nil, err
			}
			return router, errors.New(fmt.Sprintf("${%s}未上线到{%s}，上线/更新失败", globalVariable.Key, clusterName))
		}

		//已发布的环境变量
		toMap := common.SliceToMap(variablePublishVersion.ClusterVariable, func(t *variable_entry.ClusterVariable) int {
			return t.VariableId
		})

		for _, variableId := range variableIds {
			if _, ok := toMap[variableId]; !ok {
				globalVariable, err := d.globalVariableService.GetById(ctx, variableId)
				if err != nil {
					return nil, err
				}
				return router, errors.New(fmt.Sprintf("${%s}未上线到{%s}，上线/更新失败", globalVariable.Key, clusterName))
			}
		}
	}

	//获取apinto client
	client, err := d.apintoClient.GetClient(ctx, cluster.Id)
	if err != nil {
		return nil, err
	}

	//获取当前运行的版本
	runtime, err := d.discoveryRuntimeStore.GetForCluster(ctx, discoveryID, cluster.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	namespace, err := d.namespaceService.GetById(namespaceId)
	if err != nil {
		return nil, err
	}

	//判断是否是更新
	isApintoUpdate := false
	if runtime == nil {
		runtime = &discovery_entry.DiscoveryRuntime{
			NamespaceId: namespaceId,
			DiscoveryID: discoveryID,
			ClusterId:   cluster.Id,
			VersionID:   latestVersion.Id,
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
		runtime.VersionID = latestVersion.Id
		runtime.Operator = operator
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        discoveryInfo.UUID,
		Name:        discoveryName,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
		PublishType: 1,
	})

	//事务
	err = d.discoveryStore.Transaction(ctx, func(txCtx context.Context) error {

		driverManager := d.discoveryManager.GetDriver(discoveryInfo.Driver)

		discoveryConfig := driverManager.ToApinto(namespace.Name, strings.ToLower(discoveryInfo.Name), discoveryInfo.Desc, []byte(latestVersion.Config))

		if err = d.discoveryRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}

		if isApintoUpdate {
			return client.ForDiscovery().Update(strings.ToLower(discoveryInfo.Name), *discoveryConfig)
		} else {
			return client.ForDiscovery().Create(*discoveryConfig)
		}
	})

	return nil, err
}

func (d *discoveryService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	runtimes, err := d.discoveryRuntimeStore.GetByCluster(ctx, clusterId)
	if err != nil {
		log.Errorf("discoveryService-ResetOnline-getRuntimes clusterId=%d err=%s", clusterId, err.Error())
		return
	}
	client, err := d.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		log.Errorf("discoveryService-ResetOnline-getClient clusterId=%d err=%s", clusterId, err.Error())
		return
	}

	namespace, err := d.namespaceService.GetById(namespaceId)
	if err != nil {
		log.Errorf("discoveryService-ResetOnline-getNamespace clusterId=%d err=%s", clusterId, err.Error())
		return
	}

	for _, runtime := range runtimes {
		if !runtime.IsOnline {
			continue
		}

		discoveryInfo, err := d.discoveryStore.Get(ctx, runtime.DiscoveryID)
		if err != nil {
			log.Errorf("discoveryService-ResetOnline-getDiscovery clusterId=%d discoveryId=%d err=%s", clusterId, runtime.DiscoveryID, err.Error())
			continue
		}

		version, err := d.discoveryVersionStore.Get(ctx, runtime.VersionID)
		if err != nil {
			log.Errorf("discoveryService-ResetOnline-getVersion clusterId=%d versionId=%d err=%s", clusterId, runtime.VersionID, err.Error())
			continue
		}
		driverManager := d.discoveryManager.GetDriver(discoveryInfo.Driver)

		discoveryConfig := driverManager.ToApinto(namespace.Name, strings.ToLower(discoveryInfo.Name), discoveryInfo.Desc, []byte(version.Config))

		if err = client.ForDiscovery().Create(*discoveryConfig); err != nil {
			log.Errorf("discoveryService-ResetOnline-apinto clusterId=%d discoveryConfig=%v err=%s", clusterId, discoveryConfig, err.Error())
			continue
		}
	}
}

func (d *discoveryService) OfflineDiscovery(ctx context.Context, namespaceId, operator int, discoveryName, clusterName string) error {
	discoveryInfo, err := d.discoveryStore.GetByName(ctx, namespaceId, discoveryName)
	if err != nil {
		return err
	}

	//获取当前集群信息
	cluster, err := d.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	if err = d.lockService.Lock(locker_service.LockNameDiscovery, discoveryInfo.Id); err != nil {
		return err
	}
	defer d.lockService.Unlock(locker_service.LockNameDiscovery, discoveryInfo.Id)

	discoveryInfo, err = d.discoveryStore.GetByName(ctx, namespaceId, discoveryName)
	if err != nil {
		return err
	}

	//获取当前的版本
	runtime, err := d.discoveryRuntimeStore.GetForCluster(ctx, discoveryInfo.Id, cluster.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}

	//下线前确定引用了此服务发现的上游服务均已下线
	quoteSet, err := d.quoteStore.GetTargetQuote(ctx, discoveryInfo.Id, quote_entry.QuoteTargetKindTypeDiscovery)
	if err != nil {
		return err
	}
	serviceIds := quoteSet[quote_entry.QuoteKindTypeService]
	for _, serviceID := range serviceIds {
		if d.service.IsOnline(ctx, cluster.Id, serviceID) {
			info, err := d.service.GetServiceSchemaInfo(ctx, serviceID)
			if err != nil {
				return err
			}
			return fmt.Errorf("service %s is already online. ", info.Name)
		}
	}

	t := time.Now()

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        discoveryInfo.UUID,
		Name:        discoveryName,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
		PublishType: 2,
	})

	//事务
	return d.discoveryStore.Transaction(ctx, func(txCtx context.Context) error {
		if !runtime.IsOnline {
			return errors.New("已下线不可重复下线")
		}
		runtime.IsOnline = false
		runtime.UpdateTime = t
		runtime.Operator = operator
		err = d.discoveryRuntimeStore.Save(txCtx, runtime)
		if err != nil {
			return err
		}

		//发布到apinto
		client, err := d.apintoClient.GetClient(ctx, cluster.Id)
		if err != nil {
			return err
		}
		return common.CheckWorkerNotExist(client.ForDiscovery().Delete(strings.ToLower(discoveryName) + "@discovery"))
	})
}

func (d *discoveryService) GetDiscoveryName(ctx context.Context, discoveryID int) (string, error) {
	info, err := d.discoveryStore.Get(ctx, discoveryID)
	if err != nil {
		return "", err
	}

	return info.Name, nil
}
func (d *discoveryService) GetDiscoveryInfoByID(ctx context.Context, discoveryID int) (*discovery_model.DiscoveryListItem, error) {
	discovery, err := d.discoveryStore.Get(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	info := &discovery_model.DiscoveryListItem{
		Name:       discovery.Name,
		UUID:       discovery.UUID,
		Driver:     discovery.Driver,
		Desc:       discovery.Desc,
		UpdateTime: discovery.UpdateTime,
		IsDelete:   false,
	}

	return info, nil
}

func (d *discoveryService) GetDiscoveryID(ctx context.Context, namespaceID int, discoveryName string) (int, error) {
	//判断是静态服务发现则返回id 0
	if discoveryName == driver.DriverStatic {
		return 0, nil
	}

	discoveryInfo, err := d.discoveryStore.GetByName(ctx, namespaceID, discoveryName)
	if err != nil {
		return 0, err
	}
	return discoveryInfo.Id, nil
}

func (d *discoveryService) GetDiscoveryEnum(ctx context.Context, namespaceID int) ([]*discovery_model.DiscoveryEnum, error) {
	enums := make([]*discovery_model.DiscoveryEnum, 0)
	list, err := d.discoveryStore.List(ctx, map[string]interface{}{
		"namespace": namespaceID,
	})
	if err != nil {
		return nil, err
	}

	//静态节点驱动
	staticEnum := &discovery_model.DiscoveryEnum{
		Name:   driver.DriverStatic,
		Driver: driver.DriverStatic,
		Render: d.staticDriver.Render(),
	}
	enums = append(enums, staticEnum)

	for _, discovery := range list {
		discoveryDriver := d.discoveryManager.GetDriver(discovery.Driver)
		if discoveryDriver == nil {
			continue
		}
		enum := &discovery_model.DiscoveryEnum{
			Name:   discovery.Name,
			Driver: discovery.Driver,
			Render: discoveryDriver.OptionsEnum().Render(),
		}
		enums = append(enums, enum)
	}
	return enums, nil
}

func (d *discoveryService) GetDriversRender() []*driver_manager.DriverInfo {
	return d.discoveryManager.List()
}

func (d *discoveryService) GetServiceDiscoveryDriverByID(ctx context.Context, discoveryID int) (string, string, driver.IServiceDriver, error) {
	//判断是static 则返回静态驱动
	if discoveryID == 0 {
		return driver.DriverStatic, driver.DriverStatic, d.staticDriver, nil
	}
	discovery, err := d.discoveryStore.Get(ctx, discoveryID)
	if err != nil {
		return "", "", nil, err
	}
	return discovery.Name, discovery.Driver, d.discoveryManager.GetDriver(discovery.Driver).OptionsEnum(), nil
}

func (d *discoveryService) GetLatestDiscoveryVersion(ctx context.Context, discoveryID int) (*discovery_model.DiscoveryVersion, error) {
	var err error
	stat, err := d.discoveryStatStore.Get(ctx, discoveryID)
	if err != nil {
		return nil, err
	}

	var version *discovery_entry.DiscoveryVersion

	version, err = d.discoveryVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return nil, err
	}

	return (*discovery_model.DiscoveryVersion)(version), nil
}

func (d *discoveryService) GetServiceDiscoveryDriver(ctx context.Context, namespaceID int, discoveryName string) (int, string, driver.IServiceDriver, error) {
	//判断是static 则返回静态驱动
	if discoveryName == driver.DriverStatic {
		return 0, driver.DriverStatic, d.staticDriver, nil
	}

	discovery, err := d.discoveryStore.GetByName(ctx, namespaceID, discoveryName)
	if err != nil {
		return 0, "", nil, err
	}
	return discovery.Id, discovery.Driver, d.discoveryManager.GetDriver(discovery.Driver).OptionsEnum(), nil
}

func (d *discoveryService) isDiscoveryCanDelete(ctx context.Context, namespaceID, discoveryID int) (bool, error) {
	quotedSet, err := d.quoteStore.GetTargetQuote(ctx, discoveryID, quote_entry.QuoteTargetKindTypeDiscovery)
	if err != nil {
		return false, err
	}
	for _, serviceID := range quotedSet[quote_entry.QuoteKindTypeService] {
		serviceInfo, err := d.service.GetServiceSchemaInfo(ctx, serviceID)
		if err != nil {
			return false, err
		}
		return false, fmt.Errorf("Discovery is in use by service %s. ", serviceInfo.Name)
	}

	clusters, err := d.clusterService.GetByNamespaceId(ctx, namespaceID)
	if err != nil {
		return false, err
	}
	for _, cluster := range clusters {
		runtime, _ := d.discoveryRuntimeStore.GetForCluster(ctx, discoveryID, cluster.Id)
		if runtime != nil && runtime.IsOnline {
			return false, fmt.Errorf("Discovery is online in cluster %s. ", cluster.Name)
		}
	}
	return true, nil

}

func (d *discoveryService) IsOnline(ctx context.Context, clusterId, discoveryId int) bool {
	runtime, err := d.discoveryRuntimeStore.GetForCluster(ctx, discoveryId, clusterId)
	if err != nil {
		return false
	}
	return runtime.IsOnline
}
