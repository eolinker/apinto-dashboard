package plugin_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/client/v1"
	global_plugin "github.com/eolinker/apinto-dashboard/client/v1/initialize/plugin"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/modules/audit/audit-model"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/apinto-dashboard/modules/plugin"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-entry"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	"github.com/eolinker/apinto-dashboard/modules/plugin/plugin-store"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
)

var (
	clusterPluginNotFound = errors.New("plugin doesn't exist. ")
)

type clusterPluginService struct {
	clusterService                   cluster.IClusterService
	clusterNodeService               cluster.IClusterNodeService
	namespaceService                 namespace.INamespaceService
	apintoClient                     cluster.IApintoClient
	pluginService                    plugin.IPluginService
	clusterPluginStore               plugin_store.IClusterPluginStore
	clusterPluginHistoryStore        plugin_store.IClusterPluginHistoryStore
	clusterPluginPublishVersionStore plugin_store.IClusterPluginPublishVersionStore
	clusterPluginRuntimeStore        plugin_store.IClusterPluginRuntimeStore
	clusterPluginPublishHistoryStore plugin_store.IClusterPluginPublishHistoryStore
	lockService                      locker_service.IAsynLockService
	userInfoService                  user.IUserInfoService
}

func newClusterPluginService() plugin.IClusterPluginService {

	s := &clusterPluginService{}
	bean.Autowired(&s.pluginService)
	bean.Autowired(&s.clusterPluginStore)
	bean.Autowired(&s.clusterService)
	bean.Autowired(&s.clusterPluginHistoryStore)
	bean.Autowired(&s.clusterPluginPublishVersionStore)
	bean.Autowired(&s.clusterPluginRuntimeStore)
	bean.Autowired(&s.clusterPluginPublishHistoryStore)
	bean.Autowired(&s.namespaceService)
	bean.Autowired(&s.clusterNodeService)
	bean.Autowired(&s.apintoClient)
	bean.Autowired(&s.lockService)
	bean.Autowired(&s.userInfoService)

	return s
}

