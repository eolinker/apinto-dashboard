package main

import (
	"github.com/eolinker/apinto-dashboard/config"
	client "github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin-client"
	shared "github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
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
	middleware, err := handler.CreateMiddleware("test", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Any("/*tt", middleware.Middleware, func(context *gin.Context) {
		handler.ServerGin(context, nil, nil)
	})
	err = engine.Run(":8880")
	if err != nil {
		log.Info(err)
		return
	}
}
