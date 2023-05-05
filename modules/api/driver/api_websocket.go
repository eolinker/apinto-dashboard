package driver

import (
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/modules/api"
	api_dto "github.com/eolinker/apinto-dashboard/modules/api/api-dto"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
)

type apiWebsocket struct {
	apintoDriverName string
}

func (a *apiWebsocket) CheckInput(input *api_dto.APIInfo) error {
	return checkInput(input)
}

func (a *apiWebsocket) ToApinto(name, desc string, disable bool, method []string, requestPath, requestPathLabel, proxyPath, serviceName string, timeout, retry int, hosts []string, match []*api_entry.MatchConf, header []*api_entry.ProxyHeader, templateUUID string) *v1.RouterConfig {

	router := toApinto(name, desc, disable, method, requestPath, requestPathLabel, proxyPath, serviceName, timeout, retry, hosts, match, header, templateUUID)
	router.Driver = a.apintoDriverName
	router.EnableWebsocket = true
	return router
}

func CreateAPIWebsocket(apintoDriverName string) api.IAPIDriver {
	return &apiWebsocket{apintoDriverName: apintoDriverName}
}
