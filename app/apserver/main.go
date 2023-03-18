package main

import (
	"context"
	"fmt"
	"os"

	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/db_migrator"
	cluster_service "github.com/eolinker/apinto-dashboard/modules/cluster"
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
	app.Run(os.Args)
}

func run() {

	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	registerRouter(engine)

	err := bean.Check()
	if err != nil {
		log.Fatal(err)
	}

	//初始化数据库表 sql操作
	initDB()

	//初始化超管账号 和清除超管缓存
	// todo 不适合开源，后续通过插件接入

	//初始化集群插件
	initClustersPlugin()

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

func initClustersPlugin() {
	var clientService cluster_service.IApintoClient
	bean.Autowired(&clientService)

	err := clientService.InitClustersGlobalPlugin(context.Background())
	if err != nil {
		panic(err)
	}
}
