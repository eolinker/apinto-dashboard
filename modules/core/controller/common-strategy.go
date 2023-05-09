package controller

import (
	strategy_controller "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-controller"
	apinto_module "github.com/eolinker/apinto-module"
	"net/http"
)

func commonStrategy() apinto_module.RoutersInfo {
	commonStrategyController := strategy_controller.NewStrategyCommonController()
	return apinto_module.RoutersInfo{
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/filter-options",
			Handler:     "strategy-common.filterOptions",
			HandlerFunc: []apinto_module.HandlerFunc{commonStrategyController.FilterOptions},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/filter-remote/:name",
			Handler:     "strategy-common.filterRemote",
			HandlerFunc: []apinto_module.HandlerFunc{commonStrategyController.FilterRemote},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/metrics-options",
			Handler:     "strategy-common.metricsOptions",
			HandlerFunc: []apinto_module.HandlerFunc{commonStrategyController.MetricsOptions},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/content-type",
			Handler:     "strategy-common.contentType",
			HandlerFunc: []apinto_module.HandlerFunc{commonStrategyController.ContentType},
		},
		{
			Method:      http.MethodGet,
			Path:        "/api/strategy/charset",
			Handler:     "strategy-common.charset",
			HandlerFunc: []apinto_module.HandlerFunc{commonStrategyController.Charset},
		},
	}
}
