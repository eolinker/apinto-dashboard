package service

import (
	"context"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
)

var _ INamespaceService = (*namespaceService)(nil)

type INamespaceService interface {
	GetByName(namespace string) (*model.Namespace, error)
	GetById(namespaceId int) (*model.Namespace, error)
	GetAll() ([]*model.Namespace, error)
}

type namespaceService struct {
	namespaceStore store.INamespaceStore
}

func newNamespaceService() INamespaceService {
	n := &namespaceService{}
	bean.Autowired(&n.namespaceStore)
	return n
}

func (n *namespaceService) GetByName(name string) (*model.Namespace, error) {

	namespace, err := n.namespaceStore.GetByName(context.TODO(), name)
	if err != nil {
		return nil, err
	}

	return &model.Namespace{Namespace: namespace}, nil
}

func (n *namespaceService) GetById(id int) (*model.Namespace, error) {

	namespace, err := n.namespaceStore.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}
	return &model.Namespace{Namespace: namespace}, nil
}

func (n *namespaceService) GetAll() ([]*model.Namespace, error) {
	list, err := n.namespaceStore.GetAll(context.TODO())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Namespace, 0, len(list))

	for _, namespace := range list {

		result = append(result, &model.Namespace{Namespace: namespace})
	}

	return result, nil
}
