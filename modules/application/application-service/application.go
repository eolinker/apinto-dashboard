package application_service

import (
	"context"
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	v2 "github.com/eolinker/apinto-dashboard/client/v2"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	driverInfo "github.com/eolinker/apinto-dashboard/driver"
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
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

const (
	anonymousIds          = "anonymous"
	professionApplication = "app"
)

var _ application.IApplicationService = (*applicationService)(nil)

type applicationService struct {
	applicationStore        application_store.IApplicationStore
	applicationAuthStore    application_store.IApplicationAuthStore
	applicationVersionStore application_store.IApplicationVersionStore
	applicationStatStore    application_store.IApplicationStatStore
	applicationHistoryStore application_store.IApplicationHistoryStore
	publishHistoryStore     application_store.IAppPublishHistoryStore
	clusterService          cluster.IClusterService
	randomService           random.IRandomService
	apintoClient            cluster.IApintoClient
	lockService             locker_service.IAsynLockService
	userInfoService         user.IUserInfoService
	driverManager           application.IAuthDriverManager
}

func newApplicationService() application.IApplicationService {
	app := &applicationService{}
	bean.Autowired(&app.applicationStore)
	bean.Autowired(&app.applicationAuthStore)
	bean.Autowired(&app.applicationVersionStore)
	bean.Autowired(&app.applicationStatStore)
	bean.Autowired(&app.applicationHistoryStore)
	bean.Autowired(&app.publishHistoryStore)
	bean.Autowired(&app.randomService)
	bean.Autowired(&app.clusterService)
	bean.Autowired(&app.apintoClient)
	bean.Autowired(&app.lockService)
	bean.Autowired(&app.userInfoService)
	bean.Autowired(&app.driverManager)
	return app
}

func (a *applicationService) OnlineInfo(ctx context.Context, namespaceId int, uuid string) (*application_model.ApplicationBasicInfo, []*application_model.AppCluster, error) {
	appInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, uuid)
	if err != nil {
		return nil, nil, err
	}

	info := &application_model.ApplicationBasicInfo{
		Uuid: appInfo.IdStr,
		Name: appInfo.Name,
		Desc: appInfo.Desc,
	}

	_, clusters, err := a.ClustersStatus(ctx, namespaceId, appInfo.Id, appInfo.IdStr, appInfo.Version)
	if err != nil {
		return nil, nil, err
	}
	items := make([]*application_model.AppCluster, 0, len(clusters))
	for _, clu := range clusters {
		items = append(items, &application_model.AppCluster{
			Name:       clu.Name,
			Title:      clu.Title,
			Env:        clu.Env,
			Status:     clu.Status,
			Updater:    clu.Updater,
			UpdateTime: clu.UpdateTime,
		})
	}
	return info, items, nil
}

