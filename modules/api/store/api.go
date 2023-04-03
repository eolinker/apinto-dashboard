package api_store

import (
	"context"
	"fmt"
	"github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"github.com/eolinker/apinto-dashboard/store"
	"strings"
)

type IAPIStore interface {
	store.IBaseStore[api_entry.API]
	GetListByRequestPath(ctx context.Context, namespaceID int, requestPath string) ([]*api_entry.API, error)
	GetListPageByGroupIDs(ctx context.Context, namespaceID, pageNum, pageSize int, groupIDs, searchSources []string, searchName string) ([]*api_entry.API, int, error)
	GetListByGroupID(ctx context.Context, namespaceID int, groupIDs string) ([]*api_entry.API, error)
	GetCountByGroupID(ctx context.Context, namespaceID int, groupIDs string) (int64, error)
	GetListByName(ctx context.Context, namespaceID int, name string) ([]*api_entry.API, error)
	GetByUUID(ctx context.Context, namespaceID int, uuid string) (*api_entry.API, error)
	GetByUUIDs(ctx context.Context, namespaceID int, uuids []string) ([]*api_entry.API, error)
	GetByPath(ctx context.Context, namespaceID int, path string) ([]*api_entry.API, error)
	GetByIds(ctx context.Context, namespaceID int, ids []int) ([]*api_entry.API, error)
	GetListAll(ctx context.Context, namespaceID int) ([]*api_entry.API, error)
	GetSourceList(ctx context.Context) ([]*api_entry.APISource, error)
}

type apiStore struct {
	*store.BaseStore[api_entry.API]
}

func NewAPIStore(db store.IDB) IAPIStore {
	return &apiStore{BaseStore: store.CreateStore[api_entry.API](db)}
}

func (a *apiStore) GetListByRequestPath(ctx context.Context, namespaceID int, requestPath string) ([]*api_entry.API, error) {
	return a.ListQuery(ctx, "`namespace` = ? and `request_path` = ?", []interface{}{namespaceID, requestPath}, "")
}

func (a *apiStore) GetListPageByGroupIDs(ctx context.Context, namespaceID, pageNum, pageSize int, groupIDs, searchSources []string, searchName string) ([]*api_entry.API, int, error) {
	apis := make([]*api_entry.API, 0)
	db := a.DB(ctx).Where("`namespace` = ?", namespaceID)
	count := int64(0)
	if len(groupIDs) > 0 {
		db = db.Where("`group_uuid` in (?)", groupIDs)
	}

	if len(searchSources) > 0 {
		db = db.Where(fmt.Sprintf("(`source_type`,`source_id`,`source_label`) in ( %s )", strings.Join(searchSources, ",")))
	}

	if searchName != "" {
		db = db.Where("`name` like ?", "%"+searchName+"%")
	}
	if pageNum > 0 && pageSize > 0 {
		err := db.Model(apis).Count(&count).Order("update_time DESC").Limit(pageSize).Offset(store.PageIndex(pageNum, pageSize)).Find(&apis).Error
		if err != nil {
			return nil, 0, err
		}
	} else {
		err := db.Order("update_time DESC").Find(&apis).Error
		if err != nil {
			return nil, 0, err
		}
	}

	return apis, int(count), nil
}

func (a *apiStore) GetCountByGroupID(ctx context.Context, namespaceID int, groupID string) (int64, error) {
	count := int64(0)
	err := a.DB(ctx).Where("`namespace` = ? and `group_uuid` = ?", namespaceID, groupID).Model(api_entry.API{}).Count(&count).Error
	return count, err
}

func (a *apiStore) GetListByGroupID(ctx context.Context, namespaceID int, groupID string) ([]*api_entry.API, error) {
	return a.ListQuery(ctx, "`namespace` = ? and `group_uuid` = ?", []interface{}{namespaceID, groupID}, "")
}

func (a *apiStore) GetListByName(ctx context.Context, namespaceID int, name string) ([]*api_entry.API, error) {
	apis := make([]*api_entry.API, 0)
	db := a.DB(ctx).Where("`namespace` = ?", namespaceID)
	if name != "" {
		db = db.Where("`name` like ?", "%"+name+"%")
	}
	err := db.Order("update_time DESC").Find(&apis).Error
	if err != nil {
		return nil, err
	}
	return apis, nil
}

func (a *apiStore) GetByUUID(ctx context.Context, namespaceID int, uuid string) (*api_entry.API, error) {
	return a.FirstQuery(ctx, "`namespace` = ? and `uuid` = ?", []interface{}{namespaceID, uuid}, "")
}

func (a *apiStore) GetByUUIDs(ctx context.Context, namespaceID int, uuids []string) ([]*api_entry.API, error) {
	apis := make([]*api_entry.API, 0)
	err := a.DB(ctx).Where("`namespace` = ? and `uuid` in (?)", namespaceID, uuids).Find(&apis).Error
	return apis, err
}

func (a *apiStore) GetByPath(ctx context.Context, namespaceID int, path string) ([]*api_entry.API, error) {
	apis := make([]*api_entry.API, 0)
	err := a.DB(ctx).Where("`namespace` = ? and `request_path_label` = ?", namespaceID, path).Find(&apis).Error
	return apis, err
}

func (a *apiStore) GetByIds(ctx context.Context, namespaceID int, ids []int) ([]*api_entry.API, error) {
	apis := make([]*api_entry.API, 0)
	err := a.DB(ctx).Where("`namespace` = ? and `id` in (?)", namespaceID, ids).Find(&apis).Error
	return apis, err
}

func (a *apiStore) GetListAll(ctx context.Context, namespaceID int) ([]*api_entry.API, error) {
	apis := make([]*api_entry.API, 0)
	err := a.DB(ctx).Where("`namespace` = ?", namespaceID).Find(&apis).Error
	return apis, err
}

func (a *apiStore) GetSourceList(ctx context.Context) ([]*api_entry.APISource, error) {
	sourceList := make([]*api_entry.APISource, 0, 2)

	rows, err := a.DB(ctx).Raw("SELECT `source_type`,`source_id`,`source_label` FROM `api` GROUP BY `source_type`,`source_id`,`source_label`").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		sourceType, sourceLabel string
		sourceId                int
	)
	for rows.Next() {
		_ = rows.Scan(&sourceType, &sourceId, &sourceLabel)
		sourceList = append(sourceList, &api_entry.APISource{
			SourceType:  sourceType,
			SourceID:    sourceId,
			SourceLabel: sourceLabel,
		})
	}

	return sourceList, err
}
