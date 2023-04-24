package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/grpc-service"
)

var _ grpc_service.GetConsoleInfoServer = (*consoleInfoService)(nil)

type consoleInfoService struct {
}

func NewConsoleInfoService() grpc_service.GetConsoleInfoServer {
	c := &consoleInfoService{}
	return c
}

func (c *consoleInfoService) GetAllNamespaces(ctx context.Context, request *grpc_service.EmptyRequest) (*grpc_service.NamespacesListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c *consoleInfoService) GetApis(ctx context.Context, req *grpc_service.GetApisReq) (*grpc_service.ApisResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c *consoleInfoService) GetApisByUUIDs(ctx context.Context, req *grpc_service.GetApisByUUIDsReq) (*grpc_service.ApisResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c *consoleInfoService) GetAllServices(ctx context.Context, req *grpc_service.GetServicesReq) (*grpc_service.ServicesResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c *consoleInfoService) GetAllServicesByNames(ctx context.Context, req *grpc_service.GetServicesByNamesReq) (*grpc_service.ServicesResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c *consoleInfoService) GetAllApps(ctx context.Context, req *grpc_service.GetAppsReq) (*grpc_service.AppsResp, error) {
	//TODO implement me
	panic("implement me")
}

func (c *consoleInfoService) GetAppsByUuids(ctx context.Context, req *grpc_service.GetAppsByUuidsReq) (*grpc_service.AppsResp, error) {
	//TODO implement me
	panic("implement me")
}
