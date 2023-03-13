package service

import (
	"context"
	"errors"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/store"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
	"time"
)

type globalVariableService struct {
	globalVariableStore         store.IGlobalVariableStore
	clusterVariableStore        store.IClusterVariableStore
	clusterService              IClusterService
	userInfoService             IUserInfoService
	variableRuntimeStore        store.IVariableRuntimeStore
	variablePublishVersionStore store.IVariablePublishVersionStore
	variableHistoryStore        store.IVariableHistoryStore
	quoteStore                  store.IQuoteStore
}

type IGlobalVariableService interface {
	List(ctx context.Context, pageNum, pageSize, namespace int, key string, status int) ([]*model.GlobalVariableListItem, int, error)
	GetInfo(ctx context.Context, namespaceID int, key string) ([]*model.GlobalVariableDetails, error)
	Create(ctx context.Context, namespaceID, userID int, key, desc string) (int, error)
	Delete(ctx context.Context, namespaceID, userID int, key string) error
	GetByKeys(ctx context.Context, namespaceId int, keys []string) ([]*model.GlobalVariable, error)
	GetById(ctx context.Context, namespaceId int) (*model.GlobalVariable, error)
}

func newGlobalVariableService() IGlobalVariableService {
	s := &globalVariableService{}
	bean.Autowired(&s.globalVariableStore)
	bean.Autowired(&s.clusterVariableStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.variableRuntimeStore)
	bean.Autowired(&s.variablePublishVersionStore)
	bean.Autowired(&s.variableHistoryStore)
	bean.Autowired(&s.userInfoService)
	bean.Autowired(&s.quoteStore)

	return s
}

func (g *globalVariableService) GetById(ctx context.Context, id int) (*model.GlobalVariable, error) {
	variable, err := g.globalVariableStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.GlobalVariable{Variables: variable}, nil
}

func (g *globalVariableService) GetByKeys(ctx context.Context, namespaceId int, keys []string) ([]*model.GlobalVariable, error) {

	variables, err := g.globalVariableStore.GetGlobalVariableByKeys(ctx, namespaceId, keys)
	if err != nil {
		return nil, err
	}
	list := make([]*model.GlobalVariable, 0, len(variables))
	for _, variable := range variables {
		list = append(list, &model.GlobalVariable{Variables: variable})
	}
	return list, nil
}

func (g *globalVariableService) List(ctx context.Context, pageNum, pageSize, namespaceID int, key string, status int) ([]*model.GlobalVariableListItem, int, error) {
	variables, total, err := g.globalVariableStore.GetList(ctx, pageNum, pageSize, namespaceID, key)
	if err != nil {
		return nil, 0, err
	}

	userIds := common.SliceToSliceIds(variables, func(t *entry.Variables) int {
		return t.Operator
	})

	userInfoMaps, _ := g.userInfoService.GetUserInfoMaps(ctx, userIds...)

	variableList := make([]*model.GlobalVariableListItem, 0, len(variables))
	for _, variable := range variables {

		operatorName := ""
		if userInfo, ok := userInfoMaps[variable.Operator]; ok {
			operatorName = userInfo.NickName
		}

		item := &model.GlobalVariableListItem{
			Variables:   variable,
			OperatorStr: operatorName,
		}

		item.Status = 1 //空闲
		count, err := g.quoteStore.Count(ctx, variable.Id, entry.QuoteTargetKindTypeVariable)
		if err != nil {
			return nil, 0, err
		}
		if count > 0 {
			item.Status = 2 //使用中
		}

		if status != 0 && item.Status != status {
			continue
		}
		variableList = append(variableList, item)
	}
	return variableList, total, nil
}

