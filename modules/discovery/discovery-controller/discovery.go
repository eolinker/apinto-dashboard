package discovery_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	discovery_driver "github.com/eolinker/apinto-dashboard/driver-manager/driver"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/discovery"
	"github.com/eolinker/apinto-dashboard/modules/discovery/discover-dto"
	"github.com/eolinker/apinto-dashboard/modules/online/online-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type discoveryController struct {
	discoveryService discovery.IDiscoveryService
}

func RegisterDiscoveryRouter(router gin.IRouter) {
	c := &discoveryController{}
	bean.Autowired(&c.discoveryService)

	router.GET("/discoveries", controller.GenAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getList)
	router.GET("/discovery", controller.GenAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getInfo)
	router.POST("/discovery", controller.GenAccessHandler(access.DiscoveryEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindDiscovery), c.create)
	router.PUT("/discovery", controller.GenAccessHandler(access.DiscoveryEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindDiscovery), c.update)
	router.DELETE("/discovery", controller.GenAccessHandler(access.DiscoveryEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindDiscovery), c.delete)
	router.GET("/discovery/enum", c.getEnum)
	router.GET("/discovery/drivers", controller.GenAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getDrivers)

	router.PUT("/discovery/:discovery_name/online", controller.GenAccessHandler(access.DiscoveryEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindDiscovery), c.online)
	router.PUT("/discovery/:discovery_name/offline", controller.GenAccessHandler(access.DiscoveryEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindDiscovery), c.offline)
	router.GET("/discovery/:discovery_name/onlines", controller.GenAccessHandler(access.DiscoveryView, access.DiscoveryEdit), c.getOnlineList)
}

// getList 获取注册中心列表
func (d *discoveryController) getList(ginCtx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
	searchName := ginCtx.Query("name")

	listItem, err := d.discoveryService.GetDiscoveryList(ginCtx, namespaceID, searchName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetDiscoveryList fail. err:%s", err.Error())))
		return
	}
	discoveries := make([]*discover_dto.DiscoveryListItem, 0, len(listItem))

	for _, item := range listItem {
		discovery := &discover_dto.DiscoveryListItem{
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
	namespaceID := namespace_controller.GetNamespaceId(ginCtx)
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

	discovery := &discover_dto.DiscoveryInfoProxy{
		Name:   info.Name,
		UUID:   info.UUID,
		Driver: info.Driver,
		Desc:   info.Desc,
		Config: discover_dto.ConfigProxy(info.Config),
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(&discover_dto.DiscoveryInfoOutput{
		Discovery: discovery,
		Render:    dto.Render(info.Render),
	}))
}

// create 新建注册中心
func (d *discoveryController) create(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	operator := controller.GetUserId(ginCtx)

	input := new(discover_dto.DiscoveryInfoProxy)
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
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	discoveryName := ginCtx.Query("name")
	if discoveryName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateDiscovery Info fail. err: discoveryName can't be nil")))
		return
	}
	operator := controller.GetUserId(ginCtx)

	input := new(discover_dto.DiscoveryInfoProxy)
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
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	discoveryName := ginCtx.Query("name")
	if discoveryName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteDiscovery Info fail. err: discoveryName can't be nil")))
		return
	}

	userId := controller.GetUserId(ginCtx)
	err := d.discoveryService.DeleteDiscovery(ginCtx, namespaceId, userId, discoveryName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteDiscovery fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// getEnum 获取作为选项的注册中心列表
func (d *discoveryController) getEnum(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	enumList, err := d.discoveryService.GetDiscoveryEnum(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetDiscoveryEnum fail. err:%s", err.Error())))
		return
	}
	discoveries := make([]*discover_dto.DiscoveryEnum, 0, len(enumList))
	for _, enum := range enumList {
		discovery := &discover_dto.DiscoveryEnum{
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

	drivers := make([]*discover_dto.DriversItem, 0, len(driverList))
	for _, driver := range driverList {
		d := &discover_dto.DriversItem{
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
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	discoveryName := ginCtx.Param("discovery_name")
	input := &online_dto.UpdateOnlineStatusInput{}
	operator := controller.GetUserId(ginCtx)

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
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
	discoveryName := ginCtx.Param("discovery_name")
	input := &online_dto.UpdateOnlineStatusInput{}
	operator := controller.GetUserId(ginCtx)

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
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)
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

	resp := make([]*online_dto.OnlineOut, 0, len(list))
	for _, online := range list {
		updateTime := ""
		if !online.UpdateTime.IsZero() {
			updateTime = common.TimeToStr(online.UpdateTime)
		}
		resp = append(resp, &online_dto.OnlineOut{
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
