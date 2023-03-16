package application_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	driverInfo "github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/apinto-dashboard/modules/application/application-dto"
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"github.com/eolinker/apinto-dashboard/modules/application/application-model"
	"github.com/eolinker/apinto-dashboard/modules/application/application-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"time"
)

type applicationAuthService struct {
	applicationAuthStore        application_store.IApplicationAuthStore
	applicationAuthVersionStore application_store.IApplicationAuthVersionStore
	applicationAuthStatStore    application_store.IApplicationAuthStatStore
	applicationAuthRuntimeStore application_store.IApplicationAuthRuntimeStore
	applicationAuthPublishStore application_store.IApplicationAuthPublishStore
	applicationAuthHistoryStore application_store.IApplicationAuthHistoryStore
	applicationService          application.IApplicationService
	clusterService              cluster.IClusterService
	userInfoService             user.IUserInfoService
	driverManager               application.IAuthDriverManager
}

func (a *applicationAuthService) GetDriver(driver string) application.IAuthDriver {
	return a.driverManager.GetDriver(driver)
}

func newApplicationAuth() application.IApplicationAuthService {
	auth := &applicationAuthService{}

	bean.Autowired(&auth.applicationAuthStore)
	bean.Autowired(&auth.applicationAuthStatStore)
	bean.Autowired(&auth.applicationAuthRuntimeStore)
	bean.Autowired(&auth.applicationAuthPublishStore)
	bean.Autowired(&auth.applicationAuthVersionStore)
	bean.Autowired(&auth.applicationAuthHistoryStore)
	bean.Autowired(&auth.applicationService)
	bean.Autowired(&auth.clusterService)
	bean.Autowired(&auth.driverManager)
	bean.Autowired(&auth.userInfoService)
	return auth
}

func (a *applicationAuthService) GetList(ctx context.Context, namespaceId int, appId string) ([]*application_model.ApplicationAuth, error) {

	applicationInfo, err := a.applicationService.AppInfo(ctx, namespaceId, appId)
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

	resList := make([]*application_model.ApplicationAuth, 0, len(list))
	for _, auth := range list {

		operatorName := ""
		if userInfo, ok := userInfoMaps[auth.Operator]; ok {
			operatorName = userInfo.NickName
		}

		stat, err := a.applicationAuthStatStore.Get(ctx, auth.Id)
		if err != nil {
			return nil, err
		}
		version, err := a.applicationAuthVersionStore.Get(ctx, stat.VersionID)
		if err != nil {
			return nil, err
		}

		authModel := &application_model.ApplicationAuth{
			ApplicationAuth: auth,
			Operator:        operatorName,
			Config:          version.Config,
		}

		authModel.ParamPosition = auth.Position
		authModel.ParamName = auth.TokenName
		authModel.ParamInfo = a.getAuthParamInfo(authModel)
		authModel.RuleInfo = a.getRuleInfo(authModel)
		resList = append(resList, authModel)

	}
	return resList, nil
}

func (a *applicationAuthService) Create(ctx context.Context, namespaceId, userId int, appId string, input *application_dto.ApplicationAuthInput) error {
	authDriver := a.driverManager.GetDriver(input.Driver)
	if err := authDriver.CheckInput(input.Config); err != nil {
		return err
	}

	if input.ExpireTime > 0 && input.ExpireTime < time.Now().Unix() {
		return errors.New("过期时间不能小于当前时间")
	}

	applicationInfo, err := a.applicationService.AppInfo(ctx, namespaceId, appId)
	if err != nil {
		return err
	}

	if applicationInfo.IdStr == anonymousIds {
		return errors.New("匿名应用不能添加鉴权信息")
	}

	t := time.Now()
	applicationAuth := &application_entry.ApplicationAuth{
		Uuid:          uuid.New(),
		Namespace:     namespaceId,
		Application:   applicationInfo.Id,
		IsTransparent: input.IsTransparent,
		Driver:        input.Driver,
		Position:      input.Position,
		TokenName:     input.TokenName,
		ExpireTime:    input.ExpireTime,
		Operator:      userId,
		CreateTime:    t,
		UpdateTime:    t,
	}

	return a.applicationAuthStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = a.applicationAuthStore.Insert(txCtx, applicationAuth); err != nil {
			return err
		}

		config := application_entry.ApplicationAuthVersionConfig{
			Config: string(input.Config),
		}

		if err = a.applicationAuthHistoryStore.HistoryAdd(txCtx, namespaceId, applicationAuth.Id, &application_entry.ApplicationAuthHistoryInfo{
			Auth:   *applicationAuth,
			Config: config,
		}, userId); err != nil {
			return err
		}

		applicationAuthVersion := &application_entry.ApplicationAuthVersion{
			ApplicationAuthID:            applicationAuth.Id,
			NamespaceID:                  namespaceId,
			ApplicationAuthVersionConfig: config,
			Operator:                     userId,
			CreateTime:                   t,
		}

		if err = a.applicationAuthVersionStore.Save(txCtx, applicationAuthVersion); err != nil {
			return err
		}

		stat := &application_entry.ApplicationAuthStat{
			ApplicationAuthId: applicationAuthVersion.ApplicationAuthID,
			VersionID:         applicationAuthVersion.Id,
		}

		return a.applicationAuthStatStore.Save(txCtx, stat)
	})

}

