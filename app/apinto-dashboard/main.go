package main

import (
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/security"
	"github.com/eolinker/apinto-dashboard/modules/profession"
	"net/http"
)

func main() {
	config:=new(apinto.Config)
	config.DefaultZone = apinto.ZhCn

	detailsService := security.NewUserDetailsService()
	detailsService.Add(security.NewUserDetails("admin","admin", map[string]interface{}{}))
	config.UserDetailsService =detailsService
	config.Modules = append(config.Modules, &apinto.Module{
		Path:     "/discovery/list",
		Icon:     "",
		Handler:  profession.NewProfession("discovery"),
		Name:     "discovery",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn:"服务发现",
			apinto.EnUs: "discovery service",
		},
	})
	config.Modules = append(config.Modules, &apinto.Module{
		Path:     "/routers/list",
		Icon:     "",
		Handler:  profession.NewProfession("routers"),
		Name:     "routers",
		I18nName: map[apinto.ZoneName]string{
			apinto.ZhCn:"路由",
			apinto.EnUs: "routers",
		},
	})
	config.Statics = map[string]string{
		"":"./static",
		"js":"./static/js",
	}
	service,err := apinto.Create(config)

	if err!= nil{
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080",service)
}
