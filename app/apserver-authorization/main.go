package main

import (
	"github.com/eolinker/apinto-dashboard/app/apserver-authorization/cli"
	"github.com/eolinker/eosc/log"
	"os"
)

func main() {
	InitCLILog()

	app := cli.NewApp()
	app.Default()
	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
	log.Close()
}
