package flux

import "github.com/eolinker/eosc/common/bean"

func init() {
	iFluxQuery := newFluxQuery()
	bean.Injection(&iFluxQuery)

	//初始化buckets配置
	initBucketsConfig()
	//初始化tasks定时脚本配置
	initTasksConfig()
}
