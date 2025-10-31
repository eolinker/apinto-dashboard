/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package module

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/controller"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/log/dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	singleController *Controller
	once             sync.Once
)

type Controller struct {
	clusterService     cluster.IClusterService
	clusterNodeService cluster.IClusterNodeService
}

func (c *Controller) RoutersInfo() apinto_module.RoutersInfo {
	return apinto_module.RoutersInfo{
		{
			Method: http.MethodGet,
			Path:   "/api/log/files",

			HandlerFunc: c.list,
		}, {
			Method: http.MethodGet,
			Path:   "/api/log/download/:key",

			HandlerFunc: c.download,
		}, {
			Method: http.MethodGet,
			Path:   "/api/log/tail/:key",

			HandlerFunc: c.tail,
		},
	}
}
func (c *Controller) list(ginCtx *gin.Context) {
	clusterName := ginCtx.Query("cluster")
	nodes, err := c.clusterNodeService.List(ginCtx, namespace_controller.GetNamespaceId(ginCtx), clusterName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	nodeName := ginCtx.Query("node")
	var addrFor []string
	addrCount := make(map[string]int) // 对admin url 计数，用来过滤被多个节点使用的admin url
	for _, n := range nodes {
		if n.Name == nodeName {
			addrFor = n.AdminAddrs
		}
		for _, ad := range n.AdminAddrs {
			addrCount[ad]++
		}
	}
	if len(addrFor) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "无效节点")

		return
	}
	adminUrls := make([]string, 0, len(addrFor))
	for _, ad := range addrFor {
		if addrCount[ad] == 1 {
			adminUrls = append(adminUrls, ad)
		}
	}
	if len(adminUrls) == 0 {
		controller.ErrorJson(ginCtx, http.StatusOK, "节点无独占admin url，无法执行该操作")
		return
	}

	client, err := v1.NewClient(adminUrls)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ol, err := client.ForOutput().List()
	if err != nil {
		return
	}

	fileOL := make([]string, 0, len(ol))
	for _, o := range ol {
		if o.Driver == "file" {
			fileOL = append(fileOL, o.Name)
		}
	}

	output := make([]*dto.Output, 0, len(fileOL)+1)

	output = append(output, readFiles("节点日志", adminUrls, fmt.Sprintf("/apinto/log/node")))

	for _, o := range fileOL {
		output = append(output, readFiles(o, adminUrls, fmt.Sprintf("/apinto/log/access/%s", o)))
	}

	data := map[string]any{
		"output": output,
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
func NewController() *Controller {
	once.Do(func() {
		singleController = new(Controller)
		bean.Autowired(&singleController.clusterService)
		bean.Autowired(&singleController.clusterNodeService)
	})
	return singleController
}

type FileItem struct {
	Name    string `json:"name,omitempty"`
	Size    string `json:"size,omitempty"`
	ModTime string `json:"modTime,omitempty"`
}

func readFiles(name string, addrs []string, prefix string) *dto.Output {

	out := &dto.Output{
		Name: name,
	}
	for _, addr := range addrs {
		url := fmt.Sprintf("%s/%s/files", strings.Trim(addr, "/"), strings.Trim(prefix, "/"))
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Info("read file form ", url, " ", err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Info("read file form ", url, " ", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil
		}

		out.TailKey = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s/%s/tail", strings.Trim(addr, "/"), strings.Trim(prefix, "/"))))
		body, err := io.ReadAll(resp.Body)
		if err != nil && err != io.EOF {
			log.Info("read file form ", url, " ", err)
			continue
		}
		_ = resp.Body.Close()
		files := make([]*FileItem, 0)
		err = json.Unmarshal(body, &files)
		if err != nil {
			log.Info("read file form ", url, " ", err)
			break
		}
		out.Files = make([]dto.File, 0, len(files))
		for _, fi := range files {
			key := fmt.Sprintf("%s/%s/file/%s", strings.Trim(addr, "/"), strings.Trim(prefix, "/"), fi.Name)
			out.Files = append(out.Files, dto.File{
				File: fi.Name,
				Size: fi.Size,
				Mod:  fi.ModTime,
				Key:  base64.URLEncoding.EncodeToString([]byte(key)),
			})
		}
		break
	}
	return out
}
