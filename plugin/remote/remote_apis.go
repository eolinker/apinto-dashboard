package remote

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/modules/remote_storage"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func newRemotePluginController(moduleName string, cfg *Config, define *Define) *remotePluginController {
	r := &remotePluginController{
		moduleName: moduleName,
		cfg:        cfg,
		define:     define,
	}
	bean.Autowired(&r.remoteStorageService)
	return r
}

type remotePluginController struct {
	remoteStorageService remote_storage.IRemoteStorageService
	moduleName           string
	cfg                  *Config
	define               *Define
}

func (r *remotePluginController) createRemoteApis() []apinto_module.RouterInfo {
	return []apinto_module.RouterInfo{
		{
			Method:      "GET",
			Path:        fmt.Sprintf("/api/remote/%s", r.moduleName),
			Handler:     fmt.Sprintf("%s.getOpenMode", r.moduleName),
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{r.getOpenMode()},
		},
		{
			Method:      "GET",
			Path:        fmt.Sprintf("/api/remote/%s/store/:key", r.moduleName),
			Handler:     fmt.Sprintf("%s.getRemoteObject", r.moduleName),
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{r.getObject()},
		},
		{
			Method:      "PUT",
			Path:        fmt.Sprintf("/api/remote/%s/store/:key", r.moduleName),
			Handler:     fmt.Sprintf("%s.saveRemoteObject", r.moduleName),
			Labels:      apinto_module.RouterLabelApi,
			HandlerFunc: []apinto_module.HandlerFunc{r.saveObject()},
		},
	}
}

func (r *remotePluginController) getOpenMode() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		data := common.Map{}

		module := common.Map{}
		module["name"] = r.moduleName

		server := ""
		if r.define.Internet {
			server = r.define.Server
		} else {
			server = r.cfg.Server
		}

		url := strings.TrimSuffix(server, "/")
		if r.define.Path != "" {
			url = fmt.Sprintf("%s/%s", url, strings.TrimPrefix(r.define.Path, "/"))
		}
		module["url"] = url

		module["query"] = configParamToDto(r.cfg.Query, r.define.Querys)
		module["header"] = configParamToDto(r.cfg.Header, r.define.Headers)
		module["initialize"] = configParamToDto(r.cfg.Initialize, r.define.Initialize)

		data["module"] = module
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
	}
}

// configParamToDto 将config内的query，header，initialize转化成接口需要返回的格式
func configParamToDto(configParams map[string]string, defineParams []ExtendParamsRender) []*DtoOpenModeParam {
	items := make([]*DtoOpenModeParam, 0, len(configParams))

	defineParamsMap := common.SliceToMap(defineParams, func(t ExtendParamsRender) string {
		return t.Name
	})
	for name, value := range configParams {
		item := &DtoOpenModeParam{
			Name:  name,
			Value: value,
		}
		if info, has := defineParamsMap[name]; has {
			item.Type = info.Type
		}
		items = append(items, item)
	}

	return items
}

func (r *remotePluginController) saveObject() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		key := ginCtx.Param("key")

		object := new(interface{})
		if err := ginCtx.BindJSON(object); err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}

		err := r.remoteStorageService.Save(r.moduleName, key, object)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
	}
}

func (r *remotePluginController) getObject() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		key := ginCtx.Param("key")

		storageInfo, err := r.remoteStorageService.Get(r.moduleName, key)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}

		data := common.Map{}
		data[key] = storageInfo.Object
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
	}
}
