package upstream_model

import (
	"github.com/eolinker/apinto-dashboard/entry"
	"time"
)

type ServiceListItem struct {
	Name        string
	UUID        string
	Scheme      string
	DiscoveryID int
	DriverName  string
	Config      string
	UpdateTime  time.Time
	IsDelete    bool
}

type ServiceInfo struct {
	*entry.ServiceVersion
	Name      string
	Desc      string
	ServiceId int
	UUID      string
}

type ServiceOnline struct {
	ClusterID   int
	ClusterName string
	Env         string
	Status      int //1.未上线 2.已下线 3.已上线  4.待更新
	Operator    string
	UpdateTime  time.Time
}
