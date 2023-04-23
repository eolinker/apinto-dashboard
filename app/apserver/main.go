package main

import (
	"fmt"
	grpc_service "github.com/eolinker/apinto-dashboard/grpc-service"
	"github.com/eolinker/apinto-dashboard/modules/grpc-service/service"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"

	"github.com/eolinker/apinto-dashboard/initialize"
	"github.com/eolinker/apinto-dashboard/modules/core"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin_timer"

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

	// 执行内置插件初始化
	err = initialize.InitPlugins()
	if err != nil {
		log.Fatal(err)
	}
	coreService.ReloadModule()
	go plugin_timer.ExtenderTimer()
	// todo 不适合开源，后续通过插件接入

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", GetPort()))
	if err != nil {
		panic(err)
	}
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8181))
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{Handler: coreService}
	grpcServer := grpc.NewServer()

	grpc_service.RegisterGetConsoleInfoServer(grpcServer, service.NewConsoleInfoService())
	grpc_service.RegisterNoticeSendServer(grpcServer, service.NewNoticeSendService())

	console := newConsoleServer(httpServer, grpcServer)
	go func() {
		err := httpServer.Serve(httpListener)
		if err != nil {
			log.Error("listen httpServer error: ", err)
		}
	}()
	go func() {
		err := grpcServer.Serve(grpcListener)
		if err != nil {
			log.Error("listen grpcServer error: ", err)
		}
	}()

	err = console.Wait()
	if err != nil {
		log.Fatal(err)
	}

}

type Front struct {
}

func (f *Front) CreateEngine() *gin.Engine {
	engine := gin.Default()
	engine.Use(apinto_module.SetRepeatReader)
	return engine
}