func (a *applicationAuthService) Update(ctx context.Context, namespaceId, userId int, appId, uuidStr string, input *application_dto.ApplicationAuthInput) error {
	_, err := a.applicationService.AppInfo(ctx, namespaceId, appId)
	if err != nil {
		return err
	}

	authDriver := a.driverManager.GetDriver(input.Driver)
	if err = authDriver.CheckInput(input.Config); err != nil {
		return err
	}

	auth, err := a.applicationAuthStore.GetByUUID(ctx, uuidStr)
	if err != nil {
		return err
	}
	isUpdateVersion := false
	if auth.ExpireTime != input.ExpireTime {
		isUpdateVersion = true
	} else if auth.IsTransparent != input.IsTransparent {
		isUpdateVersion = true
	} else if auth.TokenName != input.TokenName {
		isUpdateVersion = true
	} else if auth.Position != input.Position {
		isUpdateVersion = true
	} else if auth.Driver != input.Driver {
		isUpdateVersion = true
	}

	//获取当前版本
	stat, err := a.applicationAuthStatStore.Get(ctx, auth.Id)
	if err != nil {
		return err
	}

	oldVersion, err := a.applicationAuthVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return err
	}

	if oldVersion.Config != string(input.Config) {
		isUpdateVersion = true
	}

	t := time.Now()

	auth.IsTransparent = input.IsTransparent
	auth.ExpireTime = input.ExpireTime
	auth.Operator = userId
	auth.UpdateTime = t
	auth.Driver = input.Driver
	auth.Position = input.Position
	auth.TokenName = input.TokenName

	return a.applicationAuthStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = a.applicationAuthStore.Save(txCtx, auth); err != nil {
			return err
		}

		newConfig := application_entry.ApplicationAuthVersionConfig{
			Config: string(input.Config),
		}

		if isUpdateVersion {
			applicationAuthVersion := &application_entry.ApplicationAuthVersion{
				ApplicationAuthID:            auth.Id,
				NamespaceID:                  namespaceId,
				ApplicationAuthVersionConfig: newConfig,
				Operator:                     userId,
				CreateTime:                   t,
			}

			if err = a.applicationAuthVersionStore.Save(txCtx, applicationAuthVersion); err != nil {
				return err
			}

			stat = &application_entry.ApplicationAuthStat{
				ApplicationAuthId: applicationAuthVersion.ApplicationAuthID,
				VersionID:         applicationAuthVersion.Id,
			}

			return a.applicationAuthStatStore.Save(txCtx, stat)
		}

		return a.applicationAuthHistoryStore.HistoryEdit(txCtx, namespaceId, auth.Id, &application_entry.ApplicationAuthHistoryInfo{
			Auth:   *auth,
			Config: oldVersion.ApplicationAuthVersionConfig,
		}, &application_entry.ApplicationAuthHistoryInfo{
			Auth:   *auth,
			Config: newConfig,
		}, userId)
	})
}

func (a *applicationAuthService) Info(ctx context.Context, namespaceId int, appId, uuid string) (*application_model.ApplicationAuth, error) {
	_, err := a.applicationService.AppInfo(ctx, namespaceId, appId)
	if err != nil {
		return nil, err
	}

	auth, err := a.applicationAuthStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	stat, err := a.applicationAuthStatStore.Get(ctx, auth.Id)
	if err != nil {
		return nil, err
	}

	version, err := a.applicationAuthVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return nil, err
	}

	resAuth := &application_model.ApplicationAuth{
		ApplicationAuth: auth,
		Operator:        "",
		Config:          version.Config,
	}
	return resAuth, nil
}

func (a *applicationAuthService) Delete(ctx context.Context, namespaceId, userId int, uuid string) error {

	auth, err := a.applicationAuthStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}

	//获取当前版本信息
	stat, err := a.applicationAuthStatStore.Get(ctx, auth.Id)
	if err != nil {
		return err
	}

	version, err := a.applicationAuthVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return err
	}

	return a.applicationAuthStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = a.applicationAuthStore.Delete(ctx, auth.Id); err != nil {
			return err
		}
		delMap := make(map[string]interface{})
		delMap["`kind`"] = "application_auth"
		delMap["`target`"] = auth.Id

		if _, err = a.applicationAuthStatStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = a.applicationAuthVersionStore.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		for _, clusterInfo := range clusters {
			delMap["`cluster`"] = clusterInfo.Id
			if _, err = a.applicationAuthRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}

			if err = a.applicationAuthPublishStore.DeleteByClusterIdAppId(txCtx, clusterInfo.Id, auth.Application); err != nil {
				return err
			}
		}

		return a.applicationAuthHistoryStore.HistoryDelete(txCtx, namespaceId, auth.Id, application_entry.ApplicationAuthHistoryInfo{
			Auth:   *auth,
			Config: version.ApplicationAuthVersionConfig,
		}, userId)
	})
}