func (c *clusterPluginService) GetList(ctx context.Context, namespaceID int, clusterName string) ([]*plugin_model.CluPluginListItem, error) {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return nil, err
	}

	list := make([]*plugin_model.CluPluginListItem, 0)

	//获取工作空间下所有插件
	globalPlugins, err := c.pluginService.GetBasicInfoList(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	//获取该集群下所有集群插件
	clusterPlugins, err := c.clusterPluginStore.GetByClusterIds(ctx, clusterInfo.Id)
	if err != nil {
		return nil, err
	}

	clusterPluginsMap := common.SliceToMap(clusterPlugins, func(t *plugin_entry.ClusterPlugin) string {
		return t.PluginName
	})

	currentVersion, err := c.GetPublishVersion(ctx, clusterInfo.Id)
	if err != nil {
		return nil, err
	}

	var releasedPluginsMap = map[string]*plugin_entry.CluPluginPublishConfig{}
	if currentVersion != nil && currentVersion.ClusterPluginPublishVersion != nil {
		//当前版本已发布的集群插件
		releasedPluginsMap = common.SliceToMap(currentVersion.PublishedPluginsList, func(t *plugin_entry.CluPluginPublishConfig) string {
			return t.ClusterPlugin.PluginName
		})
	}

	//获取用户ids
	userIds := make([]int, 0, 5)
	userSet := make(map[int]struct{}, 5)
	for _, p := range globalPlugins {
		if _, exist := userSet[p.Operator]; exist {
			continue
		}
		userSet[p.Operator] = struct{}{}
		userIds = append(userIds, p.Operator)
	}
	for _, p := range clusterPlugins {
		if _, exist := userSet[p.Operator]; exist {
			continue
		}
		userSet[p.Operator] = struct{}{}
		userIds = append(userIds, p.Operator)
	}
	userInfoMaps, _ := c.userInfoService.GetUserInfoMaps(ctx, userIds...)

	innerPlugins := getInnerPlugins()
	innerPluginsMap := common.SliceToMap(innerPlugins, func(t *plugin_model.InnerPlugin) string {
		return t.PluginName
	})
	//对比全局插件
	for _, gPlugin := range globalPlugins {
		nowSort := gPlugin.Sort
		item := &plugin_model.CluPluginListItem{}

		if current, has := clusterPluginsMap[gPlugin.Name]; has {

			item.ClusterPlugin = current
			if released, has := releasedPluginsMap[gPlugin.Name]; has {
				if current.Status == released.ClusterPlugin.Status && current.Config == released.ClusterPlugin.Config && released.Sort == nowSort {
					item.Publish = 2     //已发布
					item.ChangeState = 0 //无变更
					item.ReleasedSort = nowSort
					item.NowSort = nowSort
				} else {
					item.Publish = 1     //未发布
					item.ChangeState = 2 //修改
					item.ReleasedSort = released.Sort
					item.NowSort = nowSort
				}
			} else {
				item.Publish = 1     //未发布
				item.ChangeState = 2 //修改
				item.ReleasedSort = 0
				item.NowSort = nowSort
			}
			operatorName := ""
			if userInfo, ok := userInfoMaps[current.Operator]; ok {
				operatorName = userInfo.NickName
			}
			item.Operator = operatorName

		} else {
			item.NowSort = nowSort
			//若为内置插件
			if gPlugin.Type == 1 {
				innerPlugin, ok := innerPluginsMap[gPlugin.Name]
				if !ok {
					return nil, errors.New(fmt.Sprintf("找不到名称为%s的内置插件", gPlugin.Name))
				}
				item.ClusterPlugin = &plugin_entry.ClusterPlugin{PluginName: gPlugin.Name, Status: innerPlugin.Status, Config: innerPlugin.Config}
				item.ReleasedSort = 0
				item.IsBuiltIn = true
				//若该内置插件已发布过
				if released, has := releasedPluginsMap[gPlugin.Name]; has {
					if innerPlugin.Status == released.ClusterPlugin.Status && innerPlugin.Config == released.ClusterPlugin.Config && released.Sort == nowSort {
						item.Publish = 2     //已发布
						item.ChangeState = 0 //无变更
						item.ReleasedSort = nowSort
						item.NowSort = nowSort
					} else {
						item.Publish = 1     //未发布
						item.ChangeState = 2 //修改
						item.ReleasedSort = released.Sort
						item.NowSort = nowSort
					}
				} else {
					item.ChangeState = 1
				}
			} else {
				//当前集群插件中没有，则是刚新增的插件
				item.ClusterPlugin = &plugin_entry.ClusterPlugin{PluginName: gPlugin.Name}
				item.Status = 1
				item.ChangeState = 1 //新增
				item.Publish = 3
			}
		}
		list = append(list, item)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].NowSort < list[j].NowSort
	})

	return list, nil
}

func (c *clusterPluginService) GetPlugin(ctx context.Context, namespaceID int, clusterName string, pluginName string) (*plugin_model.ClusterPluginInfo, error) {
	//验证clusterName合法性
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return nil, common.ClusterNotExist
	}

	clusterPlugin, err := c.clusterPluginStore.GetClusterPluginByClusterIDByPluginName(ctx, clusterInfo.Id, pluginName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("%s %w", pluginName, clusterPluginNotFound)
		}
		return nil, err
	}

	return &plugin_model.ClusterPluginInfo{
		PluginName: clusterPlugin.PluginName,
		Status:     clusterPlugin.Status,
		Config:     clusterPlugin.Config,
	}, nil
}

