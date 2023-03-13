package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"time"
)

const (
	apiName     = "api"
	serviceName = "service"
)

type ICommonGroupService interface {
	CreateGroup(ctx context.Context, namespaceId int, operator int, groupType, tagName, groupName, uuid, parentUuid string) error
	UpdateGroup(ctx context.Context, namespaceId int, operator int, groupType, name, uuid string) error
	DeleteGroup(ctx context.Context, namespaceId int, operator int, groupType, uuid string) error
	GroupList(ctx context.Context, namespaceId int, groupType, tagName, parentUuid, queryName string) (*model.CommonGroupRoot, []*model.CommonGroupApi, error)
	GroupListAll(ctx context.Context, namespaceId int, groupType, tagName string) ([]*entry.CommonGroup, error)
	groupUUIDS(ctx context.Context, namespaceId int, groupType, tagName, parentUuid string) ([]string, error)
	GroupSort(ctx context.Context, namespaceId int, groupType, tagName string, input *dto.CommGroupSortInput) error
	toGroupRoot(ctx context.Context, namespaceId int, queryUUid string, groups []*entry.CommonGroup, uuidMaps map[string]string) *model.CommonGroupRoot
	parentGroupV2(parentUUID string, groupMaps map[string]*entry.CommonGroup, groupIdMaps map[int]*entry.CommonGroup, outMapUUID map[string]string)
	parentGroupName(uuid string, groupMaps map[string]*entry.CommonGroup, groupIdMaps map[int]*entry.CommonGroup, nameList *[]string)
	subGroupUUIDS(groups map[int][]*entry.CommonGroup, parentGroup *model.CommonGroup, list *[]string)
	IsGroupExist(ctx context.Context, uuid string) (bool, error)
	CheckGroupNameReduplicated(ctx context.Context, groupName string, parentID int) (bool, error)
}

type commonGroupService struct {
	commonGroupStore store.ICommonGroupStore
	service          IService
	apiService       IAPIService
}

func newCommonGroupService() ICommonGroupService {
	c := &commonGroupService{}
	bean.Autowired(&c.commonGroupStore)
	bean.Autowired(&c.service)
	bean.Autowired(&c.apiService)
	return c
}

func (c *commonGroupService) GroupSort(ctx context.Context, namespaceId int, groupType, tagName string, input *dto.CommGroupSortInput) error {
	//TODO 排序需要检查排序后分组名有无重复的情况

	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return errors.New("params error")
	}
	var (
		group *entry.CommonGroup
		err   error
	)

	if input.Root != "" {
		group, err = c.commonGroupStore.GetByUUID(ctx, input.Root)
		if err != nil {
			return err
		}
	} else {
		group = new(entry.CommonGroup)
	}

	groups, err := c.commonGroupStore.GetByUUIDS(ctx, input.Items)
	if err != nil {
		return err
	}
	groupMaps := common.SliceToMap(groups, func(t *entry.CommonGroup) string {
		return t.Uuid
	})

	list := make([]*entry.CommonGroup, 0)
	for i, groupUUID := range input.Items {
		if v, ok := groupMaps[groupUUID]; ok {
			v.Sort = i + 1
			v.ParentId = group.ParentId
			list = append(list, v)
		}
	}

	return c.commonGroupStore.Transaction(ctx, func(txCtx context.Context) error {
		return c.commonGroupStore.UpdateSort(txCtx, list)
	})
}

func (c *commonGroupService) DeleteGroup(ctx context.Context, namespaceId int, operator int, groupType, uuid string) error {
	group, err := c.commonGroupStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	switch groupType {
	case apiName:

		groupUUIDS, err := c.groupUUIDS(ctx, namespaceId, groupType, apiName, uuid)
		if err != nil {
			return err
		}
		groupUUIDS = append(groupUUIDS, uuid)
		for _, groupUUID := range groupUUIDS {

			count := c.apiService.GetAPICountByGroupUUID(ctx, namespaceId, groupUUID)
			if count > 0 {
				return errors.New(fmt.Sprintf("该分组下已有API，不可删除"))
			}

			//apiInfos, err := c.apiService.GetAPIInfoByGroupUUID(ctx, namespaceId, groupUUID)
			//if err != nil {
			//	return err
			//}
			//
			//for _, info := range apiInfos {
			//	isDelete, _ := c.apiService.isApiCanDelete(ctx, info.Id)
			//	if !isDelete {
			//		return errors.New(fmt.Sprintf("该目录下{%s}已上线", info.Name))
			//	}
			//}
		}

	default:
		return errors.New("param err")
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: uuid,
		Name: group.Name,
	})

	_, err = c.commonGroupStore.Delete(ctx, group.Id)
	return err
}

