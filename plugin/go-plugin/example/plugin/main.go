package main

import (
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/plugin"
)

func main() {

	module := NewPlugin()
	ps := plugin.NewPlugin(module)

	ps.Server()

}