func (g *globalVariableService) GetInfo(ctx context.Context, namespaceID int, key string) ([]*model.GlobalVariableDetails, error) {
	globalVariable, err := g.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return nil, err
	}
	if globalVariable == nil {
		return nil, errors.New("GlobalVariable Key is not exist. ")
	}

	//获取该全局环境变量在所有集群下的集群环境变量
	clusterVariables, err := g.clusterVariableStore.GetVariablesByGlobalVariableID(ctx, namespaceID, globalVariable.Id)
	variableDetails := make([]*model.GlobalVariableDetails, 0, len(clusterVariables))
	for _, cVariable := range clusterVariables {
		if cVariable.Status == 1 {
			continue
		}
		clusterInfo, err := g.clusterService.GetByClusterId(ctx, cVariable.ClusterId)
		if err != nil {
			return nil, err
		}

		//获取集群变量状态
		status, err := g.getClusterVariableStatus(ctx, cVariable.ClusterId, cVariable.VariableId)
		detail := &model.GlobalVariableDetails{
			ClusterVariable: cVariable,
			Status:          status,
			ClusterName:     clusterInfo.Name,
			Environment:     clusterInfo.Env,
		}
		variableDetails = append(variableDetails, detail)
	}

	return variableDetails, nil
}

func (g *globalVariableService) getClusterVariableStatus(ctx context.Context, clusterID, globalVariableID int) (int, error) {
	//获取该集群环境变量
	currentVariable, err := g.clusterVariableStore.GetClusterVariableByClusterIDByGlobalID(ctx, clusterID, globalVariableID)
	if err != nil {
		return 0, err
	}
	//若该集群下没有配置该环境变量，则返回缺失状态
	if currentVariable == nil {
		return 3, nil //缺失
	}

	//获取该集群当前版本的已发布环境变量
	runtime, err := g.variableRuntimeStore.Get(ctx, currentVariable.ClusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	var variableVersionEntry *entry.VariablePublishVersion
	if runtime != nil {
		variableVersionEntry, err = g.variablePublishVersionStore.Get(ctx, runtime.VersionId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return 0, err
		}
	}

	//当前版本已发布的集群环境变量
	versionClusterVariables := make([]*entry.ClusterVariable, 0)
	if variableVersionEntry != nil {
		versionClusterVariables = variableVersionEntry.ClusterVariable
	}
	for _, oldVariable := range versionClusterVariables {
		if oldVariable.VariableId == currentVariable.VariableId && oldVariable.Value == currentVariable.Value {
			return 2, nil //已发布
		}
	}
	return 1, nil //未发布
}

func (g *globalVariableService) Create(ctx context.Context, namespaceID, userID int, key, desc string) (int, error) {

	//验证key值在当前工作空间不存在
	globalVariable, err := g.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return 0, err
	}
	if globalVariable != nil {
		return 0, errors.New("this GlobalVariable key has already existed. ")
	}
	//在variables表中插入全局变量
	variable := &entry.Variables{
		Namespace:  namespaceID,
		Key:        key,
		Desc:       desc,
		Operator:   userID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Name: key,
	})

	if err = g.globalVariableStore.Insert(ctx, variable); err != nil {
		return 0, err
	}
	return variable.Id, err
}

func (g *globalVariableService) Delete(ctx context.Context, namespaceID, userID int, key string) error {
	//根据namespaceID和key获取该全局环境变量variable_id
	globalVariable, err := g.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return err
	}

	if globalVariable == nil {
		return errors.New("GlobalVariable Key is not exist. ")
	}

	count, err := g.quoteStore.Count(ctx, globalVariable.Id, entry.QuoteTargetKindTypeVariable)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("环境变量被引用不可删除")
	}

	//删除该命名空间下所有集群里的该全局环境变量，并且插入删除记录
	variableList, err := g.clusterVariableStore.GetVariablesByGlobalVariableID(ctx, namespaceID, globalVariable.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &model.LogObjectInfo{
		Name: key,
	})

	return g.clusterVariableStore.Transaction(ctx, func(txCtx context.Context) error {
		for _, cVariable := range variableList {
			//删除集群变量
			if _, err = g.clusterVariableStore.Delete(txCtx, cVariable.Id); err != nil {
				return err
			}
			//插入删除记录

			if err = g.variableHistoryStore.HistoryDelete(txCtx, namespaceID, cVariable.Id, &entry.VariableValue{Key: key, Value: cVariable.Value}, userID); err != nil {
				return err
			}
		}

		//删除全局环境变量
		_, err = g.globalVariableStore.Delete(txCtx, globalVariable.Id)
		return err
	})

}
