package variable_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type variablesController struct {
	globalVariableService variable.IGlobalVariableService
}

func RegisterVariablesRouter(router gin.IRoutes) {
	c := &variablesController{}
	bean.Autowired(&c.globalVariableService)

	router.GET("/variables", c.gets)
	router.GET("/variable", c.get)
	router.POST("/variable", controller.AuditLogHandler(enum.LogOperateTypeCreate, enum.LogKindGlobalVariable, c.post))
	router.DELETE("/variable", controller.AuditLogHandler(enum.LogOperateTypeDelete, enum.LogKindGlobalVariable, c.del))
}

// 获取全局环境变量列表
func (e *variablesController) gets(ginCtx *gin.Context) {
	pageNumStr := ginCtx.Query("page_num")
	pageSizeStr := ginCtx.Query("page_size")
	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 20
	}

	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	key := ginCtx.Query("key")
	status := ginCtx.Query("status")
	if status != "" && !enum.CheckVariableUsageStatus(status) {
		controller.ErrorJson(ginCtx, http.StatusOK, "status is illegal. ")
		return
	}

	variableList, total, err := e.globalVariableService.List(ginCtx, pageNum, pageSize, namespaceID, key, enum.GetStatusIndexByName(status))
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("GetGlobalVariable fail. err:%s", err.Error())))
		return
	}
	resp := make([]*variable_dto.GlobalVariableListItem, 0, len(variableList))
	for _, item := range variableList {
		data := &variable_dto.GlobalVariableListItem{
			Key:         item.Key,
			Status:      enum.VariableUsageStatus(item.Status),
			Description: item.Desc,
			Operator:    item.OperatorStr,
			CreateTime:  common.TimeToStr(item.CreateTime),
		}
		resp = append(resp, data)
	}

	data := common.Map[string, interface{}]{}
	data["variables"] = resp
	data["total"] = total
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 获取单个全局环境变量使用信息
func (e *variablesController) get(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	key := ginCtx.Query("key")
	if key == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("GetGlobalVariable Info fail. err: key can't be nil")))
		return
	}
	variableDetails, err := e.globalVariableService.GetInfo(ginCtx, namespaceID, key)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("GetGlobalVariable Info fail. err:%s", err.Error())))
		return
	}
	list := make([]*variable_dto.GlobalVariableDetailsItem, 0, len(variableDetails))

	for _, item := range variableDetails {
		data := &variable_dto.GlobalVariableDetailsItem{
			ClusterName: item.ClusterName,
			Environment: item.Environment,
			Value:       item.Value,
			Status:      enum.PublishType(item.Status),
		}
		list = append(list, data)
	}

	data := common.Map[string, interface{}]{}
	data["variables"] = list
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

// 新增全局环境变量
func (e *variablesController) post(ginCtx *gin.Context) {

	input := &variable_dto.GlobalVariableInput{}

	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if input.Key == "" {
		controller.ErrorJson(ginCtx, http.StatusOK, "parameter error")
		return
	}

	if err := common.IsMatchString(common.EnglishOrNumber_, input.Key); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	userId := controller.GetUserId(ginCtx)
	_, err := e.globalVariableService.Create(ginCtx, namespaceID, userId, input.Key, input.Desc)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Create GlobalVariable fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// 删除全局环境变量
func (e *variablesController) del(ginCtx *gin.Context) {
	key := ginCtx.Query("key")
	if key == "" {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Delete GlobalVariable fail. err: key can't be nil")))
		return
	}

	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	userID := controller.GetUserId(ginCtx)
	err := e.globalVariableService.Delete(ginCtx, namespaceID, userID, key)
	if err != nil {
		ginCtx.JSON(http.StatusOK, controller.NewErrorResult(fmt.Sprintf("Delete GlobalVariable fail. err:%s", err.Error())))
		return
	}
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
