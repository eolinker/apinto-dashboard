package plugin_template_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/frontend-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	"github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/plugin"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template/plugin-template-store"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/go-basic/uuid"
	"gorm.io/gorm"
	"sort"
	"time"
)

type pluginTemplateService struct {
	pluginTemplateStore        plugin_template_store.IPluginTemplateStore
	pluginTemplateStatStore    plugin_template_store.IPluginTemplateStatStore
	pluginTemplateRuntimeStore plugin_template_store.IPluginTemplateRuntimeStore
	pluginTemplateHistoryStore plugin_template_store.IPluginTemplateHistoryStore
	pluginTemplateVersionStore plugin_template_store.IPluginTemplateVersionStore
	quoteStore                 quote_store.IQuoteStore

	userInfoService      user.IUserInfoService
	clusterService       cluster.IClusterService
	lockService          locker_service.IAsynLockService
	apintoClient         cluster.IApintoClient
	pluginService        plugin.IPluginService
	apiService           api.IAPIService
	clusterPluginService plugin.IClusterPluginService
	namespaceService     namespace.INamespaceService
}

func newPluginTemplateService() plugin_template.IPluginTemplateService {
	n := &pluginTemplateService{}
	bean.Autowired(&n.pluginTemplateStore)
	bean.Autowired(&n.pluginTemplateStatStore)
	bean.Autowired(&n.pluginTemplateHistoryStore)
	bean.Autowired(&n.pluginTemplateRuntimeStore)
	bean.Autowired(&n.pluginTemplateVersionStore)
	bean.Autowired(&n.quoteStore)

	bean.Autowired(&n.userInfoService)
	bean.Autowired(&n.clusterService)
	bean.Autowired(&n.lockService)
	bean.Autowired(&n.pluginService)
	bean.Autowired(&n.apintoClient)
	bean.Autowired(&n.apiService)
	bean.Autowired(&n.clusterPluginService)
	bean.Autowired(&n.namespaceService)
	return n
}

// GetList 获取插件模板列表
func (p *pluginTemplateService) GetList(ctx context.Context, namespaceId int) ([]*plugin_template_model.PluginTemplate, error) {
	pluginTemplates, err := p.pluginTemplateStore.GetListByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	clusters, err := p.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(pluginTemplates, func(t *plugin_template_entry.PluginTemplate) int {
		return t.Operator
	})

	userInfoMaps, _ := p.userInfoService.GetUserInfoMaps(ctx, userIds...)

	resultList := make([]*plugin_template_model.PluginTemplate, 0, len(pluginTemplates))
	for _, template := range pluginTemplates {

		pluginTemplate := &plugin_template_model.PluginTemplate{PluginTemplate: template}
		if userInfo, uOk := userInfoMaps[template.Operator]; uOk {
			pluginTemplate.OperatorStr = userInfo.NickName
		}

		isDelete, _ := p.isDelete(ctx, clusters, template)

		pluginTemplate.IsDelete = isDelete
		resultList = append(resultList, pluginTemplate)
	}

	return resultList, nil
}

func (p *pluginTemplateService) isDelete(ctx context.Context, clusters []*cluster_model.Cluster, pluginTemplate *plugin_template_entry.PluginTemplate) (bool, error) {
	for _, clusterInfo := range clusters {
		isOnline, _ := p.IsOnline(ctx, clusterInfo.Id, pluginTemplate.UUID)
		if isOnline {
			return false, errors.New(fmt.Sprintf("插件模板已在%s集群上线,不可删除", clusterInfo.Name))
		}
	}

	quote, err := p.quoteStore.GetTargetQuote(ctx, pluginTemplate.Id, quote_entry.QuoteTargetKindTypePluginTemplate)
	if err != nil {
		return false, err
	}
	apiIds := quote[quote_entry.QuoteKindTypeAPI]
	if len(apiIds) > 0 {
		return false, errors.New("插件模板被API绑定,不可以删除")
	}

	return true, err
}

func (p *pluginTemplateService) IsOnline(ctx context.Context, clusterId int, uuid string) (bool, error) {

	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return false, err
	}

	//获取当前运行的版本
	runtime, err := p.pluginTemplateRuntimeStore.GetForCluster(ctx, pluginTemplate.Id, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	if runtime != nil && runtime.IsOnline {
		return true, nil
	}
	return false, nil
}

