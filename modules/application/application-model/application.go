package application_model

import (
	application_entry "github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"time"
)

const anonymousName = "匿名应用"

type ApplicationList []*ApplicationListItem

func (a ApplicationList) Len() int {
	return len(a)
}

func (a ApplicationList) Less(i, j int) bool {
	if a[i].Name == anonymousName {
		return true
	} else if a[j].Name == anonymousName {
		return false
	} else {
		return a[i].UpdateTime.After(a[j].UpdateTime)
	}
}

func (a ApplicationList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type ApplicationListItem struct {
	Uuid         string
	Name         string
	Desc         string
	UpdateTime   time.Time
	Operator     int
	OperatorName string
	IsDelete     bool
	Publish      []*APPListItemPublish
}

func (a *ApplicationListItem) UserId() int {
	return a.Operator

}

func (a *ApplicationListItem) Set(name string) {
	a.OperatorName = name
}

type ApplicationInfo struct {
	Name       string
	Uuid       string
	Desc       string
	CustomAttr []ApplicationCustomAttr
	Params     []ApplicationExtraParam
}

type ApplicationBasicInfo struct {
	Uuid       string
	Name       string
	Desc       string
	UpdateTime time.Time
}

type ApplicationRemoteOption struct {
	Uuid  string `json:"uuid,omitempty"`
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type ApplicationEntire struct {
	*application_entry.Application
}

type ApplicationBasicInfoList []*ApplicationBasicInfo

func (a ApplicationBasicInfoList) Len() int {
	return len(a)
}

func (a ApplicationBasicInfoList) Less(i, j int) bool {
	if a[i].Name == anonymousName {
		return true
	} else if a[j].Name == anonymousName {
		return false
	} else {
		return a[i].UpdateTime.After(a[j].UpdateTime)
	}
}

func (a ApplicationBasicInfoList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type APPListItemPublish struct {
	Name   string
	Title  string
	Status int
}

type ApplicationExtraParam struct {
	Key      string
	Value    string
	Conflict string
	Position string
}

type ApplicationCustomAttr struct {
	Key   string
	Value string
}

type ApplicationKeys struct {
	Key     string
	Values  []string
	KeyName string
}

// AppCluster 集群信息
type AppCluster struct {
	Name       string
	Title      string
	Env        string
	Status     int
	UpdaterId  int
	Updater    string
	UpdateTime string
}

func (a *AppCluster) UserId() int {

	return a.UpdaterId
}

func (a *AppCluster) Set(name string) {
	a.Updater = name
}
