package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/dto/variable-dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type variablesController struct {
	globalVariableService service.IGlobalVariableService
}

func RegisterVariablesRouter(router gin.IRoutes) {
	c := &variablesController{}
	bean.Autowired(&c.globalVariableService)

	router.GET("/variables", GenAccessHandler(access.VariableView, access.VariableEdit, access.ServiceView, access.DiscoveryView), c.gets)
	router.GET("/variable", GenAccessHandler(access.VariableView, access.VariableEdit), c.get)
	router.POST("/variable", GenAccessHandler(access.VariableEdit, access.ServiceEdit, access.DiscoveryEdit), LogHandler(enum.LogOperateTypeCreate, enum.LogKindGlobalVariable), c.post)
	router.DELETE("/variable", GenAccessHandler(access.VariableEdit), LogHandler(enum.LogOperateTypeDelete, enum.LogKindGlobalVariable), c.del)
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

	namespaceID := GetNamespaceId(ginCtx)
	key := ginCtx.Query("key")
	status := ginCtx.Query("status")
	if status != "" && !enum.CheckVariableUsageStatus(status) {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("status is illegal. "))
		return
	}

	variableList, total, err := e.globalVariableService.List(ginCtx, pageNum, pageSize, namespaceID, key, enum.GetStatusIndexByName(status))
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetGlobalVariable fail. err:%s", err.Error())))
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
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// 获取单个全局环境变量使用信息
func (e *variablesController) get(ginCtx *gin.Context) {
	namespaceID := GetNamespaceId(ginCtx)
	key := ginCtx.Query("key")
	if key == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetGlobalVariable Info fail. err: key can't be nil")))
		return
	}
	variableDetails, err := e.globalVariableService.GetInfo(ginCtx, namespaceID, key)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetGlobalVariable Info fail. err:%s", err.Error())))
		return
	}
	list := make([]*variable_dto.GlobalVariableDetailsItem, 0, len(variableDetails))

	for _, item := range variableDetails {
		data := &variable_dto.GlobalVariableDetailsItem{
			ClusterName: item.ClusterName,
			Environment: item.Environment,
			Value:       item.Value,
			Status:      enum.ClusterVariablePublish(item.Status),
		}
		list = append(list, data)
	}

	data := common.Map[string, interface{}]{}
	data["variables"] = list
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// 新增全局环境变量
func (e *variablesController) post(ginCtx *gin.Context) {

	input := &variable_dto.GlobalVariableInput{}

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if input.Key == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("parameter error"))
		return
	}

	if err := common.IsMatchString(common.EnglishOrNumber_, input.Key); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	namespaceID := GetNamespaceId(ginCtx)
	userId := GetUserId(ginCtx)
	_, err := e.globalVariableService.Create(ginCtx, namespaceID, userId, input.Key, input.Desc)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Create GlobalVariable fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// 删除全局环境变量
func (e *variablesController) del(ginCtx *gin.Context) {
	key := ginCtx.Query("key")
	if key == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Delete GlobalVariable fail. err: key can't be nil")))
		return
	}

	namespaceID := GetNamespaceId(ginCtx)
	userID := GetUserId(ginCtx)
	err := e.globalVariableService.Delete(ginCtx, namespaceID, userID, key)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("Delete GlobalVariable fail. err:%s", err.Error())))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}
