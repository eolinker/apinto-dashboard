package plugin_timer

import (
	"context"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/cache"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/client/v1/initialize/plugin"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/modules/cluster"
	"github.com/eolinker/apinto-dashboard/modules/namespace"
	plugin2 "github.com/eolinker/apinto-dashboard/modules/plugin"
	plugin_model "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-model"
	plugin_service "github.com/eolinker/apinto-dashboard/modules/plugin/plugin-service"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"time"
)

var _ IExtender = (*extender)(nil)

type IExtender interface {
	UpdateExtender()
}

type extender struct {
	pluginService    plugin2.IPluginService
	apintoClient     cluster.IApintoClient
	commonCache      cache.ICommonCache
	extenderCache    plugin_service.IExtenderCache
	clusterService   cluster.IClusterService
	namespaceService namespace.INamespaceService
}

func newExtender() IExtender {
	ex := &extender{}
	bean.Autowired(&ex.pluginService)
	bean.Autowired(&ex.apintoClient)
	bean.Autowired(&ex.commonCache)
	bean.Autowired(&ex.clusterService)
	bean.Autowired(&ex.extenderCache)
	bean.Autowired(&ex.namespaceService)

	return ex
}

func (e *extender) UpdateExtender() {
	ctx := context.TODO()
	//key+时间戳
	lockKey := fmt.Sprintf("updateExtender_%d", time.Now().Unix())
	if err := e.lock(ctx, lockKey); err != nil {
		log.Warnf("lockKey=%s error=%s", lockKey, err.Error())
		return
	}

	extenderMaps, err := e.getExtenderMaps(ctx)
	if err != nil {
		return
	}

	namespaces, err := e.namespaceService.GetAll()
	if err != nil {
		log.Errorf("extender-UpdateExtender namespaceService.GetAll error:%s", err.Error())
		return
	}
	//内置插件
	globalPlugins := plugin.GetGlobalPluginConf()

	for _, namespaceInfo := range namespaces {

		if err = e.builtPlugin(ctx, namespaceInfo.Id, globalPlugins, extenderMaps); err != nil {
			log.Errorf("extender-UpdateExtender err=%s", err.Error())
		}

	}

}

func (e *extender) builtPlugin(ctx context.Context, namespaceId int, globalPlugins []*v1.GlobalPlugin, extenderMaps map[string]*plugin_model.ExtenderInfo) error {
	plugins, err := e.pluginService.GetBasicInfoList(ctx, namespaceId)
	if err != nil {
		log.Errorf("获取插件信息失败 err=%s", err.Error())
		return err
	}

	pluginBuilt := make([]*plugin_model.PluginBuilt, 0)
	if len(plugins) > 0 {
		pluginMaps := common.SliceToMap(plugins, func(t *plugin_model.PluginBasic) string {
			return t.Extended
		})
		//如果有已经插件,去重后把新增的内置插件添加到最后
		sort := plugins[len(plugins)-1].Sort
		for _, globalPlugin := range globalPlugins {
			if _, ok := pluginMaps[globalPlugin.Id]; !ok {
				pluginBuilt = append(pluginBuilt, &plugin_model.PluginBuilt{
					Extended: globalPlugin.Id,
					Name:     globalPlugin.Name,
					Rely:     globalPlugin.Rely,
					Sort:     sort + 1,
				})
				sort++
			}
		}

	} else {

		for i, globalPlugin := range globalPlugins {
			pluginBuilt = append(pluginBuilt, &plugin_model.PluginBuilt{
				Extended: globalPlugin.Id,
				Name:     globalPlugin.Name,
				Rely:     globalPlugin.Rely,
				Sort:     i + 1,
			})
		}

	}

	for _, built := range pluginBuilt {
		if extenderInfo, extOk := extenderMaps[built.Extended]; extOk {
			built.Schema = extenderInfo.Schema
		} else {
			log.Errorf("拿不到ID为%s插件的 JSON Schema", built.Extended)
		}
	}
	if len(pluginBuilt) > 0 {
		if err = e.pluginService.InsertBuilt(ctx, namespaceId, pluginBuilt); err != nil {
			log.Errorf("新增内置插件失败 err=%s", err.Error())
			return err
		}
	}
	return nil
}

func (e *extender) getExtenderMaps(ctx context.Context) (map[string]*plugin_model.ExtenderInfo, error) {
	clusters, err := e.clusterService.GetAllCluster(ctx)
	if err != nil {
		log.Errorf("获取集群信息失败 err=%s", err.Error())
		return nil, err
	}

	var extenders = make(map[string]*v1.ExtenderListItem)
	var clients []v1.IClient
	for _, clusterInfo := range clusters {
		client, err := e.apintoClient.GetClient(ctx, clusterInfo.Id)
		if err != nil {
			log.Errorf("获取client失败 err=%s", err.Error())
			continue
		}
		clients = append(clients, client)
		tempExtenders, err := client.ForExtender().List()
		if err != nil {
			log.Errorf("获取扩展ID列表失败 err=%s", err.Error())
			continue
		}
		for _, tempExtender := range tempExtenders {
			extenders[tempExtender.Id] = tempExtender
		}
	}

	extenderInfos := make([]*plugin_model.ExtenderInfo, 0)

	for _, item := range extenders {
		for _, client := range clients {
			extenderInfo, err := client.ForExtender().Info(item.Group, item.Project, item.Name)
			if err != nil {
				log.Errorf("获取插件扩展信息失败 err=%s", err.Error())
				continue
			}
			extenderInfos = append(extenderInfos, &plugin_model.ExtenderInfo{
				Id:      item.Id,
				Group:   item.Group,
				Project: item.Project,
				Name:    item.Name,
				Version: item.Version,
				Schema:  string(extenderInfo.Render),
			})
			break
		}
	}

	log.DebugF("当前最新扩展ID列表 extenders=%v", extenders)
	if len(extenderInfos) > 0 {
		if err = e.extenderCache.SetAll(ctx, e.extenderCache.Key(), extenderInfos, time.Minute*5); err != nil {
			log.Errorf("扩展ID插入缓存失败 err=%s", err.Error())
			return nil, err
		}
	}

	extenderMaps := common.SliceToMap(extenderInfos, func(t *plugin_model.ExtenderInfo) string {
		return t.Id
	})

	return extenderMaps, nil
}

func (e *extender) lock(ctx context.Context, key string) error {
	log.DebugF("UpdateExtender-lock key=%s", key)
	b, err := e.commonCache.SetNX(ctx, key, "1", time.Minute)
	if err != nil {
		return err
	}
	if b {
		return nil
	}
	return errors.New("锁已被占用")
}
