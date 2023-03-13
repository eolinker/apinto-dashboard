package controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	discovery_driver "github.com/eolinker/apinto-dashboard/driver-manager/driver"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type discoveryController struct {
	discoveryService service.IDiscoveryService
}

func RegisterDiscoveryRouter(router gin.IRouter) {
	c := &discoveryController{}
	bean.Autowired(&c.discoveryService)

	router.GET("/discoveries", genAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getList)
	router.GET("/discovery", genAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getInfo)
	router.POST("/discovery", genAccessHandler(access.DiscoveryEdit), logHandler(enum.LogOperateTypeCreate, enum.LogKindDiscovery), c.create)
	router.PUT("/discovery", genAccessHandler(access.DiscoveryEdit), logHandler(enum.LogOperateTypeEdit, enum.LogKindDiscovery), c.update)
	router.DELETE("/discovery", genAccessHandler(access.DiscoveryEdit), logHandler(enum.LogOperateTypeDelete, enum.LogKindDiscovery), c.delete)
	router.GET("/discovery/enum", c.getEnum)
	router.GET("/discovery/drivers", genAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getDrivers)

	router.PUT("/discovery/:discovery_name/online", genAccessHandler(access.DiscoveryEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindDiscovery), c.online)
	router.PUT("/discovery/:discovery_name/offline", genAccessHandler(access.DiscoveryEdit), logHandler(enum.LogOperateTypePublish, enum.LogKindDiscovery), c.offline)
	router.GET("/discovery/:discovery_name/onlines", genAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getOnlineList)
}

// getList 获取注册中心列表
func (d *discoveryController) getList(ginCtx *gin.Context) {
	namespaceID := getNamespaceId(ginCtx)
	searchName := ginCtx.Query("name")

	listItem, err := d.discoveryService.GetDiscoveryList(ginCtx, namespaceID, searchName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetDiscoveryList fail. err:%s", err.Error())))
		return
	}
	discoveries := make([]*dto.DiscoveryListItem, 0, len(listItem))

	for _, item := range listItem {
		discovery := &dto.DiscoveryListItem{
			Name:       item.Name,
			UUID:       item.UUID,
			Driver:     item.Driver,
			Desc:       item.Desc,
			UpdateTime: common.TimeToStr(item.UpdateTime),
			IsDelete:   item.IsDelete,
		}

		discoveries = append(discoveries, discovery)
	}

	data := common.Map[string, interface{}]{}
	data["discoveries"] = discoveries
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// getInfo 获取注册中心配置信息
func (d *discoveryController) getInfo(ginCtx *gin.Context) {
	namespaceID := getNamespaceId(ginCtx)
	discoveryName := ginCtx.Query("name")
	if discoveryName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetDiscoveryInfo fail. err: discoveryName can't be nil")))
		return
	}
	info, err := d.discoveryService.GetDiscoveryVersionInfo(ginCtx, namespaceID, discoveryName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetDiscoveryInfo fail. err:%s", err.Error())))
		return
	}

	discovery := &dto.DiscoveryInfoProxy{
		Name:   info.Name,
		UUID:   info.UUID,
		Driver: info.Driver,
		Desc:   info.Desc,
		Config: dto.ConfigProxy(info.Config),
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(&dto.DiscoveryInfoOutput{
		Discovery: discovery,
		Render:    dto.Render(info.Render),
	}))
}