func (c *commonGroupService) UpdateGroup(ctx context.Context, namespaceId int, operator int, groupType, name, uuid string) error {
	group, err := c.commonGroupStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	group.Operator = operator
	group.UpdateTime = time.Now()
	group.Name = name

	//检查修改的分组名是否有重复
	groupList, err := c.commonGroupStore.GetByNameParentID(ctx, name, group.ParentId)
	if err != nil {
		return err
	}
	for _, item := range groupList {
		//处理修改的分组名不变的情况
		if item.Uuid != group.Uuid {
			return fmt.Errorf("groupName %s is reduplicated", name)
		}
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: uuid,
		Name: name,
	})
	return c.commonGroupStore.Save(ctx, group)
}

func (c *commonGroupService) GroupListAll(ctx context.Context, namespaceId int, groupType, tagName string) ([]*entry.CommonGroup, error) {

	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return nil, errors.New("params error")
	}

	groups, err := c.commonGroupStore.GetList(ctx, namespaceId, groupType, tagId)
	if err != nil {
		return nil, err
	}

	return groups, err

}

// toGroupRoot uuidMaps为其他API关联的groupUUID以及直至顶级目录的groupUUID
func (c *commonGroupService) toGroupRoot(ctx context.Context, namespaceId int, queryUUid string, groups []*entry.CommonGroup, uuidMaps map[string]string) *model.CommonGroupRoot {
	parentList := make([]*model.CommonGroup, 0)
	name := ""
	for _, group := range groups {
		if queryUUid != "" {
			if group.Uuid == queryUUid {
				name = group.Name
				parentList = append(parentList, &model.CommonGroup{
					Group: group,
				})
			}
		} else {
			if group.ParentId == 0 {
				if len(uuidMaps) > 0 {
					if _, ok := uuidMaps[group.Uuid]; !ok {
						continue
					}
				}
				parentList = append(parentList, &model.CommonGroup{
					Group: group,
				})
			}
		}
	}

	array := common.SliceToMapArray(groups, func(t *entry.CommonGroup) int {
		return t.ParentId
	})

	for _, group := range parentList {
		//采用递归查询子目录
		c.subGroupV2(ctx, namespaceId, array, uuidMaps, group)
	}

	resRoot := &model.CommonGroupRoot{
		UUID:        queryUUid,
		Name:        name,
		CommonGroup: parentList,
	}
	return resRoot
}

func (c *commonGroupService) GroupList(ctx context.Context, namespaceId int, groupType, tagName, parentUuid, queryName string) (*model.CommonGroupRoot, []*model.CommonGroupApi, error) {

	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return nil, nil, errors.New("params error")
	}

	var err error
	uuids := ""
	name := ""
	//传了某个目录的UUID 则只查该目录下的所有目录
	parentId := 0
	if len(parentUuid) > 0 {
		group, err := c.commonGroupStore.GetByUUID(ctx, parentUuid)
		if err != nil {
			return nil, nil, err
		}
		parentId = group.Id
		uuids = group.Uuid
		name = group.Name
	}

	apis := make([]*model.CommonGroupApi, 0)
	//查询API的根目录列表
	if groupType == apiName {
		apis, err = c.apiService.GetAPIListByName(ctx, namespaceId, queryName)
		if err != nil {
			return nil, nil, err
		}
	}

	groups := make([]*entry.CommonGroup, 0)

	if queryName == "" {
		//根据父级ID查询目录，如果parentId为0  查跟目录
		groups, err = c.commonGroupStore.GetByParentId(ctx, namespaceId, groupType, tagId, parentId)
		if err != nil {
			return nil, nil, err
		}
	} else {
		groupUUIDS := common.SliceToSliceIds(apis, func(t *model.CommonGroupApi) string {
			return t.GroupUUID
		})

		groupsMaps := common.SliceToMap(groups, func(t *entry.CommonGroup) string {
			return t.Uuid
		})

		commonGroups, err := c.commonGroupStore.GetByName(ctx, namespaceId, queryName)
		if err != nil {
			return nil, nil, err
		}
		for _, group := range commonGroups {
			groupUUIDS = append(groupUUIDS, group.Uuid)
		}

		//查询这些目录的上级目录 直至根目录
		for _, v := range groupUUIDS {
			group, err := c.parentGroup(ctx, v)
			if err != nil {
				return nil, nil, err
			}
			if _, ok := groupsMaps[group.Uuid]; !ok {
				groupsMaps[group.Uuid] = group
			}
		}
		groups = nil
		for _, group := range groupsMaps {
			groups = append(groups, group)
		}
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Sort < groups[j].Sort
		})
	}

	list := make([]*model.CommonGroup, 0)
	for _, group := range groups {
		list = append(list, &model.CommonGroup{
			Group: group,
		})
	}
	for _, group := range list {
		//采用递归查询子目录
		c.subGroup(ctx, namespaceId, groupType, tagId, group)
	}

	resApis := make([]*model.CommonGroupApi, 0, len(apis))
	for _, api := range apis {
		resApis = append(resApis, &model.CommonGroupApi{
			Name:      api.Name,
			UUID:      api.UUID,
			Methods:   api.Methods,
			GroupUUID: api.GroupUUID,
		})
	}

	resRoot := &model.CommonGroupRoot{
		UUID:        uuids,
		Name:        name,
		CommonGroup: list,
	}

	return resRoot, resApis, nil
}

