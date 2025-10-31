package service

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"sync/atomic"

	apinto_module "github.com/eolinker/apinto-dashboard/module"
)

var _ IProviderService = (*ProviderService)(nil)

type IProviderService interface {
	apinto_module.IProviders
	set(providers apinto_module.IProviders)
}
type ProviderService struct {
	atomic.Pointer[apinto_module.IProviders]
}

func (p *ProviderService) Status(key string, namespaceId int, cluster string) (pm3.CargoStatus, string) {
	ps := p.Load()
	if ps == nil || *ps == nil {
		return pm3.None, ""
	}
	return (*ps).Status(key, namespaceId, cluster)
}

func NewProviderService() IProviderService {
	p := &ProviderService{}

	return p
}
func (p *ProviderService) set(providers apinto_module.IProviders) {
	p.Store(&providers)
}
func (p *ProviderService) Provider(name string) (pm3.Provider, bool) {
	ps := p.Load()
	if ps == nil || *ps == nil {
		return nil, false
	}
	return (*ps).Provider(name)
}