// create 新建注册中心
func (d *discoveryController) create(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	operator := getUserId(ginCtx)

	input := new(dto.DiscoveryInfoProxy)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	//校验注册中心名是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, input.Name); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	input.Driver = strings.ToLower(input.Driver)

	if input.Name == discovery_driver.DriverStatic {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("discoveryName can't be static. "))
		return
	}

	err := d.discoveryService.CreateDiscovery(ginCtx, namespaceId, operator, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateDiscovery fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// alter 修改注册中心
func (d *discoveryController) update(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	discoveryName := ginCtx.Query("name")
	if discoveryName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateDiscovery Info fail. err: discoveryName can't be nil")))
		return
	}
	operator := getUserId(ginCtx)

	input := new(dto.DiscoveryInfoProxy)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	input.Name = discoveryName
	//input.UUID = uuid TODO 这里需要获取rest的uuid
	err := d.discoveryService.UpdateDiscovery(ginCtx, namespaceId, operator, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateDiscovery fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// delete 删除注册中心
func (d *discoveryController) delete(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	discoveryName := ginCtx.Query("name")
	if discoveryName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteDiscovery Info fail. err: discoveryName can't be nil")))
		return
	}

	userId := getUserId(ginCtx)
	err := d.discoveryService.DeleteDiscovery(ginCtx, namespaceId, userId, discoveryName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteDiscovery fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// getEnum 获取作为选项的注册中心列表
func (d *discoveryController) getEnum(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)

	enumList, err := d.discoveryService.GetDiscoveryEnum(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetDiscoveryEnum fail. err:%s", err.Error())))
		return
	}
	discoveries := make([]*dto.DiscoveryEnum, 0, len(enumList))
	for _, enum := range enumList {
		discovery := &dto.DiscoveryEnum{
			Name:   enum.Name,
			Driver: enum.Driver,
			Render: dto.Render(enum.Render),
		}
		discoveries = append(discoveries, discovery)
	}

	data := common.Map[string, interface{}]{}
	data["discoveries"] = discoveries
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// getDrivers 获取可用的驱动列表
func (d *discoveryController) getDrivers(ginCtx *gin.Context) {
	driverList := d.discoveryService.GetDriversRender()

	drivers := make([]*dto.DriversItem, 0, len(driverList))
	for _, driver := range driverList {
		d := &dto.DriversItem{
			Name:   driver.Name,
			Render: dto.Render(driver.Render),
		}
		drivers = append(drivers, d)
	}

	data := common.Map[string, interface{}]{}
	data["drivers"] = drivers
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// online 上线
func (d *discoveryController) online(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	discoveryName := ginCtx.Param("discovery_name")
	input := &dto.UpdateOnlineStatusInput{}
	operator := getUserId(ginCtx)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	router, err := d.discoveryService.OnlineDiscovery(ginCtx, namespaceId, operator, discoveryName, input.ClusterName)
	if err != nil && router == nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	} else if err == nil {
		ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
		return
	}

	msg := ""
	if err != nil {
		msg = err.Error()
	}

	m := common.Map[string, interface{}]{}
	m["router"] = router
	ginCtx.JSON(http.StatusOK, dto.Result{
		Code: -1,
		Data: m,
		Msg:  msg,
	})

}

// online 下线
func (d *discoveryController) offline(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	discoveryName := ginCtx.Param("discovery_name")
	input := &dto.UpdateOnlineStatusInput{}
	operator := getUserId(ginCtx)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := d.discoveryService.OfflineDiscovery(ginCtx, namespaceId, operator, discoveryName, input.ClusterName); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// getOnlineList 上线管理列表
func (d *discoveryController) getOnlineList(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)
	discoveryName := ginCtx.Param("discovery_name")
	if discoveryName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult("discovery_name can't be nil. "))
		return
	}
	list, err := d.discoveryService.OnlineList(ginCtx, namespaceId, discoveryName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resp := make([]*dto.OnlineOut, 0, len(list))
	for _, online := range list {
		updateTime := ""
		if !online.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(online.UpdateTime)
		}
		resp = append(resp, &dto.OnlineOut{
			Name:       online.ClusterName,
			Status:     enum.OnlineStatus(online.Status),
			Env:        online.Env,
			Operator:   online.Operator,
			UpdateTime: updateTime,
		})
	}

	m := common.Map[string, interface{}]{}
	m["clusters"] = resp
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
}
