package local

import (
	"encoding/json"
	"fmt"
	apinto_module "github.com/eolinker/apinto-module"
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

func (t *tPlugin) CreateModule(name string, config interface{}) (apinto_module.Module, error) {
	cv, err := t.checkConfig(name, config)
	if err != nil {
		return nil, err
	}

	p := NewProxyAPi(cv.Server, name, cv)
	module := &tModule{name: name}

	module.routersInfo = append(module.routersInfo, p.CreateHome(t.define.Router.Home)...)
	for _, html := range t.define.Router.Html {
		module.routersInfo = append(module.routersInfo, p.CreateHtml(html.Path, html.Label))
	}
	for _, a := range t.define.Router.Frontend {
		path := fmt.Sprintf("/%s/", strings.Trim(a, "/"))
		module.routersInfo = append(module.routersInfo, p.CreateHtml(path, apinto_module.RouterLabelAssets))
	}
	for path, ms := range t.define.Router.Api {
		for method, att := range ms {

			module.routersInfo = append(module.routersInfo, p.CreateApi(fmt.Sprintf("%s.%s", method, path), method, path, att))

		}
	}
	for path, ms := range t.define.Router.OpenApi {
		for method, att := range ms {
			module.routersInfo = append(module.routersInfo, p.CreateOpenApi(fmt.Sprintf("%s.%s", method, path), method, path, att))
		}
	}
	for _, m := range t.define.Middleware {
		rules := make([][]string, 0, len(m.Rule))
		for _, rl := range m.Rule {
			rules = append(rules, strings.Split(rl, ","))
		}
		module.middlewareHandler = append(module.middlewareHandler, p.CreateMiddleware(m.Name, m.Path, m.Life, rules))
	}

	return module, nil
}

func (t *tPlugin) CheckConfig(name string, config interface{}) error {
	_, err := t.checkConfig(name, config)
	if err != nil {
		return err
	}
	return nil
}
func (t *tPlugin) checkConfig(name string, config interface{}) (*Config, error) {
	return DecodeConfig(config)
}

func (c *tPlugin) GetPluginFrontend(moduleName string) string {
	//TODO
	return ""
}

func (c *tPlugin) IsPluginVisible() bool {
	return true
}

func (c *tPlugin) IsShowServer() bool {
	return true
}

func (c *tPlugin) IsCanUninstall() bool {
	return true
}

func (c *tPlugin) IsCanDisable() bool {
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
