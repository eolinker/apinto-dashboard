package apinto

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
	"strings"
	"sync"
)

type IAdmin interface {
	GetNode() string
}

type IClient interface {
	List(profession string) (data []byte, code int, err error)
	Get(profession string, name string) (data []byte, code int, err error)
	Create(profession string, body []byte) (data []byte, code int, err error)
	Delete(profession string, name string) (data []byte, code int, err error)
	Update(profession string, name string, body []byte) (data []byte, code int, err error)
	Patch(profession string, name string, body []byte) (data []byte, code int, err error)
	PatchPath(profession string, name string, path string, body []byte) (data []byte, code int, err error)
	Render(profession string, driver string) (data []byte, code int, err error)
	Drivers(profession string) (data []byte, code int, err error)
	Extenders() (data []byte, code int, err error)
	Extender(id string) (data []byte, code int, err error)
}

type admin struct {
	lock   sync.RWMutex
	nodes  []string
	client *http.Client
}

func NewAdmin(nodes []string) *admin {
	return &admin{
		nodes:  nodes,
		client: http.DefaultClient,
	}
}
func (a *admin) GetNode() string {
	a.lock.RLock()
	v := a.nodes[0]
	a.lock.RUnlock()
	return v
}

func (a *admin) updateNodes(nodes []string) error {
	a.lock.Lock()
	a.nodes = nodes
	a.lock.Unlock()
	return nil
}

func (a *admin) addNode(node string) error {
	if node == "" {
		return errors.New("empty node url")
	}
	node = strings.TrimSuffix(node, "/")
	err := a.ping(node)
	if err != nil {
		return err
	}
	a.lock.Lock()
	a.nodes = append(a.nodes, node)

	a.lock.Unlock()
	return nil
}

func (a *admin) ping(node string) error {
	url := fmt.Sprintf("%s/api/", node)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node addr %s can not be connected", node)
	}
	return nil
}

func (a *admin) do(method string, url string, body []byte) ([]byte, int, error) {
	req, err := a.newRequest(method, url, body)
	if err != nil {
		log.Error("new request:", err)
		return nil, 500, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.client.Do(req)
	if err != nil {
		log.Error("do request:", err)
		return nil, 500, err
	}
	data, err := ReadBody(resp.Body)
	if err != nil {
		log.Error("read body:", err)

		return nil, 500, err
	}
	return data, resp.StatusCode, nil
}

func (a *admin) newRequest(method string, url string, body []byte) (*http.Request, error) {
	if body == nil {
		return http.NewRequest(method, url, nil)
	}
	return http.NewRequest(method, url, bytes.NewReader(body))
}
