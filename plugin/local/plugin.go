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
	dv, err := decode[Define](define)
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
	module := &tModule{}

	module.routersInfo = append(module.routersInfo, p.CreateHome(t.define.Router.Home))
	for _, html := range t.define.Router.Html {
		module.routersInfo = append(module.routersInfo, p.CreateHtml(html.Path, html.Label))
	}
	for _, a := range t.define.Router.Frontend {
		module.routersInfo = append(module.routersInfo, p.CreateHtml(a, apinto_module.RouterLabelAssets))
	}
	for path, ms := range t.define.Router.Api {
		for method, att := range ms {
			module.routersInfo = append(module.routersInfo, p.CreateApi(fmt.Sprintf("%s.%s", method, path), method, path, att.Label))

		}
	}
	for path, ms := range t.define.Router.OpenApi {
		for method, att := range ms {
			module.routersInfo = append(module.routersInfo, p.CreateOpenApi(fmt.Sprintf("%s.%s", method, path), method, path, att.Label))
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
	return decode[Config](config)
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