func (a *applicationAuthService) Online(ctx context.Context, namespaceId, userId, clusterId, applicationId int) ([]*application_model.ApplicationAuth, error) {
	applicationAuths, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return nil, err
	}
	t := time.Now()

	authPublishList := make([]*application_entry.ApplicationAuthPublish, 0)
	list := make([]*application_model.ApplicationAuth, 0, len(applicationAuths))
	for _, auth := range applicationAuths {
		//获取当前版本信息
		stat, err := a.applicationAuthStatStore.Get(ctx, auth.Id)
		if err != nil {
			return nil, err
		}
		version, err := a.applicationAuthVersionStore.Get(ctx, stat.VersionID)
		if err != nil {
			return nil, err
		}

		list = append(list, &application_model.ApplicationAuth{
			ApplicationAuth: auth,
			Config:          version.Config,
		})
		runtime := &application_entry.ApplicationAuthRuntime{
			ClusterId:         clusterId,
			ApplicationAuthId: auth.Id,
			NamespaceId:       namespaceId,
			VersionId:         version.Id,
			IsOnline:          true,
			Operator:          userId,
			CreateTime:        t,
			UpdateTime:        t,
		}
		if err = a.applicationAuthRuntimeStore.Save(ctx, runtime); err != nil {
			return nil, err
		}
		authPublishList = append(authPublishList, &application_entry.ApplicationAuthPublish{
			NamespaceId:     namespaceId,
			Cluster:         clusterId,
			Application:     applicationId,
			ApplicationAuth: auth.Id,
			Operator:        userId,
			CreateTime:      t,
		})
	}
	if len(authPublishList) > 0 {
		return list, a.applicationAuthPublishStore.Reset(ctx, clusterId, applicationId, authPublishList)
	}
	return list, nil
}

func (a *applicationAuthService) Offline(ctx context.Context, clusterId, applicationId int) error {
	return a.applicationAuthPublishStore.DeleteByClusterIdAppId(ctx, clusterId, applicationId)
}

func (a *applicationAuthService) IsUpdate(ctx context.Context, clusterId, applicationId int) (bool, error) {
	applicationAuths, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return false, err
	}
	//查询发布过的鉴权信息，如果少了或者多了 则可以更新
	authPublishes, err := a.applicationAuthPublishStore.GetList(ctx, clusterId, applicationId)
	if err != nil {
		return false, err
	}
	if len(authPublishes) != len(applicationAuths) {
		return true, nil
	}

	for _, auth := range applicationAuths {

		stat, err := a.applicationAuthStatStore.Get(ctx, auth.Id)
		if err != nil {
			return false, err
		}

		//获取新版本
		newVersion, err := a.applicationAuthVersionStore.Get(ctx, stat.VersionID)
		if err != nil {
			return false, err
		}

		//获取运行的版本
		runtime, err := a.applicationAuthRuntimeStore.GetForCluster(ctx, auth.Id, clusterId)
		if err != nil {
			return false, err
		}

		runtimeVersion, err := a.applicationAuthVersionStore.Get(ctx, runtime.VersionId)
		if err != nil {
			return false, err
		}

		//对比两个版本的版本ID是否是一致，不是一致可以更新
		if newVersion.Id != runtimeVersion.Id {
			return true, nil
		}
	}
	return false, nil
}

func (a *applicationAuthService) GetListByApplicationId(ctx context.Context, applicationId int) ([]*application_model.ApplicationAuth, error) {
	list, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return nil, err
	}
	resList := make([]*application_model.ApplicationAuth, 0, len(list))
	for _, auth := range list {
		authModel := &application_model.ApplicationAuth{
			ApplicationAuth: auth,
		}
		resList = append(resList, authModel)
	}
	return resList, nil
}

func (a *applicationAuthService) GetDriversRender() []*driverInfo.DriverInfo {
	return a.driverManager.List()
}

func (a *applicationAuthService) getAuthParamInfo(auth *application_model.ApplicationAuth) string {
	return a.driverManager.GetDriver(auth.Driver).GetAuthListInfo([]byte(auth.Config))
}

func (a *applicationAuthService) getRuleInfo(auth *application_model.ApplicationAuth) string {
	switch auth.Driver {
	case application_model.AuthDriverBasic:
		return fmt.Sprintf("%s-%s", application_model.AuthDriverBasic, "")
	case application_model.AuthDriverApikey:
		return fmt.Sprintf("%s-%s", application_model.AuthDriverApikey, "")
	case application_model.AuthDriverAkSk:
		return fmt.Sprintf("%s-%s", application_model.AuthDriverAkSk, "")
	case application_model.AuthDriverJwt:
		return fmt.Sprintf("%s-%s", application_model.AuthDriverJwt, "")
	default:
		return ""
	}
}
