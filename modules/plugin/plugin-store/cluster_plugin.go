package plugin_store

import (
	"context"
	plugin_entry "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"gorm.io/gorm/clause"
)

var (
	_ IClusterPluginStore = (*clusterPluginStore)(nil)
)

type IClusterPluginStore interface {
	store.IBaseStore[plugin_entry.ClusterPlugin]
	GetByClusterIds(ctx context.Context, ids ...int) ([]*plugin_entry.ClusterPlugin, error)
	UpdatePlugins(ctx context.Context, list []*plugin_entry.ClusterPlugin) error
	GetPluginsByPluginName(ctx context.Context, namespaceID, pluginName string) ([]*plugin_entry.ClusterPlugin, error)
	GetClusterPluginByClusterIDByPluginName(ctx context.Context, clusterID int, pluginName string) (*plugin_entry.ClusterPlugin, error)
	DeleteClusterPluginByIDs(ctx context.Context, clusterPluginIDs ...int) error
}

type clusterPluginStore struct {
	*store.BaseStore[plugin_entry.ClusterPlugin]
}

func newClusterPluginStore(db store.IDB) IClusterPluginStore {
	return &clusterPluginStore{BaseStore: store.CreateStore[plugin_entry.ClusterPlugin](db)}
}

func (c *clusterPluginStore) GetByClusterIds(ctx context.Context, clusterIds ...int) ([]*plugin_entry.ClusterPlugin, error) {
	return c.ListQuery(ctx, "`cluster` in (?)", []interface{}{clusterIds}, "")
}

func (c *clusterPluginStore) UpdatePlugins(ctx context.Context, list []*plugin_entry.ClusterPlugin) error {

	columns := make([]clause.Column, 0)
	columns = append(columns, clause.Column{
		Name: "cluster",
	}, clause.Column{
		Name: "plugin_name",
	})

	for _, val := range list {
		err := c.DB(ctx).Clauses(
			clause.OnConflict{
				Columns:   columns,
				UpdateAll: true,
			},
		).Create(val).Error
		if err != nil {
			return err
		}
	}
	return nil

}

func (c *clusterPluginStore) GetPluginsByPluginName(ctx context.Context, namespaceID, pluginName string) ([]*plugin_entry.ClusterPlugin, error) {
	return c.ListQuery(ctx, "namespace = ? AND plugin_name = ?", []interface{}{namespaceID, pluginName}, "")
}

func (c *clusterPluginStore) GetClusterPluginByClusterIDByPluginName(ctx context.Context, clusterID int, pluginName string) (*plugin_entry.ClusterPlugin, error) {
	db := c.DB(ctx)
	variable := &plugin_entry.ClusterPlugin{}
	if err := db.Where("cluster = ? AND plugin_name = ?", clusterID, pluginName).Take(variable).Error; err != nil {
		return nil, err
	}
	return variable, nil
}

func (c *clusterPluginStore) DeleteClusterPluginByIDs(ctx context.Context, clusterPluginIDs ...int) error {
	db := c.DB(ctx)
	db.Where("`id` in (?)", []interface{}{clusterPluginIDs}).Delete(&plugin_entry.ClusterPlugin{})

	return nil
}
