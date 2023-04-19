package api_service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	apiservice "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/api/api-dto"
	apientry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	apiStore "github.com/eolinker/apinto-dashboard/modules/api/store"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/group"
	"github.com/eolinker/apinto-dashboard/modules/group/group-entry"
	"github.com/eolinker/apinto-dashboard/modules/group/group-model"
	"github.com/eolinker/apinto-dashboard/modules/group/group-service"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/openapp"
	"github.com/eolinker/apinto-dashboard/modules/openapp/open-app-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-model"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type apiService struct {
	apiStore   apiStore.IAPIStore
	apiStat    apiStore.IAPIStatStore
	apiVersion apiStore.IAPIVersionStore
	apiRuntime apiStore.IAPIRuntimeStore
	quoteStore quote_store.IQuoteStore
	apiHistory apiStore.IApiHistoryStore

	service          upstream.IService
	commonGroup      group.ICommonGroupService
	clusterService   cluster.IClusterService
	namespaceService namespace.INamespaceService
	apintoClient     cluster.IApintoClient
	userInfoService  user.IUserInfoService
	extAppService    openapp.IExternalApplicationService
	apiManager       apiservice.IAPIDriverManager

	pluginTemplateService plugin_template.IPluginTemplateService

	lockService    locker_service.IAsynLockService
	importApiCache IImportApiCache
	batchApiCache  IBatchOnlineApiTaskCache
}

func NewAPIService() apiservice.IAPIService {
	as := &apiService{}
	bean.Autowired(&as.apiStore)
	bean.Autowired(&as.apiStat)
	bean.Autowired(&as.apiVersion)
	bean.Autowired(&as.apiRuntime)
	bean.Autowired(&as.quoteStore)

	bean.Autowired(&as.service)
	bean.Autowired(&as.commonGroup)
	bean.Autowired(&as.clusterService)
	bean.Autowired(&as.namespaceService)
	bean.Autowired(&as.apintoClient)
	bean.Autowired(&as.apiManager)
	bean.Autowired(&as.apiHistory)
	bean.Autowired(&as.userInfoService)
	bean.Autowired(&as.extAppService)

	bean.Autowired(&as.lockService)
	bean.Autowired(&as.importApiCache)
	bean.Autowired(&as.batchApiCache)
	bean.Autowired(&as.pluginTemplateService)

	return as
}

func (a *apiService) GetAPICountByGroupUUID(ctx context.Context, namespaceID int, groupUUID string) int64 {
	count, err := a.apiStore.GetCountByGroupID(ctx, namespaceID, groupUUID)
	if err != nil {
		log.Errorf("GetAPICountByGroupUUID-apiStore.GetCountByGroupID namespaceId:%d,groupUUid:%s err=%s", namespaceID, groupUUID, err.Error())
		return 0
	}
	return count
}

func (a *apiService) APICount(ctx context.Context, namespaceId int) (int64, error) {
	return a.apiStore.APICount(ctx, map[string]interface{}{
		"namespace": namespaceId,
	})
}

func (a *apiService) APIOnlineCount(ctx context.Context, namespaceId int) (int64, error) {
	return a.apiRuntime.OnlineCountByKind(ctx)
}

func (a *apiService) GetGroups(ctx context.Context, namespaceId int, parentUuid, queryName string) (*group_model.CommonGroupRoot, []*group_model.CommonGroupApi, error) {
	groups, err := a.commonGroup.GroupListAll(ctx, namespaceId, group_service.ApiName, group_service.ApiName)
	if err != nil {
		log.Errorf("GetGroups-commonGroup.GroupListAll namespaceId:%d,parentUuid:%s,queryName=%s, err=%s", namespaceId, parentUuid, queryName, err.Error())
		return nil, nil, err
	}

	apis := make([]*group_model.CommonGroupApi, 0)
	apisAll := make([]*group_model.CommonGroupApi, 0)
	//查询API
	//apisAll, err := a.GetAPIListByName(ctx, namespaceId, "")
	//if err != nil {
	//	return nil, nil, err
	//}
	//for _, apiService := range apisAll {
	//	if queryName != "" && strings.Count(strings.ToUpper(apiService.Name), strings.ToUpper(queryName)) > 0 {
	//		apis = append(apis, apiService)
	//	}
	//}

	if queryName == "" {
		return a.commonGroup.ToGroupRoot(ctx, namespaceId, parentUuid, groups, map[string]string{}), nil, nil
	}

	//查询API的目录直至跟目录
	groupUUIDS := common.SliceToSliceIds(apis, func(t *group_model.CommonGroupApi) string {
		return t.GroupUUID
	})

	uuidMaps := common.SliceToMap(groupUUIDS, func(t string) string {
		return t
	})
	groupsMaps := common.SliceToMap(groups, func(t *group_entry.CommonGroup) string {
		return t.Uuid
	})
	groupsIdMaps := common.SliceToMap(groups, func(t *group_entry.CommonGroup) int {
		return t.Id
	})

	groupsParentIdMaps := common.SliceToMapArray(groups, func(t *group_entry.CommonGroup) int {
		return t.ParentId
	})

	for _, groupInfo := range groups {
		//模糊搜索
		if strings.Count(strings.ToUpper(groupInfo.Name), strings.ToUpper(queryName)) > 0 {

			uuidMaps[groupInfo.Uuid] = groupInfo.Uuid
			//如果绝对相等，需要把改目录下的所有根目录也查询出来
			if groupInfo.Name == queryName {
				uuids := &[]string{}
				a.commonGroup.SubGroupUUIDS(groupsParentIdMaps, &group_model.CommonGroup{
					Group: groupInfo,
				}, uuids)
				for _, s := range *uuids {
					uuidMaps[s] = s
				}
			}
		}
	}

	//拿到API的目录以及跟目录的所有UUID
	outMapUUID := make(map[string]string)
	for _, groupUuid := range uuidMaps {
		a.commonGroup.ParentGroupV2(groupUuid, groupsMaps, groupsIdMaps, outMapUUID)
	}

	groupRoot := a.commonGroup.ToGroupRoot(ctx, namespaceId, parentUuid, groups, outMapUUID)

	apiAllMaps := common.SliceToMapArray(apisAll, func(t *group_model.CommonGroupApi) string {
		return t.GroupUUID
	})

	apiMaps := common.SliceToMap(apis, func(t *group_model.CommonGroupApi) string {
		return t.GroupUUID
	})

	resApis := &[]*group_model.CommonGroupApi{}
	a.subGroup(groupRoot.CommonGroup, apiAllMaps, apiMaps, resApis)

	*resApis = append(*resApis, apis...)

	return groupRoot, *resApis, err
}

func (a *apiService) subGroup(list []*group_model.CommonGroup, apiAllMaps map[string][]*group_model.CommonGroupApi, apiMaps map[string]*group_model.CommonGroupApi, apis *[]*group_model.CommonGroupApi) {
	if len(list) == 0 {
		return
	}
	for _, groupInfo := range list {
		if _, ok := apiMaps[groupInfo.Group.Uuid]; !ok {
			*apis = append(*apis, apiAllMaps[groupInfo.Group.Uuid]...)
		}
		a.subGroup(groupInfo.Subgroup, apiAllMaps, apiMaps, apis)
	}
}

func (a *apiService) GetAPIList(ctx context.Context, namespaceID int, groupUUID, searchName string, searchSources []string, pageNum, pageSize int) ([]*apimodel.APIListItem, int, error) {
	groupList := make([]string, 0)
	var err error
	//获取传入的groupUUID下包括子分组的所有UUID
	if groupUUID != "" {
		groupList, err = a.commonGroup.GroupUUIDS(ctx, namespaceID, group_service.ApiName, group_service.ApiName, groupUUID)
		if err != nil {
			return nil, 0, err
		}
		groupList = append(groupList, groupUUID)
	}

	//获取指定目录下所有API
	apis, total, err := a.apiStore.GetListPageByGroupIDs(ctx, namespaceID, pageNum, pageSize, groupList, searchSources, searchName)
	if err != nil {
		return nil, 0, err
	}
	apiList := make([]*apimodel.APIListItem, 0, len(apis))
	source := ""

	apiIds := make([]int, 0, len(apis))
	for _, api := range apis {
		apiIds = append(apiIds, api.Id)
	}
	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, 0, err
	}

	for _, api := range apis {
		version := versionMap[api.Id]

		if api.SourceType != enum.SourceSync {
			source = enum.GetSourceTitle(api.SourceType)
		} else {
			//如果是同步 source = 应用+标签
			appName, _ := a.extAppService.GetExtAppName(ctx, api.SourceID)
			source = appName
			if api.SourceLabel != "" {
				source = appName + "-" + api.SourceLabel
			}
		}

		isDelete, _ := a.isApiCanDelete(ctx, api.Id)

		item := &apimodel.APIListItem{
			GroupUUID:   api.GroupUUID,
			APIUUID:     api.UUID,
			APIName:     api.Name,
			Method:      version.Method,
			ServiceName: version.ServiceName,
			RequestPath: version.RequestPathLabel,
			Source:      source,
			UpdateTime:  api.UpdateTime,
			IsDelete:    isDelete,
		}

		apiList = append(apiList, item)
	}

	return apiList, total, nil
}

