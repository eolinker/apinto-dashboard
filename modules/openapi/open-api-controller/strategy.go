package open_api_controller

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/controller/users"
	namespace_controller "github.com/eolinker/apinto-dashboard/modules/base/namespace-controller"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	cluster_model "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	open_strategy "github.com/eolinker/apinto-dashboard/modules/openapi/open-api-controller/strategy"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
	"github.com/go-basic/uuid"
	"net/http"
	"strconv"
	"strings"
)

func init() {
	RegisterRouter(newStrategyOpenController())
}
func newStrategyOpenController() *strategyOpenController {
	c := &strategyOpenController{profession: "strategy"}
	bean.Autowired(&c.dynamicService)
	bean.Autowired(&c.clusterService)
	return c
}

type strategyOpenController struct {
	dynamicService dynamic.IDynamicService
	clusterService cluster.IClusterService
	profession     string
}

func (c *strategyOpenController) Name() string {
	return "strategy"
}

func (c *strategyOpenController) Apis() []pm3.Api {
	return []pm3.Api{
		{
			Method:      http.MethodGet,
			Path:        "/api2/strategys",
			HandlerFunc: c.strategyList,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodPost,
			Path:        "/api2/strategy",
			HandlerFunc: c.create,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodPut,
			Path:        "/api2/strategy",
			HandlerFunc: c.save,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodPost,
			Path:        "/api2/strategy/sync",
			HandlerFunc: c.sync,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodGet,
			Path:        "/api2/strategy",
			HandlerFunc: c.info,
			Authority:   pm3.Public,
		},
		{
			Method:      http.MethodDelete,
			Path:        "/api2/strategy",
			HandlerFunc: c.delete,
			Authority:   pm3.Public,
		},
	}
}

var validStrategy = map[string]struct{}{
	"limiting": {},
	"visit":    {},
}

type SaveStrategy struct {
	*BasicInfo
	Append map[string]interface{} `json:"append"`
}

type BasicInfo struct {
	Uuid         string    `json:"uuid"`
	Name         string    `json:"name"`
	Desc         string    `json:"desc"`
	StrategyType string    `json:"strategy_type"`
	Filters      []*Filter `json:"filters"`
}

type Filter struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

//var existKeys = []string{
//	"uuid",
//	"name",
//	"desc",
//	"strategy_type",
//	"filters",
//}

func (r *SaveStrategy) UnmarshalJSON(bytes []byte) error {
	basicInfo := new(BasicInfo)
	err := json.Unmarshal(bytes, basicInfo)
	if err != nil {
		return err
	}
	validator, has := open_strategy.GetValidator(basicInfo.StrategyType)
	if !has {
		return fmt.Errorf("invalid strategy_type: %s", basicInfo.StrategyType)
	}
	bytes, err = validator.Validate(bytes)
	if err != nil {
		return err
	}
	tmp := make(map[string]interface{})
	err = json.Unmarshal(bytes, &tmp)
	if err != nil {
		return err
	}
	r.BasicInfo = basicInfo
	//for _, key := range existKeys {
	//	delete(tmp, key)
	//}
	filters := make(map[string]interface{})
	for _, filter := range basicInfo.Filters {
		filters[filter.Name] = filter.Values
	}
	tmp["filters"] = filters
	tmp[basicInfo.StrategyType] = tmp["config"]

	delete(tmp, "config")
	r.Append = tmp
	return nil
}

