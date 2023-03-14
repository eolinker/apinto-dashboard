package upstream_controller

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/dto/online-dto"
	"github.com/eolinker/apinto-dashboard/dto/service-dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	_ "github.com/eolinker/apinto-dashboard/modules/upstream/service"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type serviceController struct {
	service   upstream.IService
	discovery service.IDiscoveryService
}

func RegisterServiceRouter(router gin.IRouter) {
	c := &serviceController{}
	bean.Autowired(&c.service)
	bean.Autowired(&c.discovery)

	router.GET("/services", controller.GenAccessHandler(access.ServiceView, access.ServiceEdit), c.getList)
	router.GET("/service", controller.GenAccessHandler(access.ServiceView, access.ServiceEdit), c.getInfo)
	router.PUT("/service", controller.GenAccessHandler(access.ServiceEdit), controller.LogHandler(enum.LogOperateTypeEdit, enum.LogKindService), c.alter)
	router.POST("/service", controller.GenAccessHandler(access.ServiceEdit), controller.LogHandler(enum.LogOperateTypeCreate, enum.LogKindService), c.create)
	router.DELETE("/service", controller.GenAccessHandler(access.ServiceEdit), controller.LogHandler(enum.LogOperateTypeDelete, enum.LogKindService), c.del)
	router.GET("/service/enum", c.getEnum)

	router.PUT("/service/:service_name/online", controller.GenAccessHandler(access.ServiceEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindService), c.online)
	router.PUT("/service/:service_name/offline", controller.GenAccessHandler(access.ServiceEdit), controller.LogHandler(enum.LogOperateTypePublish, enum.LogKindService), c.offline)
	router.GET("/service/:service_name/onlines", controller.GenAccessHandler(access.ServiceView, access.ServiceEdit), c.getOnlineList)

	//router.GET("/service/:service_name/api", c.getApi)
	//router.PUT("/service/:service_name/api", c.putApi)
	//router.POST("/service/:service_name/api", c.postApi)
	//router.DELETE("/service/:service_name/api", c.delApi)
	//
	//router.GET("/service/:service_name/:group_uuid/apis", c.getApis)
}

// getList 获取服务信息列表
func (s *serviceController) getList(ginCtx *gin.Context) {
	namespaceID := controller.GetNamespaceId(ginCtx)
	searchName := ginCtx.Query("name")
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

	backgroundCtx := ginCtx
	listItems, total, err := s.service.GetServiceList(backgroundCtx, namespaceID, searchName, pageNum, pageSize)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetServiceList fail. err:%s", err.Error())))
		return
	}
	services := make([]*service_dto.ServiceListItem, 0, len(listItems))

	for _, item := range listItems {
		li := &service_dto.ServiceListItem{
			Name:        item.Name,
			UUID:        item.UUID,
			Scheme:      item.Scheme,
			ServiceType: item.DriverName, //用枚举
			Config:      item.Config,
			UpdateTime:  common.TimeToStr(item.UpdateTime),
			IsDelete:    item.IsDelete,
		}
		services = append(services, li)
	}

	data := common.Map[string, interface{}]{}
	data["services"] = services
	data["total"] = total
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// getInfo 获取服务信息
func (s *serviceController) getInfo(ginCtx *gin.Context) {
	namespaceID := controller.GetNamespaceId(ginCtx)
	serviceName := ginCtx.Query("name")
	if serviceName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetServiceInfo Info fail. err: serviceName can't be nil")))
		return
	}
	backgroundCtx := ginCtx
	info, err := s.service.GetServiceInfo(backgroundCtx, namespaceID, serviceName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetServiceInfo fail. err:%s", err.Error())))
		return
	}

	discoveryName, _, driver, err := s.discovery.GetServiceDiscoveryDriverByID(backgroundCtx, info.DiscoveryId)
	if err != nil && driver == nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetServiceInfo fail. Get discovery driver fail. %s", err.Error())))
		return
	}
	conf := driver.FormatConfig([]byte(info.Config))

	serivce := &service_dto.ServiceInfoProxy{
		Name:          info.Name,
		UUID:          info.UUID,
		Desc:          info.Desc,
		Scheme:        info.Scheme,
		DiscoveryName: discoveryName,
		Config:        service_dto.ServiceConfigProxy(conf),
		Timeout:       info.Timeout,
		Balance:       info.Balance,
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(&service_dto.ServiceInfoOutput{
		Service: serivce,
		Render:  service_dto.Render(driver.Render()),
	}))
}

