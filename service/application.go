package service

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

const anonymousIds = "anonymous"

var _ IApplicationService = (*applicationService)(nil)

type IApplicationService interface {
	CreateApp(ctx context.Context, namespaceId, userId int, input *dto.ApplicationInput) error
	UpdateApp(ctx context.Context, namespaceId, userId int, input *dto.ApplicationInput) error
	DelApp(ctx context.Context, namespaceId, userId int, id string) error
	AppList(ctx context.Context, namespaceId, userId, pageNum, pageSize int, queryName string) ([]*model.Application, int, error)
	AppListAll(ctx context.Context, namespaceId int) ([]*model.Application, error)
	AppListFilter(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*model.Application, int, error)
	AppListByUUIDS(ctx context.Context, namespaceId int, uuids []string) ([]*model.Application, error)
	AppInfoDetails(ctx context.Context, namespaceId int, id string) (*model.Application, error)
	AppInfo(ctx context.Context, namespaceId int, id string) (*model.Application, error)
	Online(ctx context.Context, namespaceId, userId int, id, clusterName string) error
	Offline(ctx context.Context, namespaceId, userId int, id, clusterName string) error
	Disable(ctx context.Context, namespaceId, userId int, id, clusterName string, disable bool) error
	OnlineList(ctx context.Context, namespaceId int, id string) ([]*model.ApplicationOnline, error)
	GetAppKeys(ctx context.Context, namespaceId int) ([]*model.ApplicationKeys, error)
	getAppVersion(ctx context.Context, appId int) (*model.ApplicationVersion, error)
	IResetOnlineService
}

type applicationService struct {
	applicationStore            store.IApplicationStore
	applicationRuntimeStore     store.IApplicationRuntimeStore
	applicationAuthRuntimeStore store.IApplicationAuthRuntimeStore
	applicationVersionStore     store.IApplicationVersionStore
	applicationStatStore        store.IApplicationStatStore
	applicationHistoryStore     store.IApplicationHistoryStore
	clusterService              IClusterService
	applicationAuthService      IApplicationAuthService
	randomService               IRandomService
	apintoClient                IApintoClient
	lockService                 IAsynLockService
	userInfoService             IUserInfoService
}

func newApplicationService() IApplicationService {
	app := &applicationService{}
	bean.Autowired(&app.applicationStore)
	bean.Autowired(&app.applicationRuntimeStore)
	bean.Autowired(&app.applicationAuthRuntimeStore)
	bean.Autowired(&app.applicationVersionStore)
	bean.Autowired(&app.applicationStatStore)
	bean.Autowired(&app.applicationHistoryStore)
	bean.Autowired(&app.randomService)
	bean.Autowired(&app.clusterService)
	bean.Autowired(&app.apintoClient)
	bean.Autowired(&app.applicationAuthService)
	bean.Autowired(&app.lockService)
	bean.Autowired(&app.userInfoService)
	return app
}

