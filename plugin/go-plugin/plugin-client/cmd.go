package plugin_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eolinker/apinto-dashboard/config"
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"io"
	"os"
	"os/exec"
)

var (
	logOutPut io.Writer
	logLevel  hclog.Level
)

func CreateClient(id, name string, path string, configs ...string) *plugin.Client {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("plugin:%s", id),
		Output:     logOutPut,
		TimeFormat: "[2006-01-02 15:04:05]",
		Level:      logLevel,
	})

	cmd := exec.Command(path)
	cmd.Args[0] = fmt.Sprintf("%s:%s", os.Args[0], id)
	buf := &bytes.Buffer{}
	json.NewEncoder(buf).Encode(config.GetConfigData())

	cmd.Stdin = buf

	cmd.Env = append(cmd.Env, "ApintoDashboardConfig=true")
	for _, cv := range configs {
		cmd.Env = append(cmd.Env, cv)
	}
	// We're a host! Start by launching the plugin process.
	c := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.HandshakeConfig,
		Plugins:         CreateClientProxy(id, name),

		Cmd:              cmd,
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		StartTimeout:     0,
	})
	return c
}

func SetLog(level string, w io.Writer) {
	logLevel = hclog.LevelFromString(level)
	logOutPut = w
}
