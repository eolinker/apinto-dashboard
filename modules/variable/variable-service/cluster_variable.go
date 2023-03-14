package variable_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/history-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	cluster_dto2 "github.com/eolinker/apinto-dashboard/modules/cluster/cluster-dto"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-entry"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	variable_entry2 "github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-model"
	variable_store2 "github.com/eolinker/apinto-dashboard/modules/variable/variable-store"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

type clusterVariableService struct {
	clusterService              cluster.IClusterService
	clusterNodeService          cluster.IClusterNodeService
	namespaceService            namespace.INamespaceService
	apintoClient                cluster.IApintoClient
	globalVariableStore         variable_store2.IGlobalVariableStore
	variableHistoryStore        variable_store2.IVariableHistoryStore
	variablePublishVersionStore variable_store2.IVariablePublishVersionStore
	variableRuntimeStore        variable_store2.IVariableRuntimeStore
	variablePublishHistoryStore variable_store2.IVariablePublishHistoryStore
	clusterVariableStore        variable_store2.IClusterVariableStore
	lockService                 locker_service.IAsynLockService
	userInfoService             user.IUserInfoService
}

func newClusterVariableService() variable.IClusterVariableService {

	s := &clusterVariableService{}
	bean.Autowired(&s.clusterVariableStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.globalVariableStore)
	bean.Autowired(&s.variableHistoryStore)
	bean.Autowired(&s.variablePublishHistoryStore)
	bean.Autowired(&s.variablePublishVersionStore)
	bean.Autowired(&s.variableRuntimeStore)
	bean.Autowired(&s.namespaceService)
	bean.Autowired(&s.clusterNodeService)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.lockService)
	bean.Autowired(&s.userInfoService)

	return s
}

func (c *clusterVariableService) GetList(ctx context.Context, namespaceID int, clusterName string) ([]*variable_model.ClusterVariableListItem, error) {
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return nil, err
	}

	list := make([]*variable_model.ClusterVariableListItem, 0)

	//获取工作空间下所有全局环境变量
	globalVariables, _, err := c.globalVariableStore.GetList(ctx, 0, 0, namespaceID, "")
	if err != nil {
		return nil, err
	}

	//获取该集群下所有环境变量
	clusterVariables, err := c.clusterVariableStore.GetByClusterIds(ctx, cluster.Id)
	if err != nil {
		return nil, err
	}

	clusterVariablesMap := common.SliceToMap(clusterVariables, func(t *variable_entry2.ClusterVariable) int {
		return t.VariableId
	})

	//获取该集群当前版本的已发布环境变量
	variablePublishVersionEntry, err := c.GetPublishVersion(ctx, cluster.Id)
	if err != nil {
		return nil, err
	}

	//当前版本已发布的集群环境变量
	versionClusterVariables := make([]*variable_entry2.ClusterVariable, 0)
	if variablePublishVersionEntry != nil {
		versionClusterVariables = variablePublishVersionEntry.ClusterVariable
	}
	versionClusterVariablesMap := common.SliceToMap(versionClusterVariables, func(t *variable_entry2.ClusterVariable) string {
		return t.Key
	})

	userIds := common.SliceToSliceIds(globalVariables, func(t *variable_entry2.Variables) int {
		return t.Operator
	})

	userIds = append(userIds, common.SliceToSliceIds(clusterVariables, func(t *variable_entry2.ClusterVariable) int {
		return t.Operator
	})...)

	userInfoMaps, _ := c.userInfoService.GetUserInfoMaps(ctx, userIds...)

	//对比全局变量
	for _, gVariable := range globalVariables {

		item := &variable_model.ClusterVariableListItem{
			ClusterVariable: nil,
			Desc:            gVariable.Desc,
			Publish:         1, //未发布
		}

		if now, has := clusterVariablesMap[gVariable.Id]; has {
			if now.Status == 0 {
				item.ClusterVariable = now
				if old, has := versionClusterVariablesMap[gVariable.Key]; has {
					if now.Value == old.Value {
						item.Publish = 2 //已发布
					} else {
						item.Publish = 1 //未发布
					}
				} else {
					item.Publish = 1 //未发布
				}
				operatorName := ""
				if userInfo, ok := userInfoMaps[now.Operator]; ok {
					operatorName = userInfo.NickName
				}
				item.Operator = operatorName
			} else {
				item.ClusterVariable = &variable_entry2.ClusterVariable{Key: gVariable.Key}
				item.Publish = 3
			}

		} else {
			item.ClusterVariable = &variable_entry2.ClusterVariable{Key: gVariable.Key}
			item.Publish = 3 //缺失
		}

		//筛选出集群版本有，而全局中没有的变量
		delete(versionClusterVariablesMap, gVariable.Key)
		list = append(list, item)
	}

	//对比集群发布版本中有，而全局中没有的变量
	//for _, vClusterVariable := range versionClusterVariablesMap {
	//
	//	if vi, has := clusterVariablesMap[vClusterVariable.VariableId]; has {
	//
	//		operator := ""
	//		if userInfo, ok := userInfoMaps[vi.Operator]; ok {
	//			operator = userInfo.Operator
	//		}
	//
	//		item := &model.ClusterVariableListItem{
	//			ClusterVariable: &entry.ClusterVariable{Key: vClusterVariable.Key, UpdateTime: vi.UpdateTime},
	//			Desc:            "", //TODO
	//			Operator:        operator,
	//			Publish:         1, //未发布
	//			Status:          int(entry.OptDelete),
	//		}
	//		list = append(list, item)
	//	}
	//
	//}

	return list, nil
}

