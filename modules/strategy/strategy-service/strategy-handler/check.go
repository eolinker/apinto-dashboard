package strategy_handler

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
)

func checkFilters(fileters []*strategy_dto.FilterInput) error {
	filterNameSet := make(map[string]struct{})
	for _, filter := range fileters {
		switch filter.Name {
		case config.FilterApplication, config.FilterApi, config.FilterPath, config.FilterService, config.FilterMethod, config.FilterIP:
		default:
			if !common.IsMatchFilterAppKey(filter.Name) {
				return fmt.Errorf("filter.Name %s is illegal. ", filter.Name)
			}
		}

		if len(filter.Values) == 0 {
			return fmt.Errorf("filter.Options can't be null. filter.Name:%s ", filter.Name)
		}

		if _, has := filterNameSet[filter.Name]; has {
			return fmt.Errorf("filterName %s is reduplicative. ", filter.Name)
		}
		filterNameSet[filter.Name] = struct{}{}
	}

	return nil
}

// checkStatusCode 校验状态码
func checkStatusCode(codes ...int) error {
	codeSet := map[int]struct{}{}
	for _, code := range codes {
		if _, has := codeSet[code]; has {
			return fmt.Errorf("status_code %d is redeplicated. ", code)
		}
		if code < 0 || code >= 1000 {
			return fmt.Errorf("status_code %d is illegal. ", code)
		}
		codeSet[code] = struct{}{}
	}
	return nil
}

func checkCharset(charset string) error {
	switch charset {
	case config.CharsetUTF8, config.CharsetGBK, config.CharsetASCII:
	default:
		return fmt.Errorf("charset %s is illegal. ", charset)
	}
	return nil
}
