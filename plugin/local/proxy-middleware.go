package local

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	plugin_client "github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin-client"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"sync"
)

const (
	MiddlewarePre  = "pre"
	MiddlewarePost = "post"
)

func (p *ProxyAPi) CreateMiddleware(name string, rule [][]string) apinto_module.MiddlewareHandler {
	once := sync.Once{}
	var middlewareHandler *plugin_client.MiddlewareClientHandler
	var err error

	return apinto_module.MiddlewareHandler{
		Name: name,
		Rule: apinto_module.CreateMiddlewareRules(rule),
		Handler: func(ginCtx *gin.Context) {
			p.initClient()
			if p.clientError != nil {
				log.Error("middleware proxy client:", p.clientError)
				return
			}

			once.Do(func() {
				p.wg.Add(1)
				defer p.wg.Done()
				middlewareHandler, err = p.handler.CreateMiddleware(name, nil, nil)
			})
			if err != nil {
				log.Error("middleware proxy handler:", p.clientError)

				return
			}
			p.wg.Add(1)
			defer p.wg.Done()
			middlewareHandler.Middleware(ginCtx)
		},
	}

}