// create 新增服务
func (s *serviceController) create(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	operator := controller.GetUserId(ginCtx)

	inputProxy := new(service_dto.ServiceInfoProxy)
	if err := ginCtx.BindJSON(inputProxy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	//校验服务名是否合法
	if err := common.IsMatchString(common.EnglishOrNumber_, inputProxy.Name); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	backgroundCtx := ginCtx
	discoveryID, driverName, driver, err := s.discovery.GetServiceDiscoveryDriver(backgroundCtx, namespaceId, inputProxy.DiscoveryName)
	if err != nil && driver == nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateService fail. Get discovery driver fail. %s", err.Error())))
		return
	}
	newConf, formatAddr, variableList, err := driver.CheckInput(inputProxy.Config)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateService fail. Err: %s", err.Error())))
		return
	}
	inputProxy.Config = newConf

	input := &service_dto.ServiceInfo{
		Name:        inputProxy.Name,
		UUID:        inputProxy.UUID,
		Desc:        inputProxy.Desc,
		Scheme:      inputProxy.Scheme,
		DiscoveryID: discoveryID,
		DriverName:  driverName,
		FormatAddr:  formatAddr,
		Config:      inputProxy.Config.String(),
		Timeout:     inputProxy.Timeout,
		Balance:     inputProxy.Balance,
	}

	_, err = s.service.CreateService(backgroundCtx, namespaceId, operator, input, variableList)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateService fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// alter 修改服务信息
func (s *serviceController) alter(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	serviceName := ginCtx.Query("name")
	if serviceName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateService Info fail. err: serviceName can't be nil")))
		return
	}
	operator := controller.GetUserId(ginCtx)
	backgroundCtx := ginCtx
	inputProxy := new(service_dto.ServiceInfoProxy)
	if err := ginCtx.BindJSON(inputProxy); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	discoveryID, driverName, driver, err := s.discovery.GetServiceDiscoveryDriver(backgroundCtx, namespaceId, inputProxy.DiscoveryName)
	if err != nil && driver == nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateService fail. Get discovery driver fail. %s", err.Error())))
		return
	}

	newConf, formatAddr, variableList, err := driver.CheckInput(inputProxy.Config)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateService fail. Err: %s", err.Error())))
		return
	}
	inputProxy.Config = newConf

	input := &service_dto.ServiceInfo{
		Name:        serviceName,
		UUID:        inputProxy.UUID, //TODO 应该用rest里的uuid
		Desc:        inputProxy.Desc,
		Scheme:      inputProxy.Scheme,
		DiscoveryID: discoveryID,
		DriverName:  driverName,
		FormatAddr:  formatAddr,
		Config:      inputProxy.Config.String(),
		Timeout:     inputProxy.Timeout,
		Balance:     inputProxy.Balance,
	}

	err = s.service.UpdateService(backgroundCtx, namespaceId, operator, input, variableList)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateService fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// del 删除服务信息
