package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	driverInfo "github.com/eolinker/apinto-dashboard/driver-manager"
	"github.com/eolinker/apinto-dashboard/driver-manager/driver"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"time"
)

type IApplicationAuthService interface {
	GetList(ctx context.Context, namespaceId int, appId string) ([]*model.ApplicationAuth, error)
	Create(ctx context.Context, namespaceId, userId int, appId string, input *dto.ApplicationAuthInput) error
	Update(ctx context.Context, namespaceId, userId int, appId, uuid string, input *dto.ApplicationAuthInput) error
	Delete(ctx context.Context, namespaceId, userId int, uuid string) error
	Info(ctx context.Context, namespaceId int, appId, uuid string) (*model.ApplicationAuth, error)
	online(ctx context.Context, namespaceId, userId, clusterId, applicationId int) ([]*model.ApplicationAuth, error)
	offline(ctx context.Context, clusterId, applicationId int) error
	isUpdate(ctx context.Context, clusterId, applicationId int) (bool, error)
	getListByApplicationId(ctx context.Context, applicationId int) ([]*model.ApplicationAuth, error)
	GetDriversRender() []*driverInfo.DriverInfo
	GetDriver(driver string) driver.IAuthDriver
}

type applicationAuthService struct {
	applicationAuthStore        store.IApplicationAuthStore
	applicationAuthVersionStore store.IApplicationAuthVersionStore
	applicationAuthStatStore    store.IApplicationAuthStatStore
	applicationAuthRuntimeStore store.IApplicationAuthRuntimeStore
	applicationAuthPublishStore store.IApplicationAuthPublishStore
	applicationAuthHistoryStore store.IApplicationAuthHistoryStore
	applicationService          IApplicationService
	clusterService              IClusterService
	userInfoService             IUserInfoService
	driverManager               driverInfo.IAuthDriverManager
}

func (a *applicationAuthService) GetDriver(driver string) driver.IAuthDriver {
	return a.driverManager.GetDriver(driver)
}

func newApplicationAuth() IApplicationAuthService {
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

func (a *applicationAuthService) GetList(ctx context.Context, namespaceId int, appId string) ([]*model.ApplicationAuth, error) {

	applicationInfo, err := a.applicationService.AppInfo(ctx, namespaceId, appId)
	if err != nil {
		return nil, err
	}

	list, err := a.applicationAuthStore.GetListByApplication(ctx, applicationInfo.Id)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(list, func(t *entry.ApplicationAuth) int {
		return t.Operator
	})

	userInfoMaps, _ := a.userInfoService.GetUserInfoMaps(ctx, userIds...)

	resList := make([]*model.ApplicationAuth, 0, len(list))
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

		authModel := &model.ApplicationAuth{
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

func (a *applicationAuthService) Create(ctx context.Context, namespaceId, userId int, appId string, input *dto.ApplicationAuthInput) error {
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
	applicationAuth := &entry.ApplicationAuth{
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

		config := entry.ApplicationAuthVersionConfig{
			Config: string(input.Config),
		}

		if err = a.applicationAuthHistoryStore.HistoryAdd(txCtx, namespaceId, applicationAuth.Id, &entry.ApplicationAuthHistoryInfo{
			Auth:   *applicationAuth,
			Config: config,
		}, userId); err != nil {
			return err
		}

		applicationAuthVersion := &entry.ApplicationAuthVersion{
			ApplicationAuthID:            applicationAuth.Id,
			NamespaceID:                  namespaceId,
			ApplicationAuthVersionConfig: config,
			Operator:                     userId,
			CreateTime:                   t,
		}

		if err = a.applicationAuthVersionStore.Save(txCtx, applicationAuthVersion); err != nil {
			return err
		}

		stat := &entry.ApplicationAuthStat{
			ApplicationAuthId: applicationAuthVersion.ApplicationAuthID,
			VersionID:         applicationAuthVersion.Id,
		}

		return a.applicationAuthStatStore.Save(txCtx, stat)
	})

}

func (a *applicationAuthService) Update(ctx context.Context, namespaceId, userId int, appId, uuidStr string, input *dto.ApplicationAuthInput) error {
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

		newConfig := entry.ApplicationAuthVersionConfig{
			Config: string(input.Config),
		}

		if isUpdateVersion {
			applicationAuthVersion := &entry.ApplicationAuthVersion{
				ApplicationAuthID:            auth.Id,
				NamespaceID:                  namespaceId,
				ApplicationAuthVersionConfig: newConfig,
				Operator:                     userId,
				CreateTime:                   t,
			}

			if err = a.applicationAuthVersionStore.Save(txCtx, applicationAuthVersion); err != nil {
				return err
			}

			stat = &entry.ApplicationAuthStat{
				ApplicationAuthId: applicationAuthVersion.ApplicationAuthID,
				VersionID:         applicationAuthVersion.Id,
			}

			return a.applicationAuthStatStore.Save(txCtx, stat)
		}

		return a.applicationAuthHistoryStore.HistoryEdit(txCtx, namespaceId, auth.Id, &entry.ApplicationAuthHistoryInfo{
			Auth:   *auth,
			Config: oldVersion.ApplicationAuthVersionConfig,
		}, &entry.ApplicationAuthHistoryInfo{
			Auth:   *auth,
			Config: newConfig,
		}, userId)
	})
}

