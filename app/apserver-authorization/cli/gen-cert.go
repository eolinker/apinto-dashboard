package cli

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"time"
)

const (
	machineCodeFlag     = "machine_code"
	companyFlag         = "company"
	editionFlag         = "edition"
	beginTimeFlag       = "begin_time"
	endTimeFlag         = "end_time"
	controllerCountFlag = "controller_count"
	nodeCountFlag       = "node_count"

	editionStandard = "standard"
	editionPremium  = "premium"
	editionUltimate = "ultimate"

	certDir = "./export/cert/%s/%s-%s"
)

func GenCert() *cli.Command {
	timeLocation, _ := time.LoadLocation("Asia/Shanghai")
	return &cli.Command{
		Name:  "gen-cert",
		Usage: "generate the cert file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     machineCodeFlag,
				Usage:    "机器码",
				Required: true,
			},
			&cli.StringFlag{
				Name:     companyFlag,
				Usage:    "公司名称",
				Required: true,
			},
			&cli.StringFlag{
				Name:     editionFlag,
				Usage:    "授权版本 standard（标准版）、premium(高级版)、ultimate(旗舰版）",
				Required: true,
			},
			&cli.TimestampFlag{
				Name:  beginTimeFlag,
				Usage: "开始时间 格式为 2006-01-02",
				//Usage:    "开始时间 格式为 2006-01-02T15:04:05",
				Layout: "2006-01-02",
				//Layout:   "2006-01-02T15:04:05",
				Timezone: timeLocation,
				Required: true,
			},
			&cli.TimestampFlag{
				Name:  endTimeFlag,
				Usage: "结束时间 格式为 2006-01-02",
				//Usage:    "结束时间 格式为 2006-01-02T15:04:05",
				Required: true,
				Layout:   "2006-01-02",
				//Layout:   "2006-01-02T15:04:05",
				Timezone: timeLocation,
			},
			&cli.IntFlag{
				Name:     controllerCountFlag,
				Usage:    "控制台（apserver）数量，不填则默认不限制，若值不为整型，则报错",
				Required: false,
			},
			&cli.IntFlag{
				Name:     nodeCountFlag,
				Usage:    "节点（node）数量，不填则默认不限制，若值不为整型，则报错",
				Required: false,
			},
		},
		Action: GenCertFunc,
	}
}

// GenCertFunc 生成授权证书
func GenCertFunc(c *cli.Context) error {
	machineCode := c.String(machineCodeFlag)
	company := c.String(companyFlag)
	edition := c.String(editionFlag)
	beginTime := c.Timestamp(beginTimeFlag)
	endTime := c.Timestamp(endTimeFlag)
	controllerCount := c.Int(controllerCountFlag)
	nodeCount := c.Int(nodeCountFlag)

	//校验edition
	switch edition {
	case editionStandard, editionPremium, editionUltimate:
	default:
		return fmt.Errorf("edition %s is invalid. ", edition)
	}

	//默认为-1, 表示不限
	if controllerCount == 0 {
		controllerCount = -1
	}

	//默认为-1, 表示不限
	if nodeCount == 0 {
		nodeCount = -1
	}

	//NewCertInfo
	certInfo := newCertInfo(machineCode, company, edition, *beginTime, *endTime, controllerCount, nodeCount)
	certDirPath := fmt.Sprintf(certDir, company, beginTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
	return buildCert(certInfo, certDirPath)
}
