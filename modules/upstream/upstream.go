package upstream

import "github.com/eolinker/apinto-dashboard/client/v1"

type IServiceDriver interface {
	Render() string
	CheckInput(config []byte) ([]byte, string, []string, error)
	FormatConfig(config []byte) []byte
	ToApinto(name, namespace, desc, scheme, balance, discoveryName, driverName string, timeout int, config []byte) *v1.ServiceConfig
}
