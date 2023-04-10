package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const fileDir = "./plugin"

var (
	filePathErr = errors.New("filePath is illegal. ")
)

type pluginFrontController struct {
}

func RegisterPluginFrontRouter(router gin.IRoutes) {
	p := &pluginFrontController{}
	router.GET("/plugin/icon/:id/:file", p.checkPluginID, p.setIConName, p.getPluginInfo)
	router.GET("/plugin/icon/:id", p.checkPluginID, p.setIConName, p.getPluginInfo)

	router.GET("/plugin/md/:id/:file", p.checkPluginID, p.setMDName, p.getPluginInfo)
	router.GET("/plugin/md/:id", p.checkPluginID, p.setMDName, p.getPluginInfo)

	router.GET("/plugin/info/:id/resources/*filepath", p.checkPluginID, p.getPluginResources)
}

func (p *pluginFrontController) checkPluginID(c *gin.Context) {
	pluginID := c.Param("id")
	//TODO 若不存在
	c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
	c.Abort()
}

func (p *pluginFrontController) setIConName(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		//TODO 获取插件配置文件内的icon
		fileName = ""
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
