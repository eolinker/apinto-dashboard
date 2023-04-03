package namespace

import (
	"github.com/eolinker/apinto-dashboard/modules/namespace/namespace-model"
)

type INamespaceService interface {
	GetByName(namespace string) (*namespace_model.Namespace, error)
	GetById(namespaceId int) (*namespace_model.Namespace, error)
	GetAll() ([]*namespace_model.Namespace, error)
}