func (a *apiService) GetAPIInfo(ctx context.Context, namespaceID int, uuid string) (*apimodel.APIInfo, error) {
	api, err := a.apiStore.GetByUUID(ctx, namespaceID, uuid)
	if err != nil {
		log.Errorf("GetAPIInfo-apiStore.GetByUUID namespaceId:%d,uuid:%s,err=%s", namespaceID, uuid, err.Error())
		return nil, err
	}
	return &apimodel.APIInfo{API: api}, nil
}

func (a *apiService) GetAPIInfoById(ctx context.Context, id int) (*apimodel.APIInfo, error) {
	api, err := a.apiStore.Get(ctx, id)
	if err != nil {
		log.Errorf("GetAPIInfo-apiStore.GetByUUID id:%s,err=%s", id, err.Error())
		return nil, err
	}
	return &apimodel.APIInfo{API: api}, nil
}

func (a *apiService) GetAPIInfoByGroupUUID(ctx context.Context, namespaceID int, groupUUID string) ([]*apimodel.APIInfo, error) {
	apis, err := a.apiStore.GetListByGroupID(ctx, namespaceID, groupUUID)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &apimodel.APIInfo{API: api})
	}

	return list, nil
}

func (a *apiService) GetAPIInfoByUUIDS(ctx context.Context, namespaceID int, uuids []string) ([]*apimodel.APIInfo, error) {
	apis, err := a.apiStore.GetByUUIDs(ctx, namespaceID, uuids)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &apimodel.APIInfo{API: api})
	}

	return list, nil
}

func (a *apiService) GetAPIInfoByPath(ctx context.Context, namespaceID int, path string) ([]*apimodel.APIInfo, error) {
	apis, err := a.apiStore.GetByPath(ctx, namespaceID, path)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &apimodel.APIInfo{API: api})
	}

	return list, nil
}

func (a *apiService) GetAPIListItemByUUIDS(ctx context.Context, namespaceID int, uuids []string) ([]*apimodel.APIListItem, error) {
	apis, err := a.apiStore.GetByUUIDs(ctx, namespaceID, uuids)
	if err != nil {
		return nil, err
	}

	apiIds := make([]int, 0, len(apis))
	for _, api := range apis {
		apiIds = append(apiIds, api.Id)
	}

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIListItem, 0, len(apis))
	for _, api := range apis {
		version := versionMap[api.Id]
		item := &apimodel.APIListItem{
			GroupUUID:   api.GroupUUID,
			APIUUID:     api.UUID,
			APIName:     api.Name,
			Method:      version.Method,
			ServiceName: version.ServiceName,
			RequestPath: version.RequestPathLabel,
			UpdateTime:  api.UpdateTime,
		}

		list = append(list, item)
	}

	return list, nil
}

func (a *apiService) GetAPIListItemAll(ctx context.Context, namespaceID int) ([]*apimodel.APIListItem, error) {
	apis, err := a.apiStore.GetListByName(ctx, namespaceID, "")
	if err != nil {
		return nil, err
	}

	apiIds := make([]int, 0, len(apis))
	for _, api := range apis {
		apiIds = append(apiIds, api.Id)
	}

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIListItem, 0, len(apis))
	for _, api := range apis {
		version := versionMap[api.Id]
		item := &apimodel.APIListItem{
			GroupUUID:   api.GroupUUID,
			APIUUID:     api.UUID,
			APIName:     api.Name,
			Method:      version.Method,
			ServiceName: version.ServiceName,
			RequestPath: version.RequestPathLabel,
			UpdateTime:  api.UpdateTime,
		}

		list = append(list, item)
	}

	return list, nil
}

func (a *apiService) GetAPIInfoAll(ctx context.Context, namespaceID int) ([]*apimodel.APIInfo, error) {
	apis, err := a.apiStore.GetListByName(ctx, namespaceID, "")
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &apimodel.APIInfo{API: api})
	}

	return list, nil
}

// GetAPIsForSync 同步api时使用
func (a *apiService) GetAPIsForSync(ctx context.Context, namespaceID int) ([]*apimodel.APIVersionInfo, error) {
	apis, err := a.apiStore.GetListByName(ctx, namespaceID, "")
	if err != nil {
		return nil, err
	}

	apiIds := make([]int, 0, len(apis))
	for _, api := range apis {
		apiIds = append(apiIds, api.Id)
	}

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIVersionInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &apimodel.APIVersionInfo{
			Api:     api,
			Version: versionMap[api.Id],
		})
	}

	return list, nil
}

func (a *apiService) getAPIVersions(ctx context.Context, apiIds []int) (map[int]*apientry.APIVersion, error) {
	versions, err := a.apiVersion.GetAPIVersionByApiIds(ctx, apiIds)
	if err != nil {
		return nil, err
	}
	return common.SliceToMap(versions, func(t *apientry.APIVersion) int {
		return t.ApiID
	}), nil
}

func (a *apiService) GetAPIVersionInfo(ctx context.Context, namespaceID int, uuid string) (*apimodel.APIVersionInfo, error) {
	api, err := a.apiStore.GetByUUID(ctx, namespaceID, uuid)
	if err != nil {
		return nil, err
	}

	version, err := a.GetLatestAPIVersion(ctx, api.Id)
	if err != nil {
		return nil, err
	}

	info := &apimodel.APIVersionInfo{
		Api:     api,
		Version: version,
	}

	return info, nil
}

func (a *apiService) CreateAPI(ctx context.Context, namespaceID int, operator int, input *api_dto.APIInfo) error {

	if err := a.CheckAPIReduplicative(ctx, namespaceID, "", input); err != nil {
		return err
	}

	if input.UUID == "" {
		input.UUID = uuid.New()
	}

	input.UUID = strings.ToLower(input.UUID)

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.ApiName,
	})

	isExist, err := a.commonGroup.IsGroupExist(ctx, input.GroupUUID)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("group doesn't. group_uuid:%s ", input.GroupUUID)
	}

	serviceID, err := a.service.GetServiceIDByName(ctx, namespaceID, input.ServiceName)
	if err != nil {
		return err
	}

	var templateID int
	if input.TemplateUUID != "" {
		templateInfo, err := a.pluginTemplateService.GetByUUID(ctx, namespaceID, input.TemplateUUID)
		if err != nil {
			return err
		}
		templateID = templateInfo.Id
	}

	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		apiInfo := &apientry.API{
			NamespaceId:      namespaceID,
			UUID:             input.UUID,
			GroupUUID:        input.GroupUUID,
			Name:             input.ApiName,
			RequestPath:      input.RequestPath,
			RequestPathLabel: input.RequestPathLabel,
			SourceType:       enum.SourceSelfBuild,
			SourceID:         -1,
			SourceLabel:      "",
			Desc:             input.Desc,
			Operator:         operator,
			CreateTime:       t,
			UpdateTime:       t,
		}
		if err = a.apiStore.Save(txCtx, apiInfo); err != nil {
			return err
		}

		//添加版本信息
		apiVersionInfo := &apientry.APIVersion{
			ApiID:       apiInfo.Id,
			NamespaceID: namespaceID,
			APIVersionConfig: apientry.APIVersionConfig{
				Driver:           input.Driver,
				RequestPath:      input.RequestPath,
				RequestPathLabel: input.RequestPathLabel,
				ServiceID:        serviceID,
				ServiceName:      input.ServiceName,
				TemplateID:       templateID,
				TemplateUUID:     input.TemplateUUID,
				Method:           input.Method,
				ProxyPath:        input.ProxyPath,
				Timeout:          input.Timeout,
				Retry:            input.Retry,
				EnableWebsocket:  input.EnableWebsocket,
				Match:            input.Match,
				Header:           input.Header,
			},
			Operator:   operator,
			CreateTime: t,
		}

		if err = a.apiVersion.Save(txCtx, apiVersionInfo); err != nil {
			return err
		}

		if err = a.apiHistory.HistoryAdd(txCtx, namespaceID, apiInfo.Id, &apientry.ApiHistoryInfo{
			Api:    *apiInfo,
			Config: apiVersionInfo.APIVersionConfig,
		}, operator); err != nil {
			return err
		}

		stat := &apientry.APIStat{
			APIID:     apiInfo.Id,
			VersionID: apiVersionInfo.Id,
		}

		//添加版本关联原表信息
		if err = a.apiStat.Save(txCtx, stat); err != nil {
			return err
		}

		//更新所引用的插件模板
		if templateID != 0 {
			err = a.quoteStore.Set(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quote_entry.QuoteTargetKindTypePluginTemplate, templateID)
			if err != nil {
				return err
			}
		}

		//quote更新所引用的服务

		return a.quoteStore.Set(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quote_entry.QuoteTargetKindTypeService, serviceID)
	})

}

