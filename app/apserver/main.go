package main

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin_timer"
	"net"
	"net/http"
	"os"

	"github.com/eolinker/apinto-dashboard/app/apserver/version"
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
	//engine := gin.Default()

	//registerRouter(engine)

	var coreService core.ICore
	bean.Autowired(&coreService)
	var front core.EngineCreate = new(Front)
	bean.Injection(&front)
	initDB()

	err := bean.Check()
	if err != nil {
		log.Fatal(err)
	}
	// todo 执行导航初始化
	// todo 执行内置插件初始化
	go plugin_timer.ExtenderTimer()
	// todo 不适合开源，后续通过插件接入
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", GetPort()))
	if err != nil {
		panic(err)
	}

	s := http.Server{Handler: coreService}
	s.Serve(listener)
}

type Front struct {
}

func (f *Front) CreateEngine() *gin.Engine {
	engine := gin.Default()
	controller.EmbedFrontend(engine)
	return engine
}
