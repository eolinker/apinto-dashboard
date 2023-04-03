package discovery

import (
	"errors"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/driver"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
)

type IDiscoveryDriverManager interface {
	driver.IDriverManager[IDiscoveryDriver]
	List() []*driver.DriverInfo
}

type IDiscoveryDriver interface {
	Render() string
	CheckInput(config []byte) ([]byte, string, []string, error)
	FormatConfig(config []byte) []byte
	CheckConfIsChange(old []byte, latest []byte) bool
	ToApinto(namespace, name, desc string, config []byte) *v1.DiscoveryConfig
	OptionsEnum() upstream.IServiceDriver
}

var (
	ErrVariableIllegal = errors.New("Variable is illegal. ")
)