func (a *applicationService) OnlineList(ctx context.Context, namespaceId int, id string) ([]*model.ApplicationOnline, error) {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}
	applicationId := application.Id

	//获取工作空间下的所有集群
	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *model.Cluster) int {
		return t.Id
	})

	//获取当前服务下集群运行的版本
	runtimes, err := a.applicationRuntimeStore.GetByTarget(ctx, applicationId)
	if err != nil {
		return nil, err
	}
	//最新版本
	lastVersion, err := a.getAppVersion(ctx, application.Id)
	if err != nil {
		return nil, err
	}

	runtimeMaps := common.SliceToMap(runtimes, func(t *entry.ApplicationRuntime) int {
		return t.ClusterId
	})

	userIds := common.SliceToSliceIds(runtimes, func(t *entry.ApplicationRuntime) int {
		return t.Operator
	})

	userInfoMaps, _ := a.userInfoService.GetUserInfoMaps(ctx, userIds...)

	list := make([]*model.ApplicationOnline, 0, len(clusters))
	for _, cluster := range clusterMaps {

		applicationOnline := &model.ApplicationOnline{
			ClusterID:   cluster.Id,
			ClusterName: cluster.Name,
			Env:         cluster.Env,
			Status:      1, //默认为未上线状态
		}

		if runtime, ok := runtimeMaps[cluster.Id]; ok {
			applicationOnline.Disable = runtime.Disable
			if runtime.IsOnline {
				applicationOnline.Status = 3
			} else {
				applicationOnline.Status = 2
			}
			applicationOnline.UpdateTime = runtime.UpdateTime

			if userInfo, uOk := userInfoMaps[runtime.Operator]; uOk {
				applicationOnline.Operator = userInfo.NickName
			}

			if applicationOnline.Status == 3 {
				currentVersion, err := a.applicationVersionStore.Get(ctx, runtime.VersionId)
				if err != nil {
					return nil, err
				}

				if currentVersion.Id != lastVersion.Id {
					applicationOnline.Status = 4
				}

				if applicationOnline.Status == 3 {
					isUpdate, err := a.applicationAuthService.isUpdate(ctx, cluster.Id, currentVersion.ApplicationID)
					if err != nil {
						return nil, err
					}
					if isUpdate {
						applicationOnline.Status = 4
					}
				}

			}
		} else {
			if application.IdStr == anonymousIds {
				applicationOnline.Disable = true
			}
		}

		list = append(list, applicationOnline)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Status > list[j].Status
	})
	return list, nil
}

func (a *applicationService) Online(ctx context.Context, namespaceId, userId int, id, clusterName string) error {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}
	//除了匿名应用以外，其他应用需要配置鉴权信息才可上线
	anonymous := true
	if application.IdStr != anonymousIds {
		auths, err := a.applicationAuthService.getListByApplicationId(ctx, application.Id)
		if err != nil {
			return err
		}
		if len(auths) == 0 {
			return errors.New("需要配置鉴权信息才可上线")
		}
		anonymous = false
	}
	//获取当前集群信息
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	applicationId := application.Id
	clusterId := cluster.Id

	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(LockNameApplication, applicationId); err != nil {
		return err
	}
	defer a.lockService.Unlock(LockNameApplication, applicationId)

	//拿到锁后需要重新获取下信息
	application, err = a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	//获取当前应用的版本
	lastVersion, err := a.getAppVersion(ctx, application.Id)
	if err != nil {
		return err
	}

	runtime, err := a.applicationRuntimeStore.GetForCluster(ctx, applicationId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	t := time.Now()
	if runtime == nil {
		runtime = &entry.ApplicationRuntime{
			NamespaceId:   namespaceId,
			ApplicationId: applicationId,
			ClusterId:     clusterId,
			VersionId:     lastVersion.Id,
			IsOnline:      true,
			Operator:      userId,
			CreateTime:    t,
			UpdateTime:    t,
		}
		if anonymous {
			runtime.Disable = true
		}
	} else {
		runtime.IsOnline = true
		runtime.UpdateTime = t
		runtime.VersionId = lastVersion.Id
		runtime.Operator = userId
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        id,
		Name:        application.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
		PublishType: 1,
	})

	return a.applicationRuntimeStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = a.applicationRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}

		auths := make([]v1.ApplicationAuth, 0)

		if !anonymous {
			//上线鉴权信息
			authList, err := a.applicationAuthService.online(txCtx, namespaceId, userId, clusterId, applicationId)
			if err != nil {
				return err
			}
			for _, auth := range authList {
				auths = append(auths, a.applicationAuthService.GetDriver(auth.Driver).ToApinto(auth.ExpireTime, auth.Position, auth.TokenName, []byte(auth.Config), auth.IsTransparent))
			}
		}

		labels := make(map[string]string)
		for _, attr := range lastVersion.CustomAttrList {
			labels[attr.Key] = attr.Value
		}
		appConfig := &v1.ApplicationConfig{
			Name:        application.IdStr,
			Driver:      "app",
			Auth:        auths,
			Disable:     runtime.Disable,
			Description: application.Desc,
			Labels:      labels,
			Additional:  a.getApplicationAdditional(lastVersion.ExtraParamList),
			Anonymous:   anonymous,
		}

		if runtime.Id > 0 {
			return client.ForApp().Update(application.IdStr, *appConfig)
		}
		return client.ForApp().Create(*appConfig)
	})
}

