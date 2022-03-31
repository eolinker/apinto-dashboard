package main

import (
	"fmt"
	apinto "github.com/eolinker/apinto-dashboard"
	"github.com/eolinker/apinto-dashboard/internal/security"
	"net/http"
)

func main() {
	config:=new(apinto.Config)
	config.DefaultZone = apinto.ZhCn

	detailsService := security.NewUserDetailsService()
	detailsService.Add(security.NewUserDetails("admin","admin", map[string]interface{}{}))
	config.UserDetailsService =detailsService

	service,err := apinto.Create(config)

	if err!= nil{
		fmt.Println(err)
		return
	}
	http.ListenAndServe(":8080",service)
}