// GetUsableList 获取模板可用列表
func (p *pluginTemplateService) GetUsableList(ctx context.Context, namespaceId int) ([]*plugin_template_model.PluginTemplate, error) {
	pluginTemplates, err := p.pluginTemplateStore.GetListByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	resultList := make([]*plugin_template_model.PluginTemplate, 0, len(pluginTemplates))
	for _, template := range pluginTemplates {

		resultList = append(resultList, &plugin_template_model.PluginTemplate{PluginTemplate: template})
	}

	return resultList, nil
}

// Create 新增插件模板
func (p *pluginTemplateService) Create(ctx context.Context, namespaceId, operator int, input *plugin_template_model.PluginTemplateDetail) error {

	if input.UUID == "" {
		input.UUID = uuid.New()
	}

	plugins := make([]*plugin_model.Plugin, 0, len(input.Plugins))
	for _, pluginInfo := range input.Plugins {

		modelPlugin, err := p.pluginService.GetByName(ctx, namespaceId, pluginInfo.Name)
		if err == gorm.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("插件名为%s的插件不存在", pluginInfo.Name))
		} else if err != nil {
			return err
		}

		bytes, _ := json.Marshal(pluginInfo.Config)
		//检测JsonSchema格式是否正确
		if !common.JsonSchemaValid(modelPlugin.Schema, string(bytes)) {
			return errors.New(fmt.Sprintf("插件名为%s的config配置格式错误", modelPlugin.Name))
		}

		plugins = append(plugins, modelPlugin)
	}

	t := time.Now()
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: input.Name,
	})

	pluginTemplate := &plugin_template_entry.PluginTemplate{
		NamespaceId: namespaceId,
		UUID:        input.UUID,
		Name:        input.Name,
		Desc:        input.Desc,
		Operator:    operator,
		CreateTime:  t,
		UpdateTime:  t,
	}

	return p.pluginTemplateStore.Transaction(ctx, func(txCtx context.Context) error {
		if err := p.pluginTemplateStore.Save(txCtx, pluginTemplate); err != nil {
			return err
		}

		newPlugins := make([]*plugin_template_entry.PluginTemplateVersionConfigDetail, 0, len(input.Plugins))
		for _, pluginInfo := range input.Plugins {
			bytes, _ := json.Marshal(pluginInfo.Config)
			newPlugins = append(newPlugins, &plugin_template_entry.PluginTemplateVersionConfigDetail{
				Name:    pluginInfo.Name,
				Config:  string(bytes),
				Disable: pluginInfo.Disable,
			})
		}

		pluginTemplateVersionConfig := plugin_template_entry.PluginTemplateVersionConfig{}
		pluginTemplateVersionConfig.Plugins = newPlugins
		version := &plugin_template_entry.PluginTemplateVersion{
			PluginTemplateId:            pluginTemplate.Id,
			NamespaceID:                 namespaceId,
			PluginTemplateVersionConfig: pluginTemplateVersionConfig,
			Operator:                    operator,
			CreateTime:                  t,
		}

		//添加版本信息
		if err := p.pluginTemplateVersionStore.Save(txCtx, version); err != nil {
			return err
		}

		//添加修改记录
		if err := p.pluginTemplateHistoryStore.HistoryAdd(txCtx, namespaceId, pluginTemplate.Id, &plugin_template_entry.PluginTemplateHistoryInfo{
			PluginTemplate: *pluginTemplate,
			Config:         pluginTemplateVersionConfig,
		}, operator); err != nil {
			return err
		}

		//添加版本绑定关系
		stat := &plugin_template_entry.PluginTemplateStat{
			PluginTemplateId: pluginTemplate.Id,
			VersionID:        version.Id,
		}

		//添加版本关联原表信息
		if err := p.pluginTemplateStatStore.Save(txCtx, stat); err != nil {
			return err
		}

		ids := common.SliceToSliceIds(plugins, func(t *plugin_model.Plugin) int {
			return t.Id
		})

		return p.quoteStore.Set(txCtx, pluginTemplate.Id, quote_entry.QuoteKindTypePluginTemplate, quote_entry.QuoteTargetKindTypePlugin, ids...)
	})
}

