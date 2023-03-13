package controller

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

/*
每使用该泛型新增一个策略需要做一下操作：
1. entry层，定好该策略所需的配置结构体
2. 在service/strategy-handler内 为该策略新建一个实现了IStrategyHandler接口的处理器handler。
3. controller层新增一个文件 strategy-xxx.go指定该策略的路由，并且New一个策略处理器handler传入泛型StrategyService
4. 新增的策略service需要添加进IStrategyCommonService， 为了重置集群时将能将所有策略的配置重置到集群中。
5. 额外需要做的是添加一些enum枚举。
*/

type strategyController[T any, K any] struct {
	strategyService service.IStrategyService[T, K]
}

func newStrategyController[T any, K any](strategyService service.IStrategyService[T, K]) *strategyController[T, K] {
	return &strategyController[T, K]{strategyService: strategyService}
}

func (s *strategyController[T, K]) list(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")

	strategies, err := s.strategyService.GetList(ginCtx, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	resList := make([]*dto.StrategyListOut, 0, len(strategies))

	for _, strategy := range strategies {
		resList = append(resList, &dto.StrategyListOut{
			UUID:       strategy.Strategy.UUID,
			Name:       strategy.Strategy.Name,
			Priority:   strategy.Strategy.Priority,
			IsStop:     strategy.Strategy.IsStop,
			IsDelete:   strategy.Strategy.IsDelete,
			Status:     enum.StrategyOnlineStatus(strategy.Status),
			Filters:    strategy.Filters,
			Conf:       strategy.Conf,
			Operator:   strategy.OperatorStr,
			UpdateTime: common.TimeToStr(strategy.Strategy.UpdateTime),
		})
	}

	data := common.Map[string, interface{}]{}
	data["strategies"] = resList
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))

}

func (s *strategyController[T, K]) get(ginCtx *gin.Context) {
	namespaceID := GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")
	if uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetStrategyInfo fail. err: uuid can't be nil")))
		return
	}

	info, extender, err := s.strategyService.GetInfo(ginCtx, namespaceID, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("GetStrategyInfo fail. err:%s", err.Error())))
		return
	}

	strategy := &dto.StrategyInfoOutput[K]{
		Name:     info.Name,
		UUID:     info.UUID,
		Desc:     info.Desc,
		Priority: info.Priority,
		Filters:  info.Filters,
		Config:   info.Config,
	}

	data := make(map[string]interface{})
	data["strategy"] = strategy
	data["extender"] = extender

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (s *strategyController[T, K]) create(ginCtx *gin.Context) {
	namespaceID := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	if clusterName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateStrategy fail. err: clusterName can't be nil")))
		return
	}

	input := new(dto.StrategyInfoInput[T])
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	operator := GetUserId(ginCtx)

	//校验参数
	if err := s.strategyService.CheckInput(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateStrategy fail. err:%s", err.Error())))
		return
	}

	err := s.strategyService.CreateStrategy(ginCtx, namespaceID, operator, clusterName, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("CreateAPI fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (s *strategyController[T, K]) update(ginCtx *gin.Context) {
	namespaceID := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	if clusterName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateStrategy fail. err: clusterName can't be nil")))
		return
	}
	uuid := ginCtx.Query("uuid")
	if uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateStrategy fail. err: uuid can't be nil")))
		return
	}

	input := new(dto.StrategyInfoInput[T])
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}
	input.Uuid = uuid

	operator := GetUserId(ginCtx)

	//校验参数
	if err := s.strategyService.CheckInput(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateStrategy fail. err:%s", err.Error())))
		return
	}

	err := s.strategyService.UpdateStrategy(ginCtx, namespaceID, operator, clusterName, input)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("UpdateStrategy fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (s *strategyController[T, K]) del(ginCtx *gin.Context) {
	namespaceID := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	if clusterName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteStrategy fail. err: clusterName can't be nil")))
		return
	}
	uuid := ginCtx.Query("uuid")
	if uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteStrategy fail. err: uuid can't be nil")))
		return
	}

	userID := GetUserId(ginCtx)
	err := s.strategyService.DeleteStrategy(ginCtx, namespaceID, userID, clusterName, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("DeleteStrategy fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (s *strategyController[T, K]) restore(ginCtx *gin.Context) {
	namespaceID := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	if clusterName == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("RestoreStrategy fail. err: clusterName can't be nil")))
		return
	}
	uuid := ginCtx.Query("uuid")
	if uuid == "" {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("RestoreStrategy fail. err: uuid can't be nil")))
		return
	}

	userID := GetUserId(ginCtx)
	err := s.strategyService.RestoreStrategy(ginCtx, namespaceID, userID, clusterName, uuid)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(fmt.Sprintf("RestoreStrategy fail. err:%s", err.Error())))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (s *strategyController[T, K]) updateStop(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	uuid := ginCtx.Query("uuid")
	clusterName := ginCtx.Query("cluster_name")

	input := new(dto.StrategyStatusInput)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	userId := GetUserId(ginCtx)
	err := s.strategyService.UpdateStop(ginCtx, namespaceId, userId, uuid, clusterName, input.IsStop)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	data := common.Map[string, interface{}]{}
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))

}

