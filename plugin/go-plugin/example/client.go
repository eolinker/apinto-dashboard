package main

import (
	"github.com/eolinker/apinto-dashboard/config"
	"github.com/eolinker/apinto-dashboard/module/builder"
	client "github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin-client"
	shared "github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
)

func main() {
	config.ReadConfig()
	// Create an hclog.Logger
	c := client.CreateClient("demo", "./plugin/plugin")
	defer c.Kill()

	// Connect via RPC
	rpcClient, err := c.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(shared.PluginHandlerName)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	handler := raw.(client.ClientHandler)

	b := builder.NewModuleBuilder(gin.New())
	b.Append(handler)
	server, _, err := b.Build()
	if err != nil {
		return
	}
	listen, err := net.Listen("tcp", ":8880")
	if err != nil {
		return
	}
	http.Serve(listen, server)

}
