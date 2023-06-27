package namespace_service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-model"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-store"
	"github.com/eolinker/eosc/common/bean"
)

var _ namespace.INamespaceService = (*namespaceService)(nil)

type namespaceService struct {
	namespaceStore       namespace_store.INamespaceStore
	namespaceCacheByName cache.IRedisCache[namespace_model.Namespace, string]
	namespaceCacheById   cache.IRedisCache[namespace_model.Namespace, int]
	namespaceCacheAll    cache.IRedisCacheNoKey[namespace_model.Namespace]
}

func newNamespaceService() *namespaceService {
	n := &namespaceService{}
	bean.Autowired(&n.namespaceStore)
	return n
}

func (n *namespaceService) GetByName(name string) (*namespace_model.Namespace, error) {
	namespaceInfo, err := n.namespaceCacheByName.Get(context.Background(), name)
	if namespaceInfo != nil {
		return namespaceInfo, nil
	}
	namespaceEntry, err := n.namespaceStore.GetByName(context.TODO(), name)
	if err != nil {
		return nil, err
	}
	namespaceInfo = &namespace_model.Namespace{
		Id:         namespaceEntry.Id,
		Name:       namespaceEntry.Name,
		CreateTime: namespaceEntry.CreateTime}
	n.namespaceCacheByName.Set(context.Background(), name, namespaceInfo)
	n.namespaceCacheById.Set(context.Background(), namespaceInfo.Id, namespaceInfo)
	return namespaceInfo, nil
}

func (n *namespaceService) GetById(id int) (*namespace_model.Namespace, error) {
	namespaceInfo, err := n.namespaceCacheById.Get(context.Background(), id)
	if namespaceInfo != nil {
		return namespaceInfo, nil
	}
	namespaceEntry, err := n.namespaceStore.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	namespaceInfo = &namespace_model.Namespace{
		Id:         namespaceEntry.Id,
		Name:       namespaceEntry.Name,
		CreateTime: namespaceEntry.CreateTime,
	}
	n.namespaceCacheById.Set(context.Background(), id, namespaceInfo)
	n.namespaceCacheByName.Set(context.Background(), namespaceInfo.Name, namespaceInfo)
	return namespaceInfo, nil
}

func (n *namespaceService) GetAll() ([]*namespace_model.Namespace, error) {
	rs, err := n.namespaceCacheAll.GetAll(context.Background())
	if rs != nil {
		return rs, nil
	}
	list, err := n.namespaceStore.GetAll(context.TODO())
	if err != nil {
		return nil, err
	}

	result := make([]*namespace_model.Namespace, 0, len(list))

	for _, namespaceEntity := range list {

		result = append(result, &namespace_model.Namespace{
			Id:         namespaceEntity.Id,
			Name:       namespaceEntity.Name,
			CreateTime: namespaceEntity.CreateTime,
		})
	}
	n.namespaceCacheAll.SetAll(context.Background(), result)
	return result, nil
}
