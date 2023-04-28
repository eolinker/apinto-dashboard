package plugin_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/api"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	quote_entry "github.com/eolinker/apinto-dashboard/modules/base/quote-entry"
	quote_store "github.com/eolinker/apinto-dashboard/modules/base/quote-store"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/cluster/cluster-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-store"
	"github.com/eolinker/apinto-dashboard/modules/plugin_template"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"sort"
	"strings"
	"time"
)

type pluginService struct {
	pluginStore        plugin_store.IPluginStore
	quoteStore         quote_store.IQuoteStore
	pluginHistoryStore plugin_store.IPluginHistoryStore
	extenderCache      IExtenderCache

	userInfoService       user.IUserInfoService
	apiService            api.IAPIService
	pluginTemplateService plugin_template.IPluginTemplateService
	clusterService        cluster.IClusterService
	lockService           locker_service.IAsynLockService
	apintoClient          cluster.IApintoClient
	clusterPluginService  plugin.IClusterPluginService
}

func newPluginService() plugin.IPluginService {
	n := &pluginService{}
	bean.Autowired(&n.pluginStore)
	bean.Autowired(&n.pluginHistoryStore)
	bean.Autowired(&n.quoteStore)
	bean.Autowired(&n.extenderCache)
	bean.Autowired(&n.userInfoService)
	bean.Autowired(&n.clusterService)
	bean.Autowired(&n.lockService)
	bean.Autowired(&n.apintoClient)
	bean.Autowired(&n.clusterPluginService)
	bean.Autowired(&n.pluginTemplateService)
	return n
}

// GetList 获取插件列表
func (p *pluginService) GetList(ctx context.Context, namespaceId int) ([]*plugin_model.Plugin, error) {
	plugins, err := p.pluginStore.GetListByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	userIds := common.SliceToSliceIds(plugins, func(t *plugin_entry.Plugin) int {
		return t.Operator
	})

	clusters, err := p.clusterService.GetAllCluster(ctx)
	if err != nil {
		return nil, err
	}

	userInfoMaps, _ := p.userInfoService.GetUserInfoMaps(ctx, userIds...)

	resultList := make([]*plugin_model.Plugin, 0, len(plugins))
	for _, info := range plugins {

		pluginInfo := &plugin_model.Plugin{Plugin: info}
		if userInfo, uOk := userInfoMaps[info.Operator]; uOk {
			pluginInfo.OperatorStr = userInfo.NickName
		}

		isDelete, _ := p.isDelete(ctx, namespaceId, clusters, plugins, pluginInfo.Plugin)
		pluginInfo.IsDelete = isDelete

		pluginInfo.IsBuilt = info.Type == 1
		resultList = append(resultList, pluginInfo)
	}

	return resultList, nil
}

func (p *pluginService) isDelete(ctx context.Context, namespaceId int, clusters []*cluster_model.Cluster, plugins []*plugin_entry.Plugin, plugin *plugin_entry.Plugin) (bool, error) {

	relyPlugin := make(map[int]struct{})

	pluginMaps := common.SliceToMap(plugins, func(t *plugin_entry.Plugin) int {
		return t.Id
	})
	for _, pluginInfo := range plugins {
		if v, ok := pluginMaps[pluginInfo.Rely]; ok {
			relyPlugin[v.Id] = struct{}{}
		}
	}
	//依赖的插件不可删除
	if _, ok := relyPlugin[plugin.Id]; ok {
		return false, errors.New("依赖插件不可删除")
	}

	quote, _ := p.quoteStore.GetTargetQuote(ctx, plugin.Id, quote_entry.QuoteTargetKindTypePlugin)
	for kindType, ids := range quote {
		switch kindType {
		case quote_entry.QuoteKindTypePluginTemplate:
			if len(ids) > 0 {
				pluginTemplateInfo, _ := p.pluginTemplateService.GetBasicInfoByID(ctx, ids[0])
				if pluginTemplateInfo != nil {
					return false, errors.New(fmt.Sprintf("该插件已被名为%s的模板引用，不可删除", pluginTemplateInfo.Name))
				}
			}
			return false, errors.New("未知引用，不可删除")
		case quote_entry.QuoteKindTypeAPI:
			if len(ids) > 0 {
				apiInfo, _ := p.apiService.GetAPIInfoById(ctx, ids[0])
				if apiInfo != nil {
					return false, errors.New(fmt.Sprintf("该插件已被名为%s的API引用，不可删除", apiInfo.Name))
				}
			}
			return false, errors.New("未知引用，不可删除")
		}
	}

	clusterNames := make([]string, 0)
	for _, clusterInfo := range clusters {
		//当插件在网关集群的状态为已发布且是禁用状态可删除，当插件在网关集群的状态为未发布可删除
		isDelete, err := p.clusterPluginService.IsDelete(ctx, namespaceId, clusterInfo.Name, plugin.Name)
		if err != nil {
			return false, err
		}
		if !isDelete {
			clusterNames = append(clusterNames, clusterInfo.Name)
		}
	}

	if len(clusterNames) > 0 {
		return false, errors.New(fmt.Sprintf("该插件已被%s集群上线启用或全局启用，不可删除!", strings.Join(clusterNames, ",")))
	}

	return true, nil
}

