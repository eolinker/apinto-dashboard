package apinto_module

import (
	"fmt"
	"github.com/eolinker/apinto-dashboard/pm3"
	"strings"
	"sync"
)

type ProviderMultiple struct {
	providers map[string]pm3.Provider
}

func NewProviderMultiple(providers map[string]pm3.Provider) pm3.Provider {
	return &ProviderMultiple{providers: providers}
}

func (p *ProviderMultiple) Provide(namespaceID int) []pm3.Cargo {
	result := make([]pm3.Cargo, 0)
	for name, provider := range p.providers {
		cargos := provider.Provide(namespaceID)
		for _, g := range cargos {
			result = append(result, pm3.Cargo{
				Value: fmt.Sprint(g.Value, "@", name),
				Title: g.Title,
			})
		}
	}
	return result

}

type IProviders interface {
	Provider(skill string) (pm3.Provider, bool)
	pm3.ProviderStatus
}
type ProvidersManager struct {
	providersForSkill map[string]pm3.Provider // map[cargoName]pm3.ProviderMultiple
	statusForModule   map[string]pm3.ProviderStatus
}

func NewProvidersManager(providersForSkill map[string]pm3.Provider, statusForModule map[string]pm3.ProviderStatus) IProviders {
	return &ProvidersManager{providersForSkill: providersForSkill, statusForModule: statusForModule}
}

func (p *ProvidersManager) Provider(name string) (pm3.Provider, bool) {
	provider, has := p.providersForSkill[name]
	return provider, has
}

func (p *ProvidersManager) Status(key string, namespaceId int, cluster string) (pm3.CargoStatus, string) {
	i := strings.LastIndex(key, "@")
	if i < 0 {
		return pm3.None, ""
	}
	providerName := key[i+1:]
	provider, has := p.statusForModule[providerName]
	if !has {
		return pm3.None, ""
	}
	return provider.Status(key[:i], namespaceId, cluster)

}

type ProviderBuilder struct {
	lock      sync.Mutex
	providers map[string]pm3.ProviderSupport // map[moduleName]map[cargoName]pm3.Provider
}

func NewProviderBuilder() *ProviderBuilder {
	return &ProviderBuilder{
		lock:      sync.Mutex{},
		providers: make(map[string]pm3.ProviderSupport),
	}
}

func (p *ProviderBuilder) Append(name string, ps pm3.ProviderSupport) *ProviderBuilder {
	p.lock.Lock()
	defer p.lock.Unlock()
	if ps != nil {
		p.providers[name] = ps
	}
	return p
}
func (p *ProviderBuilder) Build() IProviders {
	p.lock.Lock()
	defer p.lock.Unlock()
	providerForSkillOfModule := make(map[string]map[string]pm3.Provider)
	statusForModule := make(map[string]pm3.ProviderStatus)
	for moduleName, providerSupport := range p.providers {

		for name, provider := range providerSupport.Provider() {
			cps, has := providerForSkillOfModule[name]
			if !has {
				providerForSkillOfModule[name] = make(map[string]pm3.Provider)
				cps = providerForSkillOfModule[name]
			}
			cps[moduleName] = provider
		}
		statusForModule[moduleName] = providerSupport
	}
	providersForSkill := make(map[string]pm3.Provider)
	for skill, cps := range providerForSkillOfModule {
		providersForSkill[skill] = NewProviderMultiple(cps)
	}
	return NewProvidersManager(providersForSkill, statusForModule)
}