func (a *applicationService) Online(ctx context.Context, namespaceId, userId int, id string, clusterNames []string) error {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(locker_service.LockNameApplication, applicationInfo.Id); err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationInfo.Id)
	//除了匿名应用以外，其他应用需要配置鉴权信息才可上线
	anonymous := true
	if applicationInfo.IdStr != anonymousIds {
		if !a.isApplicationSetAuth(ctx, applicationInfo.Id) {
			return errors.New("需要配置鉴权信息才可上线")
		}
		anonymous = false
	}
	//获取当前集群信息
	clusterInfos, err := a.clusterService.GetByNames(ctx, namespaceId, clusterNames)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        id,
		Name:        applicationInfo.Name,
		ClusterName: strings.Join(clusterNames, ","),
		PublishType: 1,
	})

	applicationId := applicationInfo.Id
	//获取当前应用的版本
	latestVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}
	auths := make([]v1.ApplicationAuth, 0)
	if !anonymous {
		//上线鉴权信息
		authList, err := a.OnlineAuth(ctx, applicationId)
		if err != nil {
			return err
		}
		for _, auth := range authList {
			auths = append(auths, a.GetDriver(auth.Driver).ToApinto(auth.ExpireTime, auth.Position, auth.TokenName, []byte(auth.Config), auth.IsTransparent))
		}
	}
	labels := make(map[string]string)
	for _, attr := range latestVersion.CustomAttrList {
		labels[attr.Key] = attr.Value
	}
	t := time.Now()

	for _, clu := range clusterInfos {
		err = a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
			publishHistory := &application_entry.AppPublishHistory{
				VersionName:              applicationInfo.Version,
				ClusterId:                clu.Id,
				NamespaceId:              namespaceId,
				Desc:                     applicationInfo.Desc,
				VersionId:                latestVersion.Id,
				Target:                   applicationInfo.Id,
				ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
				OptType:                  1, //上线
				Operator:                 userId,
				OptTime:                  t,
			}

			if err = a.publishHistoryStore.Insert(txCtx, publishHistory); err != nil {
				return err
			}

			return v2.Online(clu.Name, clu.Addr, professionApplication, applicationInfo.IdStr, &v2.WorkerInfo[v2.BasicInfo]{
				BasicInfo: &v2.BasicInfo{
					Profession:  professionApplication,
					Name:        applicationInfo.IdStr,
					Driver:      "app",
					Description: applicationInfo.Desc,
					Version:     applicationInfo.Version,
				},
				Append: map[string]interface{}{
					"auth":       auths,
					"labels":     labels,
					"additional": a.getApplicationAdditional(latestVersion.ExtraParamList),
					"anonymous":  anonymous,
				},
			})
		})
	}

	return nil
}

func (a *applicationService) Offline(ctx context.Context, namespaceId, userId int, id string, clusterNames []string) error {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return err
	}

	if err = a.lockService.Lock(locker_service.LockNameApplication, applicationInfo.Id); err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationInfo.Id)

	//获取当前集群信息
	clusterInfos, err := a.clusterService.GetByNames(ctx, namespaceId, clusterNames)
	if err != nil {
		return err
	}
	//获取当前应用的版本
	latestVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        id,
		Name:        applicationInfo.Name,
		ClusterName: strings.Join(clusterNames, ","),
		PublishType: 2,
	})
	t := time.Now()
	for _, clu := range clusterInfos {
		err = a.applicationStore.Transaction(ctx, func(txCtx context.Context) error {
			publishHistory := &application_entry.AppPublishHistory{
				VersionName:              applicationInfo.Version,
				ClusterId:                clu.Id,
				NamespaceId:              namespaceId,
				Desc:                     applicationInfo.Desc,
				VersionId:                latestVersion.Id,
				Target:                   applicationInfo.Id,
				ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
				OptType:                  3, //下线
				Operator:                 userId,
				OptTime:                  t,
			}

			if err = a.publishHistoryStore.Insert(txCtx, publishHistory); err != nil {
				return err
			}
			return common.CheckWorkerNotExist(v2.Offline(clu.Name, clu.Addr, professionApplication, applicationInfo.IdStr))
		})
	}
	return nil
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
		AuthList:       []*application_entry.ApplicationAuth{},
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

	tmpApplicationInfo, _ := a.applicationStore.GetByName(ctx, namespaceId, input.Name)
	if tmpApplicationInfo != nil && tmpApplicationInfo.IdStr != input.Id {
		return errors.New("应用名重复")
	}

	if applicationInfo.IdStr == anonymousIds && input.Name != "匿名应用" {
		return errors.New("不能更改匿名应用的应用名")
	}

	//获取应用当前版本
	latestVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}

	isUpdateVersion := false
	oldAttrMaps := make(map[string]string)
	for _, attr := range latestVersion.CustomAttrList {
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
	for _, extra := range latestVersion.ExtraParamList {
		oldExtraMaps[extra.Key] = fmt.Sprintf("%s-%s-%s", extra.Value, extra.Position, extra.Conflict)
	}
	newExtraMaps := make(map[string]string)
	for _, extra := range input.Params {
		newExtraMaps[extra.Key] = fmt.Sprintf("%s-%s-%s", extra.Value, extra.Position, extra.Conflict)
	}

	if !common.DiffMap(oldExtraMaps, newExtraMaps) {
		isUpdateVersion = true
	}

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

		versionConfig := application_entry.ApplicationVersionConfig{
			CustomAttrList: a.dtoAttrToEntryAttr(input.CustomAttrList),
			ExtraParamList: a.dtoExtraToEntryExtra(input.Params),
			AuthList:       latestVersion.AuthList,
		}

		if err = a.applicationHistoryStore.HistoryEdit(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application:              oldApplication,
			ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
		}, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: versionConfig,
		}, userId); err != nil {
			return nil
		}

		if isUpdateVersion {
			applicationInfo.Version = common.GenVersion(t)

			applicationVersion := &application_entry.ApplicationVersion{
				ApplicationID:            applicationInfo.Id,
				NamespaceID:              namespaceId,
				ApplicationVersionConfig: versionConfig,
				Operator:                 userId,
				CreateTime:               t,
			}
			if err = a.applicationVersionStore.Save(txCtx, applicationVersion); err != nil {
				return err
			}
			stat := &application_entry.ApplicationStat{
				ApplicationID: applicationVersion.ApplicationID,
				VersionID:     applicationVersion.Id,
			}
			err = a.applicationStatStore.Save(txCtx, stat)
			if err != nil {
				return err
			}
		}

		return a.applicationStore.Save(txCtx, applicationInfo)
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
				AuthList:       version.AuthList,
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
			Desc:         applicationInfo.Desc,
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

