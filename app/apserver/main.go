package main

import (
	"context"
	"fmt"
	machine_code "github.com/eolinker/apinto-dashboard/app/apserver/machine-code"
	"github.com/eolinker/apinto-dashboard/app/apserver/version"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/db_migrator"
	"github.com/eolinker/apinto-dashboard/service/apinto-client"
	"github.com/eolinker/apinto-dashboard/service/notice-service"
	"github.com/eolinker/apinto-dashboard/service/user-service"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/apinto-dashboard/timer"
	"github.com/eolinker/apinto-dashboard/user_center/client"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"os"
)

func main() {
	app := &cli.App{
		Name:     "apserver",
		HelpName: "apserver",
		Usage:    "apinto dashboard enterprise",

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

	client.InitUserCenterClient(GetUserCenterUrl())
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
	initAdmin()

	//初始化集群插件
	initClustersPlugin()

	//初始化通知渠道驱动管理器
	initNoticeChannelDriver()

	//定时任务
	go timer.TaskTimer()

	if err = engine.Run(fmt.Sprintf(":%d", GetPort())); err != nil {
		panic(err)
	}
}

func initNoticeChannelDriver() {
	var noticeChannelService notice_service.INoticeChannelService
	bean.Autowired(&noticeChannelService)
	err := noticeChannelService.InitChannelDriver()
	if err != nil {
		panic(err)
	}
}

func initAdmin() {
	var userInfoService user_service.IUserInfoService
	bean.Autowired(&userInfoService)
	err := userInfoService.CreateAdmin()
	if err != nil {
		panic(err)
	}

	userInfoService.CleanAdminCache()
}

func initDB() {
	var iDb store.IDB

	bean.Autowired(&iDb)

	ctx := context.Background()
	db := iDb.DB(ctx)

	db_migrator.InitSql(db)

	//计算获取加密后的机器码
	machineCode := generateMachineCode(db)
	//设置机器码
	machine_code.SetMachineCode(machineCode)
}

func initClustersPlugin() {
	var clientService apinto_client.IApintoClient
	bean.Autowired(&clientService)

	err := clientService.InitClustersGlobalPlugin(context.Background())
	if err != nil {
		panic(err)
	}
}

func generateMachineCode(db *gorm.DB) string {
	sqlDb, _ := db.DB()
	//获取数据库唯一UUID
	var variableName, serverUUID string
	row := sqlDb.QueryRow(`SHOW VARIABLES LIKE "server_uuid"`)
	err := row.Scan(&variableName, &serverUUID)
	if err != nil {
		panic(err)
	}
	//加密机器码
	machineCodeRaw := fmt.Sprintf("%s.%s.%s.%s.%s.%s", GetDBIp(), GetDBPort(), GetDBUserName(), GetDBPassword(), GetDbName(), serverUUID)
	return common.Md5(fmt.Sprintf("%s%s", machine_code.Salt, machineCodeRaw))
}