func (a *apiService) UpdateAPI(ctx context.Context, namespaceID int, operator int, input *api_dto.APIInfo) error {
	if err := a.CheckAPIReduplicative(ctx, namespaceID, input.UUID, input); err != nil {
		return err
	}

	isExist, err := a.commonGroup.IsGroupExist(ctx, input.GroupUUID)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("group doesn't. group_uuid:%s ", input.GroupUUID)
	}

	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceID, input.UUID)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceID, input.UUID)
	if err != nil {
		return err
	}

	serviceID, err := a.service.GetServiceIDByName(ctx, namespaceID, input.ServiceName)
	if err != nil {
		return err
	}

	var templateID int
	if input.TemplateUUID != "" {
		templateInfo, err := a.pluginTemplateService.GetByUUID(ctx, namespaceID, input.TemplateUUID)
		if err != nil {
			return err
		}
		templateID = templateInfo.Id
	}

	stat, err := a.apiStat.Get(ctx, apiInfo.Id)
	if err != nil {
		return err
	}
	currentVersion, err := a.apiVersion.Get(ctx, stat.VersionID)
	if err != nil {
		return err
	}

	t := time.Now()

	oldValue := apientry.ApiHistoryInfo{
		Api:    *apiInfo,
		Config: currentVersion.APIVersionConfig,
	}

	apiInfo.Desc = input.Desc
	apiInfo.GroupUUID = input.GroupUUID
	apiInfo.Name = input.ApiName
	apiInfo.RequestPath = input.RequestPath
	apiInfo.RequestPathLabel = input.RequestPathLabel
	apiInfo.Operator = operator
	apiInfo.UpdateTime = t

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.ApiName,
	})

	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		//修改基础数据
		if _, err = a.apiStore.Update(txCtx, apiInfo); err != nil {
			return err
		}

		latestVersionConfig := apientry.APIVersionConfig{
			Driver:           input.Driver,
			RequestPath:      input.RequestPath,
			RequestPathLabel: input.RequestPathLabel,
			ServiceID:        serviceID,
			ServiceName:      input.ServiceName,
			TemplateID:       templateID,
			TemplateUUID:     input.TemplateUUID,
			Method:           input.Method,
			ProxyPath:        input.ProxyPath,
			Timeout:          input.Timeout,
			Retry:            input.Retry,
			EnableWebsocket:  input.EnableWebsocket,
			Match:            input.Match,
			Header:           input.Header,
		}
		//判断配置信息是否有更新
		if a.isAPIVersionConfChange(latestVersionConfig, currentVersion.APIVersionConfig) {
			apiVersionInfo := &apientry.APIVersion{
				ApiID:            apiInfo.Id,
				NamespaceID:      apiInfo.NamespaceId,
				APIVersionConfig: latestVersionConfig,
				Operator:         operator,
				CreateTime:       t,
			}
			if err = a.apiVersion.Save(txCtx, apiVersionInfo); err != nil {
				return err
			}
			//添加版本关联原表信息
			stat = &apientry.APIStat{
				APIID:     apiInfo.Id,
				VersionID: apiVersionInfo.Id,
			}
			if err = a.apiStat.Save(txCtx, stat); err != nil {
				return err
			}

			//更新所引用的插件模板
			if currentVersion.TemplateID != templateID {
				if templateID != 0 {
					err = a.quoteStore.Set(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quote_entry.QuoteTargetKindTypePluginTemplate, templateID)
					if err != nil {
						return err
					}
				} else {
					err = a.quoteStore.DelSourceTarget(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quote_entry.QuoteTargetKindTypePluginTemplate)
					if err != nil {
						return err
					}
				}
			}

			//quote更新所引用的服务
			if err = a.quoteStore.Set(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quote_entry.QuoteTargetKindTypeService, serviceID); err != nil {
				return err
			}
		}

		newValue := apientry.ApiHistoryInfo{
			Api:    *apiInfo,
			Config: latestVersionConfig,
		}

		return a.apiHistory.HistoryEdit(txCtx, namespaceID, apiInfo.Id, &oldValue, &newValue, operator)

	})

}

func (a *apiService) DeleteAPI(ctx context.Context, namespaceId, operator int, uuid string) error {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	isDel, err := a.isApiCanDelete(ctx, apiInfo.Id)
	if err != nil {
		return err
	}
	if !isDel {
		return errors.New("API is online. ")
	}

	version, err := a.GetLatestAPIVersion(ctx, apiInfo.Id)

	OldValue := apientry.ApiHistoryInfo{
		Api:    *apiInfo,
		Config: version.APIVersionConfig,
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuid,
		Name: apiInfo.Name,
	})

	err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = a.apiStore.Delete(txCtx, apiInfo.Id); err != nil {
			return err
		}
		delMap := make(map[string]interface{})
		delMap["`kind`"] = "apiService"
		delMap["`target`"] = apiInfo.Id
		if _, err = a.apiStat.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = a.apiVersion.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if _, err = a.apiRuntime.DeleteWhere(txCtx, delMap); err != nil {
			return err
		}

		if err = a.apiHistory.HistoryDelete(txCtx, namespaceId, apiInfo.Id, &OldValue, operator); err != nil {
			return err
		}

		return a.quoteStore.DelBySource(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI)
	})
	if err != nil {
		return err
	}

	a.lockService.DeleteLock(locker_service.LockNameAPI, apiInfo.Id)
	return nil
}

