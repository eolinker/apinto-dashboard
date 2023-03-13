package controller

import (
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
)

type commonGroupController struct {
	commonGroupService service.ICommonGroupService
}

func RegisterCommonGroupRouter(router gin.IRoutes) {
	c := &commonGroupController{}
	bean.Autowired(&c.commonGroupService)
	router.GET("/group/:group_type", c.groups)
	router.POST("/group/:group_type", logHandler(enum.LogOperateTypeCreate, enum.LogKindCommonGroup), c.createGroup)
	router.PUT("/group/:group_type/:uuid", logHandler(enum.LogOperateTypeEdit, enum.LogKindCommonGroup), c.updateGroup)
	router.DELETE("/group/:group_type/:uuid", logHandler(enum.LogOperateTypeDelete, enum.LogKindCommonGroup), c.delGroup)
	router.PUT("/groups/:group_type/sort", c.groupSort)
}

// groups 获取目录列表
func (c *commonGroupController) groups(ginCtx *gin.Context) {

	namespaceId := getNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	tagName := ginCtx.Query("tag_name")
	uuid := ginCtx.Query("parent_uuid")
	queryName := ginCtx.Query("query_name")

	root, apis, err := c.commonGroupService.GroupList(ginCtx, namespaceId, groupType, tagName, uuid, queryName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	resApis := make([]*dto.CommonGroupApi, 0, len(apis))
	groups := make([]*dto.CommonGroupOut, 0, len(root.CommonGroup))
	for _, group := range root.CommonGroup {
		value := &dto.CommonGroupOut{
			UUID: group.Group.Uuid,
			Name: group.Group.Name,
		}
		c.subGroup(value, group.Subgroup)
		groups = append(groups, value)
	}

	for _, api := range apis {
		resApis = append(resApis, &dto.CommonGroupApi{
			Name:      api.Name,
			UUID:      api.UUID,
			Methods:   api.Methods,
			GroupUUID: api.GroupUUID,
		})
	}

	resRoot := &dto.CommonGroupRootOut{
		UUID:   root.UUID,
		Name:   root.Name,
		Groups: groups,
	}
	m := make(map[string]interface{})
	m["root"] = resRoot
	m["apis"] = resApis
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
}
func (c *commonGroupController) subGroup(val *dto.CommonGroupOut, list []*model.CommonGroup) {
	if len(list) == 0 {
		return
	}
	for _, group := range list {
		commonGroup := &dto.CommonGroupOut{UUID: group.Group.Uuid, Name: group.Group.Name}
		val.Children = append(val.Children, commonGroup)
		c.subGroup(commonGroup, group.Subgroup)
	}
}

// updateGroup 修改目录
func (c *commonGroupController) updateGroup(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	uuid := ginCtx.Param("uuid")
	operator := getUserId(ginCtx)
	input := new(dto.CommonGroupInput)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := c.commonGroupService.UpdateGroup(ginCtx, namespaceId, operator, groupType, input.Name, uuid); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (c *commonGroupController) groupSort(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	tagName := ginCtx.Query("tag_name")

	input := &dto.CommGroupSortInput{}
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := c.commonGroupService.GroupSort(ginCtx, namespaceId, groupType, tagName, input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// delGroup 删除目录
func (c *commonGroupController) delGroup(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	uuid := ginCtx.Param("uuid")
	operator := getUserId(ginCtx)

	if err := c.commonGroupService.DeleteGroup(ginCtx, namespaceId, operator, groupType, uuid); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// createGroup 新建目录
func (c *commonGroupController) createGroup(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	groupType := ginCtx.Param("group_type")
	tagName := ginCtx.Query("tag_name")
	operator := getUserId(ginCtx)
	input := new(dto.CommonGroupInput)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := c.commonGroupService.CreateGroup(ginCtx, namespaceId, operator, groupType, tagName, input.Name, input.UUID, input.ParentUUID); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}
