package main

import (
	"fmt"
	"os"

	apinto "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/activity-log/sqlite"
	apintoClient "github.com/eolinker/apinto-dashboard/internal/apinto"
	"github.com/eolinker/apinto-dashboard/internal/security"
	activity_log "github.com/eolinker/apinto-dashboard/modules/activity-log"

	"net/http"
	"strings"

	"github.com/eolinker/apinto-dashboard/modules/extenders"
	"github.com/eolinker/apinto-dashboard/modules/monitors"
	"github.com/eolinker/apinto-dashboard/modules/plugins"
	"github.com/eolinker/apinto-dashboard/modules/routers"
	"github.com/eolinker/eosc/log"
)

func init() {
	apinto.RetTemplate("tpl", "index", "icons")
}
func main() {
	// TODO: 日志设置
	transport := log.NewTransport(os.Stderr, log.DebugLevel)
	transport.SetFormatter(&log.LineFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: nil,
	})
	log.Reset(transport)
	log.SetPrefix("[dashboard]")
	cf, err := ReadConfig("config.yml")
	if err != nil {
		log.Panic(err)
		return
	}
	detailsService := security.NewUserDetailsService()
	err = InitUserDetails(detailsService, cf.UserDetails)
	if err != nil {
		log.Panic(err)
		return
	}
	activityHandler, err := sqlite.NewActivityDao("data/activity-log.db")
	if err != nil {
		log.Panic(err)
		return
	}

	apintoClient.Init(cf.Apinto)
	config := new(apinto.Config)

	config.DefaultZone = apinto.ZoneName(strings.ToLower(cf.Zone))

	apinto.SetActivityLogAddHandler(activityHandler, cf.FilterForwarded)
	config.UserDetailsService = detailsService

	monitorsModule := monitors.NewMonitor("monitors")
	config.Modules = append(config.Modules, &apinto.Module{
		Path:    "/monitors",
		Handler: monitorsModule,
		Name:    "monitors",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn: "监控",
			apinto.EnUs: "monitors",
		},
	})
	routersModule := routers.NewRouters("routers")
	config.Modules = append(config.Modules, &apinto.Module{
		Path:    "/routers/list",
		Handler: routersModule,
		Name:    "routers",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn: "路由",
			apinto.EnUs: "Ruters",
		},
	}, &apinto.Module{
		Handler: routersModule,
		Path:    "/profession/routers/",
		NotView: true,
	}, &apinto.Module{
		NotView: true,
		Handler: routersModule,
		Path:    "/skill/routers",
	})
	ms := toModule(cf)
	config.Modules = append(config.Modules, ms...)

	pluginModule := plugins.NewPlugins("plugins")
	config.Modules = append(config.Modules, &apinto.Module{
		Path:    "/plugins",
		Handler: pluginModule,
		Name:    "plugins",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn: "全局插件",
			apinto.EnUs: "Global Plugins",
		},
	})

	activityLogModule, err := activity_log.NewActivityLog("activity-log", activityHandler)
	if err != nil {
		log.Panic(err)
		return
	}

	config.Modules = append(config.Modules, &apinto.Module{
		Path:    "/activity-log",
		Handler: activityLogModule,
		Name:    "activity-log",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn: "操作日志",
			apinto.EnUs: "Activity log",
		},
	})
	extendersModule := extenders.NewExtenders("extenders")
	config.Modules = append(config.Modules, &apinto.Module{
		Path:    "/extenders",
		Handler: extendersModule,
		Name:    "extenders",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn: "扩展管理",
			apinto.EnUs: "extenders manager",
		},
	})

	config.Statics = map[string]string{
		"":    "./static",
		"js":  "./static/js",
		"umd": "./static/umd",
		//"css":"./static/css",
	}
	service, err := apinto.Create(config)

	if err != nil {
		log.Panic(err)
		return
	}
	err = http.ListenAndServe(fmt.Sprintf(":%s", cf.Port), service)
	if err != nil {
		log.Panic(err)
		return
	}
}
