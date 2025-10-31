package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/mpm3"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/entry"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/store"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/common/bean"
)

var (
	_ mpm3.IResourcesService = (*ResourcesService)(nil)
)

type resourcesStamp struct {
	Stamp int64  `json:"stamp"`
	Hash  string `json:"hash"`
}
type resourcesST struct {
	resources *model.PluginResources
	stamp     int64
}
type ResourcesService struct {
	store store.IPluginResources

	cache     cache.IRedisCache[resourcesStamp, string]
	localData eosc.Untyped[string, *resourcesST]
}

func NewResourcesService() mpm3.IResourcesService {

	p := &ResourcesService{
		localData: eosc.BuildUntyped[string, *resourcesST](),
	}
	bean.Autowired(&p.store)
	p.cache = cache.CreateRedisCache[resourcesStamp](time.Minute*30, func(k string) string {
		return fmt.Sprintf("mpm3:resources:%s", k)
	})

	return p
}

func (r *ResourcesService) Delete(ctx context.Context, ids ...int) error {
	_, err := r.store.Delete(ctx, ids...)
	return err

}

func (r *ResourcesService) Save(ctx context.Context, id int, uuid string, resource *model.PluginResources) error {
	data, err := json.Marshal(resource)
	if err != nil {
		return err
	}

	err = r.store.Save(ctx, &entry.Resources{
		ID:        id,
		Uuid:      uuid,
		Resources: data,
	})
	if err != nil {
		return err
	}

	r.localData.Del(uuid)
	r.cache.Delete(ctx, uuid)
	return nil
}

func (r *ResourcesService) Get(ctx context.Context, uuid string) (mpm3.IPluginResources, error) {

	stamp, err := r.cache.Get(ctx, uuid)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err != nil || stamp == nil {
		stamp = &resourcesStamp{Stamp: time.Now().Unix()}
		r.cache.Set(ctx, uuid, stamp)
	}
	data, has := r.localData.Get(uuid)
	if has && data.stamp == stamp.Stamp {
		if data.resources == nil {
			return nil, ErrModulePluginNotFound
		}
		return data.resources, nil
	}
	if data == nil {
		data = &resourcesST{
			resources: nil,
			stamp:     stamp.Stamp,
		}
	}
	defer func() {
		r.localData.Set(uuid, data)
	}()
	en, err := r.store.First(ctx, map[string]interface{}{"uuid": uuid})
	if err != nil {
		return nil, err
	}

	res := &model.PluginResources{}
	err = json.Unmarshal(en.Resources, res)
	if err != nil {
		return nil, err
	}
	data.resources = res
	return res, nil
}