func (a *applicationService) AllApp(ctx context.Context, namespaceId int) ([]*application_model.ApplicationBasicInfo, error) {
	list, err := a.applicationStore.GetListByNamespace(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	applications := make([]*application_model.ApplicationBasicInfo, 0, len(list))
	for _, item := range list {
		applications = append(applications, &application_model.ApplicationBasicInfo{
			Uuid: item.IdStr,
			Name: item.Name,
			Desc: item.Desc,
		})
	}

	return applications, nil
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
		Params:     a.entryExtraToModelExtra(version.ExtraParamList),
	}
	return res, nil
}

func (a *applicationService) GetAppRemoteOptions(ctx context.Context, namespaceId, pageNum, pageSize int, keyword string) ([]any, error) {
	list, total, err := a.applicationStore.GetListPage(ctx, namespaceId, pageNum, pageSize, keyword)
	if err != nil {
		return nil, err
	}
	applications := make([]any, 0, total)
	for _, item := range list {
		applications = append(applications, application_model.ApplicationRemoteOption{
			Uuid:  item.IdStr,
			Title: item.Name,
			Desc:  item.Desc,
		})
	}
	return applications, nil
}

func (a *applicationService) AppInfo(ctx context.Context, namespaceId int, id string) (*application_model.ApplicationEntire, error) {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, id)
	if err != nil {
		return nil, err
	}

	res := &application_model.ApplicationEntire{
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
		val := &application_model.ApplicationBasicInfo{
			Uuid:       applicationInfo.IdStr,
			Name:       applicationInfo.Name,
			Desc:       applicationInfo.Desc,
			UpdateTime: applicationInfo.UpdateTime,
		}
		applications = append(applications, val)
	}

	return applications, nil
}

func (a *applicationService) AppBasicInfo(ctx context.Context, namespaceId int, uuid string) (*application_model.ApplicationBasicInfo, error) {
	info, err := a.applicationStore.GetByIdStr(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}
	return &application_model.ApplicationBasicInfo{
		Uuid: info.IdStr,
		Name: info.Name,
		Desc: info.Desc,
	}, nil
}

