package apinto

import (
	"fmt"
	"net/http"
)

var (
	client *admin
)

func Init(address []string) {
	// just for test
	client = NewAdmin(address)
}

func Client() IClient {
	return client
}

func (a *admin) List(profession string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/%s", a.GetNode(), profession)
	return a.do(http.MethodGet, url, nil)
}

func (a *admin) Get(profession string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/%s/%s", a.GetNode(), profession, name)
	return a.do(http.MethodGet, url, nil)
}

func (a *admin) Update(profession string, name string, body []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/%s/%s", a.GetNode(), profession, name)
	return a.do(http.MethodPut, url, body)
}

func (a *admin) Create(profession string, body []byte) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/%s", a.GetNode(), profession)
	return a.do(http.MethodPost, url, body)
}

func (a *admin) Delete(profession string, name string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/api/%s/%s", a.GetNode(), profession, name)
	return a.do(http.MethodDelete, url, nil)
}

func (a *admin) Patch(profession string, name string, body []byte) ([]byte, int, error) {
	// Todo 等apinto完成对接
	url := fmt.Sprintf("%s/api/%s/%s", a.GetNode(), profession, name)
	return a.do(http.MethodPatch, url, body)
}
func (a *admin) PatchPath(profession string, name string, path string, body []byte) ([]byte, int, error) {
	// Todo 等apinto完成对接
	url := fmt.Sprintf("%s/api/%s/%s/%s", a.GetNode(), profession, name, path)
	return a.do(http.MethodPatch, url, body)
}

func (a *admin) Render(profession string, driver string) ([]byte, int, error) {
	url := fmt.Sprintf("%s/profession/%s/%s", a.GetNode(), profession, driver)
	return a.do(http.MethodGet, url, nil)
}

func (a *admin) Drivers(profession string) (data []byte, code int, err error) {
	url := fmt.Sprintf("%s/profession/%s", a.GetNode(), profession)
	return a.do(http.MethodGet, url, nil)
}

func (a *admin) Extenders() (data []byte, code int, err error) {
	url := fmt.Sprintf("%s/extenders", a.GetNode())
	return a.do(http.MethodGet, url, nil)
}
func (a *admin) Extender(id string) (data []byte, code int, err error) {
	url := fmt.Sprintf("%s/extender/%s", a.GetNode(), id)
	return a.do(http.MethodGet, url, nil)
}
