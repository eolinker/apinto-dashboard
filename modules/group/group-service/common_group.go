package group_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/group"
	"github.com/eolinker/apinto-dashboard/modules/group/group-dto"
	"github.com/eolinker/apinto-dashboard/modules/group/group-entry"
	"github.com/eolinker/apinto-dashboard/modules/group/group-model"
	"github.com/eolinker/apinto-dashboard/modules/group/group-store"
	"github.com/eolinker/apinto-dashboard/modules/upstream"
	"github.com/eolinker/eosc/common/bean"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"time"
)

const (
	ApiName      = "api"
	ServiceName  = "service"
	ModulePlugin = "module_plugin"
)

type commonGroupService struct {
	commonGroupStore group_store.ICommonGroupStore
	service          upstream.IService
	apiService       api.IAPIService
}

func newCommonGroupService() group.ICommonGroupService {
	c := &commonGroupService{}
	bean.Autowired(&c.commonGroupStore)
	bean.Autowired(&c.service)
	bean.Autowired(&c.apiService)
	return c
}

func (c *commonGroupService) GroupSort(ctx context.Context, namespaceId int, groupType, tagName string, input *group_dto.CommGroupSortInput) error {
	//TODO 排序需要检查排序后分组名有无重复的情况

	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return errors.New("params error")
	}
	var (
		g   *group_entry.CommonGroup
		err error
	)

	if input.Root != "" {
		g, err = c.commonGroupStore.GetByUUID(ctx, input.Root)
		if err != nil {
			return err
		}
	} else {
		g = new(group_entry.CommonGroup)
	}

	gs, err := c.commonGroupStore.GetByUUIDS(ctx, input.Items)
	if err != nil {
		return err
	}
	gms := common.SliceToMap(gs, func(t *group_entry.CommonGroup) string {
		return t.Uuid
	})

	list := make([]*group_entry.CommonGroup, 0)
	for i, groupUUID := range input.Items {
		if v, ok := gms[groupUUID]; ok {
			v.Sort = i + 1
			v.ParentId = g.ParentId
			list = append(list, v)
		}
	}

	return c.commonGroupStore.Transaction(ctx, func(txCtx context.Context) error {
		return c.commonGroupStore.UpdateSort(txCtx, list)
	})
}

func (c *commonGroupService) DeleteGroup(ctx context.Context, namespaceId int, operator int, groupType, uuid string) error {
	groupInfo, err := c.commonGroupStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	switch groupType {
	case ApiName:

		groupUUIDS, err := c.GroupUUIDS(ctx, namespaceId, groupType, ApiName, uuid)
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
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuid,
		Name: groupInfo.Name,
	})

	_, err = c.commonGroupStore.Delete(ctx, groupInfo.Id)
	return err
}

func (c *commonGroupService) UpdateGroup(ctx context.Context, namespaceId int, operator int, groupType, name, uuid string) error {
	groupInfo, err := c.commonGroupStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	groupInfo.Operator = operator
	groupInfo.UpdateTime = time.Now()
	groupInfo.Name = name

	//检查修改的分组名是否有重复
	groupList, err := c.commonGroupStore.GetByNameParentID(ctx, name, groupInfo.ParentId)
	if err != nil {
		return err
	}
	for _, item := range groupList {
		//处理修改的分组名不变的情况
		if item.Uuid != groupInfo.Uuid {
			return fmt.Errorf("groupName %s is reduplicated", name)
		}
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuid,
		Name: name,
	})
	return c.commonGroupStore.Save(ctx, groupInfo)
}

