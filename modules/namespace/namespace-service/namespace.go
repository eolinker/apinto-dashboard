package namespace_service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-model"
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-store"
	"github.com/eolinker/eosc/common/bean"
)

var _ namespace.INamespaceService = (*namespaceService)(nil)

type namespaceService struct {
	namespaceStore namespace_store.INamespaceStore
}

func newNamespaceService() namespace.INamespaceService {
	n := &namespaceService{}
	bean.Autowired(&n.namespaceStore)
	return n
}

func (n *namespaceService) GetByName(name string) (*namespace_model.Namespace, error) {

	namespaceInfo, err := n.namespaceStore.GetByName(context.TODO(), name)
	if err != nil {
		return nil, err
	}

	return &namespace_model.Namespace{Namespace: namespaceInfo}, nil
}

func (n *namespaceService) GetById(id int) (*namespace_model.Namespace, error) {

	namespaceInfo, err := n.namespaceStore.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	return &namespace_model.Namespace{Namespace: namespaceInfo}, nil
}

func (n *namespaceService) GetAll() ([]*namespace_model.Namespace, error) {
	list, err := n.namespaceStore.GetAll(context.TODO())
	if err != nil {
		return nil, err
	}

	result := make([]*namespace_model.Namespace, 0, len(list))

	for _, namespaceInfo := range list {

		result = append(result, &namespace_model.Namespace{Namespace: namespaceInfo})
	}

	return result, nil
}