func (a *applicationService) getApplicationAdditional(extraHeader []application_entry.ApplicationExtraParam) []v1.ApplicationAdditional {
	additional := make([]v1.ApplicationAdditional, 0, len(extraHeader))
	for _, val := range extraHeader {
		additional = append(additional, v1.ApplicationAdditional{
			Key:      val.Key,
			Value:    val.Value,
			Position: val.Position,
			Conflict: val.Conflict,
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

func (a *applicationService) dtoExtraToEntryExtra(extraParamList []application_dto.ExtraParam) []application_entry.ApplicationExtraParam {
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
		v, err := a.publishHistoryStore.GetLastPublishHistory(ctx, map[string]interface{}{
			"namespace": namespaceId,
			"cluster":   c.Id,
			"target":    appId,
			"kind":      "application",
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

func (a *applicationService) GetDriver(driver string) application.IAuthDriver {
	return a.driverManager.GetDriver(driver)
}

func (a *applicationService) GetAuthList(ctx context.Context, namespaceId int, appId string) ([]*application_model.AppAuthItem, error) {
	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, appId)
	if err != nil {
		return nil, err
	}

	list, err := a.applicationAuthStore.GetListByApplication(ctx, applicationInfo.Id)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(list, func(t *application_entry.ApplicationAuth) int {
		return t.Operator
	})

	userInfoMaps, _ := a.userInfoService.GetUserInfoMaps(ctx, userIds...)

	resList := make([]*application_model.AppAuthItem, 0, len(list))
	for _, auth := range list {

		operatorName := ""
		if userInfo, ok := userInfoMaps[auth.Operator]; ok {
			operatorName = userInfo.NickName
		}

		authModel := &application_model.AppAuthItem{
			UUID:           auth.Uuid,
			Title:          auth.Title,
			Driver:         auth.Driver,
			Operator:       operatorName,
			HideCredential: auth.IsTransparent,
			ExpireTime:     auth.ExpireTime,
			UpdateTime:     auth.UpdateTime,
		}

		resList = append(resList, authModel)

	}
	return resList, nil
}

func (a *applicationService) CreateAuth(ctx context.Context, namespaceId, userId int, appId string, input *application_dto.ApplicationAuthInput) error {
	driverAuth := a.driverManager.GetDriver(input.Driver)
	if err := driverAuth.CheckInput(input.Config); err != nil {
		return err
	}

	if input.ExpireTime > 0 && input.ExpireTime < time.Now().Unix() {
		return errors.New("过期时间不能小于当前时间")
	}

	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, appId)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameApplication, applicationInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationInfo.Id)

	if applicationInfo.IdStr == anonymousIds {
		return errors.New("匿名应用不能添加鉴权信息")
	}

	t := time.Now()
	applicationAuth := &application_entry.ApplicationAuth{
		Uuid:          uuid.New(),
		Title:         input.Title,
		Namespace:     namespaceId,
		Application:   applicationInfo.Id,
		IsTransparent: input.HideCredential,
		Driver:        input.Driver,
		Position:      input.Position,
		TokenName:     input.TokenName,
		ExpireTime:    input.ExpireTime,
		Config:        string(input.Config),
		Operator:      userId,
		CreateTime:    t,
		UpdateTime:    t,
	}

	latestVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}
	applicationAuths, err := a.applicationAuthStore.GetListByApplication(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}

	return a.applicationAuthStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = a.applicationAuthStore.Insert(txCtx, applicationAuth); err != nil {
			return err
		}
		applicationAuths = append(applicationAuths, applicationAuth)
		latestVersion.ApplicationVersionConfig.AuthList = applicationAuths

		applicationInfo.Version = common.GenVersion(t)
		applicationInfo.UpdateTime = t
		err = a.applicationStore.Save(txCtx, applicationInfo)
		if err != nil {
			return err
		}

		if err := a.applicationHistoryStore.HistoryAdd(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
		}, userId); err != nil {
			return nil
		}

		applicationVersion := &application_entry.ApplicationVersion{
			ApplicationID:            applicationInfo.Id,
			NamespaceID:              namespaceId,
			ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
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

func (a *applicationService) UpdateAuth(ctx context.Context, namespaceId, userId int, appId, uuidStr string, input *application_dto.ApplicationAuthInput) error {
	driverAuth := a.driverManager.GetDriver(input.Driver)
	if err := driverAuth.CheckInput(input.Config); err != nil {
		return err
	}

	applicationInfo, err := a.applicationStore.GetByIdStr(ctx, namespaceId, appId)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameApplication, applicationInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, applicationInfo.Id)

	authInfo, err := a.applicationAuthStore.GetByUUID(ctx, uuidStr)
	if err != nil {
		return err
	}
	isUpdate := false
	if authInfo.ExpireTime != input.ExpireTime {
		isUpdate = true
	} else if authInfo.IsTransparent != input.HideCredential {
		isUpdate = true
	} else if authInfo.TokenName != input.TokenName {
		isUpdate = true
	} else if authInfo.Position != input.Position {
		isUpdate = true
	} else if authInfo.Driver != input.Driver {
		isUpdate = true
	}

	if authInfo.Config != string(input.Config) {
		isUpdate = true
	}

	t := time.Now()

	authInfo.IsTransparent = input.HideCredential
	authInfo.ExpireTime = input.ExpireTime
	authInfo.Operator = userId
	authInfo.UpdateTime = t
	authInfo.Driver = input.Driver
	authInfo.Position = input.Position
	authInfo.TokenName = input.TokenName
	authInfo.Config = string(input.Config)
	authInfo.Title = input.Title

	latestVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}
	newAppVersionCfg := latestVersion.ApplicationVersionConfig

	return a.applicationAuthStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = a.applicationAuthStore.Save(txCtx, authInfo); err != nil {
			return err
		}

		if isUpdate {
			applicationInfo.Version = common.GenVersion(t)
			applicationInfo.UpdateTime = t
			err = a.applicationStore.Save(txCtx, applicationInfo)
			if err != nil {
				return err
			}

			//替换app version里的authList
			for i, auth := range latestVersion.AuthList {
				if auth.Id == authInfo.Id {
					newAppVersionCfg.AuthList[i] = authInfo
					break
				}
			}

			applicationVersion := &application_entry.ApplicationVersion{
				ApplicationID:            applicationInfo.Id,
				NamespaceID:              namespaceId,
				ApplicationVersionConfig: newAppVersionCfg,
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

			err = a.applicationStatStore.Save(txCtx, stat)
			if err != nil {
				return err
			}
		}

		return a.applicationHistoryStore.HistoryEdit(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
		}, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: newAppVersionCfg,
		}, userId)
	})
}

