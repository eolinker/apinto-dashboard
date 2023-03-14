package strategy_service

import "github.com/eolinker/eosc/common/bean"

func init() {
	commonStrategy := newStrategyCommonService()

	bean.Injection(&commonStrategy)
}
