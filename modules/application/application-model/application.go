package application_model

import (
	"github.com/eolinker/apinto-dashboard/modules/application/application-entry"
	"time"
)

const anonymousName = "匿名应用"

type ApplicationList []*Application

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

type Application struct {
	*application_entry.Application
	OperatorName string
	IsDelete     bool
	CustomAttr   []ApplicationCustomAttr
	ExtraParam   []ApplicationExtraParam
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

type ApplicationOnline struct {
	ClusterID   int
	ClusterName string
	Env         string
	Status      int //1.未上线 2.已下线 3.已上线  4.待更新
	Disable     bool
	Operator    string
	UpdateTime  time.Time
}

type ApplicationKeys struct {
	Key     string
	Values  []string
	KeyName string
}

type ApplicationVersion application_entry.ApplicationVersion
