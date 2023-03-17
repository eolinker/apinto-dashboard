package main

import (
	_ "github.com/eolinker/apinto-dashboard/modules/api/service"
	_ "github.com/eolinker/apinto-dashboard/modules/application/application-service"
	_ "github.com/eolinker/apinto-dashboard/modules/audit/audit-service"
	_ "github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	_ "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-service"
	_ "github.com/eolinker/apinto-dashboard/modules/discovery/discovery-serivce"
	_ "github.com/eolinker/apinto-dashboard/modules/group/group-service"
	_ "github.com/eolinker/apinto-dashboard/modules/namespace/namespace-service"
	_ "github.com/eolinker/apinto-dashboard/modules/notice/notice-service"
	_ "github.com/eolinker/apinto-dashboard/modules/openapi/openapi-service"
	_ "github.com/eolinker/apinto-dashboard/modules/openapp/openapp-service"
	_ "github.com/eolinker/apinto-dashboard/modules/random/random-service"
	_ "github.com/eolinker/apinto-dashboard/modules/strategy/strategy-service"
	_ "github.com/eolinker/apinto-dashboard/modules/upstream/service"
	_ "github.com/eolinker/apinto-dashboard/modules/user/user-service"
	_ "github.com/eolinker/apinto-dashboard/modules/variable/variable-service"

	//online-service必须放到最后导入
	_ "github.com/eolinker/apinto-dashboard/modules/online/online-service"
)