func (c *clusterPluginService) EditPlugin(ctx context.Context, namespaceID int, clusterName string, userID int, pluginName string, status int, config interface{}) error {
	//验证clusterName合法性
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return common.ClusterNotExist
	}

	//验证全局plugin存不存在
	globalPlugin, err := c.pluginService.GetByName(ctx, namespaceID, pluginName)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("%s %w", pluginName, clusterPluginNotFound)
		}
		return err
	}
	//内置插件不能编辑
	if globalPlugin.Type == 1 {
		return errors.New("Can't Edit Inner Plugin. ")
	}

	pluginConfig, _ := json.Marshal(config)
	//检测JsonSchema格式是否正确
	if err = common.JsonSchemaValid(globalPlugin.Schema, string(pluginConfig)); err != nil {
		return errors.New(fmt.Sprintf("插件配置格式错误 err=%s", err.Error()))
	}

	if err = c.lockService.Lock(locker_service.LockNameClusterPlugin, clusterInfo.Id); err != nil {
		return err
	}
	defer c.lockService.Unlock(locker_service.LockNameClusterPlugin, clusterInfo.Id)

	clusterPlugin, err := c.clusterPluginStore.GetClusterPluginByClusterIDByPluginName(ctx, clusterInfo.Id, pluginName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var oldValue *plugin_entry.ClusterPluginHistoryValue
	t := time.Now()
	//判断变更历史是新增还是修改
	isAddHistory := false

	// 若该集群插件为空，则新建
	if clusterPlugin == nil {
		isAddHistory = true

		clusterPlugin = &plugin_entry.ClusterPlugin{
			NamespaceId: namespaceID,
			ClusterId:   clusterInfo.Id,
			PluginName:  pluginName,
			Status:      status,
			Config:      string(pluginConfig),
			Operator:    userID,
			CreateTime:  t,
			UpdateTime:  t,
		}
	} else {
		oldValue = &plugin_entry.ClusterPluginHistoryValue{
			PluginName: pluginName,
			Status:     clusterPlugin.Status,
			Config:     clusterPlugin.Config,
		}

		clusterPlugin.Status = status
		clusterPlugin.Config = string(pluginConfig)
		clusterPlugin.Operator = userID
		clusterPlugin.UpdateTime = t
	}

	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        pluginName,
		ClusterId:   clusterInfo.Id,
		ClusterName: clusterName,
	})

	return c.clusterPluginStore.Transaction(ctx, func(txCtx context.Context) error {

		if err = c.clusterPluginStore.Save(txCtx, clusterPlugin); err != nil {
			return err
		}
		if isAddHistory {
			return c.clusterPluginHistoryStore.HistoryAdd(txCtx, namespaceID, clusterPlugin.Id, &plugin_entry.ClusterPluginHistoryValue{
				PluginName: pluginName,
				Status:     status,
				Config:     string(pluginConfig),
			}, userID)
		}

		return c.clusterPluginHistoryStore.HistoryEdit(txCtx, namespaceID, clusterPlugin.Id, oldValue, &plugin_entry.ClusterPluginHistoryValue{
			PluginName: pluginName,
			Status:     status,
			Config:     string(pluginConfig),
		}, userID)
	})
}

// QueryHistory 集群插件的变更记录
func (c *clusterPluginService) QueryHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*plugin_model.ClusterPluginHistory, int, error) {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, 0, common.ClusterNotExist
	}

	plugins, err := c.clusterPluginStore.GetByClusterIds(ctx, clusterInfo.Id)
	if err != nil {
		return nil, 0, err
	}

	ids := common.SliceToSliceIds(plugins, func(t *plugin_entry.ClusterPlugin) int {
		return t.Id
	})

	histories, count, err := c.clusterPluginHistoryStore.Page(ctx, namespaceId, pageNum, pageSize, ids...)
	if err != nil {
		return nil, 0, err
	}

	list := make([]*plugin_model.ClusterPluginHistory, 0, len(histories))
	for _, history := range histories {
		list = append(list, &plugin_model.ClusterPluginHistory{ClusterPluginHistory: history})
	}

	return list, count, nil
}

