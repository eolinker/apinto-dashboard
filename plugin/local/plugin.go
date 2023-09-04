package local

import (
	"encoding/json"
	"fmt"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	"strings"
)

type tPlugin struct {
	define *Define
}

func newTPlugin(define interface{}) (apinto_module.Plugin, error) {
	dv, err := DecodeDefine(define)
	if err != nil {
		return nil, err
	}
	return &tPlugin{define: dv}, nil
}

func (p *tPlugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	cv, err := p.checkConfig(name, config)
	if err != nil {
		return nil, err
	}

	proxy := NewProxyAPi(p.define.Cmd, name, cv)
	module := &tModule{name: name}

	module.routersInfo = append(module.routersInfo, proxy.CreateHome(p.define.Router.Home)...)
	for _, html := range p.define.Router.Html {
		module.routersInfo = append(module.routersInfo, proxy.CreateHtml(html.Path, html.Label))
	}
	for _, a := range p.define.Router.Frontend {
		path := fmt.Sprintf("/%s/", strings.Trim(a, "/"))
		module.routersInfo = append(module.routersInfo, proxy.CreateHtml(path, apinto_module.RouterLabelAssets))
	}
	for path, ms := range p.define.Router.Api {
		for method, att := range ms {

			module.routersInfo = append(module.routersInfo, proxy.CreateApi(fmt.Sprintf("%s.%s", method, path), method, path, att))

		}
	}
	for path, ms := range p.define.Router.OpenApi {
		for method, att := range ms {
			module.routersInfo = append(module.routersInfo, proxy.CreateOpenApi(fmt.Sprintf("%s.%s", method, path), method, path, att))
		}
	}
	for _, m := range p.define.Middleware {
		rules := make([][]string, 0, len(m.Rule))
		for _, rl := range m.Rule {
			rules = append(rules, strings.Split(rl, ","))
		}
		module.middlewareHandler = append(module.middlewareHandler, proxy.CreateMiddleware(m.Name, rules))
	}

	return module, nil
}

func (p *tPlugin) CheckConfig(name string, config interface{}) error {
	_, err := p.checkConfig(name, config)
	if err != nil {
		return err
	}
	return nil
}
func (p *tPlugin) checkConfig(name string, config interface{}) (*Config, error) {
	return DecodeConfig(config)
}

func (p *tPlugin) GetPluginFrontend(moduleName string) string {
	return fmt.Sprintf("module/%s", moduleName)
}

func (p *tPlugin) IsPluginVisible() bool {
	return true
}

func (p *tPlugin) IsShowServer() bool {
	return false
}

func (p *tPlugin) IsCanUninstall() bool {
	return true
}

func (p *tPlugin) IsCanDisable() bool {
	return true
}

func DecodeDefine(define interface{}) (*Define, error) {
	return apinto_module.DecodeFor[Define](define)
}
func DecodeConfig(config interface{}) (*Config, error) {
	return apinto_module.DecodeFor[Config](config)

}
func decode[T any](v interface{}) (*T, error) {
	var err error
	var data []byte
	switch v := v.(type) {
	case *string:
		data = []byte(*v)
	case *[]byte:
		data = *v
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(v)

	}
	if err != nil {
		return nil, err
	}
	pv := new(T)
	err = json.Unmarshal(data, pv)
	if err != nil {
		return nil, err
	}
	return pv, nil
}