func (c *clusterVariableService) Create(ctx context.Context, namespaceID int, clusterName string, userID int, key, value, desc string) error {
	//验证clusterName存在，并返回clusterID
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return err
	}

	if err = c.lockService.Lock(locker_service.LockNameVariable, cluster.Id); err != nil {
		return err
	}
	defer c.lockService.Unlock(locker_service.LockNameVariable, cluster.Id)

	//验证新增的集群环境变量和当前工作空间的环境变量没有冲突
	globalVariable, err := c.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return err
	}
	if globalVariable != nil {
		return errors.New("this GlobalVariable key has already existed. ")
	}

	t := time.Now()

	//新增全局环境变量，并返回variable_id
	variable := &variable_entry2.Variables{
		Namespace:  namespaceID,
		Key:        key,
		Desc:       desc,
		Operator:   userID,
		CreateTime: t,
		UpdateTime: t,
	}

	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        key,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return c.globalVariableStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = c.globalVariableStore.Insert(txCtx, variable); err != nil {
			return err
		}
		//在variable_cluster表插入集群环境变量
		clusterVariable := &variable_entry2.ClusterVariable{
			ClusterId:   cluster.Id,
			VariableId:  variable.Id,
			NamespaceId: namespaceID,
			Key:         key,
			Value:       value,
			Operator:    userID,
			CreateTime:  t,
			UpdateTime:  t,
		}

		if err = c.clusterVariableStore.Insert(txCtx, clusterVariable); err != nil {
			return err
		}

		return c.variableHistoryStore.HistoryAdd(txCtx, namespaceID, clusterVariable.Id, &variable_entry2.VariableValue{Key: key, Value: value}, userID)
	})
}