func (c *clusterPluginService) PublishHistory(ctx context.Context, namespaceId, pageNum, pageSize int, clusterName string) ([]*plugin_model.ClusterPluginPublish, int, error) {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, 0, common.ClusterNotExist
	}

	list, count, err := c.clusterPluginPublishHistoryStore.GetByClusterPage(ctx, pageNum, pageSize, clusterInfo.Id)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]*plugin_model.ClusterPluginPublish, 0, len(list))

	userIds := make([]int, 0, len(list))
	userSet := make(map[int]struct{}, len(list))
	for _, h := range list {
		if _, exist := userSet[h.Operator]; exist {
			continue
		}
		userSet[h.Operator] = struct{}{}
		userIds = append(userIds, h.Operator)
	}
	infoMaps, _ := c.userInfoService.GetUserInfoMaps(ctx, userIds...)

	for _, publish := range list {
		data := publish.PluginToPublish
		details := make([]*plugin_model.ClusterPluginPublishDetails, 0, len(data))

		for _, val := range data {
			details = append(details, &plugin_model.ClusterPluginPublishDetails{
				Name: val.PluginName,
				OldValue: plugin_model.ClusterPluginHistoryConfig{
					Status: val.ReleasedConfig.Status,
					Config: val.ReleasedConfig.Config,
				},
				NewValue: plugin_model.ClusterPluginHistoryConfig{
					Status: val.NoReleasedConfig.Status,
					Config: val.NoReleasedConfig.Config,
				},
				ReleasedSort: val.ReleasedSort,
				NowSort:      val.NowSort,
				OptType:      val.OptType,
				CreateTime:   val.CreateTime,
			})
		}

		operator := ""
		if userInfo, ok := infoMaps[publish.Operator]; ok {
			operator = userInfo.NickName
		}
		resp = append(resp, &plugin_model.ClusterPluginPublish{
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

func (c *clusterPluginService) ToPublishes(ctx context.Context, namespaceId int, clusterName string) ([]*plugin_model.ClusterPluginToPublish, error) {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return nil, common.ClusterNotExist
	}

	//获取集群当前的集群插件
	currentCluPlugins, err := c.clusterPluginStore.GetByClusterIds(ctx, clusterInfo.Id)
	if err != nil {
		return nil, err
	}

	t := time.Now()

	//获取全局插件
	globalPlugins, err := c.pluginService.GetBasicInfoList(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	globalPluginsMap := common.SliceToMap(globalPlugins, func(t *plugin_model.PluginBasic) string {
		return t.Name
	})

	globalRelyPluginsMap := common.SliceToMap(globalPlugins, func(t *plugin_model.PluginBasic) int {
		return t.Id
	})
	//拼凑当前的待发布插件配置, 备注：内置插件不存到集群插件表里
	currentPlugins := make([]*plugin_entry.CluPluginPublishConfig, 0, len(currentCluPlugins))
	for _, cluPlugin := range currentCluPlugins {
		nowSort := 0
		extended := ""
		rely := ""
		if gPlugin, exist := globalPluginsMap[cluPlugin.PluginName]; exist {
			nowSort = gPlugin.Sort
			extended = gPlugin.Extended

			if v, ok := globalRelyPluginsMap[gPlugin.Rely]; ok {
				rely = v.Extended
			}
		} else {
			continue
		}
		currentPlugins = append(currentPlugins, &plugin_entry.CluPluginPublishConfig{
			ClusterPlugin: cluPlugin,
			Sort:          nowSort,
			Extended:      extended,
			Rely:          rely,
		})
	}
	//将内置插件加入到待发布插件配置
	for _, gPlugin := range getInnerPlugins() {
		currentPlugins = append(currentPlugins, &plugin_entry.CluPluginPublishConfig{
			ClusterPlugin: &plugin_entry.ClusterPlugin{
				PluginName: gPlugin.PluginName,
				Status:     gPlugin.Status,
				Config:     gPlugin.Config,
				Operator:   0,
				CreateTime: t,
				UpdateTime: t,
			},
			Sort:     globalPluginsMap[gPlugin.PluginName].Sort,
			Extended: gPlugin.Id,
			Rely:     gPlugin.Rely,
		})
	}

	//查询当前版本下的集群插件
	clusterRuntime, err := c.clusterPluginRuntimeStore.GetForCluster(ctx, clusterInfo.Id, clusterInfo.Id)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//表示当前集群还没发布任何环境遍历版本
	respList := make([]*plugin_model.ClusterPluginToPublish, 0)
	if clusterRuntime == nil {
		for _, pluginInfo := range currentPlugins {

			entryPlugin := &plugin_entry.PluginToPublish{
				PluginName:     pluginInfo.ClusterPlugin.PluginName,
				Extended:       pluginInfo.Extended,
				Rely:           pluginInfo.Rely,
				ReleasedConfig: plugin_entry.PluginToPublishConfig{},
				NoReleasedConfig: plugin_entry.PluginToPublishConfig{
					Status: pluginInfo.ClusterPlugin.Status,
					Config: pluginInfo.ClusterPlugin.Config,
				},
				NowSort:    pluginInfo.Sort,
				CreateTime: pluginInfo.ClusterPlugin.CreateTime,
				OptType:    1,
			}
			respList = append(respList, &plugin_model.ClusterPluginToPublish{
				PluginToPublish: entryPlugin,
			})
		}
		sort.Slice(respList, func(i, j int) bool {
			return respList[i].NowSort < respList[j].NowSort
		})
		return respList, nil
	}

	version, err := c.clusterPluginPublishVersionStore.Get(ctx, clusterRuntime.VersionId)
	if err != nil {
		return nil, err
	}

	releasedPlugins := version.PublishedPluginsList

	//差异化对比
	releasedPluginsMap := common.SliceToMap(releasedPlugins, func(t *plugin_entry.CluPluginPublishConfig) string {
		return t.ClusterPlugin.PluginName
	})

	addList, updateList, delList := common.DiffContrast(releasedPlugins, currentPlugins)

	for _, pluginInfo := range addList {
		entryPlugin := &plugin_entry.PluginToPublish{
			Rely:             pluginInfo.Rely,
			Extended:         pluginInfo.Extended,
			PluginName:       pluginInfo.ClusterPlugin.PluginName,
			ReleasedConfig:   plugin_entry.PluginToPublishConfig{},
			NoReleasedConfig: plugin_entry.PluginToPublishConfig{Status: pluginInfo.ClusterPlugin.Status, Config: pluginInfo.ClusterPlugin.Config},
			CreateTime:       pluginInfo.ClusterPlugin.CreateTime,
			ReleasedSort:     0,
			NowSort:          pluginInfo.Sort,
			OptType:          1,
		}
		respList = append(respList, &plugin_model.ClusterPluginToPublish{PluginToPublish: entryPlugin})
	}
	for _, pluginInfo := range updateList {
		var releasedValue plugin_entry.PluginToPublishConfig
		releasedSort := 0
		if p, ok := releasedPluginsMap[pluginInfo.ClusterPlugin.PluginName]; ok {
			releasedValue = plugin_entry.PluginToPublishConfig{Status: p.ClusterPlugin.Status, Config: p.ClusterPlugin.Config}
			releasedSort = p.Sort
		}

		entryPlugin := &plugin_entry.PluginToPublish{
			Extended:         pluginInfo.Extended,
			Rely:             pluginInfo.Rely,
			PluginName:       pluginInfo.ClusterPlugin.PluginName,
			ReleasedConfig:   releasedValue,
			NoReleasedConfig: plugin_entry.PluginToPublishConfig{Status: pluginInfo.ClusterPlugin.Status, Config: pluginInfo.ClusterPlugin.Config},
			ReleasedSort:     releasedSort,
			NowSort:          pluginInfo.Sort,
			CreateTime:       pluginInfo.ClusterPlugin.CreateTime,
			OptType:          2,
		}
		respList = append(respList, &plugin_model.ClusterPluginToPublish{PluginToPublish: entryPlugin})
	}

	for _, pluginInfo := range delList {
		var releasedValue plugin_entry.PluginToPublishConfig
		releasedSort := 0
		if p, ok := releasedPluginsMap[pluginInfo.ClusterPlugin.PluginName]; ok {
			releasedValue = plugin_entry.PluginToPublishConfig{Status: p.ClusterPlugin.Status, Config: p.ClusterPlugin.Config}
			releasedSort = p.Sort
		}

		entryPlugin := &plugin_entry.PluginToPublish{
			PluginName:       pluginInfo.ClusterPlugin.PluginName,
			Extended:         pluginInfo.Extended,
			Rely:             pluginInfo.Rely,
			ReleasedConfig:   releasedValue,
			NoReleasedConfig: plugin_entry.PluginToPublishConfig{},
			ReleasedSort:     releasedSort,
			NowSort:          0,
			CreateTime:       pluginInfo.ClusterPlugin.CreateTime,
			OptType:          3,
		}

		respList = append(respList, &plugin_model.ClusterPluginToPublish{PluginToPublish: entryPlugin})
	}

	sort.Slice(respList, func(i, j int) bool {
		return respList[i].NowSort < respList[j].NowSort
	})

	return respList, nil
}

func (c *clusterPluginService) Publish(ctx context.Context, namespaceId, userId int, clusterName, versionName, desc, source string) error {
	t := time.Now()

	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceId, clusterName)
	if err != nil {
		return err
	}

	clusterId := clusterInfo.Id

	if err = c.lockService.Lock(locker_service.LockNameClusterPlugin, clusterId); err != nil {
		return err
	}
	defer c.lockService.Unlock(locker_service.LockNameClusterPlugin, clusterId)

	//查询版本名称是否重复
	publishHistory, err := c.clusterPluginPublishHistoryStore.GetByVersionName(ctx, versionName, clusterId)
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

	publishes := make([]*plugin_entry.PluginToPublish, 0)
	if err = json.Unmarshal(bytes, &publishes); err != nil {
		return err
	}
	if len(publishes) == 0 {
		return errors.New("没有变化无需发布")
	}

	//获取全局插件
	globalPlugins, err := c.pluginService.GetBasicInfoList(ctx, namespaceId)
	if err != nil {
		return err
	}
	globalPluginsMap := common.SliceToMap(globalPlugins, func(t *plugin_model.PluginBasic) string {
		return t.Name
	})
	//校验要发布的配置的JsonSchema格式是否正确
	for _, p := range publishes {
		//如果是删除操作，则不校验
		if p.OptType == 3 {
			continue
		}
		gp := globalPluginsMap[p.PluginName]
		if err = common.JsonSchemaValid(gp.Schema, p.NoReleasedConfig.Config); err != nil {
			return errors.New(fmt.Sprintf("插件%s配置格式错误 err=%s", p.PluginName, err.Error()))
		}
	}

	//获取集群当前运行的版本
	currentVersion, err := c.GetPublishVersion(ctx, clusterId)
	if err != nil {
		return err
	}

	newClusterPublishPlugins := make([]*plugin_entry.CluPluginPublishConfig, 0)
	if currentVersion != nil && currentVersion.ClusterPluginPublishVersion != nil && currentVersion.Id > 0 {
		currentVersionPluginsMaps := common.SliceToMap(currentVersion.PublishedPluginsList, func(t *plugin_entry.CluPluginPublishConfig) string {
			return t.ClusterPlugin.PluginName
		})

		for _, publish := range publishes {
			if publish.OptType == 1 { //新增 直接追加
				currentVersionPluginsMaps[publish.PluginName] = &plugin_entry.CluPluginPublishConfig{
					ClusterPlugin: &plugin_entry.ClusterPlugin{
						NamespaceId: namespaceId,
						ClusterId:   clusterId,
						PluginName:  publish.PluginName,
						Status:      publish.NoReleasedConfig.Status,
						Config:      publish.NoReleasedConfig.Config,
						CreateTime:  publish.CreateTime,
					},
					Sort:     publish.NowSort,
					Extended: publish.Extended,
				}

			} else if publish.OptType == 2 { //修改 找到旧版本的plugin  然后修改配置
				if pluginInfo, ok := currentVersionPluginsMaps[publish.PluginName]; ok {
					pluginInfo.ClusterPlugin.Status = publish.NoReleasedConfig.Status
					pluginInfo.ClusterPlugin.Config = publish.NoReleasedConfig.Config
					pluginInfo.Sort = publish.NowSort
					pluginInfo.Extended = publish.Extended
				}
			} else if publish.OptType == 3 {
				delete(currentVersionPluginsMaps, publish.PluginName)
			}
		}
		for _, pluginInfo := range currentVersionPluginsMaps {
			newClusterPublishPlugins = append(newClusterPublishPlugins, pluginInfo)
		}

	} else {
		//当前没有旧版本 表示这是第一次发布  直接用source内的未发布插件
		tmpClusterPlugins := make([]*plugin_entry.CluPluginPublishConfig, 0)
		for _, publish := range publishes {
			tmpClusterPlugins = append(tmpClusterPlugins, &plugin_entry.CluPluginPublishConfig{
				ClusterPlugin: &plugin_entry.ClusterPlugin{
					NamespaceId: namespaceId,
					ClusterId:   clusterId,
					PluginName:  publish.PluginName,
					Status:      publish.NoReleasedConfig.Status,
					Config:      publish.NoReleasedConfig.Config,
					CreateTime:  publish.CreateTime,
				},
				Sort:     publish.NowSort,
				Extended: publish.Extended,
			})

		}
		newClusterPublishPlugins = tmpClusterPlugins
	}

	newVersion := &plugin_entry.ClusterPluginPublishVersion{
		ClusterId:            clusterId,
		NamespaceId:          namespaceId,
		Desc:                 desc,
		PublishedPluginsList: newClusterPublishPlugins,
		Operator:             userId,
		CreateTime:           t,
	}

	names := make([]string, 0)
	for _, pluginInfo := range newClusterPublishPlugins {
		names = append(names, pluginInfo.ClusterPlugin.PluginName)
	}

	common.SetGinContextAuditObject(ctx, &audit_model.LogObjectInfo{
		Name:        strings.Join(names, ","),
		ClusterId:   clusterId,
		ClusterName: clusterName,
	})

	return c.clusterPluginPublishVersionStore.Transaction(ctx, func(txCtx context.Context) error {
		if err = c.clusterPluginPublishVersionStore.Save(txCtx, newVersion); err != nil {
			return err
		}
		//当前集群运行的版本
		clusterPluginRuntime := &plugin_entry.ClusterPluginRuntime{
			VersionId:   newVersion.Id,
			ClusterId:   clusterId,
			NamespaceId: namespaceId,
			Operator:    userId,
			IsOnline:    true,
			CreateTime:  t,
			UpdateTime:  t,
		}

		history := &plugin_entry.ClusterPluginPublishHistory{
			VersionName: versionName,
			ClusterId:   clusterId,
			NamespaceId: namespaceId,
			Desc:        desc,
			VersionId:   newVersion.Id,
			CluPluginPublishHistoryInfo: plugin_entry.CluPluginPublishHistoryInfo{
				PluginToPublish: publishes,
			},
			OptType:  1,
			Operator: userId,
			OptTime:  t,
		}
		if err = c.clusterPluginPublishHistoryStore.Insert(txCtx, history); err != nil {
			return err
		}

		if err = c.clusterPluginRuntimeStore.Save(txCtx, clusterPluginRuntime); err != nil {
			return err
		}
		client, err := c.apintoClient.GetClient(ctx, clusterId)
		if err != nil {
			return err
		}
		plugins, err := getToSendPlugins(newClusterPublishPlugins)
		if err != nil {
			return err
		}
		//发布插件
		return client.ForGlobalPlugin().Set(plugins)
	})
}

// DeleteAll 删除某集群下所有集群插件
func (c *clusterPluginService) DeleteAll(ctx context.Context, clusterId int) error {
	clusterPlugins, err := c.clusterPluginStore.GetByClusterIds(ctx, clusterId)
	if err != nil {
		return err
	}

	pluginIDs := common.SliceToSliceIds(clusterPlugins, func(t *plugin_entry.ClusterPlugin) int {
		return t.Id
	})

	return c.clusterPluginStore.DeleteClusterPluginByIDs(ctx, pluginIDs...)
}

func (c *clusterPluginService) GetPublishVersion(ctx context.Context, clusterId int) (*plugin_model.ClusterPluginPublishVersion, error) {
	//获取集群当前运行的版本
	currentRuntime, err := c.clusterPluginRuntimeStore.GetForCluster(ctx, clusterId, clusterId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var currentVersion *plugin_entry.ClusterPluginPublishVersion
	if currentRuntime != nil {
		//获取当前集群运行版本的详细信息
		currentVersion, err = c.clusterPluginPublishVersionStore.Get(ctx, currentRuntime.VersionId)
		if err != nil {
			return nil, err
		}
	}

	return &plugin_model.ClusterPluginPublishVersion{ClusterPluginPublishVersion: currentVersion}, nil
}

func (c *clusterPluginService) ResetOnline(ctx context.Context, _, clusterId int) {
	runtime, err := c.clusterPluginRuntimeStore.GetForCluster(ctx, clusterId, clusterId)
	if err != nil {
		log.Errorf("clusterPluginService-ResetOnline-GetRuntime clusterId=%d,err=%s", clusterId, err.Error())
		return
	}

	if runtime.IsOnline {
		version, err := c.clusterPluginPublishVersionStore.Get(ctx, runtime.VersionId)
		if err != nil {
			return
		}

		client, err := c.apintoClient.GetClient(ctx, clusterId)
		if err != nil {
			return
		}

		plugins, err := getToSendPlugins(version.PublishedPluginsList)
		if err != nil {
			return
		}
		//发布插件
		_ = client.ForGlobalPlugin().Set(plugins)
	}
}

func (c *clusterPluginService) IsOnlineByName(ctx context.Context, namespaceID int, clusterName, pluginName string) (bool, error) {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return false, err
	}
	currentVersion, err := c.GetPublishVersion(ctx, clusterInfo.Id)
	if err != nil {
		return false, err
	}
	//若该集群没有发布过集群插件
	if currentVersion.ClusterPluginPublishVersion == nil {
		return false, nil
	}

	for _, pluginInfo := range currentVersion.PublishedPluginsList {
		if pluginInfo.ClusterPlugin.PluginName == pluginName {
			return true, nil
		}
	}
	return false, nil
}

// IsDelete  当插件在网关集群的状态为已发布且是禁用状态可删除，当插件在网关集群的状态为未发布可删除
func (c *clusterPluginService) IsDelete(ctx context.Context, namespaceID int, clusterName, pluginName string) (bool, error) {
	clusterInfo, err := c.clusterService.GetByNamespaceByName(ctx, namespaceID, clusterName)
	if err != nil {
		return true, err
	}
	currentVersion, err := c.GetPublishVersion(ctx, clusterInfo.Id)
	if err != nil {
		return true, err
	}
	//若该集群没有发布过集群插件
	if currentVersion.ClusterPluginPublishVersion == nil {
		return true, nil
	}
	releasedPlugins := currentVersion.ClusterPluginPublishVersion.PublishedPluginsList
	for _, pluginInfo := range releasedPlugins {
		if pluginInfo.ClusterPlugin.PluginName == pluginName && (pluginInfo.ClusterPlugin.Status == 2 || pluginInfo.ClusterPlugin.Status == 3) {
			return false, nil
		}
	}
	return true, nil
}

// getToSendPlugins 拼好往apinto发送的全局插件列表 TODO
func getToSendPlugins(toPublish []*plugin_entry.CluPluginPublishConfig) ([]*v1.GlobalPlugin, error) {
	//toPublish内包含了内置插件
	sort.Slice(toPublish, func(i, j int) bool {
		return toPublish[i].Sort < toPublish[j].Sort
	})

	globalPlugins := global_plugin.GetGlobalPluginConf()
	globalPluginsMap := common.SliceToMap(globalPlugins, func(t *v1.GlobalPlugin) string {
		return t.Name
	})

	resultList := make([]*v1.GlobalPlugin, 0)

	for _, publish := range toPublish {
		status := ""
		switch publish.ClusterPlugin.Status {
		case 1:
			status = "disable"
		case 2:
			status = "enable"
		case 3:
			status = "global"
		}

		p := &v1.GlobalPlugin{
			Config: publish.ClusterPlugin.Config,
			Id:     publish.Extended,
			Name:   publish.ClusterPlugin.PluginName,
			Status: status,
			Rely:   publish.Rely,
		}
		if v, ok := globalPluginsMap[publish.ClusterPlugin.PluginName]; ok {
			p.Config = v.Config
		} else {
			var config interface{}
			if err := json.Unmarshal([]byte(publish.ClusterPlugin.Config), &config); err != nil {
				return nil, err
			}
			p.Config = config
		}
		resultList = append(resultList, p)
	}

	return resultList, nil
}

func getInnerPlugins() []*plugin_model.InnerPlugin {
	globalPlugins := global_plugin.GetGlobalPluginConf()
	innerPlugins := make([]*plugin_model.InnerPlugin, 0, len(globalPlugins))

	for _, pluginInfo := range globalPlugins {
		var conf []byte
		if pluginInfo.Config != nil {
			conf, _ = json.Marshal(pluginInfo.Config)
		}
		status := 0
		switch strings.ToUpper(pluginInfo.Status) {
		case enum.PluginStateNameDisable:
			status = 1
		case enum.PluginStateNameEnable:
			status = 2
		case enum.PluginStateNameGlobal:
			status = 3
		}
		innerPlugins = append(innerPlugins, &plugin_model.InnerPlugin{
			Id:         pluginInfo.Id,
			PluginName: pluginInfo.Name,
			Status:     status,
			Config:     string(conf),
			Rely:       pluginInfo.Rely,
		})
	}
	return innerPlugins
}