// subGroup 递归查询子目录
func (c *commonGroupService) subGroupV2(ctx context.Context, namespaceId int, groups map[int][]*entry.CommonGroup, uuidMaps map[string]string, parentGroup *model.CommonGroup) {

	newGroups := make([]*entry.CommonGroup, 0)
	if val, ok := groups[parentGroup.Group.Id]; ok {
		newGroups = val
	}

	if len(newGroups) == 0 {
		return
	}

	for _, group := range newGroups {
		if len(uuidMaps) > 0 {
			if _, ok := uuidMaps[group.Uuid]; !ok {
				continue
			}
		}
		val := &model.CommonGroup{Group: group}
		parentGroup.Subgroup = append(parentGroup.Subgroup, val)
		c.subGroupV2(ctx, namespaceId, groups, uuidMaps, val)
	}
}

func (c *commonGroupService) subGroup(ctx context.Context, namespaceId int, groupType string, tagId int, parentGroup *model.CommonGroup) {
	groups, err := c.commonGroupStore.GetByParentId(ctx, namespaceId, groupType, tagId, parentGroup.Group.Id)
	if err != nil {
		return
	}
	if len(groups) == 0 {
		return
	}

	for _, group := range groups {
		val := &model.CommonGroup{Group: group}
		parentGroup.Subgroup = append(parentGroup.Subgroup, val)
		c.subGroup(ctx, namespaceId, groupType, tagId, val)
	}
}

// parentGroup 递归查询跟目录
func (c *commonGroupService) parentGroup(ctx context.Context, parentUUID string) (*entry.CommonGroup, error) {
	group, err := c.commonGroupStore.GetByUUID(ctx, parentUUID)
	if err != nil {
		return nil, err
	}
	if group.ParentId > 0 {
		parentGroup, err := c.commonGroupStore.Get(ctx, group.ParentId)
		if err != nil {
			return nil, err
		}
		return c.parentGroup(ctx, parentGroup.Uuid)
	}
	return group, nil
}

func (c *commonGroupService) parentGroupV2(parentUUID string, groupMaps map[string]*entry.CommonGroup, groupIdMaps map[int]*entry.CommonGroup, outMapUUID map[string]string) {

	if group, ok := groupMaps[parentUUID]; ok {
		if group.ParentId > 0 {
			if parentGroup, ok := groupIdMaps[group.ParentId]; ok {
				outMapUUID[parentGroup.Uuid] = parentGroup.Uuid
				c.parentGroupV2(parentGroup.Uuid, groupMaps, groupIdMaps, outMapUUID)
			}
		}
		outMapUUID[group.Uuid] = group.Uuid
	}
}

// 获取某个目录的名称直至顶级目录名称
func (c *commonGroupService) parentGroupName(uuid string, groupMaps map[string]*entry.CommonGroup, groupIdMaps map[int]*entry.CommonGroup, nameList *[]string) {

	if group, ok := groupMaps[uuid]; ok {
		if group.ParentId > 0 {
			if parentGroup, ok := groupIdMaps[group.ParentId]; ok {
				//去重 todo 可优化
				isAdd := true
				for _, name := range *nameList {
					if name == parentGroup.Name {
						isAdd = false
						continue
					}
				}
				if isAdd {
					*nameList = append(*nameList, parentGroup.Name)
				}
				c.parentGroupName(parentGroup.Uuid, groupMaps, groupIdMaps, nameList)
			}
		}
		//去重 todo 可优化
		isAdd := true
		for _, name := range *nameList {
			if name == group.Name {
				isAdd = false
				continue
			}
		}
		if isAdd {
			*nameList = append(*nameList, group.Name)
		}
	}
}

