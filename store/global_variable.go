package store

import (
	"context"
	"github.com/eolinker/apinto-dashboard/entry/page-entry"
	"github.com/eolinker/apinto-dashboard/entry/variable-entry"
	"gorm.io/gorm"
)

type IGlobalVariableStore interface {
	IBaseStore[variable_entry.Variables]
	GetList(ctx context.Context, pageNum, pageSize, namespaceID int, searchName string) ([]*variable_entry.Variables, int, error)
	GetGlobalVariableIDByKey(ctx context.Context, namespaceID int, key string) (*variable_entry.Variables, error)
	GetGlobalVariableByKeys(ctx context.Context, namespaceID int, keys []string) ([]*variable_entry.Variables, error)
}

type globalVariableStore struct {
	*BaseStore[variable_entry.Variables]
}

func newGlobalVariableStore(db IDB) IGlobalVariableStore {
	return &globalVariableStore{BaseStore: CreateStore[variable_entry.Variables](db)}
}

func (g *globalVariableStore) GetList(ctx context.Context, pageNum, pageSize, namespaceID int, searchName string) ([]*variable_entry.Variables, int, error) {
	variables := make([]*variable_entry.Variables, 0)
	db := g.DB(ctx).Where("namespace = ?", namespaceID)

	count := int64(0)
	if searchName != "" {
		db = db.Where("`key` like ?", "%"+searchName+"%")
	}

	if pageNum > 0 && pageSize > 0 {
		err := db.Model(variables).Count(&count).Order("create_time asc").Limit(pageSize).Offset(page_entry.PageIndex(pageNum, pageSize)).Find(&variables).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := db.Order("create_time asc").Find(&variables).Error
		if err != nil {
			return nil, 0, err
		}
	}
	return variables, int(count), nil
}

func (g *globalVariableStore) GetGlobalVariableIDByKey(ctx context.Context, namespaceID int, key string) (*variable_entry.Variables, error) {
	variable := &variable_entry.Variables{}
	if err := g.DB(ctx).Where("namespace = ? AND `key` = ?", namespaceID, key).Take(variable).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return variable, nil
}

func (g *globalVariableStore) GetGlobalVariableByKeys(ctx context.Context, namespaceID int, keys []string) ([]*variable_entry.Variables, error) {
	variables := make([]*variable_entry.Variables, 0, len(keys))
	if err := g.DB(ctx).Where("namespace = ? AND `key` in ?", namespaceID, keys).Find(&variables).Error; err != nil {
		return nil, err
	}
	return variables, nil
}