func (s *serviceController) del(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	serviceName := ginCtx.Query("name")
	if serviceName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteService Info fail. err: serviceName can't be nil")))
		return
	}
	userId := controller.GetUserId(ginCtx)

	err := s.service.DeleteService(ginCtx, namespaceId, userId, serviceName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteService fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (s *serviceController) getEnum(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	searchName := ginCtx.Query("name")

	serviceList, err := s.service.GetServiceEnum(ginCtx, namespaceId, searchName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetServiceEnum fail. err:%s", err.Error())))
		return
	}

	data := make(map[string]interface{})
	data["list"] = serviceList
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

// online 上线
func (s *serviceController) online(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	serviceName := ginCtx.Param("service_name")
	input := &online_dto.UpdateOnlineStatusInput{}
	operator := controller.GetUserId(ginCtx)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	router, err := s.service.OnlineService(ginCtx, namespaceId, operator, serviceName, input.ClusterName)
	if err != nil && router == nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	} else if err == nil {
		ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
		return
	}

	m := common.Map[string, interface{}]{}
	m["router"] = router
	ginCtx.JSON(http.StatusOK, dto.Result{
		Code: -1,
		Data: m,
		Msg:  err.Error(),
	})

}

// online 下线
func (s *serviceController) offline(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	serviceName := ginCtx.Param("service_name")
	input := &online_dto.UpdateOnlineStatusInput{}
	operator := controller.GetUserId(ginCtx)

	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	if err := s.service.OfflineService(ginCtx, namespaceId, operator, serviceName, input.ClusterName); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

// getOnlineList 上线管理列表
func (s *serviceController) getOnlineList(ginCtx *gin.Context) {
	namespaceId := controller.GetNamespaceId(ginCtx)
	serviceName := ginCtx.Param("service_name")

	list, err := s.service.OnlineList(ginCtx, namespaceId, serviceName)
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

//// getApi 获取api接口信息
//func (s *serviceController) getApi(ginCtx *gin.Context) {
//	_ = controller.GetNamespaceId(ginCtx)
//	_ = ginCtx.Param("service_name")
//	uuid := ginCtx.Query("uuid")
//	api, err := s.service.GetApi(ginCtx, uuid)
//	if err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//
//	m := common.Map[string, interface{}]{}
//	m["api"] = &dto.ServiceApiOut{
//		Method:     dto.ServiceApiMethod(api.Method),
//		Uri:        api.Uri,
//		Id:         api.Uuid,
//		Name:       api.Name,
//		IsDelete:   api.IsDelete,
//		Operator:   api.Operator,
//		UpdateTime: common.TimeToStr(api.UpdateTime),
//		Config:     api.Config,
//	}
//	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
//}
//
//// getApis 目录下下的api接口信息
//func (s *serviceController) getApis(ginCtx *gin.Context) {
//	namespaceId := controller.GetNamespaceId(ginCtx)
//	serviceName := ginCtx.Param("service_name")
//	directoryUuid := ginCtx.Param("group_uuid")
//	apis, err := s.service.GetApiListByGroupUUID(ginCtx, namespaceId, serviceName, directoryUuid)
//	if err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//
//	list := make([]*dto.ServiceApiOut, 0, len(apis))
//	for _, api := range apis {
//		list = append(list, &dto.ServiceApiOut{
//			Method:     dto.ServiceApiMethod(api.Method),
//			Uri:        api.Uri,
//			Id:         api.Uuid,
//			Name:       api.Name,
//			IsDelete:   api.IsDelete,
//			Operator:   api.Operator,
//			UpdateTime: common.TimeToStr(api.UpdateTime),
//			Config:     api.Config,
//		})
//	}
//	m := common.Map[string, interface{}]{}
//	m["apis"] = list
//	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(m))
//}
//
//// postApi 新增api接口
//func (s *serviceController) postApi(ginCtx *gin.Context) {
//	namespaceId := controller.GetNamespaceId(ginCtx)
//	serviceName := ginCtx.Param("service_name")
//	operator := controller.GetUserId(ginCtx)
//	input := new(dto.ServiceApiInput)
//	if err := ginCtx.BindJSON(input); err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//	if err := s.service.CreateApi(ginCtx, namespaceId, operator, serviceName, input); err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
//}
//
//// putApi 修改api接口
//func (s *serviceController) putApi(ginCtx *gin.Context) {
//	namespaceId := controller.GetNamespaceId(ginCtx)
//	serviceName := ginCtx.Param("service_name")
//	uuid := ginCtx.Query("api_uuid")
//	operator := controller.GetUserId(ginCtx)
//
//	input := new(dto.ServiceApiInput)
//	if err := ginCtx.BindJSON(input); err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//	if err := s.service.UpdateApi(ginCtx, namespaceId, operator, serviceName, uuid, input); err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
//}
//
//// delApi 删除api接口
//func (s *serviceController) delApi(ginCtx *gin.Context) {
//	namespaceId := controller.GetNamespaceId(ginCtx)
//	serviceName := ginCtx.Param("service_name")
//	input := new(dto.DeleteServiceApiInput)
//	if err := ginCtx.BindJSON(input); err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//	if err := s.service.DeleteApi(ginCtx, namespaceId, serviceName, input.Ids); err != nil {
//		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
//		return
//	}
//	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
//
//}
