package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panic %v", common.PanicTrace(err))
			fmt.Println(common.PanicTrace(err))
			c.JSON(http.StatusInternalServerError, controller.NewErrorResult("服务器内部错误"))
		}
	}()

	c.Next()
}