// Update 修改插件模板
func (p *pluginTemplateService) Update(ctx context.Context, namespaceId, operator int, input *plugin_template_model.PluginTemplateDetail) error {

	plugins := make([]*plugin_model.Plugin, 0, len(input.Plugins))
	for _, pluginInfo := range input.Plugins {
		modelPlugin, err := p.pluginService.GetByName(ctx, namespaceId, pluginInfo.Name)
		if err == gorm.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("插件名为%s的插件不存在", pluginInfo.Name))
		} else if err != nil {
			return err
		}
		plugins = append(plugins, modelPlugin)

		bytes, _ := json.Marshal(pluginInfo.Config)
		//检测JsonSchema格式是否正确
		if !common.JsonSchemaValid(modelPlugin.Schema, string(bytes)) {
			return errors.New(fmt.Sprintf("插件名为%s的config配置格式错误", modelPlugin.Name))
		}

	}

	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, input.UUID)
	if err != nil {
		return err
	}

	t := time.Now()
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: input.UUID,
		Name: pluginTemplate.Name,
	})

	//获取当前版本信息
	currentVersion, err := p.getPluginTemplateVersion(ctx, pluginTemplate.Id)
	if err != nil {
		return err
	}

	oldPluginTemplate := &plugin_template_entry.PluginTemplateHistoryInfo{
		PluginTemplate: *pluginTemplate,
		Config:         currentVersion.PluginTemplateVersionConfig,
	}

	pluginTemplate.UpdateTime = t
	pluginTemplate.Desc = input.Desc
	pluginTemplate.Operator = operator

	return p.pluginTemplateStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = p.pluginTemplateStore.Save(txCtx, pluginTemplate); err != nil {
			return err
		}

		newPlugins := make([]*plugin_template_entry.PluginTemplateVersionConfigDetail, 0, len(input.Plugins))
		for _, pluginInfo := range input.Plugins {
			bytes, _ := json.Marshal(pluginInfo.Config)

			newPlugins = append(newPlugins, &plugin_template_entry.PluginTemplateVersionConfigDetail{
				Name:    pluginInfo.Name,
				Config:  string(bytes),
				Disable: pluginInfo.Disable,
			})
		}

		pluginTemplateVersionConfig := plugin_template_entry.PluginTemplateVersionConfig{}
		pluginTemplateVersionConfig.Plugins = newPlugins
		newVersion := &plugin_template_entry.PluginTemplateVersion{
			PluginTemplateId:            pluginTemplate.Id,
			NamespaceID:                 namespaceId,
			PluginTemplateVersionConfig: pluginTemplateVersionConfig,
			Operator:                    operator,
			CreateTime:                  t,
		}

		//添加版本信息
		if err = p.pluginTemplateVersionStore.Save(txCtx, newVersion); err != nil {
			return err
		}

		newPluginTemplate := &plugin_template_entry.PluginTemplateHistoryInfo{
			PluginTemplate: *pluginTemplate,
			Config:         pluginTemplateVersionConfig,
		}

		//添加修改记录
		if err = p.pluginTemplateHistoryStore.HistoryEdit(txCtx, namespaceId, pluginTemplate.Id, oldPluginTemplate, newPluginTemplate, operator); err != nil {
			return err
		}

		//添加版本绑定关系
		stat := &plugin_template_entry.PluginTemplateStat{
			PluginTemplateId: pluginTemplate.Id,
			VersionID:        newVersion.Id,
		}

		//添加版本关联原表信息
		if err = p.pluginTemplateStatStore.Save(txCtx, stat); err != nil {
			return err
		}

		ids := common.SliceToSliceIds(plugins, func(t *plugin_model.Plugin) int {
			return t.Id
		})

		return p.quoteStore.Set(txCtx, pluginTemplate.Id, quote_entry.QuoteKindTypePluginTemplate, quote_entry.QuoteTargetKindTypePlugin, ids...)
	})
}

func (p *pluginTemplateService) Delete(ctx context.Context, namespaceId, operator int, uuid string) error {
	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	clusters, err := p.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}

	_, err = p.isDelete(ctx, clusters, pluginTemplate)
	if err != nil {
		return err
	}

	version, err := p.getPluginTemplateVersion(ctx, pluginTemplate.Id)
	if err != nil {
		return err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid: uuid,
		Name: pluginTemplate.Name,
	})

	return p.pluginTemplateStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = p.pluginTemplateStore.Delete(txCtx, pluginTemplate.Id); err != nil {
			return err
		}
		//删除版本关联表
		{
			delMap := map[string]interface{}{"`kind`": "plugin_template", "`target`": pluginTemplate.Id}

			if _, err = p.pluginTemplateStatStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}

			if _, err = p.pluginTemplateVersionStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}

			if _, err = p.pluginTemplateRuntimeStore.DeleteWhere(txCtx, delMap); err != nil {
				return err
			}
		}

		//添加删除记录
		if err = p.pluginTemplateHistoryStore.HistoryDelete(txCtx, namespaceId, pluginTemplate.Id, &plugin_template_entry.PluginTemplateHistoryInfo{
			PluginTemplate: *pluginTemplate,
			Config:         version.PluginTemplateVersionConfig,
		}, operator); err != nil {
			return err
		}
		return p.quoteStore.DelBySource(txCtx, pluginTemplate.Id, quote_entry.QuoteKindTypePluginTemplate)
	})
}

