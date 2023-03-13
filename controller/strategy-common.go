package controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	strategy_handler "github.com/eolinker/apinto-dashboard/service/strategy-handler"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type strategyCommonController struct {
	strategyService service.IStrategyCommonService
}

func RegisterStrategyCommonRouter(router gin.IRoutes) {
	c := &strategyCommonController{}
	bean.Autowired(&c.strategyService)

	router.GET("/strategy/filter-options", c.filterOptions)
	router.GET("/strategy/filter-remote/:name", c.filterRemote)
	router.GET("/strategy/metrics-options", c.metricsOptions)
	router.GET("/strategy/content-type", c.contentType)
	router.GET("/strategy/charset", c.charset)
}

func (s *strategyCommonController) filterOptions(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)

	options, err := s.strategyService.GetFilterOptions(ginCtx, namespaceId)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.FilterOptionsItem, 0)
	for _, option := range options {
		resList = append(resList, &dto.FilterOptionsItem{
			Name:    option.Name,
			Title:   option.Title,
			Type:    option.Type,
			Pattern: option.Pattern,
			Options: option.Options,
		})
	}

	data := common.Map[string, interface{}]{}
	data["options"] = resList
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (s *strategyCommonController) filterRemote(ginCtx *gin.Context) {
	namespaceId := getNamespaceId(ginCtx)

	name := ginCtx.Param("name")
	keyword := ginCtx.Query("keyword")
	groupUUID := ginCtx.Query("group_uuid")
	pageNumStr := ginCtx.Query("page_num")
	pageSizeStr := ginCtx.Query("page_size")
	pageNum, _ := strconv.Atoi(pageNumStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	pageNum = 1
	pageSize = 10000

	remote, count, err := s.strategyService.GetFilterRemote(ginCtx, namespaceId, name, keyword, groupUUID, pageNum, pageSize)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	res := dto.FilterRemoteOutput{
		Target:       remote.Target,
		Titles:       remote.Titles,
		Apis:         remote.Apis,
		Services:     remote.Services,
		Applications: remote.Applications,
		Total:        count,
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(res))
}

func (s *strategyCommonController) metricsOptions(ginCtx *gin.Context) {

	options, err := s.strategyService.GetMetricsOptions()
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.MetricsOptionsItem, 0)
	for _, option := range options {
		resList = append(resList, &dto.MetricsOptionsItem{
			Name:  option.Name,
			Title: option.Title,
		})
	}

	data := common.Map[string, interface{}]{}
	data["options"] = resList
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (s *strategyCommonController) contentType(ginCtx *gin.Context) {

	items := strategy_handler.GetContentTypeList()
	data := common.Map[string, interface{}]{}
	data["items"] = items
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (s *strategyCommonController) charset(ginCtx *gin.Context) {

	items := enum.GetStrategyCharsetList()
	data := common.Map[string, interface{}]{}
	data["items"] = items
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}
