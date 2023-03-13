package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/common"
	driver_manager "github.com/eolinker/apinto-dashboard/driver-manager"
	"github.com/eolinker/apinto-dashboard/driver-manager/driver"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"reflect"
	"sort"
	"strings"
	"time"
)

type IAPIService interface {
	GetAPIList(ctx context.Context, namespaceID int, groupUUID, searchName string, searchSources []string, pageNum, pageSize int) ([]*model.APIListItem, int, error)
	GetAPICountByGroupUUID(ctx context.Context, namespaceID int, groupUUID string) int64
	GetAPIVersionInfo(ctx context.Context, namespaceID int, uuid string) (*model.APIVersionInfo, error)
	GetAPIInfo(ctx context.Context, namespaceID int, uuid string) (*model.APIInfo, error)
	GetAPIInfoByGroupUUID(ctx context.Context, namespaceID int, groupUUID string) ([]*model.APIInfo, error)
	GetAPIInfoByUUIDS(ctx context.Context, namespaceID int, uuids []string) ([]*model.APIInfo, error)
	GetAPIInfoByPath(ctx context.Context, namespaceID int, path string) ([]*model.APIInfo, error)
	GetAPIListItemByUUIDS(ctx context.Context, namespaceID int, uuids []string) ([]*model.APIListItem, error)

	GetAPIInfoAll(ctx context.Context, namespaceID int) ([]*model.APIInfo, error)

	GetAPIListItemAll(ctx context.Context, namespaceID int) ([]*model.APIListItem, error)
	GetAPIsForSync(ctx context.Context, namespaceID int) ([]*model.APIVersionInfo, error)

	CreateAPI(ctx context.Context, namespaceID int, operator int, input *dto.APIInfo) error
	UpdateAPI(ctx context.Context, namespaceID int, operator int, input *dto.APIInfo) error
	DeleteAPI(ctx context.Context, namespaceId, operator int, uuid string) error
	GetGroups(ctx context.Context, namespaceId int, parentUuid, queryName string) (*model.CommonGroupRoot, []*model.CommonGroupApi, error)

	BatchOnline(ctx context.Context, namespaceId int, operator int, onlineToken string) ([]*model.BatchListItem, error)
	BatchOffline(ctx context.Context, namespaceId int, operator int, apiUUIDs, clusterNames []string) ([]*model.BatchListItem, error)
	BatchOnlineCheck(ctx context.Context, namespaceId int, operator int, apiUUIDs, clusterNames []string) ([]*model.BatchOnlineCheckListItem, string, error)

	OnlineList(ctx context.Context, namespaceId int, uuid string) ([]*model.APIOnlineListItem, error)
	OnlineAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) (*model.Router, error)
	OfflineAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error
	EnableAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error
	DisableAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error

	GetSource(ctx context.Context) ([]*model.SourceListItem, error)
	GetImportCheckList(ctx context.Context, namespaceId int, fileData []byte, groupID, serviceName, requestPrefix string) ([]*model.ImportAPIListItem, string, error)
	ImportAPI(ctx context.Context, namespaceId, operator int, input *dto.ImportAPIInfos) error

	GetAPIListByName(ctx context.Context, namespaceId int, name string) ([]*model.CommonGroupApi, error)
	GetAPIListByServiceName(ctx context.Context, namespaceId int, serviceName []string) ([]*model.APIInfo, error)
	GetLatestAPIVersion(ctx context.Context, apiId int) (*entry.APIVersion, error)
	IsAPIOnline(ctx context.Context, clusterId, apiID int) bool
	GetAPIDriver(driverName string) driver.IAPIDriver
	GetAPINameByID(ctx context.Context, apiID int) (string, error)
	GetAPIRemoteOptions(ctx context.Context, namespaceId, pageNum, pageSize int, keyword, groupUuid string) ([]*model.RemoteApis, int, error)
	GetAPIRemoteByUUIDS(ctx context.Context, namespace int, uuids []string) ([]*model.RemoteApis, error)
	ResetOnline(ctx context.Context, namespaceId, clusterId int)
}

type apiService struct {
	apiStore   store.IAPIStore
	apiStat    store.IAPIStatStore
	apiVersion store.IAPIVersionStore
	apiRuntime store.IAPIRuntimeStore
	quoteStore store.IQuoteStore
	apiHistory store.IApiHistoryStore

	service          IService
	commonGroup      ICommonGroupService
	clusterService   IClusterService
	namespaceService INamespaceService
	apintoClient     IApintoClient
	userInfoService  IUserInfoService
	extAppService    IExternalApplicationService
	apiManager       driver_manager.IAPIDriverManager

	lockService    IAsynLockService
	importApiCache cache.IImportApiCache
	batchApiCache  cache.IBatchOnlineApiTaskCache
}