// subGroup 递归查询子目录
func (c *commonGroupService) subGroupUUIDS(groups map[int][]*entry.CommonGroup, parentGroup *model.CommonGroup, list *[]string) {

	newGroups := make([]*entry.CommonGroup, 0)
	if val, ok := groups[parentGroup.Group.Id]; ok {
		newGroups = val
	}

	if len(newGroups) == 0 {
		return
	}

	for _, group := range newGroups {
		val := &model.CommonGroup{Group: group}
		*list = append(*list, group.Uuid)
		c.subGroupUUIDS(groups, val, list)
	}
}

func (c *commonGroupService) CreateGroup(ctx context.Context, namespaceId int, operator int, groupType, tagName, groupName, uuidStr, parentUuid string) error {
	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return errors.New("params error")
	}
	t := time.Now()
	var err error
	var parentServiceGroup *entry.CommonGroup
	//查询父级的目录ID 没有查到表示是创建的根目录
	if parentUuid != "" {
		parentServiceGroup, err = c.commonGroupStore.GetByUUID(ctx, parentUuid)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}

	parentId := 0
	if parentServiceGroup != nil {
		parentId = parentServiceGroup.Id
	}

	// 判断要创建的目录下有没有重名的， 有则返回报错
	isRepeated, err := c.CheckGroupNameReduplicated(ctx, groupName, parentId)
	if err != nil {
		return err
	}
	if isRepeated {
		return fmt.Errorf("groupName %s is reduplicated. ", groupName)
	}

	if uuidStr == "" {
		uuidStr = uuid.New()
	}

	maxSort, err := c.commonGroupStore.GetMaxSort(ctx, namespaceId, groupType, tagId, parentId)
	if err != nil {
		return err
	}
	group := &entry.CommonGroup{
		Uuid:        uuidStr,
		NamespaceId: namespaceId,
		TagID:       tagId,
		Type:        groupType,
		Name:        groupName,
		ParentId:    parentId,
		Operator:    operator,
		Sort:        maxSort + 1,
		CreateTime:  t,
		UpdateTime:  t,
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Uuid: uuidStr,
		Name: groupName,
	})

	return c.commonGroupStore.Save(ctx, group)

}

func (c *commonGroupService) IsGroupExist(ctx context.Context, uuid string) (bool, error) {
	_, err := c.commonGroupStore.GetByUUID(ctx, uuid)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CheckGroupNameReduplicated 检测分组名是否重复
func (c *commonGroupService) CheckGroupNameReduplicated(ctx context.Context, groupName string, parentID int) (bool, error) {
	list, err := c.commonGroupStore.GetByNameParentID(ctx, groupName, parentID)
	if err != nil {
		return false, err
	}
	if len(list) > 0 {
		return true, nil
	}
	return false, nil
}

func (c *commonGroupService) groupUUIDS(ctx context.Context, namespaceId int, groupType, tagName, parentUuid string) ([]string, error) {
	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return nil, errors.New("params error")
	}

	//查询目录
	groups, err := c.commonGroupStore.GetList(ctx, namespaceId, groupType, tagId)
	if err != nil {
		return nil, err
	}

	resList := &[]string{}
	list := make([]*model.CommonGroup, 0)
	for _, group := range groups {
		if parentUuid != "" {
			if group.Uuid == parentUuid {
				list = append(list, &model.CommonGroup{
					Group: group,
				})
				*resList = append(*resList, group.Uuid)
			}
		} else {
			if group.ParentId == 0 {
				list = append(list, &model.CommonGroup{
					Group: group,
				})
				*resList = append(*resList, group.Uuid)
			}
		}
	}

	groupMaps := common.SliceToMapArray(groups, func(t *entry.CommonGroup) int {
		return t.ParentId
	})

	for _, group := range list {
		//采用递归查询子目录
		c.subGroupUUIDS(groupMaps, group, resList)
	}
	return *resList, nil
}

func (c *commonGroupService) getTagId(ctx context.Context, namespaceId int, groupType, typeName string) int {
	switch groupType {
	case serviceName:
		serviceInfo, err := c.service.GetServiceInfo(ctx, namespaceId, typeName)
		if err != nil {
			return 0
		}
		return serviceInfo.ServiceId
	case apiName:
		return 0
	}
	return -1
}