func (p *pluginTemplateService) GetByUUID(ctx context.Context, _ int, uuid string) (*plugin_template_model.PluginTemplateDetail, error) {
	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	version, err := p.getPluginTemplateVersion(ctx, pluginTemplate.Id)
	if err != nil {
		return nil, err
	}

	plugins := make([]*plugin_template_model.PluginInfo, 0, len(version.Plugins))

	for _, pluginInfo := range version.Plugins {
		plugins = append(plugins, &plugin_template_model.PluginInfo{
			Name:    pluginInfo.Name,
			Config:  pluginInfo.Config,
			Disable: pluginInfo.Disable,
		})
	}

	detail := &plugin_template_model.PluginTemplateDetail{
		PluginTemplate: &plugin_template_model.PluginTemplate{
			PluginTemplate: pluginTemplate,
		},
		Plugins: plugins,
	}

	return detail, nil
}

func (p *pluginTemplateService) GetBasicInfoByUUID(ctx context.Context, uuid string) (*plugin_template_model.PluginTemplateBasicInfo, error) {
	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return &plugin_template_model.PluginTemplateBasicInfo{
		PluginTemplate: pluginTemplate,
	}, nil
}

func (p *pluginTemplateService) getPluginTemplateVersion(ctx context.Context, pluginTemplateId int) (*plugin_template_model.PluginTemplateVersion, error) {
	var err error

	stat, err := p.pluginTemplateStatStore.Get(ctx, pluginTemplateId)
	if err != nil {
		return nil, err
	}

	var version *plugin_template_entry.PluginTemplateVersion

	version, err = p.pluginTemplateVersionStore.Get(ctx, stat.VersionID)
	if err != nil {
		return nil, err
	}

	return (*plugin_template_model.PluginTemplateVersion)(version), nil
}

func (p *pluginTemplateService) OnlineList(ctx context.Context, namespaceId int, uuid string) ([]*plugin_template_model.PluginTemplateOnlineItem, error) {
	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	//获取工作空间下的所有集群
	clusters, err := p.clusterService.GetByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	clusterMaps := common.SliceToMap(clusters, func(t *cluster_model.Cluster) int {
		return t.Id
	})

	//获取当前服务发现下集群运行的版本
	runtimes, err := p.pluginTemplateRuntimeStore.GetByTarget(ctx, pluginTemplate.Id)
	if err != nil {
		return nil, err
	}
	runtimeMaps := common.SliceToMap(runtimes, func(t *plugin_template_entry.PluginTemplateRuntime) int {
		return t.ClusterID
	})

	//获取操作人用户列表
	operatorList := common.SliceToSliceIds(runtimes, func(t *plugin_template_entry.PluginTemplateRuntime) int {
		return t.Operator
	})

	userInfoMaps, err := p.userInfoService.GetUserInfoMaps(ctx, operatorList...)
	if err != nil {
		return nil, err
	}

	list := make([]*plugin_template_model.PluginTemplateOnlineItem, 0, len(clusters))

	latestVersion, err := p.getPluginTemplateVersion(ctx, pluginTemplate.Id)
	if err != nil {
		return nil, err
	}

	for _, clusterInfo := range clusterMaps {
		onlineItem := &plugin_template_model.PluginTemplateOnlineItem{
			ClusterName: clusterInfo.Name,
			ClusterEnv:  clusterInfo.Env,
			Status:      1, //默认为未上线状态
		}
		if runtime, ok := runtimeMaps[clusterInfo.Id]; ok {

			operator := ""
			if userInfo, uOk := userInfoMaps[runtime.Operator]; uOk {
				operator = userInfo.NickName
			}

			onlineItem.Operator = operator
			onlineItem.Disable = runtime.Disable
			onlineItem.UpdateTime = runtime.UpdateTime
			if runtime.IsOnline {
				onlineItem.Status = 3 //已上线
			} else {
				onlineItem.Status = 2 //已下线
			}
			//已上线需要对比是否更新过 服务发现信息
			if onlineItem.Status == 3 && runtime.VersionID != latestVersion.Id {
				onlineItem.Status = 4 //待更新
			}
		}

		list = append(list, onlineItem)
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Status > list[j].Status
	})
	return list, nil
}