func (a *apiService) BatchOnline(ctx context.Context, namespaceId int, operator int, onlineToken string) ([]*apimodel.BatchListItem, error) {
	//判断uuid和operator是一致的
	key := a.batchApiCache.Key(onlineToken)
	task, err := a.batchApiCache.Get(ctx, key)
	//篡改审计日志的请求body
	if err != nil {
		return nil, err
	}
	ginContext, ok := ctx.(*gin.Context)
	if ok {
		ginContext.Set("logBody", string(task.Data))
	}

	//若提交上线的操作人和成功检测的操作人不一致 则报错
	if task.Operator != operator {
		return nil, errors.New("operator is invalid. ")
	}

	conf := new(apimodel.APIBatchConf)
	_ = json.Unmarshal(task.Data, conf)

	apiList := make([]*apientry.API, 0, len(conf.ApiUUIDs))

	errorGroup, _ := errgroup.WithContext(ctx)
	errorGroup.Go(func() error {
		//确认所有apiUUID，clusterName均存在
		for _, uid := range conf.ApiUUIDs {
			api, err := a.apiStore.GetByUUID(ctx, namespaceId, uid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("apiService doesn't exist. uuid:%s", uid)
				}
				return err
			}
			apiList = append(apiList, api)
		}
		return nil
	})

	clusterList := make([]*cluster_model.Cluster, 0, len(conf.ClusterNames))
	errorGroup.Go(func() error {
		clusters, err := a.clusterService.QueryListByNamespaceId(ctx, namespaceId)
		if err != nil {
			return err
		}
		clusterMap := common.SliceToMap(clusters, func(t *cluster_model.Cluster) string {
			return t.Name
		})
		for _, clusterName := range conf.ClusterNames {
			if clusterInfo, ok := clusterMap[clusterName]; ok {
				if clusterInfo.Status == 2 || clusterInfo.Status == 3 {
					return fmt.Errorf("cluster status is abnormal. cluster_name:%s", clusterName)
				}
				clusterList = append(clusterList, clusterInfo)
			} else {
				return fmt.Errorf("cluster doesn't exist. cluster_name:%s", clusterName)
			}

		}
		return nil
	})

	if err = errorGroup.Wait(); err != nil {
		return nil, err
	}

	//逐个处理api上线
	onlineList := make([]*apimodel.BatchListItem, 0, len(apiList)*len(clusterList))
	for _, api := range apiList {

		online, err := a.online(ctx, namespaceId, operator, api, clusterList)
		if err != nil && len(online) == 0 {
			return nil, err
		}
		onlineList = append(onlineList, online...)

		//err = a.lockService.Lock(locker_service.LockNameAPI, api.Id)
		//if err != nil {
		//	for _, clusterInfo := range clusterList {
		//		item := &apimodel.BatchListItem{
		//			APIName:    api.Name,
		//			ClusterEnv: fmt.Sprintf("%s_%s", clusterInfo.Name, clusterInfo.Env),
		//			Status:     false,
		//			Result:     err.Error(),
		//		}
		//		onlineList = append(onlineList, item)
		//	}
		//	a.lockService.Unlock(locker_service.LockNameAPI, api.Id)
		//	continue
		//}
		////确保api没被删除
		//_, err = a.apiStore.Get(ctx, api.Id)
		//if err != nil {
		//	//API被删除
		//	for _, clusterInfo := range clusterList {
		//		item := &apimodel.BatchListItem{
		//			APIName:    api.Name,
		//			ClusterEnv: fmt.Sprintf("%s_%s", clusterInfo.Name, clusterInfo.Env),
		//			Status:     false,
		//			Result:     err.Error(),
		//		}
		//		onlineList = append(onlineList, item)
		//	}
		//	a.lockService.Unlock(locker_service.LockNameAPI, api.Id)
		//	continue
		//}
		//
		//for _, clusterInfo := range clusterList {
		//	item := &apimodel.BatchListItem{
		//		APIName:    api.Name,
		//		ClusterEnv: fmt.Sprintf("%s_%s", clusterInfo.Name, clusterInfo.Env),
		//		Status:     true,
		//		Result:     "",
		//	}
		//
		//	//获取当前的版本
		//	runtime, err := a.apiRuntime.GetForCluster(ctx, api.Id, clusterInfo.Id)
		//	if err != nil && err != gorm.ErrRecordNotFound {
		//		item.Status = false
		//		item.Result = err.Error()
		//		onlineList = append(onlineList, item)
		//		continue
		//	}
		//
		//	err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		//		latest, err := a.GetLatestAPIVersion(ctx, api.Id)
		//		if err != nil {
		//			return err
		//		}
		//		//判断上游服务有没有上线
		//		if !a.service.IsOnline(ctx, clusterInfo.Id, latest.ServiceID) {
		//			item.Status = false
		//			item.Result = fmt.Sprintf("绑定的%s未上线到%s", latest.ServiceName, clusterInfo.Name)
		//			return nil
		//		}
		//
		//		if runtime != nil {
		//			current, err := a.apiVersion.Get(ctx, runtime.VersionID)
		//			if err != nil {
		//				return err
		//			}
		//
		//			//若api为已上线且无更新状态
		//			if runtime.IsOnline && !a.isAPIVersionConfChange(latest.APIVersionConfig, current.APIVersionConfig) {
		//				return nil
		//			}
		//		}
		//
		//		//发布到apinto
		//		client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
		//		if err != nil {
		//			item.Status = false
		//			item.Result = fmt.Sprintf("连接集群失败, err: %s", err.Error())
		//			return nil
		//		}
		//
		//		//封装router配置
		//		apiDriverInfo := a.GetAPIDriver(latest.Driver)
		//		routerConfig := apiDriverInfo.ToApinto(api.UUID, api.Desc, false, latest.Method, latest.RequestPath, latest.RequestPathLabel, latest.ProxyPath, strings.ToLower(latest.ServiceName), latest.Timeout, latest.Retry, latest.EnableWebsocket, latest.Match, latest.Header)
		//
		//		//未上线
		//		if runtime == nil {
		//			runtime = &apientry.APIRuntime{
		//				NamespaceId: namespaceId,
		//				ApiID:       api.Id,
		//				ClusterID:   clusterInfo.Id,
		//				VersionID:   latest.Id,
		//				IsOnline:    true,
		//				Disable:     false,
		//				Operator:    operator,
		//				CreateTime:  t,
		//				UpdateTime:  t,
		//			}
		//
		//			if err = a.apiRuntime.Insert(txCtx, runtime); err != nil {
		//				return err
		//			}
		//			if err = client.ForRouter().Create(*routerConfig); err != nil {
		//				item.Status = false
		//				item.Result = fmt.Sprintf("发送配置至集群失败, err: %s", err.Error())
		//			}
		//
		//		} else { //已下线或者待更新
		//			isOnline := runtime.IsOnline //保存旧状态
		//
		//			runtime.IsOnline = true
		//			runtime.UpdateTime = t
		//			runtime.VersionID = latest.Id
		//			runtime.Operator = operator
		//
		//			routerConfig.Disable = runtime.Disable
		//
		//			if err = a.apiRuntime.Save(txCtx, runtime); err != nil {
		//				return err
		//			}
		//
		//			//若原先是下线状态
		//			if !isOnline {
		//				if err = client.ForRouter().Create(*routerConfig); err != nil {
		//					item.Status = false
		//					item.Result = fmt.Sprintf("发送配置至集群失败, err: %s", err.Error())
		//				}
		//			}
		//
		//			if err = client.ForRouter().Update(api.UUID+"@router", *routerConfig); err != nil {
		//				item.Status = false
		//				item.Result = fmt.Sprintf("发送配置至集群失败, err: %s", err.Error())
		//			}
		//		}
		//		return nil
		//	})
		//	if err != nil {
		//		item.Status = false
		//		item.Result = err.Error()
		//	}
		//
		//	onlineList = append(onlineList, item)
		//}
		//
		//a.lockService.Unlock(locker_service.LockNameAPI, api.Id)
	}
	//编写操作记录
	logApiNameList := make([]string, 0, len(apiList))
	logCLNameList := make([]string, 0, len(clusterList))
	for _, api := range apiList {
		logApiNameList = append(logApiNameList, api.Name)
	}
	for _, cl := range clusterList {
		logCLNameList = append(logCLNameList, cl.Name)
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        strings.Join(logApiNameList, ","),
		ClusterName: strings.Join(logCLNameList, ","),
		PublishType: 1,
	})

	return onlineList, nil
}

func (a *apiService) online(ctx context.Context, namespaceId, operator int, api *apientry.API, clusterList []*cluster_model.Cluster) ([]*apimodel.BatchListItem, error) {
	t := time.Now()
	onlineList := make([]*apimodel.BatchListItem, 0)
	err := a.lockService.Lock(locker_service.LockNameAPI, api.Id)
	if err != nil {
		return nil, err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, api.Id)

	//确保api没被删除
	_, err = a.apiStore.Get(ctx, api.Id)
	if err != nil {
		//API被删除
		for _, clusterInfo := range clusterList {
			item := &apimodel.BatchListItem{
				APIName:    api.Name,
				ClusterEnv: fmt.Sprintf("%s_%s", clusterInfo.Name, clusterInfo.Env),
				Status:     false,
				Result:     err.Error(),
			}
			onlineList = append(onlineList, item)
		}
		return onlineList, nil
	}

	for _, clusterInfo := range clusterList {
		item := &apimodel.BatchListItem{
			APIName:    api.Name,
			ClusterEnv: fmt.Sprintf("%s_%s", clusterInfo.Name, clusterInfo.Env),
			Status:     true,
			Result:     "",
		}

		//获取当前的版本
		runtime, err := a.apiRuntime.GetForCluster(ctx, api.Id, clusterInfo.Id)
		if err != nil && err != gorm.ErrRecordNotFound {
			item.Status = false
			item.Result = err.Error()
			onlineList = append(onlineList, item)
			continue
		}

		err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
			latest, err := a.GetLatestAPIVersion(ctx, api.Id)
			if err != nil {
				return err
			}
			//判断上游服务有没有上线
			if !a.service.IsOnline(ctx, clusterInfo.Id, latest.ServiceID) {
				item.Status = false
				item.Result = fmt.Sprintf("绑定的%s未上线到%s", latest.ServiceName, clusterInfo.Name)
				return nil
			}

			if runtime != nil {
				current, err := a.apiVersion.Get(ctx, runtime.VersionID)
				if err != nil {
					return err
				}

				//若api为已上线且无更新状态
				if runtime.IsOnline && !a.isAPIVersionConfChange(latest.APIVersionConfig, current.APIVersionConfig) {
					return nil
				}
			}

			//发布到apinto
			client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
			if err != nil {
				item.Status = false
				item.Result = fmt.Sprintf("连接集群失败, err: %s", err.Error())
				return nil
			}

			//封装router配置
			apiDriverInfo := a.GetAPIDriver(latest.Driver)
			routerConfig := apiDriverInfo.ToApinto(api.UUID, api.Desc, false, latest.Method, latest.RequestPath, latest.RequestPathLabel, latest.ProxyPath, strings.ToLower(latest.ServiceName), latest.Timeout, latest.Retry, latest.EnableWebsocket, latest.Match, latest.Header, latest.TemplateUUID)

			//未上线
			if runtime == nil {
				runtime = &apientry.APIRuntime{
					NamespaceId: namespaceId,
					ApiID:       api.Id,
					ClusterID:   clusterInfo.Id,
					VersionID:   latest.Id,
					IsOnline:    true,
					Disable:     false,
					Operator:    operator,
					CreateTime:  t,
					UpdateTime:  t,
				}

				if err = a.apiRuntime.Insert(txCtx, runtime); err != nil {
					return err
				}
				if err = client.ForRouter().Create(*routerConfig); err != nil {
					item.Status = false
					item.Result = fmt.Sprintf("发送配置至集群失败, err: %s", err.Error())
				}

			} else { //已下线或者待更新
				isOnline := runtime.IsOnline //保存旧状态

				runtime.IsOnline = true
				runtime.UpdateTime = t
				runtime.VersionID = latest.Id
				runtime.Operator = operator

				routerConfig.Disable = runtime.Disable

				if err = a.apiRuntime.Save(txCtx, runtime); err != nil {
					return err
				}

				//若原先是下线状态
				if !isOnline {
					if err = client.ForRouter().Create(*routerConfig); err != nil {
						item.Status = false
						item.Result = fmt.Sprintf("发送配置至集群失败, err: %s", err.Error())
					}
				}

				if err = client.ForRouter().Update(api.UUID+"@router", *routerConfig); err != nil {
					item.Status = false
					item.Result = fmt.Sprintf("发送配置至集群失败, err: %s", err.Error())
				}
			}
			return nil
		})
		if err != nil {
			item.Status = false
			item.Result = err.Error()
		}

		onlineList = append(onlineList, item)
	}

	return onlineList, nil
}