func (a *applicationService) Offline(ctx context.Context, namespaceId, userId int, id, clusterName string) error {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	//获取当前集群信息
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	applicationId := application.Id
	clusterId := cluster.Id

	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(LockNameApplication, applicationId); err != nil {
		return err
	}
	defer a.lockService.Unlock(LockNameApplication, applicationId)

	//拿到锁后需要重新获取下信息
	application, err = a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	runtime, err := a.applicationRuntimeStore.GetForCluster(ctx, applicationId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}
	if !runtime.IsOnline {
		return errors.New("已下线不可重复下线")
	}

	runtime.IsOnline = false
	runtime.UpdateTime = time.Now()
	runtime.Operator = userId

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        id,
		Name:        application.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
		PublishType: 2,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = a.applicationRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}
		//鉴权信息下线
		if err = a.applicationAuthService.offline(txCtx, clusterId, applicationId); err != nil {
			return err
		}

		return common.CheckWorkerNotExist(client.ForApp().Delete(application.IdStr))
	})
}

func (a *applicationService) Disable(ctx context.Context, namespaceId, userId int, id, clusterName string, disable bool) error {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	//获取当前集群信息
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	applicationId := application.Id
	clusterId := cluster.Id

	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(LockNameApplication, applicationId); err != nil {
		return err
	}
	defer a.lockService.Unlock(LockNameApplication, applicationId)

	//拿到锁后需要重新获取下信息
	application, err = a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	runtime, err := a.applicationRuntimeStore.GetForCluster(ctx, applicationId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if runtime == nil {
		return errors.New("未上线应用不可启用/禁用")
	}
	if runtime.Disable == disable {
		if !disable {
			return errors.New("已是启用状态无需重复启用")
		}
		return errors.New("已是禁用状态无需重复禁用")
	}

	runtime.Disable = disable
	runtime.UpdateTime = time.Now()
	runtime.Operator = userId

	//编写日志操作对象信息
	enableOperate := 1
	if disable {
		enableOperate = 2
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:          id,
		Name:          application.Name,
		ClusterId:     cluster.Id,
		ClusterName:   clusterName,
		EnableOperate: enableOperate,
	})

	return a.applicationRuntimeStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = a.applicationRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}
		if runtime.IsOnline { //在线状态需要把禁用状态更新到apinto
			if err = client.ForApp().Patch(application.IdStr, map[string]interface{}{"disable": disable}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (a *applicationService) CreateApp(ctx context.Context, namespaceId, userId int, input *dto.ApplicationInput) error {
	input.Id = strings.ToLower(input.Id)
	application, _ := a.applicationStore.GetByIdStr(ctx, namespaceId, input.Id)
	if application != nil {
		return errors.New("应用ID重复")
	}

	application, _ = a.applicationStore.GetByName(ctx, namespaceId, input.Name)
	if application != nil {
		return errors.New("应用名重复")
	}

	versionConfig := entry.ApplicationVersionConfig{
		CustomAttrList: a.dtoAttrToEntryAttr(input.CustomAttrList),
		ExtraParamList: a.dtoExtraToEntryExtra(input.ExtraParamList),
		Apis:           input.Apis,
	}
	t := time.Now()

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: input.Id,
		Name: input.Name,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
		application = &entry.Application{
			NamespaceId: namespaceId,
			IdStr:       input.Id,
			Name:        input.Name,
			Desc:        input.Desc,
			Operator:    userId,
			CreateTime:  t,
			UpdateTime:  t,
		}

		if err := a.applicationStore.Save(txCtx, application); err != nil {
			return err
		}

		if err := a.applicationHistoryStore.HistoryAdd(txCtx, namespaceId, application.Id, &entry.ApplicationHistoryInfo{
			Application:              *application,
			ApplicationVersionConfig: versionConfig,
		}, userId); err != nil {
			return nil
		}

		applicationVersion := &entry.ApplicationVersion{
			ApplicationID:            application.Id,
			NamespaceID:              namespaceId,
			ApplicationVersionConfig: versionConfig,
			Operator:                 userId,
			CreateTime:               t,
		}

		if err := a.applicationVersionStore.Save(txCtx, applicationVersion); err != nil {
			return err
		}
		stat := &entry.ApplicationStat{
			ApplicationID: applicationVersion.ApplicationID,
			VersionID:     applicationVersion.Id,
		}

		return a.applicationStatStore.Save(txCtx, stat)
	})

}

