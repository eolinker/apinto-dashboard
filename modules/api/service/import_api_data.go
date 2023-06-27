package api_service

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
	"time"
)

type IImportApiCache interface {
	cache.IRedisCache[apimodel.ImportAPIRedisData, string]
}

func importKey(token string) string {
	return fmt.Sprintf("import_api_token:%s", token)
}

func newImportCache() IImportApiCache {
	cacheInfo := cache.CreateRedisCache[apimodel.ImportAPIRedisData, string](time.Hour*8, importKey)

	return cacheInfo

}
