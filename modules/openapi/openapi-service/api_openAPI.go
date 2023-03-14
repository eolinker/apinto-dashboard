package openapi_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	driver_manager "github.com/eolinker/apinto-dashboard/driver-manager"
	"github.com/eolinker/apinto-dashboard/modules/api"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	store2 "github.com/eolinker/apinto-dashboard/modules/api/store"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/group"
	"github.com/eolinker/apinto-dashboard/modules/group/group-model"
	"github.com/eolinker/apinto-dashboard/modules/openapi"
	"github.com/eolinker/apinto-dashboard/modules/openapi/openapi-dto"
	"github.com/eolinker/apinto-dashboard/modules/openapi/openapi-model"
	"github.com/eolinker/apinto-dashboard/modules/openapp"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type apiOpenAPIService struct {
	apiStore   store2.IAPIStore
	apiStat    store2.IAPIStatStore
	apiVersion store2.IAPIVersionStore
	quoteStore quote_store.IQuoteStore
	apiHistory store2.IApiHistoryStore

	apiService           api.IAPIService
	service              upstream.IService
	commonGroup          group.ICommonGroupService
	extAppService        openapp.IExternalApplicationService
	apiSyncFormatManager driver_manager.IAPISyncFormatManager
}

func newAPIOpenAPIService() openapi.IAPIOpenAPIService {
	as := &apiOpenAPIService{}
	bean.Autowired(&as.apiStore)
	bean.Autowired(&as.apiStat)
	bean.Autowired(&as.apiVersion)
	bean.Autowired(&as.quoteStore)

	bean.Autowired(&as.apiService)
	bean.Autowired(&as.service)
	bean.Autowired(&as.commonGroup)
	bean.Autowired(&as.apiHistory)
	bean.Autowired(&as.extAppService)

	bean.Autowired(&as.apiSyncFormatManager)

	return as
}

func (a *apiOpenAPIService) SyncImport(ctx context.Context, namespaceID, appID int, data *openapi_dto.SyncImportData) ([]*apimodel.ImportAPIListItem, error) {
	formatDriver := a.apiSyncFormatManager.GetDriver(data.Format)
	if formatDriver == nil {
		return nil, fmt.Errorf("format %s is illegal. ", data.Format)
	}

	var err error
	if data.Prefix, err = common.CheckAndFormatPath(data.Prefix); err != nil {
		return nil, errors.New("prefix is illegal. ")
	}

	//从format driver获取apis，然后与现有api进行对比去重
	importApis, err := formatDriver.FormatAPI(data.Content, namespaceID, appID, data.GroupUUID, data.Prefix, data.Label)
	if err != nil {
		return nil, fmt.Errorf("decode file fail. err: %s", err)
	}

	//检查分组id存不存在
	isExist, err := a.commonGroup.IsGroupExist(ctx, data.GroupUUID)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.New("分组不存在")
	}

	//检查服务，没有就用server创建，服务名使用data.ServiceName
	serviceNotExit := false
	serviceID, err := a.service.GetServiceIDByName(ctx, namespaceID, data.ServiceName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		serviceNotExit = true
		if data.Server == nil {
			return nil, fmt.Errorf("Server can't be null when service doesn't exist. ")
		}
		if len(data.Server.Nodes) == 0 {
			return nil, fmt.Errorf("Server.Nodes can't be null when service doesn't exist. ")
		}
	}

	//获取所有API
	allApi, err := a.apiService.GetAPIsForSync(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	apiMap := common.SliceToMapArray(allApi, func(t *apimodel.APIVersionInfo) string {
		return t.Api.RequestPath
	})

	apiList := make([]*apimodel.APIInfo, 0, len(importApis))
	items := make([]*apimodel.ImportAPIListItem, 0)
	for _, importAPI := range importApis {
		apiMethod := importAPI.Method[0]
		item := &apimodel.ImportAPIListItem{
			Name:   importAPI.Name,
			Method: apiMethod,
			Path:   importAPI.RequestPathLabel,
			Desc:   importAPI.Desc,
			Status: 1,
		}
		//暂不支持OPTIONS和TRACE
		if apiMethod == "OPTIONS" || apiMethod == "TRACE" {
			item.Status = 2 //冲突
			items = append(items, item)
			continue
		}

		if _, err = common.CheckAndFormatPath(item.Path); err != nil {
			item.Status = 3 //无效path
			items = append(items, item)
			continue
		}

		//判断路径和method是否冲突
		if apis, ok := apiMap[importAPI.RequestPath]; ok {
		A:
			for _, api := range apis {
				for _, method := range api.Version.Method {
					if method == apiMethod {
						item.Status = 2 //冲突
						break A
					}
				}
			}
		}
		if item.Status == 1 {
			apiList = append(apiList, importAPI)
		}
		items = append(items, item)
	}

	//对返回的结果列表进行排序，冲突的排在前面
	sort.Slice(items, func(i, j int) bool {
		return items[i].Status > items[j].Status
	})

	//可导入的api为0时返回
	if len(apiList) == 0 {
		return items, nil
	}

	t := time.Now()
	return items, a.apiStore.Transaction(ctx, func(txCtx context.Context) error {
		//若服务不存在则创建
		if serviceNotExit {
			staticAddrs := make([]string, 0, len(data.Server.Nodes))
			staticConf := make([]*openapi_model.ServiceStaticConf, 0, len(data.Server.Nodes))
			for _, node := range data.Server.Nodes {
				if node.Weight == 0 {
					node.Weight = 1
				}
				staticAddrs = append(staticAddrs, fmt.Sprintf("%s weight=%d", strings.TrimSpace(node.Url), node.Weight))
				staticConf = append(staticConf, &openapi_model.ServiceStaticConf{
					Addr:   strings.TrimSpace(node.Url),
					Weight: node.Weight,
				})
			}
			config := &openapi_model.ServiceStaticDriverConf{
				UseVariable: false,
				StaticConf:  staticConf,
			}
			configData, _ := json.Marshal(config)

			serviceID, err = a.service.CreateService(txCtx, namespaceID, 0, &dto.ServiceInfo{
				Name:        data.ServiceName,
				UUID:        uuid.New(),
				Desc:        "",
				Scheme:      data.Server.Scheme,
				DiscoveryID: 0,
				DriverName:  "static", //静态
				FormatAddr:  strings.Join(staticAddrs, ","),
				Config:      string(configData),
				Timeout:     1000,
				Balance:     "round-robin", // 默认
			}, nil)
			if err != nil {
				return err
			}
		}

		//插入api
		for _, apiInfo := range apiList {
			if err = a.apiStore.Save(txCtx, apiInfo.API); err != nil {
				return err
			}
			//添加版本信息
			apiVersionInfo := &api_entry.APIVersion{
				ApiID:       apiInfo.Id,
				NamespaceID: namespaceID,
				APIVersionConfig: api_entry.APIVersionConfig{
					Driver:           "http", //默认
					RequestPath:      apiInfo.RequestPath,
					RequestPathLabel: apiInfo.RequestPathLabel,
					ServiceID:        serviceID,
					ServiceName:      data.ServiceName,
					Method:           apiInfo.Method,
					ProxyPath:        apiInfo.RequestPathLabel,
					Timeout:          10000,
					Retry:            0,
					EnableWebsocket:  false,
					Match:            []*api_entry.MatchConf{},
					Header:           []*api_entry.ProxyHeader{},
				},
				Operator:   0, //匿名
				CreateTime: t,
			}

			if err = a.apiVersion.Save(txCtx, apiVersionInfo); err != nil {
				return err
			}

			if err = a.apiHistory.HistoryAdd(txCtx, namespaceID, apiInfo.Id, &api_entry.ApiHistoryInfo{
				Api:    *apiInfo.API,
				Config: apiVersionInfo.APIVersionConfig,
			}, 0); err != nil {
				return err
			}

			stat := &api_entry.APIStat{
				APIID:     apiInfo.Id,
				VersionID: apiVersionInfo.Id,
			}

			//添加版本关联原表信息
			if err = a.apiStat.Save(txCtx, stat); err != nil {
				return err
			}

			//quote更新所引用的服务
			quoteMap := make(map[quote_entry.QuoteTargetKindType][]int)
			quoteMap[quote_entry.QuoteTargetKindTypeService] = append(quoteMap[quote_entry.QuoteTargetKindTypeService], serviceID)

			if err = a.quoteStore.Set(txCtx, apiInfo.Id, quote_entry.QuoteKindTypeAPI, quoteMap); err != nil {
				return err
			}

		}

		//修改外部应用的标签
		return a.extAppService.UpdateExtAPPTags(txCtx, namespaceID, appID, data.Label)
	})

}

