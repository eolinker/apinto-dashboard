package builder

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/log"
	"net"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestBuild(t *testing.T) {
	var demoProviderDriver = NewDemoProviderDriver()
	modules := make([]apinto_module.Module, 0, 0) //模块实例

	module, err := demoProviderDriver.Create(nil, nil)
	if err != nil {
		return
	}
	modules = append(modules, module)
	builder := NewModuleBuilder(gin.New())
	builder.Append(modules...)

	handler, providers, err := builder.Build()
	if err != nil {
		log.Error(err)
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

var (
	_ apinto_module.Driver = (*DemoProviderDriver)(nil)
	_ pm3.Module           = (*DemoProviderModule)(nil)
)

type DemoProviderDriver struct {
	providers *atomic.Pointer[apinto_module.IProviders]
}

func (d *DemoProviderDriver) Install(info *pm3.PluginDefine) (ms []pm3.PModule, acs []pm3.PAccess, fs []pm3.PFrontend, err error) {
	return pm3.ReadPluginAssembly(info)
}

func (d *DemoProviderDriver) Create(info *pm3.PluginDefine, config pm3.PluginConfig) (pm3.Module, error) {
	return &DemoProviderModule{ModuleTool: pm3.NewModuleTool(info.Id), name: info.Id}, nil
}

func NewDemoProviderDriver() *DemoProviderDriver {
	return &DemoProviderDriver{
		providers: new(atomic.Pointer[apinto_module.IProviders]),
	}
}

func (d *DemoProviderDriver) SetProviders(providers apinto_module.IProviders) {
	d.providers.Store(&providers)
}

type DemoProviderModule struct {
	*pm3.ModuleTool
	name string

	providers *atomic.Pointer[apinto_module.IProviders]
}

func (d *DemoProviderModule) Frontend() []pm3.FrontendAsset {
	return nil
}

func (d *DemoProviderModule) Middleware() []pm3.Middleware {
	return nil
}

func (d *DemoProviderModule) Apis() apinto_module.RoutersInfo {
	apis := apinto_module.RoutersInfo{{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("/api/%s/provider/:name", d.name),

		HandlerFunc: func(context *gin.Context) {

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
				result := make([]*pm3.CargoItem, 0, len(cargos))
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
	}
	d.InitAccess(apis)
	return apis
}

func (d *DemoProviderModule) Name() string {
	return d.name
}

func (d *DemoProviderModule) Support() (pm3.ProviderSupport, bool) {
	return nil, false
}
