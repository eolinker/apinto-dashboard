package controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
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
	enums := enum.EnumValueList

	list := make([]*enum.EnumEnvOut, 0, len(enums))
	for _, value := range enums {
		list = append(list, &enum.EnumEnvOut{
			Name:  value.String(),
			Value: value,
		})
	}
	m := common.Map[string, interface{}]{}
	m["envs"] = list
	context.JSON(http.StatusOK, dto.NewSuccessResult(m))
}
