package variable_store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_ IClusterVariableStore = (*clusterVariableStore)(nil)
)

type IClusterVariableStore interface {
	store.IBaseStore[variable_entry.ClusterVariable]
	GetByClusterIds(ctx context.Context, ids ...int) ([]*variable_entry.ClusterVariable, error)
	UpdateVariables(ctx context.Context, list []*variable_entry.ClusterVariable) error
	GetVariablesByGlobalVariableID(ctx context.Context, namespaceID, globalVariableID int) ([]*variable_entry.ClusterVariable, error)
	GetClusterVariableByClusterIDByGlobalID(ctx context.Context, clusterID, variableID int) (*variable_entry.ClusterVariable, error)
}

type clusterVariableStore struct {
	*store.BaseStore[variable_entry.ClusterVariable]
}

func newClusterVariableStore(db store.IDB) IClusterVariableStore {
	return &clusterVariableStore{BaseStore: store.CreateStore[variable_entry.ClusterVariable](db)}
}

func (e *clusterVariableStore) GetByClusterIds(ctx context.Context, clusterIds ...int) ([]*variable_entry.ClusterVariable, error) {
	return e.ListQuery(ctx, "`cluster` in (?)", []interface{}{clusterIds}, "")
}

func (e *clusterVariableStore) UpdateVariables(ctx context.Context, list []*variable_entry.ClusterVariable) error {

	columns := make([]clause.Column, 0)
	columns = append(columns, clause.Column{
		Name: "cluster",
	}, clause.Column{
		Name: "variable",
	})

	for _, val := range list {
		err := e.DB(ctx).Clauses(
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

func (e *clusterVariableStore) GetVariablesByGlobalVariableID(ctx context.Context, namespaceID, globalVariableID int) ([]*variable_entry.ClusterVariable, error) {
	return e.ListQuery(ctx, "namespace = ? AND variable = ?", []interface{}{namespaceID, globalVariableID}, "")
}

func (e *clusterVariableStore) GetClusterVariableByClusterIDByGlobalID(ctx context.Context, clusterID, variableID int) (*variable_entry.ClusterVariable, error) {
	db := e.DB(ctx)
	variable := &variable_entry.ClusterVariable{}
	if err := db.Where("cluster = ? AND variable = ?", clusterID, variableID).Take(variable).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return variable, nil
}
