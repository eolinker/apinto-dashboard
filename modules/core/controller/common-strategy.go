package controller

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	strategy_controller "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-controller"
	"github.com/eolinker/apinto-dashboard/pm3"
	"net/http"
)

func commonStrategy() apinto_module.RoutersInfo {
	commonStrategyController := strategy_controller.NewStrategyCommonController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/filter-options",
			HandlerFunc: commonStrategyController.FilterOptions,
			Authority:   pm3.Public,
		},
		{
			Method:    http.MethodGet,
			Path:      "/api/strategy/filter-remote/:name",
			Authority: pm3.Public,

			HandlerFunc: commonStrategyController.FilterRemote,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/metrics-options",
			HandlerFunc: commonStrategyController.MetricsOptions,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/content-type",
			HandlerFunc: commonStrategyController.ContentType,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/charset",
			HandlerFunc: commonStrategyController.Charset,
			Authority:   pm3.Public,
		},
	}
}
