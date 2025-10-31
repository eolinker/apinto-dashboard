package apinto_module

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/pm3"
	"sync"
)

var (
	ErrorDriverNameConflict = errors.New("driver conflict")
	ErrorDriverNotExist     = errors.New("driver not exists")
	ErrorModuleNameConflict = errors.New("module conflict")
	ErrorRouterConflict     = errors.New("router conflict")
	defaultDrivers          = NewTDrivers()
)

type Config struct {
	Server  string
	Headers map[string]string
}
type Drivers interface {
	GetDriver(name string) (Driver, bool)
	Register(name string, driver Driver) error
}
type Driver interface {
	//CreatePlugin(define interface{}) (Plugin, error)\

	Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error)
	Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error)
}

type AccessInfo interface {
	Name() string
	CName() string
	Dependencies() []string
}

type tDrivers struct {
	lock    sync.RWMutex
	drivers map[string]Driver
}

func NewTDrivers() Drivers {
	return &tDrivers{
		lock:    sync.RWMutex{},
		drivers: make(map[string]Driver),
	}
}

func (t *tDrivers) GetDriver(name string) (Driver, bool) {
	t.lock.RLock()
	d, h := t.drivers[name]
	t.lock.RUnlock()
	return d, h
}

func (t *tDrivers) Register(name string, driver Driver) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	_, h := t.drivers[name]
	if h {
		return fmt.Errorf("%w：%s", ErrorDriverNameConflict, name)
	}
	t.drivers[name] = driver
	return nil
}
func GetDriver(name string) (Driver, bool) {
	return defaultDrivers.GetDriver(name)
}

func Register(name string, driver Driver) error {
	return defaultDrivers.Register(name, driver)
}
