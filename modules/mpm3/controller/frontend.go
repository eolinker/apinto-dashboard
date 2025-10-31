package controller

import (
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FrontendController struct {
	resourcesService mpm3.IResourcesService
	pluginService    mpm3.IPluginService
}

func NewFrontendController() *FrontendController {
	c := &FrontendController{}
	bean.Autowired(&c.pluginService)
	bean.Autowired(&c.resourcesService)
	return c
}

func (f *FrontendController) apis() []pm3.Api {
	return []pm3.Api{
		{
			Method:      http.MethodGet,
			Path:        "/plugin/icon/:id/:file",
			HandlerFunc: f.getICon,
		},
		{
			Method: http.MethodGet,
			Path:   "/plugin/icon/:id",

			HandlerFunc: f.getICon,
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/md/:id/:file",
			HandlerFunc: f.getMD,
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/md/:id",
			HandlerFunc: f.getMD,
		},
		{
			Method:      http.MethodGet,
			Path:        "/plugin/info/:id/resources/*filepath",
			HandlerFunc: f.getPluginResources,
		},
	}
}

func (f *FrontendController) getICon(c *gin.Context) {
	pluginID := c.Param("id")

	pluginResources, err := f.resourcesService.Get(c, pluginID)
	if err != nil {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}

	icon, has := pluginResources.ICon()
	if !has {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	contentType := http.DetectContentType(icon)
	c.Data(http.StatusOK, contentType, icon)
}

func (f *FrontendController) getMD(c *gin.Context) {
	pluginID := c.Param("id")

	pluginResources, err := f.resourcesService.Get(c, pluginID)
	if err != nil || pluginResources == nil {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}

	fileName := c.Param("file")
	var md []byte
	var has bool
	if fileName != "" {
		md, has = pluginResources.ReadMe(fileName)
	} else {
		md, has = pluginResources.RM()
	}
	if !has {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}

	contentType := http.DetectContentType(md)
	c.Data(http.StatusOK, contentType, md)
}

// getPluginMD 获取插件描述中要用到的MD文件
func (p *FrontendController) getPluginResources(c *gin.Context) {
	pluginID := c.Param("id")

	pluginResources, err := p.resourcesService.Get(c, pluginID)
	if err != nil || pluginResources == nil {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	filePath := c.Param("filepath")
	data, has := pluginResources.Resources(filePath)
	if !has {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	contentType := http.DetectContentType(data)
	c.Data(http.StatusOK, contentType, data)

}
