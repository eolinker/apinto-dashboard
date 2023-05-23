package application_service

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/apinto-dashboard/modules/application/application-dto"
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"github.com/eolinker/apinto-dashboard/modules/application/application-model"
	"github.com/eolinker/apinto-dashboard/modules/application/application-store"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/random"
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

const (
	anonymousIds          = "anonymous"
	professionApplication = "application"
)

var _ application.IApplicationService = (*applicationService)(nil)

type applicationService struct {
	applicationStore            application_store.IApplicationStore
	applicationRuntimeStore     application_store.IApplicationRuntimeStore
	applicationAuthRuntimeStore application_store.IApplicationAuthRuntimeStore
	applicationVersionStore     application_store.IApplicationVersionStore
	applicationStatStore        application_store.IApplicationStatStore
	applicationHistoryStore     application_store.IApplicationHistoryStore
	clusterService              cluster.IClusterService
	applicationAuthService      application.IApplicationAuthService
	randomService               random.IRandomService
	apintoClient                cluster.IApintoClient
	lockService                 locker_service.IAsynLockService
	userInfoService             user.IUserInfoService
}

func newApplicationService() application.IApplicationService {
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

func (a *applicationService) OnlineList(ctx context.Context, namespaceId int, id string) ([]*application_model.ApplicationOnline, error) {
	app, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}
	applicationId := app.Id

	//获取工作空间下的所有集群
	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *cluster_model.Cluster) int {
		return t.Id
	})

	//获取当前应用下集群运行的版本
	runtimes, err := a.applicationRuntimeStore.GetByTarget(ctx, applicationId)
	if err != nil {
		return nil, err
	}
	//最新版本
	lastVersion, err := a.getAppVersion(ctx, app.Id)
	if err != nil {
		return nil, err
	}

	runtimeMaps := common.SliceToMap(runtimes, func(t *application_entry.ApplicationRuntime) int {
		return t.ClusterId
	})

	userIds := common.SliceToSliceIds(runtimes, func(t *application_entry.ApplicationRuntime) int {
		return t.Operator
	})

	userInfoMaps, _ := a.userInfoService.GetUserInfoMaps(ctx, userIds...)

	list := make([]*application_model.ApplicationOnline, 0, len(clusters))
	for _, clusterInfo := range clusterMaps {

		applicationOnline := &application_model.ApplicationOnline{
			ClusterID:   clusterInfo.Id,
			ClusterName: clusterInfo.Name,
			Env:         clusterInfo.Env,
			Status:      1, //默认为未上线状态
		}

		if runtime, ok := runtimeMaps[clusterInfo.Id]; ok {
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
					isUpdate, err := a.applicationAuthService.IsUpdate(ctx, clusterInfo.Id, currentVersion.ApplicationID)
					if err != nil {
						return nil, err
					}
					if isUpdate {
						applicationOnline.Status = 4
					}
				}

			}
		} else {
			if app.IdStr == anonymousIds {
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
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}
	//除了匿名应用以外，其他应用需要配置鉴权信息才可上线
	anonymous := true
	if applicationInfo.IdStr != anonymousIds {
		auths, err := a.applicationAuthService.GetListByApplicationId(ctx, applicationInfo.Id)
		if err != nil {
			return err
		}
		if len(auths) == 0 {
			return errors.New("需要配置鉴权信息才可上线")
		}
		anonymous = false
	}
	//获取当前集群信息
	clusterInfo, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	applicationId := applicationInfo.Id
	clusterId := clusterInfo.Id

	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(locker_service.LockNameApplication, applicationId); err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationId)

	//拿到锁后需要重新获取下信息
	applicationInfo, err = a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	//获取当前应用的版本
	lastVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}

	runtime, err := a.applicationRuntimeStore.GetForCluster(ctx, applicationId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	t := time.Now()
	if runtime == nil {
		runtime = &application_entry.ApplicationRuntime{
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
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        id,
		Name:        applicationInfo.Name,
		ClusterId:   clusterInfo.Id,
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
			authList, err := a.applicationAuthService.Online(txCtx, namespaceId, userId, clusterId, applicationId)
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
			Name:        applicationInfo.IdStr,
			Driver:      "app",
			Auth:        auths,
			Disable:     runtime.Disable,
			Description: applicationInfo.Desc,
			Labels:      labels,
			Additional:  a.getApplicationAdditional(lastVersion.ExtraParamList),
			Anonymous:   anonymous,
		}

		if runtime.Id > 0 {
			return client.ForApp().Update(applicationInfo.IdStr, *appConfig)
		}
		return client.ForApp().Create(*appConfig)
	})
}