func (a *applicationService) AuthInfo(ctx context.Context, namespaceId int, appId, uuid string) (*application_model.ApplicationAuth, error) {
	_, err := a.applicationStore.GetByIdStr(ctx, namespaceId, appId)
	if err != nil {
		return nil, err
	}

	auth, err := a.applicationAuthStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	userInfo, _ := a.userInfoService.GetUserInfo(ctx, auth.Operator)
	resAuth := &application_model.ApplicationAuth{
		ApplicationAuth: auth,
		Operator:        userInfo.NickName,
		Config:          auth.Config,
	}
	return resAuth, nil
}

func (a *applicationService) AuthDetails(ctx context.Context, namespaceId int, appId, uuid string) ([]application_model.AuthDetailItem, error) {
	_, err := a.applicationStore.GetByIdStr(ctx, namespaceId, appId)
	if err != nil {
		return nil, err
	}

	auth, err := a.applicationAuthStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	driver := a.GetDriver(auth.Driver)
	if driver == nil {
		return nil, errors.New("鉴权类型不存在")
	}
	//TODO 临时处理, 前端没时间
	cfgItems := driver.GetCfgDetails([]byte(auth.Config))
	details := make([]application_model.AuthDetailItem, 0, 6+len(cfgItems))
	details = append(details, application_model.AuthDetailItem{Key: "名称", Value: auth.Title})
	details = append(details, application_model.AuthDetailItem{Key: "鉴权类型", Value: auth.Driver})
	details = append(details, application_model.AuthDetailItem{Key: "参数位置", Value: auth.Position})
	details = append(details, application_model.AuthDetailItem{Key: "参数名", Value: auth.TokenName})
	details = append(details, cfgItems...)
	dateStr := "永久"
	if auth.ExpireTime != 0 {
		dateStr = time.Unix(auth.ExpireTime, 0).Format("2006-01-02")
	}
	details = append(details, application_model.AuthDetailItem{Key: "过期日期", Value: dateStr})
	hideAuthStr := "是"
	if !auth.IsTransparent {
		hideAuthStr = "否"
	}
	details = append(details, application_model.AuthDetailItem{Key: "隐藏鉴权信息", Value: hideAuthStr})
	return details, nil
}

