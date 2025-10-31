package strategy_handler

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/eosc/common/bean"
)

var (
	strategyCommonService strategy.IStrategyCommonService
)

func init() {
	bean.Autowired(&strategyCommonService)
}
func checkFilters(fileters []*strategy_dto.FilterInput, dsec string) error {
	filterNameSet := make(map[string]struct{})
	for _, filter := range fileters {
		if !strategyCommonService.FilterNameCheck(filter.Name) {
			return fmt.Errorf("%s.Name %s is illegal. ", dsec, filter.Name)
		}

		if len(filter.Values) == 0 {
			return fmt.Errorf("%s.Options can't be null. filter.Name:%s ", dsec, filter.Name)
		}

		if _, has := filterNameSet[filter.Name]; has {
			return fmt.Errorf("%s.Name %s is reduplicative. ", dsec, filter.Name)
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
