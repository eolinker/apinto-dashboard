package local

import (
	"fmt"
	apinto_module "github.com/eolinker/apinto-dashboard/module"
	client "github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin-client"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/hashicorp/go-plugin"
)

var (
	_ apinto_module.ModuleNeedKill = (*ProxyAPi)(nil)
)

type ProxyAPi struct {
	client *plugin.Client
	cmd    string
	id     string
	module string
	params map[string]string

	client.ClientHandler
}

func (p *ProxyAPi) Kill() {
	p.kill()
}

func (p *ProxyAPi) kill() {

	if p.client != nil {
		p.client.Kill()
		p.client = nil
	}

}
func (p *ProxyAPi) initClient() error {
	params := make([]string, 0, len(p.params))
	for k, v := range p.params {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	c := client.CreateClient(p.id, p.module, cmdPath(p.cmd), params...)

	rpcClient, err := c.Client()
	if err != nil {
		return err
	}
	// Request the plugin
	raw, err := rpcClient.Dispense(shared.PluginHandlerName)
	if err != nil {

		return err
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	handler := raw.(client.ClientHandler)
	p.ClientHandler = handler
	p.client = c
	return nil
}
func NewProxyAPi(cmd string, id, module string, config *Config) (*ProxyAPi, error) {

	p := &ProxyAPi{cmd: cmd, id: id, module: module, params: config.Initialize}
	err := p.initClient()
	if err != nil {
		return nil, err
	}
	return p, nil
}