func (a *applicationService) UpdateApp(ctx context.Context, namespaceId, userId int, input *dto.ApplicationInput) error {
	application, _ := a.applicationStore.GetByName(ctx, namespaceId, input.Name)
	if application != nil && application.IdStr != input.Id {
		return errors.New("应用名重复")
	}

	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, input.Id)
	if err != nil {
		return err
	}

	if application.IdStr == anonymousIds && input.Name != "匿名应用" {
		return errors.New("不能更改匿名应用的应用名")
	}

	//获取应用当前版本
	version, err := a.getAppVersion(ctx, application.Id)
	if err != nil {
		return err
	}

	isUpdateVersion := false
	oldAttrMaps := make(map[string]string)
	for _, attr := range version.CustomAttrList {
		oldAttrMaps[attr.Key] = attr.Value
	}
	newAttrMaps := make(map[string]string)
	for _, attr := range input.CustomAttrList {
		newAttrMaps[attr.Key] = attr.Value
	}

	if !common.DiffMap(oldAttrMaps, newAttrMaps) {
		isUpdateVersion = true
	}

	oldExtraMaps := make(map[string]string)
	for _, extra := range version.ExtraParamList {
		oldExtraMaps[extra.Key] = extra.Value
	}
	newExtraMaps := make(map[string]string)
	for _, extra := range input.ExtraParamList {
		newExtraMaps[extra.Key] = extra.Value
	}

	if !common.DiffMap(oldExtraMaps, newExtraMaps) {
		isUpdateVersion = true
	}

	if !slices.Equal(version.Apis, input.Apis) {
		isUpdateVersion = true
	}

	t := time.Now()
	//添加操作记录

	oldApplication := *application

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: input.Id,
		Name: input.Name,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
		application.UpdateTime = t
		application.Operator = userId
		application.Desc = input.Desc
		application.Name = input.Name

		if err = a.applicationStore.Save(txCtx, application); err != nil {
			return err
		}

		config := entry.ApplicationVersionConfig{
			CustomAttrList: a.dtoAttrToEntryAttr(input.CustomAttrList),
			ExtraParamList: a.dtoExtraToEntryExtra(input.ExtraParamList),
			Apis:           input.Apis,
		}

		applicationVersion := &entry.ApplicationVersion{
			ApplicationID:            application.Id,
			NamespaceID:              namespaceId,
			ApplicationVersionConfig: config,
			Operator:                 userId,
			CreateTime:               t,
		}

		if err = a.applicationHistoryStore.HistoryEdit(txCtx, namespaceId, application.Id, &entry.ApplicationHistoryInfo{
			Application: oldApplication,
			ApplicationVersionConfig: entry.ApplicationVersionConfig{
				CustomAttrList: version.CustomAttrList,
				ExtraParamList: version.ExtraParamList,
				Apis:           version.Apis,
			},
		}, &entry.ApplicationHistoryInfo{
			Application:              *application,
			ApplicationVersionConfig: config,
		}, userId); err != nil {
			return nil
		}

		if isUpdateVersion {
			if err = a.applicationVersionStore.Save(txCtx, applicationVersion); err != nil {
				return err
			}
			stat := &entry.ApplicationStat{
				ApplicationID: applicationVersion.ApplicationID,
				VersionID:     applicationVersion.Id,
			}
			return a.applicationStatStore.Save(txCtx, stat)
		}
		return nil
	})
}

