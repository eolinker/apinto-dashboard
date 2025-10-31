package service

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/eolinker/apinto-dashboard/modules/user"

	"gorm.io/gorm"

	api_dto "github.com/eolinker/apinto-dashboard/modules/api/api-dto"

	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	service "github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/eosc/common/bean"
)

var (
	_ grpc_service.ApiServiceServer = (*apiServiceServer)(nil)
)

type apiServiceServer struct {
	grpc_service.UnimplementedApiServiceServer

	apiService  service.IAPIService
	userService user.IUserInfoService
}

func NewApiServiceServer() grpc_service.ApiServiceServer {
	s := &apiServiceServer{}
	bean.Autowired(&s.apiService)
	bean.Autowired(&s.userService)
	return s
}
func (s *apiServiceServer) getUserId(ctx context.Context, userName string) int {
	if userName == "" {
		return 0
	}
	info, err := s.userService.GetUserInfoByName(ctx, userName)
	if err != nil {
		return 0
	}
	return info.Id
}
func (s *apiServiceServer) createFromTemplate(ctx context.Context, namespaceId int, req *grpc_service.ApiCreateByTemplateRequest) *grpc_service.ApiCreateResponse {
	templateInfo, err := s.apiService.GetAPIVersionInfo(ctx, namespaceId, req.TemplateApiUuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &grpc_service.ApiCreateResponse{
				Uuid:     req.Uuid,
				ErrorMsg: "template api not exist",
				Status:   grpc_service.CreateStatus_TemplateAPINotExist,
			}
		}
		return &grpc_service.ApiCreateResponse{
			Uuid:     req.Uuid,
			ErrorMsg: err.Error(),
			Status:   grpc_service.CreateStatus_Fail,
		}
	}
	od, err := s.apiService.GetAPIVersionInfo(ctx, namespaceId, req.Uuid)
	if err == nil && od != nil {
		if methodsEq(od.Version.Method, req.Methods) {
			return &grpc_service.ApiCreateResponse{
				Uuid:     req.Uuid,
				ErrorMsg: "api already exists",
				Status:   grpc_service.CreateStatus_APIExist,
			}
		}
		operatorId := s.getUserId(ctx, req.Operator)

		err = s.apiService.UpdateAPI(ctx, namespaceId, operatorId, &api_dto.APIInfo{
			ApiName:          od.Api.Name,
			UUID:             od.Api.UUID,
			GroupUUID:        od.Api.GroupUUID,
			Desc:             od.Api.Desc,
			IsDisable:        od.Api.IsDisable,
			Scheme:           od.Api.Scheme,
			RequestPath:      od.Api.RequestPath,
			RequestPathLabel: od.Api.RequestPathLabel,
			ServiceName:      od.Version.ServiceName,
			Method:           req.Methods,
			ProxyPath:        od.Version.ProxyPath,
			Hosts:            od.Version.Hosts,
			Timeout:          od.Version.Timeout,
			Retry:            od.Version.Retry,
			Match:            od.Version.Match,
			Header:           od.Version.Header,
			TemplateUUID:     od.Version.TemplateUUID,
		})

		if err != nil {
			return &grpc_service.ApiCreateResponse{
				Uuid:     req.Uuid,
				ErrorMsg: err.Error(),
				Status:   grpc_service.CreateStatus_Fail,
			}
		}
		return &grpc_service.ApiCreateResponse{
			Uuid:     req.Uuid,
			ErrorMsg: "",
			Status:   grpc_service.CreateStatus_SUCCESS,
		}
	}
	if err != gorm.ErrRecordNotFound {
		return &grpc_service.ApiCreateResponse{
			Uuid:     req.Uuid,
			ErrorMsg: err.Error(),
			Status:   grpc_service.CreateStatus_Fail,
		}
	}

	path := req.Path

	index := len(templateInfo.Api.RequestPath) - 1
	if templateInfo.Api.RequestPath[index] == '*' {
		path = strings.Replace(req.Path, templateInfo.Version.ProxyPath, templateInfo.Version.RequestPath[:index], 1)
	}
	operatorId := s.getUserId(ctx, req.Operator)
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("%s-未命名", templateInfo.Api.Name)
	}

	// 该api不存在，可以创建
	_, _, err = s.apiService.CreateAPI(ctx, 1, operatorId, &api_dto.APIInfo{
		ApiName:          name,
		UUID:             req.Uuid,
		GroupUUID:        templateInfo.Api.GroupUUID,
		Desc:             req.Desc,
		IsDisable:        templateInfo.Api.IsDisable,
		Scheme:           templateInfo.Api.Scheme,
		RequestPath:      path,
		RequestPathLabel: path,
		ServiceName:      templateInfo.Version.ServiceName,
		Method:           req.Methods,
		ProxyPath:        req.Path,
		Hosts:            templateInfo.Version.Hosts,
		Timeout:          templateInfo.Version.Timeout,
		Retry:            templateInfo.Version.Retry,
		Match:            templateInfo.Version.Match,
		Header:           templateInfo.Version.Header,
		TemplateUUID:     templateInfo.Version.TemplateUUID,
	})
	if err != nil {
		return &grpc_service.ApiCreateResponse{
			Uuid:     req.Uuid,
			ErrorMsg: err.Error(),
			Status:   grpc_service.CreateStatus_Fail,
		}
	}
	return &grpc_service.ApiCreateResponse{
		Uuid:     req.Uuid,
		ErrorMsg: "",
		Status:   grpc_service.CreateStatus_SUCCESS,
	}
}
func methodsEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Strings(a)
	sort.Strings(b)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func (s *apiServiceServer) CreateFromTemplate(server grpc_service.ApiService_CreateFromTemplateServer) error {

	for {
		req, err := server.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		resp := s.createFromTemplate(server.Context(), 1, req)
		err = server.Send(resp)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}

	return nil
}

func (s *apiServiceServer) ForeachInfo(req *grpc_service.ApiSearchRequest, server grpc_service.ApiService_ForeachInfoServer) error {

	count := 0
	page := 0
	pageSize := 20
	for {
		page++
		options, total, err := s.apiService.GetAPIRemoteOptions(server.Context(), 1, page, pageSize, req.Keyword, req.Group)
		if err != nil {
			return err
		}
		for _, i := range options {
			err := server.Send(&grpc_service.ApiItem{
				Uuid:  i.Uuid,
				Name:  i.Title,
				Path:  i.RequestPath,
				Group: i.Group,
			})
			if err != nil {
				return err
			}
		}
		count += pageSize
		if count >= total {
			break
		}
	}
	return nil

}
