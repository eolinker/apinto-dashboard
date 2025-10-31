package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/dto"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Plugin struct {
	pluginService mpm3.IPluginService
}

func NewPlugin() *Plugin {
	p := &Plugin{}
	bean.Autowired(&p.pluginService)
	return p
}

func (p *Plugin) Apis() []pm3.Api {
	return []pm3.Api{
		{
			Authority: 0,
			Method:    http.MethodPost,
			Path:      "/api/system/plugin/disable",

			HandlerFunc: p.disable,
		},
		{
			Authority: 0,
			Method:    http.MethodPost,
			Path:      "/api/system/plugin/enable",

			HandlerFunc: p.enable,
		},
	}
}

func (p *Plugin) enable(ginCtx *gin.Context) {
	pluginUUID := ginCtx.Query("id")

	input := new(dto.PluginEnableInfo)
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	cfg := &model.PluginEnableCfg{
		Header: common.SliceToSlice(input.Header, func(s dto.ExtendParams) *model.ExtendParams {
			return &model.ExtendParams{
				Name:  s.Name,
				Value: s.Value,
			}
		}),
		Query: common.SliceToSlice(input.Query, func(s dto.ExtendParams) *model.ExtendParams {
			return &model.ExtendParams{
				Name:  s.Name,
				Value: s.Value,
			}
		}),
		Initialize: common.SliceToSlice(input.Initialize, func(s dto.ExtendParams) *model.ExtendParams {
			return &model.ExtendParams{
				Name:  s.Name,
				Value: s.Value,
			}
		}),
	}

	err := p.pluginService.EnablePlugin(ginCtx, users.GetUserId(ginCtx), pluginUUID, cfg)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (p *Plugin) disable(ginCtx *gin.Context) {

	pluginUUID := ginCtx.Query("id")

	err := p.pluginService.DisablePlugin(ginCtx, users.GetUserId(ginCtx), pluginUUID)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, fmt.Sprintf("Disable plugin fail. err:%s", err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