func (a *applicationService) DelApp(ctx context.Context, namespaceId, userId int, id string) error {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}
	if application.IdStr == anonymousIds {
		return errors.New("匿名应用不能删除")
	}

	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}
	for _, cluster := range clusters {
		runtime, _ := a.applicationRuntimeStore.GetForCluster(ctx, application.Id, cluster.Id)
		if runtime != nil && runtime.IsOnline {
			return errors.New("应用已上线不可删除")
		}
	}

	//获取应用当前版本信息
	version, err := a.getAppVersion(ctx, application.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: id,
		Name: application.Name,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = a.applicationStore.Delete(txCtx, application.Id); err != nil {
			return err
		}

		//添加操作记录

		if err = a.applicationHistoryStore.HistoryDelete(txCtx, namespaceId, application.Id, &entry.ApplicationHistoryInfo{
			Application: *application,
			ApplicationVersionConfig: entry.ApplicationVersionConfig{
				CustomAttrList: version.CustomAttrList,
				ExtraParamList: version.ExtraParamList,
			},
		}, userId); err != nil {
			return nil
		}

		delMap := make(map[string]interface{})
		delMap["`kind`"] = "application"
		delMap["`target`"] = application.Id

		if _, err = a.applicationStatStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}
		if _, err = a.applicationVersionStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}
		for _, cluster := range clusters {
			delMap["`cluster`"] = cluster.Id
			if _, err = a.applicationRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *applicationService) AppList(ctx context.Context, namespaceId, userId, pageNum, pageSize int, queryName string) ([]*model.Application, int, error) {

	anonymousApplication, err := a.applicationStore.GetByIdStr(ctx, namespaceId, anonymousIds)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	//没有匿名应用创建一个
	if anonymousApplication == nil {
		t := time.Now()
		entryApplication := &entry.Application{
			NamespaceId: namespaceId,
			IdStr:       anonymousIds,
			Name:        "匿名应用",
			Operator:    userId,
			CreateTime:  t,
			UpdateTime:  t,
		}

		err = a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
			if err = a.applicationStore.Save(txCtx, entryApplication); err != nil {
				return err
			}
			version := &entry.ApplicationVersion{
				ApplicationID:            entryApplication.Id,
				NamespaceID:              namespaceId,
				ApplicationVersionConfig: entry.ApplicationVersionConfig{},
				Operator:                 userId,
				CreateTime:               t,
			}

			if err = a.applicationVersionStore.Save(txCtx, version); err != nil {
				return err
			}
			return a.applicationStatStore.Save(txCtx, &entry.ApplicationStat{
				ApplicationID: version.ApplicationID,
				VersionID:     version.Id,
			})
		})

		if err != nil {
			return nil, 0, err
		}

	}
	list, count, err := a.applicationStore.GetListPage(ctx, namespaceId, pageNum, pageSize, queryName)
	if err != nil {
		return nil, 0, err
	}

	applications := make([]*model.Application, 0, len(list))
	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, 0, err
	}

	userIds := common.SliceToSliceIds(list, func(t *entry.Application) int {
		return t.Operator
	})

	userInfoMaps, _ := a.userInfoService.GetUserInfoMaps(ctx, userIds...)

	for _, application := range list {

		operatorName := ""
		if userInfo, ok := userInfoMaps[application.Operator]; ok {
			operatorName = userInfo.NickName
		}

		val := &model.Application{Application: application, OperatorName: operatorName}

		isDelete := true
		if val.IdStr == anonymousIds {
			isDelete = false
		} else {
			for _, cluster := range clusters {
				runtime, _ := a.applicationRuntimeStore.GetForCluster(ctx, application.Id, cluster.Id)
				if runtime != nil && runtime.IsOnline {
					isDelete = false
					break
				}
			}
		}

		val.IsDelete = isDelete

		applications = append(applications, val)
	}

	//对列表进行排序，匿名排第一位，其余按更新时间降序
	sort.Sort(model.ApplicationList(applications))

	return applications, count, nil
}
func (a *applicationService) AppListAll(ctx context.Context, namespaceId int) ([]*model.Application, error) {

	list, err := a.applicationStore.GetList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	applications := make([]*model.Application, 0, len(list))

	for _, application := range list {
		applications = append(applications, &model.Application{Application: application})
	}

	sort.Sort(model.ApplicationList(applications))

	return applications, nil
}

