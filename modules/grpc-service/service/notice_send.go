package service

import (
	"context"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
)

var _ grpc_service.NoticeSendServer = (*noticeSendService)(nil)

type noticeSendService struct {
}

func NewNoticeSendService() grpc_service.NoticeSendServer {
	n := &noticeSendService{}
	return n
}

func (n *noticeSendService) Send(ctx context.Context, req *grpc_service.NoticeSendReq) (*grpc_service.NoticeSendResp, error) {
	//TODO implement me
	panic("implement me")
}
