package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/eosc/log"
	"net/http"
	"time"
)

var _ IClient = (*Client)(nil)

type IAdmin interface {
}

type IClient interface {
	ForDiscovery() IProfession[DiscoveryConfig, WorkerInfo]
	ForAuth() IProfession[AuthConfig, WorkerInfo]
	ForService() IProfession[ServiceConfig, WorkerInfo]
	ForOutput() IProfession[interface{}, WorkerInfo]
	ForRedisOutput() IProfession[RedisOutput, WorkerInfo]
	ForInfluxV2Output() IProfession[InfluxV2Output, WorkerInfo]
	ForApp() IProfession[ApplicationConfig, WorkerInfo]
	ForGlobalPlugin() IPlugin[GlobalPlugin]
	ForRouter() IProfession[RouterConfig, RouterInfo]
	ForCert() ICert
	ForVariable() IVariable
	ForStrategy() IStrategy
	ClusterInfo() (*ClusterInfo, error)
	addrs() []string
}

type Client struct {
	adminAddrs     []string
	discovery      IProfession[DiscoveryConfig, WorkerInfo]
	auth           IProfession[AuthConfig, WorkerInfo]
	output         IProfession[any, WorkerInfo]
	redisOutput    IProfession[RedisOutput, WorkerInfo]
	influxV2Output IProfession[InfluxV2Output, WorkerInfo]
	service        IProfession[ServiceConfig, WorkerInfo]
	router         IProfession[RouterConfig, RouterInfo]
	app            IProfession[ApplicationConfig, WorkerInfo]
	cert           ICert
	plugin         IPlugin[GlobalPlugin]
	variable       IVariable
	strategy       IStrategy
}

func (c *Client) addrs() []string {
	return c.adminAddrs
}

func (c *Client) ping() error {
	if len(c.adminAddrs) == 0 {
		return errors.New("获取不到可连接的地址")
	}
	var err error
	var resp *http.Response
	for _, node := range c.adminAddrs {
		req := http.Client{Timeout: time.Second * 1}
		url := fmt.Sprintf("%s/system/info", node)
		resp, err = req.Get(url)
		if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
			err = fmt.Errorf("node addr %s can not be connected", node)
			continue
		}
		return nil
	}

	return err
}

func NewClient(addrs []string) (IClient, error) {
	c := &Client{adminAddrs: addrs}
	err := c.ping()
	if err != nil {
		return nil, err
	}
	c.init()
	return c, nil
}
func (c *Client) init() {
	c.discovery = newIProfession[DiscoveryConfig, WorkerInfo](c, "discovery")
	c.auth = newIProfession[AuthConfig, WorkerInfo](c, "auth")
	c.output = newIProfession[interface{}, WorkerInfo](c, "output")
	c.redisOutput = newIProfession[RedisOutput, WorkerInfo](c, "output")
	c.influxV2Output = newIProfession[InfluxV2Output, WorkerInfo](c, "output")
	c.service = newIProfession[ServiceConfig, WorkerInfo](c, "service")
	c.router = newIProfession[RouterConfig, RouterInfo](c, "router")
	c.cert = newCert(c)
	c.app = newIProfession[ApplicationConfig, WorkerInfo](c, "app")
	c.plugin = newIPlugin[GlobalPlugin](c)
	c.strategy = newIStrategy(c)
	c.variable = newIVariable(c, "variable")
}
func (c *Client) ForDiscovery() IProfession[DiscoveryConfig, WorkerInfo] {
	return c.discovery
}

func (c *Client) ForApp() IProfession[ApplicationConfig, WorkerInfo] {
	return c.app
}

func (c *Client) ForCert() ICert {
	return c.cert
}

func (c *Client) ForGlobalPlugin() IPlugin[GlobalPlugin] {
	return c.plugin
}

func (c *Client) ForAuth() IProfession[AuthConfig, WorkerInfo] {
	return c.auth
}

func (c *Client) ForStrategy() IStrategy {
	return c.strategy
}

func (c *Client) ForInfluxV2Output() IProfession[InfluxV2Output, WorkerInfo] {
	return c.influxV2Output
}

func (c *Client) ForOutput() IProfession[interface{}, WorkerInfo] {
	return c.output
}

func (c *Client) ForRedisOutput() IProfession[RedisOutput, WorkerInfo] {
	return c.redisOutput
}

func (c *Client) ForService() IProfession[ServiceConfig, WorkerInfo] {
	return c.service
}

func (c *Client) ForRouter() IProfession[RouterConfig, RouterInfo] {
	return c.router
}
func (c *Client) ForVariable() IVariable {
	return c.variable
}
func (c *Client) ClusterInfo() (nodes *ClusterInfo, errRes error) {

	for _, addr := range c.adminAddrs {
		url := fmt.Sprintf("%s/system/info", addr)
		bytes, code, err := requestDo(http.MethodGet, url, nil)
		if err != nil {
			log.Errorf("client-nodes error=%s", err)
			errRes = err
			continue
		}
		if code != http.StatusOK {
			errRes = errors.New(string(bytes))
			continue
		}

		info := new(ClusterInfo)
		if err = json.Unmarshal(bytes, &info); err != nil {
			log.Errorf("client-nodes error=%s", err)
			errRes = err
			continue
		}

		nodes = info
		errRes = nil
		return
	}

	return
}