func (c *clusterVariableService) Update(ctx context.Context, namespaceID int, clusterName string, userID int, key, value string) error {
	//验证clusterName合法性
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return common.ClusterNotExist
	}

	//验证key合法性,并返回variable_id
	globalVariable, err := c.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return err
	}

	if globalVariable == nil {
		return errors.New("globalVariable key is not exist. ")
	}

	if err = c.lockService.Lock(locker_service.LockNameVariable, cluster.Id); err != nil {
		return err
	}
	defer c.lockService.Unlock(locker_service.LockNameVariable, cluster.Id)

	clusterVariable, err := c.clusterVariableStore.GetClusterVariableByClusterIDByGlobalID(ctx, cluster.Id, globalVariable.Id)
	if err != nil {
		return err
	}
	//若变量存在且未软删，而且新旧值一样
	if clusterVariable != nil && clusterVariable.Status == 0 && clusterVariable.Value == value {
		return errors.New("新值和旧值一样，无需保存")
	}

	oldValue := &variable_entry2.VariableValue{Key: key}
	if clusterVariable != nil {
		oldValue.Value = clusterVariable.Value
	}

	t := time.Now()
	// 若该环境变量为空，则新建
	if clusterVariable == nil {
		clusterVariable = &variable_entry2.ClusterVariable{
			ClusterId:   cluster.Id,
			VariableId:  globalVariable.Id,
			NamespaceId: namespaceID,
			Status:      0,
			Key:         key,
			Value:       value,
			Operator:    userID,
			CreateTime:  t,
			UpdateTime:  t,
		}
	} else {
		clusterVariable.Status = 0
		clusterVariable.Value = value
		clusterVariable.Operator = userID
		clusterVariable.UpdateTime = t
	}

	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        key,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return c.clusterVariableStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = c.clusterVariableStore.Save(txCtx, clusterVariable); err != nil {
			return err
		}

		return c.variableHistoryStore.HistoryEdit(txCtx, namespaceID, clusterVariable.Id, oldValue, &variable_entry2.VariableValue{Key: key, Value: value}, userID)
	})

}

func (c *clusterVariableService) Delete(ctx context.Context, namespaceID int, clusterName string, userID int, key string) error {
	//验证clusterName合法性
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return err
	}

	//验证key合法性,并返回variable_id
	globalVariable, err := c.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return err
	}
	if globalVariable == nil {
		return errors.New("globalVariable Key is not exist. ")
	}

	if err = c.lockService.Lock(locker_service.LockNameVariable, cluster.Id); err != nil {
		return err
	}
	defer c.lockService.Unlock(locker_service.LockNameVariable, cluster.Id)

	//检查variable_cluster表中是否有该集群环境变量，没有则直接返回
	clusterVariable, err := c.clusterVariableStore.GetClusterVariableByClusterIDByGlobalID(ctx, cluster.Id, globalVariable.Id)
	if err != nil {
		return err
	}

	if clusterVariable == nil || clusterVariable.Status != 0 {
		return nil
	}

	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        key,
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return c.clusterVariableStore.Transaction(ctx, func(txCtx context.Context) error {
		clusterVariable.Status = 1
		if _, err = c.clusterVariableStore.Update(txCtx, clusterVariable); err != nil {
			return err
		}

		return c.variableHistoryStore.HistoryDelete(txCtx, namespaceID, clusterVariable.Id, &variable_entry2.VariableValue{Key: key, Value: clusterVariable.Value}, userID)
	})

}

// DeleteAll 调用方需要开启事务 删除集群下的所有环境变量
func (c *clusterVariableService) DeleteAll(ctx context.Context, namespaceID int, clusterId, userID int) error {

	clusterVariables, err := c.clusterVariableStore.GetByClusterIds(ctx, clusterId)
	if err != nil {
		return err
	}

	for _, clusterVariable := range clusterVariables {
		clusterVariable.Status = 1
		if _, err = c.clusterVariableStore.Update(ctx, clusterVariable); err != nil {
			return err
		}

		if err = c.variableHistoryStore.HistoryDelete(ctx, namespaceID, clusterVariable.Id, &variable_entry2.VariableValue{Key: clusterVariable.Key, Value: clusterVariable.Value}, userID); err != nil {
			return err
		}
	}
	return nil

}