func (p *pluginService) GetBasicInfoList(ctx context.Context, namespaceId int) ([]*plugin_model.PluginBasic, error) {
	plugins, err := p.pluginStore.GetListByNamespaceId(ctx, namespaceId)
	if err != nil {
		return nil, err
	}

	resultList := make([]*plugin_model.PluginBasic, 0, len(plugins))
	for _, info := range plugins {

		pluginInfo := &plugin_model.PluginBasic{Plugin: info}
		resultList = append(resultList, pluginInfo)
	}

	return resultList, nil
}

// InsertBuilt 新增内置插件
func (p *pluginService) InsertBuilt(ctx context.Context, namespaceId int, pluginBuilt []*plugin_model.PluginBuilt) error {
	if err := p.lockService.Lock(locker_service.LockNamePluginNamespace, namespaceId); err != nil {
		return err
	}

	defer p.lockService.Unlock(locker_service.LockNamePluginNamespace, namespaceId)

	t := time.Now()
	plugins := make([]*plugin_entry.Plugin, 0, len(pluginBuilt))

	type PluginRely struct {
		plugin *plugin_entry.Plugin
		rely   string
	}

	//有依赖的插件
	pluginRelys := make([]*PluginRely, 0)

	for _, built := range pluginBuilt {
		pluginInfo := &plugin_entry.Plugin{
			NamespaceId: namespaceId,
			Name:        built.Name,
			Extended:    built.Extended,
			Schema:      built.Schema,
			Type:        1,
			Sort:        built.Sort,
			CreateTime:  t,
			UpdateTime:  t,
		}

		plugins = append(plugins, pluginInfo)
		if built.Rely != "" {
			pluginRelys = append(pluginRelys, &PluginRely{
				plugin: pluginInfo,
				rely:   built.Rely,
			})
		}
	}

	pluginMaps := common.SliceToMap(plugins, func(t *plugin_entry.Plugin) string {
		return t.Extended
	})

	return p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		if err := p.pluginStore.Insert(txCtx, plugins...); err != nil {
			return err
		}

		//处理插件依赖
		for _, rely := range pluginRelys {
			if v, ok := pluginMaps[rely.rely]; ok {
				rely.plugin.Rely = v.Id
				if err := p.pluginStore.Save(txCtx, rely.plugin); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// Create 新增插件
func (p *pluginService) Create(ctx context.Context, namespaceId, operator int, input *plugin_model.PluginInput) error {

	if err := p.lockService.Lock(locker_service.LockNamePluginNamespace, namespaceId); err != nil {
		return err
	}
	defer p.lockService.Unlock(locker_service.LockNamePluginNamespace, namespaceId)

	pluginInfo, _ := p.pluginStore.GetByName(ctx, namespaceId, input.Name)

	if pluginInfo != nil && pluginInfo.Name == input.Name {
		return errors.New("插件名称重复")
	}

	extender, err := p.extenderCache.GetAll(ctx)
	if err != nil {
		return err
	}

	extenderMap := common.SliceToMap(extender, func(t *plugin_model.ExtenderInfo) string {
		return t.Id
	})

	schema := ""
	if v, ok := extenderMap[input.Extended]; ok {
		schema = v.Schema
	} else {
		return errors.New("获取jsonSchema失败")
	}

	t := time.Now()
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: input.Name,
	})

	rely := 0
	if input.RelyName != "" {
		relyPlugin, err := p.pluginStore.GetByName(ctx, namespaceId, input.RelyName)
		if err != nil {
			return err
		}
		rely = relyPlugin.Id
	}

	//获取最大排序号
	index, err := p.pluginStore.MaxSort(ctx, namespaceId)
	if err != nil {
		return err
	}

	pluginInfo = &plugin_entry.Plugin{
		NamespaceId: namespaceId,
		Name:        input.Name,
		Extended:    input.Extended,
		Desc:        input.Desc,
		Schema:      schema,
		Type:        2,
		Rely:        rely,
		Sort:        index + 1,
		Operator:    operator,
		CreateTime:  t,
		UpdateTime:  t,
	}

	return p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = p.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}

		//添加修改记录
		if err = p.pluginHistoryStore.HistoryAdd(txCtx, namespaceId, pluginInfo.Id, &plugin_entry.PluginHistoryInfo{
			Plugin: *pluginInfo,
		}, operator); err != nil {
			return err
		}

		return nil
	})
}