func (a *apiService) BatchOffline(ctx context.Context, namespaceId int, operator int, apiUUIDs, clusterNames []string) ([]*apimodel.BatchListItem, error) {

	errorGroup, _ := errgroup.WithContext(ctx)

	apiList := make([]*apientry.API, 0, len(apiUUIDs))
	errorGroup.Go(func() error {
		//确认所有apiUUID，clusterName均存在
		for _, uid := range apiUUIDs {
			api, err := a.apiStore.GetByUUID(ctx, namespaceId, uid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("apiService doesn't exist. uuid:%s", uid)
				}
				return err
			}
			apiList = append(apiList, api)
		}
		return nil
	})

	clusterList := make([]*cluster_model.Cluster, 0, len(clusterNames))
	errorGroup.Go(func() error {
		clusters, err := a.clusterService.QueryListByNamespaceId(ctx, namespaceId)
		if err != nil {
			return err
		}
		clusterMap := common.SliceToMap(clusters, func(t *cluster_model.Cluster) string {
			return t.Name
		})
		for _, clusterName := range clusterNames {
			if clusterInfo, ok := clusterMap[clusterName]; ok {
				if clusterInfo.Status == 2 || clusterInfo.Status == 3 {
					return fmt.Errorf("cluster status is abnormal. cluster_name:%s", clusterName)
				}
				clusterList = append(clusterList, clusterInfo)
			} else {
				return fmt.Errorf("cluster doesn't exist. cluster_name:%s", clusterName)
			}
		}
		return nil
	})

	if err := errorGroup.Wait(); err != nil {
		return nil, err
	}

	//逐个处理api下线，已经下线或者未上线的不进行操作
	offlineList := make([]*apimodel.BatchListItem, 0, len(apiList)*len(clusterList))
	for _, api := range apiList {
		items, err := a.offline(ctx, operator, api, clusterList)
		if err != nil && len(items) == 0 {
			return nil, err
		}

		offlineList = append(offlineList, items...)
	}
	//编写操作记录
	logApiNameList := make([]string, 0, len(apiList))
	logCLNameList := make([]string, 0, len(clusterList))
	for _, api := range apiList {
		logApiNameList = append(logApiNameList, api.Name)
	}
	for _, cl := range clusterList {
		logCLNameList = append(logCLNameList, cl.Name)
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        strings.Join(logApiNameList, ","),
		ClusterName: strings.Join(logCLNameList, ","),
		PublishType: 2,
	})

	return offlineList, nil
}

func (a *apiService) offline(ctx context.Context, operator int, api *apientry.API, clusterList []*cluster_model.Cluster) ([]*apimodel.BatchListItem, error) {
	offlineList := make([]*apimodel.BatchListItem, 0)

	err := a.lockService.Lock(locker_service.LockNameAPI, api.Id)
	if err != nil {
		return nil, err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, api.Id)

	latestApi, err := a.apiStore.Get(ctx, api.Id)
	if err != nil {
		return nil, err
	}

	for _, clusterInfo := range clusterList {
		//获取当前的版本
		runtime, err := a.apiRuntime.GetForCluster(ctx, api.Id, clusterInfo.Id)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}

		item := &apimodel.BatchListItem{
			APIName:    latestApi.Name,
			ClusterEnv: fmt.Sprintf("%s_%s", clusterInfo.Name, clusterInfo.Env),
			Status:     true,
			Result:     "",
		}
		//上线状态的进行下线操作，未上线或已下线状态直接成功
		if runtime != nil && runtime.IsOnline {
			err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
				runtime.IsOnline = false
				runtime.UpdateTime = time.Now()
				runtime.Operator = operator
				err = a.apiRuntime.Save(txCtx, runtime)
				if err != nil {
					return err
				}

				//发布到apinto
				client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
				if err != nil {
					return err
				}
				return common.CheckWorkerNotExist(client.ForRouter().Delete(api.UUID + "@router"))
			})
			if err != nil {
				item.Status = false
				item.Result = err.Error()
			}
		}

		offlineList = append(offlineList, item)
	}

	return offlineList, nil
}

func (a *apiService) BatchOnlineCheck(ctx context.Context, namespaceId int, operator int, apiUUIDs, clusterNames []string) ([]*apimodel.BatchOnlineCheckListItem, string, error) {
	//确认所有apiUUID，clusterName均存在
	apiList := make([]*apientry.API, 0, len(apiUUIDs))
	apiIds := make([]int, 0, len(apiUUIDs))

	groupInfo, _ := errgroup.WithContext(ctx)
	groupInfo.Go(func() error {
		for _, uid := range apiUUIDs {
			api, err := a.apiStore.GetByUUID(ctx, namespaceId, uid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("api doesn't exist. uuid:%s", uid)
				}
				return err
			}
			apiList = append(apiList, api)
			apiIds = append(apiIds, api.Id)
		}
		return nil

	})

	clusterList := make([]*cluster_model.Cluster, 0, len(clusterNames))
	groupInfo.Go(func() error {
		clusters, err := a.clusterService.QueryListByNamespaceId(ctx, namespaceId)
		if err != nil {
			return err
		}
		clusterMap := common.SliceToMap(clusters, func(t *cluster_model.Cluster) string {
			return t.Name
		})
		for _, clusterName := range clusterNames {
			if clusterInfo, ok := clusterMap[clusterName]; ok {
				if clusterInfo.Status == 2 || clusterInfo.Status == 3 {
					return fmt.Errorf("cluster status is abnormal. cluster_name:%s", clusterName)
				}
				clusterList = append(clusterList, clusterInfo)
			} else {
				return fmt.Errorf("cluster doesn't exist. cluster_name:%s", clusterName)
			}
		}
		return nil
	})

	if err := groupInfo.Wait(); err != nil {
		return nil, "", err
	}

	isAllOnline := true
	checkServiceMap := make(map[int]string)     //serviceId集合，用于对检查列表的去重
	checkTemplateMap := make(map[string]string) //插件模板ID集合，用于对检查列表的去重
	checkList := make([]*apimodel.BatchOnlineCheckListItem, 0, len(apiList)*len(clusterList))

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, "", err
	}

	//确认每个api对应的cluster所配置的serviceID和模板的上线情况
	for _, api := range apiList {
		version := versionMap[api.Id]
		if _, has := checkServiceMap[version.ServiceID]; !has {
			checkServiceMap[version.ServiceID] = version.ServiceName
		}
		if _, has := checkTemplateMap[version.TemplateUUID]; !has && version.TemplateUUID != "" {
			templateInfo, err := a.pluginTemplateService.GetBasicInfoByUUID(ctx, version.TemplateUUID)
			if err != nil {
				return nil, "", err
			}
			checkTemplateMap[version.TemplateUUID] = templateInfo.Name
		}
	}

	for serviceID, serName := range checkServiceMap {
		for _, clusterInfo := range clusterList {
			item := &apimodel.BatchOnlineCheckListItem{
				ServiceTemplate: serName,
				ClusterEnv:      fmt.Sprintf("%s%s", clusterInfo.Name, clusterInfo.Env),
				Status:          true,
				Solution:        &frontend_model.Router{},
			}

			if isOnline := a.service.IsOnline(ctx, clusterInfo.Id, serviceID); !isOnline {
				isAllOnline = false
				item.Status = false
				item.Result = fmt.Sprintf("%s未上线到%s", serName, clusterInfo.Name)
				item.Solution.Name = frontend_model.RouterNameServiceOnline
				item.Solution.Params = map[string]string{"cluster_name": clusterInfo.Name, "service_name": serName}
			}
			checkList = append(checkList, item)
		}
	}
	for templateUuid, templateName := range checkTemplateMap {
		for _, clusterInfo := range clusterList {
			item := &apimodel.BatchOnlineCheckListItem{
				ServiceTemplate: templateName,
				ClusterEnv:      fmt.Sprintf("%s%s", clusterInfo.Name, clusterInfo.Env),
				Status:          true,
				Solution:        &frontend_model.Router{},
			}
			isOnline, err := a.pluginTemplateService.IsOnline(ctx, clusterInfo.Id, templateUuid)
			if err != nil {
				return nil, "", err
			}
			if !isOnline {
				isAllOnline = false
				item.Status = false
				item.Result = fmt.Sprintf("%s未上线到%s", templateName, clusterInfo.Name)
				item.Solution.Name = frontend_model.RouterNameTemplateOnline
				item.Solution.Params = map[string]string{"cluster_name": clusterInfo.Name, "template_uuid": templateUuid}
			}
			checkList = append(checkList, item)
		}
	}

	//若所有的API均已上线，则生成一个UUID
	onlineToken := ""
	if isAllOnline {
		onlineToken = uuid.New()

		taskData := &apimodel.APIBatchConf{
			ApiUUIDs:     apiUUIDs,
			ClusterNames: clusterNames,
		}

		data, _ := json.Marshal(taskData)
		task := &apimodel.BatchOnlineCheckTask{
			Operator: operator,
			Token:    onlineToken,
			Data:     data,
		}
		key := a.batchApiCache.Key(onlineToken)
		if err := a.batchApiCache.Set(ctx, key, task, time.Hour*8); err != nil {
			return nil, "", err
		}
	}

	return checkList, onlineToken, nil
}

