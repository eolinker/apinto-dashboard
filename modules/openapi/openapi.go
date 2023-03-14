package openapi

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/api/model"
	"github.com/eolinker/apinto-dashboard/modules/openapi/openapi-dto"
	"github.com/eolinker/apinto-dashboard/modules/openapi/openapi-model"
)

type IAPIOpenAPIService interface {
	SyncImport(ctx context.Context, namespaceID, appID int, data *openapi_dto.SyncImportData) ([]*apimodel.ImportAPIListItem, error)
	GetSyncImportInfo(ctx context.Context, namespaceID int) ([]*openapi_model.ApiOpenAPIGroups, []*openapi_model.ApiOpenAPIService, []string, error)
}
