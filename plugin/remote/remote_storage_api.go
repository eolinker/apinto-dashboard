package local

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

func (r *remotePluginController) createRemoteStorageApis() []apinto_module.RouterInfo {
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
		module["url"] = fmt.Sprintf("%s/%s", strings.TrimSuffix(r.define.Server, "/"), strings.TrimPrefix(r.define.Path, "/"))

		//TODO 补充信息

		data["module"] = module
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
	}
}

func (r *remotePluginController) saveObject() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		module := ginCtx.Param("module")
		key := ginCtx.Param("key")

		//TODO
		object, err := ginCtx.GetRawData()
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		err = r.remoteStorageService.Save(module, key, object)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
	}
}

func (r *remotePluginController) getObject() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		module := ginCtx.Param("module")
		key := ginCtx.Param("key")

		storageInfo, err := r.remoteStorageService.Get(module, key)
		if err != nil {
			controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
			return
		}

		data := common.Map{}
		data[key] = storageInfo.Object
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
	}
}
