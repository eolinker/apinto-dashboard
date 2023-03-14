package namespace_service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/model/namespace-model"
	"github.com/eolinker/apinto-dashboard/store/namespace-store"
	"github.com/eolinker/eosc/common/bean"
)

var _ INamespaceService = (*namespaceService)(nil)

type INamespaceService interface {
	GetByName(namespace string) (*namespace_model.Namespace, error)
	GetById(namespaceId int) (*namespace_model.Namespace, error)
	GetAll() ([]*namespace_model.Namespace, error)
}

type namespaceService struct {
	namespaceStore namespace_store.INamespaceStore
}

func newNamespaceService() INamespaceService {
	n := &namespaceService{}
	bean.Autowired(&n.namespaceStore)
	return n
}

func (n *namespaceService) GetByName(name string) (*namespace_model.Namespace, error) {

	namespace, err := n.namespaceStore.GetByName(context.TODO(), name)
	if err != nil {
		return nil, err
	}

	return &namespace_model.Namespace{Namespace: namespace}, nil
}

func (n *namespaceService) GetById(id int) (*namespace_model.Namespace, error) {

	namespace, err := n.namespaceStore.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	return &namespace_model.Namespace{Namespace: namespace}, nil
}

func (n *namespaceService) GetAll() ([]*namespace_model.Namespace, error) {
	list, err := n.namespaceStore.GetAll(context.TODO())
	if err != nil {
		return nil, err
	}

	result := make([]*namespace_model.Namespace, 0, len(list))

	for _, namespace := range list {

		result = append(result, &namespace_model.Namespace{Namespace: namespace})
	}

	return result, nil
}