func (s *strategyController[T, K]) toPublish(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
	list, err := s.strategyService.ToPublish(ginCtx, namespaceId, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.StrategyToPublishListOut, 0)

	for _, publish := range list {
		resList = append(resList, &dto.StrategyToPublishListOut{
			Name:     publish.Strategy.Name,
			Priority: publish.Strategy.Priority,
			Status:   enum.StrategyOnlineStatus(publish.Status),
			OptTime:  common.TimeToStr(publish.Strategy.UpdateTime),
		})
	}

	bytes, _ := json.Marshal(list)
	source := common.Base64Encode(bytes)

	data := common.Map[string, interface{}]{}
	data["is_publish"] = len(resList) > 0
	data["source"] = source
	data["strategies"] = resList
	data["version_name"] = time.Now().Format("20060102150405") + "-release"

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (s *strategyController[T, K]) publish(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")

	input := new(dto.StrategyPublish)
	if err := ginCtx.BindJSON(input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	userId := GetUserId(ginCtx)
	if err := s.strategyService.Publish(ginCtx, namespaceId, userId, clusterName, input); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}

func (s *strategyController[T, K]) publishHistory(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")
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
	list, total, err := s.strategyService.PublishHistory(ginCtx, namespaceId, pageNum, pageSize, clusterName)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	resList := make([]*dto.StrategyPublishHistory, 0, len(list))
	for _, history := range list {
		details := make([]*dto.StrategyPublishHistoryDetails, 0)
		for _, detail := range history.Details {
			details = append(details, &dto.StrategyPublishHistoryDetails{
				Name:       detail.Name,
				Priority:   detail.Priority,
				Status:     enum.StrategyOnlineStatus(detail.Status),
				CreateTime: common.TimeToStr(detail.OptTime),
			})
		}
		resList = append(resList, &dto.StrategyPublishHistory{
			Id:         history.Id,
			Name:       history.Name,
			OptType:    history.OptType,
			Operator:   history.Operator,
			CreateTime: common.TimeToStr(history.CreateTime),
			Details:    details,
		})
	}

	data := common.Map[string, interface{}]{}
	data["historys"] = resList
	data["total"] = total
	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(data))
}

func (s *strategyController[T, K]) changePriority(ginCtx *gin.Context) {
	namespaceId := GetNamespaceId(ginCtx)
	clusterName := ginCtx.Query("cluster_name")

	maps := common.Map[string, int]{}

	if err := ginCtx.BindJSON(&maps); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	priorityMap := common.Map[int, int]{}
	for _, priority := range maps {
		if priority == 0 {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult("优先级不可填空"))
			return
		}
		if priority > 999 {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult("优先级不可超过999"))
			return
		}
		priorityMap[priority] += 1
	}
	for _, v := range priorityMap {
		if v > 1 {
			ginCtx.JSON(http.StatusOK, dto.NewErrorResult("优先级不可重复"))
			return
		}
	}

	if len(maps) == 0 {
		ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
		return
	}

	userId := GetUserId(ginCtx)
	if err := s.strategyService.ChangePriority(ginCtx, namespaceId, userId, clusterName, maps); err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewErrorResult(err.Error()))
		return
	}

	ginCtx.JSON(http.StatusOK, dto.NewSuccessResult(nil))
}