func (c *strategyOpenController) create(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	var worker SaveStrategy
	err := ctx.BindJSON(&worker)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	if worker.StrategyType == "" {
		controller.ErrorJson(ctx, http.StatusOK, "strategy_type is required")
		return
	}
	_, ok := validStrategy[worker.StrategyType]
	if !ok {
		controller.ErrorJson(ctx, http.StatusOK, "invalid strategy_type")
		return
	}
	if worker.Uuid == "" {
		worker.Uuid = uuid.New()
	}
	module := fmt.Sprintf("strategy-%s", worker.StrategyType)
	body, _ := json.Marshal(worker.Append)

	err = c.dynamicService.Create(ctx, namespaceID, c.profession, module, worker.StrategyType, worker.Name, worker.Uuid, worker.StrategyType, worker.Desc, string(body), users.GetUserId(ctx))
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	data := common.Map{}
	info := common.Map{}
	info["id"] = worker.Uuid
	info["source_name"] = fmt.Sprintf("%s@%s", worker.Uuid, c.profession)
	data["info"] = info
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (c *strategyOpenController) save(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uid := ctx.Query("uuid")
	var worker SaveStrategy
	err := ctx.BindJSON(&worker)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	info, err := c.dynamicService.Info(ctx, namespaceID, c.profession, uid)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	module := fmt.Sprintf("strategy-%s", info.BasicInfo.Driver)
	body, _ := json.Marshal(worker.Append)
	err = c.dynamicService.Save(ctx, namespaceID, c.profession, module, worker.Name, uid, worker.Desc, string(body), users.GetUserId(ctx))
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (c *strategyOpenController) delete(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uid := ctx.Query("uuid")
	info, err := c.dynamicService.Info(ctx, namespaceID, c.profession, uid)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	module := fmt.Sprintf("strategy-%s", info.BasicInfo.Driver)
	err = c.dynamicService.Delete(ctx, namespaceID, c.profession, module, uid)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

type SyncCluster struct {
	Status   string `json:"status"`
	Cluster  string `json:"cluster"`
	Strategy string `json:"strategy"`
}

func (c *strategyOpenController) sync(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	var tmp SyncCluster
	err := ctx.BindJSON(&tmp)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	info, err := c.dynamicService.Info(ctx, namespaceID, c.profession, tmp.Strategy)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	module := fmt.Sprintf("strategy-%s", info.BasicInfo.Driver)
	userId := users.GetUserId(ctx)
	var failClusters []string
	switch tmp.Status {
	case "online":
		_, failClusters, err = c.dynamicService.Online(ctx, namespaceID, c.profession, module, tmp.Strategy, []string{tmp.Cluster}, userId)
		if err != nil {
			controller.ErrorJson(ctx, http.StatusOK, err.Error())
			return
		}

	case "offline":
		_, failClusters, err = c.dynamicService.Offline(ctx, namespaceID, c.profession, module, tmp.Strategy, []string{tmp.Cluster}, userId)
		if err != nil {
			controller.ErrorJson(ctx, http.StatusOK, err.Error())
			return
		}
	default:
		controller.ErrorJson(ctx, http.StatusOK, "invalid status")
		return
	}
	if len(failClusters) > 0 {
		controller.ErrorJson(ctx, http.StatusOK, strings.Join(failClusters, "\n"))
		return
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(nil))
}

func (c *strategyOpenController) info(ctx *gin.Context) {
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	uid := ctx.Query("uuid")
	info, err := c.dynamicService.Info(ctx, namespaceID, c.profession, uid)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	data := map[string]interface{}{
		"uuid":          info.BasicInfo.ID,
		"name":          info.BasicInfo.Title,
		"desc":          info.BasicInfo.Description,
		"strategy_type": info.BasicInfo.Driver,
	}
	for key, value := range info.Append {
		data[key] = value
	}
	ctx.JSON(http.StatusOK, controller.NewSuccessResult(data))
}

func (c *strategyOpenController) strategyList(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	strategyType := ctx.Query("strategy_type")
	namespaceID := namespace_controller.GetNamespaceId(ctx)
	clusterNames := ctx.Query("cluster")
	names := make([]string, 0)
	if clusterNames != "" {
		err := json.Unmarshal([]byte(clusterNames), &names)
		if err != nil {
			controller.ErrorJson(ctx, http.StatusOK, err.Error())
			return
		}
	}

	all := len(names) < 1
	var err error
	var cs []*cluster_model.Cluster
	if all {
		cs, err = c.clusterService.GetAllCluster(ctx)
	} else {
		cs, err = c.clusterService.GetByNames(ctx, namespaceID, names)
	}

	csMap := map[string]string{}
	for _, t := range cs {
		csMap[t.Name] = t.Title
	}
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}

	columns := []string{
		"id",
		"name",
		"description",
		"updater",
		"create_time",
		"update_time",
	}
	page, pageSize := c.getPage(ctx)
	drivers := []string{
		"limiting",
		"visit",
	}
	if strategyType != "" {
		drivers = []string{strategyType}
	}
	clusterInfo, err := c.dynamicService.ClusterStatuses(ctx, namespaceID, c.profession, names, drivers, keyword, page, pageSize)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}

	list, total, err := c.dynamicService.List(ctx, namespaceID, c.profession, columns, drivers, keyword, page, pageSize)
	if err != nil {
		controller.ErrorJson(ctx, http.StatusOK, err.Error())
		return
	}
	result := make([]*Strategy, 0, len(list))
	for _, l := range list {
		id := l["id"]
		clusterStatuses := make([]*ClusterStatus, 0, len(cs))
		ci, ok := clusterInfo[id]
		if ok {
			for key, status := range ci {
				clusterStatuses = append(clusterStatuses, &ClusterStatus{
					Id:     strings.TrimPrefix(key, "cluster_"),
					Name:   key,
					Title:  csMap[key],
					Status: status,
				})
			}
		} else {
			for _, t := range cs {
				clusterStatuses = append(clusterStatuses, &ClusterStatus{
					Id:     strings.TrimPrefix(t.Name, "cluster_"),
					Name:   t.Name,
					Title:  t.Title,
					Status: "未发布",
				})
			}
		}
		result = append(result, &Strategy{
			Uuid:            l["id"],
			Name:            l["title"],
			Desc:            l["description"],
			CreateTime:      l["create_time"],
			UpdateTime:      l["update_time"],
			ClusterStatuses: clusterStatuses,
		})
	}

	ctx.JSON(http.StatusOK, controller.NewSuccessResult(map[string]interface{}{
		"strategys": result,
		"total":     total,
	}))
	return
}

func (c *strategyOpenController) getPage(ctx *gin.Context) (int, int) {
	page := ctx.Query("page")
	pageSize := ctx.Query("page_size")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	pz, _ := strconv.Atoi(pageSize)
	if pz < 1 {
		pz = 1000
	}
	return p, pz
}

type Basic struct {
	Name  string   `json:"name"`
	Title string   `json:"title"`
	Attr  string   `json:"attr,omitempty"`
	Enum  []string `json:"enum,omitempty"`
}

type Strategy struct {
	Uuid            string           `json:"uuid"`
	Name            string           `json:"name"`
	Desc            string           `json:"desc"`
	CreateTime      string           `json:"create_time"`
	UpdateTime      string           `json:"update_time"`
	ClusterStatuses []*ClusterStatus `json:"cluster_statuses"`
}

type ClusterStatus struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Title  string `json:"title"`
	Status string `json:"status"`
}
