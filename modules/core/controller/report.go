package controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/report"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *Module) reportStatus(ginCtx *gin.Context) {
	consoleId, nodes, need := report.NeedReport()
	if !need {
		ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
			"report": "no",
		}))

		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"report": "yes",
		"id":     consoleId,
		"nodes":  nodes,
	}))

}

func (m *Module) updateReport(ginCtx *gin.Context) {
	report.UpdateReport()
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))

}
