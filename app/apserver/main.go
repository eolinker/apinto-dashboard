package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/eolinker/apinto-dashboard/config"
	"github.com/eolinker/apinto-dashboard/frontend"
	grpcservice "github.com/eolinker/apinto-dashboard/grpc-service"
	apintomodule "github.com/eolinker/apinto-dashboard/module"
	"github.com/eolinker/apinto-dashboard/module/builder"
	"github.com/eolinker/apinto-dashboard/modules/grpc-service/service"
	"github.com/eolinker/apinto-dashboard/pm3/embed_registry"
	"github.com/eolinker/apinto-dashboard/report"
	"github.com/soheilhy/cmux"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/modules/core"
	corecontroller "github.com/eolinker/apinto-dashboard/modules/core/controller"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin_timer"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
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
	config.ReadConfig()
	initLog()
	gin.SetMode(gin.ReleaseMode)
	builder.AddSystemModule(corecontroller.NewModule("core.apinto.com", "core"))

	builder.SetIndexHtml(frontend.IndexHtml)
	builder.AddAssetMiddleware(frontend.AddExpires, frontend.GzipHandler)

	var coreService core.ICore
	bean.Autowired(&coreService)

	var front core.EngineCreate = new(Front)
	bean.Injection(&front)
	config.InitDb()
	config.InitRedis()
	err := bean.Check()
	if err != nil {
		log.Fatal(err)
	}
	report.InitReport(config.GetLogDir(), config.DisableReport())
	// 执行内置插件初始化
	err = embed_registry.InitEmbedPlugins()
	if err != nil {
		log.Fatal(err)
	}
	err = embed_registry.LoadLocalPlugins("plugins", "plugin.yml")
	if err != nil {
		log.Fatal(err)
	}
	_ = coreService.ReloadModule()
	go plugin_timer.ExtenderTimer()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetPort()))
	if err != nil {
		panic(err)
	}
	// Create a cmux.
	m := cmux.New(listener)
	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))

	httpL := m.Match(cmux.HTTP1Fast(http.MethodPatch))

	httpServer := &http.Server{Handler: coreService, IdleTimeout: 90 * time.Second}
	grpcServer := grpc.NewServer()

	grpcservice.RegisterGetConsoleInfoServer(grpcServer, service.NewConsoleInfoService())
	grpcservice.RegisterNoticeSendServer(grpcServer, service.NewNoticeSendService())
	grpcservice.RegisterApiServiceServer(grpcServer, service.NewApiServiceServer())

	console := newConsoleServer(httpServer, grpcServer)
	go func() {
		err := httpServer.Serve(httpL)
		if err != nil {
			log.Error("listen httpServer error: ", err)
		}
	}()
	go func() {
		err := grpcServer.Serve(grpcL)
		if err != nil {
			log.Error("listen grpcServer error: ", err)
		}
	}()
	go func() {
		err := m.Serve()
		if err != nil {
			log.Error("server close: ", err)
			return
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
	engine.Use(apintomodule.SetRepeatReader)
	return engine
}
