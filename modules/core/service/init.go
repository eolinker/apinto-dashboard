package service

import "github.com/eolinker/eosc/common/bean"

func init() {
	providerService := NewProviderService()
	bean.Injection(&providerService)

	iCore := NewService()
	bean.Injection(&iCore)
}
