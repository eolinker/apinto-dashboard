package variable_service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/apinto-dashboard/modules/variable"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-entry"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-model"
	"github.com/eolinker/apinto-dashboard/modules/variable/variable-store"
	"github.com/eolinker/eosc/common/bean"
	"gorm.io/gorm"
)

type globalVariableService struct {
	clusterVariableService      variable.IClusterVariableService
	globalVariableStore         variable_store.IGlobalVariableStore
	clusterVariableStore        variable_store.IClusterVariableStore
	clusterService              cluster.IClusterService
	userInfoService             user.IUserInfoService
	variableRuntimeStore        variable_store.IVariableRuntimeStore
	variablePublishVersionStore variable_store.IVariablePublishVersionStore
	variableHistoryStore        variable_store.IVariableHistoryStore
	quoteStore                  quote_store.IQuoteStore
}

func newGlobalVariableService() variable.IGlobalVariableService {
	s := &globalVariableService{}
	bean.Autowired(&s.clusterVariableService)
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

func (g *globalVariableService) GetById(ctx context.Context, id int) (*variable_model.GlobalVariable, error) {
	variableInfo, err := g.globalVariableStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &variable_model.GlobalVariable{Variables: variableInfo}, nil
}

func (g *globalVariableService) GetByKeys(ctx context.Context, namespaceId int, keys []string) ([]*variable_model.GlobalVariable, error) {

	variables, err := g.globalVariableStore.GetGlobalVariableByKeys(ctx, namespaceId, keys)
	if err != nil {
		return nil, err
	}
	list := make([]*variable_model.GlobalVariable, 0, len(variables))
	for _, variableInfo := range variables {
		list = append(list, &variable_model.GlobalVariable{Variables: variableInfo})
	}
	return list, nil
}

func (g *globalVariableService) List(ctx context.Context, pageNum, pageSize, namespaceID int, key string, status int) ([]*variable_model.GlobalVariableListItem, int, error) {
	variables, total, err := g.globalVariableStore.GetList(ctx, pageNum, pageSize, namespaceID, key)
	if err != nil {
		return nil, 0, err
	}

	variableList := make([]*variable_model.GlobalVariableListItem, 0, len(variables))
	for _, variableInfo := range variables {

		item := &variable_model.GlobalVariableListItem{
			Variables: variableInfo,
		}

		item.Status = 1 //空闲
		count, err := g.quoteStore.Count(ctx, variableInfo.Id, quote_entry.QuoteTargetKindTypeVariable)
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
	user.SetUserName(g.userInfoService, ctx, variableList...)
	return variableList, total, nil
}

func (g *globalVariableService) GetInfo(ctx context.Context, namespaceID int, key string) ([]*variable_model.GlobalVariableDetails, error) {
	globalVariable, err := g.globalVariableStore.GetGlobalVariableIDByKey(ctx, namespaceID, key)
	if err != nil {
		return nil, err
	}
	if globalVariable == nil {
		return nil, errors.New("GlobalVariable Key is not exist. ")
	}

	//获取该全局环境变量在所有集群下的集群环境变量
	clusterVariables, err := g.clusterVariableStore.GetVariablesByGlobalVariableID(ctx, namespaceID, globalVariable.Id)
	variableDetails := make([]*variable_model.GlobalVariableDetails, 0, len(clusterVariables))
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
		detail := &variable_model.GlobalVariableDetails{
			ClusterVariable: cVariable,
			Status:          status,
			ClusterName:     clusterInfo.Title, //这里用title是因为显示需要集群名称，之前是名称是用name的
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

	var variableVersionEntry *variable_entry.VariablePublishVersion
	if runtime != nil {
		variableVersionEntry, err = g.variablePublishVersionStore.Get(ctx, runtime.VersionId)
		if err != nil && err != gorm.ErrRecordNotFound {
			return 0, err
		}
	}

	//当前版本已发布的集群环境变量
	versionClusterVariables := make([]*variable_entry.ClusterVariable, 0)
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
	variableInfo := &variable_entry.Variables{
		Namespace:  namespaceID,
		Key:        key,
		Desc:       desc,
		Operator:   userID,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	//编写日志操作对象信息
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: key,
	})

	if err = g.globalVariableStore.Insert(ctx, variableInfo); err != nil {
		return 0, err
	}
	return variableInfo.Id, err
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

	count, err := g.quoteStore.Count(ctx, globalVariable.Id, quote_entry.QuoteTargetKindTypeVariable)
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
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: key,
	})

	return g.clusterVariableStore.Transaction(ctx, func(txCtx context.Context) error {
		for _, cVariable := range variableList {
			//删除集群变量
			if _, err = g.clusterVariableStore.Delete(txCtx, cVariable.Id); err != nil {
				return err
			}
			//插入删除记录

			if err = g.variableHistoryStore.HistoryDelete(txCtx, namespaceID, cVariable.Id, &variable_entry.VariableValue{Key: key, Value: cVariable.Value}, userID); err != nil {
				return err
			}
		}

		//删除全局环境变量
		_, err = g.globalVariableStore.Delete(txCtx, globalVariable.Id)
		return err
	})

}

