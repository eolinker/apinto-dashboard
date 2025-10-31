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
	_ mpm3.IAccessService = (*AccessService)(nil)
)

type AccessCache cache.IRedisCacheNoKey[*model.Access]
type AccessService struct {
	store store.IPluginAccessStore
	cache AccessCache
}

func (s *AccessService) Clean(ctx context.Context) {
	_ = s.cache.Delete(ctx)
}

func NewAccessService() mpm3.IAccessService {
	s := &AccessService{}
	s.cache = cache.CreateRedisCacheNoKey[*model.Access](time.Minute, "mpm3:access-list")
	bean.Autowired(&s.store)

	return s
}

func (s *AccessService) Save(ctx context.Context, plugin int, as []pm3.PAccess) error {

	ens := make([]*entry.Access, 0, len(as))
	for _, mi := range as {
		ens = append(ens, &entry.Access{
			Id:     0,
			Plugin: plugin,
			Name:   mi.Name,
			CName:  mi.Cname,
			Module: mi.Module,
			Depend: mi.Depend,
		})
	}
	err := s.store.Save(ctx, plugin, ens)
	if err != nil {
		return err

	}
	_ = s.cache.Delete(ctx)
	return nil
}

func (s *AccessService) GetEnable(ctx context.Context) []*model.Access {
	modules, err := s.cache.GetAll(ctx)
	if err != nil {
		log.Error("get access list of enable cache:", err)
	}
	if len(modules) > 0 {

		return modules
	}

	ens := s.store.GetEnable(ctx)
	ms := make([]*model.Access, 0, len(ens))
	for _, e := range ens {
		ms = append(ms, &model.Access{

			Name:   e.Name,
			CName:  e.CName,
			Module: e.Module,
			Depend: e.Depend,
		})
	}

	_ = s.cache.SetAll(ctx, ms)
	return ms
}