func (a *applicationService) Offline(ctx context.Context, namespaceId, userId int, id, clusterName string) error {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	//获取当前集群信息
	clusterInfo, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	applicationId := applicationInfo.Id
	clusterId := clusterInfo.Id

	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(locker_service.LockNameApplication, applicationId); err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationId)

	//拿到锁后需要重新获取下信息
	applicationInfo, err = a.applicationStore.GetByIdStr(ctx, namespaceId, id)
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
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        id,
		Name:        applicationInfo.Name,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
		PublishType: 2,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = a.applicationRuntimeStore.Save(txCtx, runtime); err != nil {
			return err
		}
		//鉴权信息下线
		if err = a.applicationAuthService.Offline(txCtx, clusterId, applicationId); err != nil {
			return err
		}

		return common.CheckWorkerNotExist(client.ForApp().Delete(applicationInfo.IdStr))
	})
}

func (a *applicationService) CreateApp(ctx context.Context, namespaceId, userId int, input *application_dto.ApplicationInput) error {
	input.Id = strings.ToLower(input.Id)
	applicationInfo, _ := a.applicationStore.GetByIdStr(ctx, namespaceId, input.Id)
	if applicationInfo != nil {
		return errors.New("应用ID重复")
	}

	applicationInfo, _ = a.applicationStore.GetByName(ctx, namespaceId, input.Name)
	if applicationInfo != nil {
		return errors.New("应用名重复")
	}

	versionConfig := application_entry.ApplicationVersionConfig{
		CustomAttrList: a.dtoAttrToEntryAttr(input.CustomAttrList),
		ExtraParamList: []application_entry.ApplicationExtraParam{},
	}
	t := time.Now()

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.Id,
		Name: input.Name,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
		applicationInfo = &application_entry.Application{
			NamespaceId: namespaceId,
			IdStr:       input.Id,
			Name:        input.Name,
			Desc:        input.Desc,
			Version:     common.GenVersion(t),
			Operator:    userId,
			CreateTime:  t,
			UpdateTime:  t,
		}

		if err := a.applicationStore.Save(txCtx, applicationInfo); err != nil {
			return err
		}

		if err := a.applicationHistoryStore.HistoryAdd(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: versionConfig,
		}, userId); err != nil {
			return nil
		}

		applicationVersion := &application_entry.ApplicationVersion{
			ApplicationID:            applicationInfo.Id,
			NamespaceID:              namespaceId,
			ApplicationVersionConfig: versionConfig,
			Operator:                 userId,
			CreateTime:               t,
		}

		if err := a.applicationVersionStore.Save(txCtx, applicationVersion); err != nil {
			return err
		}
		stat := &application_entry.ApplicationStat{
			ApplicationID: applicationVersion.ApplicationID,
			VersionID:     applicationVersion.Id,
		}

		return a.applicationStatStore.Save(txCtx, stat)
	})

}