func (g *globalVariableService) QuoteVariables(ctx context.Context, namespaceID int, sourceID int, quoteType quote_entry.QuoteKindType, variableKeys []string) error {
	variables, err := g.globalVariableStore.GetGlobalVariableByKeys(ctx, namespaceID, variableKeys)
	if err != nil {
		return fmt.Errorf("引用环境变量失败:%w", err)
	}
	if len(variables) != len(variableKeys) {
		variablesSet := common.SliceToMap(variables, func(v *variable_entry.Variables) string {
			return v.Key
		})
		defeatVariables := make([]string, 0, 3)
		for _, v := range variableKeys {
			if _, has := variablesSet[v]; !has {
				defeatVariables = append(defeatVariables, v)
			}
		}
		return fmt.Errorf("全局环境不存在%s", strings.Join(defeatVariables, ","))
	}

	variablesIDs := make([]int, 0, len(variableKeys))
	for _, v := range variables {
		variablesIDs = append(variablesIDs, v.Id)
	}
	return g.quoteStore.Set(ctx, sourceID, quoteType, quote_entry.QuoteTargetKindTypeVariable, variablesIDs...)
}

func (g *globalVariableService) CheckQuotedVariablesOnline(ctx context.Context, clusterID int, clusterName string, sourceID int, quoteType quote_entry.QuoteKindType) error {
	//服务引用的环境变量
	quoteMaps, err := g.quoteStore.GetSourceQuote(ctx, sourceID, quoteType)
	if err != nil {
		return fmt.Errorf("获取目标引用的环境变量失败:%w", err)
	}
	variableIds := quoteMaps[quote_entry.QuoteTargetKindTypeVariable]
	if len(variableIds) > 0 {
		//获取集群正在运行的环境变量版本
		variablePublishVersion, err := g.clusterVariableService.GetPublishVersion(ctx, clusterID)
		if err != nil {
			return fmt.Errorf("获取环境变量发布状态失败:%w", err)
		}
		if variablePublishVersion == nil {
			return errors.New("环境变量尚未发布")
		}

		//已发布的环境变量
		toMap := common.SliceToMap(variablePublishVersion.ClusterVariable, func(t *variable_entry.ClusterVariable) int {
			return t.VariableId
		})

		for _, variableId := range variableIds {
			if _, ok := toMap[variableId]; !ok {
				globalVariable, err := g.GetById(ctx, variableId)
				if err != nil {
					return err
				}
				return errors.New(fmt.Sprintf("${%s}未上线到{%s}，上线/更新失败", globalVariable.Key, clusterName))
			}
		}
	}

	return nil
}

func (g *globalVariableService) DeleteVariableQuote(ctx context.Context, sourceID int, quoteType quote_entry.QuoteKindType) error {
	return g.quoteStore.DelSourceTarget(ctx, sourceID, quoteType, quote_entry.QuoteTargetKindTypeVariable)
}
