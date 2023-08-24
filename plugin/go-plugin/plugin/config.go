package plugin

import (
	"encoding/json"
	"github.com/eolinker/apinto-dashboard/config"
	"os"
)

func init() {
	initConfig()
	initLogger()
}
func initConfig() {

	if os.Getenv("ApintoDashboardConfig") != "true" {
		config.ReadConfig()
		return
	}
	c := new(config.Config)
	err := json.NewDecoder(os.Stdin).Decode(c)
	if err != nil {
		logger.Error("read config", err)
		return
	}
	config.SetConfig(c)

}
