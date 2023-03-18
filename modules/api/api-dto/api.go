package api_dto

import (
	"github.com/eolinker/apinto-dashboard/enum"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
)

type APIListItem struct {
	GroupUUID   string   `json:"group_uuid"`
	APIUUID     string   `json:"uuid"`
	APIName     string   `json:"name"`
	Method      []string `json:"method"`
	ServiceName string   `json:"service"`
	RequestPath string   `json:"request_path"`
	Source      string   `json:"source"`
	UpdateTime  string   `json:"update_time"`
	IsDelete    bool     `json:"is_delete"`
}

type APIEnum struct {
	ApiId   string `json:"api_id"`
	APIName string `json:"name"`
}

type APIInfo struct {
	ApiName          string                   `json:"name"`
	UUID             string                   `json:"uuid"`
	GroupUUID        string                   `json:"group_uuid"`
	Desc             string                   `json:"desc"`
	Driver           string                   `json:"driver"`
	RequestPath      string                   `json:"request_path"`
	RequestPathLabel string                   `json:"-"` //前端不传这个，后端存字段会使用
	ServiceName      string                   `json:"service"`
	Method           []string                 `json:"method"`
	ProxyPath        string                   `json:"proxy_path"`
	Timeout          int                      `json:"timeout"`
	Retry            int                      `json:"retry"`
	EnableWebsocket  bool                     `json:"enable_websocket"`
	Match            []*api_entry.MatchConf   `json:"match"`
	Header           []*api_entry.ProxyHeader `json:"proxy_header"`
	TemplateUUID     string                   `json:"template_uuid"`
}

type MatchConf struct {
	Position  string `json:"position"`
	MatchType string `json:"match_type"`
	Key       string `json:"key"`
	Pattern   string `json:"pattern"`
}

type ProxyHeader struct {
	OptType string `json:"opt_type"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

type ApiBatchOnlineCheckListItem struct {
	ServiceName string                 `json:"service"`
	ClusterEnv  string                 `json:"cluster"`
	Status      bool                   `json:"status"`
	Result      string                 `json:"result,omitempty"`
	Solution    *frontend_model.Router `json:"solution,omitempty"`
}

type ApiBatchCheckListItem struct {
	ApiName    string `json:"api"`
	ClusterEnv string `json:"cluster"`
	Status     bool   `json:"status"`
	Result     string `json:"result,omitempty"`
}

type ApiBatchInput struct {
	ApiUUIDs     []string `json:"api_uuids"`
	ClusterNames []string `json:"cluster_names"`
	OnlineToken  string   `json:"online_token"`
}

type ImportAPIListItem struct {
	Id     int                   `json:"id,omitempty"`
	Name   string                `json:"name"`
	Method string                `json:"method"`
	Path   string                `json:"path"`
	Desc   string                `json:"desc"`
	Status enum.ImportStatusType `json:"status"` //1正常 2冲突
}

// ImportAPIInfos 控制台导入API接口所需要的信息
type ImportAPIInfos struct {
	Apis  []ImportAPIInfoItem `json:"apis"`
	Token string              `json:"token"`
}

type ImportAPIInfoItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}