func (a *apiService) OnlineList(ctx context.Context, namespaceId int, uuid string) ([]*apimodel.APIOnlineListItem, error) {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	//获取工作空间下的所有集群
	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *cluster_model.Cluster) int {
		return t.Id
	})

	//获取当前服务发现下集群运行的版本
	runtimes, err := a.apiRuntime.GetByTarget(ctx, apiInfo.Id)
	if err != nil {
		return nil, err
	}
	runtimeMaps := common.SliceToMap(runtimes, func(t *apientry.APIRuntime) int {
		return t.ClusterID
	})

	//获取操作人用户列表
	operatorList := common.SliceToSliceIds(runtimes, func(t *apientry.APIRuntime) int {
		return t.Operator
	})

	userInfoMaps, err := a.userInfoService.GetUserInfoMaps(ctx, operatorList...)
	if err != nil {
		return nil, err
	}

	list := make([]*apimodel.APIOnlineListItem, 0, len(clusters))

	latestVersion, err := a.GetLatestAPIVersion(ctx, apiInfo.Id)
	if err != nil {
		return nil, err
	}

	for _, clusterInfo := range clusterMaps {
		apiOnline := &apimodel.APIOnlineListItem{
			ClusterName: clusterInfo.Name,
			ClusterEnv:  clusterInfo.Env,
			Status:      1, //默认为未上线状态
		}
		if runtime, ok := runtimeMaps[clusterInfo.Id]; ok {

			operator := ""
			if userInfo, uOk := userInfoMaps[runtime.Operator]; uOk {
				operator = userInfo.NickName
			}

			apiOnline.Operator = operator
			apiOnline.Disable = runtime.Disable
			apiOnline.UpdateTime = runtime.UpdateTime
			if runtime.IsOnline {
				apiOnline.Status = 3 //已上线
			} else {
				apiOnline.Status = 2 //已下线
			}
			//已上线需要对比是否更新过 服务发现信息
			if apiOnline.Status == 3 && runtime.VersionID != latestVersion.Id {
				apiOnline.Status = 4 //待更新
			}
		}

		list = append(list, apiOnline)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Status > list[j].Status
	})
	return list, nil
}

func (a *apiService) OnlineAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) (*frontend_model.Router, error) {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	err = a.lockService.Lock(locker_service.LockNameAPI, apiInfo.Id)
	if err != nil {
		return nil, err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	apiID := apiInfo.Id
	t := time.Now()

	//获取当前集群信息
	clusterInfo, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	latestVersion, err := a.GetLatestAPIVersion(ctx, apiID)
	if err != nil {
		return nil, err
	}

	//判断上游服务有没有上线
	if !a.service.IsOnline(ctx, clusterInfo.Id, latestVersion.ServiceID) {
		return &frontend_model.Router{
			Name:   frontend_model.RouterNameServiceOnline,
			Params: map[string]string{"service_name": latestVersion.ServiceName},
		}, errors.New(fmt.Sprintf("绑定的%s未上线到%s", latestVersion.ServiceName, clusterName))
	}

	//判断插件模板有没有上线
	if latestVersion.TemplateID != 0 {
		isTemplateOnline, err := a.pluginTemplateService.IsOnline(ctx, clusterInfo.Id, latestVersion.TemplateUUID)
		if err != nil {
			return nil, err
		}
		if !isTemplateOnline {
			return &frontend_model.Router{
				Name:   frontend_model.RouterNameTemplateOnline,
				Params: map[string]string{"template_uuid": latestVersion.TemplateUUID},
			}, errors.New(fmt.Sprintf("绑定的插件模板未上线到%s", clusterName))
		}
	}

	//获取当前运行的版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiID, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//发布到apinto
	client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
	if err != nil {
		return nil, err
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        uuid,
		Name:        apiInfo.Name,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
		PublishType: 1,
	})

	//事务
	err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {

		apiDriverInfo := a.GetAPIDriver(latestVersion.Driver)
		routerConfig := apiDriverInfo.ToApinto(apiInfo.UUID, apiInfo.Desc, false, latestVersion.Method, latestVersion.RequestPath, latestVersion.RequestPathLabel, latestVersion.ProxyPath, strings.ToLower(latestVersion.ServiceName), latestVersion.Timeout, latestVersion.Retry, latestVersion.EnableWebsocket, latestVersion.Match, latestVersion.Header, latestVersion.TemplateUUID)
		if runtime == nil {
			runtime = &apientry.APIRuntime{
				NamespaceId: namespaceId,
				ApiID:       apiID,
				ClusterID:   clusterInfo.Id,
				VersionID:   latestVersion.Id,
				IsOnline:    true,
				Disable:     false,
				Operator:    operator,
				CreateTime:  t,
				UpdateTime:  t,
			}

			if err = a.apiRuntime.Insert(txCtx, runtime); err != nil {
				return err
			}
			return client.ForRouter().Create(*routerConfig)
		} else {
			//保存旧状态
			isOnline := runtime.IsOnline

			runtime.IsOnline = true
			runtime.UpdateTime = t
			runtime.VersionID = latestVersion.Id
			runtime.Operator = operator

			routerConfig.Disable = runtime.Disable

			if err = a.apiRuntime.Save(txCtx, runtime); err != nil {
				return err
			}

			//若原先是下线状态
			if !isOnline {
				return client.ForRouter().Create(*routerConfig)
			}

			return client.ForRouter().Update(apiInfo.UUID+"@router", *routerConfig)
		}
	})

	return nil, err
}

func (a *apiService) ResetOnline(ctx context.Context, _, clusterId int) {
	runtimes, err := a.apiRuntime.GetByCluster(ctx, clusterId)
	if err != nil {
		log.Errorf("apiService-ResetOnline-getRuntimes clusterId=%d,err=%d", clusterId, err.Error())
		return
	}
	client, err := a.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		log.Errorf("apiService-ResetOnline-getClient clusterId=%d,err=%d", clusterId, err.Error())
		return
	}

	for _, runtime := range runtimes {
		if !runtime.IsOnline {
			continue
		}

		apiInfo, err := a.apiStore.Get(ctx, runtime.ApiID)
		if err != nil {
			log.Errorf("apiService-ResetOnline-getApiInfo apiId=%d, clusterId=%d,err=%d", runtime.ApiID, clusterId, err.Error())
			continue
		}

		version, err := a.apiVersion.Get(ctx, runtime.VersionID)
		if err != nil {
			log.Errorf("apiService-ResetOnline-getVersion versionId=%d, clusterId=%d,err=%d", runtime.VersionID, clusterId, err.Error())
			continue
		}
		routerConfig := a.GetAPIDriver(version.Driver).ToApinto(apiInfo.UUID, apiInfo.Desc, false, version.Method, version.RequestPath, version.RequestPathLabel, version.ProxyPath, strings.ToLower(version.ServiceName), version.Timeout, version.Retry, version.EnableWebsocket, version.Match, version.Header, version.TemplateUUID)

		if err = client.ForRouter().Create(*routerConfig); err != nil {
			log.Errorf("apiService-ResetOnline-apintoCreate routerConfig=%d, clusterId=%d,err=%d", routerConfig, clusterId, err.Error())
			continue
		}
	}
}

func (a *apiService) OfflineAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	//获取当前集群信息
	clusterInfo, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	//获取当前的版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiInfo.Id, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}

	t := time.Now()

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        uuid,
		Name:        apiInfo.Name,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
		PublishType: 2,
	})

	//事务
	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		if !runtime.IsOnline {
			return errors.New("已下线不可重复下线")
		}
		runtime.IsOnline = false
		runtime.UpdateTime = t
		runtime.Operator = operator
		err = a.apiRuntime.Save(txCtx, runtime)
		if err != nil {
			return err
		}

		//发布到apinto
		client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
		if err != nil {
			return err
		}

		return common.CheckWorkerNotExist(client.ForRouter().Delete(apiInfo.UUID + "@router"))
	})
}

func (a *apiService) EnableAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	//获取当前集群信息
	clusterInfo, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, apiInfo.Id)
	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}
	//获取当前版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiInfo.Id, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}
	if !runtime.IsOnline {
		return errors.New("Api must be online. ")
	}

	t := time.Now()
	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:          uuid,
		Name:          apiInfo.Name,
		ClusterId:     clusterInfo.Id,
		ClusterName:   clusterName,
		EnableOperate: 1,
	})

	//事务
	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		runtime.UpdateTime = t
		runtime.Operator = operator
		runtime.Disable = false
		err = a.apiRuntime.Save(txCtx, runtime)
		if err != nil {
			return err
		}

		//发布到apinto
		client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
		if err != nil {
			return err
		}
		return client.ForRouter().Patch(apiInfo.UUID+"@router", map[string]interface{}{"disable": false})
	})
}