func (a *applicationService) UpdateApp(ctx context.Context, namespaceId, userId int, input *application_dto.ApplicationInput) error {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, input.Id)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameApplication, applicationInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationInfo.Id)

	applicationInfo, _ = a.applicationStore.GetByName(ctx, namespaceId, input.Name)
	if applicationInfo != nil && applicationInfo.IdStr != input.Id {
		return errors.New("应用名重复")
	}

	if applicationInfo.IdStr == anonymousIds && input.Name != "匿名应用" {
		return errors.New("不能更改匿名应用的应用名")
	}

	//获取应用当前版本
	version, err := a.getAppVersion(ctx, applicationInfo.Id)
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

	//oldExtraMaps := make(map[string]string)
	//for _, extra := range version.ExtraParamList {
	//	oldExtraMaps[extra.Key] = extra.Value
	//}
	//newExtraMaps := make(map[string]string)
	//for _, extra := range input.ExtraParamList {
	//	newExtraMaps[extra.Key] = extra.Value
	//}
	//
	//if !common.DiffMap(oldExtraMaps, newExtraMaps) {
	//	isUpdateVersion = true
	//}

	t := time.Now()
	//添加操作记录

	oldApplication := *applicationInfo

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.Id,
		Name: input.Name,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
		applicationInfo.UpdateTime = t
		applicationInfo.Operator = userId
		applicationInfo.Desc = input.Desc
		applicationInfo.Name = input.Name
		if isUpdateVersion {
			applicationInfo.Version = common.GenVersion(t)
		}

		if err = a.applicationStore.Save(txCtx, applicationInfo); err != nil {
			return err
		}

		versionConfig := application_entry.ApplicationVersionConfig{
			CustomAttrList: a.dtoAttrToEntryAttr(input.CustomAttrList),
			ExtraParamList: version.ExtraParamList,
		}

		applicationVersion := &application_entry.ApplicationVersion{
			ApplicationID:            applicationInfo.Id,
			NamespaceID:              namespaceId,
			ApplicationVersionConfig: versionConfig,
			Operator:                 userId,
			CreateTime:               t,
		}

		if err = a.applicationHistoryStore.HistoryEdit(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application: oldApplication,
			ApplicationVersionConfig: application_entry.ApplicationVersionConfig{
				CustomAttrList: version.CustomAttrList,
				ExtraParamList: version.ExtraParamList,
			},
		}, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: versionConfig,
		}, userId); err != nil {
			return nil
		}

		if isUpdateVersion {
			if err = a.applicationVersionStore.Save(txCtx, applicationVersion); err != nil {
				return err
			}
			stat := &application_entry.ApplicationStat{
				ApplicationID: applicationVersion.ApplicationID,
				VersionID:     applicationVersion.Id,
			}
			return a.applicationStatStore.Save(txCtx, stat)
		}
		return nil
	})
}

func (a *applicationService) DelApp(ctx context.Context, namespaceId, userId int, id string) error {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameApplication, applicationInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationInfo.Id)

	if applicationInfo.IdStr == anonymousIds {
		return errors.New("匿名应用不能删除")
	}

	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}

	isOnline, err := a.isAppOnline(applicationInfo.IdStr, clusters)
	if err != nil {
		return err
	}

	if isOnline {
		return errors.New("应用已上线不可删除")
	}

	//获取应用当前版本信息
	version, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: id,
		Name: applicationInfo.Name,
	})

	return a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = a.applicationStore.Delete(txCtx, applicationInfo.Id); err != nil {
			return err
		}

		//添加操作记录
		if err = a.applicationHistoryStore.HistoryDelete(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application: *applicationInfo,
			ApplicationVersionConfig: application_entry.ApplicationVersionConfig{
				CustomAttrList: version.CustomAttrList,
				ExtraParamList: version.ExtraParamList,
			},
		}, userId); err != nil {
			return nil
		}

		delMap := make(map[string]interface{})
		delMap["`kind`"] = "application"
		delMap["`target`"] = applicationInfo.Id

		if _, err = a.applicationStatStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}
		if _, err = a.applicationVersionStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}
		for _, clusterInfo := range clusters {
			delMap["`cluster`"] = clusterInfo.Id
			if _, err = a.applicationRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}
		}

		return nil
	})
}

