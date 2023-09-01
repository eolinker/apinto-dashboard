package apinto_module

import (
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
)

type TestCore struct {
}

func (t *TestCore) RoutersInfo() RoutersInfo {
	return nil
}

func (t *TestCore) Name() string {
	return "core"
}

func (t *TestCore) Routers() (Routers, bool) {
	return t, true
}

func (t *TestCore) Middleware() (Middleware, bool) {
	return nil, false
}

func (t *TestCore) Support() (ProviderSupport, bool) {
	return nil, false
}

func TestBuild(t *testing.T) {
	demoProviderDriver := NewDemoProviderDriver()
	modules := make([]Module, 0, 0) //模块实例
	plugin, err := demoProviderDriver.CreatePlugin(nil)
	if err != nil {
		return
	}
	module, err := plugin.CreateModule("core", nil)
	if err != nil {
		return
	}
	modules = append(modules, module)
	builder := NewModuleBuilder(gin.New(), new(TestCore))
	builder.Append(modules...)

	handler, providers, err := builder.Build()
	if err != nil {

		return
	}
	demoProviderDriver.SetProviders(providers)

	server := http.Server{Handler: handler}
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		return
	}
	server.Serve(listen)
}

type DemoProviderDriver struct {
	providers *atomic.Pointer[IProviders]
}

func (d *DemoProviderDriver) GetPluginFrontend(moduleName string) string {
	return ""
}

func (d *DemoProviderDriver) IsPluginVisible() bool {
	return false
}

func (d *DemoProviderDriver) IsShowServer() bool {
	return false
}

func (d *DemoProviderDriver) IsCanUninstall() bool {
	return false
}

func (d *DemoProviderDriver) IsCanDisable() bool {
	return false
}

func NewDemoProviderDriver() *DemoProviderDriver {
	return &DemoProviderDriver{
		providers: new(atomic.Pointer[IProviders]),
	}
}

func (d *DemoProviderDriver) SetProviders(providers IProviders) {
	d.providers.Store(&providers)
}

func (d *DemoProviderDriver) CreateModule(name string, config interface{}) (Module, error) {
	return &DemoProviderModule{name: name}, nil
}

func (d *DemoProviderDriver) CheckConfig(name string, config interface{}) error {
	return nil
}

func (d *DemoProviderDriver) CreatePlugin(define interface{}) (Plugin, error) {
	return d, nil
}

type DemoProviderModule struct {
	name string

	providers *atomic.Pointer[IProviders]
}

func (d *DemoProviderModule) RoutersInfo() RoutersInfo {
	return RoutersInfo{{
		Method:  http.MethodGet,
		Path:    fmt.Sprintf("/api/%s/provider/:name", d.name),
		Handler: "",
		HandlerFunc: []HandlerFunc{func(context *gin.Context) {

			name := context.Param("name")

			providers := d.providers.Load()
			if providers != nil {
				provider, ok := (*providers).Provider(name)
				if !ok {
					context.JSON(200, struct {
						Code string `json:"code"`
						Msg  string `json:"msg"`
					}{
						"200", fmt.Sprintf("not support data for %s", name),
					})
				}
				cargos := provider.Provide(0)
				result := make([]*CargoItem, 0, len(cargos))
				for _, c := range cargos {
					result = append(result, c.Export())
				}
				context.JSON(200, map[string]interface{}{
					"code": "00000",
					"data": map[string]interface{}{
						"cargos": result,
					},
				})
			}
		}},
	}}
}

func (d *DemoProviderModule) Name() string {
	return d.name
}

func (d *DemoProviderModule) Access() []AccessInfo {
	return nil
}

func (d *DemoProviderModule) Routers() (Routers, bool) {
	return d, true
}

func (d *DemoProviderModule) Middleware() (Middleware, bool) {
	return nil, false
}

func (d *DemoProviderModule) Support() (ProviderSupport, bool) {
	return nil, false
}
