package dto

import (
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/model"
)

type StrategyListOut struct {
	UUID       string                    `json:"uuid,omitempty"`
	Name       string                    `json:"name,omitempty"`
	Priority   int                       `json:"priority,omitempty"`
	IsStop     bool                      `json:"is_stop"`
	IsDelete   bool                      `json:"is_deleted"`
	Status     enum.StrategyOnlineStatus `json:"status"`
	Filters    string                    `json:"filters,omitempty"`
	Conf       string                    `json:"conf,omitempty"`
	Operator   string                    `json:"operator,omitempty"`
	UpdateTime string                    `json:"update_time,omitempty"`
}

type StrategyToPublishListOut struct {
	Name     string                    `json:"name"`
	Priority int                       `json:"priority"`
	Status   enum.StrategyOnlineStatus `json:"status"`
	OptTime  string                    `json:"opt_time"`
}

type StrategyPublish struct {
	VersionName string `json:"version_name"`
	Desc        string `json:"desc"`
	Source      string `json:"source"`
}

type StrategyStatusInput struct {
	IsStop bool `json:"is_stop"`
}

type StrategyInfoInput[T any] struct {
	Name     string         `json:"name"`
	Uuid     string         `json:"uuid"`
	Desc     string         `json:"desc"`
	Priority int            `json:"priority"`
	Filters  []*FilterInput `json:"filters"`
	Config   *T             `json:"config"`
}

type FilterInput struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

type StrategyInfoOutput[T any] struct {
	Name     string                `json:"name"`
	UUID     string                `json:"uuid"`
	Desc     string                `json:"desc"`
	Priority int                   `json:"priority"`
	Filters  []*model.FilterOutput `json:"filters,omitempty"`
	Config   *T                    `json:"config"`
}

//type StrategyInfoOutput struct {
//	Name     string                `json:"name"`
//	UUID     string                `json:"uuid"`
//	Desc     string                `json:"desc"`
//	Priority int                   `json:"priority"`
//	Filters  []*model.FilterOutput `json:"filters,omitempty"`
//	Config   *model.LimitConf      `json:"config"`
//}

type FilterOptionsItem struct {
	Name    string   `json:"name"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Pattern string   `json:"pattern"`
	Options []string `json:"options,omitempty"`
}

type MetricsOptionsItem struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

type FilterRemoteOutput struct {
	Target       string                      `json:"target"`
	Titles       []*model.RemoteTitles       `json:"titles"`
	Apis         []*model.RemoteApis         `json:"apis"`
	Services     []*model.RemoteServices     `json:"services"`
	Applications []*model.RemoteApplications `json:"applications"`
	Total        int                         `json:"total"`
}

type StrategyPublishHistory struct {
	Id         int                              `json:"id"`
	Name       string                           `json:"name"`     //版本名
	OptType    int                              `json:"opt_type"` //1.发布 2.回滚
	Operator   string                           `json:"operator"`
	CreateTime string                           `json:"create_time"`
	Details    []*StrategyPublishHistoryDetails `json:"details"`
}

type StrategyPublishHistoryDetails struct {
	Name       string                    `json:"name"`
	Priority   int                       `json:"priority"`
	Status     enum.StrategyOnlineStatus `json:"status"`
	CreateTime string                    `json:"opt_time"`
}
