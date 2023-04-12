package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/initialize"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
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
		return
	}
	//若插件存在
	err = p.modulePluginService.CheckPluginISDeCompress(c, pluginID)
	if err != nil {
		log.Errorf("Decompress Plugin Package fail. pluginID:%s, err:%s", pluginID, err)
		c.Abort()
		return
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
		if fileName == "" {
			fileName = "icon.png"
		}
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

	//判断插件存不存在
	info, err := p.modulePluginService.GetPluginInfo(c, pluginID)
	if err != nil {
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		c.Abort()
		return
	}
	//若为内置插件，则从内嵌目录中获取
	if info.Type == 0 || info.Type == 1 {
		fsHandler, err := initialize.GetInnerPluginFSHandler(stripPrefix, filePath)
		//若文件不存在
		if err != nil {
			c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
			return
		}
		fsHandler.ServeHTTP(c.Writer, c.Request)
		return
	}

	pluginFs := gin.Dir(fileDir, false)
	fileServer := http.StripPrefix(stripPrefix, http.FileServer(pluginFs))
	// Check if file exists and/or if we have permission to access it
	f, err := pluginFs.Open(filePath)
	if err != nil {
		//TODO 文件不存在时
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	defer f.Close()

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