func (a *apiService) DisableAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	//获取当前集群信息
	clusterInfo, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	err = a.lockService.Lock(locker_service.LockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.Unlock(locker_service.LockNameAPI, apiInfo.Id)
	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}
	//获取当前版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiInfo.Id, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}
	if !runtime.IsOnline {
		return errors.New("Api must be online. ")
	}

	t := time.Now()

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:          uuid,
		Name:          apiInfo.Name,
		ClusterId:     clusterInfo.Id,
		ClusterName:   clusterName,
		EnableOperate: 2,
	})

	//事务
	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		runtime.UpdateTime = t
		runtime.Operator = operator
		runtime.Disable = true
		err = a.apiRuntime.Save(txCtx, runtime)
		if err != nil {
			return err
		}

		//发布到apinto
		client, err := a.apintoClient.GetClient(ctx, clusterInfo.Id)
		if err != nil {
			return err
		}
		return client.ForRouter().Patch(apiInfo.UUID+"@router", map[string]interface{}{"disable": true})
	})
}

func (a *apiService) GetSource(ctx context.Context) ([]*apimodel.SourceListItem, error) {
	items := make([]*apimodel.SourceListItem, 0, 2)

	sourceList, err := a.apiStore.GetSourceList(ctx)
	if err != nil {
		return items, err
	}

	//对来源列表进行排序
	sort.Sort(apientry.APISourceList(sourceList))

	for _, source := range sourceList {
		title := ""
		if source.SourceType == enum.SourceSync {
			//若来源为同步，通过来源id查找外部应用名
			appName, _ := a.extAppService.GetExtAppName(ctx, source.SourceID)
			title = appName
			if source.SourceLabel != "" {
				title = appName + "-" + source.SourceLabel
			}
		} else {
			title = enum.GetSourceTitle(source.SourceType)
		}
		items = append(items, &apimodel.SourceListItem{
			Id:    fmt.Sprintf("%s:%d:%s", source.SourceType, source.SourceID, source.SourceLabel),
			Title: title,
		})
	}

	return items, nil
}

func (a *apiService) GetImportCheckList(ctx context.Context, namespaceId int, fileData []byte, groupID, serviceName, requestPrefix string) ([]*apimodel.ImportAPIListItem, string, error) {
	//解析swagger3.0 TODO 写死解析3.0 等之后有其他格式再用driverManager，openAPI同步现在是用driverManager的
	swaggerConfig := new(open_app_model.SwaggerConfig)
	reader := bytes.NewReader(fileData)

	if err := common.DecodeYAML(reader, swaggerConfig); err != nil {
		reader = bytes.NewReader(fileData)
		err = common.DecodeJSON(reader, swaggerConfig)
		if err != nil {
			return nil, "", err
		}
	}

	//参数校验
	isExist, err := a.commonGroup.IsGroupExist(ctx, groupID)
	if err != nil {
		return nil, "", err
	}

	if !isExist {
		return nil, "", errors.New("分组不存在")
	}

	if _, err = a.service.GetServiceIDByName(ctx, namespaceId, serviceName); err != nil {
		return nil, "", errors.New("上游服务不存在")
	}

	//格式化requestPrefix
	if requestPrefix != "" {
		requestPrefix = "/" + strings.Trim(requestPrefix, "/")
	}

	if requestPrefix, err = common.CheckAndFormatPath(requestPrefix); err != nil {
		return nil, "", errors.New("requet_prefix is illegal. ")
	}

	items := make([]*apimodel.ImportAPIListItem, 0)
	for path, pathMap := range swaggerConfig.Paths {
		for method, info := range pathMap {
			name := info.Summary
			if name == "" {
				name = info.OperationID
			}
			if name == "" {
				name = method + "-" + requestPrefix + path
			}
			items = append(items, &apimodel.ImportAPIListItem{
				Name:   name,
				Method: strings.ToUpper(method),
				Path:   requestPrefix + path,
				Desc:   info.Description,
				Status: 1,
			})
		}
	}

	//获取所有API
	apiList, err := a.GetAPIListByName(ctx, namespaceId, "")
	if err != nil {
		return nil, "", err
	}

	apiMap := common.SliceToMapArray(apiList, func(t *group_model.CommonGroupApi) string {
		return t.Path
	})

	for _, item := range items {
		//暂不支持OPTIONS和TRACE
		if item.Method == "OPTIONS" || item.Method == "TRACE" {
			item.Status = 2
			continue
		}

		if item.Path, err = common.CheckAndFormatPath(item.Path); err != nil {
			item.Status = 3
			continue
		}

		if apis, ok := apiMap[common.ReplaceRestfulPath(item.Path, enum.RestfulLabel)]; ok {
			for _, groupApi := range apis {
				if slices.Contains(groupApi.Methods, item.Method) {
					item.Status = 2
					break
				}
			}
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Status > items[j].Status
	})

	redisDataItems := make([]*apimodel.ImportAPIRedisDataItem, 0)
	for i, item := range items {
		item.Id = i + 1
		//如果状态不为正常，则不存进redis
		if item.Status != 1 {
			continue
		}

		e := &apientry.API{
			NamespaceId:      namespaceId,
			UUID:             uuid.New(),
			GroupUUID:        groupID,
			Name:             item.Name,
			RequestPath:      common.ReplaceRestfulPath(item.Path, enum.RestfulLabel),
			RequestPathLabel: item.Path,
			SourceType:       enum.SourceImport,
			SourceID:         -1,
			SourceLabel:      "",
			Desc:             item.Desc,
			CreateTime:       time.Time{},
			UpdateTime:       time.Time{},
		}

		apiInfo := &apimodel.APIInfo{
			API:    e,
			Method: strings.Split(item.Method, ","),
		}

		redisDataItem := &apimodel.ImportAPIRedisDataItem{
			ID:  item.Id,
			Api: apiInfo,
		}
		redisDataItems = append(redisDataItems, redisDataItem)
	}

	token := uuid.New()
	//数据存储到缓存
	key := a.importApiCache.Key(token)

	importAPIRedisData := &apimodel.ImportAPIRedisData{
		Apis:        redisDataItems,
		ServiceName: serviceName,
		GroupID:     groupID,
	}

	if err = a.importApiCache.Set(ctx, key, importAPIRedisData, time.Hour*8); err != nil {
		return nil, "", err
	}

	return items, token, nil
}

func (a *apiService) ImportAPI(ctx context.Context, namespaceId, operator int, input *api_dto.ImportAPIInfos) error {

	key := a.importApiCache.Key(input.Token)
	apiData, err := a.importApiCache.Get(ctx, key)
	if err != nil {
		return err
	}
	//判断目录是否存在
	isExist, err := a.commonGroup.IsGroupExist(ctx, apiData.GroupID)
	if err != nil {
		return err
	}
	if !isExist {
		return errors.New("分组不存在,请重新导入")
	}

	//判断服务是否存在
	serviceID, err := a.service.GetServiceIDByName(ctx, namespaceId, apiData.ServiceName)
	if err != nil {
		return err
	}

	importApiMaps := common.SliceToMap(apiData.Apis, func(t *apimodel.ImportAPIRedisDataItem) int {
		return t.ID
	})

	//获取现存所有API
	existedApiList, err := a.GetAPIListByName(ctx, namespaceId, "")
	if err != nil {
		return err
	}

	existedApiMaps := common.SliceToMapArray(existedApiList, func(t *group_model.CommonGroupApi) string {
		return t.Path
	})

	createApis := make([]*apimodel.APIInfo, 0, len(input.Apis))
	logAPINames := make([]string, 0, len(input.Apis))
	for _, api := range input.Apis {
		if v, ok := importApiMaps[api.Id]; ok {
			if api.Name != "" {
				v.Api.Name = api.Name
			}
			// TODO 现在只能修改apiName， 请求路径和描述以后可能要改

			//检查api是否有冲突
			isReduplicated := false
			if rApis, ok := existedApiMaps[v.Api.RequestPath]; ok {
			A:
				for _, rApi := range rApis {
					for _, method := range v.Api.Method {
						if slices.Contains(rApi.Methods, method) {
							isReduplicated = true
							break A
						}
					}

				}
			}
			if isReduplicated {
				log.Errorf("import api %s fail. api is reduplicated. path:%s. method:%s. ", v.Api.Name, v.Api.RequestPathLabel, v.Api.Method)
			} else {
				createApis = append(createApis, v.Api)
				logAPINames = append(logAPINames, v.Api.Name)
			}
		} else {
			return errors.New(fmt.Sprintf("序号为%d的数据不存在", api.Id))
		}
	}

	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: strings.Join(logAPINames, ","),
	})

	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()

		for _, apiInfo := range createApis {
			apiInfo.CreateTime = t
			apiInfo.UpdateTime = t
			apiInfo.Operator = operator
			if err = a.apiStore.Save(txCtx, apiInfo.API); err != nil {
				return err
			}
			//添加版本信息
			apiVersionInfo := &apientry.APIVersion{
				ApiID:       apiInfo.Id,
				NamespaceID: namespaceId,
				APIVersionConfig: apientry.APIVersionConfig{
					Driver:           "http",
					RequestPath:      apiInfo.RequestPath,
					RequestPathLabel: apiInfo.RequestPathLabel,
					ServiceID:        serviceID,
					ServiceName:      apiData.ServiceName,
					Method:           apiInfo.Method,
					ProxyPath:        apiInfo.RequestPathLabel,
					Timeout:          10000,
					Retry:            0,
					EnableWebsocket:  false,
					Match:            []*apientry.MatchConf{},
					Header:           []*apientry.ProxyHeader{},
				},
				Operator:   operator,
				CreateTime: t,
			}

			if err = a.apiVersion.Save(txCtx, apiVersionInfo); err != nil {
				return err
			}

			if err = a.apiHistory.HistoryAdd(txCtx, namespaceId, apiInfo.Id, &apientry.ApiHistoryInfo{
				Api:    *apiInfo.API,
				Config: apiVersionInfo.APIVersionConfig,
			}, operator); err != nil {
				return err
			}

			stat := &apientry.APIStat{
				APIID:     apiInfo.Id,
				VersionID: apiVersionInfo.Id,
			}

			//添加版本关联原表信息
			if err = a.apiStat.Save(txCtx, stat); err != nil {
				return err
			}

			//quote更新所引用的服务
			if err = a.quoteStore.Set(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quote_entry.QuoteTargetKindTypeService, serviceID); err != nil {
				return err
			}

		}
		return nil
	})
}

