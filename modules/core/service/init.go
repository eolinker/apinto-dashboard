package service

import (
	apinto_module "github.com/eolinker/apinto-module"
	"github.com/eolinker/eosc/common/bean"
)

func init() {
	navigationService := newNavigationService()
	bean.Injection(&navigationService)

	providerService := NewProviderService()

	var iProviders apinto_module.IProviders = providerService
	bean.Injection(&iProviders)
	iCore := NewService(providerService)
	bean.Injection(&iCore)

}