func (a *applicationAuthService) Info(ctx context.Context, namespaceId int, appId, uuid string) (*model.ApplicationAuth, error) {
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

	resAuth := &model.ApplicationAuth{
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

		for _, cluster := range clusters {
			delMap["`cluster`"] = cluster.Id
			if _, err = a.applicationAuthRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}

			if err = a.applicationAuthPublishStore.DeleteByClusterIdAppId(txCtx, cluster.Id, auth.Application); err != nil {
				return err
			}
		}

		return a.applicationAuthHistoryStore.HistoryDelete(txCtx, namespaceId, auth.Id, entry.ApplicationAuthHistoryInfo{
			Auth:   *auth,
			Config: version.ApplicationAuthVersionConfig,
		}, userId)
	})
}

func (a *applicationAuthService) online(ctx context.Context, namespaceId, userId, clusterId, applicationId int) ([]*model.ApplicationAuth, error) {
	applicationAuths, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return nil, err
	}
	t := time.Now()

	authPublishList := make([]*entry.ApplicationAuthPublish, 0)
	list := make([]*model.ApplicationAuth, 0, len(applicationAuths))
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

		list = append(list, &model.ApplicationAuth{
			ApplicationAuth: auth,
			Config:          version.Config,
		})
		runtime := &entry.ApplicationAuthRuntime{
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
		authPublishList = append(authPublishList, &entry.ApplicationAuthPublish{
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

func (a *applicationAuthService) offline(ctx context.Context, clusterId, applicationId int) error {
	return a.applicationAuthPublishStore.DeleteByClusterIdAppId(ctx, clusterId, applicationId)
}

func (a *applicationAuthService) isUpdate(ctx context.Context, clusterId, applicationId int) (bool, error) {
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

func (a *applicationAuthService) getListByApplicationId(ctx context.Context, applicationId int) ([]*model.ApplicationAuth, error) {
	list, err := a.applicationAuthStore.GetListByApplication(ctx, applicationId)
	if err != nil {
		return nil, err
	}
	resList := make([]*model.ApplicationAuth, 0, len(list))
	for _, auth := range list {
		authModel := &model.ApplicationAuth{
			ApplicationAuth: auth,
		}
		resList = append(resList, authModel)
	}
	return resList, nil
}

func (a *applicationAuthService) GetDriversRender() []*driverInfo.DriverInfo {
	return a.driverManager.List()
}

func (a *applicationAuthService) getAuthParamInfo(auth *model.ApplicationAuth) string {
	return a.driverManager.GetDriver(auth.Driver).GetAuthListInfo([]byte(auth.Config))
}

func (a *applicationAuthService) getRuleInfo(auth *model.ApplicationAuth) string {
	switch auth.Driver {
	case model.AuthDriverBasic:
		return fmt.Sprintf("%s-%s", model.AuthDriverBasic, "")
	case model.AuthDriverApikey:
		return fmt.Sprintf("%s-%s", model.AuthDriverApikey, "")
	case model.AuthDriverAkSk:
		return fmt.Sprintf("%s-%s", model.AuthDriverAkSk, "")
	case model.AuthDriverJwt:
		return fmt.Sprintf("%s-%s", model.AuthDriverJwt, "")
	default:
		return ""
	}
}
