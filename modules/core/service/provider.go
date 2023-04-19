package service

import (
	apinto_module "github.com/eolinker/apinto-module"
	"sync/atomic"
)

var _ IProviderService = (*ProviderService)(nil)

type IProviderService interface {
	apinto_module.IProviders
	set(providers apinto_module.IProviders)
}
type ProviderService struct {
	atomic.Pointer[apinto_module.IProviders]
}

func NewProviderService() IProviderService {
	p := &ProviderService{}

	return p
}
func (p *ProviderService) set(providers apinto_module.IProviders) {
	p.Store(&providers)
}
func (p *ProviderService) Provider(name string) (apinto_module.Provider, bool) {
	ps := p.Load()
	if ps == nil || *ps == nil {
		return nil, false
	}
	return (*ps).Provider(name)
}
