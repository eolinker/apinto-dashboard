package main

import (
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/security"
	"log"
	"net/http"
	"strings"
)

func init() {
	apinto.RetTemplate("tpl","index","icons")
}
func main() {
	cf, err := ReadConfig("config.yml")
	if err != nil {
		log.Println("[Error]",err)
		return
	}

	config:=new(apinto.Config)

	config.DefaultZone =apinto.ZoneName(strings.ToLower(cf.Zone))

	detailsService := security.NewUserDetailsService()
	detailsService.Add(security.NewUserDetails("admin","admin", map[string]interface{}{}))
	config.UserDetailsService =detailsService

	ms:=toModule(cf)
	config.Modules = append(config.Modules, ms...)
	config.Statics = map[string]string{
		"":"./static",
		"js":"./static/js",
		//"css":"./static/css",
	}
	service,err := apinto.Create(config)

	if err!= nil{
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080",service)
}
