package discovery_model

import (
	"github.com/eolinker/apinto-dashboard/modules/discovery/discovery-entry"
	"time"
)

type DiscoveryListItem struct {
	Name       string
	UUID       string
	Driver     string
	Desc       string
	UpdateTime time.Time
	IsDelete   bool
}

type DiscoveryInfo struct {
	Name   string
	UUID   string
	Driver string
	Desc   string
	Config []byte
	Render string
}

type DiscoveryEnum struct {
	Name   string
	Driver string
	Render string
}

type DiscoveryDriver struct {
	DriverName string
	Render     interface{}
}

type DiscoveryOnline struct {
	ClusterName string
	Env         string
	Status      int //1.未上线 2.已下线 3.已上线  4.待更新
	Operator    string
	UpdateTime  time.Time
}

type DiscoveryVersion discovery_entry.DiscoveryVersion