// Update 修改插件
func (p *pluginService) Update(ctx context.Context, namespaceId, operator int, input *plugin_model.PluginInput) error {

	if err := p.lockService.Lock(locker_service.LockNamePluginNamespace, namespaceId); err != nil {
		return err
	}
	defer p.lockService.Unlock(locker_service.LockNamePluginNamespace, namespaceId)

	pluginInfo, err := p.pluginStore.GetByName(ctx, namespaceId, input.Name)
	if err != nil {
		return err
	}

	t := time.Now()
	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: input.Name,
	})

	oldPlugin := &plugin_entry.PluginHistoryInfo{
		Plugin: *pluginInfo,
	}
	pluginInfo.UpdateTime = t
	pluginInfo.Desc = input.Desc
	pluginInfo.Operator = operator

	newPlugin := &plugin_entry.PluginHistoryInfo{
		Plugin: *pluginInfo,
	}

	return p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = p.pluginStore.Save(txCtx, pluginInfo); err != nil {
			return err
		}

		//添加修改记录
		if err = p.pluginHistoryStore.HistoryEdit(txCtx, namespaceId, pluginInfo.Id, oldPlugin, newPlugin, operator); err != nil {
			return err
		}

		return nil
	})
}

func (p *pluginService) Delete(ctx context.Context, namespaceId, operator int, name string) error {
	pluginInfo, err := p.pluginStore.GetByName(ctx, namespaceId, name)
	if err != nil {
		return err
	}

	plugins, err := p.pluginStore.GetListByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}

	clusters, err := p.clusterService.GetAllCluster(ctx)
	if err != nil {
		return err
	}

	//判断是否可删除
	if _, err = p.isDelete(ctx, namespaceId, clusters, plugins, pluginInfo); err != nil {
		return err
	}

	controller.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name: name,
	})

	return p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		if _, err = p.pluginStore.Delete(txCtx, pluginInfo.Id); err != nil {
			return err
		}
		//添加删除记录
		if err = p.pluginHistoryStore.HistoryDelete(txCtx, namespaceId, pluginInfo.Id, &plugin_entry.PluginHistoryInfo{
			Plugin: *pluginInfo,
		}, operator); err != nil {
			return err
		}
		return nil
	})
}

func (p *pluginService) GetByName(ctx context.Context, namespaceId int, name string) (*plugin_model.Plugin, error) {
	pluginInfo, err := p.pluginStore.GetByName(ctx, namespaceId, name)
	if err != nil {
		return nil, err
	}

	result := &plugin_model.Plugin{
		Plugin: pluginInfo,
	}
	if result.Rely > 0 {
		relyPlugin, err := p.pluginStore.Get(ctx, result.Rely)
		if err != nil {
			return nil, err
		}
		result.RelyName = relyPlugin.Name
	}
	return result, nil
}

func (p *pluginService) Sort(ctx context.Context, namespaceId, _ int, names []string) error {
	if err := p.lockService.Lock(locker_service.LockNamePluginNamespace, namespaceId); err != nil {
		return err
	}
	defer p.lockService.Unlock(locker_service.LockNamePluginNamespace, namespaceId)

	plugins, err := p.pluginStore.GetListByNamespaceId(ctx, namespaceId)
	if err != nil {
		return err
	}
	//长度不对 不允许修改排序
	if len(plugins) != len(names) {
		return errors.New("数据错误，请重新刷新页面。")
	}

	pluginMaps := common.SliceToMap(plugins, func(t *plugin_entry.Plugin) string {
		return t.Name
	})

	pluginIdMaps := common.SliceToMap(plugins, func(t *plugin_entry.Plugin) int {
		return t.Id
	})

	for i, name := range names {
		if pluginInfo, ok := pluginMaps[name]; ok {
			pluginInfo.Sort = i + 1
		} else {
			return errors.New(fmt.Sprintf("找不到名称为%s的插件", name))
		}
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Sort < plugins[j].Sort
	})

	for _, pluginInfo := range plugins {
		if pluginInfo.Rely > 0 {
			if v, ok := pluginIdMaps[pluginInfo.Rely]; ok {
				if v.Sort > pluginInfo.Sort {
					return errors.New(fmt.Sprintf("%s依赖%s，不可修改顺序", pluginInfo.Name, v.Name))
				}
			} else {
				return errors.New(fmt.Sprintf("%s依赖的插件不存在，不可修改顺序", pluginInfo.Name))
			}
		}
	}

	return p.pluginStore.Transaction(ctx, func(txCtx context.Context) error {
		for _, pluginInfo := range plugins {

			if err = p.pluginStore.Save(txCtx, pluginInfo); err != nil {
				return err
			}

		}

		return nil
	})
}
