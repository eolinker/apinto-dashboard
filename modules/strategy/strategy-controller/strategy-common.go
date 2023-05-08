package strategy_controller

import (
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/strategy"
	"github.com/eolinker/apinto-dashboard/modules/strategy/config"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-dto"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type strategyCommonController struct {
	strategyService strategy.IStrategyCommonService
}

func NewStrategyCommonController() *strategyCommonController {
	c := &strategyCommonController{}
	bean.Autowired(&c.strategyService)
	return c
}

func (s *strategyCommonController) FilterOptions(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

	options, err := s.strategyService.GetFilterOptions(ginCtx, namespaceId)
	if err != nil {
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	resList := make([]*strategy_dto.FilterOptionsItem, 0)
	for _, option := range options {
		resList = append(resList, &strategy_dto.FilterOptionsItem{
			Name:    option.Name,
			Title:   option.Title,
			Type:    option.Type,
			Pattern: option.Pattern,
			Options: option.Options,
		})
	}

	data := common.Map{}
	data["options"] = resList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (s *strategyCommonController) FilterRemote(ginCtx *gin.Context) {
	namespaceId := namespace_controller.GetNamespaceId(ginCtx)

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
		controller.ErrorJson(ginCtx, http.StatusOK, err.Error())
		return
	}

	res := strategy_dto.FilterRemoteOutput{
		Target:       remote.Target,
		Titles:       remote.Titles,
		Apis:         remote.Apis,
		Services:     remote.Services,
		Applications: remote.Applications,
		Total:        count,
	}

	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(res))
}

func (s *strategyCommonController) MetricsOptions(ginCtx *gin.Context) {
	metrics := config.StrategyMetrics()

	resList := make([]*strategy_dto.MetricsOptionsItem, 0)
	for _, option := range metrics {
		resList = append(resList, &strategy_dto.MetricsOptionsItem{
			Name:  option.Name(),
			Title: option.Title(),
		})
	}

	data := common.Map{}
	data["options"] = resList
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (s *strategyCommonController) ContentType(ginCtx *gin.Context) {

	items := config.GetContentTypeList()
	data := common.Map{}
	data["items"] = items
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (s *strategyCommonController) Charset(ginCtx *gin.Context) {

	items := config.GetStrategyCharsetList()
	data := common.Map{}
	data["items"] = items
	ginCtx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}
