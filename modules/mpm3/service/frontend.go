package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/store"
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"time"
)

var (
	_ mpm3.IFrontendService = (*FrontendService)(nil)
)

type FrontendService struct {
	store store.IPluginFrontendStore
	cache cache.IRedisCacheNoKey[pm3.PFrontend]
}

func (s *FrontendService) Clean(ctx context.Context) {
	_ = s.cache.Delete(ctx)
}

func NewFrontendService() mpm3.IFrontendService {
	s := &FrontendService{}
	s.cache = cache.CreateRedisCacheNoKey[pm3.PFrontend](time.Minute, "mpm3:FrontendAsset-list")
	bean.Autowired(&s.store)

	return s
}

func (s *FrontendService) Save(ctx context.Context, plugin int, as []pm3.PFrontend) error {

	ens := make([]*entry.Frontend, 0, len(as))
	for _, mi := range as {
		ens = append(ens, &entry.Frontend{
			Id:      0,
			Plugin:  plugin,
			Content: mi,
		})
	}
	err := s.store.Save(ctx, plugin, ens)
	if err != nil {
		return err

	}
	_ = s.cache.Delete(ctx)
	return nil
}

func (s *FrontendService) GetEnable(ctx context.Context) []pm3.PFrontend {
	timeout, cancelFunc := context.WithTimeout(ctx, time.Second)
	defer cancelFunc()
	modules, err := s.cache.GetAll(timeout)
	if err != nil {
		log.Error("get plugin frontend list of enable cache:", err)
	}
	if len(modules) > 0 {
		return modules
	}

	ens := s.store.GetEnable(ctx)
	ms := make([]pm3.PFrontend, 0, len(ens))
	for _, e := range ens {
		ms = append(ms, e.Content)
	}

	_ = s.cache.SetAll(ctx, ms)
	return ms
}
