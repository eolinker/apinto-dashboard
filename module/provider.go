package apinto_module

import (
	"fmt"
	"strings"
	"sync"
)

type CargoItem struct {
	Value string `json:"name"`
	Title string `json:"title"`
}

type ProviderMultiple struct {
	providers map[string]Provider
}

func NewProviderMultiple(providers map[string]Provider) Provider {
	return &ProviderMultiple{providers: providers}
}

func (p *ProviderMultiple) Provide(namespaceID int) []Cargo {
	result := make([]Cargo, 0)
	for name, provider := range p.providers {
		cargos := provider.Provide(namespaceID)
		for _, g := range cargos {
			result = append(result, Cargo{
				Value: fmt.Sprint(g.Value, "@", name),
				Title: g.Title,
			})
		}
	}
	return result

}

type IProviders interface {
	Provider(skill string) (Provider, bool)
	ProviderStatus
}
type ProvidersManager struct {
	providersForSkill map[string]Provider // map[cargoName]ProviderMultiple
	statusForModule   map[string]ProviderStatus
}

func NewProvidersManager(providersForSkill map[string]Provider, statusForModule map[string]ProviderStatus) IProviders {
	return &ProvidersManager{providersForSkill: providersForSkill, statusForModule: statusForModule}
}

func (p *ProvidersManager) Provider(name string) (Provider, bool) {
	provider, has := p.providersForSkill[name]
	return provider, has
}

func (p *ProvidersManager) Status(key string, namespaceId int, cluster string) (CargoStatus, string) {
	i := strings.LastIndex(key, "@")
	if i < 0 {
		return None, ""
	}
	providerName := key[i+1:]
	provider, has := p.statusForModule[providerName]
	if !has {
		return None, ""
	}
	return provider.Status(key[:i], namespaceId, cluster)

}

type ProviderBuilder struct {
	lock      sync.Mutex
	providers map[string]ProviderSupport // map[moduleName]map[cargoName]provider
}

func NewProviderBuilder() *ProviderBuilder {
	return &ProviderBuilder{
		lock:      sync.Mutex{},
		providers: make(map[string]ProviderSupport),
	}
}

func (p *ProviderBuilder) Append(name string, ps ProviderSupport) *ProviderBuilder {
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
	providerForSkillOfModule := make(map[string]map[string]Provider)
	statusForModule := make(map[string]ProviderStatus)
	for moduleName, providerSupport := range p.providers {

		for name, provider := range providerSupport.Provider() {
			cps, has := providerForSkillOfModule[name]
			if !has {
				providerForSkillOfModule[name] = make(map[string]Provider)
				cps = providerForSkillOfModule[name]
			}
			cps[moduleName] = provider
		}
		statusForModule[moduleName] = providerSupport
	}
	providersForSkill := make(map[string]Provider)
	for skill, cps := range providerForSkillOfModule {
		providersForSkill[skill] = NewProviderMultiple(cps)
	}
	return NewProvidersManager(providersForSkill, statusForModule)
}
