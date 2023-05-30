package strategy_model

import (
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/modules/strategy/strategy-entry"
	"time"
)

type Strategy struct {
	Strategy    *strategy_entry.Strategy
	Version     *strategy_entry.StrategyVersion
	Filters     string
	Conf        string
	OperatorStr string
	Status      int //这里只能用int 不能用 enum.StrategyOnlineStatus json解析会保报错 //1.待更新 2.已上线 3.待删除 4.未上线
}

type StrategyToPublish[T any] struct {
	Status          int //这里只能用int 不能用 enum.StrategyOnlineStatus json解析会保报错
	Strategy        *strategy_entry.Strategy
	StrategyVersion *strategy_entry.StrategyVersion
}

type StrategyPublishHistory struct {
	Id         int
	Name       string //版本名
	OptType    int    //1.发布 2.回滚
	Operator   string
	CreateTime time.Time
	Details    []*StrategyPublishHistoryDetails
}

type StrategyPublishHistoryDetails struct {
	Name     string
	Priority int
	Status   int
	OptTime  time.Time
}

type StrategyInfoOutput[T any] struct {
	*strategy_entry.Strategy
	Filters []*FilterOutput
	Config  *T
}

//type ExtenderData struct {
//	Api         map[string]*RemoteApis         `json:"api"`
//	Service     map[string]*RemoteServices     `json:"service"`
//	Application map[string]*RemoteApplications `json:"application"`
//}

type FilterOutput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	Type   string   `json:"type"`
	Label  string   `json:"label"`
	Title  string   `json:"title"`
}

type FilterOptionsItem struct {
	Name    string
	Title   string
	Type    string
	Pattern string
	Options []string
}
type FilterOptionsItems []*FilterOptionsItem

func (fs FilterOptionsItems) Len() int {
	return len(fs)
}

func (fs FilterOptionsItems) Less(i, j int) bool {
	return fs[i].Title < fs[j].Title
}

func (fs FilterOptionsItems) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

type MetricsOptionsItem struct {
	Name  string
	Title string
}

type FilterRemoteOutput struct {
	Target string
	Titles []apinto_module.OptionTitle
	Key    string
	List   []any
}

type RemoteTitles struct {
	Title string `json:"title"`
	Field string `json:"field"`
}

type RemoteApis struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
	//Service     string `json:"service"`
	Group       string `json:"group"`
	RequestPath string `json:"request_path"`
}

type RemoteServices struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Scheme string `json:"scheme"`
	Desc   string `json:"desc"`
}

type RemoteApplications struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
	Uuid string `json:"uuid"`
}

// VisitInfoOutputConf 访问策略信息输出配置
type VisitInfoOutputConf struct {
	VisitRule       string          `json:"visit_rule"`
	InfluenceSphere []*FilterOutput `json:"influence_sphere"`
	Continue        bool            `json:"continue"`
	//Extender        *ExtenderData   `json:"extender"`
}
