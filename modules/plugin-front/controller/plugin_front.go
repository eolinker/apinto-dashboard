package controller

import (
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
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
	router.GET("/plugin/icon/*filepath", p.getPluginICON)
	router.GET("/plugin/md/*filepath", p.getPluginMD)
	router.GET("/plugin/info/:id/resources/*filepath", p.getPluginResources)
}

// getPluginICON 获取插件描述中要用到的图标
func (p *pluginFrontController) getPluginICON(c *gin.Context) {
	filePath := strings.Trim(c.Param("filepath"), "/")
	list := strings.Split(filePath, "/")
	pathLen := len(list)
	//若文件路径层级多于 id/file 或 filepath 为空
	if pathLen > 2 || filePath == "" {
		log.Info(fmt.Errorf("%s:%w", c.Request.URL.Path, filePathErr))
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	//防止读取resources目录
	if pathLen == 2 && list[1] == "resources" {
		log.Info(fmt.Errorf("%s:%w", c.Request.URL.Path, filePathErr))
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	if pathLen == 1 {
		//TODO 从插件配置文件获取图标名
		fileName := ""
		filePath = fmt.Sprintf("%s/%s", list[0], fileName)
	}

	pluginFs := gin.Dir(fileDir, false)
	fileServer := http.StripPrefix("/plugin/icon", http.FileServer(pluginFs))
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
func (p *pluginFrontController) getPluginMD(c *gin.Context) {
	filePath := strings.Trim(c.Param("filepath"), "/")
	list := strings.Split(filePath, "/")
	pathLen := len(list)
	//若文件路径层级多于 id/file 或 filepath 为空
	if pathLen > 2 || filePath == "" {
		log.Info(fmt.Errorf("%s:%w", c.Request.URL.Path, filePathErr))
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	//防止读取resources目录
	if pathLen == 2 && list[1] == "resources" {
		log.Info(fmt.Errorf("%s:%w", c.Request.URL.Path, filePathErr))
		c.Data(http.StatusNotFound, "application/text", []byte("404 page not found"))
		return
	}
	if pathLen == 1 {
		filePath = fmt.Sprintf("%s/%s", list[0], "README.md")
	}

	pluginFs := gin.Dir(fileDir, false)
	fileServer := http.StripPrefix("/plugin/md", http.FileServer(pluginFs))
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
	fileServer := http.StripPrefix("/plugin/md", http.FileServer(pluginFs))
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
