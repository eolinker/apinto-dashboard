package env_controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/base/env-model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type enumController struct {
}

func RegisterEnumRouter(router gin.IRoutes) {

	c := &enumController{}
	router.GET("/enum/envs", c.getEnv)
}

func (e *enumController) getEnv(context *gin.Context) {
	//enums, err := e.enumService.GetByType(entry.EnvType)
	//if err != nil {
	//	context.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
	//	return
	//}
	enums := env_model.EnumValueList

	list := make([]*env_model.EnumEnvOut, 0, len(enums))
	for _, value := range enums {
		list = append(list, &env_model.EnumEnvOut{
			Name:  value.String(),
			Value: value,
		})
	}
	m := common.Map[string, interface{}]{}
	m["envs"] = list
	context.JSON(http.StatusOK, controller.NewSuccessResult(m))
}
