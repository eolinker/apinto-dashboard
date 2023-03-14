package group_store

import (
	"context"
	group_entry "github.com/eolinker/apinto-dashboard/entry/group-entry"
	"github.com/eolinker/apinto-dashboard/store"
)

type ICommonGroupStore interface {
	store.IBaseStore[group_entry.CommonGroup]
	GetByUUID(ctx context.Context, uuid string) (*group_entry.CommonGroup, error)
	GetByUUIDS(ctx context.Context, uuid []string) ([]*group_entry.CommonGroup, error)
	GetByName(ctx context.Context, namespaceId int, name string) ([]*group_entry.CommonGroup, error)
	GetMaxSort(ctx context.Context, namespaceId int, groupType string, tagId, parentId int) (int, error)
	UpdateSort(ctx context.Context, list []*group_entry.CommonGroup) error
	GetByParentId(ctx context.Context, namespaceId int, groupType string, tagId, parentId int) ([]*group_entry.CommonGroup, error)
	GetList(ctx context.Context, namespace int, groupType string, tagId int) ([]*group_entry.CommonGroup, error)
	GetByNameParentID(ctx context.Context, groupName string, parentID int) ([]*group_entry.CommonGroup, error)
}

type commonGroupStore struct {
	*store.BaseStore[group_entry.CommonGroup]
}

func newCommonGroupStore(db store.IDB) ICommonGroupStore {
	return &commonGroupStore{BaseStore: store.CreateStore[group_entry.CommonGroup](db)}
}

func (c *commonGroupStore) GetMaxSort(ctx context.Context, namespaceId int, groupType string, tagId, parentId int) (int, error) {

	sort := 0

	db := c.DB(ctx).Table("common_group").Select("IFNULL(MAX(`sort`),0) AS `sort`") //IFNULL MAX 为了处理 N/A默认值的问题
	err := db.Where("`namespace` = ? and  `type` = ? and `tag` = ? and `parent_id` = ?", namespaceId, groupType, tagId, parentId).Order("`sort` desc").Limit(1).Row().Scan(&sort)
	if err != nil {
		return 0, err
	}

	return sort, nil
}

func (c *commonGroupStore) UpdateSort(ctx context.Context, list []*group_entry.CommonGroup) error {

	for _, group := range list {

		if err := c.DB(ctx).Exec("update `common_group` set `sort` = ?,`parent_id` = ? where id = ?", group.Sort, group.ParentId, group.Id).Error; err != nil {
			return err
		}
	}

	return nil
}

func (c *commonGroupStore) GetByUUID(ctx context.Context, uuid string) (*group_entry.CommonGroup, error) {
	return c.FirstQuery(ctx, "`uuid` = ?", []interface{}{uuid}, "")
}

func (c *commonGroupStore) GetByName(ctx context.Context, namespaceId int, name string) ([]*group_entry.CommonGroup, error) {
	return c.ListQuery(ctx, "`namespace` =? and `name` like ?", []interface{}{namespaceId, "%" + name + "%"}, "")
}

func (c *commonGroupStore) GetByUUIDS(ctx context.Context, uuid []string) ([]*group_entry.CommonGroup, error) {
	return c.ListQuery(ctx, "`uuid` in (?)", []interface{}{uuid}, "")
}

func (c *commonGroupStore) GetByParentId(ctx context.Context, namespace int, groupType string, tagId, parentId int) ([]*group_entry.CommonGroup, error) {
	return c.ListQuery(ctx, "`namespace` = ? and `type` = ? and `tag` = ? and `parent_id` = ?", []interface{}{namespace, groupType, tagId, parentId}, "`sort` asc")
}

func (c *commonGroupStore) GetList(ctx context.Context, namespace int, groupType string, tagId int) ([]*group_entry.CommonGroup, error) {
	return c.ListQuery(ctx, "`namespace` = ? and `type` = ? and `tag` = ?", []interface{}{namespace, groupType, tagId}, "`sort` asc")
}

func (c *commonGroupStore) GetByNameParentID(ctx context.Context, groupName string, parentID int) ([]*group_entry.CommonGroup, error) {
	return c.ListQuery(ctx, "`name` = ? and `parent_id` = ? ", []interface{}{groupName, parentID}, "")
}