func (p *pluginTemplateService) ResetOnline(ctx context.Context, _, clusterId int) {
	runtimes, err := p.pluginTemplateRuntimeStore.GetByCluster(ctx, clusterId)
	if err != nil {
		log.Errorf("pluginTemplateService-ResetOnline-getRuntimes clusterId=%d err=%s", clusterId, err.Error())
		return
	}
	client, err := p.apintoClient.GetClient(ctx, clusterId)
	if err != nil {
		log.Errorf("pluginTemplateService-ResetOnline-getClient clusterId=%d err=%s", clusterId, err.Error())
		return
	}

	for _, runtime := range runtimes {
		if !runtime.IsOnline {
			continue
		}

		pluginTemplate, err := p.pluginTemplateStore.Get(ctx, runtime.PluginTemplateID)
		if err != nil {
			log.Errorf("pluginTemplateService-ResetOnline-getPluginTemplate clusterId=%d pluginTemplateId=%d err=%s", clusterId, runtime.PluginTemplateID, err.Error())
			continue
		}

		version, err := p.pluginTemplateVersionStore.Get(ctx, runtime.VersionID)
		if err != nil {
			log.Errorf("pluginTemplateService-ResetOnline-getVersion clusterId=%d versionId=%d err=%s", clusterId, runtime.VersionID, err.Error())
			continue
		}

		pluginMaps := make(map[string]*v1.Plugin, 0)
		for _, plugin := range version.Plugins {
			var config interface{}
			if err = json.Unmarshal([]byte(plugin.Config), &config); err != nil {
				log.Errorf("pluginTemplateService-ResetOnline-json.Unmarshal err=%s", err.Error())
				return
			}
			pluginMaps[plugin.Name] = &v1.Plugin{
				Disable: plugin.Disable,
				Config:  config,
			}
		}
		pluginTemplateConfig := v1.PluginTemplateConfig{
			Plugins:     pluginMaps,
			Name:        pluginTemplate.UUID,
			Driver:      "plugin_template",
			Description: pluginTemplate.Desc,
		}

		if err = client.ForPluginTemplate().Create(pluginTemplateConfig); err != nil {
			log.Errorf("pluginTemplateService-ResetOnline-apinto clusterId=%d pluginTemplateConfig=%v err=%s", clusterId, pluginTemplateConfig, err.Error())
			continue
		}
	}
}
func (p *pluginTemplateService) Online(ctx context.Context, namespaceId, operator int, uuid, clusterName string) (*frontend_model.Router, error) {
	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	if err = p.lockService.Lock(locker_service.LockNamePluginTemplate, pluginTemplate.Id); err != nil {
		return nil, err
	}
	defer p.lockService.Unlock(locker_service.LockNamePluginTemplate, pluginTemplate.Id)

	pluginTemplateId := pluginTemplate.Id
	t := time.Now()

	//获取当前集群信息
	clusterInfo, err := p.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, err
	}

	latestVersion, err := p.getPluginTemplateVersion(ctx, pluginTemplateId)
	if err != nil {
		return nil, err
	}

	router := &frontend_model.Router{
		Name:   frontend_model.RouterNameServiceOnline,
		Params: make(map[string]string),
	}
	router.Params["plugin_template_name"] = pluginTemplate.Name

	for _, pluginInfo := range latestVersion.Plugins {
		isOnline, err := p.clusterPluginService.IsOnlineByName(ctx, namespaceId, clusterName, pluginInfo.Name)
		if err != nil {
			return nil, err
		}
		if !isOnline {
			return router, errors.New(fmt.Sprintf("绑定的插件%s未上线到%s集群", pluginInfo.Name, clusterName))
		}

	}

	//获取当前运行的版本
	runtime, err := p.pluginTemplateRuntimeStore.GetForCluster(ctx, pluginTemplateId, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//发布到apinto
	client, err := p.apintoClient.GetClient(ctx, clusterInfo.Id)
	if err != nil {
		return nil, err
	}

	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        uuid,
		Name:        pluginTemplate.Name,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
		PublishType: 1,
	})

	pluginMaps := make(map[string]*v1.Plugin, 0)
	for _, pluginInfo := range latestVersion.Plugins {
		var config interface{}
		if err = json.Unmarshal([]byte(pluginInfo.Config), &config); err != nil {
			return nil, err
		}
		pluginMaps[pluginInfo.Name] = &v1.Plugin{
			Disable: pluginInfo.Disable,
			Config:  config,
		}
	}
	pluginTemplateConfig := v1.PluginTemplateConfig{
		Plugins:     pluginMaps,
		Name:        pluginTemplate.UUID,
		Driver:      "plugin_template",
		Description: pluginTemplate.Desc,
	}

	//事务
	err = p.pluginTemplateStore.Transaction(ctx, func(txCtx context.Context) error {

		if runtime == nil {
			runtime = &plugin_template_entry.PluginTemplateRuntime{
				NamespaceId:      namespaceId,
				PluginTemplateID: pluginTemplateId,
				ClusterID:        clusterInfo.Id,
				VersionID:        latestVersion.Id,
				IsOnline:         true,
				Disable:          false,
				Operator:         operator,
				CreateTime:       t,
				UpdateTime:       t,
			}

			if err = p.pluginTemplateRuntimeStore.Insert(txCtx, runtime); err != nil {
				return err
			}
			return client.ForPluginTemplate().Create(pluginTemplateConfig)
		} else {
			//保存旧状态
			isOnline := runtime.IsOnline

			runtime.IsOnline = true
			runtime.UpdateTime = t
			runtime.VersionID = latestVersion.Id
			runtime.Operator = operator

			if err = p.pluginTemplateRuntimeStore.Save(txCtx, runtime); err != nil {
				return err
			}

			//若原先是下线状态
			if !isOnline {
				return client.ForPluginTemplate().Create(pluginTemplateConfig)
			}

			return client.ForPluginTemplate().Update(pluginTemplate.UUID, pluginTemplateConfig)
		}
	})

	return nil, err
}

