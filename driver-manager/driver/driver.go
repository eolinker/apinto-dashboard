package driver

import (
	"errors"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/modules/api/api-dto"
	apimodel "github.com/eolinker/apinto-dashboard/modules/api/model"
)

const (
	DriverStatic = "static"
	DriverConsul = "consul"
	DriverNacos  = "nacos"
	DriverEureka = "eureka"
	DriverBasic  = "basic"
	DriverApikey = "apikey"
	DriverAKsK   = "aksk"
	DriverJwt    = "jwt"

	DriverApiHTTP = "http"
)

var (
	ErrVariableIllegal = errors.New("Variable is illegal. ")
)

type IServiceDriver interface {
	Render() string
	CheckInput(config []byte) ([]byte, string, []string, error)
	FormatConfig(config []byte) []byte
	ToApinto(name, namespace, desc, scheme, balance, discoveryName, driverName string, timeout int, config []byte) *v1.ServiceConfig
}

type IDriver interface {
	Render() string
	CheckInput(config []byte) ([]byte, string, []string, error)
	FormatConfig(config []byte) []byte
	CheckConfIsChange(old []byte, latest []byte) bool
	ToApinto(namespace, name, desc string, config []byte) *v1.DiscoveryConfig
}

type IDiscoveryDriver interface {
	IDriver
	OptionsEnum() IServiceDriver
}

type IAuthDriver interface {
	Render() string
	CheckInput(config []byte) error
	//GetAuthListInfo 获取健全列表展示需要用的参数信息
	GetAuthListInfo(config []byte) string
	ToApinto(expire int64, position string, tokenName string, config []byte, hideCredential bool) v1.ApplicationAuth
}

type IAPIDriver interface {
	CheckInput(input *api_dto.APIInfo) error
	ToApinto(name, desc string, disable bool, method []string, requestPath, requestPathLabel, proxyPath, serviceName string, timeout, retry int, enableWebsocket bool, match []*api_entry.MatchConf, header []*api_entry.ProxyHeader) *v1.RouterConfig
}

// IAPISyncFormatDriver 同步api所需的文件格式驱动
type IAPISyncFormatDriver interface {
	FormatAPI(data []byte, namespaceID, appID int, groupID, prefix, label string) ([]*apimodel.APIInfo, error)
}

type ICLConfigDriver interface {
	CheckInput(config []byte) error
	ToApinto(name string, config []byte) interface{}
	FormatOut(operator string, config *entry.ClusterConfig) interface{}
	InitConfig(config []byte) error
}

// IMonitorSourceDriver 监控数据源驱动
type IMonitorSourceDriver interface {
	CheckInput(config []byte) ([]byte, error)
}
