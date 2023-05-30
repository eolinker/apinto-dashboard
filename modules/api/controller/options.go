package controller

import (
	"context"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/eosc/common/bean"
)

var _ apinto_module.IFilterOptionHandler = (*FilterOption)(nil)

type FilterOption struct {
	apiService api.IAPIService
	config     apinto_module.FilterOptionConfig
}

func newFilterOption() *FilterOption {
	f := &FilterOption{config: apinto_module.FilterOptionConfig{
		Title: "API",
		Titles: []apinto_module.OptionTitle{
			{
				Field: "title",
				Title: "API名称",
			},
			{
				Field: "group",
				Title: "所属目录",
			},
			{
				Field: "request_path",
				Title: "请求路径",
			},
		},
		Key: "uuid",
	}}
	bean.Autowired(&f.apiService)
	return f
}

func (f *FilterOption) Name() string {
	return "api"
}

func (f *FilterOption) Config() apinto_module.FilterOptionConfig {
	return f.config
}

func (f *FilterOption) GetOptions(namespaceId int, keyword, groupUUID string, pageNum, pageSize int) ([]any, int) {
	list, i, err := f.apiService.GetAPIRemoteOptions(context.Background(), namespaceId, pageNum, pageSize, keyword, groupUUID)
	if err != nil {
		return nil, 0
	}
	return list, i
}

func (f *FilterOption) Labels(namespaceId int, values ...string) []string {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return []string{f.Label(namespaceId, values[0])}
	}
	infos, err := f.apiService.GetAPIRemoteByUUIDS(context.Background(), namespaceId, values)
	if err != nil {
		return nil
	}

	rs := make([]string, 0, len(infos))
	for _, v := range infos {
		rs = append(rs, v.Title)
	}
	return rs
}

func (f *FilterOption) Label(namespaceId int, value string) string {
	info, err := f.apiService.GetAPIInfo(context.Background(), namespaceId, value)
	if err != nil || info == nil {
		return ""
	}
	return info.Name
}