// SyncConf 同步配置
func (c *clusterVariableService) SyncConf(ctx context.Context, namespaceId, userId int, clusterName string, conf *cluster_dto2.SyncConf) error {
	_, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return common.ClusterNotExist
	}

	//查询被同步的集群ID
	clusterIds := common.SliceToSliceIds(conf.Clusters, func(t *cluster_dto2.ClusterInput) int {
		return t.Id
	})

	//查询被同步集群的环境变量原来的信息
	variables, err := c.clusterVariableStore.GetByClusterIds(ctx, clusterIds...)
	if err != nil {
		return err
	}

	variableMap := common.SliceToMap(variables, func(t *variable_entry2.ClusterVariable) string {
		return fmt.Sprintf("%d_%s", t.VariableId, t.Key)
	})

	clusterVariableList := make([]*variable_entry2.ClusterVariable, 0)
	variableHistoryList := make([]*variable_entry2.VariableHistory, 0)

	t := time.Now()
	for _, cluster := range conf.Clusters {
		for _, variable := range conf.Variables {
			clusterVariableList = append(clusterVariableList, &variable_entry2.ClusterVariable{
				ClusterId:   cluster.Id,
				VariableId:  variable.VariableId,
				NamespaceId: namespaceId,
				Key:         variable.Key,
				Value:       variable.Value,
				Operator:    userId,
				CreateTime:  t,
				UpdateTime:  t,
			})
			optType := history_entry.OptAdd //默认新增
			oldValue := ""
			if v, ok := variableMap[fmt.Sprintf("%d_%s", variable.VariableId, variable.Key)]; ok {
				//修改
				optType = history_entry.OptEdit
				oldValue = v.Value
			}
			variableHistoryList = append(variableHistoryList, &variable_entry2.VariableHistory{

				VariableId: variable.VariableId,
				OldValue:   variable_entry2.VariableValue{Key: variable.Key, Value: oldValue},
				NewValue:   variable_entry2.VariableValue{Key: variable.Key, Value: variable.Value},
				OptType:    optType,
			})
		}
	}

	return c.clusterVariableStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = c.clusterVariableStore.UpdateVariables(txCtx, clusterVariableList); err != nil {
			return err
		}
		// 查询变更记录表
		for _, h := range variableHistoryList {

			if err = c.variableHistoryStore.History(txCtx, namespaceId, h.VariableId, &h.OldValue, &h.NewValue, h.OptType, userId); err != nil {
				return err
			}
		}
		return nil
	})
}

// QueryHistory 环境变量的变更记录查询
func (c *clusterVariableService) QueryHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*variable_model.ClusterVariableHistory, int, error) {
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, 0, common.ClusterNotExist
	}

	variables, err := c.clusterVariableStore.GetByClusterIds(ctx, cluster.Id, cluster.Id)
	if err != nil {
		return nil, 0, err
	}

	ids := common.SliceToSliceIds(variables, func(t *variable_entry2.ClusterVariable) int {
		return t.Id
	})

	histories, count, err := c.variableHistoryStore.Page(ctx, namespaceId, pageNum, pageSize, ids...)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*variable_model.ClusterVariableHistory, 0, len(histories))
	for _, history := range histories {
		list = append(list, &variable_model.ClusterVariableHistory{VariableHistory: history})
	}

	return list, count, nil
}

