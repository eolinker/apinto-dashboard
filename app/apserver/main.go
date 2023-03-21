package main

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin_timer"
	"os"

	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/db_migrator"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "apserver",
		HelpName: "apserver",
		Usage:    "apinto dashboard",

		Version:     version.Version,
		Description: "",
		Commands:    []*cli.Command{version.Build()},
		Flags:       nil,
		Action: func(context *cli.Context) error {
			run()
			return nil
		},
	}
	_ = app.Run(os.Args)
}

func run() {

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	registerRouter(engine)

	//初始化数据库表 sql操作
	initDB()

	err := bean.Check()
	if err != nil {
		log.Fatal(err)
	}

	go plugin_timer.ExtenderTimer()
	// todo 不适合开源，后续通过插件接入

	if err = engine.Run(fmt.Sprintf(":%d", GetPort())); err != nil {
		panic(err)
	}
}

func initDB() {
	var iDb store.IDB

	bean.Autowired(&iDb)

	ctx := context.Background()
	db := iDb.DB(ctx)

	db_migrator.InitSql(db)

}