func (a *apiService) GetAPIListByName(ctx context.Context, namespaceId int, name string) ([]*group_model.CommonGroupApi, error) {
	apiList, err := a.apiStore.GetListByName(ctx, namespaceId, name)
	if err != nil {
		return nil, err
	}
	groupAPIs := make([]*group_model.CommonGroupApi, 0, len(apiList))

	apiIds := make([]int, 0, len(apiList))
	for _, api := range apiList {
		apiIds = append(apiIds, api.Id)
	}

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, err
	}

	for _, api := range apiList {
		version := versionMap[api.Id]
		groupApi := &group_model.CommonGroupApi{
			Path:      api.RequestPath,
			PathLabel: api.RequestPathLabel,
			Name:      api.Name,
			UUID:      api.UUID,
			Methods:   version.Method,
			GroupUUID: api.GroupUUID,
		}
		groupAPIs = append(groupAPIs, groupApi)
	}

	return groupAPIs, nil
}

func (a *apiService) GetAPIListByServiceName(ctx context.Context, namespaceId int, serviceNames []string) ([]*apimodel.APIInfo, error) {

	var err error
	groupAPIs := make([]*apimodel.APIInfo, 0)

	for _, serviceName := range serviceNames {

		target := 0
		if serviceName != "" {
			serviceId, err := a.service.GetServiceIDByName(ctx, namespaceId, serviceName)
			if err != nil {
				return nil, err
			}
			target = serviceId
		}

		apiList := make([]*apientry.API, 0)

		if target > 0 {
			quote, err := a.quoteStore.GetTargetQuote(ctx, target, quote_entry.QuoteTargetKindTypeService)
			if err != nil {
				return nil, err
			}
			apiList, err = a.apiStore.GetByIds(ctx, namespaceId, quote[quote_entry.QuoteKindTypeAPI])
			if err != nil {
				return nil, err
			}
		} else {
			apiList, err = a.apiStore.GetListAll(ctx, namespaceId)
			if err != nil {
				return nil, err
			}
		}

		for _, api := range apiList {
			groupApi := &apimodel.APIInfo{
				API: api,
			}
			groupAPIs = append(groupAPIs, groupApi)
		}

	}

	return groupAPIs, nil
}

func (a *apiService) isApiCanDelete(ctx context.Context, apiId int) (bool, error) {
	count, err := a.apiRuntime.OnlineCount(ctx, apiId)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil
	}

	return true, nil
}

func (a *apiService) GetLatestAPIVersion(ctx context.Context, apiId int) (*apientry.APIVersion, error) {
	stat, err := a.apiStat.Get(ctx, apiId)
	if err != nil {
		return nil, err
	}
	return a.apiVersion.Get(ctx, stat.VersionID)
}

// CheckAPIReduplicative 检测API配置是否重复，不可同名同request_url同method
func (a *apiService) CheckAPIReduplicative(ctx context.Context, namespaceID int, uuid string, input *api_dto.APIInfo) error {
	//获取相同requestPath的API
	apiList, err := a.apiStore.GetListByRequestPath(ctx, namespaceID, input.RequestPath)
	if err != nil {
		return err
	}
	inputLen := len(input.Method)
	for _, api := range apiList {
		//筛去api本身
		if api.UUID == uuid {
			continue
		}

		apiVersion, err := a.GetLatestAPIVersion(ctx, api.Id)
		if err != nil {
			return err
		}
		//查重Method  空Method数组表示ALL，ALL和其它method不重复，但ALL和ALL会重复

		//若已有API的method为ALL
		if len(apiVersion.Method) == 0 && inputLen == 0 {
			return fmt.Errorf("requestPath %s and method ALL is reduplicative. ", input.RequestPathLabel)
		} else {
			//若已有API的method不为为ALL
			if inputLen == 0 {
				continue
			}
			currentMap := common.SliceToMap(input.Method, func(method string) string {
				return method
			})
			for _, m := range input.Method {
				if _, has := currentMap[m]; has {
					return fmt.Errorf("requestPath %s and method %s is reduplicative. ", input.RequestPathLabel, m)
				}
			}
		}

	}

	return nil
}

func (a *apiService) IsAPIOnline(ctx context.Context, clusterId, apiID int) bool {
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiID, clusterId)
	if err != nil {
		return false
	}
	return runtime.IsOnline
}

func (a *apiService) isAPIVersionConfChange(latest apientry.APIVersionConfig, current apientry.APIVersionConfig) bool {
	return !reflect.DeepEqual(latest, current)
}

func (a *apiService) GetAPIDriver(driverName string) apiservice.IAPIDriver {
	return a.apiManager.GetDriver(driverName)
}

func (a *apiService) GetAPINameByID(ctx context.Context, apiID int) (string, error) {
	info, err := a.apiStore.Get(ctx, apiID)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}

func (a *apiService) GetAPIRemoteOptions(ctx context.Context, namespaceID, pageNum, pageSize int, keyword, groupUuid string) ([]*strategy_model.RemoteApis, int, error) {
	groupList := make([]string, 0)
	var err error
	//获取传入的groupUUID下包括子分组的所有UUID
	if groupUuid != "" {
		groupList, err = a.commonGroup.GroupUUIDS(ctx, namespaceID, group_service.ApiName, group_service.ApiName, groupUuid)
		if err != nil {
			return nil, 0, err
		}
		groupList = append(groupList, groupUuid)
	}

	groups, err := a.commonGroup.GroupListAll(ctx, namespaceID, group_service.ApiName, group_service.ApiName)
	if err != nil {
		return nil, 0, err
	}

	//获取指定目录下所有API
	apis, total, err := a.apiStore.GetListPageByGroupIDs(ctx, namespaceID, pageNum, pageSize, groupList, nil, keyword)
	if err != nil {
		return nil, 0, err
	}
	apiList := make([]*strategy_model.RemoteApis, 0, len(apis))

	groupUUIDMap := common.SliceToMap(groups, func(t *group_entry.CommonGroup) string {
		return t.Uuid
	})
	groupIdMap := common.SliceToMap(groups, func(t *group_entry.CommonGroup) int {
		return t.Id
	})

	for _, api := range apis {
		//version, err := a.GetLatestAPIVersion(ctx, apiService.Id)
		//if err != nil {
		//	return nil, 0, err
		//}
		parentGroupName := &[]string{}

		a.commonGroup.ParentGroupName(api.GroupUUID, groupUUIDMap, groupIdMap, parentGroupName)

		item := &strategy_model.RemoteApis{
			Uuid: api.UUID,
			Name: api.Name,
			//Service:     version.ServiceName,
			Group:       strings.Join(*parentGroupName, "/"), //TODO
			RequestPath: api.RequestPathLabel,
		}

		apiList = append(apiList, item)
	}

	return apiList, total, nil
}

func (a *apiService) GetAPIRemoteByUUIDS(ctx context.Context, namespace int, uuids []string) ([]*strategy_model.RemoteApis, error) {

	groups, err := a.commonGroup.GroupListAll(ctx, namespace, group_service.ApiName, group_service.ApiName)
	if err != nil {
		return nil, err
	}

	groupUUIDMap := common.SliceToMap(groups, func(t *group_entry.CommonGroup) string {
		return t.Uuid
	})
	groupIdMap := common.SliceToMap(groups, func(t *group_entry.CommonGroup) int {
		return t.Id
	})

	apis, err := a.apiStore.GetByUUIDs(ctx, namespace, uuids)
	if err != nil {
		return nil, err
	}

	apiIds := make([]int, 0, len(apis))
	for _, api := range apis {
		apiIds = append(apiIds, api.Id)
	}

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, err
	}

	apiList := make([]*strategy_model.RemoteApis, 0, len(apis))
	for _, api := range apis {
		version := versionMap[api.Id]

		parentGroupName := &[]string{}
		a.commonGroup.ParentGroupName(api.GroupUUID, groupUUIDMap, groupIdMap, parentGroupName)

		item := &strategy_model.RemoteApis{
			Uuid:        api.UUID,
			Name:        api.Name,
			Service:     version.ServiceName,
			Group:       strings.Join(*parentGroupName, "/"), //TODO
			RequestPath: version.RequestPathLabel,
		}

		apiList = append(apiList, item)
	}

	return apiList, nil
}
