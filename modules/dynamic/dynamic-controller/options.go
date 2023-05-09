package dynamic_controller

import (
	"context"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	dynamic_model "github.com/eolinker/apinto-dashboard/modules/dynamic/dynamic-model"
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
)

var _ apinto_module.IFilterOptionHandler = (*FilterOption)(nil)

type FilterOption struct {
	dynamicService dynamic.IDynamicService
	name           string
	config         apinto_module.FilterOptionConfig
	columns        []string
	profession     string
	drivers        []string
}

func NewFilterOption(name string, config apinto_module.FilterOptionConfig, columns []string, profession string, drivers []string) *FilterOption {
	f := &FilterOption{name: name, config: config, columns: columns, profession: profession, drivers: drivers}
	bean.Autowired(&f.dynamicService)
	return f
}

func (f *FilterOption) Name() string {
	return f.name
}

func (f *FilterOption) Config() apinto_module.FilterOptionConfig {
	return f.config
}

func (f *FilterOption) GetOptions(namespaceId int, keyword, groupUUID string, pageNum, pageSize int) ([]any, int) {
	list, i, err := f.dynamicService.List(context.Background(), namespaceId, f.profession, f.columns, f.drivers, keyword, pageNum, pageSize)
	if err != nil {
		return nil, 0
	}

	return common.SliceToSlice(list, func(t map[string]string) any {
		return t
	}), i
}

func (f *FilterOption) Labels(namespaceId int, values ...string) []string {
	if len(values) == 0 {
		return nil
	}
	if len(values) == 1 {
		return []string{f.Label(namespaceId, values[0])}
	}
	infos, err := f.dynamicService.ListByNames(context.Background(), namespaceId, f.profession, values)
	if err != nil {
		return nil
	}
	infoMap := common.SliceToMap(infos, func(t *dynamic_model.DynamicBasicInfo) string {
		return t.ID
	})
	rs := make([]string, 0, len(values))
	for _, v := range values {
		if i, has := infoMap[v]; has {
			rs = append(rs, i.Title)
		} else {
			rs = append(rs, "")
		}
	}
	return rs
}

func (f *FilterOption) Label(namespaceId int, value string) string {
	info, err := f.dynamicService.Info(context.Background(), namespaceId, f.profession, value)
	if err != nil {
		return ""
	}
	return info.BasicInfo.Title
}
