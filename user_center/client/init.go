package client

import (
	"github.com/eolinker/eosc/common/bean"
)

func InitUserCenterClient(url string) {
	client := newIUserCenterClient(url)

	bean.Injection(&client)
}
