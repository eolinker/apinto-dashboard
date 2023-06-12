package controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/resources_manager"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type pluginFrontController struct {
	modulePluginService module_plugin.IModulePluginService
}

func newPluginFrontController() *pluginFrontController {
	p := &pluginFrontController{}
	bean.Autowired(&p.modulePluginService)

	return p
}

func (p *pluginFrontController) checkPluginID(c *gin.Context) {
	pluginID := c.Param("id")

	isExist, err := p.modulePluginService.CheckPluginInstalled(c, pluginID)
	if err != nil {
		controller.ErrorJson(c, http.StatusOK, err.Error())
		c.Abort()
	}

	if !isExist {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		c.Abort()
		return
	}
	//若插件存在
	err = p.modulePluginService.CheckExternPluginInCache(c, pluginID)
	if err != nil {
		log.Errorf("CheckPluginInCache fail. pluginID:%s, err:%s", pluginID, err)
		c.Abort()
		return
	}
}

func (p *pluginFrontController) getICon(c *gin.Context) {
	pluginID := c.Param("id")
	pluginResources, has := resources_manager.GetEmbedPluginResources(pluginID)
	if !has {
		pluginResources, has = resources_manager.GetExternPluginResources(pluginID)
		if !has {
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			return
		}
	}

	icon, has := pluginResources.ICon()
	if !has {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	contentType := http.DetectContentType(icon)
	c.Data(http.StatusOK, contentType, icon)
}

func (p *pluginFrontController) getMD(c *gin.Context) {
	pluginID := c.Param("id")
	pluginResources, has := resources_manager.GetEmbedPluginResources(pluginID)
	if !has {
		pluginResources, has = resources_manager.GetExternPluginResources(pluginID)
		if !has {
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			return
		}
	}
	fileName := c.Param("file")
	var md []byte
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
func (p *pluginFrontController) getPluginResources(c *gin.Context) {
	pluginID := c.Param("id")
	pluginResources, has := resources_manager.GetEmbedPluginResources(pluginID)
	if !has {
		pluginResources, has = resources_manager.GetExternPluginResources(pluginID)
		if !has {
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			return
		}
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
