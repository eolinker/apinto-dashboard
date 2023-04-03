package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
)

type IPlugin[T any] interface {
	List() ([]*T, error)
	Set(config []*T) error
}

type plugin[T any] struct {
	IClient
}
type PluginConfig[T any] struct {
	Plugin []*T `json:"plugins,omitempty"`
}

func newIPlugin[T any](IClient IClient) IPlugin[T] {
	return &plugin[T]{IClient: IClient}
}

func (p *plugin[T]) List() ([]*T, error) {
	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s", addr, "setting/plugin")
		res, code, err := requestDo(http.MethodGet, url, nil)

		if err != nil {
			log.Errorf("plugin-api-list err=%s", err.Error())
			resErr = err
			continue
		}
		if code != http.StatusOK {
			log.Errorf("plugin-api-list err=%s", string(res))
			resErr = errors.New(string(res))
			continue
		}

		infos := new(PluginConfig[T])

		if err = json.Unmarshal(res, &infos); err != nil {
			log.Errorf("plugin-json.Unmarshal err=%s", err)
			resErr = err
			continue
		}

		return infos.Plugin, nil
	}
	return nil, resErr
}

//Set T 包含所有组件 修改就是一次性修改所有的全局插件
func (p *plugin[T]) Set(t []*T) error {
	cf := PluginConfig[T]{
		Plugin: t,
	}
	bytes, _ := json.Marshal(cf)

	var resErr error
	for _, addr := range p.addrs() {
		url := fmt.Sprintf("%s/api/%s", addr, "setting/plugin")

		res, code, err := requestDo(http.MethodPut, url, bytes)

		if err != nil {
			resErr = err
			log.Errorf("profession-api-set url=%s,err=%s,req=%s", url, err, string(bytes))
			continue
		}
		if code != http.StatusOK {
			log.Errorf("profession-api-set url=%s, err=%s,req=%s", url, string(res), string(bytes))
			continue
		}
		return nil
	}
	return resErr
}
