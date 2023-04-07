package group

import (
	"context"
	"github.com/eolinker/apinto-dashboard/modules/group/group-dto"
	"github.com/eolinker/apinto-dashboard/modules/group/group-entry"
	"github.com/eolinker/apinto-dashboard/modules/group/group-model"
)

type ICommonGroupService interface {
	CreateGroup(ctx context.Context, namespaceId int, operator int, groupType, tagName, groupName, uuid, parentUuid string) error
	UpdateGroup(ctx context.Context, namespaceId int, operator int, groupType, name, uuid string) error
	DeleteGroup(ctx context.Context, namespaceId int, operator int, groupType, uuid string) error
	GroupList(ctx context.Context, namespaceId int, groupType, tagName, parentUuid, queryName string) (*group_model.CommonGroupRoot, []*group_model.CommonGroupApi, error)
	GroupListAll(ctx context.Context, namespaceId int, groupType, tagName string) ([]*group_entry.CommonGroup, error)
	GroupUUIDS(ctx context.Context, namespaceId int, groupType, tagName, parentUuid string) ([]string, error)
	GroupSort(ctx context.Context, namespaceId int, groupType, tagName string, input *group_dto.CommGroupSortInput) error
	ToGroupRoot(ctx context.Context, namespaceId int, queryUUid string, groups []*group_entry.CommonGroup, uuidMaps map[string]string) *group_model.CommonGroupRoot
	ParentGroupV2(parentUUID string, groupMaps map[string]*group_entry.CommonGroup, groupIdMaps map[int]*group_entry.CommonGroup, outMapUUID map[string]string)
	ParentGroupName(uuid string, groupMaps map[string]*group_entry.CommonGroup, groupIdMaps map[int]*group_entry.CommonGroup, nameList *[]string)
	SubGroupUUIDS(groups map[int][]*group_entry.CommonGroup, parentGroup *group_model.CommonGroup, list *[]string)
	IsGroupExist(ctx context.Context, uuid string) (bool, error)
	CheckGroupNameReduplicated(ctx context.Context, groupName string, parentID int) (bool, error)
	GetGroupInfo(ctx context.Context, uuid string) (*group_entry.CommonGroup, error)
}
