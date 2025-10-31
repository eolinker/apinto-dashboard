package controller

import (
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/eosc/common/bean"
)

type AccessController struct {
	accessService mpm3.IAccessService
}

func NewAccessController() *AccessController {
	c := &AccessController{}
	bean.Autowired(&c.accessService)
	return c
}

//
//func (c *AccessController) AccessList(ginCtx *gin.Context) {
//	accessM := c.accessService.GetEnable(ginCtx)
//
//	access := make([]string, 0, len(accessM))
//	for _, m := range accessM {
//		access = append(access, m.Name)
//	}
//	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
//		"access": access,
//	}))
//}