// ToPublishs 待发布的环境变量
func (c *clusterVariableService) ToPublishs(ctx context.Context, namespaceId int, clusterName string) ([]*variable_model.VariableToPublish, error) {
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, common.ClusterNotExist
	}

	//查询现在数据库中的环境变量
	currentClusterVariables, err := c.clusterVariableStore.GetByClusterIds(ctx, cluster.Id)
	if err != nil {
		return nil, err
	}

	newClusterVariables := make([]*variable_entry2.ClusterVariable, 0)
	for _, variable := range currentClusterVariables {
		if variable.Status == 0 {
			newClusterVariables = append(newClusterVariables, variable)
		}
	}

	//查询当前版本下的环境变量
	clusterRuntime, err := c.variableRuntimeStore.GetForCluster(ctx, cluster.Id, cluster.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//表示当前集群还没发布任何环境遍历版本
	respList := make([]*variable_model.VariableToPublish, 0)
	if clusterRuntime == nil {
		for _, value := range newClusterVariables {
			entryVariable := variable_entry2.VariableToPublish{
				Key:             value.Key,
				VariableId:      value.VariableId,
				NoReleasedValue: value.Value,
				CreateTime:      value.CreateTime,
				OptType:         1,
			}
			respList = append(respList, &variable_model.VariableToPublish{VariableToPublish: entryVariable})
		}
		return respList, nil
	}

	version, err := c.variablePublishVersionStore.Get(ctx, clusterRuntime.VersionId)
	if err != nil {
		return nil, err
	}

	clusterVariable := version.ClusterVariable

	//差异化对比
	versionClusterVariableMap := common.SliceToMap(clusterVariable, func(t *variable_entry2.ClusterVariable) string {
		return t.Key
	})

	addList, updateList, delList := common.DiffContrast(clusterVariable, newClusterVariables)

	for _, val := range addList {
		entryVariable := variable_entry2.VariableToPublish{
			VariableId:      val.VariableId,
			Key:             val.Key,
			NoReleasedValue: val.Value,
			CreateTime:      val.CreateTime,
			OptType:         1,
		}
		respList = append(respList, &variable_model.VariableToPublish{VariableToPublish: entryVariable})
	}
	for _, val := range updateList {
		finishValue := ""
		if variable, ok := versionClusterVariableMap[val.Key]; ok {
			finishValue = variable.Value
		}

		entryVariable := variable_entry2.VariableToPublish{
			VariableId:      val.VariableId,
			Key:             val.Key,
			NoReleasedValue: val.Value,
			FinishValue:     finishValue,
			CreateTime:      val.CreateTime,
			OptType:         2,
		}
		respList = append(respList, &variable_model.VariableToPublish{VariableToPublish: entryVariable})
	}

	for _, val := range delList {
		finishValue := ""
		if cVariable, ok := versionClusterVariableMap[val.Key]; ok {
			finishValue = cVariable.Value
		}
		entryVariable := variable_entry2.VariableToPublish{
			VariableId:  val.VariableId,
			Key:         val.Key,
			FinishValue: finishValue,
			CreateTime:  val.CreateTime,
			OptType:     3,
		}

		respList = append(respList, &variable_model.VariableToPublish{VariableToPublish: entryVariable})
	}

	sort.Slice(respList, func(i, j int) bool {
		return respList[i].CreateTime.Unix() < respList[j].CreateTime.Unix()
	})

	return respList, nil
}