func (c *commonGroupService) GroupListAll(ctx context.Context, namespaceId int, groupType, tagName string) ([]*group_entry.CommonGroup, error) {

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

// ToGroupRoot  uuidMaps为其他API关联的groupUUID以及直至顶级目录的groupUUID
func (c *commonGroupService) ToGroupRoot(ctx context.Context, namespaceId int, queryUUid string, groups []*group_entry.CommonGroup, uuidMaps map[string]string) *group_model.CommonGroupRoot {
	parentList := make([]*group_model.CommonGroup, 0)
	name := ""
	for _, groupInfo := range groups {
		if queryUUid != "" {
			if groupInfo.Uuid == queryUUid {
				name = groupInfo.Name
				parentList = append(parentList, &group_model.CommonGroup{
					Group: groupInfo,
				})
			}
		} else {
			if groupInfo.ParentId == 0 {
				if len(uuidMaps) > 0 {
					if _, ok := uuidMaps[groupInfo.Uuid]; !ok {
						continue
					}
				}
				parentList = append(parentList, &group_model.CommonGroup{
					Group: groupInfo,
				})
			}
		}
	}

	array := common.SliceToMapArray(groups, func(t *group_entry.CommonGroup) int {
		return t.ParentId
	})

	for _, groupInfo := range parentList {
		//采用递归查询子目录
		c.subGroupV2(ctx, namespaceId, array, uuidMaps, groupInfo)
	}

	resRoot := &group_model.CommonGroupRoot{
		UUID:        queryUUid,
		Name:        name,
		CommonGroup: parentList,
	}
	return resRoot
}

func (c *commonGroupService) GroupList(ctx context.Context, namespaceId int, groupType, tagName, parentUuid, queryName string) (*group_model.CommonGroupRoot, []*group_model.CommonGroupApi, error) {

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
		groupInfo, err := c.commonGroupStore.GetByUUID(ctx, parentUuid)
		if err != nil {
			return nil, nil, err
		}
		parentId = groupInfo.Id
		uuids = groupInfo.Uuid
		name = groupInfo.Name
	}

	apis := make([]*group_model.CommonGroupApi, 0)
	//查询API的根目录列表
	if groupType == ApiName {
		apis, err = c.apiService.GetAPIListByName(ctx, namespaceId, queryName)
		if err != nil {
			return nil, nil, err
		}
	}

	groups := make([]*group_entry.CommonGroup, 0)

	if queryName == "" {
		//根据父级ID查询目录，如果parentId为0  查跟目录
		groups, err = c.commonGroupStore.GetByParentId(ctx, namespaceId, groupType, tagId, parentId)
		if err != nil {
			return nil, nil, err
		}
	} else {
		groupUUIDS := common.SliceToSliceIds(apis, func(t *group_model.CommonGroupApi) string {
			return t.GroupUUID
		})

		groupsMaps := common.SliceToMap(groups, func(t *group_entry.CommonGroup) string {
			return t.Uuid
		})

		commonGroups, err := c.commonGroupStore.GetByName(ctx, namespaceId, queryName)
		if err != nil {
			return nil, nil, err
		}
		for _, groupInfo := range commonGroups {
			groupUUIDS = append(groupUUIDS, groupInfo.Uuid)
		}

		//查询这些目录的上级目录 直至根目录
		for _, v := range groupUUIDS {
			groupInfo, err := c.parentGroup(ctx, v)
			if err != nil {
				return nil, nil, err
			}
			if _, ok := groupsMaps[groupInfo.Uuid]; !ok {
				groupsMaps[groupInfo.Uuid] = groupInfo
			}
		}
		groups = nil
		for _, groupInfo := range groupsMaps {
			groups = append(groups, groupInfo)
		}
		sort.Slice(groups, func(i, j int) bool {
			return groups[i].Sort < groups[j].Sort
		})
	}

	list := make([]*group_model.CommonGroup, 0)
	for _, groupInfo := range groups {
		list = append(list, &group_model.CommonGroup{
			Group: groupInfo,
		})
	}
	for _, groupInfo := range list {
		//采用递归查询子目录
		c.subGroup(ctx, namespaceId, groupType, tagId, groupInfo)
	}

	resApis := make([]*group_model.CommonGroupApi, 0, len(apis))
	for _, apiInfo := range apis {
		resApis = append(resApis, &group_model.CommonGroupApi{
			Name:      apiInfo.Name,
			UUID:      apiInfo.UUID,
			Methods:   apiInfo.Methods,
			GroupUUID: apiInfo.GroupUUID,
		})
	}

	resRoot := &group_model.CommonGroupRoot{
		UUID:        uuids,
		Name:        name,
		CommonGroup: list,
	}

	return resRoot, resApis, nil
}

// subGroup 递归查询子目录
func (c *commonGroupService) subGroupV2(ctx context.Context, namespaceId int, groups map[int][]*group_entry.CommonGroup, uuidMaps map[string]string, parentGroup *group_model.CommonGroup) {

	newGroups := make([]*group_entry.CommonGroup, 0)
	if val, ok := groups[parentGroup.Group.Id]; ok {
		newGroups = val
	}

	if len(newGroups) == 0 {
		return
	}

	for _, groupInfo := range newGroups {
		if len(uuidMaps) > 0 {
			if _, ok := uuidMaps[groupInfo.Uuid]; !ok {
				continue
			}
		}
		val := &group_model.CommonGroup{Group: groupInfo}
		parentGroup.Subgroup = append(parentGroup.Subgroup, val)
		c.subGroupV2(ctx, namespaceId, groups, uuidMaps, val)
	}
}

func (c *commonGroupService) subGroup(ctx context.Context, namespaceId int, groupType string, tagId int, parentGroup *group_model.CommonGroup) {
	groups, err := c.commonGroupStore.GetByParentId(ctx, namespaceId, groupType, tagId, parentGroup.Group.Id)
	if err != nil {
		return
	}
	if len(groups) == 0 {
		return
	}

	for _, groupInfo := range groups {
		val := &group_model.CommonGroup{Group: groupInfo}
		parentGroup.Subgroup = append(parentGroup.Subgroup, val)
		c.subGroup(ctx, namespaceId, groupType, tagId, val)
	}
}

