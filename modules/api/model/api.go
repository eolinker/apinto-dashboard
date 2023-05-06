package apimodel

import (
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	"time"
)

type APIListItem struct {
	GroupUUID   string
	APIUUID     string
	APIName     string
	Method      []string
	ServiceName string
	RequestPath string
	Source      string
	UpdateTime  time.Time
	IsDelete    bool
	Scheme      string
	Publish     []*APIListItemPublish
	IsDisable   bool
}

type APIListItemPublish struct {
	Name   string
	Status int
}

type APIInfo struct {
	*api_entry.API
	Method []string
}

type APIVersionInfo struct {
	Api     *api_entry.API
	Version *api_entry.APIVersion
}

type BatchOnlineCheckListItem struct {
	ServiceTemplate string
	ClusterEnv      string
	Status          bool
	Result          string
	Solution        *frontend_model.Router
}

type BatchListItem struct {
	APIName    string
	ClusterEnv string
	Status     bool
	Result     string
}

type APIOnlineListItem struct {
	ClusterName string
	ClusterEnv  string
	Status      int
	Operator    string
	UpdateTime  string
}

// SourceListItem 来源列表项
type SourceListItem struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// ImportAPIListItem 导入API的检测列表项
type ImportAPIListItem struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Method string `json:"method"`
	Path   string `json:"path"`
	Desc   string `json:"desc"`
	Status int    `json:"status"` //1正常 2冲突 3无效路径
}

// ImportAPIRedisData 导入API检测时临时存入redis的数据结构
type ImportAPIRedisData struct {
	Apis        []*ImportAPIRedisDataItem `json:"apis"`
	ServiceName string                    `json:"service_name"`
	GroupID     string                    `json:"group_id"`
}

type ImportAPIRedisDataItem struct {
	ID  int      `json:"id"`
	Api *APIInfo `json:"api"`
}

// BatchOnlineCheckTask 批量上线api的临时数据，用于存入redis
type BatchOnlineCheckTask struct {
	Id       int    `json:"id"`
	Operator int    `json:"operator"`
	Token    string `json:"token"`
	Data     []byte `json:"data"`
}

// APIBatchConf BatchOnlineCheckTask.Data的数据结构
type APIBatchConf struct {
	ApiUUIDs     []string `json:"api_uuids"`
	ClusterNames []string `json:"cluster_names"`
}

// ApiCluster 集群信息
type ApiCluster struct {
	Name       string
	Title      string
	Env        string
	Status     int
	Updater    string
	UpdateTime string
}
