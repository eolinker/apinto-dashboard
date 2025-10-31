package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	navigation_service "github.com/eolinker/apinto-dashboard/modules/navigation"
	"github.com/eolinker/apinto-dashboard/modules/user"
	userModel "github.com/eolinker/apinto-dashboard/modules/user/user-model"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

var _ grpc_service.GetConsoleInfoServer = (*consoleInfoService)(nil)

const (
	professionService = "service"
	professionApp     = "app"
)

type consoleInfoService struct {
	namespaceService    namespace.INamespaceService
	apiService          api.IAPIService
	modulePluginService mpm3.IModuleService
	navigationService   navigation_service.INavigationService
	accessService       mpm3.IAccessService
	applicationService  application.IApplicationService
	dynamicService      dynamic.IDynamicService
	clusterService      cluster.IClusterService
	userService         user.IUserInfoService
	modulesCache        INavigationModulesCache
	grpc_service.UnimplementedGetConsoleInfoServer
}

func (c *consoleInfoService) SaveUserInfo(ctx context.Context, request *grpc_service.UserInfoRequest) (*grpc_service.UserInfoResponse, error) {
	userInfo := &userModel.UserInfo{
		Id:           int(request.Id),
		Sex:          int(request.Sex),
		UserName:     request.UserName,
		NoticeUserId: request.NoticeUserId,
		NickName:     request.NickName,
		Email:        request.Email,
		Phone:        request.Phone,
		Avatar:       request.Avatar,
	}
	err := c.userService.SaveUserInfo(ctx, userInfo)
	if err != nil {
		return nil, err
	}
	return &grpc_service.UserInfoResponse{}, nil
}

func NewConsoleInfoService() grpc_service.GetConsoleInfoServer {
	c := &consoleInfoService{}
	bean.Autowired(&c.namespaceService)
	bean.Autowired(&c.apiService)
	bean.Autowired(&c.modulePluginService)
	bean.Autowired(&c.navigationService)
	bean.Autowired(&c.dynamicService)
	bean.Autowired(&c.applicationService)
	bean.Autowired(&c.clusterService)
	bean.Autowired(&c.accessService)

	bean.Autowired(&c.modulesCache)
	bean.Autowired(&c.userService)
	return c
}

