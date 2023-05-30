package group_controller

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/group"
	"github.com/eolinker/apinto-dashboard/modules/group/group-dto"
	"github.com/eolinker/apinto-dashboard/modules/group/group-model"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type commonGroupController struct {
	commonGroupService group.ICommonGroupService
}

var (
	locker             sync.Mutex
	controllerInstance *commonGroupController
)

func NewCommonGroupController() *commonGroupController {
	if controllerInstance == nil {
		locker.Lock()
		defer locker.Unlock()
		if controllerInstance == nil {
			controllerInstance = &commonGroupController{}
			bean.Autowired(&controllerInstance.commonGroupService)
		}
	}
	return controllerInstance

}

// groups 获取目录列表
func (c *commonGroupController) groups(ginCtx *gin.Context) {

	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	tagName := ginCtx.Query("tag_name")
	uuid := ginCtx.Query("parent_uuid")
	queryName := ginCtx.Query("query_name")

	root, apis, err := c.commonGroupService.GroupList(ginCtx, namespaceId, groupType, tagName, uuid, queryName)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}
	resApis := make([]*group_dto.CommonGroupApi, 0, len(apis))
	groups := make([]*group_dto.CommonGroupOut, 0, len(root.CommonGroup))
	for _, group := range root.CommonGroup {
		value := &group_dto.CommonGroupOut{
			UUID: group.Group.Uuid,
			Name: group.Group.Name,
		}
		c.subGroup(value, group.Subgroup)
		groups = append(groups, value)
	}

	for _, api := range apis {
		resApis = append(resApis, &group_dto.CommonGroupApi{
			Name:      api.Name,
			UUID:      api.UUID,
			Methods:   api.Methods,
			GroupUUID: api.GroupUUID,
		})
	}

	resRoot := &group_dto.CommonGroupRootOut{
		UUID:   root.UUID,
		Name:   root.Name,
		Groups: groups,
	}
	m := make(map[string]interface{})
	m["root"] = resRoot
	m["apis"] = resApis
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(m))
}
func (c *commonGroupController) subGroup(val *group_dto.CommonGroupOut, list []*group_model.CommonGroup) {
	if len(list) == 0 {
		return
	}
	for _, group := range list {
		commonGroup := &group_dto.CommonGroupOut{UUID: group.Group.Uuid, Name: group.Group.Name}
		val.Children = append(val.Children, commonGroup)
		c.subGroup(commonGroup, group.Subgroup)
	}
}

// UpdateGroup 修改目录
func (c *commonGroupController) UpdateGroup(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	uuid := ginCtx.Param("uuid")
	operator := users.GetUserId(ginCtx)
	input := new(group_dto.CommonGroupInput)

	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := c.commonGroupService.UpdateGroup(ginCtx, namespaceId, operator, groupType, input.Name, uuid); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (c *commonGroupController) GroupSort(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	tagName := ginCtx.Query("tag_name")

	input := &group_dto.CommGroupSortInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if err := c.commonGroupService.GroupSort(ginCtx, namespaceId, groupType, tagName, input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// DelGroup 删除目录
func (c *commonGroupController) DelGroup(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	uuid := ginCtx.Param("uuid")
	operator := users.GetUserId(ginCtx)

	if err := c.commonGroupService.DeleteGroup(ginCtx, namespaceId, operator, groupType, uuid); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

// CreateGroup 新建目录
func (c *commonGroupController) CreateGroup(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	tagName := ginCtx.Query("tag_name")
	operator := users.GetUserId(ginCtx)
	input := new(group_dto.CommonGroupInput)

	if err := ginCtx.BindJSON(input); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	if _, err := c.commonGroupService.CreateGroup(ginCtx, namespaceId, operator, groupType, tagName, input.Name, input.UUID, input.ParentUUID); err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}
