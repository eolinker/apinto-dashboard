package main

import (
	_ "github.com/eolinker/apinto-dashboard/modules/api/service"
	_ "github.com/eolinker/apinto-dashboard/modules/application/application-service"
	_ "github.com/eolinker/apinto-dashboard/modules/audit/audit-service"
	_ "github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	_ "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-service"
	_ "github.com/eolinker/apinto-dashboard/modules/discovery/discovery-serivce"
	_ "github.com/eolinker/apinto-dashboard/modules/group/group-service"
)