func (c *consoleInfoService) GetAllNamespaces(ctx context.Context, request *grpc_service.EmptyRequest) (*grpc_service.NamespacesListResp, error) {
	namespaces, err := c.namespaceService.GetAll()
	if err != nil {
		return nil, fmt.Errorf("获取命名空间列表报错. err: %s", err)
	}
	items := make([]*grpc_service.NamespaceItem, 0, len(namespaces))
	for _, ns := range namespaces {
		items = append(items, &grpc_service.NamespaceItem{
			NamespaceId:   int32(ns.Id),
			NamespaceName: ns.Name,
		})
	}
	return &grpc_service.NamespacesListResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetApis(ctx context.Context, req *grpc_service.GetApisReq) (*grpc_service.ApisResp, error) {
	apis, err := c.apiService.GetAPIInfoByPath(ctx, int(req.NamespaceId), req.Path)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取API列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ApisItem, 0, len(apis))
	for _, apiInfo := range apis {
		items = append(items, &grpc_service.ApisItem{
			Uuid:      apiInfo.UUID,
			Name:      apiInfo.Name,
			GroupUuid: apiInfo.GroupUUID,
			Path:      apiInfo.RequestPathLabel,
			Desc:      apiInfo.Desc,
			Methods:   apiInfo.Method,
		})
	}
	return &grpc_service.ApisResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetApisByUUIDs(ctx context.Context, req *grpc_service.GetApisByUUIDsReq) (*grpc_service.ApisResp, error) {
	apis, err := c.apiService.GetAPIInfoByUUIDS(ctx, int(req.NamespaceId), req.Uuids)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("通过uuids获取API列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ApisItem, 0, len(apis))
	for _, apiInfo := range apis {
		items = append(items, &grpc_service.ApisItem{
			Uuid:      apiInfo.UUID,
			Name:      apiInfo.Name,
			GroupUuid: apiInfo.GroupUUID,
			Path:      apiInfo.RequestPathLabel,
			Desc:      apiInfo.Desc,
			Methods:   apiInfo.Method,
		})
	}
	return &grpc_service.ApisResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetApisByServices(ctx context.Context, req *grpc_service.GetApisByServicesReq) (*grpc_service.ApisResp, error) {
	apis, err := c.apiService.GetAPIListByServiceName(ctx, int(req.NamespaceId), req.Services)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("通过服务列表获取API列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ApisItem, 0, len(apis))
	for _, apiInfo := range apis {
		items = append(items, &grpc_service.ApisItem{
			Uuid:      apiInfo.UUID,
			Name:      apiInfo.Name,
			GroupUuid: apiInfo.GroupUUID,
			Path:      apiInfo.RequestPathLabel,
			Desc:      apiInfo.Desc,
			Methods:   apiInfo.Method,
		})
	}
	return &grpc_service.ApisResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAllServices(ctx context.Context, req *grpc_service.GetServicesReq) (*grpc_service.ServicesResp, error) {
	basicInfos, err := c.dynamicService.ListByKeyword(ctx, int(req.NamespaceId), professionService, req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取Service列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ServicesItem, 0, len(basicInfos))
	for _, info := range basicInfos {
		items = append(items, &grpc_service.ServicesItem{
			Name:  info.ID,
			Title: info.Title,
			Desc:  info.Description,
		})
	}
	return &grpc_service.ServicesResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAllServicesByNames(ctx context.Context, req *grpc_service.GetServicesByNamesReq) (*grpc_service.ServicesResp, error) {
	basicInfos, err := c.dynamicService.ListByNames(ctx, int(req.NamespaceId), professionService, req.Names)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取Service列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ServicesItem, 0, len(basicInfos))
	for _, info := range basicInfos {
		items = append(items, &grpc_service.ServicesItem{
			Name:  info.ID,
			Title: info.Title,
			Desc:  info.Description,
		})
	}
	return &grpc_service.ServicesResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAllApps(ctx context.Context, req *grpc_service.GetAppsReq) (*grpc_service.AppsResp, error) {
	basicInfos, err := c.applicationService.AllApp(ctx, int(req.NamespaceId))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取应用列表报错. err: %s", err)
	}
	items := make([]*grpc_service.AppsItem, 0, len(basicInfos))
	for _, info := range basicInfos {
		items = append(items, &grpc_service.AppsItem{
			Uuid: info.Uuid,
			Name: info.Name,
			Desc: info.Desc,
		})
	}
	return &grpc_service.AppsResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAppsByUuids(ctx context.Context, req *grpc_service.GetAppsByUuidsReq) (*grpc_service.AppsResp, error) {
	basicInfos, err := c.applicationService.AppListByUUIDS(ctx, int(req.NamespaceId), req.Uuids)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取应用列表报错. err: %s", err)
	}
	items := make([]*grpc_service.AppsItem, 0, len(basicInfos))
	for _, info := range basicInfos {
		items = append(items, &grpc_service.AppsItem{
			Uuid: info.Uuid,
			Name: info.Name,
			Desc: info.Desc,
		})
	}
	return &grpc_service.AppsResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetNavigationModules(ctx context.Context, req *grpc_service.EmptyRequest) (*grpc_service.NavigationModulesResp, error) {
	modules := c.modulePluginService.GetEnable(ctx) //这个接口已从缓存中拿了

	navigations := c.navigationService.List()
	accessList := c.accessService.GetEnable(ctx)
	accessOfModule := common.SliceToMapArray(accessList, func(t *model.Access) string {
		return t.Module
	})

	navigationItems := make([]*grpc_service.NavigationItem, 0, len(navigations))
	modulesItems := make([]*grpc_service.ModuleItem, 0, len(modules))

	for _, item := range navigations {
		navigationItems = append(navigationItems, &grpc_service.NavigationItem{
			Id:    item.Uuid,
			Cname: item.Title,
		})
	}
	for _, item := range modules {

		mi := &grpc_service.ModuleItem{
			Name:         item.Name,
			Cname:        item.CName,
			NavigationId: item.Navigation,
			Access:       nil,
		}
		if aom, h := accessOfModule[item.Name]; h {
			asl := make([]*grpc_service.AccessItem, 0, len(aom))
			for _, it := range aom {
				asl = append(asl, &grpc_service.AccessItem{
					Name:   it.Name,
					Cname:  it.CName,
					Depend: it.Depend,
				})
			}
			mi.Access = asl
		}
		modulesItems = append(modulesItems, mi)
	}
	return &grpc_service.NavigationModulesResp{
		NavigationItems: navigationItems,
		ModulesItems:    modulesItems,
	}, nil
}

func (c *consoleInfoService) GetClusters(ctx context.Context, req *grpc_service.GetClustersReq) (*grpc_service.ClusterInfoResp, error) {
	clusters, err := c.clusterService.GetByNamespaceId(ctx, int(req.NamespaceId))
	if err != nil {
		return nil, errors.New("获取集群列表失败")
	}
	items := make([]*grpc_service.ClusterInfo, 0, len(clusters))
	for _, clu := range clusters {
		items = append(items, &grpc_service.ClusterInfo{
			Name:  clu.Name,
			Title: clu.Title,
			Uuid:  clu.UUID,
			Env:   clu.Env,
			Desc:  clu.Desc,
		})
	}
	return &grpc_service.ClusterInfoResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetClustersByNames(ctx context.Context, req *grpc_service.GetClustersReq) (*grpc_service.ClusterInfoResp, error) {
	clusters, err := c.clusterService.GetByNames(ctx, int(req.NamespaceId), req.Names)
	if err != nil {
		return nil, errors.New("获取集群列表失败")
	}
	items := make([]*grpc_service.ClusterInfo, 0, len(clusters))
	for _, clu := range clusters {
		items = append(items, &grpc_service.ClusterInfo{
			Name:  clu.Name,
			Title: clu.Title,
			Uuid:  clu.UUID,
			Env:   clu.Env,
			Desc:  clu.Desc,
		})
	}
	return &grpc_service.ClusterInfoResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetClustersByUUIDs(ctx context.Context, req *grpc_service.GetClustersReq) (*grpc_service.ClusterInfoResp, error) {
	clusters, err := c.clusterService.GetByUUIDs(ctx, int(req.NamespaceId), req.Uuids)
	if err != nil {
		return nil, errors.New("获取集群列表失败")
	}
	items := make([]*grpc_service.ClusterInfo, 0, len(clusters))
	for _, clu := range clusters {
		items = append(items, &grpc_service.ClusterInfo{
			Name:  clu.Name,
			Title: clu.Title,
			Uuid:  clu.UUID,
			Env:   clu.Env,
			Desc:  clu.Desc,
		})
	}
	return &grpc_service.ClusterInfoResp{
		Items: items,
	}, nil
}
