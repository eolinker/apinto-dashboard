package driver

import (
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/api/api-dto"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
)

type apiHTTP struct {
	apintoDriverName string
}

func (a *apiHTTP) CheckInput(input *api_dto.APIInfo) error {
	for _, m := range input.Method {
		switch m {
		case enum.MethodGET, enum.MethodPOST, enum.MethodPUT, enum.MethodDELETE, enum.MethodPATCH, enum.MethodHEAD, enum.MethodOPTIONS:
		default:
			return fmt.Errorf("method %s is illegal. ", m)
		}
	}

	return checkInput(input)
}

func (a *apiHTTP) ToApinto(name, desc string, disable bool, method []string, requestPath, requestPathLabel, proxyPath, serviceName string, timeout, retry int, hosts []string, match []*api_entry.MatchConf, header []*api_entry.ProxyHeader, templateUUID string) *v1.RouterConfig {
	router := toApinto(name, desc, disable, method, requestPath, requestPathLabel, proxyPath, serviceName, timeout, retry, hosts, match, header, templateUUID)
	router.Append["websocket"] = false
	router.Driver = a.apintoDriverName
	return router
}

func CreateAPIHttp(apintoDriverName string) api.IAPIDriver {
	return &apiHTTP{apintoDriverName: apintoDriverName}
}