// Publish 发布环境变量
func (c *clusterVariableService) Publish(ctx context.Context, namespaceId, userId int, clusterName, versionName, desc, source string) error {
	t := time.Now()

	namespace, err := c.namespaceService.GetById(namespaceId)
	if err != nil {
		return err
	}

	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	if err = c.lockService.Lock(locker_service.LockNameVariable, cluster.Id); err != nil {
		return err
	}
	defer c.lockService.Unlock(locker_service.LockNameVariable, cluster.Id)

	//查询版本名称是否重复
	publishHistory, err := c.variablePublishHistoryStore.GetByVersionName(ctx, versionName, cluster.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if publishHistory != nil && publishHistory.Id > 0 {
		return errors.New("该版本名称已存在")
	}

	bytes, err := common.Base64Decode(source)
	if err != nil {
		return err
	}

	publishes := make([]*variable_entry2.VariableToPublish, 0)

	if err = json.Unmarshal(bytes, &publishes); err != nil {
		return err
	}

	if len(publishes) == 0 {
		return errors.New("没有变化无需发布")
	}

	//获取集群当前运行的版本
	currentVersion, err := c.GetPublishVersion(ctx, cluster.Id)
	if err != nil {
		return err
	}

	newClusterVariables := make([]*variable_entry2.ClusterVariable, 0)
	insertClusterVariables := make([]*variable_entry2.ClusterVariable, 0)

	for _, publish := range publishes {
		newClusterVariables = append(newClusterVariables, &variable_entry2.ClusterVariable{
			ClusterId:   cluster.Id,
			VariableId:  publish.VariableId,
			NamespaceId: namespaceId,
			Key:         publish.Key,
			Value:       publish.NoReleasedValue,
			CreateTime:  publish.CreateTime,
		})
	}

	if currentVersion != nil && currentVersion.Id > 0 {
		currentVersionClusterVariables := currentVersion.ClusterVariable

		currentVersionClusterVariablesMaps := common.SliceToMap(currentVersionClusterVariables, func(t *variable_entry2.ClusterVariable) string {
			return t.Key
		})

		for _, publish := range publishes {
			if publish.OptType == 1 { //新增 直接追加
				currentVersionClusterVariablesMaps[publish.Key] = &variable_entry2.ClusterVariable{
					ClusterId:   cluster.Id,
					VariableId:  publish.VariableId,
					NamespaceId: namespaceId,
					Key:         publish.Key,
					Value:       publish.NoReleasedValue,
				}
			} else if publish.OptType == 2 { //修改 找到旧版本的key  然后改value值
				if val, ok := currentVersionClusterVariablesMaps[publish.Key]; ok {
					val.Value = publish.NoReleasedValue
				}
			} else if publish.OptType == 3 {
				delete(currentVersionClusterVariablesMaps, publish.Key)
			}
		}
		for _, variable := range currentVersionClusterVariablesMaps {
			insertClusterVariables = append(insertClusterVariables, variable)
		}

	} else {
		//当前没有旧版本 表示这是第一次发布  直接用发布的变量存储
		insertClusterVariables = newClusterVariables
	}

	variablePublishVersionConfig := variable_entry2.VariablePublishVersionConfig{
		ClusterVariable: insertClusterVariables,
	}

	newVersion := &variable_entry2.VariablePublishVersion{
		ClusterId:                    cluster.Id,
		NamespaceId:                  namespaceId,
		Desc:                         desc,
		VariablePublishVersionConfig: variablePublishVersionConfig,
		Operator:                     userId,
		CreateTime:                   t,
	}

	names := make([]string, 0)
	for _, variable := range variablePublishVersionConfig.ClusterVariable {

		names = append(names, variable.Key)
	}

	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        strings.Join(names, ","),
		ClusterId:   cluster.Id,
		ClusterName: clusterName,
	})

	return c.variablePublishVersionStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.variablePublishVersionStore.Save(txCtx, newVersion); err != nil {
			return err
		}
		//当前集群运行的版本
		variableRuntime := &variable_entry2.VariableRuntime{
			VersionId:   newVersion.Id,
			ClusterId:   cluster.Id,
			NamespaceId: namespaceId,
			Operator:    userId,
			IsOnline:    true,
			CreateTime:  t,
			UpdateTime:  t,
		}

		history := &variable_entry2.VariablePublishHistory{
			VersionName: versionName,
			ClusterId:   cluster.Id,
			NamespaceId: namespaceId,
			Desc:        desc,
			VersionId:   newVersion.Id,
			VariablePublishHistoryInfo: variable_entry2.VariablePublishHistoryInfo{
				VariableToPublish: publishes,
			},
			OptType:  1,
			Operator: userId,
			OptTime:  t,
		}
		if err = c.variablePublishHistoryStore.Insert(txCtx, history); err != nil {
			return err
		}

		if err = c.variableRuntimeStore.Save(txCtx, variableRuntime); err != nil {
			return err
		}
		client, err := c.apintoClient.GetClient(ctx, cluster.Id)
		if err != nil {
			return err
		}

		variableMaps := make(map[string]string)
		for _, variable := range insertClusterVariables {
			variableMaps[variable.Key] = variable.Value
		}

		return client.ForVariable().Publish(namespace.Name, variableMaps)
	})
}

