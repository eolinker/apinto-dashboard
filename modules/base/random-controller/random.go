package random_controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/random"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type RandomController struct {
	randomService random.IRandomService
}

func NewRandomController() *RandomController {
	r := &RandomController{}
	bean.Autowired(&r.randomService)
	return r
}

func (r *RandomController) GET(ginCtx *gin.Context) {
	template := ginCtx.Param("template")
	randomStr := r.randomService.RandomStr(template)
	m := common.Map[string, interface{}]{}

	m["id"] = strings.ToLower(randomStr)

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}