func (a *applicationService) AppListFilter(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*model.Application, int, error) {

	list, count, err := a.applicationStore.GetListPage(ctx, namespaceId, pageNum, pageSize, queryName)
	if err != nil {
		return nil, 0, err
	}

	applications := make([]*model.Application, 0, len(list))

	for _, application := range list {
		val := &model.Application{Application: application}
		applications = append(applications, val)
	}

	return applications, count, nil
}

func (a *applicationService) AppInfoDetails(ctx context.Context, namespaceId int, id string) (*model.Application, error) {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}

	version, err := a.getAppVersion(ctx, application.Id)
	if err != nil {
		return nil, err
	}

	res := &model.Application{
		Application:  application,
		OperatorName: "",
		CustomAttr:   a.entryAttrToModelAttr(version.CustomAttrList),
		ExtraParam:   a.entryExtraToModelExtra(version.ExtraParamList),
	}
	return res, nil
}

func (a *applicationService) AppInfo(ctx context.Context, namespaceId int, id string) (*model.Application, error) {
	application, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}

	res := &model.Application{
		Application: application,
	}
	return res, nil
}

func (a *applicationService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	runtimes, err := a.applicationRuntimeStore.GetByCluster(ctx, clusterId)
	if err != nil {
		log.Errorf("applicationService-ResetOnline-getRuntimes clusterId=%d, err=%s", clusterId, err.Error())
		return
	}

	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		log.Errorf("applicationService-ResetOnline-getClient clusterId=%d, err=%s", clusterId, err.Error())
		return
	}

	for _, runtime := range runtimes {
		if !runtime.IsOnline {
			continue
		}

		application, err := a.applicationStore.Get(ctx, runtime.ApplicationId)
		if err != nil {
			log.Errorf("applicationService-ResetOnline-getApplication appId=%d,clusterId=%d, err=%s", runtime.ApplicationId, clusterId, err.Error())
			continue
		}

		version, err := a.applicationVersionStore.Get(ctx, runtime.VersionId)
		if err != nil {
			log.Errorf("applicationService-ResetOnline-getVersion versionId=%d,clusterId=%d, err=%s", runtime.VersionId, clusterId, err.Error())
			continue
		}

		//上线鉴权信息
		authList, err := a.applicationAuthService.GetList(ctx, namespaceId, application.IdStr)
		if err != nil {
			log.Errorf("applicationService-ResetOnline-getAuthList appIds=%s, err=%s", application.IdStr, err.Error())
			continue
		}
		auths := make([]v1.ApplicationAuth, 0)
		for _, auth := range authList {
			auths = append(auths, a.applicationAuthService.GetDriver(auth.Driver).ToApinto(auth.ExpireTime, auth.Position, auth.TokenName, []byte(auth.Config), auth.IsTransparent))
		}

		labels := make(map[string]string)
		for _, attr := range version.CustomAttrList {
			labels[attr.Key] = attr.Value
		}

		appConfig := &v1.ApplicationConfig{
			Name:        application.IdStr,
			Driver:      "app",
			Auth:        auths,
			Labels:      labels,
			Description: application.Desc,
			Additional:  a.getApplicationAdditional(version.ExtraParamList),
		}
		if err = client.ForApp().Create(*appConfig); err != nil {
			log.Errorf("applicationService-ResetOnline-apinto appConfig=%v, clusterId=%d  err=%s", appConfig, clusterId, err.Error())
			continue
		}
	}
}

