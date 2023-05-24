package application_controller

import (
	"context"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/modules/application"
	"github.com/eolinker/eosc/common/bean"
)

var _ apinto_module.IFilterOptionHandler = (*FilterOption)(nil)

type FilterOption struct {
	appService application.IApplicationService
	config     apinto_module.FilterOptionConfig
}

func newFilterOption() *FilterOption {
	f := &FilterOption{config: apinto_module.FilterOptionConfig{
		Title: "应用",
		Titles: []apinto_module.OptionTitle{
			{
				Field: "title",
				Title: "应用名称",
			},
			{
				Field: "uuid",
				Title: "应用ID",
			},
			{
				Field: "desc",
				Title: "应用描述",
			},
		},
		Key: "uuid",
	}}
	bean.Autowired(&f.appService)
	return f
}

func (f *FilterOption) Name() string {
	return "application"
}

func (f *FilterOption) Config() apinto_module.FilterOptionConfig {
	return f.config
}

func (f *FilterOption) GetOptions(namespaceId int, keyword, groupUUID string, pageNum, pageSize int) ([]any, int) {
	list, err := f.appService.GetAppRemoteOptions(context.Background(), namespaceId, pageNum, pageSize, keyword)
	if err != nil {
		return nil, 0
	}
	return list, len(list)
}

func (f *FilterOption) Labels(namespaceId int, values ...string) []string {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return []string{f.Label(namespaceId, values[0])}
	}
	infos, err := f.appService.AppListByUUIDS(context.Background(), namespaceId, values)
	if err != nil {
		return nil
	}

	rs := make([]string, 0, len(infos))
	for _, v := range infos {
		rs = append(rs, v.Name)
	}
	return rs
}

func (f *FilterOption) Label(namespaceId int, value string) string {
	info, err := f.appService.AppBasicInfo(context.Background(), namespaceId, value)
	if err != nil {
		return ""
	}
	return info.Name
}
