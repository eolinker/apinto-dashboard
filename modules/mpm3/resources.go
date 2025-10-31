package mpm3

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/mpm3/model"
)

type IPluginResources interface {
	ICon() ([]byte, bool)
	RM() ([]byte, bool)
	ReadMe(name string) ([]byte, bool)
	Resources(path string) ([]byte, bool)
}
type IResourcesService interface {
	Save(ctx context.Context, id int, uuid string, resource *model.PluginResources) error
	Get(ctx context.Context, uuid string) (IPluginResources, error)
	Delete(ctx context.Context, ids ...int) error
}