func newAPIService() IAPIService {
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

func (a *apiService) GetGroups(ctx context.Context, namespaceId int, parentUuid, queryName string) (*model.CommonGroupRoot, []*model.CommonGroupApi, error) {
	groups, err := a.commonGroup.GroupListAll(ctx, namespaceId, apiName, apiName)
	if err != nil {
		log.Errorf("GetGroups-commonGroup.GroupListAll namespaceId:%d,parentUuid:%s,queryName=%s, err=%s", namespaceId, parentUuid, queryName, err.Error())
		return nil, nil, err
	}

	apis := make([]*model.CommonGroupApi, 0)
	apisAll := make([]*model.CommonGroupApi, 0)
	//查询API
	//apisAll, err := a.GetAPIListByName(ctx, namespaceId, "")
	//if err != nil {
	//	return nil, nil, err
	//}
	//for _, api := range apisAll {
	//	if queryName != "" && strings.Count(strings.ToUpper(api.Name), strings.ToUpper(queryName)) > 0 {
	//		apis = append(apis, api)
	//	}
	//}

	if queryName == "" {
		return a.commonGroup.toGroupRoot(ctx, namespaceId, parentUuid, groups, map[string]string{}), nil, nil
	}

	//查询API的目录直至跟目录
	groupUUIDS := common.SliceToSliceIds(apis, func(t *model.CommonGroupApi) string {
		return t.GroupUUID
	})

	uuidMaps := common.SliceToMap(groupUUIDS, func(t string) string {
		return t
	})
	groupsMaps := common.SliceToMap(groups, func(t *entry.CommonGroup) string {
		return t.Uuid
	})
	groupsIdMaps := common.SliceToMap(groups, func(t *entry.CommonGroup) int {
		return t.Id
	})

	groupsParentIdMaps := common.SliceToMapArray(groups, func(t *entry.CommonGroup) int {
		return t.ParentId
	})

	for _, group := range groups {
		//模糊搜索
		if strings.Count(strings.ToUpper(group.Name), strings.ToUpper(queryName)) > 0 {

			uuidMaps[group.Uuid] = group.Uuid
			//如果绝对相等，需要把改目录下的所有根目录也查询出来
			if group.Name == queryName {
				uuids := &[]string{}
				a.commonGroup.subGroupUUIDS(groupsParentIdMaps, &model.CommonGroup{
					Group: group,
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
		a.commonGroup.parentGroupV2(groupUuid, groupsMaps, groupsIdMaps, outMapUUID)
	}

	groupRoot := a.commonGroup.toGroupRoot(ctx, namespaceId, parentUuid, groups, outMapUUID)

	apiAllMaps := common.SliceToMapArray(apisAll, func(t *model.CommonGroupApi) string {
		return t.GroupUUID
	})

	apiMaps := common.SliceToMap(apis, func(t *model.CommonGroupApi) string {
		return t.GroupUUID
	})

	resApis := &[]*model.CommonGroupApi{}
	a.subGroup(groupRoot.CommonGroup, apiAllMaps, apiMaps, resApis)

	*resApis = append(*resApis, apis...)

	return groupRoot, *resApis, err
}

func (a *apiService) subGroup(list []*model.CommonGroup, apiAllMaps map[string][]*model.CommonGroupApi, apiMaps map[string]*model.CommonGroupApi, apis *[]*model.CommonGroupApi) {
	if len(list) == 0 {
		return
	}
	for _, group := range list {
		if _, ok := apiMaps[group.Group.Uuid]; !ok {
			*apis = append(*apis, apiAllMaps[group.Group.Uuid]...)
		}
		a.subGroup(group.Subgroup, apiAllMaps, apiMaps, apis)
	}
}

func (a *apiService) GetAPIList(ctx context.Context, namespaceID int, groupUUID, searchName string, searchSources []string, pageNum, pageSize int) ([]*model.APIListItem, int, error) {
	groupList := make([]string, 0)
	var err error
	//获取传入的groupUUID下包括子分组的所有UUID
	if groupUUID != "" {
		groupList, err = a.commonGroup.groupUUIDS(ctx, namespaceID, apiName, apiName, groupUUID)
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
	apiList := make([]*model.APIListItem, 0, len(apis))
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
			appName, _ := a.extAppService.getExtAppName(ctx, api.SourceID)
			source = appName
			if api.SourceLabel != "" {
				source = appName + "-" + api.SourceLabel
			}
		}

		isDelete, _ := a.isApiCanDelete(ctx, api.Id)

		item := &model.APIListItem{
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

func (a *apiService) GetAPIInfo(ctx context.Context, namespaceID int, uuid string) (*model.APIInfo, error) {
	api, err := a.apiStore.GetByUUID(ctx, namespaceID, uuid)
	if err != nil {
		log.Errorf("GetAPIInfo-apiStore.GetByUUID namespaceId:%d,uuid:%s,err=%s", namespaceID, uuid, err.Error())
		return nil, err
	}
	return &model.APIInfo{API: api}, nil
}

func (a *apiService) GetAPIInfoByGroupUUID(ctx context.Context, namespaceID int, groupUUID string) ([]*model.APIInfo, error) {
	apis, err := a.apiStore.GetListByGroupID(ctx, namespaceID, groupUUID)
	if err != nil {
		return nil, err
	}

	list := make([]*model.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &model.APIInfo{API: api})
	}

	return list, nil
}

func (a *apiService) GetAPIInfoByUUIDS(ctx context.Context, namespaceID int, uuids []string) ([]*model.APIInfo, error) {
	apis, err := a.apiStore.GetByUUIDs(ctx, namespaceID, uuids)
	if err != nil {
		return nil, err
	}

	list := make([]*model.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &model.APIInfo{API: api})
	}

	return list, nil
}

func (a *apiService) GetAPIInfoByPath(ctx context.Context, namespaceID int, path string) ([]*model.APIInfo, error) {
	apis, err := a.apiStore.GetByPath(ctx, namespaceID, path)
	if err != nil {
		return nil, err
	}

	list := make([]*model.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &model.APIInfo{API: api})
	}

	return list, nil
}

func (a *apiService) GetAPIListItemByUUIDS(ctx context.Context, namespaceID int, uuids []string) ([]*model.APIListItem, error) {
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

	list := make([]*model.APIListItem, 0, len(apis))
	for _, api := range apis {
		version := versionMap[api.Id]
		item := &model.APIListItem{
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

func (a *apiService) GetAPIListItemAll(ctx context.Context, namespaceID int) ([]*model.APIListItem, error) {
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

	list := make([]*model.APIListItem, 0, len(apis))
	for _, api := range apis {
		version := versionMap[api.Id]
		item := &model.APIListItem{
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

func (a *apiService) GetAPIInfoAll(ctx context.Context, namespaceID int) ([]*model.APIInfo, error) {
	apis, err := a.apiStore.GetListByName(ctx, namespaceID, "")
	if err != nil {
		return nil, err
	}

	list := make([]*model.APIInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &model.APIInfo{API: api})
	}

	return list, nil
}

// GetAPIsForSync 同步api时使用
func (a *apiService) GetAPIsForSync(ctx context.Context, namespaceID int) ([]*model.APIVersionInfo, error) {
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

	list := make([]*model.APIVersionInfo, 0, len(apis))
	for _, api := range apis {
		list = append(list, &model.APIVersionInfo{
			Api:     api,
			Version: versionMap[api.Id],
		})
	}

	return list, nil
}

func (a *apiService) getAPIVersions(ctx context.Context, apiIds []int) (map[int]*entry.APIVersion, error) {
	versions, err := a.apiVersion.GetAPIVersionByApiIds(ctx, apiIds)
	if err != nil {
		return nil, err
	}
	return common.SliceToMap(versions, func(t *entry.APIVersion) int {
		return t.ApiID
	}), nil
}

func (a *apiService) GetAPIVersionInfo(ctx context.Context, namespaceID int, uuid string) (*model.APIVersionInfo, error) {
	api, err := a.apiStore.GetByUUID(ctx, namespaceID, uuid)
	if err != nil {
		return nil, err
	}

	version, err := a.GetLatestAPIVersion(ctx, api.Id)
	if err != nil {
		return nil, err
	}

	info := &model.APIVersionInfo{
		Api:     api,
		Version: version,
	}

	return info, nil
}

func (a *apiService) CreateAPI(ctx context.Context, namespaceID int, operator int, input *dto.APIInfo) error {

	if err := a.CheckAPIReduplicative(ctx, namespaceID, "", input); err != nil {
		return err
	}

	if input.UUID == "" {
		input.UUID = uuid.New()
	}

	input.UUID = strings.ToLower(input.UUID)

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
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

	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		t := time.Now()
		apiInfo := &entry.API{
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

		serviceID, err := a.service.GetServiceIDByName(txCtx, namespaceID, input.ServiceName)
		if err != nil {
			return err
		}
		//添加版本信息
		apiVersionInfo := &entry.APIVersion{
			ApiID:       apiInfo.Id,
			NamespaceID: namespaceID,
			APIVersionConfig: entry.APIVersionConfig{
				Driver:           input.Driver,
				RequestPath:      input.RequestPath,
				RequestPathLabel: input.RequestPathLabel,
				ServiceID:        serviceID,
				ServiceName:      input.ServiceName,
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

		if err = a.apiHistory.HistoryAdd(txCtx, namespaceID, apiInfo.Id, &entry.ApiHistoryInfo{
			Api:    *apiInfo,
			Config: apiVersionInfo.APIVersionConfig,
		}, operator); err != nil {
			return err
		}

		stat := &entry.APIStat{
			APIID:     apiInfo.Id,
			VersionID: apiVersionInfo.Id,
		}

		//添加版本关联原表信息
		if err = a.apiStat.Save(txCtx, stat); err != nil {
			return err
		}

		//quote更新所引用的服务
		quoteMap := make(map[entry.QuoteTargetKindType][]int)
		quoteMap[entry.QuoteTargetKindTypeService] = append(quoteMap[entry.QuoteTargetKindTypeService], serviceID)

		return a.quoteStore.Set(txCtx, apiInfo.Id, entry.QuoteKindTypeAPI, quoteMap)
	})

}

func (a *apiService) UpdateAPI(ctx context.Context, namespaceID int, operator int, input *dto.APIInfo) error {
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

	err = a.lockService.lock(lockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.unlock(lockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceID, input.UUID)
	if err != nil {
		return err
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

	oldValue := entry.ApiHistoryInfo{
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

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.ApiName,
	})

	return a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		//修改基础数据
		if _, err = a.apiStore.Update(txCtx, apiInfo); err != nil {
			return err
		}
		serviceID, err := a.service.GetServiceIDByName(txCtx, namespaceID, input.ServiceName)
		if err != nil {
			return err
		}

		latestVersionConfig := entry.APIVersionConfig{
			Driver:           input.Driver,
			RequestPath:      input.RequestPath,
			RequestPathLabel: input.RequestPathLabel,
			ServiceID:        serviceID,
			ServiceName:      input.ServiceName,
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
			apiVersionInfo := &entry.APIVersion{
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
			stat = &entry.APIStat{
				APIID:     apiInfo.Id,
				VersionID: apiVersionInfo.Id,
			}
			if err = a.apiStat.Save(txCtx, stat); err != nil {
				return err
			}

			//quote更新所引用的服务
			quoteMap := make(map[entry.QuoteTargetKindType][]int)
			quoteMap[entry.QuoteTargetKindTypeService] = append(quoteMap[entry.QuoteTargetKindTypeService], serviceID)
			if err = a.quoteStore.Set(txCtx, apiInfo.Id, entry.QuoteKindTypeAPI, quoteMap); err != nil {
				return err
			}
		}

		newValue := entry.ApiHistoryInfo{
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

	err = a.lockService.lock(lockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.unlock(lockNameAPI, apiInfo.Id)

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

	OldValue := entry.ApiHistoryInfo{
		Api:    *apiInfo,
		Config: version.APIVersionConfig,
	}

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: uuid,
		Name: apiInfo.Name,
	})

	err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = a.apiStore.Delete(txCtx, apiInfo.Id); err != nil {
			return err
		}
		delMap := make(map[string]interface{})
		delMap["`kind`"] = "api"
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

		return a.quoteStore.DelBySource(txCtx, apiInfo.Id, entry.QuoteKindTypeAPI)
	})
	if err != nil {
		return err
	}

	a.lockService.deleteLock(lockNameAPI, apiInfo.Id)
	return nil
}

func (a *apiService) BatchOnline(ctx context.Context, namespaceId int, operator int, onlineToken string) ([]*model.BatchListItem, error) {
	//判断uuid和operator是一致的
	key := a.batchApiCache.Key(onlineToken)
	task, err := a.batchApiCache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	//篡改审计日志的请求body
	ginContext, ok := ctx.(*gin.Context)
	if ok {
		ginContext.Set("logBody", string(task.Data))
	}

	//若提交上线的操作人和成功检测的操作人不一致 则报错
	if task.Operator != operator {
		return nil, errors.New("operator is invalid. ")
	}

	conf := new(model.APIBatchConf)
	_ = json.Unmarshal(task.Data, conf)

	apiList := make([]*entry.API, 0, len(conf.ApiUUIDs))

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
		//确认所有apiUUID，clusterName均存在
		for _, uid := range conf.ApiUUIDs {
			api, err := a.apiStore.GetByUUID(ctx, namespaceId, uid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("api doesn't exist. uuid:%s", uid)
				}
				return err
			}
			apiList = append(apiList, api)
		}
		return nil
	})

	clusterList := make([]*model.Cluster, 0, len(conf.ClusterNames))
	group.Go(func() error {
		clusters, err := a.clusterService.QueryListByNamespaceId(ctx, namespaceId)
		if err != nil {
			return err
		}
		clusterMap := common.SliceToMap(clusters, func(t *model.Cluster) string {
			return t.Name
		})
		for _, clusterName := range conf.ClusterNames {
			if cluster, ok := clusterMap[clusterName]; ok {
				if cluster.Status == 2 || cluster.Status == 3 {
					return fmt.Errorf("cluster status is abnormal. cluster_name:%s", clusterName)
				}
				clusterList = append(clusterList, cluster)
			} else {
				return fmt.Errorf("cluster doesn't exist. cluster_name:%s", clusterName)
			}

		}
		return nil
	})

	if err = group.Wait(); err != nil {
		return nil, err
	}

	t := time.Now()

	//逐个处理api上线
	onlineList := make([]*model.BatchListItem, 0, len(apiList)*len(clusterList))
	for _, api := range apiList {
		err = a.lockService.lock(lockNameAPI, api.Id)
		if err != nil {
			for _, cluster := range clusterList {
				item := &model.BatchListItem{
					APIName:    api.Name,
					ClusterEnv: fmt.Sprintf("%s_%s", cluster.Name, cluster.Env),
					Status:     false,
					Result:     err.Error(),
				}
				onlineList = append(onlineList, item)
			}
			a.lockService.unlock(lockNameAPI, api.Id)
			continue
		}
		//确保api没被删除
		_, err = a.apiStore.Get(ctx, api.Id)
		if err != nil {
			//API被删除
			for _, cluster := range clusterList {
				item := &model.BatchListItem{
					APIName:    api.Name,
					ClusterEnv: fmt.Sprintf("%s_%s", cluster.Name, cluster.Env),
					Status:     false,
					Result:     err.Error(),
				}
				onlineList = append(onlineList, item)
			}
			a.lockService.unlock(lockNameAPI, api.Id)
			continue
		}

		for _, cluster := range clusterList {
			item := &model.BatchListItem{
				APIName:    api.Name,
				ClusterEnv: fmt.Sprintf("%s_%s", cluster.Name, cluster.Env),
				Status:     true,
				Result:     "",
			}

			//获取当前的版本
			runtime, err := a.apiRuntime.GetForCluster(ctx, api.Id, cluster.Id)
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
				if !a.service.IsOnline(ctx, cluster.Id, latest.ServiceID) {
					item.Status = false
					item.Result = fmt.Sprintf("绑定的%s未上线到%s", latest.ServiceName, cluster.Name)
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
				client, err := a.apintoClient.GetClient(ctx, cluster.Id)
				if err != nil {
					item.Status = false
					item.Result = fmt.Sprintf("连接集群失败, err: %s", err.Error())
					return nil
				}

				//封装router配置
				apiDriver := a.GetAPIDriver(latest.Driver)
				routerConfig := apiDriver.ToApinto(api.UUID, api.Desc, false, latest.Method, latest.RequestPath, latest.RequestPathLabel, latest.ProxyPath, strings.ToLower(latest.ServiceName), latest.Timeout, latest.Retry, latest.EnableWebsocket, latest.Match, latest.Header)

				//未上线
				if runtime == nil {
					runtime = &entry.APIRuntime{
						NamespaceId: namespaceId,
						ApiID:       api.Id,
						ClusterID:   cluster.Id,
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

		a.lockService.unlock(lockNameAPI, api.Id)
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

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Name:        strings.Join(logApiNameList, ","),
		ClusterName: strings.Join(logCLNameList, ","),
		PublishType: 1,
	})

	return onlineList, nil
}

func (a *apiService) BatchOffline(ctx context.Context, namespaceId int, operator int, apiUUIDs, clusterNames []string) ([]*model.BatchListItem, error) {

	group, _ := errgroup.WithContext(ctx)

	apiList := make([]*entry.API, 0, len(apiUUIDs))
	group.Go(func() error {
		//确认所有apiUUID，clusterName均存在
		for _, uid := range apiUUIDs {
			api, err := a.apiStore.GetByUUID(ctx, namespaceId, uid)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return fmt.Errorf("api doesn't exist. uuid:%s", uid)
				}
				return err
			}
			apiList = append(apiList, api)
		}
		return nil
	})

	clusterList := make([]*model.Cluster, 0, len(clusterNames))
	group.Go(func() error {
		clusters, err := a.clusterService.QueryListByNamespaceId(ctx, namespaceId)
		if err != nil {
			return err
		}
		clusterMap := common.SliceToMap(clusters, func(t *model.Cluster) string {
			return t.Name
		})
		for _, clusterName := range clusterNames {
			if cluster, ok := clusterMap[clusterName]; ok {
				if cluster.Status == 2 || cluster.Status == 3 {
					return fmt.Errorf("cluster status is abnormal. cluster_name:%s", clusterName)
				}
				clusterList = append(clusterList, cluster)
			} else {
				return fmt.Errorf("cluster doesn't exist. cluster_name:%s", clusterName)
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	//逐个处理api下线，已经下线或者未上线的不进行操作
	offlineList := make([]*model.BatchListItem, 0, len(apiList)*len(clusterList))
	for _, api := range apiList {
		err := a.lockService.lock(lockNameAPI, api.Id)
		if err != nil {
			for _, cluster := range clusterList {
				item := &model.BatchListItem{
					APIName:    api.Name,
					ClusterEnv: fmt.Sprintf("%s_%s", cluster.Name, cluster.Env),
					Status:     false,
					Result:     err.Error(),
				}
				offlineList = append(offlineList, item)
			}
			a.lockService.unlock(lockNameAPI, api.Id)
			continue
		}
		latestApi, err := a.apiStore.Get(ctx, api.Id)
		if err != nil {
			a.lockService.unlock(lockNameAPI, api.Id)
			return nil, err
		}

		for _, cluster := range clusterList {
			//获取当前的版本
			runtime, err := a.apiRuntime.GetForCluster(ctx, api.Id, cluster.Id)
			if err != nil && err != gorm.ErrRecordNotFound {
				a.lockService.unlock(lockNameAPI, api.Id)
				return nil, err
			}

			item := &model.BatchListItem{
				APIName:    latestApi.Name,
				ClusterEnv: fmt.Sprintf("%s_%s", cluster.Name, cluster.Env),
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
					client, err := a.apintoClient.GetClient(ctx, cluster.Id)
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

		a.lockService.unlock(lockNameAPI, api.Id)
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

	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Name:        strings.Join(logApiNameList, ","),
		ClusterName: strings.Join(logCLNameList, ","),
		PublishType: 2,
	})

	return offlineList, nil
}

func (a *apiService) BatchOnlineCheck(ctx context.Context, namespaceId int, operator int, apiUUIDs, clusterNames []string) ([]*model.BatchOnlineCheckListItem, string, error) {
	//确认所有apiUUID，clusterName均存在
	apiList := make([]*entry.API, 0, len(apiUUIDs))
	apiIds := make([]int, 0, len(apiUUIDs))

	group, _ := errgroup.WithContext(ctx)
	group.Go(func() error {
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

	clusterList := make([]*model.Cluster, 0, len(clusterNames))
	group.Go(func() error {
		clusters, err := a.clusterService.QueryListByNamespaceId(ctx, namespaceId)
		if err != nil {
			return err
		}
		clusterMap := common.SliceToMap(clusters, func(t *model.Cluster) string {
			return t.Name
		})
		for _, clusterName := range clusterNames {
			if cluster, ok := clusterMap[clusterName]; ok {
				if cluster.Status == 2 || cluster.Status == 3 {
					return fmt.Errorf("cluster status is abnormal. cluster_name:%s", clusterName)
				}
				clusterList = append(clusterList, cluster)
			} else {
				return fmt.Errorf("cluster doesn't exist. cluster_name:%s", clusterName)
			}
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, "", err
	}

	isAllOnline := true
	checkSet := make(map[int]struct{}) //serviceId集合，用于对检查列表的去重
	checkList := make([]*model.BatchOnlineCheckListItem, 0, len(apiList)*len(clusterList))

	versionMap, err := a.getAPIVersions(ctx, apiIds)
	if err != nil {
		return nil, "", err
	}

	//确认每个api对应的cluster所配置的serviceID的上线情况
	for _, api := range apiList {
		version := versionMap[api.Id]
		if _, has := checkSet[version.ServiceID]; has {
			continue
		}
		checkSet[version.ServiceID] = struct{}{}

		for _, cluster := range clusterList {
			item := &model.BatchOnlineCheckListItem{
				ServiceName: version.ServiceName,
				ClusterEnv:  fmt.Sprintf("%s%s", cluster.Name, cluster.Env),
				Status:      true,
				Solution:    &model.Router{},
			}

			if isOnline := a.service.IsOnline(ctx, cluster.Id, version.ServiceID); !isOnline {
				isAllOnline = false
				item.Status = false
				item.Result = fmt.Sprintf("%s未上线到%s", version.ServiceName, cluster.Name)
				item.Solution.Name = model.RouterNameServiceOnline
				item.Solution.Params = map[string]string{"service_name": version.ServiceName}
			}
			checkList = append(checkList, item)
		}
	}

	//若所有的API均已上线，则生成一个UUID
	onlineToken := ""
	if isAllOnline {
		onlineToken = uuid.New()

		taskData := &model.APIBatchConf{
			ApiUUIDs:     apiUUIDs,
			ClusterNames: clusterNames,
		}

		data, _ := json.Marshal(taskData)
		task := &model.BatchOnlineCheckTask{
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

func (a *apiService) OnlineList(ctx context.Context, namespaceId int, uuid string) ([]*model.APIOnlineListItem, error) {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	//获取工作空间下的所有集群
	clusters, err := a.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *model.Cluster) int {
		return t.Id
	})

	//获取当前服务发现下集群运行的版本
	runtimes, err := a.apiRuntime.GetByTarget(ctx, apiInfo.Id)
	if err != nil {
		return nil, err
	}
	runtimeMaps := common.SliceToMap(runtimes, func(t *entry.APIRuntime) int {
		return t.ClusterID
	})

	//获取操作人用户列表
	operatorList := common.SliceToSliceIds(runtimes, func(t *entry.APIRuntime) int {
		return t.Operator
	})

	userInfoMaps, err := a.userInfoService.GetUserInfoMaps(ctx, operatorList...)
	if err != nil {
		return nil, err
	}

	list := make([]*model.APIOnlineListItem, 0, len(clusters))

	latestVersion, err := a.GetLatestAPIVersion(ctx, apiInfo.Id)
	if err != nil {
		return nil, err
	}

	for _, cluster := range clusterMaps {
		apiOnline := &model.APIOnlineListItem{
			ClusterName: cluster.Name,
			ClusterEnv:  cluster.Env,
			Status:      1, //默认为未上线状态
		}
		if runtime, ok := runtimeMaps[cluster.Id]; ok {

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

func (a *apiService) OnlineAPI(ctx context.Context, namespaceId, operator int, uuid, clusterName string) (*model.Router, error) {
	apiInfo, err := a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	err = a.lockService.lock(lockNameAPI, apiInfo.Id)
	if err != nil {
		return nil, err
	}
	defer a.lockService.unlock(lockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return nil, err
	}

	apiID := apiInfo.Id
	t := time.Now()

	//获取当前集群信息
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	latestVersion, err := a.GetLatestAPIVersion(ctx, apiID)
	if err != nil {
		return nil, err
	}

	router := &model.Router{
		Name:   model.RouterNameServiceOnline,
		Params: make(map[string]string),
	}
	router.Params["service_name"] = latestVersion.ServiceName

	//判断上游服务有没有上线
	if !a.service.IsOnline(ctx, cluster.Id, latestVersion.ServiceID) {
		return router, errors.New(fmt.Sprintf("绑定的%s未上线到%s", latestVersion.ServiceName, clusterName))
	}

	//获取当前运行的版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiID, cluster.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//发布到apinto
	client, err := a.apintoClient.GetClient(ctx, cluster.Id)
	if err != nil {
		return nil, err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        uuid,
		Name:        apiInfo.Name,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
		PublishType: 1,
	})

	//事务
	err = a.apiStore.Transaction(ctx, func(txCtx context.Context) error {

		apiDriver := a.GetAPIDriver(latestVersion.Driver)
		routerConfig := apiDriver.ToApinto(apiInfo.UUID, apiInfo.Desc, false, latestVersion.Method, latestVersion.RequestPath, latestVersion.RequestPathLabel, latestVersion.ProxyPath, strings.ToLower(latestVersion.ServiceName), latestVersion.Timeout, latestVersion.Retry, latestVersion.EnableWebsocket, latestVersion.Match, latestVersion.Header)
		if runtime == nil {
			runtime = &entry.APIRuntime{
				NamespaceId: namespaceId,
				ApiID:       apiID,
				ClusterID:   cluster.Id,
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

func (a *apiService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
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
		routerConfig := a.GetAPIDriver(version.Driver).ToApinto(apiInfo.UUID, apiInfo.Desc, false, version.Method, version.RequestPath, version.RequestPathLabel, version.ProxyPath, strings.ToLower(version.ServiceName), version.Timeout, version.Retry, version.EnableWebsocket, version.Match, version.Header)

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

	err = a.lockService.lock(lockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.unlock(lockNameAPI, apiInfo.Id)

	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}

	//获取当前集群信息
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	//获取当前的版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiInfo.Id, cluster.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}

	t := time.Now()

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:        uuid,
		Name:        apiInfo.Name,
		ClusterId:   cluster.Id,
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
		client, err := a.apintoClient.GetClient(ctx, cluster.Id)
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
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	err = a.lockService.lock(lockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.unlock(lockNameAPI, apiInfo.Id)
	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}
	//获取当前版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiInfo.Id, cluster.Id)
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
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:          uuid,
		Name:          apiInfo.Name,
		ClusterId:     cluster.Id,
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
		client, err := a.apintoClient.GetClient(ctx, cluster.Id)
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
	cluster, err := a.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	err = a.lockService.lock(lockNameAPI, apiInfo.Id)
	if err != nil {
		return err
	}
	defer a.lockService.unlock(lockNameAPI, apiInfo.Id)
	apiInfo, err = a.apiStore.GetByUUID(ctx, namespaceId, uuid)
	if err != nil {
		return err
	}
	//获取当前版本
	runtime, err := a.apiRuntime.GetForCluster(ctx, apiInfo.Id, cluster.Id)
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
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid:          uuid,
		Name:          apiInfo.Name,
		ClusterId:     cluster.Id,
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
		client, err := a.apintoClient.GetClient(ctx, cluster.Id)
		if err != nil {
			return err
		}
		return client.ForRouter().Patch(apiInfo.UUID+"@router", map[string]interface{}{"disable": true})
	})
}

func (a *apiService) GetSource(ctx context.Context) ([]*model.SourceListItem, error) {
	items := make([]*model.SourceListItem, 0, 2)

	sourceList, err := a.apiStore.GetSourceList(ctx)
	if err != nil {
		return items, err
	}

	//对来源列表进行排序
	sort.Sort(entry.APISourceList(sourceList))

	for _, source := range sourceList {
		title := ""
		if source.SourceType == enum.SourceSync {
			//若来源为同步，通过来源id查找外部应用名
			appName, _ := a.extAppService.getExtAppName(ctx, source.SourceID)
			title = appName
			if source.SourceLabel != "" {
				title = appName + "-" + source.SourceLabel
			}
		} else {
			title = enum.GetSourceTitle(source.SourceType)
		}
		items = append(items, &model.SourceListItem{
			Id:    fmt.Sprintf("%s:%d:%s", source.SourceType, source.SourceID, source.SourceLabel),
			Title: title,
		})
	}

	return items, nil
}

func (a *apiService) GetImportCheckList(ctx context.Context, namespaceId int, fileData []byte, groupID, serviceName, requestPrefix string) ([]*model.ImportAPIListItem, string, error) {
	//解析swagger3.0 TODO 写死解析3.0 等之后有其他格式再用driverManager，openAPI同步现在是用driverManager的
	swaggerConfig := new(model.SwaggerConfig)
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

	items := make([]*model.ImportAPIListItem, 0)
	for path, pathMap := range swaggerConfig.Paths {
		for method, info := range pathMap {
			name := info.Summary
			if name == "" {
				name = info.OperationID
			}
			if name == "" {
				name = method + "-" + requestPrefix + path
			}
			items = append(items, &model.ImportAPIListItem{
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

	apiMap := common.SliceToMapArray(apiList, func(t *model.CommonGroupApi) string {
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
			for _, api := range apis {
				for _, method := range api.Methods {
					if method == item.Method {
						item.Status = 2
						break
					}
				}
			}

		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Status > items[j].Status
	})

	redisDataItems := make([]*model.ImportAPIRedisDataItem, 0)
	for i, item := range items {
		item.Id = i + 1
		//如果状态不为正常，则不存进redis
		if item.Status != 1 {
			continue
		}

		api := &entry.API{
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

		apiInfo := &model.APIInfo{
			API:    api,
			Method: strings.Split(item.Method, ","),
		}

		redisDataItem := &model.ImportAPIRedisDataItem{
			ID:  item.Id,
			Api: apiInfo,
		}
		redisDataItems = append(redisDataItems, redisDataItem)
	}

	token := uuid.New()
	//数据存储到缓存
	key := a.importApiCache.Key(token)

	importAPIRedisData := &model.ImportAPIRedisData{
		Apis:        redisDataItems,
		ServiceName: serviceName,
		GroupID:     groupID,
	}

	if err = a.importApiCache.Set(ctx, key, importAPIRedisData, time.Hour*8); err != nil {
		return nil, "", err
	}

	return items, token, nil
}

func (a *apiService) ImportAPI(ctx context.Context, namespaceId, operator int, input *dto.ImportAPIInfos) error {

	key := a.importApiCache.Key(input.Token)
	apiData, err := a.importApiCache.Get(ctx, key)
	if err != nil {
		return err
	}

	serviceID, err := a.service.GetServiceIDByName(ctx, namespaceId, apiData.ServiceName)
	if err != nil {
		return err
	}

	maps := common.SliceToMap(apiData.Apis, func(t *model.ImportAPIRedisDataItem) int {
		return t.ID
	})

	createApis := make([]*model.APIInfo, 0, len(input.Apis))
	logAPINames := make([]string, 0, len(input.Apis))
	for _, api := range input.Apis {
		if v, ok := maps[api.Id]; ok {
			if api.Name != "" {
				v.Api.Name = api.Name
			}
			// TODO 现在只能修改apiName， 请求路径和描述以后可能要改
			createApis = append(createApis, v.Api)
			logAPINames = append(logAPINames, v.Api.Name)
		} else {
			return errors.New(fmt.Sprintf("序号为%d的数据不存在", api.Id))
		}
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
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
			apiVersionInfo := &entry.APIVersion{
				ApiID:       apiInfo.Id,
				NamespaceID: namespaceId,
				APIVersionConfig: entry.APIVersionConfig{
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
					Match:            []*entry.MatchConf{},
					Header:           []*entry.ProxyHeader{},
				},
				Operator:   operator,
				CreateTime: t,
			}

			if err = a.apiVersion.Save(txCtx, apiVersionInfo); err != nil {
				return err
			}

			if err = a.apiHistory.HistoryAdd(txCtx, namespaceId, apiInfo.Id, &entry.ApiHistoryInfo{
				Api:    *apiInfo.API,
				Config: apiVersionInfo.APIVersionConfig,
			}, operator); err != nil {
				return err
			}

			stat := &entry.APIStat{
				APIID:     apiInfo.Id,
				VersionID: apiVersionInfo.Id,
			}

			//添加版本关联原表信息
			if err = a.apiStat.Save(txCtx, stat); err != nil {
				return err
			}

			//quote更新所引用的服务
			quoteMap := make(map[entry.QuoteTargetKindType][]int)
			quoteMap[entry.QuoteTargetKindTypeService] = append(quoteMap[entry.QuoteTargetKindTypeService], serviceID)

			if err = a.quoteStore.Set(txCtx, apiInfo.Id, entry.QuoteKindTypeAPI, quoteMap); err != nil {
				return err
			}

		}
		return nil
	})
}

func (a *apiService) GetAPIListByName(ctx context.Context, namespaceId int, name string) ([]*model.CommonGroupApi, error) {
	apiList, err := a.apiStore.GetListByName(ctx, namespaceId, name)
	if err != nil {
		return nil, err
	}
	groupAPIs := make([]*model.CommonGroupApi, 0, len(apiList))

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
		groupApi := &model.CommonGroupApi{
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

func (a *apiService) GetAPIListByServiceName(ctx context.Context, namespaceId int, serviceNames []string) ([]*model.APIInfo, error) {

	var err error
	groupAPIs := make([]*model.APIInfo, 0)

	for _, serviceName := range serviceNames {

		target := 0
		if serviceName != "" {
			serviceId, err := a.service.GetServiceIDByName(ctx, namespaceId, serviceName)
			if err != nil {
				return nil, err
			}
			target = serviceId
		}

		apiList := make([]*entry.API, 0)

		if target > 0 {
			quote, err := a.quoteStore.GetTargetQuote(ctx, target, entry.QuoteTargetKindTypeService)
			if err != nil {
				return nil, err
			}
			apiList, err = a.apiStore.GetByIds(ctx, namespaceId, quote[entry.QuoteKindTypeAPI])
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
			groupApi := &model.APIInfo{
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

func (a *apiService) GetLatestAPIVersion(ctx context.Context, apiId int) (*entry.APIVersion, error) {
	stat, err := a.apiStat.Get(ctx, apiId)
	if err != nil {
		return nil, err
	}
	return a.apiVersion.Get(ctx, stat.VersionID)
}

// CheckAPIReduplicative 检测API配置是否重复，不可同名同request_url同method
func (a *apiService) CheckAPIReduplicative(ctx context.Context, namespaceID int, uuid string, input *dto.APIInfo) error {
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

func (a *apiService) isAPIVersionConfChange(latest entry.APIVersionConfig, current entry.APIVersionConfig) bool {
	return !reflect.DeepEqual(latest, current)
}

func (a *apiService) GetAPIDriver(driverName string) driver.IAPIDriver {
	return a.apiManager.GetDriver(driverName)
}

func (a *apiService) GetAPINameByID(ctx context.Context, apiID int) (string, error) {
	info, err := a.apiStore.Get(ctx, apiID)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}

func (a *apiService) GetAPIRemoteOptions(ctx context.Context, namespaceID, pageNum, pageSize int, keyword, groupUuid string) ([]*model.RemoteApis, int, error) {
	groupList := make([]string, 0)
	var err error
	//获取传入的groupUUID下包括子分组的所有UUID
	if groupUuid != "" {
		groupList, err = a.commonGroup.groupUUIDS(ctx, namespaceID, apiName, apiName, groupUuid)
		if err != nil {
			return nil, 0, err
		}
		groupList = append(groupList, groupUuid)
	}

	groups, err := a.commonGroup.GroupListAll(ctx, namespaceID, apiName, apiName)
	if err != nil {
		return nil, 0, err
	}

	//获取指定目录下所有API
	apis, total, err := a.apiStore.GetListPageByGroupIDs(ctx, namespaceID, pageNum, pageSize, groupList, nil, keyword)
	if err != nil {
		return nil, 0, err
	}
	apiList := make([]*model.RemoteApis, 0, len(apis))

	groupUUIDMap := common.SliceToMap(groups, func(t *entry.CommonGroup) string {
		return t.Uuid
	})
	groupIdMap := common.SliceToMap(groups, func(t *entry.CommonGroup) int {
		return t.Id
	})

	for _, api := range apis {
		//version, err := a.GetLatestAPIVersion(ctx, api.Id)
		//if err != nil {
		//	return nil, 0, err
		//}
		parentGroupName := &[]string{}

		a.commonGroup.parentGroupName(api.GroupUUID, groupUUIDMap, groupIdMap, parentGroupName)

		item := &model.RemoteApis{
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

func (a *apiService) GetAPIRemoteByUUIDS(ctx context.Context, namespace int, uuids []string) ([]*model.RemoteApis, error) {

	groups, err := a.commonGroup.GroupListAll(ctx, namespace, apiName, apiName)
	if err != nil {
		return nil, err
	}

	groupUUIDMap := common.SliceToMap(groups, func(t *entry.CommonGroup) string {
		return t.Uuid
	})
	groupIdMap := common.SliceToMap(groups, func(t *entry.CommonGroup) int {
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

	apiList := make([]*model.RemoteApis, 0, len(apis))
	for _, api := range apis {
		version := versionMap[api.Id]

		parentGroupName := &[]string{}
		a.commonGroup.parentGroupName(api.GroupUUID, groupUUIDMap, groupIdMap, parentGroupName)

		item := &model.RemoteApis{
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
