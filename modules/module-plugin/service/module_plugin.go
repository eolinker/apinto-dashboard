package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/eolinker/apinto-dashboard/modules/base/locker-service"
	"github.com/eolinker/apinto-dashboard/modules/group"
	group_service "github.com/eolinker/apinto-dashboard/modules/group/group-service"
	module_plugin "github.com/eolinker/apinto-dashboard/modules/module-plugin"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/dto"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/model"
	"github.com/eolinker/apinto-dashboard/modules/module-plugin/store"
	"github.com/eolinker/eosc/common/bean"
)

var (
	modulePluginNotFound = errors.New("plugin doesn't exist. ")
)

type modulePluginService struct {
	pluginStore        store.IModulePluginStore
	pluginEnableStore  store.IModulePluginEnableStore
	pluginPackageStore store.IModulePluginPackageStore

	commonGroup group.ICommonGroupService
	lockService locker_service.IAsynLockService
}

func newModulePluginService() module_plugin.IModulePluginService {

	s := &modulePluginService{}
	bean.Autowired(&s.pluginStore)
	bean.Autowired(&s.pluginEnableStore)
	bean.Autowired(&s.pluginPackageStore)

	bean.Autowired(&s.commonGroup)
	bean.Autowired(&s.lockService)
	return s
}

func (m *modulePluginService) GetPlugins(ctx context.Context, groupUUID, searchName string) ([]*model.ModulePluginItem, error) {
	groupID := -1
	if groupUUID != "" {
		groupInfo, err := m.commonGroup.GetGroupInfo(ctx, groupUUID)
		if err != nil {
			return nil, err
		}
		groupID = groupInfo.Id
	}
	pluginEntries, err := m.pluginStore.GetPluginList(ctx, groupID, searchName)
	if err != nil {
		return nil, err
	}
	plugins := make([]*model.ModulePluginItem, 0, len(pluginEntries))
	for _, entry := range pluginEntries {
		plugin := &model.ModulePluginItem{
			ModulePlugin: entry,
			IsEnable:     false,
			IsInner:      true,
		}
		enableEntry, err := m.pluginEnableStore.Get(ctx, entry.Id)
		if err != nil {
			return nil, err
		}
		//若为非内置
		if entry.Type == 2 {
			plugin.IsInner = false
		}
		//若插件已启用
		if enableEntry.IsEnable == 2 {
			plugin.IsEnable = true
		}
		plugins = append(plugins, plugin)
	}
	return plugins, nil
}

func (m *modulePluginService) GetPluginInfo(ctx context.Context, pluginUUID string) (*model.ModulePluginInfo, error) {
	plugin, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}

	info := &model.ModulePluginInfo{
		ModulePlugin: plugin,
		IsEnable:     false,
		Uninstall:    true,
	}

	enableEntry, err := m.pluginEnableStore.Get(ctx, plugin.Id)
	if err != nil {
		return nil, err
	}
	//若为非内置
	if plugin.Type == 2 {
		info.Uninstall = false
	}
	//若插件已启用
	if enableEntry.IsEnable == 2 {
		info.IsEnable = true
	}
	return info, nil
}

func (m *modulePluginService) GetPluginGroups(ctx context.Context) ([]*model.PluginGroup, error) {
	groupEntries, err := m.commonGroup.GroupListAll(ctx, -1, group_service.ModulePlugin, group_service.ModulePlugin)
	if err != nil {
		return nil, err
	}
	groups := make([]*model.PluginGroup, 0, len(groupEntries))
	for _, entry := range groupEntries {
		groups = append(groups, &model.PluginGroup{
			UUID: entry.Uuid,
			Name: entry.Name,
		})
	}
	return groups, nil
}

func (m *modulePluginService) GetPluginEnableInfo(ctx context.Context, pluginUUID string) (*model.PluginEnableInfo, error) {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}
	enableEntry, err := m.pluginEnableStore.Get(ctx, pluginInfo.Id)
	if err != nil {
		return nil, err
	}

	//TODO 通过导航id获取导航信息

	enableCfg := new(model.PluginEnableCfg)

	info := &model.PluginEnableInfo{
		Name:       enableEntry.Name,
		Navigation: "",
		ApiGroup:   enableCfg.APIGroup,
		Server:     enableCfg.Server,
		Header:     enableCfg.Header,
		Query:      enableCfg.Query,
		Initialize: enableCfg.Initialize,
	}

	return info, nil
}

func (m *modulePluginService) GetPluginEnableRender(ctx context.Context, pluginUUID string) (*model.PluginEnableRender, error) {
	pluginInfo, err := m.pluginStore.GetPluginInfo(ctx, pluginUUID)
	if err != nil {
		return nil, err
	}

	renderCfg := &model.PluginEnableRender{
		Internet:  false,
		Invisible: true,
		ApiGroup:  false,
	}
	switch pluginInfo.Driver {
	case "remote":
		remoteDefine := new(model.RemoteDefine)
		_ = json.Unmarshal([]byte(pluginInfo.Details), remoteDefine)
		if !remoteDefine.Internet {
			renderCfg.Internet = true
		}
		renderCfg.Querys = remoteDefine.Querys
		renderCfg.Initialize = remoteDefine.Initialize
	case "local":
		renderCfg.ApiGroup = true
		localDefine := new(model.LocalDefine)
		_ = json.Unmarshal([]byte(pluginInfo.Details), localDefine)
		renderCfg.Headers = localDefine.Headers
		renderCfg.Querys = localDefine.Querys
		renderCfg.Initialize = localDefine.Initialize
		renderCfg.Invisible = localDefine.Invisible
	}
	return renderCfg, nil
}

func (m *modulePluginService) InstallPlugin(ctx context.Context, groupName string, packageContent []byte) error {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) EnablePlugin(ctx context.Context, pluginUUID string, enableInfo *dto.PluginEnableInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) DisablePlugin(ctx context.Context, pluginUUID string) error {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) GetEnabledPlugins(ctx context.Context) ([]*model.InstalledPlugin, error) {
	//TODO implement me
	panic("implement me")
}

func (m *modulePluginService) GetMiddlewareList(ctx context.Context) ([]*model.MiddlewareItem, error) {
	//TODO implement me
	panic("implement me")
}