func (a *applicationService) AppList(ctx context.Context, namespaceId, userId, pageNum, pageSize int, queryName string, clusters []string) ([]*application_model.ApplicationListItem, int, error) {
	anonymousApplication, err := a.applicationStore.GetByIdStr(ctx, namespaceId, anonymousIds)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	//没有匿名应用创建一个
	if anonymousApplication == nil {
		t := time.Now()
		entryApplication := &application_entry.Application{
			NamespaceId: namespaceId,
			IdStr:       anonymousIds,
			Name:        "匿名应用",
			Version:     common.GenVersion(t),
			Operator:    userId,
			CreateTime:  t,
			UpdateTime:  t,
		}

		err = a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
			if err = a.applicationStore.Save(txCtx, entryApplication); err != nil {
				return err
			}
			version := &application_entry.ApplicationVersion{
				ApplicationID:            entryApplication.Id,
				NamespaceID:              namespaceId,
				ApplicationVersionConfig: application_entry.ApplicationVersionConfig{},
				Operator:                 userId,
				CreateTime:               t,
			}

			if err = a.applicationVersionStore.Save(txCtx, version); err != nil {
				return err
			}
			return a.applicationStatStore.Save(txCtx, &application_entry.ApplicationStat{
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

	userIds := common.SliceToSliceIds(list, func(t *application_entry.Application) int {
		return t.Operator
	})
	userInfoMaps, _ := a.userInfoService.GetUserInfoMaps(ctx, userIds...)

	var clusterInfos []*cluster_model.Cluster
	if len(clusters) > 0 {
		clusterInfos, err = a.clusterService.GetByNames(ctx, namespaceId, clusters)
		if err != nil {
			return nil, 0, err
		}
	} else {
		clusterInfos, err = a.clusterService.GetByNamespaceId(ctx, namespaceId)
		if err != nil {
			return nil, 0, err
		}
	}
	clusterVersions := a.getApintoAPPVersions(clusterInfos)

	applications := make([]*application_model.ApplicationListItem, 0, len(list))
	for _, applicationInfo := range list {
		operatorName := ""
		if userInfo, ok := userInfoMaps[applicationInfo.Operator]; ok {
			operatorName = userInfo.NickName
		}

		isOnline := false
		publish := make([]*application_model.APPListItemPublish, 0, len(clusterInfos))
		for _, clu := range clusterInfos {
			cluVersions, has := clusterVersions[clu.Name]
			if !has {
				publish = append(publish, &application_model.APPListItemPublish{
					Name:   clu.Name,
					Title:  clu.Title,
					Status: 1, //未发布
				})
				continue
			}

			vers, has := cluVersions[applicationInfo.IdStr]
			if !has {
				publish = append(publish, &application_model.APPListItemPublish{
					Name:   clu.Name,
					Title:  clu.Title,
					Status: 1, //未发布
				})
				continue
			}

			isOnline = true
			status := 4 //待更新
			if vers == applicationInfo.Version {
				status = 3 //上线
			}

			publish = append(publish, &application_model.APPListItemPublish{
				Name:   clu.Name,
				Title:  clu.Title,
				Status: status,
			})
		}

		item := &application_model.ApplicationListItem{
			Uuid:         applicationInfo.IdStr,
			Name:         applicationInfo.Name,
			Desc:         applicationInfo.Name,
			UpdateTime:   applicationInfo.UpdateTime,
			OperatorName: operatorName,
			IsDelete:     !isOnline,
			Publish:      publish,
		}

		//匿名应用默认不可以删除
		if item.Uuid == anonymousIds {
			item.IsDelete = false
		}

		applications = append(applications, item)
	}

	//对列表进行排序，匿名排第一位，其余按更新时间降序
	sort.Sort(application_model.ApplicationList(applications))

	return applications, count, nil
}
func (a *applicationService) AppEnumList(ctx context.Context, namespaceId int) ([]*application_model.ApplicationBasicInfo, error) {
	list, err := a.applicationStore.GetList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	applications := make([]*application_model.ApplicationBasicInfo, 0, len(list))

	for _, item := range list {
		applications = append(applications, &application_model.ApplicationBasicInfo{
			Uuid:       item.IdStr,
			Name:       item.Name,
			UpdateTime: item.UpdateTime,
		})
	}

	sort.Sort(application_model.ApplicationBasicInfoList(applications))

	return applications, nil
}

func (a *applicationService) AppListFilter(ctx context.Context, namespaceId, pageNum, pageSize int, queryName string) ([]*application_model.ApplicationBasicInfo, int, error) {

	list, count, err := a.applicationStore.GetListPage(ctx, namespaceId, pageNum, pageSize, queryName)
	if err != nil {
		return nil, 0, err
	}

	applications := make([]*application_model.ApplicationBasicInfo, 0, len(list))

	for _, item := range list {
		applications = append(applications, &application_model.ApplicationBasicInfo{
			Uuid: item.IdStr,
			Name: item.Name,
			Desc: item.Desc,
		})
	}

	return applications, count, nil
}

func (a *applicationService) AppInfoDetails(ctx context.Context, namespaceId int, id string) (*application_model.ApplicationInfo, error) {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}

	version, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return nil, err
	}

	res := &application_model.ApplicationInfo{
		Name:       applicationInfo.Name,
		Uuid:       applicationInfo.IdStr,
		Desc:       applicationInfo.Desc,
		CustomAttr: a.entryAttrToModelAttr(version.CustomAttrList),
	}
	return res, nil
}

func (a *applicationService) AppInfo(ctx context.Context, namespaceId int, id string) (*application_model.ApplicationBasicInfo, error) {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}

	res := &application_model.ApplicationBasicInfo{
		Application: applicationInfo,
	}
	return res, nil
}