func (a *applicationService) DeleteAuth(ctx context.Context, namespaceId, userId int, uuid string) error {
	authInfo, err := a.applicationAuthStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameApplication, authInfo.Application)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameApplication, authInfo.Application)

	applicationInfo, err := a.applicationStore.Get(ctx, authInfo.Application)
	if err != nil {
		return err
	}

	latestVersion, err := a.getAppVersion(ctx, applicationInfo.Id)
	if err != nil {
		return err
	}
	newAppVersionCfg := latestVersion.ApplicationVersionConfig
	t := time.Now()

	return a.applicationAuthStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = a.applicationAuthStore.Delete(ctx, authInfo.Id); err != nil {
			return err
		}

		applicationInfo.Version = common.GenVersion(t)
		applicationInfo.UpdateTime = t
		err = a.applicationStore.Save(txCtx, applicationInfo)
		if err != nil {
			return err
		}

		//替换app version里的authList
		for i, auth := range latestVersion.AuthList {
			if auth.Id == authInfo.Id {
				newAppVersionCfg.AuthList = append(newAppVersionCfg.AuthList[:i], newAppVersionCfg.AuthList[i+1:]...)
				break
			}
		}

		applicationVersion := &application_entry.ApplicationVersion{
			ApplicationID:            applicationInfo.Id,
			NamespaceID:              namespaceId,
			ApplicationVersionConfig: newAppVersionCfg,
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

		err = a.applicationStatStore.Save(txCtx, stat)
		if err != nil {
			return err
		}

		return a.applicationHistoryStore.HistoryEdit(txCtx, namespaceId, applicationInfo.Id, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: latestVersion.ApplicationVersionConfig,
		}, &application_entry.ApplicationHistoryInfo{
			Application:              *applicationInfo,
			ApplicationVersionConfig: newAppVersionCfg,
		}, userId)
	})
}

func (a *applicationService) OnlineAuth(ctx context.Context, applicationId int) ([]*application_model.ApplicationAuth, error) {
	applicationAuths, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return nil, err
	}

	list := make([]*application_model.ApplicationAuth, 0, len(applicationAuths))
	for _, auth := range applicationAuths {
		list = append(list, &application_model.ApplicationAuth{
			ApplicationAuth: auth,
			Config:          auth.Config,
		})
	}

	return list, nil
}

func (a *applicationService) isApplicationSetAuth(ctx context.Context, applicationId int) bool {
	list, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return false
	}
	if len(list) > 0 {
		return true
	}
	return false
}

func (a *applicationService) GetDriversRender() []*driverInfo.DriverInfo {
	return a.driverManager.List()
}

func (a *applicationService) getAuthParamInfo(auth *application_model.ApplicationAuth) string {
	return a.driverManager.GetDriver(auth.Driver).GetAuthListInfo([]byte(auth.Config))
}