func (a *apiOpenAPIService) GetSyncImportInfo(ctx context.Context, namespaceID int) ([]*openapi_model.ApiOpenAPIGroups, []*openapi_model.ApiOpenAPIService, []string, error) {
	//组装api分组
	groups, _, err := a.apiService.GetGroups(ctx, namespaceID, "", "")
	if err != nil {
		return nil, nil, nil, err
	}

	apiGroups := make([]*openapi_model.ApiOpenAPIGroups, 0, len(groups.CommonGroup))
	for _, commonGroup := range groups.CommonGroup {
		group := &openapi_model.ApiOpenAPIGroups{
			Uuid: commonGroup.Group.Uuid,
			Name: commonGroup.Group.Name,
		}
		apiGroups = append(apiGroups, group)
		a.subGroup(group, commonGroup.Subgroup)
	}

	//组装服务列表
	services, st, err := a.service.GetServiceRemoteOptions(ctx, namespaceID, 0, 0, "")
	if err != nil {
		return nil, nil, nil, err
	}

	serviceList := make([]*openapi_model.ApiOpenAPIService, 0, st)
	for _, s := range services {
		item := &openapi_model.ApiOpenAPIService{
			Name: s.Name,
			Desc: s.Desc,
		}
		serviceList = append(serviceList, item)
	}

	//获取支持的文件格式列表
	formatList := a.apiSyncFormatManager.List()

	return apiGroups, serviceList, formatList, nil
}

func (a *apiOpenAPIService) subGroup(val *openapi_model.ApiOpenAPIGroups, list []*group_model.CommonGroup) {
	if len(list) == 0 {
		return
	}
	for _, group := range list {
		commonGroup := &openapi_model.ApiOpenAPIGroups{
			Uuid: group.Group.Uuid,
			Name: group.Group.Name,
		}
		val.Children = append(val.Children, commonGroup)
		a.subGroup(commonGroup, group.Subgroup)
	}
}
