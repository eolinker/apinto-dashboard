package service

import (
	"fmt"

	"github.com/eolinker/apinto-dashboard/modules/middleware/model"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	middlewareService := newMiddlewareService()

	bean.Injection(&middlewareService)
	for i := 0; i < 30; i++ {
		demoMiddlewares = append(demoMiddlewares, &model.MiddlewareInfo{
			Name: fmt.Sprintf("apinto-%d.login", i+1),
			Desc: fmt.Sprintf("apinto-%d 用户登录", i+1),
		})
	}
}

var demoMiddlewares = make([]*model.MiddlewareInfo, 0, 30)