// parentGroup 递归查询跟目录
func (c *commonGroupService) parentGroup(ctx context.Context, parentUUID string) (*group_entry.CommonGroup, error) {
	groupInfo, err := c.commonGroupStore.GetByUUID(ctx, parentUUID)
	if err != nil {
		return nil, err
	}
	if groupInfo.ParentId > 0 {
		parentGroup, err := c.commonGroupStore.Get(ctx, groupInfo.ParentId)
		if err != nil {
			return nil, err
		}
		return c.parentGroup(ctx, parentGroup.Uuid)
	}
	return groupInfo, nil
}

func (c *commonGroupService) ParentGroupV2(parentUUID string, groupMaps map[string]*group_entry.CommonGroup, groupIdMaps map[int]*group_entry.CommonGroup, outMapUUID map[string]string) {

	if groupInfo, ok := groupMaps[parentUUID]; ok {
		if groupInfo.ParentId > 0 {
			if parentGroup, ok := groupIdMaps[groupInfo.ParentId]; ok {
				outMapUUID[parentGroup.Uuid] = parentGroup.Uuid
				c.ParentGroupV2(parentGroup.Uuid, groupMaps, groupIdMaps, outMapUUID)
			}
		}
		outMapUUID[groupInfo.Uuid] = groupInfo.Uuid
	}
}

// ParentGroupName 获取某个目录的名称直至顶级目录名称
func (c *commonGroupService) ParentGroupName(uuid string, groupMaps map[string]*group_entry.CommonGroup, groupIdMaps map[int]*group_entry.CommonGroup, nameList *[]string) {

	if groupInfo, ok := groupMaps[uuid]; ok {
		if groupInfo.ParentId > 0 {
			if parentGroup, ok := groupIdMaps[groupInfo.ParentId]; ok {
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
				c.ParentGroupName(parentGroup.Uuid, groupMaps, groupIdMaps, nameList)
			}
		}
		//去重 todo 可优化
		isAdd := true
		for _, name := range *nameList {
			if name == groupInfo.Name {
				isAdd = false
				continue
			}
		}
		if isAdd {
			*nameList = append(*nameList, groupInfo.Name)
		}
	}
}

// SubGroupUUIDS subGroup 递归查询子目录
func (c *commonGroupService) SubGroupUUIDS(groups map[int][]*group_entry.CommonGroup, parentGroup *group_model.CommonGroup, list *[]string) {

	newGroups := make([]*group_entry.CommonGroup, 0)
	if val, ok := groups[parentGroup.Group.Id]; ok {
		newGroups = val
	}

	if len(newGroups) == 0 {
		return
	}

	for _, groupInfo := range newGroups {
		val := &group_model.CommonGroup{Group: groupInfo}
		*list = append(*list, groupInfo.Uuid)
		c.SubGroupUUIDS(groups, val, list)
	}
}

func (c *commonGroupService) CreateGroup(ctx context.Context, namespaceId int, operator int, groupType, tagName, groupName, uuidStr, parentUuid string) error {
	tagId := c.getTagId(ctx, namespaceId, groupType, tagName)
	if tagId == -1 {
		return errors.New("params error")
	}
	t := time.Now()
	var err error
	var parentServiceGroup *group_entry.CommonGroup
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
	groupInfo := &group_entry.CommonGroup{
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
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuidStr,
		Name: groupName,
	})

	return c.commonGroupStore.Save(ctx, groupInfo)

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

func (c *commonGroupService) GroupUUIDS(ctx context.Context, namespaceId int, groupType, tagName, parentUuid string) ([]string, error) {
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
	list := make([]*group_model.CommonGroup, 0)
	for _, groupInfo := range groups {
		if parentUuid != "" {
			if groupInfo.Uuid == parentUuid {
				list = append(list, &group_model.CommonGroup{
					Group: groupInfo,
				})
				*resList = append(*resList, groupInfo.Uuid)
			}
		} else {
			if groupInfo.ParentId == 0 {
				list = append(list, &group_model.CommonGroup{
					Group: groupInfo,
				})
				*resList = append(*resList, groupInfo.Uuid)
			}
		}
	}

	groupMaps := common.SliceToMapArray(groups, func(t *group_entry.CommonGroup) int {
		return t.ParentId
	})

	for _, groupInfo := range list {
		//采用递归查询子目录
		c.SubGroupUUIDS(groupMaps, groupInfo, resList)
	}
	return *resList, nil
}

func (c *commonGroupService) getTagId(ctx context.Context, namespaceId int, groupType, typeName string) int {
	switch groupType {
	case ServiceName:
		serviceInfo, err := c.service.GetServiceInfo(ctx, namespaceId, typeName)
		if err != nil {
			return 0
		}
		return serviceInfo.ServiceId
	case ApiName:
		return 0
	case ModulePlugin:
		return 0
	}
	return -1
}

func (c *commonGroupService) GetGroupInfo(ctx context.Context, uuid string) (*group_entry.CommonGroup, error) {
	return c.commonGroupStore.GetByUUID(ctx, uuid)
}