func (a *applicationService) getAppVersion(ctx context.Context, appId int) (*model.ApplicationVersion, error) {
	var err error

	stat, err := a.applicationStatStore.Get(ctx, appId)
	if err != nil {
		return nil, err
	}

	var version *entry.ApplicationVersion

	version, err = a.applicationVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return nil, err
	}

	return (*model.ApplicationVersion)(version), nil
}

func (a *applicationService) GetAppKeys(ctx context.Context, namespaceId int) ([]*model.ApplicationKeys, error) {
	applications, err := a.applicationStore.GetList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	list := make([]*model.ApplicationKeys, 0)

	keys := common.Map[string, []string]{}

	for _, application := range applications {

		version, err := a.getAppVersion(ctx, application.Id)
		if err != nil {
			return nil, err
		}

		for _, val := range version.CustomAttrList {
			keys[val.Key] = append(keys[val.Key], val.Value)
		}

	}

	if len(keys) == 0 {
		return nil, err
	}

	for k, v := range keys {

		newValues := make([]string, 0)
		newValues = append(newValues, enum.FilterValuesALL)
		newValues = append(newValues, v...)

		list = append(list, &model.ApplicationKeys{
			Key:     k,
			Values:  newValues,
			KeyName: fmt.Sprintf("%s(应用)", k),
		})
	}

	return list, nil
}

func (a *applicationService) AppListByUUIDS(ctx context.Context, namespaceId int, uuids []string) ([]*model.Application, error) {
	list, err := a.applicationStore.GetList(ctx, namespaceId, uuids...)
	if err != nil {
		return nil, err
	}

	applications := make([]*model.Application, 0, len(list))

	for _, application := range list {
		val := &model.Application{Application: application}
		applications = append(applications, val)
	}

	return applications, nil
}

func (a *applicationService) getApplicationAdditional(extraHeader []entry.ApplicationExtraParam) []v1.ApplicationAdditional {
	additional := make([]v1.ApplicationAdditional, 0, len(extraHeader))
	for _, val := range extraHeader {
		position := "header"
		if val.Position != "" {
			position = val.Position
		}
		additional = append(additional, v1.ApplicationAdditional{
			Key:      val.Key,
			Value:    val.Value,
			Position: position,
		})
	}
	return additional
}

func (a *applicationService) entryExtraToModelExtra(extraParamList []entry.ApplicationExtraParam) []model.ApplicationExtraParam {
	extraParam := make([]model.ApplicationExtraParam, 0, len(extraParamList))
	for _, param := range extraParamList {
		extraParam = append(extraParam, model.ApplicationExtraParam{
			Key:      param.Key,
			Value:    param.Value,
			Conflict: param.Conflict,
			Position: param.Position,
		})
	}
	return extraParam
}

func (a *applicationService) entryAttrToModelAttr(attrs []entry.ApplicationCustomAttr) []model.ApplicationCustomAttr {
	customAttr := make([]model.ApplicationCustomAttr, 0, len(attrs))
	for _, attr := range attrs {
		customAttr = append(customAttr, model.ApplicationCustomAttr{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}
	return customAttr
}

func (a *applicationService) dtoExtraToEntryExtra(extraParamList []dto.ApplicationExtraParam) []entry.ApplicationExtraParam {
	extraParam := make([]entry.ApplicationExtraParam, 0, len(extraParamList))
	for _, param := range extraParamList {
		extraParam = append(extraParam, entry.ApplicationExtraParam{
			Key:      param.Key,
			Value:    param.Value,
			Conflict: param.Conflict,
			Position: param.Position,
		})
	}
	return extraParam
}

func (a *applicationService) dtoAttrToEntryAttr(attrs []dto.ApplicationCustomAttr) []entry.ApplicationCustomAttr {
	customAttr := make([]entry.ApplicationCustomAttr, 0, len(attrs))
	for _, attr := range attrs {
		customAttr = append(customAttr, entry.ApplicationCustomAttr{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}
	return customAttr
}
