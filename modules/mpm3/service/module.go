package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/store"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"time"
)

var (
	_ mpm3.IModuleService = (*ModuleService)(nil)
)

type ModuleServiceCache cache.IRedisCacheNoKey[*model.Module]
type ModuleService struct {
	store store.IModuleStore
	cache ModuleServiceCache
}

func (m *ModuleService) Clean(ctx context.Context) {
	_ = m.cache.Delete(ctx)
}

func NewModuleService() mpm3.IModuleService {
	s := &ModuleService{}
	s.cache = cache.CreateRedisCacheNoKey[*model.Module](time.Minute, "mpm3:module-list")
	bean.Autowired(&s.store)
	return s
}

func (m *ModuleService) Save(ctx context.Context, plugin int, module []pm3.PModule) error {

	ens := make([]*entry.Module, 0, len(module))
	for _, mi := range module {
		ens = append(ens, &entry.Module{
			Id:         0,
			Plugin:     plugin,
			Navigation: mi.Navigation,
			Name:       mi.Name,
			CName:      mi.Cname,
			Router:     mi.Router,
		})
	}
	err := m.store.Save(ctx, plugin, ens)
	if err != nil {
		return err

	}
	_ = m.cache.Delete(ctx)
	return nil
}

func (m *ModuleService) GetEnable(ctx context.Context) []*model.Module {
	modules, err := m.cache.GetAll(ctx)
	if err != nil {
		log.Error("get module enable cache:", err)
	}
	if len(modules) > 0 {

		return modules
	}

	ens := m.store.GetEnable(ctx)
	ms := make([]*model.Module, 0, len(ens))
	for _, e := range ens {
		ms = append(ms, &model.Module{

			Navigation: e.Navigation,
			Name:       e.Name,
			CName:      e.CName,
			Router:     e.Router,
		})
	}

	m.cache.SetAll(ctx, ms)
	return ms
}