func (p *pluginTemplateService) Offline(ctx context.Context, namespaceId, operator int, uuid, clusterName string) error {
	pluginTemplate, err := p.pluginTemplateStore.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	if err = p.lockService.Lock(locker_service.LockNamePluginTemplate, pluginTemplate.Id); err != nil {
		return err
	}
	defer p.lockService.Unlock(locker_service.LockNamePluginTemplate, pluginTemplate.Id)

	//获取当前集群信息
	clusterInfo, err := p.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	//获取当前的版本
	runtime, err := p.pluginTemplateRuntimeStore.GetForCluster(ctx, pluginTemplate.Id, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if runtime == nil {
		return errors.New("invalid version")
	}

	quote, err := p.quoteStore.GetTargetQuote(ctx, pluginTemplate.Id, quote_entry.QuoteTargetKindTypePluginTemplate)
	if err != nil {
		return err
	}
	apiIds := quote[quote_entry.QuoteKindTypeAPI]

	for _, apiId := range apiIds {
		//判断API有没有下线
		if p.apiService.IsAPIOnline(ctx, clusterInfo.Id, apiId) {
			apiInfo, err := p.apiService.GetAPIInfoById(ctx, apiId)
			if err != nil {
				return err
			}
			return errors.New(fmt.Sprintf("名称为%s的API引用了该插件模板，不可下线", apiInfo.Name))
		}
	}

	t := time.Now()
	//编写日志操作对象信息
	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Uuid:        uuid,
		Name:        pluginTemplate.Name,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
		PublishType: 2,
	})

	//事务
	return p.pluginTemplateStore.Transaction(ctx, func(txCtx context.Context) error {
		if !runtime.IsOnline {
			return errors.New("已下线不可重复下线")
		}
		runtime.IsOnline = false
		runtime.UpdateTime = t
		runtime.Operator = operator
		err = p.pluginTemplateRuntimeStore.Save(txCtx, runtime)
		if err != nil {
			return err
		}

		//发布到apinto
		client, err := p.apintoClient.GetClient(ctx, clusterInfo.Id)
		if err != nil {
			return err
		}

		return common.CheckWorkerNotExist(client.ForPluginTemplate().Delete(pluginTemplate.UUID))
	})

}
