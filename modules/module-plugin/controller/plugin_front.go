package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const fileDir = "./plugin"

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
	}
}

func (p *pluginFrontController) setIConName(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		//获取插件配置的icon
		pluginID := c.Param("id")
		info, err := p.modulePluginService.GetPluginInfo(c, pluginID)
		if err != nil {
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			c.Abort()
			return
		}
		fileName = info.ICon
	}
	c.Set("file", fileName)
	c.Set("strip_prefix", "/plugin/icon")
}

func (p *pluginFrontController) setMDName(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		fileName = "README.md"
	}
	c.Set("file", fileName)
	c.Set("strip_prefix", "/plugin/md")
}

func (p *pluginFrontController) getPluginInfo(c *gin.Context) {
	pluginID := c.Param("id")
	fileName := c.GetString("file")
	stripPrefix := c.GetString("strip_prefix")

	filePath := fmt.Sprintf("%s/%s", pluginID, fileName)

	pluginFs := gin.Dir(fileDir, false)
	fileServer := http.StripPrefix(stripPrefix, http.FileServer(pluginFs))
	// Check if file exists and/or if we have permission to access it
	f, err := pluginFs.Open(filePath)
	if err != nil {
		//TODO 文件不存在时
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	f.Close()

	fileServer.ServeHTTP(c.Writer, c.Request)
}

// getPluginMD 获取插件描述中要用到的MD文件
func (p *pluginFrontController) getPluginResources(c *gin.Context) {
	pluginID := c.Param("id")

	filePath := fmt.Sprintf("%s/resources/%s", pluginID, strings.Trim(c.Param("filepath"), "/"))

	pluginFs := gin.Dir(fileDir, false)
	fileServer := http.StripPrefix("/plugin/info", http.FileServer(pluginFs))
	// Check if file exists and/or if we have permission to access it
	f, err := pluginFs.Open(filePath)
	if err != nil {
		//TODO 文件不存在时
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	f.Close()

	fileServer.ServeHTTP(c.Writer, c.Request)
}