func (a *applicationService) getAppVersion(ctx context.Context, appId int) (*application_entry.ApplicationVersion, error) {
	stat, err := a.applicationStatStore.Get(ctx, appId)
	if err != nil {
		return nil, err
	}

	version, err := a.applicationVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func (a *applicationService) GetAppKeys(ctx context.Context, namespaceId int) ([]*application_model.ApplicationKeys, error) {
	applications, err := a.applicationStore.GetList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	list := make([]*application_model.ApplicationKeys, 0)

	keys := map[string][]string{}

	for _, applicationInfo := range applications {

		version, err := a.getAppVersion(ctx, applicationInfo.Id)
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
		newValues = append(newValues, config.FilterValuesALL)
		newValues = append(newValues, v...)

		list = append(list, &application_model.ApplicationKeys{
			Key:     k,
			Values:  newValues,
			KeyName: fmt.Sprintf("%s(应用)", k),
		})
	}

	return list, nil
}

func (a *applicationService) AppListByUUIDS(ctx context.Context, namespaceId int, uuids []string) ([]*application_model.ApplicationBasicInfo, error) {
	list, err := a.applicationStore.GetList(ctx, namespaceId, uuids...)
	if err != nil {
		return nil, err
	}

	applications := make([]*application_model.ApplicationBasicInfo, 0, len(list))

	for _, applicationInfo := range list {
		val := &application_model.ApplicationBasicInfo{Application: applicationInfo}
		applications = append(applications, val)
	}

	return applications, nil
}

func (a *applicationService) getApplicationAdditional(extraHeader []application_entry.ApplicationExtraParam) []v1.ApplicationAdditional {
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

func (a *applicationService) entryExtraToModelExtra(extraParamList []application_entry.ApplicationExtraParam) []application_model.ApplicationExtraParam {
	extraParam := make([]application_model.ApplicationExtraParam, 0, len(extraParamList))
	for _, param := range extraParamList {
		extraParam = append(extraParam, application_model.ApplicationExtraParam{
			Key:      param.Key,
			Value:    param.Value,
			Conflict: param.Conflict,
			Position: param.Position,
		})
	}
	return extraParam
}

func (a *applicationService) entryAttrToModelAttr(attrs []application_entry.ApplicationCustomAttr) []application_model.ApplicationCustomAttr {
	customAttr := make([]application_model.ApplicationCustomAttr, 0, len(attrs))
	for _, attr := range attrs {
		customAttr = append(customAttr, application_model.ApplicationCustomAttr{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}
	return customAttr
}

func (a *applicationService) dtoExtraToEntryExtra(extraParamList []application_dto.ApplicationExtraParam) []application_entry.ApplicationExtraParam {
	extraParam := make([]application_entry.ApplicationExtraParam, 0, len(extraParamList))
	for _, param := range extraParamList {
		extraParam = append(extraParam, application_entry.ApplicationExtraParam{
			Key:      param.Key,
			Value:    param.Value,
			Conflict: param.Conflict,
			Position: param.Position,
		})
	}
	return extraParam
}

func (a *applicationService) dtoAttrToEntryAttr(attrs []application_dto.ApplicationCustomAttr) []application_entry.ApplicationCustomAttr {
	customAttr := make([]application_entry.ApplicationCustomAttr, 0, len(attrs))
	for _, attr := range attrs {
		customAttr = append(customAttr, application_entry.ApplicationCustomAttr{
			Key:   attr.Key,
			Value: attr.Value,
		})
	}
	return customAttr
}

func (a *applicationService) getApintoAPPVersions(clusters []*cluster_model.Cluster) map[string]map[string]string {
	results := make(map[string]map[string]string, len(clusters))

	for _, c := range clusters {
		client, err := v2.GetClusterClient(c.Name, c.Addr)
		if err != nil {
			log.Errorf("get cluster %s Client error: %v", c.Name, err)
			continue
		}
		versions, err := client.Versions(professionApplication)
		if err != nil {
			log.Errorf("get cluster status error: %v", err)
			continue
		}
		results[c.Name] = versions
	}
	return results
}

func (a *applicationService) isAppOnline(appUUID string, clusters []*cluster_model.Cluster) (bool, error) {
	online := false
	for _, c := range clusters {
		client, err := v2.GetClusterClient(c.Name, c.Addr)
		if err != nil {
			log.Errorf("get cluster status error: %v", err)
			continue
		}

		_, err = client.Version(professionApplication, appUUID)
		if err != nil {
			continue
		}
		online = true
		break
	}
	return online, nil
}

func (a *applicationService) ClustersStatus(ctx context.Context, namespaceId, appId int, appUUID, appVersion string) (bool, []*application_model.AppCluster, error) {
	clusters, err := a.clusterService.GetAllCluster(ctx)
	if err != nil {
		return false, nil, err
	}
	result := make([]*application_model.AppCluster, 0, len(clusters))
	online := false
	for _, c := range clusters {
		var operator int
		var updateTime string
		v, err := a.apiPublishHistory.GetLastPublishHistory(ctx, map[string]interface{}{
			"namespace": namespaceId,
			"cluster":   c.Id,
			"target":    appId,
			"kind":      "api",
		})
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				result = append(result, &application_model.AppCluster{
					Name:   c.Name,
					Title:  c.Title,
					Env:    c.Env,
					Status: 1, //未发布
				})
				continue
			}
			// 可能存在id不相同，但是控制台已经发布的情况
		} else {
			operator = v.Operator
			updateTime = v.OptTime.Format("2006-01-02 15:04:05")
		}

		client, err := v2.GetClusterClient(c.Name, c.Addr)
		if err != nil {
			result = append(result, &application_model.AppCluster{
				Name:   c.Name,
				Title:  c.Title,
				Env:    c.Env,
				Status: 1, //未发布
			})
			log.Errorf("get cluster status error: %v", err)
			continue
		}

		updater := ""
		if operator > 0 {
			u, err := a.userInfoService.GetUserInfo(ctx, operator)
			if err == nil {
				updater = u.UserName
			}
		}
		version, err := client.Version(professionApplication, appUUID)
		if err != nil {
			result = append(result, &application_model.AppCluster{
				Name:       c.Name,
				Title:      c.Title,
				Env:        c.Env,
				Status:     1, //未发布
				Updater:    updater,
				UpdateTime: updateTime,
			})
			continue
		}
		online = true
		status := 4 //待更新
		if version == appVersion {
			status = 3 //上线
		}

		result = append(result, &application_model.AppCluster{
			Name:       c.Name,
			Title:      c.Title,
			Env:        c.Env,
			Status:     status,
			Updater:    updater,
			UpdateTime: updateTime,
		})
	}
	return online, result, nil
}
