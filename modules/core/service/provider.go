package service

import (
	"github.com/eolinker/apinto-dashboard/modules/core"
	apinto_module "github.com/eolinker/apinto-module"
	"sync/atomic"
)

var _ core.IProviders = (*ProviderService)(nil)

type ProviderService struct {
	atomic.Pointer[apinto_module.IProviders]
}

func NewProviderService() *ProviderService {
	p := &ProviderService{}
	return p
}
func (p *ProviderService) Set(providers apinto_module.IProviders) {
	p.Store(&providers)
}
func (p *ProviderService) Provider(name string) (apinto_module.Provider, bool) {
	ps := p.Load()
	if ps == nil || *ps == nil {
		return nil, false
	}
	return (*ps).Provider(name)
}
