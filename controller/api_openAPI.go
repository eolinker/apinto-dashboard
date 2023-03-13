package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type apiOpenAPIController struct {
	apiOpenAPIService service.IAPIOpenAPIService
	extAPPService     service.IExternalApplicationService
}

func RegisterApiOpenAPIRouter(router gin.IRoutes) {
	a := &apiOpenAPIController{}
	bean.Autowired(&a.apiOpenAPIService)
	bean.Autowired(&a.extAPPService)

	router.GET("/apis/import", a.getImportInfo)
	router.POST("/apis/import", a.syncAPI)
}

func (a *apiOpenAPIController) getImportInfo(ginCtx *gin.Context) {
	//检测openAPI token
	token := ginCtx.GetHeader("Authorization")
	namespaceID := getNamespaceId(ginCtx)
	_, err := a.extAPPService.CheckExtAPPToken(ginCtx, namespaceID, token)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetAPIImportInfos fail. err:%s", err)))
		return
	}

	groups, services, formats, err := a.apiOpenAPIService.GetSyncImportInfo(ginCtx, namespaceID)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetAPIImportInfos fail. err:%s", err)))
		return
	}

	data := make(map[string]interface{})
	data["groups"] = groups
	data["upstreams"] = services
	data["formats"] = formats
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (a *apiOpenAPIController) syncAPI(ginCtx *gin.Context) {
	//检测openAPI token 并获取相应外部应用的id
	token := ginCtx.GetHeader("Authorization")
	namespaceID := getNamespaceId(ginCtx)
	appID, err := a.extAPPService.CheckExtAPPToken(ginCtx, namespaceID, token)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI fail. err:%s", err)))
		return
	}

	inputData := new(dto.SyncImportData)
	//组装同步信息
	contentType := ginCtx.ContentType()

	switch contentType {
	case "multipart/form-data":
		fileInfo, err := ginCtx.FormFile("file")
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI get file fail. err: %s. ", err)))
			return
		}
		file, err := fileInfo.Open()
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI open file fail. err: %s. ", err)))
			return
		}
		fileData := make([]byte, fileInfo.Size)
		if _, err = file.Read(fileData); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI  read file fail. err: %s. ", err)))
			return
		}
		defer file.Close()

		inputData.Content = fileData
		inputData.Format = ginCtx.PostForm("format")
		inputData.ServiceName = ginCtx.PostForm("upstream")
		inputData.GroupUUID = ginCtx.PostForm("group")
		inputData.Label = ginCtx.PostForm("label")
		inputData.Prefix = ginCtx.PostForm("prefix")

		nodesForm := ginCtx.PostForm("nodes")
		if nodesForm != "" {
			server, err := a.getServer(nodesForm)
			if err != nil {
				ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
				return
			}
			inputData.Server = server
		}
	case "application/x-www-form-urlencoded":
		contentForm := ginCtx.PostForm("content")
		isBase64Encode := ginCtx.PostForm("encode") == "base64"
		content := []byte(contentForm)
		if isBase64Encode {
			content, err = common.Base64Decode(contentForm)
			if err != nil {
				ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI decode content fail. err: %s. ", err)))
				return
			}
		}

		inputData.Content = content
		inputData.Format = ginCtx.PostForm("format")
		inputData.ServiceName = ginCtx.PostForm("upstream")
		inputData.GroupUUID = ginCtx.PostForm("group")
		inputData.Label = ginCtx.PostForm("label")
		inputData.Prefix = ginCtx.PostForm("prefix")

		nodesForm := ginCtx.PostForm("nodes")
		if nodesForm != "" {
			server, err := a.getServer(nodesForm)
			if err != nil {
				ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
				return
			}
			inputData.Server = server
		}
	case "application/json":
		if err = ginCtx.BindJSON(inputData); err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
			return
		}
		//检查Server
		server := inputData.Server
		if server != nil {
			schme := strings.ToUpper(server.Scheme)
			server.Scheme = schme
			if schme != "HTTP" && schme != "HTTPS" {
				ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI fail. err: server.scheme is illegal. ")))
				return
			}
			if len(server.Nodes) == 0 {
				ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI fail. err: server.nodes is null. ")))
				return
			}
			for _, node := range server.Nodes {
				if !common.IsMatchDomainPort(node.Url) && !common.IsMatchIpPort(node.Url) {
					ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI fail. err: server.nodes.url %s is illegal. ", node.Url)))
					return
				}
			}
		}
	default:
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI fail. err: Content-Type illegal. ")))
		return
	}

	//校验服务名是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, inputData.ServiceName); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI upstream 't is illegal. err: %s. ", err)))
		return
	}

	checkList, err := a.apiOpenAPIService.SyncImport(ginCtx, namespaceID, appID, inputData)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("syncAPI fail. err:%s", err)))
		return
	}

	resultList := make([]dto.ImportAPIListItem, 0)
	for _, item := range checkList {
		resultList = append(resultList, dto.ImportAPIListItem{
			Name:   item.Name,
			Method: item.Method,
			Path:   item.Path,
			Desc:   item.Desc,
			Status: enum.ImportStatusType(item.Status),
		})
	}
	data := make(map[string]interface{})
	data["apis"] = resultList

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// getServer 解析scheme://url weight;url weight
func (a *apiOpenAPIController) getServer(nodesForm string) (*dto.ImportServerInfo, error) {
	nodes := strings.Split(nodesForm, ";")
	schemeIdx := strings.Index(nodes[0], "://")
	if schemeIdx < 0 {
		return nil, fmt.Errorf("syncAPI decode nodes fail. err: nodes %s is illegal. ", nodesForm)
	}

	server := new(dto.ImportServerInfo)
	scheme := strings.ToUpper(nodes[0][:schemeIdx])
	server.Scheme = scheme
	if scheme != "HTTP" && scheme != "HTTPS" {
		return nil, fmt.Errorf("syncAPI decode nodes fail. err: nodes %s is illegal. ", nodesForm)
	}

	upstreams := make([]*dto.ImportNodesInfo, 0, len(nodes))
	nodeInfo, err := a.getNodeInfo(nodes[0][schemeIdx+3:])
	if err != nil {
		return nil, err
	}
	upstreams = append(upstreams, nodeInfo)
	if len(nodes) > 1 {
		for _, node := range nodes[1:] {
			if node == "" {
				continue
			}
			nodeInfo, err = a.getNodeInfo(node)
			if err != nil {
				return nil, err
			}
			upstreams = append(upstreams, nodeInfo)
		}
	}
	server.Nodes = upstreams
	return server, nil
}

func (a *apiOpenAPIController) getNodeInfo(nodeStr string) (*dto.ImportNodesInfo, error) {
	idx := strings.Index(nodeStr, " ")
	url := nodeStr[:idx]
	weight, err := strconv.Atoi(nodeStr[idx+1:])
	if err != nil {
		weight = 1
	}

	//若同时不符合ip:port 或者域名 域名:port则报错
	if !common.IsMatchIpPort(url) && !common.IsMatchDomainPort(url) {
		return nil, fmt.Errorf("syncAPI decode nodes fail. err: url %s is illegal. ", url)
	}

	return &dto.ImportNodesInfo{
		Url:    url,
		Weight: weight,
	}, nil
}