func (c *clusterVariableService) ResetOnline(ctx context.Context, namespaceId, clusterId int) {
	runtime, err := c.variableRuntimeStore.GetForCluster(ctx, clusterId, clusterId)
	if err != nil {
		log.Errorf("clusterVariableService-ResetOnline-GetRuntime clusterId=%d,err=%s", clusterId, err.Error())
		return
	}

	if runtime.IsOnline {
		version, err := c.variablePublishVersionStore.Get(ctx, runtime.VersionId)
		if err != nil {
			return
		}

		client, err := c.apintoClient.GetClient(ctx, clusterId)
		if err != nil {
			return
		}

		variableMaps := make(map[string]string)
		for _, variable := range version.ClusterVariable {
			variableMaps[variable.Key] = variable.Value
		}

		namespace, err := c.namespaceService.GetById(namespaceId)
		if err != nil {
			return
		}

		client.ForVariable().Publish(namespace.Name, variableMaps)

	}
}

func (c *clusterVariableService) PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*variable_model.VariablePublish, int, error) {
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, 0, common.ClusterNotExist
	}

	list, count, err := c.variablePublishHistoryStore.GetByClusterPage(ctx, pageNum, pageSize, cluster.Id)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]*variable_model.VariablePublish, 0, len(list))

	userIds := common.SliceToSliceIds(list, func(t *variable_entry2.VariablePublishHistory) int {
		return t.Operator
	})

	infoMaps, _ := c.userInfoService.GetUserInfoMaps(ctx, userIds...)

	for _, publish := range list {

		data := publish.VariableToPublish

		details := make([]*variable_model.VariablePublishDetails, 0, len(data))

		for _, val := range data {
			details = append(details, &variable_model.VariablePublishDetails{
				Key:        val.Key,
				OldValue:   val.FinishValue,
				NewValue:   val.NoReleasedValue,
				OptType:    val.OptType,
				CreateTime: val.CreateTime,
			})
		}

		operator := ""
		if userInfo, ok := infoMaps[publish.Operator]; ok {
			operator = userInfo.NickName
		}
		resp = append(resp, &variable_model.VariablePublish{
			Id:         publish.Id,
			Name:       publish.VersionName,
			OptType:    publish.OptType,
			Operator:   operator,
			CreateTime: publish.OptTime,
			Details:    details,
		})
	}

	return resp, count, nil
}

func (c *clusterVariableService) GetPublishVersion(ctx context.Context, clusterId int) (*variable_model.VariablePublishVersion, error) {
	//获取集群当前运行的版本
	currentRuntime, err := c.variableRuntimeStore.GetForCluster(ctx, clusterId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var currentVersion *variable_entry2.VariablePublishVersion
	if currentRuntime != nil {
		//获取当前集群运行版本的详细信息
		currentVersion, err = c.variablePublishVersionStore.Get(ctx, currentRuntime.VersionId)
		if err != nil {
			return nil, err
		}
	}

	return (*variable_model.VariablePublishVersion)(currentVersion), nil
}

func (c *clusterVariableService) GetSyncConf(ctx context.Context, namespaceId int, clusterName string) (*variable_model.ClustersVariables, error) {
	cluster, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, common.ClusterNotExist
	}

	clusters, err := c.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	variables, err := c.clusterVariableStore.GetByClusterIds(ctx, cluster.Id)
	if err != nil {
		return nil, err
	}

	newCluster := make([]*cluster_entry.Cluster, 0)
	for _, val := range clusters {
		//过滤自己的集群信息
		if val.Id == cluster.Id {
			continue
		}
		newCluster = append(newCluster, val.Cluster)
	}
	resp := &variable_model.ClustersVariables{
		Clusters:  newCluster,
		Variables: variables,
	}

	return resp, nil
}
