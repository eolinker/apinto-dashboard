package service

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

var _ grpc_service.GetConsoleInfoServer = (*consoleInfoService)(nil)

type consoleInfoService struct {
	namespaceService namespace.INamespaceService
	apiService       api.IAPIService
	upstreamService  upstream.IService
	appService       application.IApplicationService
}

func NewConsoleInfoService() grpc_service.GetConsoleInfoServer {
	c := &consoleInfoService{}
	bean.Autowired(&c.namespaceService)
	bean.Autowired(&c.apiService)
	bean.Autowired(&c.upstreamService)
	bean.Autowired(&c.appService)
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

func (c *consoleInfoService) GetAllServices(ctx context.Context, req *grpc_service.GetServicesReq) (*grpc_service.ServicesResp, error) {
	services, err := c.upstreamService.GetServiceListAll(ctx, int(req.NamespaceId), req.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取Service列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ServicesItem, 0, len(services))
	for _, info := range services {
		items = append(items, &grpc_service.ServicesItem{
			Name: info.Name,
			Desc: info.Desc,
		})
	}
	return &grpc_service.ServicesResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAllServicesByNames(ctx context.Context, req *grpc_service.GetServicesByNamesReq) (*grpc_service.ServicesResp, error) {
	services, err := c.upstreamService.GetServiceListByNames(ctx, int(req.NamespaceId), req.Names)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取Service列表报错. err: %s", err)
	}
	items := make([]*grpc_service.ServicesItem, 0, len(services))
	for _, info := range services {
		items = append(items, &grpc_service.ServicesItem{
			Name: info.Name,
			Desc: info.Desc,
		})
	}
	return &grpc_service.ServicesResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAllApps(ctx context.Context, req *grpc_service.GetAppsReq) (*grpc_service.AppsResp, error) {
	apps, err := c.appService.AppListAll(ctx, int(req.NamespaceId))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取应用列表报错. err: %s", err)
	}
	items := make([]*grpc_service.AppsItem, 0, len(apps))
	for _, info := range apps {
		items = append(items, &grpc_service.AppsItem{
			Uuid: info.IdStr,
			Name: info.Name,
			Desc: info.Desc,
		})
	}
	return &grpc_service.AppsResp{
		Items: items,
	}, nil
}

func (c *consoleInfoService) GetAppsByUuids(ctx context.Context, req *grpc_service.GetAppsByUuidsReq) (*grpc_service.AppsResp, error) {
	apps, err := c.appService.AppListByUUIDS(ctx, int(req.NamespaceId), req.Uuids)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("获取应用列表报错. err: %s", err)
	}
	items := make([]*grpc_service.AppsItem, 0, len(apps))
	for _, info := range apps {
		items = append(items, &grpc_service.AppsItem{
			Uuid: info.IdStr,
			Name: info.Name,
			Desc: info.Desc,
		})
	}
	return &grpc_service.AppsResp{
		Items: items,
	}, nil
}
