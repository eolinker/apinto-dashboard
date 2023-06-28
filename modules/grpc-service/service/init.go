package service

import (
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	iInstalledCache := newNavigationModulesCache()
	bean.Injection(&iInstalledCache)

}
